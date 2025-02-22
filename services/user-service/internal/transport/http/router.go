package httpTransport

import (
	"context"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/config"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/transport/http/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log/slog"
)

type Router struct {
	e       *echo.Echo
	handler *Handler
	logger  *slog.Logger
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
		logger:  logger,
	}

	router.registerRoutes()

	return router
}

func (r *Router) registerRoutes() {
	authMiddleware := r.handler.GetAuthMiddleware()

	// API Version 1 group
	v1 := r.e.Group("/api/v1")
	{
		// Authentication endpoints
		auth := v1.Group("/auth")
		{
			// Public
			auth.POST("/register", r.handler.RegisterUser)
			auth.POST("/login", r.handler.LoginUser)

			// Protected
			protectedAuth := auth.Group("", authMiddleware)
			{
				protectedAuth.POST("/logout", r.handler.LogoutUser)
			}
		}

		// User endpoints
		users := v1.Group("/users")
		{
			// Public
			users.GET("/{userId}", r.handler.GetUserByID)

			// Protected
			protectedUsers := users.Group("", authMiddleware)
			{
				protectedUsers.GET("/me", r.handler.GetCurrentUser)
				protectedUsers.PUT("/me", r.handler.UpdateCurrentUser)
				protectedUsers.PATCH("/me/password", r.handler.ChangePassword)
			}
		}
	}

	// Service endpoints
	r.e.GET("/health", func(c echo.Context) error { return c.NoContent(200) })
	r.e.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (r *Router) Run(address string) error {
	return r.e.Start(address)
}

func (r *Router) Shutdown(ctx context.Context) error {
	r.logger.Info("Shutting down HTTP server")
	return r.e.Shutdown(ctx)
}
