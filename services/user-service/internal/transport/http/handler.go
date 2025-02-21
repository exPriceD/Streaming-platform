package router

import (
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/entity"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/transport/http/middleware"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type UserService interface {
	RegisterUser(username, email, password, confirmPassword string, consent bool) (string, string, string, *entity.User, error)
	LoginUser(loginIdentifier, password string) (string, string, string, error)
	ValidateToken(accessToken string) (bool, error)
	RefreshToken(refreshToken string) (string, string, error)
	Logout(refreshToken string) (bool, error)
}

type Handler struct {
	userService UserService
	logger      *slog.Logger
}

func NewHandler(userService UserService, logger *slog.Logger) *Handler {
	return &Handler{userService: userService, logger: logger}
}

func (h *Handler) GetAuthMiddleware() echo.MiddlewareFunc {
	return (&middleware.AuthMiddleware{AuthService: h.userService}).UserIdentity
}

func (h *Handler) RegisterUser(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Failed to bind request", slog.String("error", err.Error()))
		c.Set("error_message", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	userId, accessToken, refreshToken, _, err := h.userService.RegisterUser(req.Username, req.Email, req.Password, req.ConfirmPassword, req.Consent)
	if err != nil {
		h.logger.Error("Failed to register user", slog.String("error", err.Error()), slog.String("email", req.Email))
		c.Set("error_message", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		// Должно быть включено при использовании HTTPS
		// Secure:   false,
		// SameSite: http.SameSiteStrictMode,
	})

	h.logger.Info("User registered successfully", slog.String("userId", userId), slog.String("email", req.Email))
	c.Set("log_message", "User registered successfully")
	return c.JSON(http.StatusOK, echo.Map{
		"message":      "User registered successfully",
		"user_id":      userId,
		"access_token": accessToken,
	})
}

func (h *Handler) LoginUser(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Failed to bind login request", slog.String("error", err.Error()))
		c.Set("error_message", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	userId, accessToken, refreshToken, err := h.userService.LoginUser(req.LoginIdentifier, req.Password)
	if err != nil {
		h.logger.Error("Failed to login user", slog.String("error", err.Error()), slog.String("login", req.LoginIdentifier))
		c.Set("error_message", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		// Должно быть включено при использовании HTTPS
		// Secure:   false,
		// SameSite: http.SameSiteStrictMode,
	})

	h.logger.Info("User logged in successfully", slog.String("userId", userId), slog.String("login", req.LoginIdentifier))
	c.Set("log_message", "User logged in successfully")
	return c.JSON(http.StatusOK, echo.Map{
		"message":      "User logged in successfully",
		"user_id":      userId,
		"access_token": accessToken,
	})
}

func (h *Handler) LogoutUser(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		h.logger.Warn("No refresh token found during logout")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No refresh token found"})
	}

	ok, err := h.userService.Logout(refreshToken.Value)
	if err != nil {
		h.logger.Error("Failed to logout user", slog.String("error", err.Error()))
		c.Set("error_message", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if !ok {
		h.logger.Error("Logout failed unexpectedly")
		c.Set("error_message", "Logout failed unexpectedly")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Logout failed"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		HttpOnly: true,
		// Должно быть включено при использовании HTTPS
		// Secure:   false,
		// SameSite: http.SameSiteStrictMode,
		MaxAge: -1, // Удаляем cookie
	})

	h.logger.Info("User logged out successfully")
	c.Set("log_message", "User logged out successfully")
	return c.JSON(http.StatusOK, echo.Map{"message": "User logged out successfully"})
}

func (h *Handler) GetCurrentUser(c echo.Context) error {
	return nil
}

func (h *Handler) GetUserData(c echo.Context) error {
	return nil
}

func (h *Handler) ChangePassword(c echo.Context) error {
	return nil
}

func (h *Handler) ForgotPassword(c echo.Context) error {
	return nil
}

func (h *Handler) UpdateCurrentUser(c echo.Context) error {
	return nil
}

func (h *Handler) GetUserByID(c echo.Context) error {
	return nil
}

func (h *Handler) UpdateUser(c echo.Context) error {
	return nil
}

func (h *Handler) DeleteUser(c echo.Context) error {
	return nil
}
