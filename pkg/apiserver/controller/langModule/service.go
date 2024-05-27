package langmodule

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type LangService struct {
	store store.Store
}

func NewService(store store.Store) *LangService {
	return &LangService{
		store: store,
	}
}

func (s *LangService) getLangList() ([]*model.Lang, error) {
	return s.store.Lang().GetLangList()
}

func (s *LangService) getLang(langReq *dto.LangReq) (*model.Lang, error) {
	language := &model.Lang{
		Code: langReq.Code,
	}

	_, err := s.store.Lang().GetLang(language)

	return language, err
}
