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
	Token      string `json:"token"`
	RoleID     int64  `json:"role_id"`
}

type QueryStringCustomer struct {
	Search    string `form:"search"`
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	OrderBy   string `form:"sort"`
	OrderType string `form:"order"`
}
