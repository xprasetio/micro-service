package request

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

type SignUpRequest struct {
	Name                 string `json:"name" validate:"required,min=3,max=32"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=6,max=32"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,max=32"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdatePasswordRequest struct {
	CurrentPassword         string `json:"password,omitempty"`
	NewPassword             string `json:"password_new" validate:"required,min=6,max=32"`
	NewPasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,max=32"`
}

type UpdateDataUserRequest struct {
	Name                 string `json:"name" validate:"required,min=3,max=32"`
	Email                string `json:"email" validate:"required,email"`
	Address              string `json:"address" validate:"required,min=3,max=32"`
	Phone                string `json:"phone" validate:"required,min=3,max=32"`
	Photo                string `json:"photo" validate:"required,min=3,max=32"`
	Lat                  string `json:"lat" validate:"required,min=3,max=32"`
	Lng                  string `json:"lng" validate:"required,min=3,max=32"`
}
