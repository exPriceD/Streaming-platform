package handler

import (
	"net/http"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/service"

	"github.com/gin-gonic/gin"
)

// StreamHandler управляет HTTP-запросами для стримов
type StreamHandler struct {
	streamService *service.StreamService
}

// NewStreamHandler создает новый обработчик стримов
func NewStreamHandler(router *gin.Engine, streamService *service.StreamService) {
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
func (h *StreamHandler) StartStream(c *gin.Context) {
	var request struct {
		UserID      string `json:"user_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	stream, err := h.streamService.StartStream(request.UserID, request.Title, request.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stream)
}

// StopStream обрабатывает завершение стрима
func (h *StreamHandler) StopStream(c *gin.Context) {
	streamID := c.Param("id")

	err := h.streamService.StopStream(streamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Стрим завершён"})
}

// GetStream получает информацию о стриме
func (h *StreamHandler) GetStream(c *gin.Context) {
	streamID := c.Param("id")

	stream, err := h.streamService.GetStream(streamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if stream == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Стрим не найден"})
		return
	}

	c.JSON(http.StatusOK, stream)
}
