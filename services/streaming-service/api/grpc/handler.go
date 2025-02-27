package grpc

import (
	"context"
	"time"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/service"
	"github.com/exPriceD/Streaming-platform/services/streaming-service/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// StreamHandler подрубает gRPC интерфейс для StreamingService
type StreamHandler struct {
	proto.UnimplementedStreamingServiceServer
	streamService *service.StreamService
}

// NewStreamHandler создаёт новый gRPC хендлер
func NewStreamHandler(streamService *service.StreamService) *StreamHandler {
	return &StreamHandler{streamService: streamService}
}

// StartStream запускает новый стрим
func (h *StreamHandler) StartStream(ctx context.Context, req *proto.StartStreamRequest) (*proto.StreamResponse, error) {
	stream, err := h.streamService.StartStream(req.UserId, req.Title, req.Description)
	if err != nil {
		return nil, err
	}

	return &proto.StreamResponse{
		Id:          stream.ID.String(),
		UserId:      stream.UserID,
		Title:       stream.Title,
		Description: stream.Description,
		Status:      stream.Status,
		CreatedAt:   timestamppb.New(stream.CreatedAt),
		UpdatedAt:   timestamppb.New(stream.UpdatedAt),
	}, nil
}

// StopStream останавливает стрим
func (h *StreamHandler) StopStream(ctx context.Context, req *proto.StopStreamRequest) (*proto.StreamResponse, error) {
	err := h.streamService.StopStream(uuid.MustParse(req.StreamId).String())
	if err != nil {
		return nil, err
	}

	stream, err := h.streamService.GetStream(uuid.MustParse(req.StreamId).String())
	if err != nil {
		return nil, err
	}

	return &proto.StreamResponse{
		Id:        stream.ID.String(),
		UserId:    stream.UserID,
		Title:     stream.Title,
		Status:    stream.Status,
		UpdatedAt: timestamppb.New(time.Now()),
	}, nil
}

// GetStream получает информацию о стриме
func (h *StreamHandler) GetStream(ctx context.Context, req *proto.GetStreamRequest) (*proto.StreamResponse, error) {
	stream, err := h.streamService.GetStream(uuid.MustParse(req.StreamId).String())
	if err != nil {
		return nil, err
	}

	return &proto.StreamResponse{
		Id:        stream.ID.String(),
		UserId:    stream.UserID,
		Title:     stream.Title,
		Status:    stream.Status,
		CreatedAt: timestamppb.New(stream.CreatedAt),
		UpdatedAt: timestamppb.New(stream.UpdatedAt),
	}, nil
}

// GenerateStreamKey создаёт новый stream-key для пользователя
func (h *StreamHandler) GenerateStreamKey(ctx context.Context, req *proto.GenerateStreamKeyRequest) (*proto.GenerateStreamKeyResponse, error) {
	streamKey, err := h.streamService.GenerateStreamKey(req.UserId)
	if err != nil {
		return nil, err
	}

	return &proto.GenerateStreamKeyResponse{
		UserId:    req.UserId,
		StreamKey: streamKey,
	}, nil
}

// GetStreamKey получает существующий stream-key пользователя
func (h *StreamHandler) GetStreamKey(ctx context.Context, req *proto.GetStreamKeyRequest) (*proto.GetStreamKeyResponse, error) {
	streamKey, exists, err := h.streamService.GetStreamKey(req.UserId)
	if err != nil {
		return nil, err
	}

	return &proto.GetStreamKeyResponse{
		UserId:    req.UserId,
		StreamKey: streamKey,
		Exists:    exists,
	}, nil
}

// RegenerateStreamKey пересоздаёт stream-key пользователя
func (h *StreamHandler) RegenerateStreamKey(ctx context.Context, req *proto.RegenerateStreamKeyRequest) (*proto.RegenerateStreamKeyResponse, error) {
	streamKey, err := h.streamService.RegenerateStreamKey(req.UserId)
	if err != nil {
		return nil, err
	}

	return &proto.RegenerateStreamKeyResponse{
		UserId:    req.UserId,
		StreamKey: streamKey,
	}, nil
}
