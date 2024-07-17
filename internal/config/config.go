package config

import (
	"authgo/internal/utils"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Application applocationConf
	Db          DbEnv
	Mail        mailConf
}

type applocationConf struct {
	Port      string
	Domain    string
	JwtSecter string
}

type mailConf struct {
	Username string
	Password string
	Host     string
	Port     string
}

type DbEnv struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

var cfg *Config

func New() (*Config, error) {
	err := godotenv.Load(".env.local")
	utils.ErrorHandler(err)

	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")

	if dbName == "" || dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" {
		return nil, fmt.Errorf("some environment variables are missing")
	}

	cfg = &Config{
		Db: DbEnv{
			Host:     dbHost,
			Port:     dbPort,
			Name:     dbName,
			User:     dbUser,
			Password: dbPassword,
		},
		Application: applocationConf{
			Port:      os.Getenv("APP_PORT"),
			Domain:    os.Getenv("APP_DOMAIN"),
			JwtSecter: os.Getenv("JWT_SECTER"),
		},
		Mail: mailConf{
			Username: os.Getenv("SMTPFrom"),
			Password: os.Getenv("SMTPPassword"),
			Host:     os.Getenv("SMTPHost"),
			Port:     os.Getenv("SMTPPort"),
		},
	}

	return cfg, nil
}

func GetConfig() *Config {
	if cfg == nil {
		log.Panic("Config not initialized. Call New() before GetConfig()")
	}
	return cfg
}
