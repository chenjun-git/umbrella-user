// Code generated by protoc-gen-go. DO NOT EDIT.
// source: token.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	token.proto

It has these top-level messages:
	Contact
	Address
	VerifyTokenReq
	VerifyTokenResp
	RefreshTokenReq
	RefreshTokenResp
	CreatePersonalTokenReq
	CreatePersonalTokenResp
	DeletePersonalTokenReq
	DeletePersonalTokenResp
	Empty
*/
package pb

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

// Token check status
type VerifyResult int32

const (
	VerifyResult__          VerifyResult = 0
	VerifyResult_Allow      VerifyResult = 1
	VerifyResult_Deny       VerifyResult = 2
	VerifyResult_ExpireSoon VerifyResult = 3
	VerifyResult_Expired    VerifyResult = 4
)

var VerifyResult_name = map[int32]string{
	0: "_",
	1: "Allow",
	2: "Deny",
	3: "ExpireSoon",
	4: "Expired",
}
var VerifyResult_value = map[string]int32{
	"_":          0,
	"Allow":      1,
	"Deny":       2,
	"ExpireSoon": 3,
	"Expired":    4,
}

func (x VerifyResult) String() string {
	return proto.EnumName(VerifyResult_name, int32(x))
}
func (VerifyResult) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Contact struct {
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Phone string `protobuf:"bytes,2,opt,name=phone" json:"phone,omitempty"`
	Email string `protobuf:"bytes,3,opt,name=email" json:"email,omitempty"`
}

func (m *Contact) Reset()                    { *m = Contact{} }
func (m *Contact) String() string            { return proto.CompactTextString(m) }
func (*Contact) ProtoMessage()               {}
func (*Contact) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Contact) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Contact) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Contact) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

// 地址
type Address struct {
	Street      string `protobuf:"bytes,1,opt,name=street" json:"street,omitempty"`
	City        string `protobuf:"bytes,2,opt,name=city" json:"city,omitempty"`
	State       string `protobuf:"bytes,3,opt,name=state" json:"state,omitempty"`
	PostCode    string `protobuf:"bytes,4,opt,name=post_code" json:"post_code,omitempty"`
	Country     string `protobuf:"bytes,5,opt,name=country" json:"country,omitempty"`
	CountryCode string `protobuf:"bytes,6,opt,name=country_code" json:"country_code,omitempty"`
}

func (m *Address) Reset()                    { *m = Address{} }
func (m *Address) String() string            { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()               {}
func (*Address) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Address) GetStreet() string {
	if m != nil {
		return m.Street
	}
	return ""
}

func (m *Address) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *Address) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *Address) GetPostCode() string {
	if m != nil {
		return m.PostCode
	}
	return ""
}

func (m *Address) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *Address) GetCountryCode() string {
	if m != nil {
		return m.CountryCode
	}
	return ""
}

type VerifyTokenReq struct {
	AccessToken string `protobuf:"bytes,1,opt,name=access_token" json:"access_token,omitempty"`
}

func (m *VerifyTokenReq) Reset()                    { *m = VerifyTokenReq{} }
func (m *VerifyTokenReq) String() string            { return proto.CompactTextString(m) }
func (*VerifyTokenReq) ProtoMessage()               {}
func (*VerifyTokenReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *VerifyTokenReq) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

type VerifyTokenResp struct {
	Result  VerifyResult `protobuf:"varint,1,opt,name=Result,enum=pb.VerifyResult" json:"Result,omitempty"`
	UserId  string       `protobuf:"bytes,2,opt,name=user_id" json:"user_id,omitempty"`
	Device  string       `protobuf:"bytes,4,opt,name=device" json:"device,omitempty"`
	App     string       `protobuf:"bytes,5,opt,name=app" json:"app,omitempty"`
	Message string       `protobuf:"bytes,6,opt,name=message" json:"message,omitempty"`
}

func (m *VerifyTokenResp) Reset()                    { *m = VerifyTokenResp{} }
func (m *VerifyTokenResp) String() string            { return proto.CompactTextString(m) }
func (*VerifyTokenResp) ProtoMessage()               {}
func (*VerifyTokenResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *VerifyTokenResp) GetResult() VerifyResult {
	if m != nil {
		return m.Result
	}
	return VerifyResult__
}

func (m *VerifyTokenResp) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *VerifyTokenResp) GetDevice() string {
	if m != nil {
		return m.Device
	}
	return ""
}

func (m *VerifyTokenResp) GetApp() string {
	if m != nil {
		return m.App
	}
	return ""
}

func (m *VerifyTokenResp) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type RefreshTokenReq struct {
	RefreshToken string `protobuf:"bytes,1,opt,name=refresh_token" json:"refresh_token,omitempty"`
}

