// Code generated by protoc-gen-go. DO NOT EDIT.
// source: captain.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type InstallChartRequest struct {
	Chart *Chart `protobuf:"bytes,1,opt,name=chart,proto3" json:"chart,omitempty"`
	// timeout specifies the max amount of time any caplet can run.
	Timeout int64 `protobuf:"varint,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
	// sync forces all node to pull image
	Sync                 bool     `protobuf:"varint,3,opt,name=sync,proto3" json:"sync,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstallChartRequest) Reset()         { *m = InstallChartRequest{} }
func (m *InstallChartRequest) String() string { return proto.CompactTextString(m) }
func (*InstallChartRequest) ProtoMessage()    {}
func (*InstallChartRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_75f80b18f5aaa93a, []int{0}
}

func (m *InstallChartRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallChartRequest.Unmarshal(m, b)
}
func (m *InstallChartRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallChartRequest.Marshal(b, m, deterministic)
}
func (m *InstallChartRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallChartRequest.Merge(m, src)
}
func (m *InstallChartRequest) XXX_Size() int {
	return xxx_messageInfo_InstallChartRequest.Size(m)
}
func (m *InstallChartRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallChartRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InstallChartRequest proto.InternalMessageInfo

func (m *InstallChartRequest) GetChart() *Chart {
	if m != nil {
		return m.Chart
	}
	return nil
}

func (m *InstallChartRequest) GetTimeout() int64 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *InstallChartRequest) GetSync() bool {
	if m != nil {
		return m.Sync
	}
	return false
}

