package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

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

	// Открываем двусторонний поток
	stream, err := client.BidirectionalStreamGreetings(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	done := make(chan struct{})

	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Println("Server closed the connection.")
					close(done)
					return
				}
				log.Fatalf("error recieving response: %v", err)
			}
			fmt.Println("Server response:", resp.Message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter names (type 'exit' to quit):")
	for scanner.Scan() {
		name := scanner.Text()
		if name == "exit" {
			break
		}

		err := stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("error sending message: %v", err)
		}
	}

	stream.CloseSend()
	<-done
}
