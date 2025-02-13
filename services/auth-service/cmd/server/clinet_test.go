package main

import (
	"context"
	"fmt"
	pb "github.com/exPriceD/Streaming-platform/services/auth-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"testing"
)

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

	// Тестируем регистрацию
	registerReq := &pb.RegisterRequest{
		Username:                "admin_tester",
		Email:                   "tester@example.com",
		Password:                "password123",
		Birthday:                timestamppb.Now(),
		Gender:                  pb.Gender_GENDER_MALE,
		ConsentToDataProcessing: true,
	}
	registerResp, err := client.Register(context.Background(), registerReq)
	if err != nil {
		log.Fatalf("Ошибка регистрации: %v", err)
	}
	fmt.Println("✅ Успешная регистрация:", registerResp)

	// Тестируем выход
	logoutReq := &pb.LogoutRequest{
		RefreshToken: registerResp.RefreshToken,
	}
	logoutResp, err := client.Logout(context.Background(), logoutReq)
	if err != nil {
		log.Fatalf("Ошибка выхода: %v", err)
	}
	fmt.Println("✅ Успешный выход:", logoutResp)

	// Тестируем логин
	loginReq := &pb.LoginRequest{
		LoginIdentifier: &pb.LoginRequest_Email{Email: "tester@example.com"},
		Password:        "password123",
	}
	loginResp, err := client.Login(context.Background(), loginReq)
	if err != nil {
		log.Fatalf("Ошибка логина: %v", err)
	}
	fmt.Println("✅ Успешный вход:", loginResp)

	// Тестируем валидацию токена
	validateReq := &pb.ValidateTokenRequest{
		AccessToken: loginResp.AccessToken,
	}
	validateResp, err := client.ValidateToken(context.Background(), validateReq)
	if err != nil {
		log.Fatalf("Ошибка валидации токена: %v", err)
	}
	fmt.Println("✅ Проверка access_token:", validateResp)

	// Тестируем обновление токена
	refreshReq := &pb.RefreshTokenRequest{
		RefreshToken: loginResp.RefreshToken,
	}
	refreshResp, err := client.RefreshToken(context.Background(), refreshReq)
	if err != nil {
		log.Fatalf("Ошибка обновления токена: %v", err)
	}
	fmt.Println("✅ Обновление access_token:", refreshResp)
}
