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
