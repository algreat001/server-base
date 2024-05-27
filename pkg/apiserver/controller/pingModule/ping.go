package pingmodule

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
)

type PingSocketController struct {
	server model.SocketServer
	auth   *middleware.Auth
}

func InitController(socketService model.SocketServer, auth *middleware.Auth) *PingSocketController {
	service := &PingSocketController{
		server: socketService,
		auth:   auth,
	}
	service.subscribe()
	return service
}

func (c *PingSocketController) subscribe() {
	c.auth.AddFreeRoute("ping")
	c.server.On("ping", c.pingHandler)
}

func (c *PingSocketController) pingHandler(event string, ctx *model.SocketContext, _ []byte) {
	c.server.SendMessageToSocket(ctx.GetSocket(), &model.SocketResponseMessage{Event: event, Data: map[string]any{}})
}
