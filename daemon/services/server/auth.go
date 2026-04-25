package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"unbalance/daemon/common"
	"unbalance/daemon/logger"
)

const (
	sessionCookieName = "unbalanced_session"
	sessionDuration   = 180 * 24 * time.Hour
	sessionRefreshAge = 24 * time.Hour
	maxLoginAttempts  = 5
	loginLockDuration = 5 * time.Minute
	pruneInterval     = 15 * time.Minute
)

type session struct {
	Username string
	CSRF     string
	Expires  time.Time
}

type loginAttempt struct {
	Count       int
	LastFailure time.Time
	LockedUntil time.Time
}

type sessionStore struct {
	Items map[string]session `json:"items"`
}

type authStatus struct {
	Enabled       bool   `json:"enabled"`
	Configured    bool   `json:"configured"`
	Authenticated bool   `json:"authenticated"`
	Username      string `json:"username"`
	CSRFToken     string `json:"csrfToken"`
}

func (s *Server) authConfigured() bool {
	return s.ctx.AuthPassword != ""
}

func (s *Server) authRequired() bool {
	return s.ctx.AuthEnabled
}

func (s *Server) authStatus(c echo.Context) error {
	info, ok := s.currentSession(c)
	username := s.ctx.AuthUsername
	csrfToken := ""
	if ok {
		username = info.Username
		csrfToken = info.CSRF
	}

	return c.JSON(http.StatusOK, authStatus{
		Enabled:       s.authRequired(),
		Configured:    s.authConfigured(),
		Authenticated: ok,
		Username:      username,
		CSRFToken:     csrfToken,
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

	clientKey := clientKey(c)
	if err := s.checkLoginThrottle(clientKey); err != nil {
		return err
	}

	if payload.Username != s.ctx.AuthUsername {
		s.recordLoginFailure(clientKey)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(s.ctx.AuthPassword), []byte(payload.Password)); err != nil {
		s.recordLoginFailure(clientKey)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	s.clearLoginFailures(clientKey)

	if err := s.createSession(c, payload.Username); err != nil {
		logger.Red("unable to create auth session: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to create session")
	}

	info, _ := s.currentSession(c)

	return c.JSON(http.StatusOK, authStatus{
		Enabled:       true,
		Configured:    true,
		Authenticated: true,
		Username:      payload.Username,
		CSRFToken:     info.CSRF,
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

	s.clearAllSessions()

	if err := s.createSession(c, payload.Username); err != nil {
		logger.Red("unable to create setup session: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to create session")
	}

	info, _ := s.currentSession(c)

	return c.JSON(http.StatusOK, authStatus{
		Enabled:       true,
		Configured:    true,
		Authenticated: true,
		Username:      payload.Username,
		CSRFToken:     info.CSRF,
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

		info, ok := s.currentSession(c)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
		}

		c.Set("session", info)

		return next(c)
	}
}

func (s *Server) requireCSRF(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !s.authRequired() {
			return next(c)
		}

		info, ok := s.currentSession(c)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
		}

		token := c.Request().Header.Get("X-CSRF-Token")
		if token == "" || token != info.CSRF {
			return echo.NewHTTPError(http.StatusForbidden, "invalid csrf token")
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

	info, ok := s.currentSession(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	token := c.QueryParam("csrf")
	if token == "" || token != info.CSRF {
		return echo.NewHTTPError(http.StatusForbidden, "invalid websocket csrf token")
	}

	return nil
}

func (s *Server) currentSession(c echo.Context) (session, bool) {
	if !s.authConfigured() {
		return session{}, false
	}

	cookie, err := c.Cookie(sessionCookieName)
	if err != nil {
		return session{}, false
	}

	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()

	info, ok := s.sessions[cookie.Value]
	if !ok {
		return session{}, false
	}

	if time.Now().After(info.Expires) {
		delete(s.sessions, cookie.Value)
		return session{}, false
	}

	lastRefresh := info.Expires.Add(-sessionDuration)
	if time.Since(lastRefresh) >= sessionRefreshAge {
		info.Expires = time.Now().Add(sessionDuration)
		s.sessions[cookie.Value] = info
		if err := s.saveSessionsLocked(); err != nil {
			logger.Red("unable to refresh persisted auth session: %s", err)
		}
	}

	return info, true
}

func (s *Server) createSession(c echo.Context, username string) error {
	id, err := randomToken(32)
	if err != nil {
		return err
	}

	csrfToken, err := randomToken(32)
	if err != nil {
		return err
	}

	expiry := time.Now().Add(sessionDuration)

	s.sessionMu.Lock()
	s.sessions[id] = session{
		Username: username,
		CSRF:     csrfToken,
		Expires:  expiry,
	}
	err = s.saveSessionsLocked()
	s.sessionMu.Unlock()
	if err != nil {
		return err
	}

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
		if saveErr := s.saveSessionsLocked(); saveErr != nil {
			logger.Red("unable to persist cleared auth session: %s", saveErr)
		}
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

func (s *Server) clearAllSessions() {
	s.sessionMu.Lock()
	s.sessions = newSessionStore()
	if err := s.saveSessionsLocked(); err != nil {
		logger.Red("unable to clear persisted auth sessions: %s", err)
	}
	s.sessionMu.Unlock()
}

func (s *Server) pruneSessions() {
	ticker := time.NewTicker(pruneInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		s.sessionMu.Lock()
		dirty := false
		for key, info := range s.sessions {
			if now.After(info.Expires) {
				delete(s.sessions, key)
				dirty = true
			}
		}
		if dirty {
			if err := s.saveSessionsLocked(); err != nil {
				logger.Red("unable to prune persisted auth sessions: %s", err)
			}
		}
		s.sessionMu.Unlock()

		s.limiterMu.Lock()
		for key, attempt := range s.limiter {
			if attempt.LockedUntil.IsZero() && now.Sub(attempt.LastFailure) > loginLockDuration {
				delete(s.limiter, key)
				continue
			}

			if !attempt.LockedUntil.IsZero() && now.After(attempt.LockedUntil) {
				delete(s.limiter, key)
			}
		}
		s.limiterMu.Unlock()
	}
}

func (s *Server) checkLoginThrottle(key string) error {
	s.limiterMu.Lock()
	defer s.limiterMu.Unlock()

	attempt, ok := s.limiter[key]
	if !ok {
		return nil
	}

	if !attempt.LockedUntil.IsZero() && time.Now().Before(attempt.LockedUntil) {
		seconds := int(time.Until(attempt.LockedUntil).Seconds())
		if seconds < 1 {
			seconds = 1
		}

		return echo.NewHTTPError(http.StatusTooManyRequests, fmt.Sprintf("too many login attempts, retry in %d seconds", seconds))
	}

	if !attempt.LockedUntil.IsZero() && time.Now().After(attempt.LockedUntil) {
		delete(s.limiter, key)
	}

	return nil
}

func (s *Server) recordLoginFailure(key string) {
	s.limiterMu.Lock()
	defer s.limiterMu.Unlock()

	attempt := s.limiter[key]
	attempt.Count++
	attempt.LastFailure = time.Now()
	if attempt.Count >= maxLoginAttempts {
		attempt.LockedUntil = time.Now().Add(loginLockDuration)
	}
	s.limiter[key] = attempt
}

func (s *Server) clearLoginFailures(key string) {
	s.limiterMu.Lock()
	delete(s.limiter, key)
	s.limiterMu.Unlock()
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

func newLoginLimiter() map[string]loginAttempt {
	return make(map[string]loginAttempt)
}

func clientKey(c echo.Context) string {
	ip := c.RealIP()
	if ip == "" {
		ip = c.Request().RemoteAddr
	}

	return ip
}

func (s *Server) sessionFile() string {
	return filepath.Join(common.PluginLocation, common.SessionFilename)
}

func (s *Server) loadSessions() error {
	location := s.sessionFile()
	file, err := os.Open(location)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}
	defer file.Close()

	var store sessionStore
	if err := json.NewDecoder(file).Decode(&store); err != nil {
		return err
	}

	if store.Items == nil {
		store.Items = newSessionStore()
	}

	s.sessionMu.Lock()
	s.sessions = store.Items
	s.sessionMu.Unlock()

	if !s.authConfigured() {
		s.clearAllSessions()
	}

	return nil
}

func (s *Server) saveSessionsLocked() error {
	location := s.sessionFile()
	tmpName := location + ".tmp"

	file, err := os.Create(tmpName)
	if err != nil {
		return err
	}

	store := sessionStore{
		Items: s.sessions,
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(store); err != nil {
		_ = file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, location)
}
