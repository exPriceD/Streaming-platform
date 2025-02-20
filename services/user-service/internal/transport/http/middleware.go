package router

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	userService UserService
}

func (h *AuthMiddleware) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Authorization header is required"})
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid Authorization header format"})
		}
		accessToken := strings.TrimPrefix(authHeader, bearerPrefix)

		// TODO: Добавить проверку токена в кэше (например, Redis)
		// if cached, err := h.checkTokenInCache(accessToken); cached {
		//     return next(c)
		// }

		valid, err := h.userService.ValidateToken(accessToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Failed to validate token"})
		}

		if valid {
			return next(c)
		}

		refreshCookie, err := c.Cookie("refreshToken")
		if err != nil || refreshCookie == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Refresh token is required"})
		}

		newAccessToken, newRefreshToken, err := h.userService.RefreshToken(refreshCookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Failed to refresh token"})
		}

		c.SetCookie(&http.Cookie{
			Name:     "refreshToken",
			Value:    newRefreshToken,
			HttpOnly: true,
			// Должно быть включено при использовании HTTPS
			// Secure:   false,
			// SameSite: http.SameSiteStrictMode,
		})

		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error":        "Access token expired",
			"accessToken":  newAccessToken,
			"refreshToken": newRefreshToken,
		})
	}
}
