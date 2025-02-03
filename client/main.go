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

	// Отправляем поток сообщений
	stream, err := client.ClientStreamGreetings(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	names := []string{"Alice", "Bob", "Charlie", "Dave", "Eve"}

	for _, name := range names {
		fmt.Println("Sending:", name)
		err := stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
		time.Sleep(time.Second)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error recieving response: %v", err)
	}

	fmt.Println("Server response:", response.Message)
}
