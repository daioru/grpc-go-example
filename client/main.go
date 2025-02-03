package main

import (
	"context"
	"log"

	pb "github.com/daioru/grpc-go-example/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	// Создаём gRPC-клиент
	client := pb.NewGreeterClient(cc)

	// Отправляем запрос
	response, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Alice"})
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}

	log.Println("Server response:", response.Message)
}
