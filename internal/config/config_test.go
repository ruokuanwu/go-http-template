package config

import "testing"

func TestSessionConfiguration(t *testing.T) {
	for _, test := range []struct{ name, key, value string }{
		{"invalid secure flag", "COOKIE_SECURE", "invalid"},
		{"zero lifetime", "SESSION_LIFETIME", "0s"},
		{"invalid lifetime", "SESSION_LIFETIME", "tomorrow"},
		{"negative idle timeout", "SESSION_IDLE_TIMEOUT", "-1h"},
		{"idle exceeds lifetime", "SESSION_IDLE_TIMEOUT", "999h"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("COOKIE_SECURE", "true")
			t.Setenv("SESSION_LIFETIME", "168h")
			t.Setenv("SESSION_IDLE_TIMEOUT", "24h")
			t.Setenv(test.key, test.value)
			if _, err := Load(); err == nil {
				t.Fatal("expected invalid configuration to fail")
			}
		})
	}
}
