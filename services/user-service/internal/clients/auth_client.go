package clients

import (
	"context"
	"errors"
	authpb "github.com/exPriceD/Streaming-platform/services/auth-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type AuthClient struct {
	client authpb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(authServiceAddr string) (*AuthClient, error) {
	conn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Couldn't connect to AuthService: %v", err)
		return nil, err
	}
	client := authpb.NewAuthServiceClient(conn)
	return &AuthClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close закрывает gRPC-соединение
func (ac *AuthClient) Close() error {
	if ac.conn != nil {
		return ac.conn.Close()
	}
	return errors.New("connection is already closed or nil")
}

// Authenticate вызывает gRPC-метод регистрации пользователя
func (ac *AuthClient) Authenticate(ctx context.Context, req *authpb.AuthenticateRequest) (*authpb.AuthenticateResponse, error) {
	return ac.client.Authenticate(ctx, req)
}

// ValidateToken вызывает gRPC-метод проверки токена
func (ac *AuthClient) ValidateToken(ctx context.Context, token string) (*authpb.ValidateTokenResponse, error) {
	req := &authpb.ValidateTokenRequest{AccessToken: token}
	resp, err := ac.client.ValidateToken(ctx, req)
	if err != nil {
		log.Printf("ValidateToken error: %v", err)
	}
	return resp, err
}

// RefreshToken вызывает gRPC-метод обновления токена
func (ac *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (*authpb.RefreshTokenResponse, error) {
	req := &authpb.RefreshTokenRequest{RefreshToken: refreshToken}
	resp, err := ac.client.RefreshToken(ctx, req)
	if err != nil {
		log.Printf("RefreshToken error: %v", err)
	}
	return resp, err
}

// Logout вызывает gRPC-метод выхода из системы
func (ac *AuthClient) Logout(ctx context.Context, refreshToken string) (*authpb.LogoutResponse, error) {
	req := &authpb.LogoutRequest{RefreshToken: refreshToken}
	resp, err := ac.client.Logout(ctx, req)
	if err != nil {
		log.Printf("Logout error: %v", err)
	}
	return resp, err
}
