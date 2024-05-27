package sqlstore

import (
	"database/sql"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
)

type LangRepository struct {
	store *Store
}

func (lr *LangRepository) GetLangList() ([]*model.Lang, error) {
	langRows, err := lr.store.db.Query(
		"SELECT code, name FROM public.lang_get_list()",
	)
	defer langRows.Close()
	if err != nil {
		return nil, err
	}
	var languages []*model.Lang
	for langRows.Next() {
		l := &model.Lang{}
		if err := langRows.Scan(
			&l.Code,
			&l.Name,
		); err != nil {
			return nil, err
		}
		languages = append(languages, l)
	}
	return languages, nil
}

func (lr *LangRepository) GetLang(language *model.Lang) (*model.Lang, error) {
	if err := lr.store.db.QueryRow(
		"SELECT public.lang_get($1)",
		language.Code,
	).Scan(
		&language.Meta,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, servererrors.ErrorRecordNotFound
		}
		return nil, err
	}

	return language, nil
}
