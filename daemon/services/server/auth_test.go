package server

import (
	"net/http"
	"testing"
	"time"
)

func TestSessionStillValid(t *testing.T) {
	s := &Server{
		sessions: newSessionStore(),
	}

	now := time.Now()
	s.sessions["good"] = session{Username: "u", CSRF: "c", Expires: now.Add(1 * time.Hour)}
	s.sessions["expired"] = session{Username: "u", CSRF: "c", Expires: now.Add(-1 * time.Minute)}

	cases := []struct {
		name string
		id   string
		want bool
	}{
		{"valid", "good", true},
		{"expired", "expired", false},
		{"missing", "nope", false},
		{"empty", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := s.sessionStillValid(tc.id)
			if got != tc.want {
				t.Fatalf("sessionStillValid(%q) = %v, want %v", tc.id, got, tc.want)
			}
		})
	}
}

func TestUpgraderCheckOrigin(t *testing.T) {
	cases := []struct {
		name   string
		origin string
		host   string
		want   bool
	}{
		{"empty origin rejected", "", "host.example", false},
		{"matching host accepted", "https://host.example", "host.example", true},
		{"matching host with port", "http://host.example:8080", "host.example:8080", true},
		{"scheme difference ignored (proxy)", "http://host.example", "host.example", true},
		{"different host rejected", "https://attacker.example", "host.example", false},
		{"malformed origin rejected", "://broken", "host.example", false},
		{"opaque origin rejected", "null", "host.example", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := &http.Request{
				Host:   tc.host,
				Header: http.Header{},
			}
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			got := upgrader.CheckOrigin(req)
			if got != tc.want {
				t.Fatalf("CheckOrigin(origin=%q,host=%q) = %v, want %v", tc.origin, tc.host, got, tc.want)
			}
		})
	}
}
