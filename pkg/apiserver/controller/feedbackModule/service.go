package feedbackmodule

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
	"net/url"
)

type FeedbackService struct {
	store store.Store
}

func NewService(store store.Store) *FeedbackService {
	return &FeedbackService{
		store: store,
	}
}

func (s *FeedbackService) send(sendReq *dto.SendReq) error {
	messageText := url.PathEscape(sendReq.Message)
	return model.Send2Telegram(messageText)
}
