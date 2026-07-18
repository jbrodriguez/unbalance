package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func capture(t *testing.T, fn func()) string {
	t.Helper()

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	fn()

	return buf.String()
}

func TestRedWritesErrorTag(t *testing.T) {
	out := capture(t, func() { Red("boom: %s", "disk1") })

	if !strings.Contains(out, "[error] boom: disk1") {
		t.Fatalf("expected error tag, got %q", out)
	}
}

func TestYellowWritesWarnTag(t *testing.T) {
	out := capture(t, func() { Yellow("careful") })

	if !strings.Contains(out, "[warn] careful") {
		t.Fatalf("expected warn tag, got %q", out)
	}
}

func TestBlueWritesInfoTag(t *testing.T) {
	out := capture(t, func() { Blue("hello") })

	if !strings.Contains(out, "[info] hello") {
		t.Fatalf("expected info tag, got %q", out)
	}
}
