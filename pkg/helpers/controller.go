package helpers

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Response[TResponse any] struct {
	Status string     `json:"status"`
	Error  string     `json:"error,omitempty"`
	Data   *TResponse `json:"data,omitempty"`
}

type StandardResponse[TResponse any] struct {
	Response      *Response[TResponse] `json:"data,omitempty"`
	UpdateMeta    *[]byte              `json:"update-mate,omitempty"`
	UpdateCommand string               `json:"update-command,omitempty"`
}

type GinRequestContext[TRequest any, TResponse any] struct {
	GinContext                            *gin.Context
	StandardHandler                       func(request *TRequest) (*StandardResponse[TResponse], error)
	HandlerWithRequestWithResponse        func(request *TRequest) (*TResponse, error)
	HandlerWithRequestWithoutResponse     func(request *TRequest) error
	HandlerWithoutRequestWithResponse     func() (*TResponse, error)
	UserHandlerWithRequestWithResponse    func(user *model.User, request *TRequest) (*TResponse, error)
	UserHandlerWithRequestWithoutResponse func(user *model.User, request *TRequest) error
	UserHandlerWithoutRequestWithResponse func(user *model.User) (*TResponse, error)
}

func executeGinRequest[TRequest any, TResponse any](ctx *GinRequestContext[TRequest, TResponse]) (*TRequest, *TResponse, error) {
	if ctx.HandlerWithRequestWithResponse != nil {
		return ginRequestWithResponseItem(ctx)
	}
	if ctx.HandlerWithRequestWithoutResponse != nil {
		return ginRequestWithoutResponseItem(ctx)
	}
	if ctx.HandlerWithoutRequestWithResponse != nil {
		return ginRequestWithoutRequestBody(ctx)
	}
	if ctx.UserHandlerWithRequestWithResponse != nil {
		return ginRequestUserWithResponseItem(ctx)
	}
	if ctx.UserHandlerWithRequestWithoutResponse != nil {
		return ginRequestUserWithoutResponseItem(ctx)
	}
	if ctx.UserHandlerWithoutRequestWithResponse != nil {
		return ginRequestUserWithoutRequestBody(ctx)
	}
	return nil, nil, servererrors.ErrorUnknownTypeRequest
}

func ExecuteGinRequest[TRequest any, TResponse any](ctx *GinRequestContext[TRequest, TResponse]) (*TResponse, error) {
	_, res, err := executeGinRequest(ctx)
	return res, err
}

func ExecuteGinStandardRequest[TRequest any, TResponse any](ctx *GinRequestContext[TRequest, TResponse]) (*TRequest, *StandardResponse[TResponse], error) {
	request := new(TRequest)
	if err := ctx.GinContext.ShouldBindJSON(request); err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": servererrors.ErrorInvalidRequest.Error()})
		return nil, nil, err
	}
	response, err := ctx.StandardHandler(request)
	if err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": servererrors.ErrorInternal.Error()})
		return request, nil, err
	}
	if response.Response.Status == "error" {
		logrus.Infof("logical execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusOK, gin.H{"status": "error", "error": response.Response.Error, "data": response.Response.Data})
	} else {
		logrus.Infof("successful execution (handler name: %s)", ctx.GinContext.HandlerName())
		ctx.GinContext.JSON(http.StatusOK, gin.H{"status": response.Response.Status, "data": response.Response.Data})
	}
	return request, response, nil
}

