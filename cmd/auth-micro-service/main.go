package main

import (
	"auth-micro-service/internal/di"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

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

	go func() {
		lis, err := net.Listen("tcp", container.Config().GRPCPort)
		if err != nil {
			container.Logger().Fatal(error.Error(err))
		}

		log.Println("grpc server started on ", container.Config().GRPCPort)

		if err = grpcServer.Serve(lis); err != nil {
			container.Logger().Fatal(error.Error(err))
		}

		container.Logger().Info("grpc server started on :50051")
	}()

	<-stop

	container.Logger().Info("Shutting down server...")

	container.ShotDown()

	container.Logger().Info("server was shut down")
}
