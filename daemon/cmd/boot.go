package cmd

import (
	"unbalance/daemon/domain"
	"unbalance/daemon/services"
)

type Boot struct {
}

func (b *Boot) Run(ctx *domain.Context) error {
	return services.CreateOrchestrator(ctx).Run()
}
