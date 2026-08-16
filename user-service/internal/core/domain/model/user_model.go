package model

import "time"

type User struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Address    string     `json:"address"`
	Phone      string     `json:"phone"`
	Photo      string     `json:"photo"`
	Lat        string     `json:"lat"`
	Lng        string     `json:"lng"`
	IsVerified bool       `json:"is_verified"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	Roles      []Role     `gorm:"many2many:user_role;" json:"roles"`
}
