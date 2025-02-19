package router

type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Consent         bool   `json:"consent"`
}

type LoginRequest struct {
	LoginIdentifier string `json:"loginIdentifier"`
	Password        string `json:"password"`
}
