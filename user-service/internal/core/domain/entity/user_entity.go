package entity

type UserEntity struct {
	ID         int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	Photo      string `json:"photo"`
	Lat        string `json:"lat"`
	Lng        string `json:"lng"`
	IsVerified bool   `json:"is_verified"`
	RoleName   string `json:"role_name"`
}
