package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type ApiConfig struct {
	ApiHost     string
	ApiPort     string
	FrontEndApi string
}

type SMTPConfig struct {
	SMTPHost       string
	SMTPPort       string
	SMTPSenderName string
	SMTPEmail      string
	SMTPPassword   string
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     string
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}

type RedisConfig struct {
	Address   string
	RedisPort string
	Password  string
	Db        int
}

type GoogleOAuthConfig struct {
	GoogleClientID         string
	GoogleClientSecret     string
	GoogleOAuthRedirectUrl string
}

type WhatsAppConfig struct {
	URL        string
	TemplateID string
	ApiKey     string
	ShortenUrl string
}

type EncryptionConfig struct {
	SecretKeyEncryption string
}

type Config struct {
	DbConfig
	ApiConfig
	SMTPConfig
	TokenConfig
	RedisConfig
	GoogleOAuthConfig
	WhatsAppConfig
	EncryptionConfig
}

func (c *Config) ReadConfigFile() error {
	// Nyalakan untuk local saja, kalau sudah di docker matikan

	env := os.Getenv("ENV")
	if env != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println(err)
			return errors.New("Failed to load .env file")
		}
	}

	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
	}

	c.ApiConfig = ApiConfig{
		ApiHost:     os.Getenv("API_HOST"),
		ApiPort:     os.Getenv("API_PORT"),
		FrontEndApi: os.Getenv("FRONT_END_API"),
	}

	if os.Getenv("PORT") != "" {
		c.ApiConfig.ApiPort = os.Getenv("PORT")
	}

	c.SMTPConfig = SMTPConfig{
		SMTPHost:       os.Getenv("SMTP_HOST"),
		SMTPPort:       os.Getenv("SMTP_PORT"),
		SMTPSenderName: os.Getenv("SMTP_SENDER"),
		SMTPEmail:      os.Getenv("SMTP_EMAIL"),
		SMTPPassword:   os.Getenv("SMTP_PASS"),
	}

	c.TokenConfig = TokenConfig{
		ApplicationName:     "CalibrationSystem",
		JwtSignatureKey:     "x/A?D(G+KaPdSgVkYp3s6v9y$B&E)H@M",
		JwtSigningMethod:    jwt.SigningMethodHS256,
		AccessTokenLifeTime: time.Hour * 2,
	}

	c.RedisConfig = RedisConfig{
		Address:   os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		Db:        0,
	}

	c.GoogleOAuthConfig.GoogleClientID = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	c.GoogleOAuthConfig.GoogleClientSecret = os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	c.GoogleOAuthConfig.GoogleOAuthRedirectUrl = os.Getenv("GOOGLE_OAUTH_REDIRECT_URL")

	c.WhatsAppConfig = WhatsAppConfig{
		URL:        os.Getenv("WA_API_URL"),
		TemplateID: os.Getenv("WA_TEMPLATE_ID"),
		ApiKey:     os.Getenv("WA_API_KEY"),
		ShortenUrl: os.Getenv("WA_SHORTEN_URL"),
	}

	c.EncryptionConfig = EncryptionConfig{
		SecretKeyEncryption: os.Getenv("SECRET_KEY_ENCRYPTION"),
	}

	if c.SMTPEmail == "" || c.SMTPHost == "" || c.SMTPPassword == "" || c.SMTPPort == "" || c.SMTPSenderName == "" ||
		c.DbConfig.Host == "" || c.DbConfig.Name == "" || c.DbConfig.Password == "" || c.DbConfig.Port == "" || c.DbConfig.User == "" {
		return errors.New("Missing required field")
	}
	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cfg.ReadConfigFile()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
