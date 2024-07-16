package types

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        *uuid.UUID  `gorm:"type:uuid;not null;primaryKey" json:"id"`
	Password  string      `gorm:"varchar(255);not null" json:"password"`
	AvatarURL string      `gorm:"type:varchar(255);not null" json:"avatar_url"`
	FirstName string      `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName  string      `gorm:"type:varchar(255);not null" json:"last_name"`
	CreatedAt *time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt *time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	Emails    []UserEmail `gorm:"foreignKey:UserID" json:"emails"`
}

type UserEmail struct {
	ID                *uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"id"`
	UserID            *uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	Email             string     `gorm:"type:varchar(255);unique;not null" json:"email"`
	IsPrimary         bool       `gorm:"not null;default:false" json:"is_primary"`
	IsConfirmed       bool       `gorm:"not null;default:false" json:"is_confirmed"`
	ConfirmationToken *uuid.UUID `gorm:"type:uuid" json:"confirmation_token"`
}

type RegisterBody struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponce struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type SendConfirmLinkBody struct {
	Email string `json:"email"`
}
