package guidemodule

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type GuideService struct {
	store    store.Store
	executor *model.User
}

func NewService(store store.Store) *GuideService {
	return &GuideService{
		store: store,
	}
}

func (s *GuideService) loadGuide(_ *model.SocketContext, getGuideReq *dto.GetGuideReq) (*dto.GetGuideRes, error) {
	return model.ReadGuideFromFile(getGuideReq.Lang, getGuideReq.GuideName)
}

func (s *GuideService) loadArticle(_ *model.SocketContext, getArticleReq *dto.GetArticleReq) (*dto.GetArticleRes, error) {
	return model.GetArticlePath(getArticleReq.Lang, getArticleReq.GuideName, getArticleReq.ArticleName)
}

func (s *GuideService) firstArticleName(_ *model.SocketContext, getGuideReq *dto.GetFirstArticleNameReq) (*dto.GetFirstArticleNameRes, error) {
	return model.GetFirstArticleName(getGuideReq.Lang, getGuideReq.GuideName)
}
