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

	v1 := r.e.Group("/api/v1")
	{
		// Authentication endpoints
		auth := v1.Group("/auth")
		{
			// Public
			auth.POST("/register", r.handler.RegisterUser)
			auth.POST("/login", r.handler.LoginUser)

			// Protected
			auth.POST("/logout", r.handler.LogoutUser, authMiddleware)
		}

		// Protected user endpoints
		users := v1.Group("/users", authMiddleware)
		{
			users.GET("/me", r.handler.GetCurrentUser)
			users.PUT("/me", r.handler.UpdateCurrentUser)
			users.PATCH("/me/password", r.handler.ChangePassword)
			users.GET("/{userId}", r.handler.GetUserByID)
			users.PUT("/{userId}", r.handler.UpdateUser)
			users.DELETE("/{userId}", r.handler.DeleteUser)
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
