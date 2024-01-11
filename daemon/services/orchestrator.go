package services

import (
	"os"
	"os/signal"
	"syscall"

	"unbalance/daemon/domain"
	"unbalance/daemon/logger"
	"unbalance/daemon/services/core"
	"unbalance/daemon/services/server"
)

type Orchestrator struct {
	ctx *domain.Context
}

func CreateOrchestrator(ctx *domain.Context) *Orchestrator {
	return &Orchestrator{
		ctx: ctx,
	}
}

func (o *Orchestrator) Run() error {
	logger.Blue("starting unbalanced %s ...", o.ctx.Version)

	core := core.Create(o.ctx)
	server := server.Create(o.ctx, core)

	err := server.Start()
	if err != nil {
		return err
	}

	err = core.Start()
	if err != nil {
		return err
	}

	w := make(chan os.Signal, 1)
	signal.Notify(w, syscall.SIGTERM, syscall.SIGINT)
	logger.Blue("received %s signal. shutting down the app ...", <-w)

	return nil
}
