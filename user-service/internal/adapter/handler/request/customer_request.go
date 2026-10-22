package request

type CustomerRequest struct {
	Email                string  `json:"email" validate:"required,email"`
	Password             string  `json:"password" validate:"required,min=8"`
	PasswordConfirmation string  `json:"password_confirmation" validate:"required,eqfield=Password"`
	Name                 string  `json:"name" validate:"required"`
	Phone                string  `json:"phone" validate:"required"`
	Address              string  `json:"address" validate:"required"`
	Lat                  float64 `json:"lat" validate:"required"`
	Lng                  float64 `json:"lng" validate:"required"`
	RoleID               int64   `json:"role_id" validate:"required"`
}
