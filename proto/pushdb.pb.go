// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pushdb.proto

/*
Package pushdb is a generated protocol buffer package.

It is generated from these files:
	pushdb.proto

It has these top-level messages:
	Value
	Key
	SetRequest
	SetResponse
	GetRequest
	GetResponse
	WatchRequest
	WatchResponse
*/
package pushdb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/golang/protobuf/ptypes/empty"

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

type Value_Type int32

const (
	Value_DOUBLE Value_Type = 0
	Value_INT32  Value_Type = 1
	Value_BOOL   Value_Type = 2
	Value_STRING Value_Type = 3
)

var Value_Type_name = map[int32]string{
	0: "DOUBLE",
	1: "INT32",
	2: "BOOL",
	3: "STRING",
}
var Value_Type_value = map[string]int32{
	"DOUBLE": 0,
	"INT32":  1,
	"BOOL":   2,
	"STRING": 3,
}

func (x Value_Type) String() string {
	return proto.EnumName(Value_Type_name, int32(x))
}
func (Value_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Value struct {
	Type Value_Type `protobuf:"varint,2,opt,name=type,enum=pushdb.Value_Type" json:"type,omitempty"`
	// Types that are valid to be assigned to ValueOneof:
	//	*Value_DoubleValue
	//	*Value_Int32Value
	//	*Value_BoolValue
	//	*Value_StringValue
	ValueOneof isValue_ValueOneof `protobuf_oneof:"value_oneof"`
}

func (m *Value) Reset()                    { *m = Value{} }
func (m *Value) String() string            { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()               {}
func (*Value) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type isValue_ValueOneof interface {
	isValue_ValueOneof()
}

type Value_DoubleValue struct {
	DoubleValue float64 `protobuf:"fixed64,3,opt,name=double_value,json=doubleValue,oneof"`
}
type Value_Int32Value struct {
	Int32Value int32 `protobuf:"varint,4,opt,name=int32_value,json=int32Value,oneof"`
}
type Value_BoolValue struct {
	BoolValue bool `protobuf:"varint,5,opt,name=bool_value,json=boolValue,oneof"`
}
type Value_StringValue struct {
	StringValue string `protobuf:"bytes,6,opt,name=string_value,json=stringValue,oneof"`
}

func (*Value_DoubleValue) isValue_ValueOneof() {}
func (*Value_Int32Value) isValue_ValueOneof()  {}
func (*Value_BoolValue) isValue_ValueOneof()   {}
func (*Value_StringValue) isValue_ValueOneof() {}

func (m *Value) GetValueOneof() isValue_ValueOneof {
	if m != nil {
		return m.ValueOneof
	}
	return nil
}

func (m *Value) GetType() Value_Type {
	if m != nil {
		return m.Type
	}
	return Value_DOUBLE
}

func (m *Value) GetDoubleValue() float64 {
	if x, ok := m.GetValueOneof().(*Value_DoubleValue); ok {
		return x.DoubleValue
	}
	return 0
}

func (m *Value) GetInt32Value() int32 {
	if x, ok := m.GetValueOneof().(*Value_Int32Value); ok {
		return x.Int32Value
	}
	return 0
}

func (m *Value) GetBoolValue() bool {
	if x, ok := m.GetValueOneof().(*Value_BoolValue); ok {
		return x.BoolValue
	}
	return false
}

func (m *Value) GetStringValue() string {
	if x, ok := m.GetValueOneof().(*Value_StringValue); ok {
		return x.StringValue
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Value) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Value_OneofMarshaler, _Value_OneofUnmarshaler, _Value_OneofSizer, []interface{}{
		(*Value_DoubleValue)(nil),
		(*Value_Int32Value)(nil),
		(*Value_BoolValue)(nil),
		(*Value_StringValue)(nil),
	}
}

func _Value_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Value)
	// value_oneof
	switch x := m.ValueOneof.(type) {
	case *Value_DoubleValue:
		b.EncodeVarint(3<<3 | proto.WireFixed64)
		b.EncodeFixed64(math.Float64bits(x.DoubleValue))
	case *Value_Int32Value:
		b.EncodeVarint(4<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.Int32Value))
	case *Value_BoolValue:
		t := uint64(0)
		if x.BoolValue {
			t = 1
		}
		b.EncodeVarint(5<<3 | proto.WireVarint)
		b.EncodeVarint(t)
	case *Value_StringValue:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.StringValue)
	case nil:
	default:
		return fmt.Errorf("Value.ValueOneof has unexpected type %T", x)
	}
	return nil
}

