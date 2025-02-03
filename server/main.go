package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/daioru/grpc-go-example/proto"

	"google.golang.org/grpc"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

// Unary RPC
func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Recieved request from: %s", req.Name)
	return &pb.HelloResponse{Message: "Hello, " + req.Name + "!"}, nil
}

// Server streaming RPC
func (s *greeterServer) StreamGreetings(req *pb.HelloRequest, stream pb.Greeter_StreamGreetingsServer) error {
	log.Printf("Reciever request from: %s", req.Name)

	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Hello, %s! Message #%d", req.Name, i)
		err := stream.Send(&pb.HelloResponse{Message: message})
		if err != nil {
			return fmt.Errorf("failed to send message: %v", err)
		}
		time.Sleep(time.Second) //Задержка для иммитации потоковой передачи
	}
	return nil
}

// Client streaming RPC
func (s *greeterServer) ClientStreamGreetings(stream pb.Greeter_ClientStreamGreetingsServer) error {
	log.Println("Client started streaming...")

	var names []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			response := fmt.Sprintf("Hello to all: %v!", names)
			return stream.SendAndClose(&pb.HelloResponse{Message: response})
		}
		if err != nil {
			return fmt.Errorf("error recieving stream: %v", err)
		}

		log.Printf("Recieved: %s", req.Name)
		names = append(names, req.Name)
	}
}

// Bidirectional streaming RPC
func (s *greeterServer) BidirectionalStreamGreetings(stream pb.Greeter_BidirectionalStreamGreetingsServer) error {
	log.Println("Client started bidirectional streaming...")

	for {
		// Получаем сообщение от клиента
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Client finished sending messages.")
			return nil
		}
		if err != nil {
			return fmt.Errorf("error recieving message: %v", err)
		}

		log.Printf("Recieved: %s", req.Name)

		// Отправляем ответ клиенту
		response := fmt.Sprintf("Hello, %s!", req.Name)
		err = stream.Send(&pb.HelloResponse{Message: response})
		if err != nil {
			return fmt.Errorf("error sending response: %v", err)
		}

		time.Sleep(time.Second) // Иммитация задержки ответа
	}
}

func main() {
	// Создаём слушателя на порту 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Создаём gRPC-сервер
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &greeterServer{})

	log.Println("gRPC server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
