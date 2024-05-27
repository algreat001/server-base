package feedbackmodule

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type FeedbackController struct {
	store   store.Store
	service *FeedbackService
}

func InitController(router *gin.Engine, store store.Store, auth *middleware.Auth) *FeedbackController {
	fc := &FeedbackController{
		store:   store,
		service: NewService(store),
	}

	router.POST("/api/v1/feedback/send", fc.send)

	return fc
}

func (fc *FeedbackController) send(c *gin.Context) {
	req := &dto.SendReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Info("error parse json (send) - ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": servererrors.ErrorInvalidRequest.Error()})
		return
	}

	err := fc.service.send(req)
	if err != nil {
		logrus.Info("error send message to telegram - ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": servererrors.ErrorInternal.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})

}
