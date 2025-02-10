package handler

import (
	"context"
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
