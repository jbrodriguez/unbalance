package core

import (
	"testing"
	"time"

	"unbalance/daemon/domain"

	"github.com/cskr/pubsub"
)

func TestUpdateSamplesSkipsCounterReset(t *testing.T) {
	c := &Core{
		ctx: &domain.Context{
			Config: domain.Config{SpeedWindow: "90s"},
			Hub:    pubsub.New(1),
		},
	}
	started := time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC)
	operation := &domain.Operation{Started: started}

	c.updateSamplesAt(operation, 100, started.Add(1*time.Second))
	c.updateSamplesAt(operation, 250, started.Add(3*time.Second))

	if len(operation.Samples) != 2 {
		t.Fatalf("samples = %d, want 2", len(operation.Samples))
	}
	if operation.Samples[0].Bytes != 100 {
		t.Fatalf("first sample bytes = %d, want 100", operation.Samples[0].Bytes)
	}
	if operation.Samples[0].SampledAt != started.Add(1*time.Second) {
		t.Fatalf("first sample time = %s, want %s", operation.Samples[0].SampledAt, started.Add(1*time.Second))
	}
	if operation.Samples[1].Bytes != 250 {
		t.Fatalf("second sample bytes = %d, want 250", operation.Samples[1].Bytes)
	}
	if operation.Samples[1].SampledAt != started.Add(3*time.Second) {
		t.Fatalf("second sample time = %s, want %s", operation.Samples[1].SampledAt, started.Add(3*time.Second))
	}

	c.updateSamplesAt(operation, 10, started.Add(4*time.Second))

	if operation.PrevSample != 10 {
		t.Fatalf("prev sample = %d, want 10", operation.PrevSample)
	}
	if operation.SampleIndex != 2 {
		t.Fatalf("sample index = %d, want unchanged sample count 2", operation.SampleIndex)
	}
	if len(operation.Samples) != 2 {
		t.Fatalf("samples = %d, want unchanged count 2", len(operation.Samples))
	}

	c.updateSamplesAt(operation, 60, started.Add(5*time.Second))

	if len(operation.Samples) != 3 {
		t.Fatalf("samples = %d, want 3", len(operation.Samples))
	}
	if operation.Samples[2].Bytes != 60 {
		t.Fatalf("post-reset sample bytes = %d, want 60", operation.Samples[2].Bytes)
	}
	if operation.Samples[2].SampledAt != started.Add(5*time.Second) {
		t.Fatalf("post-reset sample time = %s, want %s", operation.Samples[2].SampledAt, started.Add(5*time.Second))
	}
}

func TestCalculateSpeedUsesConfiguredTimeWindow(t *testing.T) {
	c := &Core{
		ctx: &domain.Context{
			Config: domain.Config{SpeedWindow: "90s"},
		},
	}
	now := time.Now()
	operation := &domain.Operation{
		Samples: []domain.SpeedSample{
			{Bytes: 100 * 1024 * 1024, SampledAt: now.Add(-2 * time.Minute)},
			{Bytes: 200 * 1024 * 1024, SampledAt: now.Add(-80 * time.Second)},
			{Bytes: 500 * 1024 * 1024, SampledAt: now.Add(-20 * time.Second)},
		},
	}

	speed := c.calculateSpeed(operation)
	if speed < 4.9 || speed > 5.1 {
		t.Fatalf("speed = %.2f MB/s, want recent-window speed near 5 MB/s", speed)
	}
}

func TestCalculateSpeedNeedsTwoWindowSamples(t *testing.T) {
	c := &Core{
		ctx: &domain.Context{
			Config: domain.Config{SpeedWindow: "90s"},
		},
	}
	now := time.Now()
	operation := &domain.Operation{
		Samples: []domain.SpeedSample{
			{Bytes: 100 * 1024 * 1024, SampledAt: now.Add(-2 * time.Minute)},
			{Bytes: 500 * 1024 * 1024, SampledAt: now.Add(-20 * time.Second)},
		},
	}

	if speed := c.calculateSpeed(operation); speed != 0 {
		t.Fatalf("speed = %.2f MB/s, want 0 with fewer than two in-window samples", speed)
	}
}

func TestRemainingAtSpeedUsesRecentSpeed(t *testing.T) {
	remaining := remainingAtSpeed(700*1024*1024, 100*1024*1024, 10)
	if remaining != "1m" {
		t.Fatalf("remaining = %s, want 1m", remaining)
	}
}

func TestRemainingAtSpeedRoundsUpAndFormatsForDisplay(t *testing.T) {
	remaining := remainingAtSpeed(10_700_000_000, 1_100_000_000, 62.01)
	if remaining != "2m 28s" {
		t.Fatalf("remaining = %s, want 2m 28s", remaining)
	}

	remaining = remainingAtSpeed(10_700_000_000, 7_360_000_000, 65.36)
	if remaining != "49s" {
		t.Fatalf("remaining = %s, want 49s", remaining)
	}
}

func TestRemainingAtSpeedUnknownWithoutSpeed(t *testing.T) {
	remaining := remainingAtSpeed(700*1024*1024, 100*1024*1024, 0)
	if remaining != "unknown" {
		t.Fatalf("remaining = %s, want unknown", remaining)
	}
}

