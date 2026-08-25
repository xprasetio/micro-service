package model

import "time"

type VerificationToken struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64      `gorm:"index;not null" json:"user_id"`
	Token     string     `json:"token"`
	TokenType string     `json:"token_type"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Users     []User     `gorm:"many2many:user_verification_token; foreignKey:UserID" json:"users"`
}