func ginRequestWithResponseItem[TRequest any, TResponse any](ctx *GinRequestContext[TRequest, TResponse]) (*TRequest, *TResponse, error) {
	request := new(TRequest)
	if err := ctx.GinContext.ShouldBindJSON(request); err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": servererrors.ErrorInvalidRequest.Error()})
		return nil, nil, err
	}
	data, err := ctx.HandlerWithRequestWithResponse(request)
	if err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": servererrors.ErrorInternal.Error()})
		return request, nil, err
	}
	logrus.Infof("successful execution (handler name: %s)", ctx.GinContext.HandlerName())
	ctx.GinContext.JSON(http.StatusOK, gin.H{"status": "ok", "data": data})
	return request, data, nil
}
func ginRequestWithoutResponseItem[TRequest any, TResponse any](ctx *GinRequestContext[TRequest, TResponse]) (*TRequest, *TResponse, error) {
	request := new(TRequest)
	if err := ctx.GinContext.ShouldBindJSON(request); err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": servererrors.ErrorInvalidRequest.Error()})
		return nil, nil, err
	}
	if err := ctx.HandlerWithRequestWithoutResponse(request); err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": servererrors.ErrorInternal.Error()})
		return request, nil, err
	}
	logrus.Infof("successful execution (handler name: %s)", ctx.GinContext.HandlerName())
	ctx.GinContext.JSON(http.StatusOK, gin.H{"status": "ok"})
	return request, nil, nil
}
func ginRequestWithoutRequestBody[TRequest any, TResponse any](ctx *GinRequestContext[TRequest, TResponse]) (*TRequest, *TResponse, error) {
	data, err := ctx.HandlerWithoutRequestWithResponse()
	if err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": servererrors.ErrorInternal.Error()})
		return nil, nil, err
	}
	logrus.Infof("successful execution (handler name: %s)", ctx.GinContext.HandlerName())
	ctx.GinContext.JSON(http.StatusOK, gin.H{"status": "ok", "data": data})
	return nil, data, nil
}

// //////////////////////
func ginRequestUserWithResponseItem[TRequest any, TResponse any](ctx *GinRequestContext[TRequest, TResponse]) (*TRequest, *TResponse, error) {
	request := new(TRequest)
	if err := ctx.GinContext.ShouldBindJSON(request); err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": servererrors.ErrorInvalidRequest.Error()})
		return nil, nil, err
	}
	user := ctx.GinContext.MustGet("user").(*model.User)
	data, err := ctx.UserHandlerWithRequestWithResponse(user, request)
	if err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": servererrors.ErrorInternal.Error()})
		return request, nil, err
	}
	logrus.Infof("successful execution (handler name: %s)", ctx.GinContext.HandlerName())
	ctx.GinContext.JSON(http.StatusOK, gin.H{"status": "ok", "data": data})
	return request, data, nil
}
func ginRequestUserWithoutResponseItem[TRequest any, TResponse any](ctx *GinRequestContext[TRequest, TResponse]) (*TRequest, *TResponse, error) {
	request := new(TRequest)
	if err := ctx.GinContext.ShouldBindJSON(request); err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": servererrors.ErrorInvalidRequest.Error()})
		return nil, nil, err
	}
	user := ctx.GinContext.MustGet("user").(*model.User)
	if err := ctx.UserHandlerWithRequestWithoutResponse(user, request); err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": servererrors.ErrorInternal.Error()})
		return request, nil, err
	}
	logrus.Infof("successful execution (handler name: %s)", ctx.GinContext.HandlerName())
	ctx.GinContext.JSON(http.StatusOK, gin.H{"status": "ok"})
	return request, nil, nil
}
func ginRequestUserWithoutRequestBody[TRequest any, TResponse any](ctx *GinRequestContext[TRequest, TResponse]) (*TRequest, *TResponse, error) {
	user := ctx.GinContext.MustGet("user").(*model.User)
	data, err := ctx.UserHandlerWithoutRequestWithResponse(user)
	if err != nil {
		logrus.Infof("execution error (handler name: %s, error: %s)", ctx.GinContext.HandlerName(), err)
		ctx.GinContext.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": servererrors.ErrorInternal.Error()})
		return nil, nil, err
	}
	logrus.Infof("successful execution (handler name: %s)", ctx.GinContext.HandlerName())
	ctx.GinContext.JSON(http.StatusOK, gin.H{"status": "ok", "data": data})
	return nil, data, nil
}
