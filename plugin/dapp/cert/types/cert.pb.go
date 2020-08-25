// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cert.proto

package types

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

type Cert struct {
	CertId               []byte   `protobuf:"bytes,1,opt,name=certId,proto3" json:"certId,omitempty"`
	CreateTime           int64    `protobuf:"varint,2,opt,name=createTime,proto3" json:"createTime,omitempty"`
	Key                  string   `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Cert) Reset()         { *m = Cert{} }
func (m *Cert) String() string { return proto.CompactTextString(m) }
func (*Cert) ProtoMessage()    {}
func (*Cert) Descriptor() ([]byte, []int) {
	return fileDescriptor_a142e29cbef9b1cf, []int{0}
}

func (m *Cert) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Cert.Unmarshal(m, b)
}
func (m *Cert) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Cert.Marshal(b, m, deterministic)
}
func (m *Cert) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Cert.Merge(m, src)
}
func (m *Cert) XXX_Size() int {
	return xxx_messageInfo_Cert.Size(m)
}
func (m *Cert) XXX_DiscardUnknown() {
	xxx_messageInfo_Cert.DiscardUnknown(m)
}

var xxx_messageInfo_Cert proto.InternalMessageInfo

func (m *Cert) GetCertId() []byte {
	if m != nil {
		return m.CertId
	}
	return nil
}

func (m *Cert) GetCreateTime() int64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

func (m *Cert) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Cert) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type CertAction struct {
	// Types that are valid to be assigned to Value:
	//	*CertAction_New
	//	*CertAction_Update
	//	*CertAction_Normal
	Value                isCertAction_Value `protobuf_oneof:"value"`
	Ty                   int32              `protobuf:"varint,4,opt,name=ty,proto3" json:"ty,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *CertAction) Reset()         { *m = CertAction{} }
func (m *CertAction) String() string { return proto.CompactTextString(m) }
func (*CertAction) ProtoMessage()    {}
func (*CertAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_a142e29cbef9b1cf, []int{1}
}

func (m *CertAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CertAction.Unmarshal(m, b)
}
func (m *CertAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CertAction.Marshal(b, m, deterministic)
}
func (m *CertAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CertAction.Merge(m, src)
}
func (m *CertAction) XXX_Size() int {
	return xxx_messageInfo_CertAction.Size(m)
}
func (m *CertAction) XXX_DiscardUnknown() {
	xxx_messageInfo_CertAction.DiscardUnknown(m)
}

var xxx_messageInfo_CertAction proto.InternalMessageInfo

type isCertAction_Value interface {
	isCertAction_Value()
}

type CertAction_New struct {
	New *CertNew `protobuf:"bytes,1,opt,name=new,proto3,oneof"`
}

type CertAction_Update struct {
	Update *CertUpdate `protobuf:"bytes,2,opt,name=update,proto3,oneof"`
}

type CertAction_Normal struct {
	Normal *CertNormal `protobuf:"bytes,3,opt,name=normal,proto3,oneof"`
}

func (*CertAction_New) isCertAction_Value() {}

func (*CertAction_Update) isCertAction_Value() {}

func (*CertAction_Normal) isCertAction_Value() {}

func (m *CertAction) GetValue() isCertAction_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *CertAction) GetNew() *CertNew {
	if x, ok := m.GetValue().(*CertAction_New); ok {
		return x.New
	}
	return nil
}

func (m *CertAction) GetUpdate() *CertUpdate {
	if x, ok := m.GetValue().(*CertAction_Update); ok {
		return x.Update
	}
	return nil
}

func (m *CertAction) GetNormal() *CertNormal {
	if x, ok := m.GetValue().(*CertAction_Normal); ok {
		return x.Normal
	}
	return nil
}

func (m *CertAction) GetTy() int32 {
	if m != nil {
		return m.Ty
	}
	return 0
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*CertAction) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*CertAction_New)(nil),
		(*CertAction_Update)(nil),
		(*CertAction_Normal)(nil),
	}
}

type CertNew struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CertNew) Reset()         { *m = CertNew{} }
func (m *CertNew) String() string { return proto.CompactTextString(m) }
func (*CertNew) ProtoMessage()    {}
func (*CertNew) Descriptor() ([]byte, []int) {
	return fileDescriptor_a142e29cbef9b1cf, []int{2}
}

func (m *CertNew) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CertNew.Unmarshal(m, b)
}
func (m *CertNew) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CertNew.Marshal(b, m, deterministic)
}
func (m *CertNew) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CertNew.Merge(m, src)
}
func (m *CertNew) XXX_Size() int {
	return xxx_messageInfo_CertNew.Size(m)
}
func (m *CertNew) XXX_DiscardUnknown() {
	xxx_messageInfo_CertNew.DiscardUnknown(m)
}

