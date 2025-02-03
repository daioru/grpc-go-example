package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/daioru/grpc-go-example/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Создаём клиентское соединение с сервером
	cc, err := grpc.NewClient(
		"127.0.0.1:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer cc.Close()

	// Создаём gRPC-клиент Greeter
	client := pb.NewGreeterClient(cc)

	// Отправляем запрос
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Alice"})
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}

	fmt.Println("Server response:", response.Message)
}
