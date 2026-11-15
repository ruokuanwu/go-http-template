package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Address            string
	DatabasePath       string
	WebDir             string
	DEBUG              bool
	CookieSecure       bool
	SessionLifetime    time.Duration
	SessionIdleTimeout time.Duration
	TrustedOrigins     []string
}

func Load() (Config, error) {
	secure, err := strconv.ParseBool(value("COOKIE_SECURE", "true"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid COOKIE_SECURE: %w", err)
	}
	lifetime, err := time.ParseDuration(value("SESSION_LIFETIME", "168h"))
	if err != nil || lifetime <= 0 {
		return Config{}, fmt.Errorf("SESSION_LIFETIME must be a positive duration")
	}
	idle, err := time.ParseDuration(value("SESSION_IDLE_TIMEOUT", "24h"))
	if err != nil || idle <= 0 || idle > lifetime {
		return Config{}, fmt.Errorf("SESSION_IDLE_TIMEOUT must be positive and not exceed SESSION_LIFETIME")
	}
	var origins []string
	if raw := os.Getenv("AUTH_TRUSTED_ORIGINS"); raw != "" {
		for _, origin := range strings.Split(raw, ",") {
			origins = append(origins, strings.TrimSpace(origin))
		}
	}
	return Config{
		Address:            value("HTTP_ADDR", ":8080"),
		DatabasePath:       value("DB_PATH", "data/app.db"),
		WebDir:             value("WEB_DIR", "web/dist"),
		DEBUG:              value("DEBUG", "false") == "true",
		CookieSecure:       secure,
		SessionLifetime:    lifetime,
		SessionIdleTimeout: idle,
		TrustedOrigins:     origins,
	}, nil
}

func value(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
