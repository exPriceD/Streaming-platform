package httpTransport

import (
	"context"
	"errors"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/entity"
	customErrors "github.com/exPriceD/Streaming-platform/services/user-service/internal/errors"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/transport/http/middleware"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"time"
)

const (
	refreshTokenCookieName = "refreshToken"
	cookieMaxAgeDelete     = -1
)

type UserService interface {
	RegisterUser(ctx context.Context, username, email, password, confirmPassword string, consent bool) (string, string, string, error)
	LoginUser(ctx context.Context, loginIdentifier, password string) (string, string, string, error)
	ValidateToken(ctx context.Context, accessToken string) (bool, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) (bool, error)
	GetUserByID(ctx context.Context, userId string) (*entity.User, error)
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

// RegisterUser godoc
// @Summary Регистрация нового пользователя
// @Description Создаёт нового пользователя и возвращает идентификатор пользователя и токен доступа. Если пользователь с таким email уже существует, возвращает ошибку.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные пользователя для регистрации"
// @Success 200 {object} RegisterResponse "Успешная регистрация"
// @Failure 400 {object} ErrorResponse "Неверный формат запроса или данные"
// @Failure 409 {object} ErrorResponse "Пользователь с таким email уже существует"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /api/v1/auth/register [post]
func (h *Handler) RegisterUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Failed to bind request", slog.String("error", err.Error()))
		c.Set("error_message", err.Error())
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
	}

	userId, accessToken, refreshToken, err := h.userService.RegisterUser(ctx, req.Username, req.Email, req.Password, req.ConfirmPassword, req.Consent)
	if err != nil {
		c.Set("error_message", err.Error())
		if errors.Is(err, customErrors.ErrInvalidInput) {
			h.logger.Warn("Invalid registration data", slog.String("email", req.Email), slog.String("error", err.Error()))
			return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
		if errors.Is(err, customErrors.ErrUserAlreadyExist) {
			h.logger.Info("User already exists", slog.String("email", req.Email))
			return c.JSON(http.StatusConflict, ErrorResponse{Error: "User with this email already exists"})
		}
		h.logger.Error("Failed to register user", slog.String("email", req.Email), slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to register user"})
	}

	h.setRefreshTokenCookie(c, refreshToken)

	h.logger.Info("User registered successfully", slog.String("userId", userId), slog.String("email", req.Email))
	c.Set("log_message", "User registered successfully")
	return c.JSON(http.StatusOK, RegisterResponse{
		Message:     "User registered successfully",
		UserId:      userId,
		AccessToken: accessToken,
	})
}

// LoginUser godoc
// @Summary Авторизация пользователя
// @Description Аутентифицирует пользователя по логину и паролю, возвращает идентификатор и access токен. При неверных данных возвращает ошибку.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для авторизации"
// @Success 200 {object} LoginResponse "Успешная авторизация"
// @Failure 400 {object} ErrorResponse "Неверный формат запроса или данные"
// @Failure 401 {object} ErrorResponse "Неверные учетные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /api/v1/auth/login [post]
func (h *Handler) LoginUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Failed to bind login request", slog.String("error", err.Error()))
		c.Set("error_message", err.Error())
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
	}

	userId, accessToken, refreshToken, err := h.userService.LoginUser(ctx, req.LoginIdentifier, req.Password)
	if err != nil {
		c.Set("error_message", err.Error())
		if errors.Is(err, customErrors.ErrInvalidInput) {
			h.logger.Warn("Invalid login data", slog.String("login", req.LoginIdentifier), slog.String("error", err.Error()))
			return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
		if errors.Is(err, customErrors.ErrUnauthorized) {
			h.logger.Info("Unauthorized login attempt", slog.String("login", req.LoginIdentifier))
			return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		}
		if errors.Is(err, customErrors.ErrUserNotFound) {
			h.logger.Info("User not found for login", slog.String("login", req.LoginIdentifier))
			return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not found"})
		}
		h.logger.Error("Failed to login user", slog.String("login", req.LoginIdentifier), slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to login"})
	}

	h.setRefreshTokenCookie(c, refreshToken)

	h.logger.Info("User logged in successfully", slog.String("userId", userId), slog.String("login", req.LoginIdentifier))
	c.Set("log_message", "User logged in successfully")
	return c.JSON(http.StatusOK, LoginResponse{
		Message:     "User logged in successfully",
		UserId:      userId,
		AccessToken: accessToken,
	})
}

// LogoutUser godoc
// @Summary Выход пользователя
// @Description Завершает сессию пользователя, удаляя refresh-токен из cookies. Требует refresh-токен в куках.
// @Tags Auth
// @Produce json
// @Security CookieAuth
// @Success 200 {object} LogoutResponse "Успешный выход"
// @Failure 401 {object} ErrorResponse "Refresh-токен отсутствует или недействителен"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /api/v1/auth/logout [post]
func (h *Handler) LogoutUser(c echo.Context) error {
	ctx := c.Request().Context()

	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		h.logger.Warn("No refresh token found during logout")
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "No refresh token found"})
	}

	success, err := h.userService.Logout(ctx, refreshToken.Value)
	if err != nil {
		h.logger.Error("Failed to logout user", slog.String("error", err.Error()))
		c.Set("error_message", err.Error())
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to logout"})
	}
	if !success {
		h.logger.Error("Logout failed unexpectedly")
		c.Set("error_message", "Logout failed unexpectedly")
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid refresh token"})
	}

	h.deleteRefreshTokenCookie(c)

	h.logger.Info("User logged out successfully")
	c.Set("log_message", "User logged out successfully")
	return c.JSON(http.StatusOK, LogoutResponse{Message: "User logged out successfully"})
}

