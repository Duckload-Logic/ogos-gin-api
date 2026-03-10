package config

import "os"

type Config struct {
	DBHost  string
	DBPort  string
	DBUser  string
	DBPass  string
	DBName  string
	DBTLS   bool
	DBTLSCA string

	JWTSecret string

	WebsitesPort string

	LocalUploadDIR               string
	AzureStorageConnectionString string
	AzureContainerName           string

	IsProduction bool
}

func LoadConfig() *Config {
	config := &Config{
		DBHost:  os.Getenv("DB_HOST"),
		DBPort:  os.Getenv("DB_PORT"),
		DBUser:  os.Getenv("DB_USER"),
		DBPass:  os.Getenv("DB_PASS"),
		DBName:  os.Getenv("DB_NAME"),
		DBTLS:   os.Getenv("DB_TLS") == "true",
		DBTLSCA: os.Getenv("DB_TLS_CA"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		WebsitesPort: os.Getenv("WEBSITES_PORT"),

		LocalUploadDIR:               os.Getenv("UPLOAD_DIR"),
		AzureStorageConnectionString: os.Getenv("AZURE_STORAGE_CONNECTION_STRING"),
		AzureContainerName:           os.Getenv("AZURE_CONTAINER_NAME"),

		IsProduction: os.Getenv("IS_PRODUCTION") == "true",
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
	if config.DBPass == "" {
		panic("DB_PASSWORD is required")
	}
	if config.DBName == "" {
		panic("DB_NAME is required")
	}
	if config.JWTSecret == "" {
		panic("JWT_SECRET is required")
	}
	if config.WebsitesPort == "" {
		panic("WEBSITES_PORT is required")
	}
	if config.LocalUploadDIR == "" && !config.IsProduction {
		panic("UPLOAD_DIR is required for local storage")
	}
	if config.IsProduction {
		if config.AzureStorageConnectionString == "" {
			panic("AZURE_STORAGE_CONNECTION_STRING is required for Azure Blob Storage")
		}
		if config.AzureContainerName == "" {
			panic("AZURE_CONTAINER_NAME is required for Azure Blob Storage")
		}
	}
}
