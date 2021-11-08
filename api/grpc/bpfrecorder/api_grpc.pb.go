// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api_bpfrecorder

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BpfRecorderClient is the client API for BpfRecorder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BpfRecorderClient interface {
	Start(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	Stop(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	SyscallsForProfile(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*SyscallsResponse, error)
}

type bpfRecorderClient struct {
	cc grpc.ClientConnInterface
}

func NewBpfRecorderClient(cc grpc.ClientConnInterface) BpfRecorderClient {
	return &bpfRecorderClient{cc}
}

func (c *bpfRecorderClient) Start(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/api_bpfrecorder.BpfRecorder/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bpfRecorderClient) Stop(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/api_bpfrecorder.BpfRecorder/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bpfRecorderClient) SyscallsForProfile(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*SyscallsResponse, error) {
	out := new(SyscallsResponse)
	err := c.cc.Invoke(ctx, "/api_bpfrecorder.BpfRecorder/SyscallsForProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BpfRecorderServer is the server API for BpfRecorder service.
// All implementations must embed UnimplementedBpfRecorderServer
// for forward compatibility
type BpfRecorderServer interface {
	Start(context.Context, *EmptyRequest) (*EmptyResponse, error)
	Stop(context.Context, *EmptyRequest) (*EmptyResponse, error)
	SyscallsForProfile(context.Context, *ProfileRequest) (*SyscallsResponse, error)
	mustEmbedUnimplementedBpfRecorderServer()
}

// UnimplementedBpfRecorderServer must be embedded to have forward compatible implementations.
type UnimplementedBpfRecorderServer struct {
}

func (UnimplementedBpfRecorderServer) Start(context.Context, *EmptyRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
}
func (UnimplementedBpfRecorderServer) Stop(context.Context, *EmptyRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (UnimplementedBpfRecorderServer) SyscallsForProfile(context.Context, *ProfileRequest) (*SyscallsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyscallsForProfile not implemented")
}
func (UnimplementedBpfRecorderServer) mustEmbedUnimplementedBpfRecorderServer() {}

// UnsafeBpfRecorderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BpfRecorderServer will
// result in compilation errors.
type UnsafeBpfRecorderServer interface {
	mustEmbedUnimplementedBpfRecorderServer()
}

func RegisterBpfRecorderServer(s grpc.ServiceRegistrar, srv BpfRecorderServer) {
	s.RegisterService(&BpfRecorder_ServiceDesc, srv)
}

func _BpfRecorder_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BpfRecorderServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api_bpfrecorder.BpfRecorder/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BpfRecorderServer).Start(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BpfRecorder_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BpfRecorderServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api_bpfrecorder.BpfRecorder/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BpfRecorderServer).Stop(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BpfRecorder_SyscallsForProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BpfRecorderServer).SyscallsForProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api_bpfrecorder.BpfRecorder/SyscallsForProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BpfRecorderServer).SyscallsForProfile(ctx, req.(*ProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BpfRecorder_ServiceDesc is the grpc.ServiceDesc for BpfRecorder service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BpfRecorder_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api_bpfrecorder.BpfRecorder",
	HandlerType: (*BpfRecorderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Start",
			Handler:    _BpfRecorder_Start_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _BpfRecorder_Stop_Handler,
		},
		{
			MethodName: "SyscallsForProfile",
			Handler:    _BpfRecorder_SyscallsForProfile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/grpc/bpfrecorder/api.proto",
}