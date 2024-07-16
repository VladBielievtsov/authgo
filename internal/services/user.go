package services

import (
	"authgo/db"
	"authgo/internal/config"
	"authgo/internal/types"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServices struct{}

func NewUserServices() *UserServices {
	return &UserServices{}
}

func (s *UserServices) RegisterByEmail(id uuid.UUID, email, avatarUrl, firstName, lastName, password string) (types.User, error) {
	cfg := config.GetConfig()
	var mailServices = NewMailServices(cfg)

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
	confirmationToken := uuid.New()

	user := types.User{
		ID:        &id,
		AvatarURL: avatarUrl,
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Password:  string(hashedPassword),
		Emails: []types.UserEmail{{
			ID:                &emailId,
			Email:             strings.TrimSpace(strings.ToLower(email)),
			IsPrimary:         true,
			IsConfirmed:       false,
			ConfirmationToken: &confirmationToken,
		}},
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return types.User{}, fmt.Errorf("could not create user: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return types.User{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	auth := mailServices.New()

	if err := mailServices.Send(
		"Subject: AuthGo - Email Confirmation\r\nContent-Type: text/html; charset=utf-8\r\n\r\n<html><body><p>Confirm Email: <a href='http://"+cfg.Application.Domain+"/confirm/"+confirmationToken.String()+"'>Confirmation Link</a></p></body></html>\r\n",
		email,
		auth,
	); err != nil {
		return types.User{}, fmt.Errorf("failed to send confirmation email: %w", err)
	}

	return user, nil
}

func (s *UserServices) LoginByEmail(email, password string) (types.LoginResponce, error) {
	var user types.User
	cfg := config.GetConfig()

	err := db.DB.Joins("JOIN user_emails ON user_emails.user_id = users.id").
		Where("user_emails.email = ? AND user_emails.is_primary = ?", email, true).
		Preload("Emails").
		First(&user).Error
	if err != nil {
		return types.LoginResponce{}, fmt.Errorf("could not find user: %v", err)
	}

	for _, userEmail := range user.Emails {
		if userEmail.IsPrimary && !userEmail.IsConfirmed {
			return types.LoginResponce{}, fmt.Errorf("confirm your email, verification link has been sent to your email")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return types.LoginResponce{}, fmt.Errorf("invalid password")
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = user.ID
	claims["exp"] = now.Add(120 * time.Minute).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(cfg.Application.JwtSecter))
	if err != nil {
		return types.LoginResponce{}, fmt.Errorf("generating JWT Token failed")
	}

	return types.LoginResponce{
		User:  user,
		Token: tokenString,
	}, nil
}

func (s *UserServices) GetAllUsers() ([]types.User, error) {
	var users []types.User

	result := db.DB.Preload("Emails").Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users: %v", result.Error)
	}

	return users, nil
}
