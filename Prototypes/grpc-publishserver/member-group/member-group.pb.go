// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chat-group.proto

package membergroup

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

type Member struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Ip                   string   `protobuf:"bytes,2,opt,name=ip" json:"ip,omitempty"`
	Port                 string   `protobuf:"bytes,3,opt,name=port" json:"port,omitempty"`
	Leader               bool     `protobuf:"varint,4,opt,name=leader" json:"leader,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Member) Reset()         { *m = Member{} }
func (m *Member) String() string { return proto.CompactTextString(m) }
func (*Member) ProtoMessage()    {}
func (*Member) Descriptor() ([]byte, []int) {
	return fileDescriptor_member_group_399f4c669f55fff0, []int{0}
}
func (m *Member) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Member.Unmarshal(m, b)
}
func (m *Member) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Member.Marshal(b, m, deterministic)
}
func (dst *Member) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Member.Merge(dst, src)
}
func (m *Member) XXX_Size() int {
	return xxx_messageInfo_Member.Size(m)
}
func (m *Member) XXX_DiscardUnknown() {
	xxx_messageInfo_Member.DiscardUnknown(m)
}

var xxx_messageInfo_Member proto.InternalMessageInfo

func (m *Member) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Member) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *Member) GetPort() string {
	if m != nil {
		return m.Port
	}
	return ""
}

func (m *Member) GetLeader() bool {
	if m != nil {
		return m.Leader
	}
	return false
}

type MemberList struct {
	// creates a slice of Member
	Member               []*Member `protobuf:"bytes,1,rep,name=member" json:"member,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *MemberList) Reset()         { *m = MemberList{} }
func (m *MemberList) String() string { return proto.CompactTextString(m) }
func (*MemberList) ProtoMessage()    {}
func (*MemberList) Descriptor() ([]byte, []int) {
	return fileDescriptor_member_group_399f4c669f55fff0, []int{1}
}
func (m *MemberList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MemberList.Unmarshal(m, b)
}
func (m *MemberList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MemberList.Marshal(b, m, deterministic)
}
func (dst *MemberList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MemberList.Merge(dst, src)
}
func (m *MemberList) XXX_Size() int {
	return xxx_messageInfo_MemberList.Size(m)
}
func (m *MemberList) XXX_DiscardUnknown() {
	xxx_messageInfo_MemberList.DiscardUnknown(m)
}

var xxx_messageInfo_MemberList proto.InternalMessageInfo

func (m *MemberList) GetMember() []*Member {
	if m != nil {
		return m.Member
	}
	return nil
}

// Empty message type used for List method
type Void struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Void) Reset()         { *m = Void{} }
func (m *Void) String() string { return proto.CompactTextString(m) }
func (*Void) ProtoMessage()    {}
func (*Void) Descriptor() ([]byte, []int) {
	return fileDescriptor_member_group_399f4c669f55fff0, []int{2}
}
func (m *Void) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Void.Unmarshal(m, b)
}
func (m *Void) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Void.Marshal(b, m, deterministic)
}
func (dst *Void) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Void.Merge(dst, src)
}
func (m *Void) XXX_Size() int {
	return xxx_messageInfo_Void.Size(m)
}
func (m *Void) XXX_DiscardUnknown() {
	xxx_messageInfo_Void.DiscardUnknown(m)
}

var xxx_messageInfo_Void proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Member)(nil), "membergroup.Member")
	proto.RegisterType((*MemberList)(nil), "membergroup.MemberList")
	proto.RegisterType((*Void)(nil), "membergroup.Void")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MembersClient is the client API for Members service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MembersClient interface {
	List(ctx context.Context, in *Void, opts ...grpc.CallOption) (*MemberList, error)
	Register(ctx context.Context, in *Member, opts ...grpc.CallOption) (*Member, error)
}

type membersClient struct {
	cc *grpc.ClientConn
}

func NewMembersClient(cc *grpc.ClientConn) MembersClient {
	return &membersClient{cc}
}

