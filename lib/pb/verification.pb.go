// Code generated by protoc-gen-go. DO NOT EDIT.
// source: verification.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type VerificationRequest struct {
	AccessToken          string   `protobuf:"bytes,1,opt,name=accessToken,proto3" json:"accessToken,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VerificationRequest) Reset()         { *m = VerificationRequest{} }
func (m *VerificationRequest) String() string { return proto.CompactTextString(m) }
func (*VerificationRequest) ProtoMessage()    {}
func (*VerificationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_69b5d5d3b04d10d4, []int{0}
}

func (m *VerificationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerificationRequest.Unmarshal(m, b)
}
func (m *VerificationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerificationRequest.Marshal(b, m, deterministic)
}
func (m *VerificationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerificationRequest.Merge(m, src)
}
func (m *VerificationRequest) XXX_Size() int {
	return xxx_messageInfo_VerificationRequest.Size(m)
}
func (m *VerificationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VerificationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VerificationRequest proto.InternalMessageInfo

func (m *VerificationRequest) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

type VerificationResponse struct {
	UserEmail            string   `protobuf:"bytes,1,opt,name=userEmail,proto3" json:"userEmail,omitempty"`
	Valid                bool     `protobuf:"varint,2,opt,name=valid,proto3" json:"valid,omitempty"`
	Role                 int32    `protobuf:"varint,3,opt,name=role,proto3" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VerificationResponse) Reset()         { *m = VerificationResponse{} }
func (m *VerificationResponse) String() string { return proto.CompactTextString(m) }
func (*VerificationResponse) ProtoMessage()    {}
func (*VerificationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_69b5d5d3b04d10d4, []int{1}
}

func (m *VerificationResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerificationResponse.Unmarshal(m, b)
}
func (m *VerificationResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerificationResponse.Marshal(b, m, deterministic)
}
func (m *VerificationResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerificationResponse.Merge(m, src)
}
func (m *VerificationResponse) XXX_Size() int {
	return xxx_messageInfo_VerificationResponse.Size(m)
}
func (m *VerificationResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VerificationResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VerificationResponse proto.InternalMessageInfo

func (m *VerificationResponse) GetUserEmail() string {
	if m != nil {
		return m.UserEmail
	}
	return ""
}

func (m *VerificationResponse) GetValid() bool {
	if m != nil {
		return m.Valid
	}
	return false
}

func (m *VerificationResponse) GetRole() int32 {
	if m != nil {
		return m.Role
	}
	return 0
}

func init() {
	proto.RegisterType((*VerificationRequest)(nil), "pb.VerificationRequest")
	proto.RegisterType((*VerificationResponse)(nil), "pb.VerificationResponse")
}

func init() {
	proto.RegisterFile("verification.proto", fileDescriptor_69b5d5d3b04d10d4)
}

var fileDescriptor_69b5d5d3b04d10d4 = []byte{
	// 183 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x4b, 0x2d, 0xca,
	0x4c, 0xcb, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x2a, 0x48, 0x52, 0x32, 0xe7, 0x12, 0x0e, 0x43, 0x92, 0x09, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e,
	0x11, 0x52, 0xe0, 0xe2, 0x4e, 0x4c, 0x4e, 0x4e, 0x2d, 0x2e, 0x0e, 0xc9, 0xcf, 0x4e, 0xcd, 0x93,
	0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x42, 0x16, 0x52, 0x8a, 0xe3, 0x12, 0x41, 0xd5, 0x58, 0x5c,
	0x90, 0x9f, 0x57, 0x9c, 0x2a, 0x24, 0xc3, 0xc5, 0x59, 0x5a, 0x9c, 0x5a, 0xe4, 0x9a, 0x9b, 0x98,
	0x99, 0x03, 0xd5, 0x87, 0x10, 0x10, 0x12, 0xe1, 0x62, 0x2d, 0x4b, 0xcc, 0xc9, 0x4c, 0x91, 0x60,
	0x52, 0x60, 0xd4, 0xe0, 0x08, 0x82, 0x70, 0x84, 0x84, 0xb8, 0x58, 0x8a, 0xf2, 0x73, 0x52, 0x25,
	0x98, 0x15, 0x18, 0x35, 0x58, 0x83, 0xc0, 0x6c, 0x23, 0x6f, 0x2e, 0x1e, 0x64, 0xf3, 0x85, 0xac,
	0xb9, 0xd8, 0xc0, 0xfc, 0x4a, 0x21, 0x71, 0xbd, 0x82, 0x24, 0x3d, 0x2c, 0x8e, 0x96, 0x92, 0xc0,
	0x94, 0x80, 0x38, 0x2a, 0x89, 0x0d, 0xec, 0x61, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8d,
	0xc3, 0xe9, 0x47, 0x06, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// VerificationClient is the client API for Verification service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VerificationClient interface {
	Verify(ctx context.Context, in *VerificationRequest, opts ...grpc.CallOption) (*VerificationResponse, error)
}

type verificationClient struct {
	cc grpc.ClientConnInterface
}

func NewVerificationClient(cc grpc.ClientConnInterface) VerificationClient {
	return &verificationClient{cc}
}

func (c *verificationClient) Verify(ctx context.Context, in *VerificationRequest, opts ...grpc.CallOption) (*VerificationResponse, error) {
	out := new(VerificationResponse)
	err := c.cc.Invoke(ctx, "/pb.Verification/Verify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VerificationServer is the server API for Verification service.
type VerificationServer interface {
	Verify(context.Context, *VerificationRequest) (*VerificationResponse, error)
}

// UnimplementedVerificationServer can be embedded to have forward compatible implementations.
type UnimplementedVerificationServer struct {
}

func (*UnimplementedVerificationServer) Verify(ctx context.Context, req *VerificationRequest) (*VerificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Verify not implemented")
}

func RegisterVerificationServer(s *grpc.Server, srv VerificationServer) {
	s.RegisterService(&_Verification_serviceDesc, srv)
}

func _Verification_Verify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VerificationServer).Verify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Verification/Verify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VerificationServer).Verify(ctx, req.(*VerificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Verification_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Verification",
	HandlerType: (*VerificationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Verify",
			Handler:    _Verification_Verify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "verification.proto",
}