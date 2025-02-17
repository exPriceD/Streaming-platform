package handler

import (
	"context"
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/service"
	pb "github.com/exPriceD/Streaming-platform/services/auth-service/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
	log         *slog.Logger
}

func NewAuthHandler(authService *service.AuthService, log *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
	}
}

func (h *AuthHandler) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	h.log.Info("Received Authenticate request", slog.String("user_id", req.UserId))

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.log.Warn("Invalid User ID in Authenticate request", slog.String("user_id", req.UserId))
		return &pb.AuthenticateResponse{
			Error: &pb.Error{
				Code:    pb.ErrorCode_INVALID_ARGUMENT,
				Message: "Invalid User ID",
			},
		}, err
	}

	accessToken, refreshToken, expiresIn, expiresAt, err := h.authService.Authenticate(userID)
	if err != nil {
		h.log.Error("Failed to authenticate user", slog.String("user_id", userID.String()), slog.String("error", err.Error()))
		return &pb.AuthenticateResponse{
			Error: &pb.Error{
				Code:    pb.ErrorCode_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, err
	}

	h.log.Info("User authenticated", slog.String("user_id", userID.String()))
	return &pb.AuthenticateResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	h.log.Info("Received ValidateToken request")

	userID, err := h.authService.ValidateToken(req.AccessToken)
	if err != nil {
		h.log.Warn("Invalid access token", slog.String("error", err.Error()))
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: mapErrorToProto(err),
		}, nil
	}

	h.log.Info("Access token validated", slog.String("user_id", userID.String()))
	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID.String(),
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	h.log.Info("Received RefreshToken request")

	newAccessToken, newRefreshToken, expiresIn, expiresAt, err := h.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		h.log.Warn("Failed to refresh tokens", slog.String("error", err.Error()))
		return &pb.RefreshTokenResponse{
			Error: mapErrorToProto(err),
		}, err
	}

	h.log.Info("Tokens refreshed successfully")
	return &pb.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	h.log.Info("Received Logout request")

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		h.log.Warn("Failed to logout user", slog.String("error", err.Error()))
		return &pb.LogoutResponse{
			Success: false,
			Error: &pb.Error{
				Code:    pb.ErrorCode_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, err
	}

	h.log.Info("User logged out successfully")
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
