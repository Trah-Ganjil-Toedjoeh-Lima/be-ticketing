package config

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

type VaultConfig struct {
	VaultHost  string
	VaultPort  string
	VaultAuth  string
	VaultToken string
	VaultPath  string
}

type EmailConfig struct {
	MailMailer      string
	MailHost        string
	MailPort        string
	MailUsername    string
	MailPassword    string
	MailEncryption  string
	MailFromAddress string
	MailFromName    string
}

type AppConfig struct {
	AppName        string
	IsProduction   bool
	AppUrl         string
	AppPort        string
	EndpointPrefix string

	DBHost                    string
	DBUser                    string
	DBPassword                string
	DBName                    string
	DBPort                    string
	DBMaxIdleConnection       int
	DBMaxOpenConnection       int
	DBConnectionMaxLifeMinute time.Duration

	RedisPassword string
	RedisHost     string
	RedisPort     string

	MinioHost          string
	MinioPort          string
	MinioLocation      string
	MinioRootUser      string
	MinioRootPassword  string
	MinioSecure        bool
	MinioTicketsBucket string

	AccessSecret  string
	RefreshSecret string
	AccessMinute  time.Duration
	RefreshMinute time.Duration

	MerchId              string
	MidtransIsProduction bool
	ClientKeySandbox     string
	ServerKeySandbox     string
	ClientKeyProduction  string
	ServerKeyProduction  string

	Email1 EmailConfig
	Email2 EmailConfig
	Email3 EmailConfig

	TransactionMinute time.Duration
	TotpPeriod        uint

	AdminName  string
	AdminEmail string
	AdminPhone string

	IsOpenGate      bool
	IsOpenAuth      bool
	QrScanBehaviour string
}

func NewAppConfig() *AppConfig {
	midtransIsProduction, _ := strconv.ParseBool(getEnv("MIDTRANS_IS_PRODUCTION", "0"))
	isProduction, _ := strconv.ParseBool(getEnv("IS_PRODUCTION", "0"))

	isOpenGate, _ := strconv.ParseBool(getEnv("IS_OPEN_GATE", "true"))
	isOpenAuth, _ := strconv.ParseBool(getEnv("IS_OPEN_AUTH", "true"))

	email1 := EmailConfig{
		getEnv("MAIL_MAILER_1", "smtp"),
		getEnv("MAIL_HOST_1", "smtp.gmail.com"),
		getEnv("MAIL_PORT_1", "465"),
		getEnv("MAIL_USERNAME_1", ""),
		getEnv("MAIL_PASSWORD_1", ""),
		getEnv("MAIL_ENCRYPTION_1", "ssl"),
		getEnv("MAIL_FROM_ADDRESS_1", ""),
		getEnv("MAIL_FROM_NAME_1", "gmco"),
	}

	email2 := EmailConfig{
		getEnv("MAIL_MAILER_2", "smtp"),
		getEnv("MAIL_HOST_2", "smtp.gmail.com"),
		getEnv("MAIL_PORT_2", "465"),
		getEnv("MAIL_USERNAME_2", ""),
		getEnv("MAIL_PASSWORD_2", ""),
		getEnv("MAIL_ENCRYPTION_2", "ssl"),
		getEnv("MAIL_FROM_ADDRESS_2", ""),
		getEnv("MAIL_FROM_NAME_2", "gmco"),
	}

	email3 := EmailConfig{
		getEnv("MAIL_MAILER_3", "smtp"),
		getEnv("MAIL_HOST_3", "smtp.gmail.com"),
		getEnv("MAIL_PORT_3", "465"),
		getEnv("MAIL_USERNAME_3", ""),
		getEnv("MAIL_PASSWORD_3", ""),
		getEnv("MAIL_ENCRYPTION_3", "ssl"),
		getEnv("MAIL_FROM_ADDRESS_3", ""),
		getEnv("MAIL_FROM_NAME_3", "gmco"),
	}

	minioSecure, _ := strconv.ParseBool(getEnv("MINIO_SECURE", "false"))

	accessMinute, _ := time.ParseDuration(getEnv("ACCESS_MINUTE", "15m"))
	refreshMinute, _ := time.ParseDuration(getEnv("REFRESH_MINUTE", "120m"))
	transactionMinute, _ := time.ParseDuration(getEnv("TRANSACTION_MINUTE", "15m"))
	totpPeriodSecond, _ := time.ParseDuration(getEnv("TOTP_PERIOD", "5m"))

	dbMaxIdleConnection, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNECTION", "10"))
	dbMaxOpenConnection, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNECTION", "10"))
	dbConnectionMaxLifeMinute, _ := time.ParseDuration(getEnv("DB_CONNECTION_MAX_LIFE_MINUTE", "60m"))

	var appConfig = AppConfig{
		AppName:        getEnv("APP_NAME", "gmcgo"),
		IsProduction:   isProduction,
		AppUrl:         getEnv("APP_URL", "127.0.0.1"),
		AppPort:        getEnv("APP_PORT", "8080"),
		EndpointPrefix: getEnv("ENDPOINT_PREFIX", "/api/v1/"),

		DBHost:                    getEnv("DB_HOST", "localhost"),
		DBUser:                    getEnv("DB_USER", "root"),
		DBPassword:                getEnv("DB_PASSWORD", "root"),
		DBName:                    getEnv("DB_NAME", "gmcgo"),
		DBPort:                    getEnv("DB_PORT", "5432"),
		DBMaxIdleConnection:       dbMaxIdleConnection,
		DBMaxOpenConnection:       dbMaxOpenConnection,
		DBConnectionMaxLifeMinute: dbConnectionMaxLifeMinute,

		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisHost:     getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),

		MinioHost:          getEnv("MINIO_HOST", ""),
		MinioPort:          getEnv("MINIO_PORT", ""),
		MinioLocation:      getEnv("MINIO_LOCATION", ""),
		MinioRootUser:      getEnv("MINIO_ROOT_USER", ""),
		MinioRootPassword:  getEnv("MINIO_ROOT_PASSWORD", ""),
		MinioSecure:        minioSecure,
		MinioTicketsBucket: getEnv("MINIO_TICKETS_BUCKET", ""),

		AccessSecret:  getEnv("ACCESS_SECRET", ""),
		RefreshSecret: getEnv("REFRESH_SECRET", ""),
		AccessMinute:  accessMinute,
		RefreshMinute: refreshMinute,

		MerchId:              getEnv("MERCH_ID", ""),
		MidtransIsProduction: midtransIsProduction,
		ClientKeySandbox:     getEnv("CLIENT_KEY_SANDBOX", ""),
		ServerKeySandbox:     getEnv("SERVER_KEY_SANDBOX", ""),
		ClientKeyProduction:  getEnv("CLIENT_KEY_PRODUCTION", ""),
		ServerKeyProduction:  getEnv("SERVER_KEY_PRODUCTION", ""),

		Email1: email1,
		Email2: email2,
		Email3: email3,

		TransactionMinute: transactionMinute,
		TotpPeriod:        uint(totpPeriodSecond.Seconds()),

		AdminName:  getEnv("ADMIN_NAME", ""),
		AdminEmail: getEnv("ADMIN_EMAIL", ""),
		AdminPhone: getEnv("ADMIN_PHONE", ""),

		IsOpenGate:      isOpenGate,
		IsOpenAuth:      isOpenAuth,
		QrScanBehaviour: "default", //open_gate, ticket_exchanging, default
	}
	return &appConfig
}

