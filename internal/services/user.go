package services

import (
	"authgo/db"
	"authgo/internal/types"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServices struct{}

func NewUserServices() *UserServices {
	return &UserServices{}
}

func (s *UserServices) RegisterByEmail(id uuid.UUID, email, avatarUrl, firstName, lastName, password string) (types.User, error) {

	if email == "" || avatarUrl == "" || firstName == "" || lastName == "" || password == "" {
		return types.User{}, fmt.Errorf("some fields are empty")
	}

	var existingUser types.User
	if err := db.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return types.User{}, fmt.Errorf("email already in use")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return types.User{}, fmt.Errorf("error checking existing user: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return types.User{}, fmt.Errorf("failed to hash password")
	}

	user := types.User{
		ID:        &id,
		Email:     email,
		AvatarURL: avatarUrl,
		FirstName: firstName,
		LastName:  lastName,
		Password:  string(hashedPassword),
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		return types.User{}, fmt.Errorf("could not begin transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return types.User{}, fmt.Errorf("could not create user")
	}

	if err := tx.Commit().Error; err != nil {
		return types.User{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	return user, nil
}
