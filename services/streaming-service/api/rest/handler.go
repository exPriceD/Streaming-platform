package rest

import (
	"net/http"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type StreamHandler struct {
	streamService *service.StreamService
}

func NewStreamHandler(streamService *service.StreamService) *StreamHandler {
	return &StreamHandler{streamService: streamService}
}

func (h *StreamHandler) StartStream(c echo.Context) error {
	var req struct {
		UserID      string `json:"user_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверный формат запроса"})
	}

	stream, err := h.streamService.StartStream(req.UserID, req.Title, req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stream)
}

func (h *StreamHandler) StopStream(c echo.Context) error {
	streamID := c.Param("id")
	err := h.streamService.StopStream(streamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Стрим завершён"})
}

func (h *StreamHandler) GetStream(c echo.Context) error {
	streamID := c.Param("id")
	stream, err := h.streamService.GetStream(uuid.MustParse(streamID).String())
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Стрим не найден"})
	}

	return c.JSON(http.StatusOK, stream)
}
