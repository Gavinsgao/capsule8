// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config.proto

package v0

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This message is just a wrapper around the actual config contained w/in the payload
type Config struct {
	PackageName string `protobuf:"bytes,1,opt,name=package_name,json=packageName" json:"package_name,omitempty"`
	Payload     []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	Descriptor_ []byte `protobuf:"bytes,3,opt,name=descriptor,proto3" json:"descriptor,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *Config) GetPackageName() string {
	if m != nil {
		return m.PackageName
	}
	return ""
}

func (m *Config) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Config) GetDescriptor_() []byte {
	if m != nil {
		return m.Descriptor_
	}
	return nil
}

func init() {
	proto.RegisterType((*Config)(nil), "capsule8.v0.Config")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0x4b,
	0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4e, 0x4e, 0x2c, 0x28, 0x2e, 0xcd,
	0x49, 0xb5, 0xd0, 0x2b, 0x33, 0x50, 0x4a, 0xe5, 0x62, 0x73, 0x06, 0x4b, 0x0a, 0x29, 0x72, 0xf1,
	0x14, 0x24, 0x26, 0x67, 0x27, 0xa6, 0xa7, 0xc6, 0xe7, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30,
	0x6a, 0x70, 0x06, 0x71, 0x43, 0xc5, 0xfc, 0x12, 0x73, 0x53, 0x85, 0x24, 0xb8, 0xd8, 0x0b, 0x12,
	0x2b, 0x73, 0xf2, 0x13, 0x53, 0x24, 0x98, 0x14, 0x18, 0x35, 0x78, 0x82, 0x60, 0x5c, 0x21, 0x39,
	0x2e, 0xae, 0x94, 0xd4, 0xe2, 0xe4, 0xa2, 0xcc, 0x82, 0x92, 0xfc, 0x22, 0x09, 0x66, 0xb0, 0x24,
	0x92, 0x88, 0x13, 0x4b, 0x14, 0x53, 0x99, 0x41, 0x12, 0x1b, 0xd8, 0x01, 0xc6, 0x80, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xff, 0x0d, 0x15, 0x7b, 0x90, 0x00, 0x00, 0x00,
}