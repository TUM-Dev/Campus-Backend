package admin

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/coreos/go-oidc/v3/oidc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

//go:embed templates/*.html style.css
var assets embed.FS

const (
	sessionCookie    = "__Host-campus-admin"
	loginCookie      = "__Host-campus-admin-login"
	loginLifetime    = 10 * time.Minute
	sessionLifetime  = time.Hour
	maxPendingLogins = 1024
)

type loginFlow struct {
	nonce, verifier string
	expires         time.Time
}

type server struct {
	db        *gorm.DB
	cfg       Config
	cron      *cron.CronService
	version   string
	started   time.Time
	oauth     oauth2.Config
	verifier  *oidc.IDTokenVerifier
	client    *http.Client
	templates *template.Template
	mu        sync.Mutex
	logins    map[string]loginFlow
}

// New returns the complete admin handler. Disabled configurations serve only 404s.
func New(db *gorm.DB, jobs *cron.CronService, version string, cfg Config) (http.Handler, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return http.NotFoundHandler(), nil
	}
	return newHandler(db, jobs, version, cfg, &http.Client{Timeout: 10 * time.Second})
}

func newHandler(db *gorm.DB, jobs *cron.CronService, version string, cfg Config, client *http.Client) (http.Handler, error) {
	templates, err := template.New("admin").Funcs(template.FuncMap{"interval": intervalLabel}).ParseFS(assets, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse admin templates: %w", err)
	}
	keyContext := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	s := &server{
		db: db, cfg: cfg, cron: jobs, version: version, started: time.Now(),
		client: client, templates: templates, logins: make(map[string]loginFlow),
		oauth: oauth2.Config{
			ClientID: cfg.ClientID, ClientSecret: cfg.ClientSecret,
			RedirectURL: cfg.BaseURL + "/admin/callback", Scopes: []string{oidc.ScopeOpenID},
			Endpoint: oauth2.Endpoint{
				AuthURL: cfg.AuthorizationURL, TokenURL: cfg.TokenURL,
				AuthStyle: oauth2.AuthStyleInHeader,
			},
		},
		verifier: oidc.NewVerifier(cfg.Issuer, oidc.NewRemoteKeySet(keyContext, cfg.JWKSURL), &oidc.Config{
			ClientID: cfg.ClientID, SupportedSigningAlgs: []string{oidc.RS256},
		}),
	}
	return s.handler(), nil
}

func (s *server) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/login", s.login)
	mux.HandleFunc("GET /admin/callback", s.callback)
	mux.HandleFunc("GET /admin/signed-out", func(w http.ResponseWriter, r *http.Request) {
		s.render(w, "signed-out", nil)
	})
	mux.HandleFunc("GET /admin/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		body, err := assets.ReadFile("style.css")
		if err != nil {
			s.internalError(w, err)
			return
		}
		_, _ = w.Write(body)
	})
	mux.HandleFunc("GET /admin/{$}", s.protected(s.dashboard))
	mux.HandleFunc("GET /admin/feeds", s.protected(s.feeds))
	mux.HandleFunc("POST /admin/feeds", s.protected(s.saveFeed))
	mux.HandleFunc("GET /admin/groups", s.protected(s.groups))
	mux.HandleFunc("GET /admin/councils", s.protected(s.councils))
	mux.HandleFunc("GET /admin/ratings", s.protected(s.ratings))
	mux.HandleFunc("GET /admin/audit", s.protected(s.audit))
	mux.HandleFunc("POST /admin/logout", s.protected(s.logout))
	mux.Handle("GET /admin", http.RedirectHandler("/admin/", http.StatusSeeOther))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; img-src 'self'; form-action 'self'; frame-ancestors 'none'; base-uri 'none'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "same-origin")
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		mux.ServeHTTP(w, r.WithContext(ctx))
	})
}

func secret() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func digest(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func setCookie(w http.ResponseWriter, name, value string, lifetime time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name: name, Value: value, Path: "/", Secure: true, HttpOnly: true,
		SameSite: http.SameSiteLaxMode, MaxAge: int(lifetime.Seconds()),
	})
}

func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name: name, Path: "/", Secure: true, HttpOnly: true,
		SameSite: http.SameSiteLaxMode, MaxAge: -1,
	})
}

func (s *server) discardExpiredLogins(now time.Time) {
	for state, flow := range s.logins {
		if !now.Before(flow.expires) {
			delete(s.logins, state)
		}
	}
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	state := secret()
	flow := loginFlow{nonce: secret(), verifier: oauth2.GenerateVerifier(), expires: time.Now().Add(loginLifetime)}
	s.mu.Lock()
	s.discardExpiredLogins(time.Now())
	if old, err := r.Cookie(loginCookie); err == nil {
		delete(s.logins, old.Value)
	}
	if len(s.logins) >= maxPendingLogins {
		s.mu.Unlock()
		log.Warn("admin pending login limit reached")
		http.Error(w, "Too many pending logins; try again later.", http.StatusTooManyRequests)
		return
	}
	s.logins[state] = flow
	s.mu.Unlock()
	setCookie(w, loginCookie, state, loginLifetime)
	http.Redirect(w, r, s.oauth.AuthCodeURL(state, oidc.Nonce(flow.nonce), oauth2.S256ChallengeOption(flow.verifier)), http.StatusSeeOther)
}

