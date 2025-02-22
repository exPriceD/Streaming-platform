package httpTransport

type RegisterRequest struct {
	Username        string `json:"username" example:"john_doe" description:"The username that the user wants to register with"`
	Email           string `json:"email" example:"john@example.com" description:"The email address that the user will use for registration"`
	Password        string `json:"password" example:"password123" description:"The password that the user wants to set for their account"`
	ConfirmPassword string `json:"confirm_password" example:"password123" description:"The password confirmation to ensure the user typed it correctly"`
	Consent         bool   `json:"consent" example:"true" description:"Indicates whether the user agrees to the terms and conditions"`
}

type LoginRequest struct {
	LoginIdentifier string `json:"login_identifier" example:"john_doe" description:"The login identifier (either username or email) of the user"`
	Password        string `json:"password" example:"password123" description:"The password that the user will use to log in"`
}