var xxx_messageInfo_CertNew proto.InternalMessageInfo

func (m *CertNew) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *CertNew) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type CertUpdate struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CertUpdate) Reset()         { *m = CertUpdate{} }
func (m *CertUpdate) String() string { return proto.CompactTextString(m) }
func (*CertUpdate) ProtoMessage()    {}
func (*CertUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_a142e29cbef9b1cf, []int{3}
}

func (m *CertUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CertUpdate.Unmarshal(m, b)
}
func (m *CertUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CertUpdate.Marshal(b, m, deterministic)
}
func (m *CertUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CertUpdate.Merge(m, src)
}
func (m *CertUpdate) XXX_Size() int {
	return xxx_messageInfo_CertUpdate.Size(m)
}
func (m *CertUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_CertUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_CertUpdate proto.InternalMessageInfo

func (m *CertUpdate) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *CertUpdate) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type CertNormal struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CertNormal) Reset()         { *m = CertNormal{} }
func (m *CertNormal) String() string { return proto.CompactTextString(m) }
func (*CertNormal) ProtoMessage()    {}
func (*CertNormal) Descriptor() ([]byte, []int) {
	return fileDescriptor_a142e29cbef9b1cf, []int{4}
}

func (m *CertNormal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CertNormal.Unmarshal(m, b)
}
func (m *CertNormal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CertNormal.Marshal(b, m, deterministic)
}
func (m *CertNormal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CertNormal.Merge(m, src)
}
func (m *CertNormal) XXX_Size() int {
	return xxx_messageInfo_CertNormal.Size(m)
}
func (m *CertNormal) XXX_DiscardUnknown() {
	xxx_messageInfo_CertNormal.DiscardUnknown(m)
}

var xxx_messageInfo_CertNormal proto.InternalMessageInfo

func (m *CertNormal) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *CertNormal) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type Authority struct {
	Enable               bool     `protobuf:"varint,1,opt,name=enable,proto3" json:"enable,omitempty"`
	CryptoPath           string   `protobuf:"bytes,2,opt,name=cryptoPath,proto3" json:"cryptoPath,omitempty"`
	SignType             string   `protobuf:"bytes,3,opt,name=signType,proto3" json:"signType,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Authority) Reset()         { *m = Authority{} }
func (m *Authority) String() string { return proto.CompactTextString(m) }
func (*Authority) ProtoMessage()    {}
func (*Authority) Descriptor() ([]byte, []int) {
	return fileDescriptor_a142e29cbef9b1cf, []int{5}
}

func (m *Authority) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Authority.Unmarshal(m, b)
}
func (m *Authority) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Authority.Marshal(b, m, deterministic)
}
func (m *Authority) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Authority.Merge(m, src)
}
func (m *Authority) XXX_Size() int {
	return xxx_messageInfo_Authority.Size(m)
}
func (m *Authority) XXX_DiscardUnknown() {
	xxx_messageInfo_Authority.DiscardUnknown(m)
}

var xxx_messageInfo_Authority proto.InternalMessageInfo

func (m *Authority) GetEnable() bool {
	if m != nil {
		return m.Enable
	}
	return false
}

func (m *Authority) GetCryptoPath() string {
	if m != nil {
		return m.CryptoPath
	}
	return ""
}

func (m *Authority) GetSignType() string {
	if m != nil {
		return m.SignType
	}
	return ""
}

type CertSignature struct {
	Signature            []byte   `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	Cert                 []byte   `protobuf:"bytes,2,opt,name=cert,proto3" json:"cert,omitempty"`
	Uid                  []byte   `protobuf:"bytes,3,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CertSignature) Reset()         { *m = CertSignature{} }
func (m *CertSignature) String() string { return proto.CompactTextString(m) }
func (*CertSignature) ProtoMessage()    {}
func (*CertSignature) Descriptor() ([]byte, []int) {
	return fileDescriptor_a142e29cbef9b1cf, []int{6}
}

func (m *CertSignature) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CertSignature.Unmarshal(m, b)
}
func (m *CertSignature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CertSignature.Marshal(b, m, deterministic)
}
func (m *CertSignature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CertSignature.Merge(m, src)
}
func (m *CertSignature) XXX_Size() int {
	return xxx_messageInfo_CertSignature.Size(m)
}
func (m *CertSignature) XXX_DiscardUnknown() {
	xxx_messageInfo_CertSignature.DiscardUnknown(m)
}

var xxx_messageInfo_CertSignature proto.InternalMessageInfo

func (m *CertSignature) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *CertSignature) GetCert() []byte {
	if m != nil {
		return m.Cert
	}
	return nil
}

func (m *CertSignature) GetUid() []byte {
	if m != nil {
		return m.Uid
	}
	return nil
}

type ReqQueryValidCertSN struct {
	Addr                 string   `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqQueryValidCertSN) Reset()         { *m = ReqQueryValidCertSN{} }
