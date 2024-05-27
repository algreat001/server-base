package usermodule

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type UserController struct {
	store   store.Store
	server  model.SocketServer
	service *UserService
}

func InitController(socketService model.SocketServer, store store.Store, auth *middleware.Auth) *UserController {
	uc := &UserController{
		store:   store,
		server:  socketService,
		service: NewService(store),
	}

	socketService.On("api/v1/user/add", uc.addUser)
	socketService.On("api/v1/user/update", uc.updateUser)
	socketService.On("api/v1/user/remove", uc.removeUser)
	socketService.On("api/v1/user/page", uc.getOnePageUsers)

	auth.AddFreeRoute("api/v1/profile/groups")
	socketService.On("api/v1/profile/groups", uc.getProfileGroups)

	return uc
}

func (uc *UserController) sendError(event string, ctx *model.SocketContext, error string) {
	uc.server.SendMessagesToUser(ctx.GetUser().Id, &model.SocketResponseMessage{
		Event: event,
		Data:  map[string]any{"status": "error", "error": error},
	})
}

func (uc *UserController) getOnePageUsers(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.UsersPageReq, *dto.UserPage]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                dto.UsersPageReqFromJson,
		HandlerWithRequestWithResponse: uc.service.getUsers,
		SocketServer:                   uc.server,
	})
}

func (uc *UserController) addUser(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.UserUpdateReq, *model.User]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                dto.UserUpdateReqFromJson,
		HandlerWithRequestWithResponse: uc.service.addUser,
		SocketServer:                   uc.server,
	})
}

func (uc *UserController) updateUser(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.UserUpdateReq, *model.User]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                dto.UserUpdateReqFromJson,
		HandlerWithRequestWithResponse: uc.service.updateUser,
		SocketServer:                   uc.server,
	})
}

func (uc *UserController) removeUser(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.UserUpdateReq, any]{
		Event:                             event,
		SocketContext:                     ctx,
		Args:                              args,
		FromJsonHandler:                   dto.UserUpdateReqFromJson,
		HandlerWithRequestWithoutResponse: uc.service.removeUser,
		SocketServer:                      uc.server,
		OkMessage:                         "user removed",
	})
}

func (uc *UserController) getProfileGroups(event string, ctx *model.SocketContext, _ []byte) {
	args, err := ctx.GetUser().ToJson()
	if err != nil {
		uc.sendError(event, ctx, servererrors.ErrorInternal.Error())
		return
	}

	model.ExecuteSocketRequest(&model.SocketRequestContext[*model.User, []*model.AceGroup]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                model.UserFromJson,
		HandlerWithRequestWithResponse: uc.service.getGroups,
		SocketServer:                   uc.server,
	})

}
