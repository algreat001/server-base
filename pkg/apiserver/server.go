package apiserver

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/controller"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/config"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type Server struct {
	store       store.Store
	router      *gin.Engine
	controllers *controller.Controllers
}

func NewServer(store store.Store) *Server {

	s := &Server{
		store:  store,
		router: gin.Default(),
	}

	s.router.Use(ginlogrus.Logger(logrus.StandardLogger()), gin.Recovery())
	s.ConfigureRouter()
	return s
}

func (s *Server) Start(bindAddr string) error {
	if config.GetInstance().SSL.CertPath != "" {
		return s.router.RunTLS(bindAddr, config.GetInstance().SSL.CertPath, config.GetInstance().SSL.KeyPath)
	}
	return s.router.Run(bindAddr)
}

func GetCORSConfig() cors.Config {
	return cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
			"Options",
			"Origin",
			"Accept",
			"X-Requested-With",
			"Authorization",
			"Content-Length",
			"Content-Type",
			"Accept-Encoding",
			"Accept-Language",
			"Cache-Control",
			"Connection",
			"Cookie",
			"Host",
			"Pragma",
			"User-Agent",
			"Upgrade",
			"Sec-Websocket-Extensions",
			"Sec-Websocket-Key",
			"Sec-Websocket-Version",
			"Sec-Websocket-Protocol",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}
}

func (s *Server) ConfigureRouter() {
	err := s.router.SetTrustedProxies(nil)

	if err != nil {
		logrus.Fatal(err)
	}

	s.router.Use(cors.New(GetCORSConfig()))

	auth := middleware.NewAuth(s.store, s.router)
	//auth.Apply()

	middleware.NewStatic(config.GetInstance().StaticPath, s.store, s.router).Apply()

	s.controllers = controller.New(s.router, s.store, auth)

}
