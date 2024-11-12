package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"tinder-geo/internal/config"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	serv   *grpc.Server
	config *config.GRPCConfig
	logger *slog.Logger
}

func NewGRPCServer(config *config.GRPCConfig, logger *slog.Logger) *GRPCServer {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall,
			logging.PayloadReceived,
			logging.PayloadSent,
			logging.FinishCall,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("Recovered from panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(interceptorLogger(logger), loggingOpts...),
		selector.UnaryServerInterceptor(
			auth.UnaryServerInterceptor(authenticator),
			selector.MatchFunc(authMatcher),
		),
	))
	reflection.Register(grpcServer)

	Register(grpcServer)

	return &GRPCServer{
		serv:   grpcServer,
		config: config,
		logger: logger,
	}
}

func (s *GRPCServer) Run() error {
	s.logger.Info("GRPC server is running", slog.Int("port", s.config.Port))

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return err
	}

	err = s.serv.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (s *GRPCServer) Stop() {
	s.logger.Info("Stopping GRPC server...")
	s.serv.GracefulStop()
}

func interceptorLogger(logger *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		logger.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func authenticator(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	// TODO: This is example only, perform proper Oauth/OIDC verification!
	if token != "test" {
		return nil, status.Error(codes.Unauthenticated, "invalid auth token")
	}
	// NOTE: You can also pass the token in the context for further interceptors or gRPC service code.
	return ctx, nil
}

func authMatcher(ctx context.Context, callMeta interceptors.CallMeta) bool {
	// return healthpb.Health_ServiceDesc.ServiceName != callMeta.Service
	return true // Какие методы требуют аутентификацию?
}
