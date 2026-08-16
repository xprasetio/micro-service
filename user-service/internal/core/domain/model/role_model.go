package model

import "time"

type Role struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `json:"name"`
	Users     []User     `gorm:"many2many:user_role;" json:"users"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
