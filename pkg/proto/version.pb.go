// Code generated by protoc-gen-go. DO NOT EDIT.
// source: version.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type VersionRequest struct {
	Full                 bool     `protobuf:"varint,1,opt,name=full,proto3" json:"full,omitempty"`
	Timeout              int64    `protobuf:"varint,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Color                bool     `protobuf:"varint,3,opt,name=color,proto3" json:"color,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionRequest) Reset()         { *m = VersionRequest{} }
func (m *VersionRequest) String() string { return proto.CompactTextString(m) }
func (*VersionRequest) ProtoMessage()    {}
func (*VersionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d2c07d79758f814, []int{0}
}

func (m *VersionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionRequest.Unmarshal(m, b)
}
func (m *VersionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionRequest.Marshal(b, m, deterministic)
}
func (m *VersionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionRequest.Merge(m, src)
}
func (m *VersionRequest) XXX_Size() int {
	return xxx_messageInfo_VersionRequest.Size(m)
}
func (m *VersionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VersionRequest proto.InternalMessageInfo

func (m *VersionRequest) GetFull() bool {
	if m != nil {
		return m.Full
	}
	return false
}

func (m *VersionRequest) GetTimeout() int64 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *VersionRequest) GetColor() bool {
	if m != nil {
		return m.Color
	}
	return false
}

type VersionResponse struct {
	Hostname             string   `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionResponse) Reset()         { *m = VersionResponse{} }
func (m *VersionResponse) String() string { return proto.CompactTextString(m) }
func (*VersionResponse) ProtoMessage()    {}
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d2c07d79758f814, []int{1}
}

func (m *VersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionResponse.Unmarshal(m, b)
}
func (m *VersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionResponse.Marshal(b, m, deterministic)
}
func (m *VersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionResponse.Merge(m, src)
}
func (m *VersionResponse) XXX_Size() int {
	return xxx_messageInfo_VersionResponse.Size(m)
}
func (m *VersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VersionResponse proto.InternalMessageInfo

func (m *VersionResponse) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *VersionResponse) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func init() {
	proto.RegisterType((*VersionRequest)(nil), "proto.VersionRequest")
	proto.RegisterType((*VersionResponse)(nil), "proto.VersionResponse")
}

func init() { proto.RegisterFile("version.proto", fileDescriptor_7d2c07d79758f814) }

var fileDescriptor_7d2c07d79758f814 = []byte{
	// 153 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4b, 0x2d, 0x2a,
	0xce, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0x21, 0x5c,
	0x7c, 0x61, 0x10, 0xf1, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x21, 0x2e, 0x96, 0xb4,
	0xd2, 0x9c, 0x1c, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x8e, 0x20, 0x30, 0x5b, 0x48, 0x82, 0x8b, 0xbd,
	0x24, 0x33, 0x37, 0x35, 0xbf, 0xb4, 0x44, 0x82, 0x49, 0x81, 0x51, 0x83, 0x39, 0x08, 0xc6, 0x15,
	0x12, 0xe1, 0x62, 0x4d, 0xce, 0xcf, 0xc9, 0x2f, 0x92, 0x60, 0x06, 0x2b, 0x87, 0x70, 0x94, 0xdc,
	0xb9, 0xf8, 0xe1, 0xa6, 0x16, 0x17, 0xe4, 0xe7, 0x15, 0xa7, 0x0a, 0x49, 0x71, 0x71, 0x64, 0xe4,
	0x17, 0x97, 0xe4, 0x25, 0xe6, 0xa6, 0x82, 0x8d, 0xe6, 0x0c, 0x82, 0xf3, 0x41, 0xc6, 0x43, 0x1d,
	0x07, 0x36, 0x9e, 0x33, 0x08, 0xc6, 0x4d, 0x62, 0x03, 0xbb, 0xd2, 0x18, 0x10, 0x00, 0x00, 0xff,
	0xff, 0xf7, 0x2d, 0xe6, 0x34, 0xbd, 0x00, 0x00, 0x00,
}