// GetCurrentUser godoc
// @Summary Получить данные текущего пользователя
// @Description Возвращает информацию о текущем пользователе на основе аутентификации через токены.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Security CookieAuth
// @Success 200 {object} UserResponse "Данные пользователя успешно получены"
// @Failure 401 {object} ErrorResponse "Неавторизован или токен истёк"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /api/v1/users/me [get]
func (h *Handler) GetCurrentUser(c echo.Context) error {
	ctx := c.Request().Context()

	userId, ok := c.Get("userId").(string)
	if !ok {
		h.logger.Warn("No userId found in context")
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	response, err := h.getUserData(ctx, userId)
	if err != nil {
		c.Set("error_message", err.Error())
		return h.handleServiceError(c, err, userId)
	}
	return c.JSON(http.StatusOK, response)
}

// GetUserByID godoc
// @Summary Получить данные пользователя по Id
// @Description Возвращает информацию о пользователе по его идентификатору. Доступно только аутентифицированным пользователям.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Security CookieAuth
// @Param userId path string true "Идентификатор пользователя"
// @Success 200 {object} UserResponse "Данные пользователя успешно получены"
// @Failure 400 {object} ErrorResponse "Неверный формат userId"
// @Failure 401 {object} ErrorResponse "Неавторизован или токен истёк"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /api/v1/users/{userId} [get]
func (h *Handler) GetUserByID(c echo.Context) error {
	ctx := c.Request().Context()

	userId := c.Param("userId")

	response, err := h.getUserData(ctx, userId)
	if err != nil {
		c.Set("error_message", err.Error())
		return h.handleServiceError(c, err, userId)
	}
	return c.JSON(http.StatusOK, response)
}

// getUserData извлекает данные пользователя по ID
func (h *Handler) getUserData(ctx context.Context, userId string) (*UserResponse, error) {
	user, err := h.userService.GetUserByID(ctx, userId)
	if err != nil {
		h.logger.Error("Failed to get user data", slog.String("userId", userId), slog.String("error", err.Error()))
		return nil, err
	}
	h.logger.Info("User data retrieved successfully", slog.String("userId", userId))
	return &UserResponse{
		UserId:    user.Id.String(),
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}

// handleServiceError обрабатывает ошибки от сервиса и возвращает соответствующий HTTP-ответ
func (h *Handler) handleServiceError(c echo.Context, err error, userId string) error {
	if errors.Is(err, customErrors.ErrUserNotFound) {
		h.logger.Info("User not found", slog.String("userId", userId))
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
	}
	if errors.Is(err, customErrors.ErrInvalidInput) {
		h.logger.Warn("Invalid userId", slog.String("userId", userId))
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
	}
	h.logger.Error("Failed to get user data", slog.String("userId", userId), slog.String("error", err.Error()))
	return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve user data"})
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

func (h *Handler) UpdateUser(c echo.Context) error {
	return nil
}

func (h *Handler) DeleteUser(c echo.Context) error {
	return nil
}

func (h *Handler) setRefreshTokenCookie(c echo.Context, refreshToken string) {
	c.SetCookie(&http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false, // Включить для HTTPS в продакшене
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
}

func (h *Handler) deleteRefreshTokenCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   false, // Включить для HTTPS в продакшене
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cookieMaxAgeDelete,
	})
}