type InstallChartResponse struct {
	Out                  string   `protobuf:"bytes,1,opt,name=out,proto3" json:"out,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstallChartResponse) Reset()         { *m = InstallChartResponse{} }
func (m *InstallChartResponse) String() string { return proto.CompactTextString(m) }
func (*InstallChartResponse) ProtoMessage()    {}
func (*InstallChartResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_75f80b18f5aaa93a, []int{1}
}

func (m *InstallChartResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallChartResponse.Unmarshal(m, b)
}
func (m *InstallChartResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallChartResponse.Marshal(b, m, deterministic)
}
func (m *InstallChartResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallChartResponse.Merge(m, src)
}
func (m *InstallChartResponse) XXX_Size() int {
	return xxx_messageInfo_InstallChartResponse.Size(m)
}
func (m *InstallChartResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallChartResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InstallChartResponse proto.InternalMessageInfo

func (m *InstallChartResponse) GetOut() string {
	if m != nil {
		return m.Out
	}
	return ""
}

type GenerateScriptRequest struct {
	Chart                *Chart   `protobuf:"bytes,1,opt,name=chart,proto3" json:"chart,omitempty"`
	Pull                 bool     `protobuf:"varint,2,opt,name=pull,proto3" json:"pull,omitempty"`
	Retag                *ReTag   `protobuf:"bytes,3,opt,name=retag,proto3" json:"retag,omitempty"`
	Save                 bool     `protobuf:"varint,4,opt,name=save,proto3" json:"save,omitempty"`
	Load                 bool     `protobuf:"varint,5,opt,name=load,proto3" json:"load,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenerateScriptRequest) Reset()         { *m = GenerateScriptRequest{} }
func (m *GenerateScriptRequest) String() string { return proto.CompactTextString(m) }
func (*GenerateScriptRequest) ProtoMessage()    {}
func (*GenerateScriptRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_75f80b18f5aaa93a, []int{2}
}

func (m *GenerateScriptRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateScriptRequest.Unmarshal(m, b)
}
func (m *GenerateScriptRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateScriptRequest.Marshal(b, m, deterministic)
}
func (m *GenerateScriptRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateScriptRequest.Merge(m, src)
}
func (m *GenerateScriptRequest) XXX_Size() int {
	return xxx_messageInfo_GenerateScriptRequest.Size(m)
}
func (m *GenerateScriptRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateScriptRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateScriptRequest proto.InternalMessageInfo

func (m *GenerateScriptRequest) GetChart() *Chart {
	if m != nil {
		return m.Chart
	}
	return nil
}

func (m *GenerateScriptRequest) GetPull() bool {
	if m != nil {
		return m.Pull
	}
	return false
}

func (m *GenerateScriptRequest) GetRetag() *ReTag {
	if m != nil {
		return m.Retag
	}
	return nil
}

func (m *GenerateScriptRequest) GetSave() bool {
	if m != nil {
		return m.Save
	}
	return false
}

func (m *GenerateScriptRequest) GetLoad() bool {
	if m != nil {
		return m.Load
	}
	return false
}

type GenerateScriptResponse struct {
	Out                  string   `protobuf:"bytes,1,opt,name=out,proto3" json:"out,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenerateScriptResponse) Reset()         { *m = GenerateScriptResponse{} }
func (m *GenerateScriptResponse) String() string { return proto.CompactTextString(m) }
func (*GenerateScriptResponse) ProtoMessage()    {}
func (*GenerateScriptResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_75f80b18f5aaa93a, []int{3}
}

func (m *GenerateScriptResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateScriptResponse.Unmarshal(m, b)
}
func (m *GenerateScriptResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateScriptResponse.Marshal(b, m, deterministic)
}
func (m *GenerateScriptResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateScriptResponse.Merge(m, src)
}
func (m *GenerateScriptResponse) XXX_Size() int {
	return xxx_messageInfo_GenerateScriptResponse.Size(m)
}
func (m *GenerateScriptResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateScriptResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateScriptResponse proto.InternalMessageInfo

func (m *GenerateScriptResponse) GetOut() string {
	if m != nil {
		return m.Out
	}
	return ""
}

func init() {
	proto.RegisterType((*InstallChartRequest)(nil), "proto.InstallChartRequest")
	proto.RegisterType((*InstallChartResponse)(nil), "proto.InstallChartResponse")
	proto.RegisterType((*GenerateScriptRequest)(nil), "proto.GenerateScriptRequest")
	proto.RegisterType((*GenerateScriptResponse)(nil), "proto.GenerateScriptResponse")
}

func init() { proto.RegisterFile("captain.proto", fileDescriptor_75f80b18f5aaa93a) }

var fileDescriptor_75f80b18f5aaa93a = []byte{
	// 290 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xcf, 0x4e, 0xc3, 0x30,
	0x0c, 0xc6, 0x09, 0x6d, 0xd9, 0xf0, 0x06, 0x42, 0xe6, 0x8f, 0xa2, 0x02, 0x52, 0x95, 0x53, 0xc5,
	0xa1, 0x87, 0xf1, 0x08, 0x3b, 0xa0, 0x9d, 0x90, 0x02, 0x2f, 0x10, 0x8a, 0x55, 0x2a, 0x95, 0xb6,
	0x34, 0xe9, 0x24, 0x9e, 0x85, 0x33, 0xef, 0x89, 0x92, 0x74, 0x82, 0xa1, 0xee, 0xb0, 0x53, 0x3f,
	0xfb, 0xb3, 0xe3, 0x9f, 0x5d, 0x38, 0xc9, 0x55, 0x6b, 0x54, 0x59, 0x67, 0x6d, 0xd7, 0x98, 0x06,
	0x23, 0xf7, 0x89, 0x67, 0x1d, 0x19, 0x55, 0x64, 0x43, 0x90, 0xbf, 0xa9, 0xce, 0xf8, 0x40, 0x14,
	0x70, 0xbe, 0xaa, 0xb5, 0x51, 0x55, 0xb5, 0xb4, 0x59, 0x49, 0x1f, 0x3d, 0x69, 0x83, 0x02, 0x22,
	0x57, 0xc5, 0x59, 0xc2, 0xd2, 0xd9, 0x62, 0xee, 0xab, 0x33, 0x5f, 0xe3, 0x2d, 0xe4, 0x30, 0x31,
	0xe5, 0x3b, 0x35, 0xbd, 0xe1, 0x87, 0x09, 0x4b, 0x03, 0xb9, 0x09, 0x11, 0x21, 0xd4, 0x9f, 0x75,
	0xce, 0x83, 0x84, 0xa5, 0x53, 0xe9, 0xb4, 0x48, 0xe1, 0x62, 0x7b, 0x90, 0x6e, 0x9b, 0x5a, 0x13,
	0x9e, 0x41, 0x60, 0x5f, 0xb0, 0x73, 0x8e, 0xa5, 0x95, 0xe2, 0x8b, 0xc1, 0xe5, 0x03, 0xd5, 0xd4,
	0x29, 0x43, 0x4f, 0x79, 0x57, 0xb6, 0x7b, 0x51, 0x21, 0x84, 0x6d, 0x5f, 0x55, 0x0e, 0x69, 0x2a,
	0x9d, 0xb6, 0x7d, 0xee, 0x00, 0x0e, 0xe8, 0xb7, 0x4f, 0xd2, 0xb3, 0x2a, 0xa4, 0xb7, 0x1c, 0xb3,
	0x5a, 0x13, 0x0f, 0x07, 0x66, 0xb5, 0x26, 0x9b, 0xab, 0x1a, 0xf5, 0xca, 0x23, 0x9f, 0xb3, 0x5a,
	0xdc, 0xc1, 0xd5, 0x7f, 0xb8, 0x5d, 0x9b, 0x2c, 0xbe, 0x19, 0x4c, 0x96, 0xfe, 0x7f, 0xe0, 0x0a,
	0xe6, 0x7f, 0xf7, 0xc7, 0x78, 0x80, 0x18, 0xb9, 0x7e, 0x7c, 0x3d, 0xea, 0xf9, 0x31, 0xe2, 0x00,
	0x1f, 0xe1, 0x74, 0x1b, 0x01, 0x6f, 0x86, 0x86, 0xd1, 0xb3, 0xc5, 0xb7, 0x3b, 0xdc, 0xcd, 0x83,
	0x2f, 0x47, 0xce, 0xbf, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff, 0x70, 0xe8, 0x45, 0x48, 0x3d, 0x02,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CaptainClient is the client API for Captain service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CaptainClient interface {
	InstallChart(ctx context.Context, in *InstallChartRequest, opts ...grpc.CallOption) (*InstallChartResponse, error)
	GenerateScript(ctx context.Context, in *GenerateScriptRequest, opts ...grpc.CallOption) (*GenerateScriptResponse, error)
}

type captainClient struct {
	cc *grpc.ClientConn
}

func NewCaptainClient(cc *grpc.ClientConn) CaptainClient {
	return &captainClient{cc}
}

func (c *captainClient) InstallChart(ctx context.Context, in *InstallChartRequest, opts ...grpc.CallOption) (*InstallChartResponse, error) {
	out := new(InstallChartResponse)
	err := c.cc.Invoke(ctx, "/proto.Captain/InstallChart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *captainClient) GenerateScript(ctx context.Context, in *GenerateScriptRequest, opts ...grpc.CallOption) (*GenerateScriptResponse, error) {
	out := new(GenerateScriptResponse)
	err := c.cc.Invoke(ctx, "/proto.Captain/GenerateScript", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CaptainServer is the server API for Captain service.
type CaptainServer interface {
	InstallChart(context.Context, *InstallChartRequest) (*InstallChartResponse, error)
	GenerateScript(context.Context, *GenerateScriptRequest) (*GenerateScriptResponse, error)
}

func RegisterCaptainServer(s *grpc.Server, srv CaptainServer) {
	s.RegisterService(&_Captain_serviceDesc, srv)
}

func _Captain_InstallChart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstallChartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaptainServer).InstallChart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Captain/InstallChart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaptainServer).InstallChart(ctx, req.(*InstallChartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Captain_GenerateScript_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateScriptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaptainServer).GenerateScript(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Captain/GenerateScript",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaptainServer).GenerateScript(ctx, req.(*GenerateScriptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Captain_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Captain",
	HandlerType: (*CaptainServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InstallChart",
			Handler:    _Captain_InstallChart_Handler,
		},
		{
			MethodName: "GenerateScript",
			Handler:    _Captain_GenerateScript_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "captain.proto",
}