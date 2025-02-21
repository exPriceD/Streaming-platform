package httpTransport

import (
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/config"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/transport/http/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"log/slog"
)

type Router struct {
	e       *echo.Echo
	handler *Handler
}

func NewRouter(handler *Handler, logger *slog.Logger, CORS config.CORSConfig) *Router {
	e := echo.New()

	e.Use(middleware.NewLoggerMiddleware(middleware.LoggerMiddlewareConfig{
		Logger:       logger,
		ConsoleLevel: slog.LevelInfo,
		FileLevel:    slog.LevelDebug,
		LogFilePath:  "logs/requests.log",
	}))

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     CORS.AllowOrigins,
		AllowMethods:     CORS.AllowMethods,
		AllowHeaders:     CORS.AllowHeaders,
		AllowCredentials: CORS.AllowCredentials,
		MaxAge:           CORS.MaxAge,
	}))

	e.Use(echoMiddleware.Recover())

	router := &Router{
		e:       e,
		handler: handler,
	}

	router.registerRoutes()

	return router
}

func (r *Router) registerRoutes() {
	authMiddleware := r.handler.GetAuthMiddleware()

	r.e.POST("/register", r.handler.RegisterUser)
	r.e.POST("/login", r.handler.LoginUser)

	r.e.POST("/logout", r.handler.LogoutUser, authMiddleware)

	v1 := r.e.Group("/api/v1", authMiddleware)
	{
		v1.GET("/users/me", r.handler.GetCurrentUser)            // Получение информации о текущем пользователе
		v1.PUT("/users/me", r.handler.UpdateCurrentUser)         // Обновление данных профиля текущего пользователя
		v1.PATCH("/users/me/password", r.handler.ChangePassword) // Изменение пароля пользователя
		v1.GET("/users/{userId}", r.handler.GetUserByID)         // Получение информации о другом пользователе по ID
		v1.PUT("/users/{userId}", r.handler.UpdateUser)          // Обновление данных пользователя (например, администратор)
		v1.DELETE("/users/{userId}", r.handler.DeleteUser)       // Удаление пользователя (например, администратор)
	}

	r.e.GET("/health", func(c echo.Context) error { return c.NoContent(200) })
}

func (r *Router) Run(address string) error {
	return r.e.Start(address)
}
