package types

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        *uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"id,omitempty"`
	Email     string     `gorm:"type:varchar(255);unique;not null" json:"email,omitempty"`
	Password  string     `gorm:"varchar(255);not null" json:"password"`
	AvatarURL string     `gorm:"type:varchar(255);not null" json:"avatar_url,omitempty"`
	FirstName string     `gorm:"type:varchar(255);not null" json:"first_name,omitempty"`
	LastName  string     `gorm:"type:varchar(255);not null" json:"last_name,omitempty"`
	CreatedAt *time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"not null;default:now()" json:"updatedAt"`
}

type RegisterBody struct {
	Email     string `json:"email,omitempty"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}
