package router

import (
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/entity"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserService interface {
	RegisterUser(username, email, password, confirmPassword string, consent bool) (string, *string, *string, *entity.User, error)
	LoginUser(loginIdentifier, password string) (string, *string, *string, error)
}

type Handler struct {
	userService UserService
}

func NewHandler(userService UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) RegisterUser(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	userID, accessToken, refreshToken, _, err := h.userService.RegisterUser(req.Username, req.Email, req.Password, req.ConfirmPassword, req.Consent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "User registered successfully",
		"userID":       userID,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) LoginUser(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	userID, accessToken, refreshToken, err := h.userService.LoginUser(req.LoginIdentifier, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "User logged in successfully",
		"userID":       userID,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
