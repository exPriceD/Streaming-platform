package grpcTransport

import (
	"context"
	pb "github.com/exPriceD/Streaming-platform/pkg/proto/v1/user"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type UserUsecase interface {
	GetUserByID(ctx context.Context, userId string) (domain.User, error)
}

type Handler struct {
	pb.UnimplementedUserServiceServer
	userUsecase UserUsecase
	logger      *slog.Logger
}

func NewHandler(userUsecase UserUsecase, logger *slog.Logger) *Handler {
	return &Handler{
		userUsecase: userUsecase,
		logger:      logger,
	}
}

func (h *Handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Info("GetUser called", slog.String("userId", req.UserId))

	user, err := h.userUsecase.GetUserByID(ctx, req.UserId)
	if err != nil {
		h.logger.Error("Failed to get user", slog.String("error", err.Error()), slog.String("userId", req.UserId))
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.GetUserResponse{
		UserId:    user.ID(),
		Username:  user.Username(),
		Email:     user.Email(),
		AvatarUrl: user.AvatarURL(),
	}, nil
}
