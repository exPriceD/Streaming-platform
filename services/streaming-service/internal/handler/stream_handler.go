package handler

import (
	"net/http"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/service"
	auth "github.com/exPriceD/Streaming-platform/services/streaming-service/internal/token"
	"github.com/labstack/echo/v4"
)

// StreamHandler управляет HTTP-запросами для стримов
type StreamHandler struct {
	streamService *service.StreamService
	jwtManager    *auth.JWTManager
}

// NewStreamHandler создаёт новый обработчик стримов
func NewStreamHandler(e *echo.Echo, streamService *service.StreamService, jwtManager *auth.JWTManager) {
	handler := &StreamHandler{streamService: streamService, jwtManager: jwtManager}

	streams := e.Group("/streams")
	{
		streams.POST("/start", handler.StartStream)
		streams.POST("/stop/:id", handler.StopStream)
		streams.GET("/:id", handler.GetStream)
	}
}

// StartStream обрабатывает запуск стрима (требуется JWT)
func (h *StreamHandler) StartStream(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "токен отсутствует"})
	}

	userID, err := h.jwtManager.VerifyToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "неверный токен"})
	}

	var request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверный формат запроса"})
	}

	stream, err := h.streamService.StartStream(userID, request.Title, request.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stream)
}

// StopStream завершает стрим (требуется JWT)
func (h *StreamHandler) StopStream(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "токен отсутствует"})
	}

	_, err := h.jwtManager.VerifyToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "неверный токен"})
	}

	streamID := c.Param("id")
	err = h.streamService.StopStream(streamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Стрим завершён"})
}

// GetStream получает информацию о стриме по ID
func (h *StreamHandler) GetStream(c echo.Context) error {
	streamID := c.Param("id")

	stream, err := h.streamService.GetStream(streamID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Стрим не найден"})
	}

	return c.JSON(http.StatusOK, stream)
}
