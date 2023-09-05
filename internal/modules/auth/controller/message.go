package controller

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	RetypePassword string `json:"retype_password"`
}

type LoginData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
