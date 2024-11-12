package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"tinder-geo/internal/config"
	"tinder-geo/internal/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// $env:CONFIG_PATH = '././config/config.yaml'; go run cmd/tinder-geo/main.go

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	// code := geohash.EncodeWithPrecision(55.7893, 37.7717, 5)
	// fmt.Println(code)
	// // clean arh: app в коде Олега Козырева и видео art of development
	// // grpc server
	// // kafka client
	// // mongo client
	// // business logic
	// // observability - 2 видео art of development
	// fmt.Println("!")
}

func run() error {
	config := config.GetConfig()

	logger := setupLogger(config.Env)
	_ = logger

	if err := runGRPCServer(&config.GRPC); err != nil {
		log.Fatal(err)
	}

	return nil
}

// grpcurl -plaintext 172.24.48.1:2342 list tinder.GeoService
func runGRPCServer(config *config.GRPCConfig) error {
	log.Printf("GRPC server is running on *:%d", config.Port)

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	server.Register(grpcServer)

	err = grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

// func runGRPCServer() {
// 	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
// 	serverOptions := []grpc.ServerOption{
// 		grpc.UnaryInterceptor(interceptor.Unary()),
// 		grpc.StreamInterceptor(interceptor.Stream()),
// 	}
// 	grpcServer := grpc.NewServer(serverOptions...)
// }

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
