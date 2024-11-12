package server

// import (
// 	"context"
// 	"log"

// 	"google.golang.org/grpc"
// )

// type AuthInterceptor struct {
// }

// func NewAuthInterceptor() *AuthInterceptor {
// 	return &AuthInterceptor{}
// }

// func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
// 	return func(
// 		ctx context.Context,
// 		req interface{},
// 		info *grpc.UnaryServerInfo,
// 		handler grpc.UnaryHandler,
// 	) (interface{}, error) {
// 		log.Println("--> unary interceptor: ", info.FullMethod)

// 		// err := interceptor.authorize(ctx, info.FullMethod)
// 		// if err != nil {
// 		// 	return nil, err
// 		// }

// 		return handler(ctx, req)
// 	}
// }
