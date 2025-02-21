package handler

import (
	"context"
	"errors"
	"github.com/exPriceD/Streaming-platform/pkg/proto/v1/auth"
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

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	authService AuthService
	log         *slog.Logger
}

func NewAuthHandler(authService AuthService, log *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
	}
}

func (h *AuthHandler) Authenticate(ctx context.Context, req *auth.AuthenticateRequest) (*auth.AuthenticateResponse, error) {
	h.log.Info("Received Authenticate request", slog.String("user_id", req.UserId))

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.log.Warn("Invalid User ID in Authenticate request", slog.String("user_id", req.UserId))
		return &auth.AuthenticateResponse{
			Error: &auth.Error{
				Code:    auth.ErrorCode_INVALID_ARGUMENT,
				Message: "Invalid User ID",
			},
		}, err
	}

	accessToken, refreshToken, expiresIn, expiresAt, err := h.authService.Authenticate(ctx, userID)
	if err != nil {
		h.log.Error("Failed to authenticate user", slog.String("user_id", userID.String()), slog.String("error", err.Error()))
		return &auth.AuthenticateResponse{
			Error: &auth.Error{
				Code:    auth.ErrorCode_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, err
	}

	h.log.Info("User authenticated", slog.String("user_id", userID.String()))
	return &auth.AuthenticateResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	h.log.Info("Received ValidateToken request")

	userID, err := h.authService.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		h.log.Warn("Invalid access token", slog.String("error", err.Error()))
		return &auth.ValidateTokenResponse{
			Valid: false,
			Error: mapErrorToProto(err),
		}, nil
	}

	h.log.Info("Access token validated", slog.String("user_id", userID.String()))
	return &auth.ValidateTokenResponse{
		Valid:  true,
		UserId: userID.String(),
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	h.log.Info("Received RefreshToken request")

	newAccessToken, newRefreshToken, expiresIn, expiresAt, err := h.authService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		h.log.Warn("Failed to refresh tokens", slog.String("error", err.Error()))
		return &auth.RefreshTokenResponse{
			Error: mapErrorToProto(err),
		}, err
	}

	h.log.Info("Tokens refreshed successfully")
	return &auth.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	h.log.Info("Received Logout request")

	err := h.authService.Logout(ctx, req.RefreshToken)
	if err != nil {
		h.log.Warn("Failed to logout user", slog.String("error", err.Error()))
		return &auth.LogoutResponse{
			Success: false,
			Error: &auth.Error{
				Code:    auth.ErrorCode_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, err
	}

	h.log.Info("User logged out successfully")
	return &auth.LogoutResponse{Success: true}, nil
}

func mapErrorToProto(err error) *auth.Error {
	switch {
	case errors.Is(err, service.ErrTokenExpired):
		return &auth.Error{
			Code:    auth.ErrorCode_TOKEN_EXPIRED,
			Message: "The token has expired",
		}
	case errors.Is(err, service.ErrTokenInvalid):
		return &auth.Error{
			Code:    auth.ErrorCode_TOKEN_INVALID,
			Message: "Invalid token",
		}
	case errors.Is(err, service.ErrRefreshTokenRevoked):
		return &auth.Error{
			Code:    auth.ErrorCode_REFRESH_TOKEN_REVOKED,
			Message: "Refresh token has been revoked",
		}
	default:
		return &auth.Error{
			Code:    auth.ErrorCode_INTERNAL_ERROR,
			Message: "Internal server error",
		}
	}
}
