package sqlstore

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
)

type AccessControlEntityRepository struct {
	store *Store
}

func (r *AccessControlEntityRepository) GetGroups() ([]*model.AceGroup, error) {
	groups := []*model.AceGroup{}
	rows, err := r.store.db.Query(
		"SELECT id, name FROM public.access_group_get_all($1)",
		nil,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*model.AceGroup{}, servererrors.ErrorRecordNotFound
		}
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		addGroup := &model.AceGroup{}
		err := rows.Scan(&addGroup.Id, &addGroup.Name)
		if err != nil {
			return groups, err
		}
		groups = append(groups, addGroup)
	}
	return groups, nil
}

func (r *AccessControlEntityRepository) GetGroup(groupId *uuid.UUID) (*model.AceGroup, error) {
	group := &model.AceGroup{}
	if err := r.store.db.QueryRow(
		"SELECT id, name FROM public.access_group_get_by_id($1,$2)",
		nil,
		groupId,
	).Scan(
		&group.Id,
		&group.Name,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, servererrors.ErrorRecordNotFound
		}
		return nil, err
	}
	return group, nil
}

func (r *AccessControlEntityRepository) SetGroup(executor *model.User, newGroup *model.AceGroup) (*model.AceGroup, error) {
	if err := r.store.db.QueryRow(
		"SELECT public.access_group_create($1,$2)",
		executor.Id,
		newGroup.Name,
	).Scan(&newGroup.Id); err != nil {
		return newGroup, err
	}
	return newGroup, nil
}
func (r *AccessControlEntityRepository) RemoveGroup(executor *model.User, group *model.AceGroup) error {
	if _, err := r.store.db.Exec(
		"SELECT public.access_group_delete($1,$2)",
		executor.Id,
		group.Id,
	); err != nil {
		return err
	}
	return nil
}
func (r *AccessControlEntityRepository) UpdateGroup(executor *model.User, group *model.AceGroup) error {
	if _, err := r.store.db.Exec(
		"SELECT public.access_group_update($1,$2,$3)",
		executor.Id,
		group.Id,
		group.Name,
	); err != nil {
		return err
	}
	return nil
}

func groupsContains(groups []*model.AceGroup, group *model.AceGroup) bool {
	for _, value := range groups {
		if value.Id == group.Id {
			return true
		}
	}
	return false
}

func (r *AccessControlEntityRepository) UpdateUserGroups(executor *model.User, user *model.User) error {
	groups, err := r.GetUserGroups(user)
	if err != nil {
		return err
	}
	// remove user from groups
	for _, group := range groups {
		if !groupsContains(user.Groups, group) {
			logrus.Info("Remove group [", user.Email, "]: ", group.Name, ",executor:", executor.Id)
			if err := r.RemoveUserFromGroup(executor, user, group); err != nil {
				return err
			}
		}
	}
	// add user to groups
	for _, group := range user.Groups {
		if !groupsContains(groups, group) {
			logrus.Info("Add group: [", user.Email, "]: ", group.Name, ",executor:", executor.Id)
			if err := r.SetUserGroup(executor, user, group); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *AccessControlEntityRepository) GetUserGroups(user *model.User) ([]*model.AceGroup, error) {
	groups := []*model.AceGroup{}
	rows, err := r.store.db.Query(
		"SELECT id, name FROM public.access_group_get_for_user($1,$2)",
		nil,
		user.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*model.AceGroup{}, servererrors.ErrorRecordNotFound
		}
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		addGroup := &model.AceGroup{}
		err := rows.Scan(&addGroup.Id, &addGroup.Name)
		if err != nil {
			return groups, err
		}
		groups = append(groups, addGroup)
	}
	return groups, nil
}

func (r *AccessControlEntityRepository) SetUserGroup(executor *model.User, user *model.User, group *model.AceGroup) error {
	result := false
	if err := r.store.db.QueryRow(
		"SELECT public.access_group_add_user($1,$2,$3)",
		executor.Id,
		user.Id,
		group.Id,
	).Scan(&result); err != nil {
		return err
	}
	return nil
}
func (r *AccessControlEntityRepository) RemoveUserFromGroup(executor *model.User, user *model.User, group *model.AceGroup) error {
	result := false // false - если пользователя не было в группе, true - если был и удален
	if err := r.store.db.QueryRow(
		"SELECT public.access_group_remove_user($1,$2,$3)",
		executor.Id,
		user.Id,
		group.Id,
	).Scan(&result); err != nil {
		return err
	}
	return nil
}

func (r *AccessControlEntityRepository) GetGroupAces(group *model.AceGroup) ([]*model.Ace, error) {
	paths := []*model.Ace{}
	rows, err := r.store.db.Query(
		"SELECT id, path FROM public.ace_get_for_access_group($1,$2)",
		nil,
		group.Id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return paths, servererrors.ErrorRecordNotFound
		}
		return paths, err
	}
	defer rows.Close()
	for rows.Next() {
		addAce := &model.Ace{}
		err := rows.Scan(&addAce.Id, &addAce.Path)
		if err != nil {
			return paths, err
		}
		paths = append(paths, addAce)
	}
	return paths, nil
}

