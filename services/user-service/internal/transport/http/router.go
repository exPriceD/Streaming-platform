package router

import (
	"github.com/labstack/echo/v4"
	"log/slog"
)

type Router struct {
	e       *echo.Echo
	handler *Handler
	log     *slog.Logger
}

func NewRouter(handler *Handler, log *slog.Logger) *Router {
	e := echo.New()

	router := &Router{
		e:       e,
		handler: handler,
		log:     log,
	}
	router.registerRoutes()

	return router
}

func (r *Router) registerRoutes() {
	r.e.POST("/register", r.handler.RegisterUser)
	r.e.POST("/login", r.handler.LoginUser)
}

func (r *Router) Run(address string) error {
	return r.e.Start(address)
}
