package config

import (
	"os"
	"strconv"
)

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

	IDPClientID     string
	IDPClientSecret string
	IDPRedirectURI  string
	IDPLoginURL     string
	IDPLogoutURL    string
	IDPTokenURL     string
	IDPUserinfoURL  string
	IDPRefreshURL   string
	IDPSessionURL   string

	RedisHost string
	RedisPort string
	RedisPass string
	RedisDB   int

	GotenbergURL string
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

		LocalUploadDIR: os.Getenv("UPLOAD_DIR"),
		AzureStorageConnectionString: os.Getenv(
			"AZURE_STORAGE_CONNECTION_STRING",
		),
		AzureContainerName: os.Getenv("AZURE_CONTAINER_NAME"),

		IsProduction: os.Getenv("IS_PRODUCTION") == "true",

		IDPClientID:     os.Getenv("IDP_CLIENT_ID"),
		IDPClientSecret: os.Getenv("IDP_CLIENT_SECRET"),
		IDPRedirectURI:  os.Getenv("IDP_REDIRECT_URI"),
		IDPLoginURL:     os.Getenv("IDP_LOGIN_ENDPOINT"),
		IDPLogoutURL:    os.Getenv("IDP_LOGOUT_ENDPOINT"),
		IDPTokenURL:     os.Getenv("IDP_TOKEN_ENDPOINT"),
		IDPUserinfoURL:  os.Getenv("IDP_USERINFO_ENDPOINT"),
		IDPRefreshURL:   os.Getenv("IDP_REFRESH_ENDPOINT"),
		IDPSessionURL:   os.Getenv("IDP_SESSION_ENDPOINT"),

		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),
		RedisPass: os.Getenv("REDIS_PASS"),
		RedisDB: func() int {
			db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
			if err != nil {
				return 0
			}
			return db
		}(),

		GotenbergURL: os.Getenv("GOTENBERG_URL"),
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
			panic(
				"AZURE_STORAGE_CONNECTION_STRING is required for Azure Blob Storage",
			)
		}
		if config.AzureContainerName == "" {
			panic("AZURE_CONTAINER_NAME is required for Azure Blob Storage")
		}
	}
	if config.IDPClientID == "" {
		panic("IDP_CLIENT_ID is required")
	}
	if config.IDPClientSecret == "" {
		panic("IDP_CLIENT_SECRET is required")
	}
	if config.IDPRedirectURI == "" {
		panic("IDP_REDIRECT_URI is required")
	}
	if config.IDPLoginURL == "" {
		panic("IDP_AUTHORIZE_ENDPOINT is required")
	}
	if config.IDPTokenURL == "" {
		panic("IDP_TOKEN_ENDPOINT is required")
	}
	if config.IDPUserinfoURL == "" {
		panic("IDP_USERINFO_ENDPOINT is required")
	}
	if config.IDPRefreshURL == "" {
		panic("IDP_REFRESH_ENDPOINT is required")
	}
	if config.IDPSessionURL == "" {
		panic("IDP_SESSION_ENDPOINT is required")
	}

	if config.RedisHost == "" {
		panic("REDIS_HOST is required")
	}
	if config.RedisPort == "" {
		panic("REDIS_PORT is required")
	}
}
