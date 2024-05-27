package acemodule

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type AceController struct {
	store   store.Store
	server  model.SocketServer
	service *AceService
}

func InitController(socketService model.SocketServer, store store.Store, _ *middleware.Auth) *AceController {
	ac := &AceController{
		store:   store,
		server:  socketService,
		service: NewService(store),
	}

	socketService.On("api/v1/ace/group/list", ac.getACEGroups)
	socketService.On("api/v1/ace/group/add", ac.addACEGroup)
	socketService.On("api/v1/ace/group/remove", ac.removeACEGroup)
	socketService.On("api/v1/ace/group/update", ac.updateACEGroup)

	// CRUD пользователей в группах безопасности
	socketService.On("api/v1/ace/group/user/get", ac.getUserACEGroups)
	socketService.On("api/v1/ace/group/user/add", ac.addUserToACEGroup)
	socketService.On("api/v1/ace/group/user/remove", ac.removeUserFromACEGroup)

	// CRUD объектов безопасности
	socketService.On("api/v1/ace/list", ac.getACL)
	socketService.On("api/v1/ace/get", ac.getACEsForGroup)
	socketService.On("api/v1/ace/add", ac.addACEToGroup)
	socketService.On("api/v1/ace/update", ac.updateACE)
	socketService.On("api/v1/ace/remove", ac.removeACE)

	return ac
}

// CRUD групп безопасности
func (ac *AceController) getACEGroups(event string, ctx *model.SocketContext, _ []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[any, []*model.AceGroup]{
		Event:                             event,
		SocketContext:                     ctx,
		HandlerWithoutRequestWithResponse: ac.service.getACEGroups,
		SocketServer:                      ac.server,
	})
}

func (ac *AceController) addACEGroup(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.ACEGroupReq, any]{
		Event:                             event,
		SocketContext:                     ctx,
		Args:                              args,
		FromJsonHandler:                   dto.ACEGroupReqFromJson,
		HandlerWithRequestWithoutResponse: ac.service.addACEGroup,
		SocketServer:                      ac.server,
		OkMessage:                         "group added",
	})
}
func (ac *AceController) updateACEGroup(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.ACEGroupReq, any]{
		Event:                             event,
		SocketContext:                     ctx,
		Args:                              args,
		FromJsonHandler:                   dto.ACEGroupReqFromJson,
		HandlerWithRequestWithoutResponse: ac.service.updateACEGroup,
		SocketServer:                      ac.server,
		OkMessage:                         "group updated",
	})
}
func (ac *AceController) removeACEGroup(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.ACEGroupReq, any]{
		Event:                             event,
		SocketContext:                     ctx,
		Args:                              args,
		FromJsonHandler:                   dto.ACEGroupReqFromJson,
		HandlerWithRequestWithoutResponse: ac.service.removeACEGroup,
		SocketServer:                      ac.server,
		OkMessage:                         "group removed",
	})
}

// CRUD пользователей в группах безопасности
func (ac *AceController) getUserACEGroups(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.UserReq, []*model.AceGroup]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                dto.UserReqFromJson,
		HandlerWithRequestWithResponse: ac.service.getUserACEGroups,
		SocketServer:                   ac.server,
	})
}
func (ac *AceController) addUserToACEGroup(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.UserInAccessGroupReq, []*model.AceGroup]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                dto.UserInAccessGroupReqFromJson,
		HandlerWithRequestWithResponse: ac.service.addUserToACEGroup,
		SocketServer:                   ac.server,
	})
}

func (ac *AceController) removeUserFromACEGroup(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.UserInAccessGroupReq, []*model.AceGroup]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                dto.UserInAccessGroupReqFromJson,
		HandlerWithRequestWithResponse: ac.service.removeUserFromACEGroup,
		SocketServer:                   ac.server,
	})
}

// CRUD объектов безопасности
func (ac *AceController) getACL(event string, ctx *model.SocketContext, _ []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[any, []*model.Ace]{
		Event:                             event,
		SocketContext:                     ctx,
		HandlerWithoutRequestWithResponse: ac.service.getACL,
		SocketServer:                      ac.server,
	})
}
func (ac *AceController) getACEsForGroup(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.ACEGroupIdReq, []*model.Ace]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                dto.ACEGroupIdReqFromJson,
		HandlerWithRequestWithResponse: ac.service.getACEsForGroup,
		SocketServer:                   ac.server,
	})
}
func (ac *AceController) addACEToGroup(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.ACEToGroupReq, any]{
		Event:                             event,
		SocketContext:                     ctx,
		Args:                              args,
		FromJsonHandler:                   dto.ACEToGroupReqFromJson,
		HandlerWithRequestWithoutResponse: ac.service.addACEToGroup,
		SocketServer:                      ac.server,
		OkMessage:                         "ace added",
	})
}

func (ac *AceController) updateACE(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.ACEToGroupReq, any]{
		Event:                             event,
		SocketContext:                     ctx,
		Args:                              args,
		FromJsonHandler:                   dto.ACEToGroupReqFromJson,
		HandlerWithRequestWithoutResponse: ac.service.updateACE,
		SocketServer:                      ac.server,
		OkMessage:                         "ace updated",
	})
}

func (ac *AceController) removeACE(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.ACEIdReq, any]{
		Event:                             event,
		SocketContext:                     ctx,
		Args:                              args,
		FromJsonHandler:                   dto.ACEIdReqFromJson,
		HandlerWithRequestWithoutResponse: ac.service.removeACE,
		SocketServer:                      ac.server,
		OkMessage:                         "ace updated",
	})
}
