package handler

import (
	"net/http"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/service"
	"github.com/labstack/echo/v4"
)

type RTMPHandler struct {
	streamService *service.StreamService
}

func NewRTMPHandler(e *echo.Echo, streamService *service.StreamService) {
	handler := &RTMPHandler{streamService: streamService}

	e.POST("/rtmp/hook", handler.HandleRTMPEvent)
}

type RTMPEvent struct {
	Name   string `json:"name"`   // название события (on_publish, on_publish_done)
	Stream string `json:"stream"` // ID стрима
}

// HandleRTMPEvent обрабатывает события от Nginx RTMP
func (h *RTMPHandler) HandleRTMPEvent(c echo.Context) error {
	var event RTMPEvent
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверный формат данных"})
	}

	switch event.Name {
	case "on_publish":
		err := h.streamService.UpdateStreamStatus(event.Stream, "live")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось обновить статус стрима"})
		}
	case "on_publish_done":
		err := h.streamService.UpdateStreamStatus(event.Stream, "stopped")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось завершить стрим"})
		}
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неизвестное событие"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "событие обработано"})
}
