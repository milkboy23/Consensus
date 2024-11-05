// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: Consensus.proto

package __

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
	ConsensusPeer_PassCriticalPassPermission_FullMethodName = "/ConsensusPeer/PassCriticalPassPermission"
)

// ConsensusPeerClient is the client API for ConsensusPeer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConsensusPeerClient interface {
	PassCriticalPassPermission(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type consensusPeerClient struct {
	cc grpc.ClientConnInterface
}

func NewConsensusPeerClient(cc grpc.ClientConnInterface) ConsensusPeerClient {
	return &consensusPeerClient{cc}
}

func (c *consensusPeerClient) PassCriticalPassPermission(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, ConsensusPeer_PassCriticalPassPermission_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConsensusPeerServer is the server API for ConsensusPeer service.
// All implementations must embed UnimplementedConsensusPeerServer
// for forward compatibility.
type ConsensusPeerServer interface {
	PassCriticalPassPermission(context.Context, *Empty) (*Empty, error)
	mustEmbedUnimplementedConsensusPeerServer()
}

// UnimplementedConsensusPeerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedConsensusPeerServer struct{}

func (UnimplementedConsensusPeerServer) PassCriticalPassPermission(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PassCriticalPassPermission not implemented")
}
func (UnimplementedConsensusPeerServer) mustEmbedUnimplementedConsensusPeerServer() {}
func (UnimplementedConsensusPeerServer) testEmbeddedByValue()                       {}

// UnsafeConsensusPeerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConsensusPeerServer will
// result in compilation errors.
type UnsafeConsensusPeerServer interface {
	mustEmbedUnimplementedConsensusPeerServer()
}

func RegisterConsensusPeerServer(s grpc.ServiceRegistrar, srv ConsensusPeerServer) {
	// If the following call pancis, it indicates UnimplementedConsensusPeerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ConsensusPeer_ServiceDesc, srv)
}

func _ConsensusPeer_PassCriticalPassPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusPeerServer).PassCriticalPassPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConsensusPeer_PassCriticalPassPermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusPeerServer).PassCriticalPassPermission(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// ConsensusPeer_ServiceDesc is the grpc.ServiceDesc for ConsensusPeer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConsensusPeer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ConsensusPeer",
	HandlerType: (*ConsensusPeerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PassCriticalPassPermission",
			Handler:    _ConsensusPeer_PassCriticalPassPermission_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Consensus.proto",
}