func (m *ReqQueryValidCertSN) String() string { return proto.CompactTextString(m) }
func (*ReqQueryValidCertSN) ProtoMessage()    {}
func (*ReqQueryValidCertSN) Descriptor() ([]byte, []int) {
	return fileDescriptor_a142e29cbef9b1cf, []int{7}
}

func (m *ReqQueryValidCertSN) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqQueryValidCertSN.Unmarshal(m, b)
}
func (m *ReqQueryValidCertSN) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqQueryValidCertSN.Marshal(b, m, deterministic)
}
func (m *ReqQueryValidCertSN) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqQueryValidCertSN.Merge(m, src)
}
func (m *ReqQueryValidCertSN) XXX_Size() int {
	return xxx_messageInfo_ReqQueryValidCertSN.Size(m)
}
func (m *ReqQueryValidCertSN) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqQueryValidCertSN.DiscardUnknown(m)
}

var xxx_messageInfo_ReqQueryValidCertSN proto.InternalMessageInfo

func (m *ReqQueryValidCertSN) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type RepQueryValidCertSN struct {
	Sn                   []byte   `protobuf:"bytes,1,opt,name=sn,proto3" json:"sn,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RepQueryValidCertSN) Reset()         { *m = RepQueryValidCertSN{} }
func (m *RepQueryValidCertSN) String() string { return proto.CompactTextString(m) }
func (*RepQueryValidCertSN) ProtoMessage()    {}
func (*RepQueryValidCertSN) Descriptor() ([]byte, []int) {
	return fileDescriptor_a142e29cbef9b1cf, []int{8}
}

func (m *RepQueryValidCertSN) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RepQueryValidCertSN.Unmarshal(m, b)
}
func (m *RepQueryValidCertSN) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RepQueryValidCertSN.Marshal(b, m, deterministic)
}
func (m *RepQueryValidCertSN) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RepQueryValidCertSN.Merge(m, src)
}
func (m *RepQueryValidCertSN) XXX_Size() int {
	return xxx_messageInfo_RepQueryValidCertSN.Size(m)
}
func (m *RepQueryValidCertSN) XXX_DiscardUnknown() {
	xxx_messageInfo_RepQueryValidCertSN.DiscardUnknown(m)
}

var xxx_messageInfo_RepQueryValidCertSN proto.InternalMessageInfo

func (m *RepQueryValidCertSN) GetSn() []byte {
	if m != nil {
		return m.Sn
	}
	return nil
}

func init() {
	proto.RegisterType((*Cert)(nil), "types.Cert")
	proto.RegisterType((*CertAction)(nil), "types.CertAction")
	proto.RegisterType((*CertNew)(nil), "types.CertNew")
	proto.RegisterType((*CertUpdate)(nil), "types.CertUpdate")
	proto.RegisterType((*CertNormal)(nil), "types.CertNormal")
	proto.RegisterType((*Authority)(nil), "types.Authority")
	proto.RegisterType((*CertSignature)(nil), "types.CertSignature")
	proto.RegisterType((*ReqQueryValidCertSN)(nil), "types.ReqQueryValidCertSN")
	proto.RegisterType((*RepQueryValidCertSN)(nil), "types.RepQueryValidCertSN")
}

func init() { proto.RegisterFile("cert.proto", fileDescriptor_a142e29cbef9b1cf) }

var fileDescriptor_a142e29cbef9b1cf = []byte{
	// 380 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0x5d, 0x4b, 0xeb, 0x40,
	0x10, 0x6d, 0x92, 0x7e, 0x65, 0xda, 0x5b, 0xee, 0xdd, 0x2b, 0x12, 0x44, 0xa4, 0x04, 0x84, 0x8a,
	0x50, 0xb0, 0xfa, 0x07, 0xaa, 0x2f, 0xf5, 0xa5, 0xe8, 0xb6, 0xfa, 0x2a, 0xdb, 0x66, 0x6c, 0x83,
	0xe9, 0x26, 0x6e, 0x36, 0x96, 0xfd, 0x3d, 0xfe, 0x51, 0xd9, 0xcd, 0x56, 0x23, 0xfa, 0xa0, 0x6f,
	0x33, 0x3b, 0xe7, 0xec, 0x39, 0x87, 0x19, 0x80, 0x25, 0x0a, 0x39, 0xcc, 0x44, 0x2a, 0x53, 0xd2,
	0x90, 0x2a, 0xc3, 0x3c, 0x7c, 0x84, 0xfa, 0x15, 0x0a, 0x49, 0xf6, 0xa1, 0xa9, 0x87, 0xd7, 0x51,
	0xe0, 0xf4, 0x9d, 0x41, 0x97, 0xda, 0x8e, 0x1c, 0x01, 0x2c, 0x05, 0x32, 0x89, 0xf3, 0x78, 0x83,
	0x81, 0xdb, 0x77, 0x06, 0x1e, 0xad, 0xbc, 0x90, 0xbf, 0xe0, 0x3d, 0xa1, 0x0a, 0xbc, 0xbe, 0x33,
	0xf0, 0xa9, 0x2e, 0xc9, 0x1e, 0x34, 0x5e, 0x58, 0x52, 0x60, 0x50, 0x37, 0x1f, 0x95, 0x4d, 0xf8,
	0xea, 0x00, 0x68, 0xa1, 0xf1, 0x52, 0xc6, 0x29, 0x27, 0x21, 0x78, 0x1c, 0xb7, 0x46, 0xab, 0x33,
	0xea, 0x0d, 0x8d, 0x97, 0xa1, 0x9e, 0x4f, 0x71, 0x3b, 0xa9, 0x51, 0x3d, 0x24, 0xa7, 0xd0, 0x2c,
	0xb2, 0x88, 0xc9, 0x52, 0xb6, 0x33, 0xfa, 0x57, 0x81, 0xdd, 0x99, 0xc1, 0xa4, 0x46, 0x2d, 0x44,
	0x83, 0x79, 0x2a, 0x36, 0x2c, 0x31, 0x56, 0x3e, 0x83, 0xa7, 0x66, 0xa0, 0xc1, 0x25, 0x84, 0xf4,
	0xc0, 0x95, 0xca, 0xf8, 0x6b, 0x50, 0x57, 0xaa, 0xcb, 0x96, 0xb5, 0x1c, 0x9e, 0x41, 0xcb, 0x9a,
	0xd8, 0x05, 0x73, 0xbe, 0x09, 0xe6, 0x56, 0x83, 0x5d, 0x94, 0xb9, 0x4a, 0x43, 0xbf, 0x65, 0x95,
	0xce, 0x7e, 0xcc, 0x7a, 0x00, 0x7f, 0x5c, 0xc8, 0x75, 0x2a, 0x62, 0xa9, 0xf4, 0xc6, 0x90, 0xb3,
	0x45, 0x82, 0x86, 0xd7, 0xa6, 0xb6, 0x2b, 0x37, 0xa6, 0x32, 0x99, 0xde, 0x30, 0xb9, 0x36, 0x7c,
	0x9f, 0x56, 0x5e, 0xc8, 0x01, 0xb4, 0xf3, 0x78, 0xc5, 0xe7, 0x2a, 0x43, 0xbb, 0xb6, 0xf7, 0x3e,
	0x9c, 0xc1, 0x1f, 0x6d, 0x6b, 0x16, 0xaf, 0x38, 0x93, 0x85, 0x40, 0x72, 0x08, 0x7e, 0xbe, 0x6b,
	0xec, 0x65, 0x7c, 0x3c, 0x10, 0x02, 0x75, 0x7d, 0x26, 0xd6, 0xa4, 0xa9, 0x75, 0x96, 0x22, 0x8e,
	0xcc, 0xcf, 0x5d, 0xaa, 0xcb, 0xf0, 0x04, 0xfe, 0x53, 0x7c, 0xbe, 0x2d, 0x50, 0xa8, 0x7b, 0x96,
	0xc4, 0x91, 0x51, 0x98, 0x6a, 0x32, 0x8b, 0x22, 0x61, 0x53, 0x9b, 0x3a, 0x3c, 0xd6, 0xd0, 0xec,
	0x0b, 0xb4, 0x07, 0x6e, 0xce, 0xad, 0xbc, 0x9b, 0xf3, 0x45, 0xd3, 0x9c, 0xf0, 0xf9, 0x5b, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x7e, 0x07, 0x28, 0x2a, 0xd0, 0x02, 0x00, 0x00,
}
