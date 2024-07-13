package db

import (
	"authgo/internal/config"
	"authgo/internal/types"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() error {
	cfg := config.GetConfig()

	host := cfg.Db.Host
	port := cfg.Db.Port
	name := cfg.Db.Name
	user := cfg.Db.User
	password := cfg.Db.Password
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, name, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connetc to the Database: %v", err)
	}
	log.Println("🚀 Connected Successfully to the Database")
	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(&types.User{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database migrated successfully")
	return nil
}
