// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0--rc1
// source: geo_service.proto

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
	GeoService_GetProfilesByLocation_FullMethodName = "/tinder.GeoService/GetProfilesByLocation"
	GeoService_ChangeLocation_FullMethodName        = "/tinder.GeoService/ChangeLocation"
)

// GeoServiceClient is the client API for GeoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GeoServiceClient interface {
	GetProfilesByLocation(ctx context.Context, in *GetProfilesByLocationRequest, opts ...grpc.CallOption) (*GetProfilesByLocationResponse, error)
	ChangeLocation(ctx context.Context, in *ChangeLocationRequest, opts ...grpc.CallOption) (*ChangeLocationResponse, error)
}

type geoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGeoServiceClient(cc grpc.ClientConnInterface) GeoServiceClient {
	return &geoServiceClient{cc}
}

func (c *geoServiceClient) GetProfilesByLocation(ctx context.Context, in *GetProfilesByLocationRequest, opts ...grpc.CallOption) (*GetProfilesByLocationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProfilesByLocationResponse)
	err := c.cc.Invoke(ctx, GeoService_GetProfilesByLocation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *geoServiceClient) ChangeLocation(ctx context.Context, in *ChangeLocationRequest, opts ...grpc.CallOption) (*ChangeLocationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ChangeLocationResponse)
	err := c.cc.Invoke(ctx, GeoService_ChangeLocation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GeoServiceServer is the server API for GeoService service.
// All implementations must embed UnimplementedGeoServiceServer
// for forward compatibility.
type GeoServiceServer interface {
	GetProfilesByLocation(context.Context, *GetProfilesByLocationRequest) (*GetProfilesByLocationResponse, error)
	ChangeLocation(context.Context, *ChangeLocationRequest) (*ChangeLocationResponse, error)
	mustEmbedUnimplementedGeoServiceServer()
}

// UnimplementedGeoServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGeoServiceServer struct{}

func (UnimplementedGeoServiceServer) GetProfilesByLocation(context.Context, *GetProfilesByLocationRequest) (*GetProfilesByLocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfilesByLocation not implemented")
}
func (UnimplementedGeoServiceServer) ChangeLocation(context.Context, *ChangeLocationRequest) (*ChangeLocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeLocation not implemented")
}
func (UnimplementedGeoServiceServer) mustEmbedUnimplementedGeoServiceServer() {}
func (UnimplementedGeoServiceServer) testEmbeddedByValue()                    {}

// UnsafeGeoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GeoServiceServer will
// result in compilation errors.
type UnsafeGeoServiceServer interface {
	mustEmbedUnimplementedGeoServiceServer()
}

func RegisterGeoServiceServer(s grpc.ServiceRegistrar, srv GeoServiceServer) {
	// If the following call pancis, it indicates UnimplementedGeoServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GeoService_ServiceDesc, srv)
}

func _GeoService_GetProfilesByLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfilesByLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoServiceServer).GetProfilesByLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GeoService_GetProfilesByLocation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoServiceServer).GetProfilesByLocation(ctx, req.(*GetProfilesByLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GeoService_ChangeLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoServiceServer).ChangeLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GeoService_ChangeLocation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoServiceServer).ChangeLocation(ctx, req.(*ChangeLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GeoService_ServiceDesc is the grpc.ServiceDesc for GeoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GeoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tinder.GeoService",
	HandlerType: (*GeoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProfilesByLocation",
			Handler:    _GeoService_GetProfilesByLocation_Handler,
		},
		{
			MethodName: "ChangeLocation",
			Handler:    _GeoService_ChangeLocation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "geo_service.proto",
}