func TestFormatRemainingDurationKeepsLargerUnitsReadable(t *testing.T) {
	remaining := formatRemainingDuration(2*time.Hour + 3*time.Minute + 4*time.Second + 20*time.Millisecond)
	if remaining != "2h 3m 5s" {
		t.Fatalf("remaining = %s, want 2h 3m 5s", remaining)
	}
}

func TestSpeedWindowParsesConfiguredDuration(t *testing.T) {
	c := &Core{
		ctx: &domain.Context{
			Config: domain.Config{SpeedWindow: "2m"},
		},
	}

	if window := c.speedWindow(); window != 2*time.Minute {
		t.Fatalf("speed window = %s, want 2m", window)
	}
}

func TestSpeedWindowTrimsWhitespace(t *testing.T) {
	c := &Core{
		ctx: &domain.Context{
			Config: domain.Config{SpeedWindow: " 90s "},
		},
	}

	if window := c.speedWindow(); window != defaultSpeedWindow {
		t.Fatalf("speed window = %s, want %s", window, defaultSpeedWindow)
	}
}

func TestSpeedWindowFallsBackOnMistypedDuration(t *testing.T) {
	c := &Core{
		ctx: &domain.Context{
			Config: domain.Config{SpeedWindow: "90sec"},
		},
	}

	if window := c.speedWindow(); window != defaultSpeedWindow {
		t.Fatalf("speed window = %s, want default %s", window, defaultSpeedWindow)
	}
}

func TestSpeedWindowFallsBackOnNonPositiveDuration(t *testing.T) {
	c := &Core{
		ctx: &domain.Context{
			Config: domain.Config{SpeedWindow: "0s"},
		},
	}

	if window := c.speedWindow(); window != defaultSpeedWindow {
		t.Fatalf("speed window = %s, want default %s", window, defaultSpeedWindow)
	}
}

func TestSpeedWindowCapsLargeDuration(t *testing.T) {
	c := &Core{
		ctx: &domain.Context{
			Config: domain.Config{SpeedWindow: "24h"},
		},
	}

	if window := c.speedWindow(); window != maxSpeedWindow {
		t.Fatalf("speed window = %s, want cap %s", window, maxSpeedWindow)
	}
}

func TestUpdateSamplesPrunesOutsideSpeedWindow(t *testing.T) {
	c := &Core{
		ctx: &domain.Context{
			Config: domain.Config{SpeedWindow: "3s"},
		},
	}
	started := time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC)
	operation := &domain.Operation{Started: started}

	c.updateSamplesAt(operation, 100, started.Add(1*time.Second))
	c.updateSamplesAt(operation, 200, started.Add(2*time.Second))
	c.updateSamplesAt(operation, 300, started.Add(6*time.Second))

	if len(operation.Samples) != 1 {
		t.Fatalf("samples = %d, want only the current-window sample", len(operation.Samples))
	}
	if operation.Samples[0].Bytes != 300 {
		t.Fatalf("remaining sample bytes = %d, want 300", operation.Samples[0].Bytes)
	}
}

func TestCommandCompletedSamplesFinalLifetimeProgress(t *testing.T) {
	c := &Core{
		ctx: &domain.Context{
			Config: domain.Config{SpeedWindow: "90s"},
			Hub:    pubsub.New(1),
		},
	}
	started := time.Now().Add(-time.Minute)
	sampleTime := time.Now().Add(-2 * time.Second)
	operation := &domain.Operation{
		Started:          started,
		BytesToTransfer:  1_000,
		BytesTransferred: 250,
		Samples: []domain.SpeedSample{
			{Bytes: 100, SampledAt: sampleTime},
			{Bytes: 250, SampledAt: sampleTime.Add(time.Second)},
		},
		SampleIndex:   2,
		PrevSample:    400,
		PrevSampleAt:  sampleTime.Add(time.Second),
		DeltaTransfer: 150,
		RsyncArgs:     []string{},
		Commands:      []*domain.Command{},
	}
	command := &domain.Command{Size: 250}

	c.commandCompleted(operation, command)

	if operation.BytesTransferred != 500 {
		t.Fatalf("bytes transferred = %d, want 500", operation.BytesTransferred)
	}
	if operation.PrevSample != operation.BytesTransferred {
		t.Fatalf("prev sample = %d, want lifetime baseline %d", operation.PrevSample, operation.BytesTransferred)
	}
	if operation.SampleIndex != 3 {
		t.Fatalf("sample index = %d, want three samples", operation.SampleIndex)
	}
	if operation.Samples[0].Bytes != 100 || operation.Samples[1].Bytes != 250 {
		t.Fatalf("samples were reset: %v", operation.Samples[:3])
	}
	if operation.Samples[2].Bytes != operation.BytesTransferred {
		t.Fatalf("final sample bytes = %d, want lifetime bytes %d", operation.Samples[2].Bytes, operation.BytesTransferred)
	}
	if !operation.Samples[2].SampledAt.After(operation.Samples[1].SampledAt) {
		t.Fatalf("final sample time = %s, want after %s", operation.Samples[2].SampledAt, operation.Samples[1].SampledAt)
	}
}
