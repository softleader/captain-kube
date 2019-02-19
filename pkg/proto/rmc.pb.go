// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rmc.proto

// 依照 https://cloud.google.com/apis/design/naming_convention 規範

package captainkube_v2

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

type RmcRequest struct {
	Chart                *Chart   `protobuf:"bytes,1,opt,name=chart,proto3" json:"chart,omitempty"`
	Retag                *ReTag   `protobuf:"bytes,2,opt,name=retag,proto3" json:"retag,omitempty"`
	Set                  []string `protobuf:"bytes,3,rep,name=set,proto3" json:"set,omitempty"`
	Constraint           string   `protobuf:"bytes,4,opt,name=constraint,proto3" json:"constraint,omitempty"`
	Verbose              bool     `protobuf:"varint,5,opt,name=verbose,proto3" json:"verbose,omitempty"`
	Timeout              int64    `protobuf:"varint,6,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Color                bool     `protobuf:"varint,7,opt,name=color,proto3" json:"color,omitempty"`
	Force                bool     `protobuf:"varint,8,opt,name=force,proto3" json:"force,omitempty"`
	DryRun               bool     `protobuf:"varint,9,opt,name=dryRun,proto3" json:"dryRun,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RmcRequest) Reset()         { *m = RmcRequest{} }
func (m *RmcRequest) String() string { return proto.CompactTextString(m) }
func (*RmcRequest) ProtoMessage()    {}
func (*RmcRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_35e41f379f3a71cc, []int{0}
}

func (m *RmcRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RmcRequest.Unmarshal(m, b)
}
func (m *RmcRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RmcRequest.Marshal(b, m, deterministic)
}
func (m *RmcRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RmcRequest.Merge(m, src)
}
func (m *RmcRequest) XXX_Size() int {
	return xxx_messageInfo_RmcRequest.Size(m)
}
func (m *RmcRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RmcRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RmcRequest proto.InternalMessageInfo

func (m *RmcRequest) GetChart() *Chart {
	if m != nil {
		return m.Chart
	}
	return nil
}

func (m *RmcRequest) GetRetag() *ReTag {
	if m != nil {
		return m.Retag
	}
	return nil
}

func (m *RmcRequest) GetSet() []string {
	if m != nil {
		return m.Set
	}
	return nil
}

func (m *RmcRequest) GetConstraint() string {
	if m != nil {
		return m.Constraint
	}
	return ""
}

func (m *RmcRequest) GetVerbose() bool {
	if m != nil {
		return m.Verbose
	}
	return false
}

func (m *RmcRequest) GetTimeout() int64 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *RmcRequest) GetColor() bool {
	if m != nil {
		return m.Color
	}
	return false
}

func (m *RmcRequest) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
}

func (m *RmcRequest) GetDryRun() bool {
	if m != nil {
		return m.DryRun
	}
	return false
}

func init() {
	proto.RegisterType((*RmcRequest)(nil), "softleader.captainkube.v2.RmcRequest")
}

func init() { proto.RegisterFile("rmc.proto", fileDescriptor_35e41f379f3a71cc) }

var fileDescriptor_35e41f379f3a71cc = []byte{
	// 260 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x49, 0x63, 0xd2, 0x66, 0x02, 0x22, 0x8b, 0xc8, 0xda, 0x83, 0x2c, 0x9e, 0x72, 0xda,
	0x43, 0x04, 0x1f, 0x40, 0xdf, 0x60, 0xf1, 0xe4, 0x6d, 0xb3, 0x9d, 0xc6, 0x60, 0x93, 0xad, 0x93,
	0x49, 0xc5, 0x17, 0xf2, 0x39, 0x25, 0xbb, 0x16, 0xf4, 0x50, 0x7a, 0xdb, 0xef, 0xdf, 0xef, 0x1f,
	0x98, 0x81, 0x82, 0x7a, 0xa7, 0xf7, 0xe4, 0xd9, 0x8b, 0xdb, 0xd1, 0x6f, 0x79, 0x87, 0x76, 0x83,
	0xa4, 0x9d, 0xdd, 0xb3, 0xed, 0x86, 0xf7, 0xa9, 0x41, 0x7d, 0xa8, 0xd7, 0xa5, 0x7b, 0xb3, 0xc4,
	0xd1, 0x5b, 0x97, 0x5d, 0x6f, 0x5b, 0x8c, 0x70, 0xff, 0xbd, 0x00, 0x30, 0xbd, 0x33, 0xf8, 0x31,
	0xe1, 0xc8, 0xe2, 0x11, 0xb2, 0xa0, 0xca, 0x44, 0x25, 0x55, 0x59, 0x2b, 0x7d, 0x72, 0xa6, 0x7e,
	0x9e, 0x3d, 0x13, 0xf5, 0xb9, 0x47, 0xc8, 0xb6, 0x95, 0x8b, 0xb3, 0x3d, 0x83, 0x2f, 0xb6, 0x35,
	0x51, 0x17, 0x57, 0x90, 0x8e, 0xc8, 0x32, 0x55, 0x69, 0x55, 0x98, 0xf9, 0x29, 0xee, 0x00, 0x9c,
	0x1f, 0x46, 0x26, 0xdb, 0x0d, 0x2c, 0x2f, 0x54, 0x52, 0x15, 0xe6, 0x4f, 0x22, 0x24, 0x2c, 0x0f,
	0x48, 0x8d, 0x1f, 0x51, 0x66, 0x2a, 0xa9, 0x56, 0xe6, 0x88, 0xf3, 0x0f, 0x77, 0x3d, 0xfa, 0x89,
	0x65, 0xae, 0x92, 0x2a, 0x35, 0x47, 0x14, 0xd7, 0x90, 0x39, 0xbf, 0xf3, 0x24, 0x97, 0xa1, 0x11,
	0x61, 0x4e, 0xb7, 0x9e, 0x1c, 0xca, 0x55, 0x4c, 0x03, 0x88, 0x1b, 0xc8, 0x37, 0xf4, 0x65, 0xa6,
	0x41, 0x16, 0x21, 0xfe, 0xa5, 0xa7, 0x1a, 0x14, 0x7f, 0x6a, 0xe7, 0xfb, 0xd3, 0xab, 0xbd, 0x5e,
	0xfe, 0xe7, 0x26, 0x0f, 0x37, 0x7e, 0xf8, 0x09, 0x00, 0x00, 0xff, 0xff, 0x3c, 0x78, 0xa8, 0xf2,
	0xa5, 0x01, 0x00, 0x00,
}