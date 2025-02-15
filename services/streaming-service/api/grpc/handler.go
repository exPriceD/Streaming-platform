package grpc

import (
	"context"
	"time"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/service"
	"github.com/exPriceD/Streaming-platform/services/streaming-service/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// StreamHandler implements the gRPC interface for StreamingService
type StreamHandler struct {
	proto.UnimplementedStreamingServiceServer
	streamService *service.StreamService
}

// NewStreamHandler creates a new gRPC handler
func NewStreamHandler(streamService *service.StreamService) *StreamHandler {
	return &StreamHandler{streamService: streamService}
}

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
