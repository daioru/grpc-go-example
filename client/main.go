package main

import (
	"context"
	"log"

	pb "github.com/daioru/grpc-go-example/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	// Загружаем сертификат сервера
	creds, err := credentials.NewClientTLSFromFile("server.crt", "")
	if err != nil {
		log.Fatalf("Failed to load TLS certificate: %v", err)
	}

	// Подключаемся к TLS серверу
	cc, err := grpc.NewClient(
		"127.0.0.1:50051",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer cc.Close()

	// Создаём клиент AuthService
	authClient := pb.NewAuthServiceClient(cc)
	greeterClient := pb.NewGreeterClient(cc)

	loginResp, err := authClient.Login(context.Background(), &pb.LoginRequest{
		Username: "admin",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	token := loginResp.Token
	log.Println("Recieved JWT Token:", token)

	md := metadata.New(map[string]string{"authorization": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := greeterClient.SayHello(ctx, &pb.HelloRequest{Name: "John"})
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}

	log.Println("Server Response:", resp.Message)
}
