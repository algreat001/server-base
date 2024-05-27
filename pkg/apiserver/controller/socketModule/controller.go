package socketmodule

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type SocketController struct {
	store    store.Store
	upgrader websocket.Upgrader
	Service  model.SocketServer
}

func InitSocketServer(router *gin.Engine, store store.Store, auth *middleware.Auth) model.SocketServer {
	sc := &SocketController{
		store:   store,
		Service: NewService(auth),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}

	router.GET("api/s1/connect", auth.MiddlewareFunc(), sc.websocketHandler)
	router.POST("api/s1/connect", auth.MiddlewareFunc(), sc.websocketHandler)

	return sc.Service
}

func (sc *SocketController) websocketHandler(c *gin.Context) {
	conn, err := sc.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Error("Failed to upgrade to WebSocket: ", err.Error())
		return
	}
	userRec, _ := c.Get("user")
	user, ok := userRec.(*model.User)
	if !ok {
		logrus.Error("Failed to get user from context")
		return
	}
	user.Groups, err = sc.store.AccessControlEntity().GetUserGroups(user)
	if err != nil {
		logrus.Error("Failed to get user groups on connection to socket: ", user.Id, err.Error())
		return
	}
	sc.Service.StartSocketHandler(conn, user)
}
