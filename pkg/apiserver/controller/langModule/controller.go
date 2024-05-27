package langmodule

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"net/http"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type LangController struct {
	store   store.Store
	service *LangService
}

func InitController(router *gin.Engine, store store.Store, auth *middleware.Auth) *LangController {
	lc := &LangController{
		store:   store,
		service: NewService(store),
	}

	auth.AddFreeRoute("/api/v1/lang/get")
	auth.AddFreeRoute("/api/v1/lang/list")

	router.POST("/api/v1/lang/get", lc.getLang)
	router.POST("/api/v1/lang/list", lc.getLangList)

	return lc
}

func (lc *LangController) getLangList(c *gin.Context) {
	result, err := lc.service.getLangList()
	if err != nil {
		logrus.Info("error get language list - ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": servererrors.ErrorRecordNotAdd.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (lc *LangController) getLang(c *gin.Context) {
	req := &dto.LangReq{}

	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Info("error parse json (getLang) - ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": servererrors.ErrorInvalidRequest.Error()})
		return
	}

	result, err := lc.service.getLang(req)
	if err != nil {
		logrus.Info("error get language - ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": servererrors.ErrorRecordNotAdd.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
