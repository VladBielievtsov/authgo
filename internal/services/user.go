package services

import (
	"authgo/db"
	"authgo/internal/config"
	"authgo/internal/types"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServices struct {
	cfg          *config.Config
	mailServices *MailServices
}

func NewUserServices(cfg *config.Config) *UserServices {
	return &UserServices{
		cfg:          cfg,
		mailServices: NewMailServices(cfg),
	}
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
	confirmationToken := rand.Intn(900000) + 100000

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

	if s.mailServices == nil {
		return types.User{}, fmt.Errorf("mail services not initialized")
	}

	auth := s.mailServices.New()

	if err := s.mailServices.Send(
		"Subject: AuthGo - Email Confirmation\r\nContent-Type: text/html; charset=utf-8\r\n\r\n<html><body><p>Verification code: "+strconv.Itoa(confirmationToken)+"</p></body></html>\r\n",
		email,
		auth,
	); err != nil {
		return types.User{}, fmt.Errorf("failed to send confirmation email: %w", err)
	}

	return user, nil
}

func (s *UserServices) LoginByEmail(email, password string) (types.LoginResponce, error) {
	var user types.User

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

	tokenString, err := tokenByte.SignedString([]byte(s.cfg.Application.JwtSecter))
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

func (s *UserServices) SendConfirmCode(email string) (string, error) {
	var user types.User

	err := db.DB.Joins("JOIN user_emails ON user_emails.user_id = users.id").
		Where("user_emails.email = ?", email).
		Preload("Emails").
		First(&user).Error
	if err != nil {
		return "", fmt.Errorf("could not find user: %v", err)
	}

	for _, userEmail := range user.Emails {
		if userEmail.Email == email {
			if userEmail.IsConfirmed {
				return "", fmt.Errorf("email is already confirmed")
			}

			newToken := rand.Intn(900000) + 100000
			userEmail.ConfirmationToken = &newToken
			now := time.Now()
			userEmail.CreatedAt = &now

			err = db.DB.Save(&userEmail).Error
			if err != nil {
				return "", fmt.Errorf("could not update confirmation token: %v", err)
			}

			auth := s.mailServices.New()

			if err := s.mailServices.Send(
				"Subject: AuthGo - Email Confirmation\r\nContent-Type: text/html; charset=utf-8\r\n\r\n<html><body><p>Verification code: "+strconv.Itoa(newToken)+"</p></body></html>\r\n",
				email,
				auth,
			); err != nil {
				return "", fmt.Errorf("failed to send confirmation email: %w", err)
			}
			return "confirmation code has been sent to your email", nil
		}
	}

	return "", fmt.Errorf("email %s not found", email)
}

func (s *UserServices) ConfirmEmail(code int, email string) (string, error) {
	var user types.User
	tx := db.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Joins("JOIN user_emails ON user_emails.user_id = users.id").
		Where("user_emails.email = ?", email).
		Preload("Emails").
		First(&user).Error
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("could not find user: %v", err)
	}

	for _, userEmail := range user.Emails {
		if userEmail.Email == email {
			if *userEmail.ConfirmationToken != code {
				tx.Rollback()
				return "", fmt.Errorf("invalid confirmation code")
			}

			if userEmail.IsConfirmed {
				tx.Rollback()
				return "", fmt.Errorf("email is already confirmed")
			}

			if time.Since(*userEmail.CreatedAt) > 1*time.Hour {
				tx.Rollback()
				return "", fmt.Errorf("confirmation code has expired")
			}

			userEmail.IsConfirmed = true
			userEmail.ConfirmationToken = nil

			err = tx.Save(&userEmail).Error
			if err != nil {
				tx.Rollback()
				return "", fmt.Errorf("could not confirm email: %v", err)
			}

			if err := tx.Commit().Error; err != nil {
				return "", fmt.Errorf("could not commit transaction: %w", err)
			}

			return "email has been confirmed", nil
		}
	}

	tx.Rollback()
	return "", fmt.Errorf("email %s not found", email)
}
