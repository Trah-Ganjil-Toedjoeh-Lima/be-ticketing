package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

//TODO: make this a singleton => use injection strategy?

type AppConfig struct {
	AppName       string
	IsProduction  string
	AppUrl        string
	AppPort       string
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	APISecret     string
	TokenDuration string
	RedisPassword string
	RedisHost     string
	RedisPort     string
	AccessSecret  string
	RefreshSecret string
	AccessMinute  time.Duration
	RefreshMinute time.Duration
}

func NewAppConfig() *AppConfig {
	var appConfig = AppConfig{
		AppName:       getEnv("APP_NAME", "gmcgo"),
		IsProduction:  getEnv("IS_PRODUCTION", "false"),
		AppUrl:        getEnv("APP_URL", "127.0.0.1"),
		AppPort:       getEnv("APP_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", "root"),
		DBName:        getEnv("DB_NAME", "gmcgo"),
		DBPort:        getEnv("DB_PORT", "5432"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisHost:     getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		AccessSecret:  getEnv("ACCESS_MINUTE", "15"),
		RefreshSecret: getEnv("SECRET_MINUTE", "240"),
		AccessMinute:  15,
		RefreshMinute: 60,
	}
	return &appConfig

}

func getEnv(key, fallback string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error on loading .env file")
	}
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
