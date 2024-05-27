package logmodule

import (
	"github.com/sirupsen/logrus"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type LogService struct {
	store store.Store
}

func NewService(store store.Store) *LogService {
	return &LogService{
		store: store,
	}
}

func (s *LogService) getLogListPage(ctx *model.SocketContext, requestLog *dto.PageReq) (*dto.LogPage, error) {
	data, count, err := s.store.Log().GetLogList(ctx.GetUser(), requestLog.Page, requestLog.PageSize, requestLog.Order, requestLog.Descending, requestLog.Filter)
	if err != nil {
		logrus.Info("error get logs - ", err)
		return nil, servererrors.ErrorInternal
	}
	if data == nil {
		data = make([]*dto.LogReq, 0)
	}

	return &dto.LogPage{
		Data:           data,
		CountAllRecord: count,
	}, nil
}
