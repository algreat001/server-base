package model

import (
	"strings"

	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
)

type ProcessorFunctionContext struct {
	PrepareFn   PrepareProcessorFn
	IterationFn ProcessorFn
	CompleteFn  ProcessorCompleteFn
}

type SocketProcessorManager struct {
	server     SocketServer
	processors map[string]*Processor
}

func NewSocketProcessorManager(server SocketServer) *SocketProcessorManager {
	return &SocketProcessorManager{
		server:     server,
		processors: make(map[string]*Processor),
	}
}

func getNameProcess(event string, ctx *SocketContext) string {
	return event + ";" + ctx.GetId().String()
}

func (spm *SocketProcessorManager) StopProcess(event string, ctx *SocketContext) {
	baseEvent := strings.Replace(event, "/stop", "", 1)
	err := spm.stopProcessAndRemoveFromMap(getNameProcess(baseEvent, ctx))
	if err != nil {
		logrus.Info("error stop process - ", err)
		spm.server.SendError(event, ctx.GetSocket(), servererrors.ErrorInvalidRequest.Error())
		return
	}
	return
}

func (spm *SocketProcessorManager) stopProcessAndRemoveFromMap(name string) error {
	if spm.processors[name] != nil && spm.processors[name].GetContext().IsRun {
		spm.processors[name].Stop()
		return nil
	}
	return servererrors.ErrorProcessNotFound
}

func (spm *SocketProcessorManager) vacuumProcessors() {
	for name, processor := range spm.processors {
		if !processor.GetContext().IsRun {
			delete(spm.processors, name)
		}
	}
}

func (spm *SocketProcessorManager) getSendMessageFn(event string, ctx *SocketContext) func(msg any) {
	return func(msg any) {
		spm.server.SendMessageToSocket(ctx.GetSocket(), &SocketResponseMessage{
			Event: event,
			Data:  &SocketResponse{Status: "ok", Data: msg},
		})
	}
}
func (spm *SocketProcessorManager) getSendErrorFn(event string, ctx *SocketContext) func(msg any) {
	return func(msg any) {
		logrus.Info(msg)
		spm.server.SendError(event, ctx.GetSocket(), servererrors.ErrorInternal.Error())
	}
}

func (spm *SocketProcessorManager) CreateProcessor(event string, ctx *SocketContext, nameProcessor string, funcCtx *ProcessorFunctionContext) (*Processor, error) {
	name := getNameProcess(event, ctx)

	if spm.processors[name] != nil && spm.processors[name].GetContext().IsRun {
		return nil, servererrors.ErrorProcessAlreadyRunning
	}

	spm.processors[name] = NewProcessor(nameProcessor, funcCtx)
	return spm.processors[name], nil
}

func (spm *SocketProcessorManager) RunProcessor(event string, ctx *SocketContext) error {
	name := getNameProcess(event, ctx)

	if spm.processors[name] == nil {
		spm.server.SendError(event, ctx.GetSocket(), servererrors.ErrorProcessNotFound.Error())
		return servererrors.ErrorProcessNotFound
	}

	if err := spm.processors[name].Run(spm.getSendMessageFn(event, ctx), spm.getSendErrorFn(event, ctx)); err != nil {
		spm.server.SendError(event, ctx.GetSocket(), err.Error())
		return err
	}
	spm.vacuumProcessors()
	return nil
}

func StartProcessorHelper[TItem any](
	spm *SocketProcessorManager,
	getModel func(args []byte) (TItem, error),
	getProcessorFunctionContext func(event string, ctx *SocketContext, item TItem) *ProcessorFunctionContext,
	nameProcessor string,
	event string,
	ctx *SocketContext,
	args []byte,
) {
	csv, err := getModel(args)
	if err != nil {
		logrus.Info(event, " - error parse json - ", err)
		spm.server.SendError(event, ctx.GetSocket(), servererrors.ErrorInvalidRequest.Error())
		return
	}
	funcCtx := getProcessorFunctionContext(event, ctx, csv)

	_, err = spm.CreateProcessor(event, ctx, nameProcessor, funcCtx)
	if err != nil {
		logrus.Info(event, " - error create processor - ", err)
		spm.server.SendError(event, ctx.GetSocket(), servererrors.ErrorInternal.Error())
	}

	if err := spm.RunProcessor(event, ctx); err != nil {
		logrus.Info(event, " - error run processor - ", err)
		spm.server.SendError(event, ctx.GetSocket(), servererrors.ErrorInternal.Error())
	}
}