func _Value_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Value)
	switch tag {
	case 3: // value_oneof.double_value
		if wire != proto.WireFixed64 {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeFixed64()
		m.ValueOneof = &Value_DoubleValue{math.Float64frombits(x)}
		return true, err
	case 4: // value_oneof.int32_value
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.ValueOneof = &Value_Int32Value{int32(x)}
		return true, err
	case 5: // value_oneof.bool_value
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.ValueOneof = &Value_BoolValue{x != 0}
		return true, err
	case 6: // value_oneof.string_value
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.ValueOneof = &Value_StringValue{x}
		return true, err
	default:
		return false, nil
	}
}

func _Value_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Value)
	// value_oneof
	switch x := m.ValueOneof.(type) {
	case *Value_DoubleValue:
		n += proto.SizeVarint(3<<3 | proto.WireFixed64)
		n += 8
	case *Value_Int32Value:
		n += proto.SizeVarint(4<<3 | proto.WireVarint)
		n += proto.SizeVarint(uint64(x.Int32Value))
	case *Value_BoolValue:
		n += proto.SizeVarint(5<<3 | proto.WireVarint)
		n += 1
	case *Value_StringValue:
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.StringValue)))
		n += len(x.StringValue)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Key struct {
	Name    string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Value   *Value `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	Version int32  `protobuf:"varint,3,opt,name=version" json:"version,omitempty"`
}

func (m *Key) Reset()                    { *m = Key{} }
func (m *Key) String() string            { return proto.CompactTextString(m) }
func (*Key) ProtoMessage()               {}
func (*Key) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Key) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Key) GetValue() *Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Key) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

type SetRequest struct {
	Key *Key `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *SetRequest) Reset()                    { *m = SetRequest{} }
func (m *SetRequest) String() string            { return proto.CompactTextString(m) }
func (*SetRequest) ProtoMessage()               {}
func (*SetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SetRequest) GetKey() *Key {
	if m != nil {
		return m.Key
	}
	return nil
}

type SetResponse struct {
	Key *Key `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *SetResponse) Reset()                    { *m = SetResponse{} }
func (m *SetResponse) String() string            { return proto.CompactTextString(m) }
func (*SetResponse) ProtoMessage()               {}
func (*SetResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *SetResponse) GetKey() *Key {
	if m != nil {
		return m.Key
	}
	return nil
}

type GetRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GetRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GetResponse struct {
	Key *Key `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *GetResponse) Reset()                    { *m = GetResponse{} }
