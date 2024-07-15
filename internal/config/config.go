package config

import (
	"authgo/internal/utils"
	"fmt"
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
	Identity string
	Username string
	Password string
	Host     string
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
			Identity: "",
			Username: os.Getenv("SMTPFrom"),
			Password: os.Getenv("SMTPPassword"),
			Host:     "smtp.gmail.com",
		},
	}

	return cfg, nil
}

func GetConfig() *Config {
	return cfg
}
