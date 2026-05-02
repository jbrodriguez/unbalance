package core

import (
	"testing"

	"unbalance/daemon/domain"
)

func TestSetRsyncArgsRejectsBlockedArgAndKeepsConfig(t *testing.T) {
	c := &Core{ctx: &domain.Context{Config: domain.Config{RsyncArgs: []string{"-X"}}}}

	_, err := c.SetRsyncArgs([]string{"--remove-source-files"})
	if err == nil {
		t.Fatalf("expected blocked arg to be rejected")
	}

	if len(c.ctx.Config.RsyncArgs) != 1 || c.ctx.Config.RsyncArgs[0] != "-X" {
		t.Fatalf("config changed after rejected arg: %v", c.ctx.Config.RsyncArgs)
	}
}
