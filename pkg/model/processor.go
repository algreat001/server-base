package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
)

type SendInfoFn func(msg any)
type PrepareProcessorFn func(ctx *ProcessorContext) ([][]string, map[string]any, error)
type ProcessorFn func(ctx *ProcessorContext, item []string, args map[string]any) string
type ProcessorCompleteFn func(ctx *ProcessorContext, args map[string]any) string

type ProcessorInfo struct {
	Msg     string  `json:"msg"`
	Percent float64 `json:"percent"`
}

type ProcessorContext struct {
	Id      uuid.UUID
	Name    string
	Result  chan error
	Info    chan ProcessorInfo
	Context context.Context
	cancel  context.CancelFunc
	Percent float64
	IsRun   bool
	Timeout time.Duration
	FuncCtx *ProcessorFunctionContext
}

func NewProcessorContext(name string, funcCtx *ProcessorFunctionContext, timeOut time.Duration) *ProcessorContext {
	return &ProcessorContext{
		Id:      uuid.New(),
		Name:    name,
		Result:  make(chan error, 1),
		Info:    make(chan ProcessorInfo, 1),
		Percent: 0,
		IsRun:   false,
		Timeout: timeOut,
		FuncCtx: funcCtx,
	}
}

type Processor struct {
	context *ProcessorContext
}

func NewProcessor(name string, funcCtx *ProcessorFunctionContext) *Processor {
	return &Processor{
		context: NewProcessorContext(name, funcCtx, 60*time.Minute),
	}
}

func (p *Processor) GetContext() *ProcessorContext {
	return p.context
}

func (p *Processor) Run(sendMessageFn SendInfoFn, sendErrorFn SendInfoFn) error {
	if p.context.IsRun {
		return servererrors.ErrorProcessAlreadyRunning
	}
	p.context.Result = make(chan error, 1)
	p.context.Info = make(chan ProcessorInfo, 1)
	p.context.IsRun = true
	p.context.Context, p.context.cancel = context.WithTimeout(context.Background(), p.context.Timeout)

	go p.process(p.context.FuncCtx)

	go p.monitor(sendMessageFn, sendErrorFn)
	return nil
}

func (p *Processor) Stop() {
	if !p.context.IsRun || p.context.cancel == nil {
		logrus.Info("Processor [" + p.context.Name + "] is not run")
		return
	}
	logrus.Info("Send signal 'stop' to processor [" + p.context.Name + "]")
	p.context.cancel()
}

func (p *Processor) monitor(sendMessageFn SendInfoFn, sendErrorFn SendInfoFn) {
	sendMessageFn(p.context.Name + " process started")

	defer func() {
		close(p.context.Result)
		close(p.context.Info)
		p.context.cancel()
		p.context.cancel = nil
		p.context.IsRun = false
	}()

	for {
		select {
		case err := <-p.context.Result:
			{
				if err != nil {
					sendErrorFn(p.context.Name + " error - " + err.Error())
				}
				sendMessageFn(p.context.Name + " process finished")
				return
			}

		case info := <-p.context.Info:
			sendMessageFn(info)
		}
	}
}

func (p *Processor) getCorrectionFactor(count int) float64 {
	if count == 0 {
		return 1
	}
	return 1 / float64(count)
}

func (p *Processor) process(funcCtx *ProcessorFunctionContext) {
	processorContext := p.context

	defer func() {
		if r := recover(); r != nil {
			logrus.Error("Process ["+processorContext.Name+"] panic: ", r)
			processorContext.Result <- nil
		}
	}()

	data, args, err := funcCtx.PrepareFn(processorContext)
	if err != nil {
		processorContext.Result <- err
		return
	}

	correction := p.getCorrectionFactor(len(data))

	for index, item := range data {
		select {
		case <-processorContext.Context.Done():
			logrus.Info("Process [" + processorContext.Name + "] canceled")
			processorContext.Result <- nil
			return
		default:
			processorContext.Percent = float64(index) * correction
			processorContext.Info <- ProcessorInfo{Msg: funcCtx.IterationFn(processorContext, item, args), Percent: processorContext.Percent}
		}
	}
	processorContext.Info <- ProcessorInfo{Msg: funcCtx.CompleteFn(processorContext, args), Percent: processorContext.Percent}

	processorContext.Result <- nil
}
