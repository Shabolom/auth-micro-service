package main

import (
	"auth-micro-service/internal/di"
	"context"
	"log"
	"net"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("./build/local/.env"); err != nil {
		panic(err)
	}

	container := di.New(ctx)
	container.Logger()

	grpcServer := container.NewAuthGRPCServer(
		container.Logger(),
		container.GetGRPCAuthHandlers(),
		container.GetGRPCUsersHandlers(),
	)

	lis, err := net.Listen("tcp", container.Config().GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("grpc server started on ", container.Config().GRPCPort)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

	log.Println("grpc server started on :50051")
}