func (s *server) callback(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(loginCookie)
	state := r.URL.Query().Get("state")
	if err != nil || len(state) != 64 || subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(state)) != 1 {
		http.Error(w, "Invalid login state.", http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	flow, ok := s.logins[state]
	delete(s.logins, state)
	s.discardExpiredLogins(time.Now())
	s.mu.Unlock()
	clearCookie(w, loginCookie)
	if !ok || !time.Now().Before(flow.expires) || r.URL.Query().Get("error") != "" || r.URL.Query().Get("code") == "" {
		http.Error(w, "Login expired or was declined. Start a new login.", http.StatusBadRequest)
		return
	}
	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, s.client)
	token, err := s.oauth.Exchange(ctx, r.URL.Query().Get("code"), oauth2.VerifierOption(flow.verifier))
	if err != nil {
		log.Warn("admin OIDC code exchange failed")
		http.Error(w, "SSO login failed.", http.StatusBadGateway)
		return
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Warn("admin OIDC response did not contain an ID token")
		http.Error(w, "SSO login failed.", http.StatusBadGateway)
		return
	}
	idToken, err := s.verifier.Verify(ctx, rawIDToken)
	if err != nil || subtle.ConstantTimeCompare([]byte(idToken.Nonce), []byte(flow.nonce)) != 1 {
		log.Warn("admin OIDC ID token validation failed")
		http.Error(w, "SSO login failed.", http.StatusUnauthorized)
		return
	}
	var claims struct {
		AuthorizedParty string `json:"azp"`
	}
	if err := idToken.Claims(&claims); err != nil || (len(idToken.Audience) > 1 && claims.AuthorizedParty != s.cfg.ClientID) ||
		(claims.AuthorizedParty != "" && claims.AuthorizedParty != s.cfg.ClientID) {
		log.Warn("admin OIDC authorized party validation failed")
		http.Error(w, "SSO login failed.", http.StatusUnauthorized)
		return
	}
	if !s.allowed(idToken.Subject) {
		log.Warn("admin login denied by subject allowlist")
		http.Error(w, "This account does not have admin access.", http.StatusForbidden)
		return
	}
	sessionToken := secret()
	expires := minTime(time.Now().Add(sessionLifetime), idToken.Expiry)
	session := model.AdminSession{
		TokenHash: digest(sessionToken), Issuer: s.cfg.Issuer, Subject: idToken.Subject,
		CSRFToken: secret(), ExpiresAt: expires.UTC(),
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if old, err := r.Cookie(sessionCookie); err == nil {
			if err := tx.Delete(&model.AdminSession{}, "token_hash = ?", digest(old.Value)).Error; err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.AdminSession{}, "expires_at <= ?", time.Now().UTC()).Error; err != nil {
			return err
		}
		return tx.Create(&session).Error
	})
	if err != nil {
		s.internalError(w, err)
		return
	}
	setCookie(w, sessionCookie, sessionToken, time.Until(expires))
	http.Redirect(w, r, "/admin/", http.StatusSeeOther)
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func (s *server) allowed(subject string) bool {
	if subject == "" {
		return false
	}
	for _, allowed := range s.cfg.AllowedSubjects {
		if subject == allowed {
			return true
		}
	}
	return false
}

type adminHandler func(http.ResponseWriter, *http.Request, model.AdminSession)

func (s *server) protected(next adminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookie)
		if err != nil || len(cookie.Value) != 64 {
			s.unauthenticated(w, r)
			return
		}
		var session model.AdminSession
		err = s.db.WithContext(r.Context()).Where("token_hash = ? AND expires_at > ?", digest(cookie.Value), time.Now().UTC()).First(&session).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.unauthenticated(w, r)
			return
		}
		if err != nil {
			s.internalError(w, err)
			return
		}
		if session.Issuer != s.cfg.Issuer || !s.allowed(session.Subject) {
			http.Error(w, "Admin access has been revoked.", http.StatusForbidden)
			return
		}
		if r.Method == http.MethodPost {
			r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid or oversized form.", http.StatusBadRequest)
				return
			}
			if origin := r.Header.Get("Origin"); origin != "" && origin != s.cfg.BaseURL {
				http.Error(w, "Invalid request origin.", http.StatusForbidden)
				return
			}
			if subtle.ConstantTimeCompare([]byte(r.PostForm.Get("csrf")), []byte(session.CSRFToken)) != 1 {
				http.Error(w, "Invalid CSRF token.", http.StatusForbidden)
				return
			}
		}
		next(w, r, session)
	}
}

func (s *server) unauthenticated(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, sessionCookie)
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	http.Error(w, "Admin login required.", http.StatusUnauthorized)
}

func (s *server) logout(w http.ResponseWriter, r *http.Request, session model.AdminSession) {
	if err := s.db.WithContext(r.Context()).Delete(&model.AdminSession{}, "token_hash = ?", session.TokenHash).Error; err != nil {
		s.internalError(w, err)
		return
	}
	clearCookie(w, sessionCookie)
	http.Redirect(w, r, "/admin/signed-out", http.StatusSeeOther)
}

func (s *server) internalError(w http.ResponseWriter, err error) {
	log.WithError(err).Error("admin request failed")
	http.Error(w, "Admin request failed; see server logs.", http.StatusInternalServerError)
}
