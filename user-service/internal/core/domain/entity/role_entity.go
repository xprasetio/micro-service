package entity

type RoleEntity struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"`
}

func (r *RoleEntity) TableName() string {
	return "roles"
}