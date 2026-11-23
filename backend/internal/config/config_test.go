package config

import "testing"

func TestLoadConfigRejectsUnsafeProductionDatabase(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("JWT_SECRET", "a-production-secret-with-more-than-32-characters")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("REQUIRE_EMAIL_VERIFICATION", "false")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected production config with postgres user to be rejected")
	}
}

func TestLoadConfigRequiresSMTPForProductionVerification(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("JWT_SECRET", "a-production-secret-with-more-than-32-characters")
	t.Setenv("DB_USER", "ecommerce_app")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("REQUIRE_EMAIL_VERIFICATION", "true")
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SMTP_USERNAME", "")
	t.Setenv("SMTP_PASSWORD", "")
	t.Setenv("SMTP_FROM", "")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected production config without SMTP settings to be rejected")
	}
}

func setValidProductionEnv(t *testing.T) {
	t.Helper()
	t.Setenv("APP_ENV", "production")
	t.Setenv("CLIENT_BASE_URL", "https://bookverse.example")
	t.Setenv("JWT_SECRET", "a-production-secret-with-more-than-32-characters")
	t.Setenv("DB_USER", "ecommerce_app")
	t.Setenv("DB_PASSWORD", "strong-prod-db-password-2026")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("REDIS_PASSWORD", "strong-prod-redis-password-2026")
	t.Setenv("MIDTRANS_SERVER_KEY", "Mid-server-prod-valid-key-example")
	t.Setenv("MIDTRANS_CLIENT_KEY", "Mid-client-prod-valid-key-example")
	t.Setenv("REQUIRE_EMAIL_VERIFICATION", "false")
}

func TestLoadConfigRejectsPlaceholderMidtransKeysInProduction(t *testing.T) {
	setValidProductionEnv(t)
	t.Setenv("MIDTRANS_SERVER_KEY", "SB-Mid-server-TEST-KEY-SAMPLE")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected production config with placeholder MIDTRANS_SERVER_KEY to be rejected")
	}

	setValidProductionEnv(t)
	t.Setenv("MIDTRANS_CLIENT_KEY", "SB-Mid-client-TEST-KEY-SAMPLE")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected production config with placeholder MIDTRANS_CLIENT_KEY to be rejected")
	}
}

func TestLoadConfigRejectsDefaultDBAndEmptyRedisPasswordInProduction(t *testing.T) {
	setValidProductionEnv(t)
	t.Setenv("DB_PASSWORD", "postgres")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected production config with default DB_PASSWORD to be rejected")
	}

	setValidProductionEnv(t)
	t.Setenv("REDIS_PASSWORD", "")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected production config with empty REDIS_PASSWORD to be rejected")
	}
}

func TestLoadConfigAcceptsValidProductionSecrets(t *testing.T) {
	setValidProductionEnv(t)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected valid production config to pass, got error: %v", err)
	}
	if cfg.Midtrans.ServerKey != "Mid-server-prod-valid-key-example" {
		t.Fatalf("unexpected MIDTRANS_SERVER_KEY: %s", cfg.Midtrans.ServerKey)
	}
}

func TestLoadConfigRejectsInsecureCookieInProduction(t *testing.T) {
	setValidProductionEnv(t)
	t.Setenv("COOKIE_SECURE", "false")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected production config with COOKIE_SECURE=false to be rejected")
	}
}

func TestLoadConfigRejectsNonHTTPSClientURLInProduction(t *testing.T) {
	setValidProductionEnv(t)
	t.Setenv("CLIENT_BASE_URL", "http://bookverse.example")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected production config with http:// CLIENT_BASE_URL to be rejected")
	}
}
