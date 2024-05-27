package controller

import (
	"github.com/gin-gonic/gin"

	acemodule "gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/controller/aceModule"
	busmodule "gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/controller/busModule"
	feedbackmodule "gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/controller/feedbackModule"
	guidemodule "gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/controller/guideModule"
	langmodule "gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/controller/langModule"
	logmodule "gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/controller/logModule"
	pingmodule "gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/controller/pingModule"
	socketmodule "gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/controller/socketModule"
	usermodule "gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/controller/userModule"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/config"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type Controller interface {
	InitRouters(*gin.Engine, store.Store)
}

type Controllers struct {
	router       *gin.Engine
	store        store.Store
	auth         *middleware.Auth
	SocketServer model.SocketServer
	ServiceBus   *model.ServiceBus
	BusModule    *busmodule.BusController
}

func New(router *gin.Engine, store store.Store, auth *middleware.Auth) *Controllers {
	c := &Controllers{
		router: router,
		store:  store,
		auth:   auth,
	}

	langmodule.InitController(router, store, auth)
	feedbackmodule.InitController(router, store, auth)
	c.SocketServer = socketmodule.InitSocketServer(router, store, auth)

	// ping logic
	pingmodule.InitController(c.SocketServer, auth)

	// admin logic
	usermodule.InitController(c.SocketServer, store, auth)
	acemodule.InitController(c.SocketServer, store, auth)
	logmodule.InitController(c.SocketServer, store, auth)
	guidemodule.InitController(c.SocketServer, store, auth)
	// init bus
	c.ServiceBus = model.NewServiceBus()
	c.BusModule = busmodule.InitBusController(store, c.ServiceBus, auth)

	err := c.ServiceBus.Listen(config.GetInstance().ServiceId)
	if err != nil {
		panic(err)
	}

	return c
}
