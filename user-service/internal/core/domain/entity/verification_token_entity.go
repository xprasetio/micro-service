package entity

import "time"

type VerificationToken struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64      `gorm:"index;not null" json:"user_id"`
	Token     string     `json:"token"`
	TokenType string     `json:"token_type"`
	ExpiresAt time.Time  `json:"expires_at"`
	User      UserEntity `gorm:"foreignKey:UserID" json:"user"`
}
