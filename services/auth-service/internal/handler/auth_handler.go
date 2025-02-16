package handler

import (
	"context"
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/service"
	pb "github.com/exPriceD/Streaming-platform/services/auth-service/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &pb.AuthenticateResponse{
			Error: &pb.Error{
				Code:    pb.ErrorCode_INVALID_ARGUMENT,
				Message: "Invalid User ID",
			},
		}, err
	}

	accessToken, refreshToken, expiresIn, expiresAt, err := h.authService.Authenticate(userID)
	if err != nil {
		return &pb.AuthenticateResponse{
			Error: &pb.Error{
				Code:    pb.ErrorCode_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, err
	}

	return &pb.AuthenticateResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userID, err := h.authService.ValidateAccessToken(req.AccessToken)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: mapErrorToProto(err),
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID.String(),
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	newAccessToken, newRefreshToken, expiresIn, expiresAt, err := h.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		return &pb.RefreshTokenResponse{
			Error: mapErrorToProto(err),
		}, err
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		return &pb.LogoutResponse{
			Success: false,
			Error: &pb.Error{
				Code:    pb.ErrorCode_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, err
	}

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
