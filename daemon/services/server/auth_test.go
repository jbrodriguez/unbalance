package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"

	"unbalance/daemon/domain"
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
		{"forwarded host accepted", "https://public.example", "internal:7090", true},
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
			if tc.name == "forwarded host accepted" {
				req.Header.Set("X-Forwarded-Host", "public.example")
			}
			got := upgrader.CheckOrigin(req)
			if got != tc.want {
				t.Fatalf("CheckOrigin(origin=%q,host=%q) = %v, want %v", tc.origin, tc.host, got, tc.want)
			}
		})
	}
}

func TestValidateWebsocketRequestAllowsForwardedHTTPSOrigin(t *testing.T) {
	s := &Server{
		ctx: &domain.Context{
			Config: domain.Config{
				AuthEnabled:  true,
				AuthUsername: "admin",
				AuthPassword: "configured",
			},
		},
		sessions: newSessionStore(),
	}
	s.sessions["sid"] = session{Username: "admin", CSRF: "token", Expires: time.Now().Add(time.Hour)}

	req := httptest.NewRequest(http.MethodGet, "/ws?csrf=token", nil)
	req.Host = "internal:7090"
	req.Header.Set("Origin", "https://public.example")
	req.Header.Set("X-Forwarded-Host", "public.example")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "sid"})
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	if err := s.validateWebsocketRequest(c); err != nil {
		t.Fatalf("validateWebsocketRequest() returned error: %v", err)
	}
}

func TestRequestExternalSecureHonorsForwardedProto(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.Header.Set("X-Forwarded-Proto", "https")

	if !requestExternalSecure(req) {
		t.Fatal("requestExternalSecure() did not honor X-Forwarded-Proto=https")
	}
}
