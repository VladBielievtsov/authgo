package main

import (
	"authgo/db"
	"authgo/internal/config"
	"authgo/internal/types"
	"authgo/internal/utils"
	"fmt"

	"gorm.io/gorm"
)

func main() {
	_, err := config.New()
	if err != nil {
		utils.ErrorHandler(err)
	}

	if err := db.ConnectDatabase(); err != nil {
		utils.ErrorHandler(err)
	}

	clearTable(db.DB, &types.User{})

	fmt.Println("Database cleared successfully.")
}

func clearTable(db *gorm.DB, model interface{}) {
	if err := db.Unscoped().Where("1 = 1").Delete(model).Error; err != nil {
		panic(err)
	}
}
