package services

import (
	"authgo/db"
	"authgo/internal/types"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServices struct{}

func NewUserServices() *UserServices {
	return &UserServices{}
}

func (s *UserServices) RegisterByEmail(id uuid.UUID, email, avatarUrl, firstName, lastName, password string) (types.User, error) {

	var count int64
	if err := db.DB.Model(&types.UserEmail{}).Where("LOWER(email) = LOWER(?)", email).Count(&count).Error; err != nil {
		return types.User{}, fmt.Errorf("error checking existing user: %w", err)
	}

	if count > 0 {
		return types.User{}, fmt.Errorf("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(password)), bcrypt.DefaultCost)
	if err != nil {
		return types.User{}, fmt.Errorf("failed to hash password")
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

	emailId := uuid.New()

	user := types.User{
		ID:        &id,
		AvatarURL: avatarUrl,
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Password:  string(hashedPassword),
		Emails:    []types.UserEmail{{ID: &emailId, Email: strings.TrimSpace(strings.ToLower(email)), IsPrimary: true}},
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return types.User{}, fmt.Errorf("could not create user: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return types.User{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	return user, nil
}
