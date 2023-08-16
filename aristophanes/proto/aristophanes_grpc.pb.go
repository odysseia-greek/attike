// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/aristophanes.proto

package proto

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

// TraceServiceClient is the client API for TraceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TraceServiceClient interface {
	StartTrace(ctx context.Context, in *StartTraceRequest, opts ...grpc.CallOption) (*TraceResponse, error)
	Trace(ctx context.Context, in *TraceRequest, opts ...grpc.CallOption) (*TraceResponse, error)
	StartNewSpan(ctx context.Context, in *StartSpanRequest, opts ...grpc.CallOption) (*TraceResponse, error)
	Span(ctx context.Context, in *SpanRequest, opts ...grpc.CallOption) (*TraceResponse, error)
	DatabaseSpan(ctx context.Context, in *DatabaseSpanRequest, opts ...grpc.CallOption) (*TraceResponse, error)
	CloseTrace(ctx context.Context, in *CloseTraceRequest, opts ...grpc.CallOption) (*TraceResponse, error)
}

type traceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTraceServiceClient(cc grpc.ClientConnInterface) TraceServiceClient {
	return &traceServiceClient{cc}
}

func (c *traceServiceClient) StartTrace(ctx context.Context, in *StartTraceRequest, opts ...grpc.CallOption) (*TraceResponse, error) {
	out := new(TraceResponse)
	err := c.cc.Invoke(ctx, "/proto.TraceService/StartTrace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceServiceClient) Trace(ctx context.Context, in *TraceRequest, opts ...grpc.CallOption) (*TraceResponse, error) {
	out := new(TraceResponse)
	err := c.cc.Invoke(ctx, "/proto.TraceService/Trace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceServiceClient) StartNewSpan(ctx context.Context, in *StartSpanRequest, opts ...grpc.CallOption) (*TraceResponse, error) {
	out := new(TraceResponse)
	err := c.cc.Invoke(ctx, "/proto.TraceService/StartNewSpan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceServiceClient) Span(ctx context.Context, in *SpanRequest, opts ...grpc.CallOption) (*TraceResponse, error) {
	out := new(TraceResponse)
	err := c.cc.Invoke(ctx, "/proto.TraceService/Span", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceServiceClient) DatabaseSpan(ctx context.Context, in *DatabaseSpanRequest, opts ...grpc.CallOption) (*TraceResponse, error) {
	out := new(TraceResponse)
	err := c.cc.Invoke(ctx, "/proto.TraceService/DatabaseSpan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceServiceClient) CloseTrace(ctx context.Context, in *CloseTraceRequest, opts ...grpc.CallOption) (*TraceResponse, error) {
	out := new(TraceResponse)
	err := c.cc.Invoke(ctx, "/proto.TraceService/CloseTrace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TraceServiceServer is the server API for TraceService service.
// All implementations must embed UnimplementedTraceServiceServer
// for forward compatibility
type TraceServiceServer interface {
	StartTrace(context.Context, *StartTraceRequest) (*TraceResponse, error)
	Trace(context.Context, *TraceRequest) (*TraceResponse, error)
	StartNewSpan(context.Context, *StartSpanRequest) (*TraceResponse, error)
	Span(context.Context, *SpanRequest) (*TraceResponse, error)
	DatabaseSpan(context.Context, *DatabaseSpanRequest) (*TraceResponse, error)
	CloseTrace(context.Context, *CloseTraceRequest) (*TraceResponse, error)
	mustEmbedUnimplementedTraceServiceServer()
}

// UnimplementedTraceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTraceServiceServer struct {
}

func (UnimplementedTraceServiceServer) StartTrace(context.Context, *StartTraceRequest) (*TraceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartTrace not implemented")
}
func (UnimplementedTraceServiceServer) Trace(context.Context, *TraceRequest) (*TraceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Trace not implemented")
}
func (UnimplementedTraceServiceServer) StartNewSpan(context.Context, *StartSpanRequest) (*TraceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartNewSpan not implemented")
}
func (UnimplementedTraceServiceServer) Span(context.Context, *SpanRequest) (*TraceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Span not implemented")
}
func (UnimplementedTraceServiceServer) DatabaseSpan(context.Context, *DatabaseSpanRequest) (*TraceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DatabaseSpan not implemented")
}
func (UnimplementedTraceServiceServer) CloseTrace(context.Context, *CloseTraceRequest) (*TraceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloseTrace not implemented")
}
func (UnimplementedTraceServiceServer) mustEmbedUnimplementedTraceServiceServer() {}

// UnsafeTraceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TraceServiceServer will
// result in compilation errors.
type UnsafeTraceServiceServer interface {
	mustEmbedUnimplementedTraceServiceServer()
}

func RegisterTraceServiceServer(s grpc.ServiceRegistrar, srv TraceServiceServer) {
	s.RegisterService(&TraceService_ServiceDesc, srv)
}

func _TraceService_StartTrace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartTraceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceServiceServer).StartTrace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TraceService/StartTrace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceServiceServer).StartTrace(ctx, req.(*StartTraceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceService_Trace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TraceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceServiceServer).Trace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TraceService/Trace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceServiceServer).Trace(ctx, req.(*TraceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceService_StartNewSpan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartSpanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceServiceServer).StartNewSpan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TraceService/StartNewSpan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceServiceServer).StartNewSpan(ctx, req.(*StartSpanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceService_Span_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SpanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceServiceServer).Span(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TraceService/Span",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceServiceServer).Span(ctx, req.(*SpanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceService_DatabaseSpan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DatabaseSpanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceServiceServer).DatabaseSpan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TraceService/DatabaseSpan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceServiceServer).DatabaseSpan(ctx, req.(*DatabaseSpanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceService_CloseTrace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseTraceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceServiceServer).CloseTrace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TraceService/CloseTrace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceServiceServer).CloseTrace(ctx, req.(*CloseTraceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TraceService_ServiceDesc is the grpc.ServiceDesc for TraceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TraceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TraceService",
	HandlerType: (*TraceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartTrace",
			Handler:    _TraceService_StartTrace_Handler,
		},
		{
			MethodName: "Trace",
			Handler:    _TraceService_Trace_Handler,
		},
		{
			MethodName: "StartNewSpan",
			Handler:    _TraceService_StartNewSpan_Handler,
		},
		{
			MethodName: "Span",
			Handler:    _TraceService_Span_Handler,
		},
		{
			MethodName: "DatabaseSpan",
			Handler:    _TraceService_DatabaseSpan_Handler,
		},
		{
			MethodName: "CloseTrace",
			Handler:    _TraceService_CloseTrace_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/aristophanes.proto",
}