package core

import (
	"testing"

	"unbalance/daemon/common"
	"unbalance/daemon/domain"
)

func TestCreateReplayOperationDoesNotMutateOriginal(t *testing.T) {
	c := &Core{}
	originalCommand := &domain.Command{
		ID:          "original-command",
		Src:         "/mnt/disk1",
		Dst:         "/mnt/disk2/",
		Entry:       "share/file.bin",
		Size:        123,
		Transferred: 123,
		Status:      common.CmdCompleted,
	}
	original := domain.Operation{
		ID:              "original-operation",
		OpKind:          common.OpScatterMove,
		BytesToTransfer: 123,
		RsyncArgs:       []string{"-avPR", "-X"},
		Commands:        []*domain.Command{originalCommand},
	}

	replay := c.createReplayOperation(original)
	if replay.ID == original.ID {
		t.Fatalf("replay operation reused original id")
	}
	if len(replay.Commands) != 1 {
		t.Fatalf("replay command count = %d, want 1", len(replay.Commands))
	}

	replayCommand := replay.Commands[0]
	if replayCommand == originalCommand {
		t.Fatalf("replay command reused original command pointer")
	}
	if replayCommand.ID == originalCommand.ID {
		t.Fatalf("replay command reused original command id")
	}
	if replayCommand.Transferred != 0 {
		t.Fatalf("replay transferred = %d, want 0", replayCommand.Transferred)
	}
	if replayCommand.Status != common.CmdPending {
		t.Fatalf("replay status = %d, want pending", replayCommand.Status)
	}

	if originalCommand.ID != "original-command" {
		t.Fatalf("original command id mutated to %q", originalCommand.ID)
	}
	if originalCommand.Transferred != 123 {
		t.Fatalf("original transferred mutated to %d", originalCommand.Transferred)
	}
	if originalCommand.Status != common.CmdCompleted {
		t.Fatalf("original status mutated to %d", originalCommand.Status)
	}

	replay.RsyncArgs[1] = "--changed"
	if original.RsyncArgs[1] != "-X" {
		t.Fatalf("original rsync args mutated to %v", original.RsyncArgs)
	}
}
