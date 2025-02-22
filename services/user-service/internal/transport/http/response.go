package httpTransport

type RegisterResponse struct {
	Message     string `json:"message" example:"User registered successfully" description:"Success Message"`
	UserID      string `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" description:"User ID"`
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"Access token"`
}

type LoginResponse struct {
	Message     string `json:"message" example:"User logged in successfully" description:"Success Message"`
	UserID      string `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" description:"User ID"`
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"Access token"`
}

type LogoutResponse struct {
	Message string `json:"message" example:"User logged out successfully" description:"Success Message"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request payload" description:"Error description"`
}
