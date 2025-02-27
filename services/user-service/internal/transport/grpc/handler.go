package grpcTransport

import (
	"context"
	pb "github.com/exPriceD/Streaming-platform/pkg/proto/v1/user"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type UserService interface {
	GetUserByID(ctx context.Context, userId string) (*entity.User, error)
}

type Handler struct {
	pb.UnimplementedUserServiceServer
	userService UserService
	logger      *slog.Logger
}

func NewHandler(userService UserService, logger *slog.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}

func (h *Handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Info("GetUser called", slog.String("userId", req.UserId))

	user, err := h.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		h.logger.Error("Failed to get user", slog.String("error", err.Error()), slog.String("userId", req.UserId))
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.GetUserResponse{
		UserId:    user.Id.String(),
		Username:  user.Username,
		Email:     user.Email,
		AvatarUrl: user.AvatarURL,
	}, nil
}
