package config

import "os"

type Config struct {
	Address      string
	DatabasePath string
	WebDir       string
}

func Load() Config {
	return Config{
		Address:      value("HTTP_ADDR", ":8080"),
		DatabasePath: value("DB_PATH", "data/app.db"),
		WebDir:       value("WEB_DIR", "web/dist"),
	}
}

func value(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
