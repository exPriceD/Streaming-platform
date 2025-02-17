package main

import (
	"context"
	"fmt"
	pb "github.com/exPriceD/Streaming-platform/services/auth-service/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
)

const authServiceAddr = "localhost:50051"

func TestClient(t *testing.T) {
	// Подключаемся к gRPC-серверу
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Ошибка закрытия соединения: %v", err)
		}
	}(conn)

	client := pb.NewAuthServiceClient(conn)

	// 🔹 Шаг 1: Генерация токенов
	fmt.Println("\n🔹 Генерация токенов")

	userID := "550e8400-e29b-41d4-a716-446655440000"
	generateReq := &pb.AuthenticateRequest{UserId: userID}
	generateResp, err := runTestCase(t, "GenerateTokens", func() (*pb.AuthenticateResponse, error) {
		return client.Authenticate(context.Background(), generateReq)
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, generateResp.AccessToken)
	assert.NotEmpty(t, generateResp.RefreshToken)

	// 🔹 Шаг 2: Валидация access_token
	fmt.Println("\n🔹 Валидация access_token")
	validateReq := &pb.ValidateTokenRequest{AccessToken: generateResp.AccessToken}

	validateResp, err := runTestCase(t, "ValidateToken", func() (*pb.ValidateTokenResponse, error) {
		return client.ValidateToken(context.Background(), validateReq)
	})
	assert.NoError(t, err)
	assert.True(t, validateResp.Valid)
	assert.Equal(t, userID, validateResp.UserId)

	// 🔹 Шаг 3: Обновление токенов
	fmt.Println("\n🔹 Обновление access_token")
	refreshReq := &pb.RefreshTokenRequest{RefreshToken: generateResp.RefreshToken}

	refreshResp, err := runTestCase(t, "RefreshToken", func() (*pb.RefreshTokenResponse, error) {
		return client.RefreshToken(context.Background(), refreshReq)
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshResp.AccessToken)
	assert.NotEmpty(t, refreshResp.RefreshToken)

	// 🔹 Шаг 4: Выход (Logout)
	fmt.Println("\n🔹 Logout")
	logoutReq := &pb.LogoutRequest{RefreshToken: generateResp.RefreshToken}

	_, err = runTestCase(t, "Logout", func() (*pb.LogoutResponse, error) {
		return client.Logout(context.Background(), logoutReq)
	})
	assert.NoError(t, err)
}

func runTestCase[T any](t *testing.T, testName string, fn func() (T, error)) (T, error) {
	fmt.Printf("🔄 Тестируем %s...\n", testName)
	result, err := fn()
	if err != nil {
		fmt.Printf("❌ Ошибка в %s: %v\n", testName, err)
		t.Fatalf("❌ %s провален: %v", testName, err)
	} else {
		fmt.Printf("✅ %s прошел успешно!\n", testName)
	}
	return result, err
}
