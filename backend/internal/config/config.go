package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                   string
	AppPort                  string
	Database                 DatabaseConfig
	Redis                    RedisConfig
	JWT                      JWTConfig
	Midtrans                 MidtransConfig
	Client                   ClientConfig
	Email                    EmailConfig
	RequireEmailVerification bool
	CORSAllowedOrigins       string
	RunSeeder                bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

type JWTConfig struct {
	Secret                  string
	AccessExpirationMinutes int
	RefreshTokenDays        int
	CookieSecure            bool
}

type MidtransConfig struct {
	ServerKey    string
	ClientKey    string
	IsProduction bool
}

type ClientConfig struct {
	BaseURL string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	From         string
}

func LoadConfig() (*Config, error) {
	// Try loading .env file (ignore error if not found in production)
	_ = godotenv.Load()

	appEnv := getEnv("APP_ENV", "development")
	jwtMinutes, err := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRATION_MINUTES", "15"))
	if err != nil {
		jwtMinutes = 15
	}
	if jwtMinutes < 5 {
		jwtMinutes = 5
	}
	refreshTokenDays, err := strconv.Atoi(getEnv("REFRESH_TOKEN_DAYS", "30"))
	if err != nil || refreshTokenDays < 1 {
		refreshTokenDays = 30
	}

	jwtSecret := getEnv("JWT_SECRET", "super-secret-ecommerce-jwt-token-key-2026")
	if appEnv == "production" {
		if jwtSecret == "super-secret-ecommerce-jwt-token-key-2026" || isPlaceholderSecret(jwtSecret) || len(jwtSecret) < 32 {
			return nil, fmt.Errorf("FATAL SECURITY: Insecure default or short JWT_SECRET detected in production mode. Set a random secret of at least 32 characters in .env")
		}
	}

	isProduction := getEnv("MIDTRANS_IS_PRODUCTION", "false") == "true"
	cookieSecure := appEnv == "production"
	if configuredCookieSecure, ok := os.LookupEnv("COOKIE_SECURE"); ok {
		cookieSecure = configuredCookieSecure == "true" || configuredCookieSecure == "1"
	}
	if appEnv == "production" && !cookieSecure {
		return nil, fmt.Errorf("FATAL SECURITY: production requires COOKIE_SECURE=true")
	}

	// Seeder configuration: enabled by default in development, disabled by default in production
	runSeeder := appEnv != "production"
	if customSeeder, ok := os.LookupEnv("RUN_SEEDER"); ok {
		runSeeder = customSeeder == "true" || customSeeder == "1"
	}

	clientBaseURL := getEnv("CLIENT_BASE_URL", "http://localhost:3000")
	if appEnv == "production" && !strings.HasPrefix(strings.ToLower(strings.TrimSpace(clientBaseURL)), "https://") {
		return nil, fmt.Errorf("FATAL SECURITY: production requires CLIENT_BASE_URL to use https://")
	}
	requireEmailVerification := appEnv == "production"
	if configuredVerification, ok := os.LookupEnv("REQUIRE_EMAIL_VERIFICATION"); ok {
		requireEmailVerification = configuredVerification == "true" || configuredVerification == "1"
	}
	emailConfig := EmailConfig{
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		From:         getEnv("SMTP_FROM", ""),
	}
	if requireEmailVerification && (emailConfig.SMTPHost == "" || emailConfig.SMTPUsername == "" || emailConfig.SMTPPassword == "" || emailConfig.From == "") {
		return nil, fmt.Errorf("FATAL SECURITY: production email verification requires SMTP_HOST, SMTP_USERNAME, SMTP_PASSWORD, and SMTP_FROM")
	}
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "")
	if corsOrigins == "" {
		if appEnv == "production" {
			corsOrigins = clientBaseURL
		} else {
			corsOrigins = fmt.Sprintf("%s,http://localhost:3000,http://127.0.0.1:3000", clientBaseURL)
		}
	}
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	midtransServerKey := getEnv("MIDTRANS_SERVER_KEY", "SB-Mid-server-TEST-KEY-SAMPLE")
	midtransClientKey := getEnv("MIDTRANS_CLIENT_KEY", "SB-Mid-client-TEST-KEY-SAMPLE")

	if appEnv == "production" {
		if dbUser == "postgres" {
			return nil, fmt.Errorf("FATAL SECURITY: production must use a dedicated least-privilege DB_USER")
		}
		if dbPassword == "" || dbPassword == "postgres" || isPlaceholderSecret(dbPassword) {
			return nil, fmt.Errorf("FATAL SECURITY: production requires a strong non-default DB_PASSWORD")
		}
		if dbSSLMode == "disable" {
			return nil, fmt.Errorf("FATAL SECURITY: production requires DB_SSLMODE to be require, verify-ca, or verify-full")
		}
		if redisPassword == "" || isPlaceholderSecret(redisPassword) {
			return nil, fmt.Errorf("FATAL SECURITY: production requires a strong REDIS_PASSWORD")
		}
		if midtransServerKey == "" || isPlaceholderSecret(midtransServerKey) {
			return nil, fmt.Errorf("FATAL SECURITY: production requires a valid MIDTRANS_SERVER_KEY")
		}
		if midtransClientKey == "" || isPlaceholderSecret(midtransClientKey) {
			return nil, fmt.Errorf("FATAL SECURITY: production requires a valid MIDTRANS_CLIENT_KEY")
		}
	}

	return &Config{
		AppEnv:  appEnv,
		AppPort: getEnv("APP_PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     dbUser,
			Password: dbPassword,
			Name:     getEnv("DB_NAME", "ecommerce_db"),
			SSLMode:  dbSSLMode,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: redisPassword,
		},
		JWT: JWTConfig{
			Secret:                  jwtSecret,
			AccessExpirationMinutes: jwtMinutes,
			RefreshTokenDays:        refreshTokenDays,
			CookieSecure:            cookieSecure,
		},
		Midtrans: MidtransConfig{
			ServerKey:    midtransServerKey,
			ClientKey:    midtransClientKey,
			IsProduction: isProduction,
		},
		Client: ClientConfig{
			BaseURL: clientBaseURL,
		},
		Email:                    emailConfig,
		RequireEmailVerification: requireEmailVerification,
		CORSAllowedOrigins:       corsOrigins,
		RunSeeder:                runSeeder,
	}, nil
}

func isPlaceholderSecret(val string) bool {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return true
	}
	upper := strings.ToUpper(trimmed)
	placeholders := []string{
		"TEST-KEY-SAMPLE",
		"YOUR_SERVER_KEY",
		"YOUR_CLIENT_KEY",
		"YOUR_STRONG_",
		"CHANGE_ME",
	}
	for _, p := range placeholders {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
