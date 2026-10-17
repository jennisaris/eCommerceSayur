package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	AppPort            string
	Database           DatabaseConfig
	Redis              RedisConfig
	JWT                JWTConfig
	Midtrans           MidtransConfig
	Client             ClientConfig
	CORSAllowedOrigins string
	RunSeeder          bool
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
	Secret          string
	ExpirationHours int
}

type MidtransConfig struct {
	ServerKey    string
	ClientKey    string
	IsProduction bool
}

type ClientConfig struct {
	BaseURL string
}

func LoadConfig() (*Config, error) {
	// Try loading .env file (ignore error if not found in production)
	_ = godotenv.Load()

	appEnv := getEnv("APP_ENV", "development")
	jwtHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "72"))
	if err != nil {
		jwtHours = 72
	}

	jwtSecret := getEnv("JWT_SECRET", "super-secret-ecommerce-jwt-token-key-2026")
	if appEnv == "production" {
		if jwtSecret == "super-secret-ecommerce-jwt-token-key-2026" || len(jwtSecret) < 32 {
			return nil, fmt.Errorf("FATAL SECURITY: Insecure default or short JWT_SECRET detected in production mode. Set a random secret of at least 32 characters in .env")
		}
	}

	isProduction := getEnv("MIDTRANS_IS_PRODUCTION", "false") == "true"

	// Seeder configuration: enabled by default in development, disabled by default in production
	runSeeder := appEnv != "production"
	if customSeeder, ok := os.LookupEnv("RUN_SEEDER"); ok {
		runSeeder = customSeeder == "true" || customSeeder == "1"
	}

	clientBaseURL := getEnv("CLIENT_BASE_URL", "http://localhost:3000")
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "")
	if corsOrigins == "" {
		if appEnv == "production" {
			corsOrigins = clientBaseURL
		} else {
			corsOrigins = fmt.Sprintf("%s,http://localhost:3000,http://127.0.0.1:3000", clientBaseURL)
		}
	}

	return &Config{
		AppEnv:  appEnv,
		AppPort: getEnv("APP_PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "ecommerce_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret:          jwtSecret,
			ExpirationHours: jwtHours,
		},
		Midtrans: MidtransConfig{
			ServerKey:    getEnv("MIDTRANS_SERVER_KEY", "SB-Mid-server-TEST-KEY-SAMPLE"),
			ClientKey:    getEnv("MIDTRANS_CLIENT_KEY", "SB-Mid-client-TEST-KEY-SAMPLE"),
			IsProduction: isProduction,
		},
		Client: ClientConfig{
			BaseURL: clientBaseURL,
		},
		CORSAllowedOrigins: corsOrigins,
		RunSeeder:          runSeeder,
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
