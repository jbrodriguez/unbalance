package core

import (
	"testing"
)

func TestClampLogLines(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{"unset falls back to default", 0, 100},
		{"negative falls back to default", -5, 100},
		{"below minimum is raised", 5, 10},
		{"within range is kept", 500, 500},
		{"above maximum is capped", 50000, 10000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clampLogLines(tt.value); got != tt.want {
				t.Fatalf("clampLogLines(%d) = %d, want %d", tt.value, got, tt.want)
			}
		})
	}
}
