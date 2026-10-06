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

type ProfileResponse struct {
	RoleName string `json:"role"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
	Lat      string `json:"lat"`
	Lng      string `json:"lng"`
}

type CustomerResponse struct {
	RoleName  string `json:"role"`
	RoleID    int64  `json:"role_id"`
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Photo     string `json:"photo"`
	Lat       string `json:"lat"`
	Lng       string `json:"lng"`
	BirthDate string `json:"birth_date"`
}
