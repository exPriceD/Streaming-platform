package httpTransport

type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Consent         bool   `json:"consent"`
}

type LoginRequest struct {
	LoginIdentifier string `json:"login_identifier"`
	Password        string `json:"password"`
}
