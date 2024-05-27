package guidemodule

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type GuideController struct {
	store   store.Store
	server  model.SocketServer
	service *GuideService
}

func InitController(socketService model.SocketServer, store store.Store, auth *middleware.Auth) *GuideController {
	gc := &GuideController{
		store:   store,
		server:  socketService,
		service: NewService(store),
	}

	socketService.On("api/v1/guide/content", gc.getContent)
	socketService.On("api/v1/guide/article", gc.getArticle)
	socketService.On("api/v1/guide/firstarticlename", gc.getFirstArticleName)

	auth.AddFreeRoute("api/v1/guide/content")
	auth.AddFreeRoute("api/v1/guide/article")
	auth.AddFreeRoute("api/v1/guide/firstarticlename")

	return gc
}

func (gc *GuideController) getContent(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.GetGuideReq, *dto.GetGuideRes]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                dto.FromJson[dto.GetGuideReq],
		HandlerWithRequestWithResponse: gc.service.loadGuide,
		SocketServer:                   gc.server,
	})
}
func (gc *GuideController) getFirstArticleName(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.GetFirstArticleNameReq, *dto.GetFirstArticleNameRes]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                dto.FromJson[dto.GetFirstArticleNameReq],
		HandlerWithRequestWithResponse: gc.service.firstArticleName,
		SocketServer:                   gc.server,
	})
}
func (gc *GuideController) getArticle(event string, ctx *model.SocketContext, args []byte) {
	model.ExecuteSocketRequest(&model.SocketRequestContext[*dto.GetArticleReq, *dto.GetArticleRes]{
		Event:                          event,
		SocketContext:                  ctx,
		Args:                           args,
		FromJsonHandler:                dto.FromJson[dto.GetArticleReq],
		HandlerWithRequestWithResponse: gc.service.loadArticle,
		SocketServer:                   gc.server,
	})
}
