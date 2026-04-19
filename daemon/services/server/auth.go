package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"unbalance/daemon/logger"
)

const (
	sessionCookieName = "unbalanced_session"
	sessionDuration   = 24 * time.Hour
)

type session struct {
	Username string
	Expires  time.Time
}

type authStatus struct {
	Enabled       bool   `json:"enabled"`
	Configured    bool   `json:"configured"`
	Authenticated bool   `json:"authenticated"`
	Username      string `json:"username"`
}

func (s *Server) authConfigured() bool {
	return s.ctx.AuthPassword != ""
}

func (s *Server) authRequired() bool {
	return s.ctx.AuthEnabled
}

func (s *Server) authStatus(c echo.Context) error {
	authenticated, username := s.currentSession(c)

	return c.JSON(http.StatusOK, authStatus{
		Enabled:       s.authRequired(),
		Configured:    s.authConfigured(),
		Authenticated: authenticated,
		Username:      username,
	})
}

func (s *Server) login(c echo.Context) error {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid login payload")
	}

	if !s.authRequired() {
		return c.JSON(http.StatusOK, authStatus{
			Enabled:       false,
			Configured:    true,
			Authenticated: true,
			Username:      s.ctx.AuthUsername,
		})
	}

	if !s.authConfigured() {
		return echo.NewHTTPError(http.StatusConflict, "authentication setup is incomplete")
	}

	if payload.Username != s.ctx.AuthUsername {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(s.ctx.AuthPassword), []byte(payload.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	if err := s.createSession(c, payload.Username); err != nil {
		logger.Red("unable to create auth session: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to create session")
	}

	return c.JSON(http.StatusOK, authStatus{
		Enabled:       true,
		Configured:    true,
		Authenticated: true,
		Username:      payload.Username,
	})
}

func (s *Server) logout(c echo.Context) error {
	s.clearSession(c)

	return c.JSON(http.StatusOK, authStatus{
		Enabled:       s.authRequired(),
		Configured:    s.authConfigured(),
		Authenticated: false,
		Username:      s.ctx.AuthUsername,
	})
}

func (s *Server) setup(c echo.Context) error {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid setup payload")
	}

	if !s.authRequired() {
		return echo.NewHTTPError(http.StatusBadRequest, "authentication is disabled")
	}

	if s.authConfigured() {
		return echo.NewHTTPError(http.StatusConflict, "authentication is already configured")
	}

	if payload.Username == "" {
		payload.Username = s.ctx.AuthUsername
	}

	if payload.Username == "" {
		payload.Username = "admin"
	}

	if len(payload.Password) < 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Red("unable to hash auth password: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to save password")
	}

	if err := s.core.SetAuth(payload.Username, string(hash)); err != nil {
		logger.Red("unable to persist auth credentials: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to persist password")
	}

	if err := s.createSession(c, payload.Username); err != nil {
		logger.Red("unable to create setup session: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to create session")
	}

	return c.JSON(http.StatusOK, authStatus{
		Enabled:       true,
		Configured:    true,
		Authenticated: true,
		Username:      payload.Username,
	})
}

func (s *Server) requireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !s.authRequired() {
			return next(c)
		}

		if !s.authConfigured() {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication setup is incomplete")
		}

		ok, _ := s.currentSession(c)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
		}

		return next(c)
	}
}

func (s *Server) validateWebsocketRequest(c echo.Context) error {
	if !s.authRequired() {
		return nil
	}

	if !s.authConfigured() {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication setup is incomplete")
	}

	allowedOrigin := "http://" + c.Request().Host
	if c.IsTLS() {
		allowedOrigin = "https://" + c.Request().Host
	}

	origin := c.Request().Header.Get("Origin")
	if origin == "" {
		return echo.NewHTTPError(http.StatusForbidden, "missing websocket origin")
	}

	if origin != allowedOrigin {
		return echo.NewHTTPError(http.StatusForbidden, "invalid websocket origin")
	}

	ok, _ := s.currentSession(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	return nil
}

func (s *Server) currentSession(c echo.Context) (bool, string) {
	cookie, err := c.Cookie(sessionCookieName)
	if err != nil {
		return false, ""
	}

	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()

	info, ok := s.sessions[cookie.Value]
	if !ok {
		return false, ""
	}

	if time.Now().After(info.Expires) {
		delete(s.sessions, cookie.Value)
		return false, ""
	}

	return true, info.Username
}

func (s *Server) createSession(c echo.Context, username string) error {
	id, err := randomToken(32)
	if err != nil {
		return err
	}

	expiry := time.Now().Add(sessionDuration)

	s.sessionMu.Lock()
	s.sessions[id] = session{
		Username: username,
		Expires:  expiry,
	}
	s.sessionMu.Unlock()

	c.SetCookie(&http.Cookie{
		Name:     sessionCookieName,
		Value:    id,
		Path:     "/",
		Expires:  expiry,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   c.IsTLS(),
	})

	return nil
}

func (s *Server) clearSession(c echo.Context) {
	cookie, err := c.Cookie(sessionCookieName)
	if err == nil {
		s.sessionMu.Lock()
		delete(s.sessions, cookie.Value)
		s.sessionMu.Unlock()
	}

	c.SetCookie(&http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   c.IsTLS(),
	})
}

func randomToken(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("invalid token length")
	}

	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

func newSessionStore() map[string]session {
	return make(map[string]session)
}
