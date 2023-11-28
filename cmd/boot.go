package cmd

import (
	"unbalance/domain"
	"unbalance/services"
)

type Boot struct {
}

func (b *Boot) Run(ctx *domain.Context) error {
	return services.CreateOrchestrator(ctx).Run()
}
