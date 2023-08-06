package controller

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	RetypePassword string `json:"retype_password"`
}

type AuthResponse struct {
	Success   bool      `json:"success"`
	ErrorCode int       `json:"error_code,omitempty"`
	Data      LoginData `json:"data"`
}

type LoginData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}