func getEnv(key string, fallback string) string {

	chVault := make(chan string)
	chdotEnv := make(chan string)
	go getVaultEnv(key, chVault)
	go getdotEnv(key, chdotEnv)
	vaultVal := <-chVault
	fallbackdotEnv := <-chdotEnv
	//close(chVault)
	//close(chdotEnv)

	fallbackosEnv := os.Getenv(key)

	var fallbackVal string
	if fallbackdotEnv != "error_failed_to_get_key" {
		fallbackVal = fallbackdotEnv
	} else if fallbackosEnv != "" {
		fallbackVal = fallbackosEnv
	} else {
		fallbackVal = fallback
	}

	if vaultVal != "error_failed_to_get_key" {
		log.Printf("getEnv: Vault key %s found", key)
		return vaultVal
	} else if fallbackVal != "" {
		log.Printf("getEnv: Vault key %s using fallback", key)
		return fallbackVal
	} else {
		log.Fatalf("getEnv: Can't set key %s", key)
		return "error_failed_to_get_key"
	}
}

func getVaultEnv(key string, ch chan string) int {
	vaultConfig := VaultConfig{
		os.Getenv("VAULT_HOST"),
		os.Getenv("VAULT_PORT"),
		os.Getenv("VAULT_AUTH"),
		os.Getenv("VAULT_TOKEN"),
		os.Getenv("VAULT_PATH"),
	}

	var vaultURL string
	if vaultConfig.VaultHost == "" || vaultConfig.VaultPort == "" || vaultConfig.VaultAuth == "" || vaultConfig.VaultToken == "" || vaultConfig.VaultPath == "" {
		log.Printf("getVault: Vault config not found")
		ch <- "error_failed_to_get_key"
		return 1
	} else {
		vaultURL = "http://" + vaultConfig.VaultHost + ":" + vaultConfig.VaultPort
		start := time.Now()

		resp, err := http.Get(vaultURL)
		if err != nil {
			log.Printf("getVault: unable to connect to Vault: %v", err)
			ch <- "error_failed_to_get_key"
			return 1
		}

		defer resp.Body.Close()
		elapsed := time.Since(start)

		log.Printf("getVault: HTTP GET request to %s took %s\n", vaultURL, elapsed)
	}

	config := vault.DefaultConfig()

	config.Address = vaultURL

	client, err := vault.NewClient(config)
	if err != nil {
		log.Printf("getVault: unable to initialize Vault client: %v", err)
		ch <- "error_failed_to_get_key"
		return 1
	}

	// Authenticate
	client.SetToken(vaultConfig.VaultToken)

	// Read a secret from the default mount path for KV v2 in dev mode, "secret"
	secret, err := client.KVv2("kv").Get(context.Background(), vaultConfig.VaultPath)
	if err != nil {
		log.Printf("getVault: unable to read secret %v", err)
		ch <- "error_failed_to_get_key"
		return 1
	}

	value, ok := secret.Data[key].(string)
	if !ok {
		log.Printf("getVault: failed to get key %s", key)
		ch <- "error_failed_to_get_key"
		return 1
	} else {
		log.Printf("getVault: Success get env %s from vault", key)
		ch <- value
		return 0
	}
}

func getdotEnv(key string, chdotEnv chan string) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("getdotEnv: Error loading .env file")
	}

	if valuedotenv, ok := os.LookupEnv(key); ok {
		log.Printf("getdotEnv: Success get env %s from .env", key)
		chdotEnv <- valuedotenv
	} else {
		log.Printf("getdotEnv: Failed to get env %s from .env", key)
		chdotEnv <- "error_failed_to_get_key"
	}
}
