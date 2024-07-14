package types

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        *uuid.UUID  `gorm:"type:uuid;not null;primaryKey" json:"id,omitempty"`
	Password  string      `gorm:"varchar(255);not null" json:"password"`
	AvatarURL string      `gorm:"type:varchar(255);not null" json:"avatar_url,omitempty"`
	FirstName string      `gorm:"type:varchar(255);not null" json:"first_name,omitempty"`
	LastName  string      `gorm:"type:varchar(255);not null" json:"last_name,omitempty"`
	CreatedAt *time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt *time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	Emails    []UserEmail `gorm:"foreignKey:UserID" json:"emails,omitempty"`
}

type UserEmail struct {
	ID        *uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"id,omitempty"`
	UserID    *uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	Email     string     `gorm:"type:varchar(255);unique;not null" json:"email,omitempty"`
	IsPrimary bool       `gorm:"not null;default:false" json:"is_primary,omitempty"`
}

type RegisterBody struct {
	Email     string `json:"email,omitempty"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponce struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
