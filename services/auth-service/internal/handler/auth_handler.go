package handler

import (
	"context"
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/service"
	pb "github.com/exPriceD/Streaming-platform/services/auth-service/proto"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, accessToken, refreshToken, err := h.authService.Register(req.Username, req.Email, req.Password, req.ConsentToDataProcessing)
	if err != nil {
		return &pb.RegisterResponse{
			Error: &pb.Error{
				Code:    pb.ErrorCode_INVALID_ARGUMENT,
				Message: err.Error(),
			},
		}, err
	}

	return &pb.RegisterResponse{
		UserId:       user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var identifier string
	var isEmail bool

	switch v := req.LoginIdentifier.(type) {
	case *pb.LoginRequest_Email:
		identifier = v.Email
		isEmail = true
	case *pb.LoginRequest_Username:
		identifier = v.Username
		isEmail = false
	default:
		return &pb.LoginResponse{
			Error: &pb.Error{
				Code:    pb.ErrorCode_INVALID_ARGUMENT,
				Message: "invalid login identifier",
			},
		}, nil
	}

	user, accessToken, refreshToken, err := h.authService.Authenticate(identifier, req.Password, isEmail)
	if err != nil {
		return &pb.LoginResponse{
			Error: &pb.Error{
				Code:    pb.ErrorCode_UNAUTHORIZED,
				Message: err.Error(),
			},
		}, err
	}

	return &pb.LoginResponse{
		UserId:       user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
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
		UserId: userID,
	}, nil
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
	default:
		return &pb.Error{
			Code:    pb.ErrorCode_INTERNAL_ERROR,
			Message: "Internal server error",
		}
	}
}
