package model

import "time"

type UserRole struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64      `gorm:"index;not null" json:"user_id"`
	RoleID    int64      `gorm:"index;not null" json:"role_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Tabler interface {
	TableName() string
}

func (UserRole) TableName() string {
	return "user_role"
}