func (c *membersClient) List(ctx context.Context, in *Void, opts ...grpc.CallOption) (*MemberList, error) {
	out := new(MemberList)
	err := c.cc.Invoke(ctx, "/membergroup.Members/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *membersClient) Register(ctx context.Context, in *Member, opts ...grpc.CallOption) (*Member, error) {
	out := new(Member)
	err := c.cc.Invoke(ctx, "/membergroup.Members/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MembersServer is the server API for Members service.
type MembersServer interface {
	List(context.Context, *Void) (*MemberList, error)
	Register(context.Context, *Member) (*Member, error)
}

func RegisterMembersServer(s *grpc.Server, srv MembersServer) {
	s.RegisterService(&_Members_serviceDesc, srv)
}

func _Members_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MembersServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/membergroup.Members/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MembersServer).List(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _Members_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Member)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MembersServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/membergroup.Members/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MembersServer).Register(ctx, req.(*Member))
	}
	return interceptor(ctx, in, info, handler)
}

var _Members_serviceDesc = grpc.ServiceDesc{
	ServiceName: "membergroup.Members",
	HandlerType: (*MembersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _Members_List_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _Members_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat-group.proto",
}

func init() { proto.RegisterFile("chat-group.proto", fileDescriptor_member_group_399f4c669f55fff0) }

var fileDescriptor_member_group_399f4c669f55fff0 = []byte{
	// 208 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x6a, 0xc4, 0x20,
	0x10, 0x86, 0x63, 0x22, 0x36, 0x9d, 0x40, 0xa1, 0x53, 0x68, 0x25, 0xa7, 0xe0, 0x29, 0x50, 0x9a,
	0x43, 0x5a, 0x0a, 0x7d, 0x87, 0xf6, 0xe2, 0xa1, 0xec, 0x35, 0x21, 0x12, 0x84, 0xcd, 0x2a, 0xc6,
	0x65, 0x5f, 0x7f, 0x51, 0x73, 0xd8, 0x85, 0xdc, 0x66, 0x3e, 0x3f, 0xff, 0x5f, 0x04, 0x5c, 0xd4,
	0x32, 0x2a, 0xf7, 0x31, 0x3b, 0x73, 0xb6, 0x9d, 0x75, 0xc6, 0x1b, 0xac, 0x12, 0x8b, 0x48, 0x1c,
	0x80, 0xfd, 0xc5, 0x15, 0x11, 0xe8, 0x69, 0x58, 0x14, 0x27, 0x0d, 0x69, 0x1f, 0x65, 0x9c, 0xf1,
	0x09, 0x72, 0x6d, 0x79, 0x1e, 0x49, 0xae, 0x6d, 0x70, 0xac, 0x71, 0x9e, 0x17, 0xc9, 0x09, 0x33,
	0xbe, 0x02, 0x3b, 0xaa, 0x61, 0x52, 0x8e, 0xd3, 0x86, 0xb4, 0xa5, 0xdc, 0x36, 0xf1, 0x03, 0x90,
	0x92, 0x7f, 0xf5, 0xea, 0xf1, 0x1d, 0x58, 0xaa, 0xe5, 0xa4, 0x29, 0xda, 0xaa, 0x7f, 0xe9, 0x6e,
	0x5e, 0xd1, 0x25, 0x51, 0x6e, 0x8a, 0x60, 0x40, 0xff, 0x8d, 0x9e, 0xfa, 0x0b, 0x3c, 0xa4, 0x93,
	0x15, 0xbf, 0x80, 0xc6, 0x9c, 0xe7, 0xbb, 0x7b, 0xc1, 0xaa, 0xdf, 0x76, 0xa2, 0x82, 0x2b, 0x32,
	0xfc, 0x86, 0x52, 0xaa, 0x59, 0xaf, 0x5e, 0x39, 0xdc, 0x6b, 0xac, 0xf7, 0xa0, 0xc8, 0x46, 0x16,
	0x7f, 0xea, 0xf3, 0x1a, 0x00, 0x00, 0xff, 0xff, 0xde, 0xe0, 0xb2, 0x5f, 0x3f, 0x01, 0x00, 0x00,
}
