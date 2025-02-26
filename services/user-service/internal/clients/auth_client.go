package clients

import (
	"context"
	"fmt"
	authpb "github.com/exPriceD/Streaming-platform/pkg/proto/v1/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type AuthClientConfig struct {
	Address     string
	DialTimeout time.Duration
	UseTLS      bool
}

type AuthClient struct {
	conn   *grpc.ClientConn
	client authpb.AuthServiceClient
}

func NewAuthClient(cfg AuthClientConfig, opts ...grpc.DialOption) (*AuthClient, error) {
	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if cfg.UseTLS {
		// Здесь можно добавить TLS-креденшалы для продакшена
		// defaultOpts = append(defaultOpts, grpc.WithTransportCredentials(credentials.NewTLS(...)))
	}

	conn, err := grpc.NewClient(cfg.Address, append(defaultOpts, opts...)...)
	if err != nil {
		return nil, fmt.Errorf("dial auth-service at %s: %w", cfg.Address, err)
	}

	return &AuthClient{
		conn:   conn,
		client: authpb.NewAuthServiceClient(conn),
	}, nil
}

// Authenticate вызывает gRPC-метод регистрации пользователя
func (c *AuthClient) Authenticate(ctx context.Context, userId string) (string, string, error) {
	resp, err := c.client.Authenticate(ctx, &authpb.AuthenticateRequest{UserId: userId})
	if err != nil {
		return "", "", fmt.Errorf("authenticate: %w", err)
	}
	if resp.Error != nil {
		return "", "", fmt.Errorf("authenticate failed: %s", resp.Error.Message)
	}

	return resp.AccessToken, resp.RefreshToken, nil
}

// ValidateToken вызывает gRPC-метод проверки токена
func (c *AuthClient) ValidateToken(ctx context.Context, accessToken string) (bool, string, error) {
	resp, err := c.client.ValidateToken(ctx, &authpb.ValidateTokenRequest{AccessToken: accessToken})
	if err != nil {
		return false, "", fmt.Errorf("validate token: %w", err)
	}
	if resp.Error != nil {
		return false, "", nil
	}

	return resp.Valid, resp.UserId, nil
}

// RefreshToken вызывает gRPC-метод обновления токена
func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	resp, err := c.client.RefreshToken(ctx, &authpb.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		return "", "", fmt.Errorf("refresh token: %w", err)
	}
	if resp.Error != nil {
		return "", "", fmt.Errorf("refresh token failed: %s", resp.Error.Message)
	}

	return resp.AccessToken, resp.RefreshToken, nil
}

// Logout вызывает gRPC-метод выхода из системы
func (c *AuthClient) Logout(ctx context.Context, refreshToken string) (bool, error) {
	resp, err := c.client.Logout(ctx, &authpb.LogoutRequest{RefreshToken: refreshToken})
	if err != nil {
		return false, fmt.Errorf("logout: %w", err)
	}
	if resp.Error != nil {
		return false, nil
	}

	return resp.Success, nil
}

// Close закрывает gRPC-соединение
func (c *AuthClient) Close() error {
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil

	return err
}
