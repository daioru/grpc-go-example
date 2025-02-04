package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/daioru/grpc-go-example/proto"
)

type authService struct {
	pb.UnimplementedAuthServiceServer
}

// Секретный ключ для подписи токена (в реале хранится в .env)
var jwtKey = []byte("my-secret-key")

func generateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	log.Println("Generated JWT Token:", signedToken)
	return signedToken, nil
}

// Заглушка с данными пользователей
var users = map[string]string{
	"admin": "password123",
	"user":  "mypassword",
}

func (s *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if password, ok := users[req.Username]; !ok || password != req.Password {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
	}

	token, err := generateJWT(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not generate token")
	}

	return &pb.LoginResponse{Token: token}, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtKey, nil
	})
}
