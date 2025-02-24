package grpcTransport

import (
	"context"
	"errors"
	pb "github.com/exPriceD/Streaming-platform/pkg/proto/v1/auth"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/service"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"time"
)

type AuthService interface {
	Authenticate(ctx context.Context, userID uuid.UUID) (string, string, int64, time.Time, error)
	ValidateToken(ctx context.Context, accessToken string) (uuid.UUID, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, int64, time.Time, error)
	Logout(ctx context.Context, refreshToken string) error
}

type Handler struct {
	pb.UnimplementedAuthServiceServer
	authService AuthService
	logger      *slog.Logger
}

func NewAuthHandler(authService AuthService, logger *slog.Logger) *Handler {
	return &Handler{
		authService: authService,
		logger:      logger,
	}
}

func (h *Handler) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	h.logger.Info("Received Authenticate request", slog.String("user_id", req.UserId))

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.logger.Warn("Invalid User ID in Authenticate request", slog.String("user_id", req.UserId))
		return &pb.AuthenticateResponse{
			Error: &pb.Error{
				Code:    pb.ErrorCode_INVALID_ARGUMENT,
				Message: "Invalid User ID",
			},
		}, err
	}

	accessToken, refreshToken, expiresIn, expiresAt, err := h.authService.Authenticate(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to authenticate user", slog.String("user_id", userID.String()), slog.String("error", err.Error()))
		return &pb.AuthenticateResponse{
			Error: &pb.Error{
				Code:    pb.ErrorCode_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, err
	}

	h.logger.Info("User authenticated", slog.String("user_id", userID.String()))
	return &pb.AuthenticateResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

func (h *Handler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	h.logger.Info("Received ValidateToken request")

	userID, err := h.authService.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		h.logger.Warn("Invalid access token", slog.String("error", err.Error()))
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: mapErrorToProto(err),
		}, nil
	}

	h.logger.Info("Access token validated", slog.String("user_id", userID.String()))
	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID.String(),
	}, nil
}

func (h *Handler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	h.logger.Info("Received RefreshToken request")

	newAccessToken, newRefreshToken, expiresIn, expiresAt, err := h.authService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		h.logger.Warn("Failed to refresh tokens", slog.String("error", err.Error()))
		return &pb.RefreshTokenResponse{
			Error: mapErrorToProto(err),
		}, err
	}

	h.logger.Info("Tokens refreshed successfully")
	return &pb.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

func (h *Handler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	h.logger.Info("Received Logout request")

	err := h.authService.Logout(ctx, req.RefreshToken)
	if err != nil {
		h.logger.Warn("Failed to logout user", slog.String("error", err.Error()))
		return &pb.LogoutResponse{
			Success: false,
			Error: &pb.Error{
				Code:    pb.ErrorCode_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, err
	}

	h.logger.Info("User logged out successfully")
	return &pb.LogoutResponse{Success: true}, nil
}

func mapErrorToProto(err error) *pb.Error {
	switch {
	case errors.Is(err, service.ErrTokenExpired):
		return &pb.Error{
			Code:    pb.ErrorCode_TOKEN_EXPIRED,
			Message: "The token has expired",
		}
	case errors.Is(err, service.ErrTokenInvalid):
		return &pb.Error{
			Code:    pb.ErrorCode_TOKEN_INVALID,
			Message: "Invalid token",
		}
	case errors.Is(err, service.ErrRefreshTokenRevoked):
		return &pb.Error{
			Code:    pb.ErrorCode_REFRESH_TOKEN_REVOKED,
			Message: "Refresh token has been revoked",
		}
	default:
		return &pb.Error{
			Code:    pb.ErrorCode_INTERNAL_ERROR,
			Message: "Internal server error",
		}
	}
}