func (m *RefreshTokenReq) Reset()                    { *m = RefreshTokenReq{} }
func (m *RefreshTokenReq) String() string            { return proto.CompactTextString(m) }
func (*RefreshTokenReq) ProtoMessage()               {}
func (*RefreshTokenReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *RefreshTokenReq) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

type RefreshTokenResp struct {
	AccessToken  string `protobuf:"bytes,1,opt,name=access_token" json:"access_token,omitempty"`
	RefreshToken string `protobuf:"bytes,2,opt,name=refresh_token" json:"refresh_token,omitempty"`
}

func (m *RefreshTokenResp) Reset()                    { *m = RefreshTokenResp{} }
func (m *RefreshTokenResp) String() string            { return proto.CompactTextString(m) }
func (*RefreshTokenResp) ProtoMessage()               {}
func (*RefreshTokenResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *RefreshTokenResp) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *RefreshTokenResp) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

type CreatePersonalTokenReq struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *CreatePersonalTokenReq) Reset()                    { *m = CreatePersonalTokenReq{} }
func (m *CreatePersonalTokenReq) String() string            { return proto.CompactTextString(m) }
func (*CreatePersonalTokenReq) ProtoMessage()               {}
func (*CreatePersonalTokenReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CreatePersonalTokenReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CreatePersonalTokenResp struct {
	PersonalToken string `protobuf:"bytes,1,opt,name=personal_token" json:"personal_token,omitempty"`
}

func (m *CreatePersonalTokenResp) Reset()                    { *m = CreatePersonalTokenResp{} }
func (m *CreatePersonalTokenResp) String() string            { return proto.CompactTextString(m) }
func (*CreatePersonalTokenResp) ProtoMessage()               {}
func (*CreatePersonalTokenResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *CreatePersonalTokenResp) GetPersonalToken() string {
	if m != nil {
		return m.PersonalToken
	}
	return ""
}

type DeletePersonalTokenReq struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *DeletePersonalTokenReq) Reset()                    { *m = DeletePersonalTokenReq{} }
func (m *DeletePersonalTokenReq) String() string            { return proto.CompactTextString(m) }
func (*DeletePersonalTokenReq) ProtoMessage()               {}
func (*DeletePersonalTokenReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *DeletePersonalTokenReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type DeletePersonalTokenResp struct {
}

func (m *DeletePersonalTokenResp) Reset()                    { *m = DeletePersonalTokenResp{} }
func (m *DeletePersonalTokenResp) String() string            { return proto.CompactTextString(m) }
func (*DeletePersonalTokenResp) ProtoMessage()               {}
func (*DeletePersonalTokenResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func init() {
	proto.RegisterType((*Contact)(nil), "pb.Contact")
	proto.RegisterType((*Address)(nil), "pb.Address")
	proto.RegisterType((*VerifyTokenReq)(nil), "pb.VerifyTokenReq")
	proto.RegisterType((*VerifyTokenResp)(nil), "pb.VerifyTokenResp")
	proto.RegisterType((*RefreshTokenReq)(nil), "pb.RefreshTokenReq")
	proto.RegisterType((*RefreshTokenResp)(nil), "pb.RefreshTokenResp")
	proto.RegisterType((*CreatePersonalTokenReq)(nil), "pb.CreatePersonalTokenReq")
	proto.RegisterType((*CreatePersonalTokenResp)(nil), "pb.CreatePersonalTokenResp")
	proto.RegisterType((*DeletePersonalTokenReq)(nil), "pb.DeletePersonalTokenReq")
	proto.RegisterType((*DeletePersonalTokenResp)(nil), "pb.DeletePersonalTokenResp")
	proto.RegisterType((*Empty)(nil), "pb.Empty")
	proto.RegisterEnum("pb.VerifyResult", VerifyResult_name, VerifyResult_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Token service

type TokenClient interface {
	// 验证 token
	VerifyToken(ctx context.Context, in *VerifyTokenReq, opts ...grpc.CallOption) (*VerifyTokenResp, error)
	// 刷新token
	RefreshToken(ctx context.Context, in *RefreshTokenReq, opts ...grpc.CallOption) (*RefreshTokenResp, error)
	// 删除 所有 token
	// rpc RemoveAllToken (RemoveAllTokenReq) returns (RemoveAllTokenResp) {
	// }
	// 生成PersonalToken
	CreatePersonalToken(ctx context.Context, in *CreatePersonalTokenReq, opts ...grpc.CallOption) (*CreatePersonalTokenResp, error)
	// 删除PersonalToken
	DeletePersonalToken(ctx context.Context, in *DeletePersonalTokenReq, opts ...grpc.CallOption) (*DeletePersonalTokenResp, error)
}

type tokenClient struct {
	cc *grpc.ClientConn
}

func NewTokenClient(cc *grpc.ClientConn) TokenClient {
	return &tokenClient{cc}
}

func (c *tokenClient) VerifyToken(ctx context.Context, in *VerifyTokenReq, opts ...grpc.CallOption) (*VerifyTokenResp, error) {
	out := new(VerifyTokenResp)
	err := grpc.Invoke(ctx, "/pb.Token/VerifyToken", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenClient) RefreshToken(ctx context.Context, in *RefreshTokenReq, opts ...grpc.CallOption) (*RefreshTokenResp, error) {
	out := new(RefreshTokenResp)
	err := grpc.Invoke(ctx, "/pb.Token/RefreshToken", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenClient) CreatePersonalToken(ctx context.Context, in *CreatePersonalTokenReq, opts ...grpc.CallOption) (*CreatePersonalTokenResp, error) {
	out := new(CreatePersonalTokenResp)
	err := grpc.Invoke(ctx, "/pb.Token/CreatePersonalToken", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenClient) DeletePersonalToken(ctx context.Context, in *DeletePersonalTokenReq, opts ...grpc.CallOption) (*DeletePersonalTokenResp, error) {
	out := new(DeletePersonalTokenResp)
	err := grpc.Invoke(ctx, "/pb.Token/DeletePersonalToken", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Token service

type TokenServer interface {
	// 验证 token
	VerifyToken(context.Context, *VerifyTokenReq) (*VerifyTokenResp, error)
	// 刷新token
	RefreshToken(context.Context, *RefreshTokenReq) (*RefreshTokenResp, error)
	// 删除 所有 token
	// rpc RemoveAllToken (RemoveAllTokenReq) returns (RemoveAllTokenResp) {
	// }
	// 生成PersonalToken
	CreatePersonalToken(context.Context, *CreatePersonalTokenReq) (*CreatePersonalTokenResp, error)
	// 删除PersonalToken
	DeletePersonalToken(context.Context, *DeletePersonalTokenReq) (*DeletePersonalTokenResp, error)
}

func RegisterTokenServer(s *grpc.Server, srv TokenServer) {
	s.RegisterService(&_Token_serviceDesc, srv)
}

func _Token_VerifyToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServer).VerifyToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Token/VerifyToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServer).VerifyToken(ctx, req.(*VerifyTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Token_RefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServer).RefreshToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Token/RefreshToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServer).RefreshToken(ctx, req.(*RefreshTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Token_CreatePersonalToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePersonalTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServer).CreatePersonalToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Token/CreatePersonalToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServer).CreatePersonalToken(ctx, req.(*CreatePersonalTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Token_DeletePersonalToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePersonalTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServer).DeletePersonalToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Token/DeletePersonalToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServer).DeletePersonalToken(ctx, req.(*DeletePersonalTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Token_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Token",
	HandlerType: (*TokenServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VerifyToken",
			Handler:    _Token_VerifyToken_Handler,
		},
		{
			MethodName: "RefreshToken",
			Handler:    _Token_RefreshToken_Handler,
		},
		{
			MethodName: "CreatePersonalToken",
			Handler:    _Token_CreatePersonalToken_Handler,
		},
		{
			MethodName: "DeletePersonalToken",
			Handler:    _Token_DeletePersonalToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "token.proto",
}

func init() { proto.RegisterFile("token.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 457 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x4d, 0x1c, 0x7f, 0x90, 0x89, 0xeb, 0x98, 0x6d, 0x49, 0x4d, 0xb8, 0x54, 0x3e, 0x54, 0x15,
	0x87, 0x48, 0x14, 0x21, 0x21, 0x71, 0x40, 0x55, 0x5b, 0x71, 0xad, 0x02, 0xe2, 0x1a, 0x6d, 0xec,
	0x29, 0xb5, 0xb0, 0xbd, 0xcb, 0xce, 0x06, 0xf0, 0xdf, 0xe1, 0x97, 0x22, 0x7b, 0xed, 0xb6, 0x26,
	0xae, 0x38, 0xee, 0xf3, 0x9b, 0x37, 0x6f, 0xde, 0x93, 0x61, 0xa6, 0xc5, 0x77, 0x2c, 0x57, 0x52,
	0x09, 0x2d, 0x98, 0x25, 0xb7, 0xf1, 0x3b, 0xf0, 0x2e, 0x45, 0xa9, 0x79, 0xa2, 0x99, 0x0f, 0x76,
	0xc9, 0x0b, 0x8c, 0xc6, 0x27, 0xe3, 0xb3, 0x29, 0x3b, 0x00, 0x47, 0xde, 0x89, 0x12, 0x23, 0xab,
	0x7b, 0x62, 0xc1, 0xb3, 0x3c, 0x9a, 0xd4, 0xcf, 0x58, 0x82, 0x77, 0x91, 0xa6, 0x0a, 0x89, 0x58,
	0x00, 0x2e, 0x69, 0x85, 0xa8, 0xdb, 0x41, 0x1f, 0xec, 0x24, 0xd3, 0xd5, 0xc3, 0x1c, 0x69, 0xae,
	0xd1, 0xcc, 0xb1, 0xe7, 0x30, 0x95, 0x82, 0xf4, 0x26, 0x11, 0x29, 0x46, 0x76, 0x03, 0xcd, 0xc1,
	0x4b, 0xc4, 0xae, 0xd4, 0xaa, 0x8a, 0x9c, 0x06, 0x38, 0x02, 0xbf, 0x05, 0x0c, 0xcd, 0x6d, 0x36,
	0x9e, 0x42, 0xf0, 0x15, 0x55, 0x76, 0x5b, 0x7d, 0xa9, 0x2f, 0x58, 0xe3, 0x8f, 0x9a, 0xc7, 0x93,
	0x04, 0x89, 0x36, 0xcd, 0x51, 0x66, 0x7d, 0xac, 0x60, 0xde, 0xe3, 0x91, 0x64, 0x27, 0xe0, 0xae,
	0x91, 0x76, 0xb9, 0x71, 0x18, 0x9c, 0x87, 0x2b, 0xb9, 0x5d, 0x19, 0x92, 0xc1, 0x6b, 0x0f, 0x3b,
	0x42, 0xb5, 0xc9, 0xd2, 0xd6, 0x76, 0x00, 0x6e, 0x8a, 0x3f, 0xb3, 0xa4, 0x33, 0x39, 0x83, 0x09,
	0x97, 0xb2, 0x35, 0x38, 0x07, 0xaf, 0x40, 0x22, 0xfe, 0xad, 0xf3, 0x76, 0x06, 0xf3, 0x35, 0xde,
	0x2a, 0xa4, 0xbb, 0x7b, 0x73, 0x2f, 0xe0, 0x40, 0x19, 0xa8, 0xe7, 0xee, 0x23, 0x84, 0x7d, 0x26,
	0xc9, 0xe1, 0x3b, 0xf6, 0x05, 0xac, 0x36, 0x86, 0xc5, 0xa5, 0x42, 0xae, 0xf1, 0x06, 0x15, 0x89,
	0x92, 0xe7, 0xf7, 0x1b, 0x7b, 0xf5, 0xc5, 0x6f, 0xe0, 0x78, 0x90, 0x47, 0x92, 0x2d, 0x20, 0x90,
	0x2d, 0xd8, 0xf3, 0x76, 0x0a, 0x8b, 0x2b, 0xcc, 0xf1, 0xbf, 0xd2, 0x2f, 0xe1, 0x78, 0x90, 0x47,
	0x32, 0xf6, 0xc0, 0xb9, 0x2e, 0xa4, 0xae, 0x5e, 0x7f, 0x02, 0xbf, 0x17, 0xb0, 0x03, 0xe3, 0x4d,
	0x38, 0x62, 0x53, 0x70, 0x2e, 0xf2, 0x5c, 0xfc, 0x0a, 0xc7, 0xec, 0x19, 0xd8, 0x57, 0x58, 0x56,
	0xa1, 0xc5, 0x02, 0x80, 0xeb, 0xdf, 0x32, 0x53, 0xf8, 0x59, 0x88, 0x32, 0x9c, 0xb0, 0x19, 0x78,
	0xe6, 0x9d, 0x86, 0xf6, 0xf9, 0x1f, 0x0b, 0x9c, 0x46, 0x9f, 0xbd, 0x87, 0xd9, 0xa3, 0x62, 0x19,
	0x7b, 0x28, 0xb1, 0xf3, 0xb9, 0x3c, 0xdc, 0xc3, 0x48, 0xc6, 0x23, 0xf6, 0x01, 0xfc, 0xc7, 0xa1,
	0xb3, 0x86, 0xf6, 0x4f, 0x61, 0xcb, 0xa3, 0x7d, 0xb0, 0x19, 0xbe, 0x81, 0xc3, 0x81, 0x20, 0xd9,
	0xb2, 0xa6, 0x0f, 0x37, 0xb1, 0x7c, 0xf5, 0xe4, 0xb7, 0x4e, 0x71, 0x20, 0x3f, 0xa3, 0x38, 0x5c,
	0x80, 0x51, 0x7c, 0x2a, 0xf4, 0xd1, 0xd6, 0x6d, 0xfe, 0xe7, 0xb7, 0x7f, 0x03, 0x00, 0x00, 0xff,
	0xff, 0xe4, 0xb8, 0xe8, 0xd0, 0xde, 0x03, 0x00, 0x00,
}
