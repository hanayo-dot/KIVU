package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds configuration parameters for the KIVU backend service.
type Config struct {
	ServerPort             string
	Environment            string
	DatabaseURL            string
	DBHost                 string
	DBPort                 string
	DBUser                 string
	DBPassword             string
	DBName                 string
	DBSSLMode              string
	JWTSecret              string
	JWTAccessExpiryMinutes int
	JWTRefreshExpiryDays   int
	AIServiceURL           string
	CopernicusClientID     string
	CopernicusClientSecret string
	UseCopernicusLive      bool
}

// LoadConfig reads values from environment variables with fallback defaults.
func LoadConfig() *Config {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "kivu_db")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	defaultDBURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	useCopernicus, _ := strconv.ParseBool(getEnv("USE_COPERNICUS_LIVE", "false"))
	accessExp, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXP_MINUTES", "30"))
	refreshExp, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXP_DAYS", "30"))

	return &Config{
		ServerPort:             getEnv("SERVER_PORT", "8080"),
		Environment:            getEnv("APP_ENV", "development"),
		DatabaseURL:            getEnv("DATABASE_URL", defaultDBURL),
		DBHost:                 dbHost,
		DBPort:                 dbPort,
		DBUser:                 dbUser,
		DBPassword:             dbPassword,
		DBName:                 dbName,
		DBSSLMode:              dbSSLMode,
		JWTSecret:              getEnv("JWT_SECRET", "kivu-hackathon-super-secret-jwt-key-2026"),
		JWTAccessExpiryMinutes: accessExp,
		JWTRefreshExpiryDays:   refreshExp,
		AIServiceURL:           getEnv("AI_SERVICE_URL", "http://localhost:8000"),
		CopernicusClientID:     getEnv("COPERNICUS_CLIENT_ID", ""),
		CopernicusClientSecret: getEnv("COPERNICUS_CLIENT_SECRET", ""),
		UseCopernicusLive:      useCopernicus,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
