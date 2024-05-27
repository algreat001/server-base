package logmodule

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type LogController struct {
	store   store.Store
	server  model.SocketServer
	service *LogService
}

func InitController(socketService model.SocketServer, store store.Store, _ *middleware.Auth) *LogController {
	lc := &LogController{
		store:   store,
		server:  socketService,
		service: NewService(store),
	}

	socketService.On("api/v1/log/list", lc.getLogListPage)

	return lc
}

func (lc *LogController) getLogListPage(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.PageReq, *dto.LogPage]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                dto.PageReqFromJson,
		HandlerWithRequestWithResponse: lc.service.getLogListPage,
		SocketServer:                   lc.server,
	})

}
