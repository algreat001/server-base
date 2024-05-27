package middleware

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type Static struct {
	store  store.Store
	router *gin.Engine
	path   string
}

func NewStatic(path string, store store.Store, router *gin.Engine) *Static {
	return &Static{
		store:  store,
		router: router,
		path:   path,
	}

}

func (s *Static) Apply() {
	logrus.Info("Static web server is started. local path = ", s.path)
	s.router.NoRoute(func(c *gin.Context) {
		c.File(s.path + "index.html")
	})
	s.router.Use(s.middlewareFunc())
}

func (s *Static) middlewareFunc() gin.HandlerFunc {
	return static.Serve("/", static.LocalFile(s.path, false))
}
