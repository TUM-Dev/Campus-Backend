package admin

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Enabled          bool
	BaseURL          string
	ClientID         string
	ClientSecret     string
	Issuer           string
	AuthorizationURL string
	TokenURL         string
	JWKSURL          string
	AllowedSubjects  []string
}

func ConfigFromEnv() (Config, error) {
	cfg := Config{
		BaseURL:          os.Getenv("ADMIN_BASE_URL"),
		ClientID:         os.Getenv("ADMIN_OIDC_CLIENT_ID"),
		ClientSecret:     os.Getenv("ADMIN_OIDC_CLIENT_SECRET"),
		Issuer:           envOr("ADMIN_OIDC_ISSUER", "https://tumidp.lrz.de/idp/shibboleth"),
		AuthorizationURL: envOr("ADMIN_OIDC_AUTH_URL", "https://login.tum.de/idp/profile/oidc/authorize"),
		TokenURL:         envOr("ADMIN_OIDC_TOKEN_URL", "https://login.tum.de/idp/profile/oidc/token"),
		JWKSURL:          envOr("ADMIN_OIDC_JWKS_URL", "https://login.tum.de/idp/profile/oidc/keyset"),
	}
	if value := os.Getenv("ADMIN_ENABLED"); value != "" {
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return Config{}, fmt.Errorf("invalid ADMIN_ENABLED: %w", err)
		}
		cfg.Enabled = enabled
	}
	for subject := range strings.SplitSeq(os.Getenv("ADMIN_ALLOWED_SUBJECTS"), ",") {
		if subject = strings.TrimSpace(subject); subject != "" {
			cfg.AllowedSubjects = append(cfg.AllowedSubjects, subject)
		}
	}
	return cfg, cfg.validate()
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (c Config) validate() error {
	if !c.Enabled {
		return nil
	}
	if c.ClientID == "" || c.ClientSecret == "" || len(c.AllowedSubjects) == 0 {
		return fmt.Errorf("admin requires an OIDC client ID, client secret and explicit subject allowlist")
	}
	for _, subject := range c.AllowedSubjects {
		if subject == "" || len(subject) > 255 {
			return fmt.Errorf("invalid admin allowlisted subject")
		}
	}
	for name, value := range map[string]string{
		"ADMIN_BASE_URL": c.BaseURL, "ADMIN_OIDC_ISSUER": c.Issuer,
		"ADMIN_OIDC_AUTH_URL": c.AuthorizationURL, "ADMIN_OIDC_TOKEN_URL": c.TokenURL,
		"ADMIN_OIDC_JWKS_URL": c.JWKSURL,
	} {
		u, err := url.Parse(value)
		if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
			return fmt.Errorf("%s must be an absolute HTTPS URL without credentials, query or fragment", name)
		}
		if name == "ADMIN_BASE_URL" && (u.Path != "" || u.RawPath != "") {
			return fmt.Errorf("ADMIN_BASE_URL must be an origin without a path or trailing slash")
		}
	}
	if len(c.Issuer) > 255 {
		return fmt.Errorf("admin OIDC issuer is too long")
	}
	return nil
}
