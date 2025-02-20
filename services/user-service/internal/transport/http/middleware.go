package router

import (
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthMiddleware struct {
	userService *service.UserService
}

func (h *AuthMiddleware) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := c.Request().Header.Get("Authorization")
		if accessToken == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Access token is required"})
		}

		// Добавить проверку токена в кэше

		valid, err := h.userService.ValidateToken(accessToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token validate failed"})
		}

		if valid {
			return next(c)
		}

		refreshToken, err := c.Cookie("refreshToken")
		if err != nil || refreshToken == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Refresh token is required"})
		}

		newAccessToken, newRefreshToken, err := h.userService.RefreshToken(refreshToken.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token refresh failed"})
		}

		c.SetCookie(&http.Cookie{
			Name:     "refreshToken",
			Value:    newRefreshToken,
			HttpOnly: true,
			// Должно быть включено при использовании HTTPS
			// Secure:   false,
			// SameSite: http.SameSiteStrictMode,
		})

		c.Response().Header().Set("Authorization", "Bearer "+newAccessToken)

		return next(c)
	}
}
