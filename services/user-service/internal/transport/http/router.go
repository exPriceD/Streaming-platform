package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"log"
	"os"
)

type Router struct {
	e       *echo.Echo
	handler *Handler
}

func NewRouter(handler *Handler) *Router {
	e := echo.New()

	mw := io.MultiWriter(os.Stdout, setLogsFile())

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: mw,
	}))

	e.Use(middleware.Recover())

	router := &Router{
		e:       e,
		handler: handler,
	}

	router.registerRoutes()

	return router
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func (r *Router) registerRoutes() {
	r.e.POST("/register", r.handler.RegisterUser)
	r.e.POST("/login", r.handler.LoginUser)
	r.e.POST("/logout", r.handler.LogoutUser)
	r.e.POST("/forgot-password", r.handler.ForgotPassword) // Запрос на восстановление пароля

	v1 := r.e.Group("/api/v1", r.handler.GetAuthMiddleware())
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
