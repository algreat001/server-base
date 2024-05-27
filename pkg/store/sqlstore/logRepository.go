package sqlstore

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
)

type LogRepository struct {
	store *Store
}

func (lr *LogRepository) getCountAllLog(filter string) (int, error) {
	var count int
	if err := lr.store.db.QueryRow(
		"SELECT public.log_get_number_all_record($1)",
		filter,
	).Scan(
		&count,
	); err != nil {
		return 0, err
	}
	return count, nil
}

func getOrder(order string) string {
	switch order {
	case "operation":
		return "l.operation"
	case "createAt":
		return "l.create_at"
	case "email":
		return "u.email"
	}
	return "l.operation"
}

func (lr *LogRepository) GetLogList(executor *model.User, numPage int, sizePage int, order string, descending bool, filter string) ([]*dto.LogReq, int, error) {
	if sizePage == 0 {
		sizePage = 10000
	}

	count, err := lr.getCountAllLog(filter)
	if err != nil {
		return nil, 0, err
	}

	logRows, err := lr.store.db.Query(
		"SELECT id, operation, meta, create_at, user_id, user_email FROM public.log_get_page($1,$2,$3,$4,$5,$6)",
		executor.Id,
		getOrder(order),
		descending,
		(numPage-1)*sizePage,
		sizePage,
		filter,
	)
	defer logRows.Close()
	if err != nil {
		logrus.Info(err)
		return nil, count, err
	}
	var logs []*dto.LogReq
	for logRows.Next() {
		l := &dto.LogReq{}
		var id *uuid.UUID
		var email *string
		l.Executor = &dto.UserReq{}
		if err := logRows.Scan(
			&l.Id,
			&l.Operation,
			&l.Meta,
			&l.CreatedAt,
			&id,
			&email,
		); err != nil {
			return nil, count, err
		}
		if id != nil {
			l.Executor.Id = id
			l.Executor.Email = email
		}
		logs = append(logs, l)
	}

	return logs, count, nil
}