func (r *AccessControlEntityRepository) GetACL() ([]*model.Ace, error) {
	aces := []*model.Ace{}
	rows, err := r.store.db.Query(
		"SELECT * FROM public.ace_get_all($1)",
		nil,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return aces, servererrors.ErrorRecordNotFound
		}
		return aces, err
	}
	defer rows.Close()
	for rows.Next() {
		addAce := &model.Ace{Group: &model.AceGroup{}}
		err := rows.Scan(&addAce.Id, &addAce.Path, &addAce.Group.Id, &addAce.Group.Name)
		if err != nil {
			return aces, err
		}
		aces = append(aces, addAce)
	}
	return aces, nil
}

func (r *AccessControlEntityRepository) SetGroupAce(executor *model.User, group *model.AceGroup, ace *model.Ace) error {
	result := false // false - если access_group уже в ace, true - если был добавлен
	if err := r.store.db.QueryRow(
		"SELECT public.ace_add_access_group($1,$2,$3)",
		nil,
		group.Id,
		ace.Path,
	).Scan(&result); err != nil {
		if err == sql.ErrNoRows {
			return servererrors.ErrorRecordNotFound
		}
		return err
	}
	return nil
}
func (r *AccessControlEntityRepository) RemoveGroupAce(executor *model.User, group *model.AceGroup, ace *model.Ace) error {
	result := false // false - если пользователя не было в группе, true - если был и удален
	if err := r.store.db.QueryRow(
		"SELECT public.ace_remove_access_group($1,$2,$3)",
		executor.Id,
		group.Id,
		ace.Path,
	).Scan(&result); err != nil {
		return err
	}
	return nil
}

func (r *AccessControlEntityRepository) RemoveAce(executor *model.User, ace *model.Ace) error {
	adminUUID, _ := uuid.Parse("9fd8bbfa-5e15-493d-9033-7a68c5858484")
	if *ace.Id == adminUUID {
		logrus.Info("Trying to remove admin ace")
		return servererrors.ErrorACERemoveAdmin
	}
	result := false
	if err := r.store.db.QueryRow(
		"SELECT public.ace_remove_by_id($1,$2)",
		executor.Id,
		ace.Id,
	).Scan(&result); err != nil {
		return err
	}
	return nil
}

func (r *AccessControlEntityRepository) UpdateAce(executor *model.User, ace *model.Ace) error {
	result := false
	if err := r.store.db.QueryRow(
		"SELECT public.ace_update($1,$2,$3,$4)",
		executor.Id,
		ace.Id,
		ace.Group.Id,
		ace.Path,
	).Scan(&result); err != nil || !result {
		return err
	}
	return nil
}
func (r *AccessControlEntityRepository) SetAce(executor *model.User, ace *model.Ace) (*model.Ace, error) {
	if err := r.store.db.QueryRow(
		"SELECT public.ace_create($1,$2,$3)",
		executor.Id,
		ace.Group.Id,
		ace.Path,
	).Scan(
		&ace.Id,
	); err != nil {
		return nil, err
	}
	return ace, nil
}

func (r *AccessControlEntityRepository) GetUserRightForPath(user *model.User, path string) (bool, error) {
	accessGranted := true
	if err := r.store.db.QueryRow(
		"SELECT public.ace_user_right_for_path($1, $2, $3)",
		path,
		user.Id,
		user.Email,
	).Scan(&accessGranted); err != nil {
		return false, err
	}
	return accessGranted, nil
}
