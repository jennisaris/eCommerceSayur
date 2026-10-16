package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv   string
	AppPort  string
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Midtrans MidtransConfig
	Client   ClientConfig
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

	jwtHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "72"))
	if err != nil {
		jwtHours = 72
	}

	isProduction := getEnv("MIDTRANS_IS_PRODUCTION", "false") == "true"

	return &Config{
		AppEnv:  getEnv("APP_ENV", "development"),
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
			Secret:          getEnv("JWT_SECRET", "super-secret-ecommerce-jwt-token-key-2026"),
			ExpirationHours: jwtHours,
		},
		Midtrans: MidtransConfig{
			ServerKey:    getEnv("MIDTRANS_SERVER_KEY", "SB-Mid-server-TEST-KEY-SAMPLE"),
			ClientKey:    getEnv("MIDTRANS_CLIENT_KEY", "SB-Mid-client-TEST-KEY-SAMPLE"),
			IsProduction: isProduction,
		},
		Client: ClientConfig{
			BaseURL: getEnv("CLIENT_BASE_URL", "http://localhost:3000"),
		},
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
