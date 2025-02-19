package router

import (
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/entity"
	"github.com/labstack/echo/v4"
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
	return nil
}

func (h *Handler) LoginUser(c echo.Context) error {
	return nil
}
