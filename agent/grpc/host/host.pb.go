// Code generated by protoc-gen-go.
// source: host.proto
// DO NOT EDIT!

/*
Package host is a generated protocol buffer package.

It is generated from these files:
	host.proto

It has these top-level messages:
	StartRequest
	StartResponse
	StopRequest
	StopResponse
	StatusRequest
	StatusResponse
*/
package host

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StartRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *StartRequest) Reset()                    { *m = StartRequest{} }
func (m *StartRequest) String() string            { return proto.CompactTextString(m) }
func (*StartRequest) ProtoMessage()               {}
func (*StartRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type StartResponse struct {
}

func (m *StartResponse) Reset()                    { *m = StartResponse{} }
func (m *StartResponse) String() string            { return proto.CompactTextString(m) }
func (*StartResponse) ProtoMessage()               {}
func (*StartResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type StopRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *StopRequest) Reset()                    { *m = StopRequest{} }
func (m *StopRequest) String() string            { return proto.CompactTextString(m) }
func (*StopRequest) ProtoMessage()               {}
func (*StopRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type StopResponse struct {
}

func (m *StopResponse) Reset()                    { *m = StopResponse{} }
func (m *StopResponse) String() string            { return proto.CompactTextString(m) }
func (*StopResponse) ProtoMessage()               {}
func (*StopResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type StatusRequest struct {
}

func (m *StatusRequest) Reset()                    { *m = StatusRequest{} }
func (m *StatusRequest) String() string            { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()               {}
func (*StatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type StatusResponse struct {
	Name   string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
}

func (m *StatusResponse) Reset()                    { *m = StatusResponse{} }
func (m *StatusResponse) String() string            { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()               {}
func (*StatusResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func init() {
	proto.RegisterType((*StartRequest)(nil), "host.StartRequest")
	proto.RegisterType((*StartResponse)(nil), "host.StartResponse")
	proto.RegisterType((*StopRequest)(nil), "host.StopRequest")
	proto.RegisterType((*StopResponse)(nil), "host.StopResponse")
	proto.RegisterType((*StatusRequest)(nil), "host.StatusRequest")
	proto.RegisterType((*StatusResponse)(nil), "host.StatusResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Host service

type HostClient interface {
	Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error)
	Stop(ctx context.Context, in *StopRequest, opts ...grpc.CallOption) (*StopResponse, error)
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (Host_StatusClient, error)
}

type hostClient struct {
	cc *grpc.ClientConn
}

func NewHostClient(cc *grpc.ClientConn) HostClient {
	return &hostClient{cc}
}

func (c *hostClient) Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error) {
	out := new(StartResponse)
	err := grpc.Invoke(ctx, "/host.Host/Start", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hostClient) Stop(ctx context.Context, in *StopRequest, opts ...grpc.CallOption) (*StopResponse, error) {
	out := new(StopResponse)
	err := grpc.Invoke(ctx, "/host.Host/Stop", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hostClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (Host_StatusClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Host_serviceDesc.Streams[0], c.cc, "/host.Host/Status", opts...)
	if err != nil {
		return nil, err
	}
	x := &hostStatusClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Host_StatusClient interface {
	Recv() (*StatusResponse, error)
	grpc.ClientStream
}

type hostStatusClient struct {
	grpc.ClientStream
}

func (x *hostStatusClient) Recv() (*StatusResponse, error) {
	m := new(StatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Host service

type HostServer interface {
	Start(context.Context, *StartRequest) (*StartResponse, error)
	Stop(context.Context, *StopRequest) (*StopResponse, error)
	Status(*StatusRequest, Host_StatusServer) error
}

func RegisterHostServer(s *grpc.Server, srv HostServer) {
	s.RegisterService(&_Host_serviceDesc, srv)
}

func _Host_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HostServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/host.Host/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HostServer).Start(ctx, req.(*StartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Host_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HostServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/host.Host/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HostServer).Stop(ctx, req.(*StopRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Host_Status_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StatusRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(HostServer).Status(m, &hostStatusServer{stream})
}

type Host_StatusServer interface {
	Send(*StatusResponse) error
	grpc.ServerStream
}

type hostStatusServer struct {
	grpc.ServerStream
}

func (x *hostStatusServer) Send(m *StatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _Host_serviceDesc = grpc.ServiceDesc{
	ServiceName: "host.Host",
	HandlerType: (*HostServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Start",
			Handler:    _Host_Start_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _Host_Stop_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Status",
			Handler:       _Host_Status_Handler,
			ServerStreams: true,
		},
	},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("host.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 203 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x90, 0x4d, 0x0e, 0x82, 0x30,
	0x14, 0x84, 0xad, 0xa9, 0x24, 0x3e, 0x15, 0xe3, 0xd3, 0x18, 0xc2, 0x4a, 0xbb, 0x72, 0x85, 0x06,
	0x17, 0x6e, 0x3c, 0x80, 0x6b, 0x38, 0x01, 0x26, 0x4d, 0xdc, 0x48, 0x91, 0x3e, 0x6e, 0xe4, 0x41,
	0xb5, 0xa5, 0xfc, 0x25, 0xc6, 0xdd, 0x63, 0x98, 0x99, 0x7c, 0x53, 0x80, 0x87, 0xd2, 0x14, 0x15,
	0xa5, 0x22, 0x85, 0xdc, 0xdc, 0x42, 0xc0, 0x3c, 0xa5, 0xac, 0xa4, 0x44, 0xbe, 0x2a, 0xa9, 0x09,
	0x11, 0x78, 0x9e, 0x3d, 0x65, 0xc0, 0x76, 0xec, 0x30, 0x4d, 0xec, 0x2d, 0x96, 0xb0, 0x70, 0x1e,
	0x5d, 0xa8, 0x5c, 0x4b, 0xb1, 0x87, 0x59, 0x4a, 0xaa, 0xf8, 0x97, 0xf1, 0x4d, 0xaf, 0xb1, 0xb8,
	0x48, 0xdd, 0x41, 0x95, 0x76, 0x21, 0x71, 0x05, 0xbf, 0x11, 0x6a, 0xcb, 0xaf, 0x1a, 0xdc, 0x82,
	0xa7, 0xad, 0x2b, 0x18, 0x5b, 0xd5, 0x7d, 0xc5, 0x6f, 0x06, 0xfc, 0xf6, 0xe5, 0xc7, 0x18, 0x26,
	0x96, 0x0d, 0x31, 0xb2, 0xdb, 0xfa, 0x63, 0xc2, 0xf5, 0x40, 0x73, 0x24, 0x23, 0x3c, 0x02, 0x37,
	0x6c, 0xb8, 0x6a, 0x7e, 0xb7, 0x53, 0x42, 0xec, 0x4b, 0x6d, 0xe0, 0x02, 0x5e, 0xcd, 0x8a, 0x5d,
	0x63, 0x37, 0x25, 0xdc, 0x0c, 0xc5, 0x26, 0x76, 0x62, 0x77, 0xcf, 0x3e, 0xf5, 0xf9, 0x13, 0x00,
	0x00, 0xff, 0xff, 0x1c, 0xd9, 0xa4, 0xc6, 0x78, 0x01, 0x00, 0x00,
}