// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0--rc1
// source: reaction_service.proto

package tinderpbv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ReactionService_Like_FullMethodName    = "/tinder.ReactionService/Like"
	ReactionService_Dislike_FullMethodName = "/tinder.ReactionService/Dislike"
)

// ReactionServiceClient is the client API for ReactionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReactionServiceClient interface {
	Like(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*LikeResponse, error)
	Dislike(ctx context.Context, in *DislikeRequest, opts ...grpc.CallOption) (*DislikeResponse, error)
}

type reactionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReactionServiceClient(cc grpc.ClientConnInterface) ReactionServiceClient {
	return &reactionServiceClient{cc}
}

func (c *reactionServiceClient) Like(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*LikeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LikeResponse)
	err := c.cc.Invoke(ctx, ReactionService_Like_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reactionServiceClient) Dislike(ctx context.Context, in *DislikeRequest, opts ...grpc.CallOption) (*DislikeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DislikeResponse)
	err := c.cc.Invoke(ctx, ReactionService_Dislike_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReactionServiceServer is the server API for ReactionService service.
// All implementations must embed UnimplementedReactionServiceServer
// for forward compatibility.
type ReactionServiceServer interface {
	Like(context.Context, *LikeRequest) (*LikeResponse, error)
	Dislike(context.Context, *DislikeRequest) (*DislikeResponse, error)
	mustEmbedUnimplementedReactionServiceServer()
}

// UnimplementedReactionServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedReactionServiceServer struct{}

func (UnimplementedReactionServiceServer) Like(context.Context, *LikeRequest) (*LikeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Like not implemented")
}
func (UnimplementedReactionServiceServer) Dislike(context.Context, *DislikeRequest) (*DislikeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Dislike not implemented")
}
func (UnimplementedReactionServiceServer) mustEmbedUnimplementedReactionServiceServer() {}
func (UnimplementedReactionServiceServer) testEmbeddedByValue()                         {}

// UnsafeReactionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReactionServiceServer will
// result in compilation errors.
type UnsafeReactionServiceServer interface {
	mustEmbedUnimplementedReactionServiceServer()
}

func RegisterReactionServiceServer(s grpc.ServiceRegistrar, srv ReactionServiceServer) {
	// If the following call pancis, it indicates UnimplementedReactionServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ReactionService_ServiceDesc, srv)
}

func _ReactionService_Like_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReactionServiceServer).Like(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReactionService_Like_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReactionServiceServer).Like(ctx, req.(*LikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReactionService_Dislike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DislikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReactionServiceServer).Dislike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReactionService_Dislike_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReactionServiceServer).Dislike(ctx, req.(*DislikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReactionService_ServiceDesc is the grpc.ServiceDesc for ReactionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReactionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tinder.ReactionService",
	HandlerType: (*ReactionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Like",
			Handler:    _ReactionService_Like_Handler,
		},
		{
			MethodName: "Dislike",
			Handler:    _ReactionService_Dislike_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "reaction_service.proto",
}