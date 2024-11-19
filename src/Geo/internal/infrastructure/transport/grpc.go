package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"tinder-geo/internal/config"
	"tinder-geo/internal/server"

	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type userClaims struct {
	ProfileId string `json:"ProfileId"`
	Gender    string `json:"Gender"`
	jwt.StandardClaims
}

type Server struct {
	serv   *grpc.Server
	config config.GRPCConfig
	logger *slog.Logger
}

func NewServer(config config.GRPCConfig, logger *slog.Logger, service server.Service) Server {
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
			logger.Error("recovered from panic", slog.Any("panic", p))
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

	server.Register(grpcServer, service)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	return Server{
		serv:   grpcServer,
		config: config,
		logger: logger,
	}
}

func (s Server) Run() error {
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

func (s Server) GracefulStop() {
	s.serv.GracefulStop()
	s.logger.Info("GRPC server stopped")
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
	parsedToken, _ := jwt.ParseWithClaims(token, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		// todo: move to secret file
		return []byte("fjg847sdjvnjxcFHdsag38d_d8sj3aqQwfdsph3456v0bjz45ty54gpo3vhjs7234f09Odp"), nil
	})

	claims := parsedToken.Claims.(*userClaims)

	if !parsedToken.Valid {
		return nil, status.Error(codes.Unauthenticated, "invalid auth token")
	}

	ctx = context.WithValue(ctx, server.ProfileIdCtxKey, claims.ProfileId)
	ctx = context.WithValue(ctx, server.GenderCtxKey, claims.Gender)
	return ctx, nil
}

func authMatcher(ctx context.Context, callMeta interceptors.CallMeta) bool {
	return grpc_health_v1.Health_ServiceDesc.ServiceName != callMeta.Service
}
