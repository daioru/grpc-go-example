package main

import (
	"context"
	"fmt"
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

func main() {
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
