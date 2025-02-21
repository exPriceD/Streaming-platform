package grpc

import (
	"context"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/service"
	"log/slog"
)

type Handler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
	logger      *slog.Logger
}

func NewHandler(userService *service.UserService, logger *slog.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}

func (h *Handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Info("GetUser called", slog.String("user_id", req.UserId))

	user, err := h.userService.GetUser(ctx, req.UserId)
	if err != nil {
		h.logger.Error("Failed to get user", slog.String("error", err.Error()), slog.String("user_id", req.UserId))
		return nil, err
	}

	return &pb.GetUserResponse{
		UserId:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Consent:   user.Consent,
		AvatarUrl: user.AvatarURL,
	}, nil
}
