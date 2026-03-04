package config

import "os"

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret string

	APIPort   string
	UploadDir string
}

func LoadConfig() *Config {
	config := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASS"),
		DBName:     os.Getenv("DB_NAME"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		APIPort:   os.Getenv("API_PORT"),
		UploadDir: os.Getenv("UPLOAD_DIR"),
	}

	validateConfig(config)

	return config
}

func validateConfig(config *Config) {
	if config.DBHost == "" {
		panic("DB_HOST is required")
	}
	if config.DBPort == "" {
		panic("DB_PORT is required")
	}
	if config.DBUser == "" {
		panic("DB_USER is required")
	}
	if config.DBPassword == "" {
		panic("DB_PASSWORD is required")
	}
	if config.DBName == "" {
		panic("DB_NAME is required")
	}
	if config.JWTSecret == "" {
		panic("JWT_SECRET is required")
	}
	if config.APIPort == "" {
		panic("API_PORT is required")
	}
	if config.UploadDir == "" {
		panic("UPLOAD_DIR is required")
	}
}