func (m *GetResponse) String() string            { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()               {}
func (*GetResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GetResponse) GetKey() *Key {
	if m != nil {
		return m.Key
	}
	return nil
}

type WatchRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *WatchRequest) Reset()                    { *m = WatchRequest{} }
func (m *WatchRequest) String() string            { return proto.CompactTextString(m) }
func (*WatchRequest) ProtoMessage()               {}
func (*WatchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *WatchRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type WatchResponse struct {
	Key *Key `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *WatchResponse) Reset()                    { *m = WatchResponse{} }
func (m *WatchResponse) String() string            { return proto.CompactTextString(m) }
func (*WatchResponse) ProtoMessage()               {}
func (*WatchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *WatchResponse) GetKey() *Key {
	if m != nil {
		return m.Key
	}
	return nil
}

func init() {
	proto.RegisterType((*Value)(nil), "pushdb.Value")
	proto.RegisterType((*Key)(nil), "pushdb.Key")
	proto.RegisterType((*SetRequest)(nil), "pushdb.SetRequest")
	proto.RegisterType((*SetResponse)(nil), "pushdb.SetResponse")
	proto.RegisterType((*GetRequest)(nil), "pushdb.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "pushdb.GetResponse")
	proto.RegisterType((*WatchRequest)(nil), "pushdb.WatchRequest")
	proto.RegisterType((*WatchResponse)(nil), "pushdb.WatchResponse")
	proto.RegisterEnum("pushdb.Value_Type", Value_Type_name, Value_Type_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for PushdbService service

type PushdbServiceClient interface {
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (PushdbService_WatchClient, error)
}

type pushdbServiceClient struct {
	cc *grpc.ClientConn
}

func NewPushdbServiceClient(cc *grpc.ClientConn) PushdbServiceClient {
	return &pushdbServiceClient{cc}
}

func (c *pushdbServiceClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := grpc.Invoke(ctx, "/pushdb.PushdbService/Set", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pushdbServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := grpc.Invoke(ctx, "/pushdb.PushdbService/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pushdbServiceClient) Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (PushdbService_WatchClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_PushdbService_serviceDesc.Streams[0], c.cc, "/pushdb.PushdbService/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &pushdbServiceWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PushdbService_WatchClient interface {
	Recv() (*WatchResponse, error)
	grpc.ClientStream
}

type pushdbServiceWatchClient struct {
	grpc.ClientStream
}

func (x *pushdbServiceWatchClient) Recv() (*WatchResponse, error) {
	m := new(WatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for PushdbService service

type PushdbServiceServer interface {
	Set(context.Context, *SetRequest) (*SetResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Watch(*WatchRequest, PushdbService_WatchServer) error
}

func RegisterPushdbServiceServer(s *grpc.Server, srv PushdbServiceServer) {
	s.RegisterService(&_PushdbService_serviceDesc, srv)
}

func _PushdbService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PushdbServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pushdb.PushdbService/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PushdbServiceServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PushdbService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PushdbServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pushdb.PushdbService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PushdbServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PushdbService_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PushdbServiceServer).Watch(m, &pushdbServiceWatchServer{stream})
}

type PushdbService_WatchServer interface {
	Send(*WatchResponse) error
	grpc.ServerStream
}

type pushdbServiceWatchServer struct {
	grpc.ServerStream
}

func (x *pushdbServiceWatchServer) Send(m *WatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _PushdbService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pushdb.PushdbService",
	HandlerType: (*PushdbServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _PushdbService_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _PushdbService_Get_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _PushdbService_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pushdb.proto",
}

func init() { proto.RegisterFile("pushdb.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 437 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x93, 0xdf, 0x6b, 0x9b, 0x50,
	0x14, 0xc7, 0xbd, 0x51, 0xb3, 0xe6, 0x98, 0x8c, 0x70, 0xb7, 0x41, 0xc8, 0x18, 0x73, 0x06, 0x86,
	0xb0, 0x91, 0x16, 0x7d, 0xd9, 0x73, 0xd8, 0x30, 0xa5, 0xa5, 0x19, 0xd7, 0xec, 0xc7, 0x5b, 0x89,
	0xed, 0x69, 0x2a, 0xb3, 0x5e, 0xa7, 0xd7, 0x80, 0x8f, 0xfb, 0x83, 0xf6, 0x3f, 0x8e, 0x7b, 0xaf,
	0xae, 0xcd, 0x28, 0x6b, 0xdf, 0xbc, 0xdf, 0xf3, 0x39, 0x9f, 0x73, 0x3d, 0x28, 0x0c, 0x8b, 0xba,
	0xba, 0xbe, 0x4c, 0xe6, 0x45, 0xc9, 0x05, 0xa7, 0x7d, 0x7d, 0x9a, 0xbe, 0xdc, 0x72, 0xbe, 0xcd,
	0xf0, 0x50, 0xa5, 0x49, 0x7d, 0x75, 0x88, 0x37, 0x85, 0x68, 0x34, 0xe4, 0xfd, 0xea, 0x81, 0xfd,
	0x75, 0x93, 0xd5, 0x48, 0xdf, 0x82, 0x25, 0x9a, 0x02, 0x27, 0x3d, 0x97, 0xf8, 0x4f, 0x03, 0x3a,
	0x6f, 0x5d, 0xaa, 0x38, 0x5f, 0x37, 0x05, 0x32, 0x55, 0xa7, 0x33, 0x18, 0x5e, 0xf2, 0x3a, 0xc9,
	0xf0, 0x7c, 0x27, 0x4b, 0x13, 0xd3, 0x25, 0x3e, 0x59, 0x1a, 0xcc, 0xd1, 0xa9, 0x96, 0xbd, 0x01,
	0x27, 0xcd, 0x45, 0x18, 0xb4, 0x8c, 0xe5, 0x12, 0xdf, 0x5e, 0x1a, 0x0c, 0x54, 0xa8, 0x91, 0xd7,
	0x00, 0x09, 0xe7, 0x59, 0x4b, 0xd8, 0x2e, 0xf1, 0x0f, 0x96, 0x06, 0x1b, 0xc8, 0x4c, 0x03, 0x33,
	0x18, 0x56, 0xa2, 0x4c, 0xf3, 0x6d, 0x8b, 0xf4, 0x5d, 0xe2, 0x0f, 0xe4, 0x20, 0x9d, 0x2a, 0xc8,
	0x0b, 0xc1, 0x92, 0x77, 0xa3, 0x00, 0xfd, 0x8f, 0xab, 0x2f, 0x8b, 0xd3, 0x4f, 0x63, 0x83, 0x0e,
	0xc0, 0x3e, 0x3e, 0x5b, 0x87, 0xc1, 0x98, 0xd0, 0x03, 0xb0, 0x16, 0xab, 0xd5, 0xe9, 0xb8, 0x27,
	0x81, 0x78, 0xcd, 0x8e, 0xcf, 0xa2, 0xb1, 0xb9, 0x18, 0x81, 0xa3, 0x94, 0xe7, 0x3c, 0x47, 0x7e,
	0xe5, 0x7d, 0x07, 0xf3, 0x04, 0x1b, 0x4a, 0xc1, 0xca, 0x37, 0x37, 0x38, 0x21, 0x72, 0x0e, 0x53,
	0xcf, 0x74, 0x06, 0xb6, 0x1e, 0x2e, 0xb7, 0xe2, 0x04, 0xa3, 0xbd, 0xad, 0x30, 0x5d, 0xa3, 0x13,
	0x78, 0xb2, 0xc3, 0xb2, 0x4a, 0x79, 0xae, 0x96, 0x61, 0xb3, 0xee, 0xe8, 0xbd, 0x03, 0x88, 0x51,
	0x30, 0xfc, 0x59, 0x63, 0x25, 0xe8, 0x2b, 0x30, 0x7f, 0x60, 0xa3, 0xfc, 0x4e, 0xe0, 0x74, 0xaa,
	0x13, 0x6c, 0x98, 0xcc, 0xbd, 0xf7, 0xe0, 0x28, 0xb8, 0x2a, 0x78, 0x5e, 0xe1, 0x43, 0xb4, 0x0b,
	0x10, 0xdd, 0xaa, 0xef, 0xb9, 0xbb, 0xf4, 0x45, 0x8f, 0xf7, 0x79, 0x30, 0xfc, 0xb6, 0x11, 0x17,
	0xd7, 0xff, 0x33, 0xce, 0x61, 0xd4, 0x32, 0x8f, 0x72, 0x06, 0xbf, 0x09, 0x8c, 0x3e, 0xab, 0x2c,
	0xc6, 0x72, 0x97, 0x5e, 0x20, 0x3d, 0x02, 0x33, 0x46, 0x41, 0xff, 0x7e, 0x5d, 0xb7, 0xdb, 0x99,
	0x3e, 0xdb, 0xcb, 0xf4, 0x00, 0xcf, 0x90, 0x1d, 0xd1, 0xdd, 0x8e, 0xe8, 0x9e, 0x8e, 0x68, 0xaf,
	0xe3, 0x03, 0xd8, 0xea, 0x96, 0xf4, 0x79, 0x57, 0xbf, 0xfb, 0x62, 0xd3, 0x17, 0xff, 0xa4, 0x5d,
	0xdf, 0x11, 0x49, 0xfa, 0xea, 0x9f, 0x08, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x8a, 0x95, 0x74,
	0x0f, 0x48, 0x03, 0x00, 0x00,
}
