package response

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	RoleName     string `json:"role_name"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Photo        string `json:"photo"`
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
	IsVerified   bool   `json:"is_verified"`
}
