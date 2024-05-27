package sqlstore

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/helpers"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
)

type UserRepository struct {
	store *Store
}

func (ur *UserRepository) AddUser(user *model.User) error {
	_, err := ur.store.db.Exec(
		"SELECT public.user_add($1, $2)",
		user.Id,
		user.Email,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) Update(user *model.User) error {
	return ur.AddUser(user)
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := ur.store.db.QueryRow(
		"SELECT id, email FROM public.user_get_by_email($1,$2)",
		nil,
		email,
	).Scan(
		&u.Id,
		&u.Email,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, servererrors.ErrorRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (ur *UserRepository) Find(id *uuid.UUID) (*model.User, error) {
	u := &model.User{}
	if err := ur.store.db.QueryRow(
		"SELECT id, email FROM public.user_get_by_id($1,$2)",
		nil,
		id,
	).Scan(
		&u.Id,
		&u.Email,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, servererrors.ErrorRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) Remove(executor *model.User, user *model.User) error {
	if err := ur.store.db.QueryRow(
		"SELECT public.user_delete_by_id($1,$2)",
		executor.Id,
		user.Id,
	).Err(); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) HardRemove(executor *model.User, user *model.User) error {
	if err := ur.store.db.QueryRow(
		"SELECT public.user_delete_by_id($1,$2)",
		executor.Id,
		user.Id,
	).Err(); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetNumberAllUsers() (int, error) {
	result := 0
	err := ur.store.db.QueryRow(
		"SELECT public.user_get_number_all_users($1)",
		nil,
	).Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, err

}

func (ur *UserRepository) GetUsers(numPage int, sizePage int, order string, descending bool, filter string) ([]*model.User, error) {
	if sizePage == 0 {
		sizePage = 10000
	}

	userRows, err := ur.store.db.Query("SELECT id, email FROM public.user_get_page($1,$2,$3,$4,$5,$6)",
		nil,
		helpers.GetSafeString(order),
		descending,
		(numPage-1)*sizePage,
		sizePage,
		helpers.GetSafeString(filter),
	)
	defer userRows.Close()
	var users []*model.User
	if err != nil {
		logrus.Info(err)
		return users, err
	}
	for userRows.Next() {
		u := &model.User{}
		if err := userRows.Scan(
			&u.Id,
			&u.Email,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil

}
