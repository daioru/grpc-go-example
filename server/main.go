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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

// Интерсептор для проверки токена
func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not provided")
	}

	token := md["authorization"]
	if len(token) == 0 || token[0] != "Bearer my-secret-token" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
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
	creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
	}

	// Создаём gRPC-сервер с TLS
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(authInterceptor),
	)

	// Регистрируем сервис
	pb.RegisterGreeterServer(s, &greeterServer{})

	// Запускаем сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
