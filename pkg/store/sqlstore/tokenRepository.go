package sqlstore

import (
	"github.com/sirupsen/logrus"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
)

type TokenRepository struct {
	store *Store
}

func (tr *TokenRepository) VerifyToken(token string) error {
	var ok bool
	if err := tr.store.db.QueryRow(
		"SELECT public.token_verify($1)",
		token,
	).Scan(
		&ok,
	); err != nil {
		logrus.Info("error verify token: ", token)
		return err
	}
	if !ok {
		logrus.Info("Token is not valid: ", token)
		return servererrors.ErrorAccessDenied
	}
	return nil
}

func (tr *TokenRepository) Create(token *model.Token) error {
	_, err := tr.store.db.Exec(
		"SELECT public.token_create($1)",
		token.Token,
	)
	if err != nil {
		logrus.Info("error create token: ", err)
		return err
	}
	return nil
}

func (tr *TokenRepository) Delete(token string) error {
	_, err := tr.store.db.Exec(
		"SELECT public.token_delete($1)",
		token,
	)
	if err != nil {
		logrus.Info("error delete token: ", err)
		return err
	}
	return nil
}
