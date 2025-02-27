package rest

import (
	"log"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/service"
	"github.com/labstack/echo/v4"
)

// StartRESTServer запускает HTTP сервер
func StartRESTServer(streamService *service.StreamService, addr string) error {
	e := echo.New()
	handler := NewStreamHandler(streamService)

	e.POST("/streams/start", handler.StartStream)
	e.POST("/streams/stop/:id", handler.StopStream)
	e.GET("/streams/:id", handler.GetStream)

	log.Printf("REST API запущен на %s", addr)
	return e.Start(addr)
}
