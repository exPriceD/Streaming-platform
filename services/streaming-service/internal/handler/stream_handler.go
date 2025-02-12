package handler

import (
	"net/http"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/service"
	"github.com/labstack/echo/v4"
)

// StreamHandler управляет HTTP-запросами для стримов
type StreamHandler struct {
	streamService *service.StreamService
}

// NewStreamHandler создает новый обработчик стримов
func NewStreamHandler(router *echo.Echo, streamService *service.StreamService) {
	handler := &StreamHandler{streamService: streamService}

	// Регистрируем маршруты API
	streams := router.Group("/streams")
	{
		streams.POST("/start", handler.StartStream)
		streams.POST("/stop/:id", handler.StopStream)
		streams.GET("/:id", handler.GetStream)
	}
}

// StartStream обрабатывает запуск стрима
func (h *StreamHandler) StartStream(c echo.Context) error {
	var request struct {
		UserID      string `json:"user_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат запроса"})
	}

	stream, err := h.streamService.StartStream(request.UserID, request.Title, request.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stream)
}

// StopStream обрабатывает завершение стрима
func (h *StreamHandler) StopStream(c echo.Context) error {
	streamID := c.Param("id")

	err := h.streamService.StopStream(streamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Стрим завершён"})
}

// GetStream получает информацию о стриме
func (h *StreamHandler) GetStream(c echo.Context) error {
	streamID := c.Param("id")

	stream, err := h.streamService.GetStream(streamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if stream == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Стрим не найден"})
	}

	return c.JSON(http.StatusOK, stream)
}
