package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	AppEnv           string
	AppPort          string
	DatabaseDSN      string
	JWTSecret        string
	JWTIssuer        string
	JWTAudience      string
	CORSAllowOrigins []string
}

func Load() Config {
	loadDotEnv(".env")
	return Config{
		AppEnv:           getenv("APP_ENV", "development"),
		AppPort:          getenv("APP_PORT", "8082"),
		DatabaseDSN:      getenv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=pos_001 port=5432 sslmode=disable TimeZone=Asia/Bangkok"),
		JWTSecret:        getenv("JWT_SECRET", "change-me-use-32-plus-chars"),
		JWTIssuer:        getenv("JWT_ISSUER", "https://global-commerce.com"),
		JWTAudience:      getenv("JWT_AUDIENCE", "pos_module"),
		CORSAllowOrigins: splitCSV(getenv("CORS_ALLOW_ORIGINS", "http://localhost:3000,http://localhost:8082")),
	}
}

func getenv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		pair := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(pair[0])
		value := strings.Trim(strings.TrimSpace(pair[1]), "\"")
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}
