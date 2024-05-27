package middleware

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type Auth struct {
	store       store.Store
	router      *gin.Engine
	freeRoutes  *helpers.SpecialRoute
	tokenRoutes *helpers.SpecialRoute
}

func NewAuth(store store.Store, router *gin.Engine) *Auth {
	return &Auth{
		store:       store,
		router:      router,
		freeRoutes:  helpers.NewSpecialRoute(),
		tokenRoutes: helpers.NewSpecialRoute(),
	}
}

func (a *Auth) Apply() {
	a.router.Use(a.MiddlewareFunc())
}

func (a *Auth) MiddlewareFunc() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {

		path := c.Request.URL.Path
		if a.freeRoutes.IsContained(path) {
			logrus.Info("Access to free route - ", path)
			c.Next()
			return
		}

		if a.tokenRoutes.IsContained(path) {
			token := c.Query("token")
			if err := a.store.Token().VerifyToken(token); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": servererrors.ErrorAccessDenied.Error()})
				c.Abort()
				return
			}
			logrus.Info("Access to token route - ", path)
			c.Next()
			return
		}

		user, err := model.AuthTokenFromContext(c).GetUserFromVerifyAuthToken()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		logrus.Info("user", user)

		ok := a.TestAuthRoute(path, user)

		if !ok {
			logrus.Warn("Access denied - ", user.Id, ", email: ", user.Email, " to ", path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": servererrors.ErrorAccessDeined.Error()})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	})
}

func (a *Auth) TestAuthRoute(path string, user *model.User) bool {
	if a.freeRoutes.IsContained(path) {
		logrus.Info("Access to free route - ", path)
		return true
	}
	ok, err := a.store.AccessControlEntity().GetUserRightForPath(user, path)
	if err != nil {
		logrus.Warn("invalid data base query for ACE - ", err)
		return false
	}
	return ok
}

func (a *Auth) AddFreeRoute(path string) {
	a.freeRoutes.Add(path)
	return
}

func (a *Auth) AddTokenRoute(path string) {
	a.tokenRoutes.Add(strings.ToLower(path))
	return
}
