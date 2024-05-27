package model

import (
	"errors"
	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
)

type SocketRequestContext[TRequestItem any, TResponse any] struct {
	Event                             string
	SocketContext                     *SocketContext
	Args                              []byte
	SocketServer                      SocketServer
	FromJsonHandler                   func(args []byte) (TRequestItem, error)
	HandlerWithRequestWithoutResponse func(ctx *SocketContext, item TRequestItem) error
	HandlerWithRequestWithResponse    func(ctx *SocketContext, item TRequestItem) (TResponse, error)
	HandlerWithoutRequestWithResponse func(ctx *SocketContext) (TResponse, error)
	OkMessage                         string
}

func ExecuteSocketRequest[TRequestItem any, TResponse any](ctx *SocketRequestContext[TRequestItem, TResponse]) error {
	if ctx.HandlerWithRequestWithoutResponse != nil {
		return socketRequestWithoutResponseItem(ctx)
	}
	if ctx.HandlerWithRequestWithResponse != nil {
		return socketRequestWithResponseItem(ctx)
	}
	if ctx.HandlerWithoutRequestWithResponse != nil {
		return socketRequestWithoutRequestBody(ctx)
	}
	return servererrors.ErrorSocketUnknownHandler
}

func socketRequestWithoutResponseItem[TRequestItem any, TResponse any](ctx *SocketRequestContext[TRequestItem, TResponse]) error {
	if ctx.FromJsonHandler == nil || (ctx.HandlerWithRequestWithoutResponse == nil) || ctx.SocketServer == nil || ctx.SocketContext == nil {
		logrus.Info(ctx.Event, " - error initialize socket request context")
		return servererrors.ErrorSocketReqContext
	}
	item, err := ctx.FromJsonHandler(ctx.Args)
	if err != nil {
		logrus.Info(ctx.Event, " - error parse json - ", err)
		ctx.SocketServer.SendError(ctx.Event, ctx.SocketContext.GetSocket(), servererrors.ErrorInvalidRequest.Error())
		return err
	}

	err = ctx.HandlerWithRequestWithoutResponse(ctx.SocketContext, item)
	if err != nil {
		logrus.Info(ctx.Event, " - error execute handler - ", err)
		if errors.Is(err, servererrors.ErrorInternal) {
			ctx.SocketServer.SendError(ctx.Event, ctx.SocketContext.GetSocket(), err.Error())
		} else {
			ctx.SocketServer.SendError(ctx.Event, ctx.SocketContext.GetSocket(), servererrors.ErrorInternal.Error())
		}
		return err
	}

	ctx.SocketServer.SendMessageToSocket(ctx.SocketContext.GetSocket(), &SocketResponseMessage{
		Event: ctx.Event,
		Data:  &SocketResponse{Status: "ok", Data: ctx.OkMessage},
	})
	return nil
}

func socketRequestWithResponseItem[TRequestItem any, TResponse any](ctx *SocketRequestContext[TRequestItem, TResponse]) error {
	if ctx.FromJsonHandler == nil || ctx.HandlerWithRequestWithResponse == nil || ctx.SocketServer == nil || ctx.SocketContext == nil {
		logrus.Info(ctx.Event, " - error initialize socket request context")
		return servererrors.ErrorSocketReqContext
	}

	item, err := ctx.FromJsonHandler(ctx.Args)
	if err != nil {
		logrus.Info(ctx.Event, " - error parse json - ", err)
		ctx.SocketServer.SendError(ctx.Event, ctx.SocketContext.GetSocket(), servererrors.ErrorInvalidRequest.Error())
		return err
	}

	result, err := ctx.HandlerWithRequestWithResponse(ctx.SocketContext, item)

	if err != nil {
		logrus.Info(ctx.Event, " - error execute handler - ", err)
		if errors.Is(err, servererrors.ErrorInternal) {
			ctx.SocketServer.SendError(ctx.Event, ctx.SocketContext.GetSocket(), err.Error())
		} else {
			ctx.SocketServer.SendError(ctx.Event, ctx.SocketContext.GetSocket(), servererrors.ErrorInternal.Error())
		}
		return err
	}

	ctx.SocketServer.SendMessageToSocket(ctx.SocketContext.GetSocket(), &SocketResponseMessage{
		Event: ctx.Event,
		Data:  &SocketResponse{Status: "ok", Data: result},
	})
	return nil
}

func socketRequestWithoutRequestBody[TRequestItem any, TResponse any](ctx *SocketRequestContext[TRequestItem, TResponse]) error {
	if ctx.HandlerWithoutRequestWithResponse == nil || ctx.SocketServer == nil || ctx.SocketContext == nil {
		logrus.Info(ctx.Event, " - error initialize socket request context")
		return servererrors.ErrorSocketReqContext
	}

	result, err := ctx.HandlerWithoutRequestWithResponse(ctx.SocketContext)

	if err != nil {
		logrus.Info(ctx.Event, " - error execute handler - ", err)
		if errors.Is(err, servererrors.ErrorInternal) {
			ctx.SocketServer.SendError(ctx.Event, ctx.SocketContext.GetSocket(), err.Error())
		} else {
			ctx.SocketServer.SendError(ctx.Event, ctx.SocketContext.GetSocket(), servererrors.ErrorInternal.Error())
		}
		return err
	}

	ctx.SocketServer.SendMessageToSocket(ctx.SocketContext.GetSocket(), &SocketResponseMessage{
		Event: ctx.Event,
		Data:  &SocketResponse{Status: "ok", Data: result},
	})
	return nil
}
