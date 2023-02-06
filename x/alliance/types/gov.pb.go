// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: alliance/gov.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MsgCreateAllianceProposal struct {
	// the title of the update proposal
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// the description of the proposal
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// Denom of the asset. It could either be a native token or an IBC token
	Denom string `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty" yaml:"denom"`
	// The reward weight specifies the ratio of rewards that will be given to each alliance asset
	// It does not need to sum to 1. rate = weight / total_weight
	// Native asset is always assumed to have a weight of 1.
	RewardWeight github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=reward_weight,json=rewardWeight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"reward_weight"`
	// set a bound of weight range to limit how much reward weights can scale.
	RewardWeightRange RewardWeightRange `protobuf:"bytes,5,opt,name=reward_weight_range,json=rewardWeightRange,proto3" json:"reward_weight_range"`
	// A positive take rate is used for liquid staking derivatives. It defines an annualized reward rate that
	// will be redirected to the distribution rewards pool
	TakeRate             github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=take_rate,json=takeRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"take_rate"`
	RewardChangeRate     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=reward_change_rate,json=rewardChangeRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"reward_change_rate"`
	RewardChangeInterval time.Duration                          `protobuf:"bytes,8,opt,name=reward_change_interval,json=rewardChangeInterval,proto3,stdduration" json:"reward_change_interval"`
}

func (m *MsgCreateAllianceProposal) Reset()         { *m = MsgCreateAllianceProposal{} }
func (m *MsgCreateAllianceProposal) String() string { return proto.CompactTextString(m) }
func (*MsgCreateAllianceProposal) ProtoMessage()    {}
func (*MsgCreateAllianceProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_5518a6f5c90c8452, []int{0}
}
func (m *MsgCreateAllianceProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreateAllianceProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateAllianceProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreateAllianceProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateAllianceProposal.Merge(m, src)
}
func (m *MsgCreateAllianceProposal) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreateAllianceProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateAllianceProposal.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateAllianceProposal proto.InternalMessageInfo

type MsgUpdateAllianceProposal struct {
	// the title of the update proposal
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// the description of the proposal
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// Denom of the asset. It could either be a native token or an IBC token
	Denom string `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty" yaml:"denom"`
	// The reward weight specifies the ratio of rewards that will be given to each alliance asset
	// It does not need to sum to 1. rate = weight / total_weight
	// Native asset is always assumed to have a weight of 1.
	RewardWeight         github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=reward_weight,json=rewardWeight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"reward_weight"`
	TakeRate             github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=take_rate,json=takeRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"take_rate"`
	RewardChangeRate     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=reward_change_rate,json=rewardChangeRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"reward_change_rate"`
	RewardChangeInterval time.Duration                          `protobuf:"bytes,7,opt,name=reward_change_interval,json=rewardChangeInterval,proto3,stdduration" json:"reward_change_interval"`
}

func (m *MsgUpdateAllianceProposal) Reset()         { *m = MsgUpdateAllianceProposal{} }
func (m *MsgUpdateAllianceProposal) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateAllianceProposal) ProtoMessage()    {}
func (*MsgUpdateAllianceProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_5518a6f5c90c8452, []int{1}
}
func (m *MsgUpdateAllianceProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateAllianceProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateAllianceProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateAllianceProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateAllianceProposal.Merge(m, src)
}
func (m *MsgUpdateAllianceProposal) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateAllianceProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateAllianceProposal.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateAllianceProposal proto.InternalMessageInfo

type MsgDeleteAllianceProposal struct {
	// the title of the update proposal
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// the description of the proposal
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Denom       string `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty" yaml:"denom"`
}

func (m *MsgDeleteAllianceProposal) Reset()         { *m = MsgDeleteAllianceProposal{} }
func (m *MsgDeleteAllianceProposal) String() string { return proto.CompactTextString(m) }
func (*MsgDeleteAllianceProposal) ProtoMessage()    {}
func (*MsgDeleteAllianceProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_5518a6f5c90c8452, []int{2}
}
func (m *MsgDeleteAllianceProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDeleteAllianceProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDeleteAllianceProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDeleteAllianceProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDeleteAllianceProposal.Merge(m, src)
}
func (m *MsgDeleteAllianceProposal) XXX_Size() int {
	return m.Size()
}
func (m *MsgDeleteAllianceProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDeleteAllianceProposal.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDeleteAllianceProposal proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgCreateAllianceProposal)(nil), "alliance.alliance.MsgCreateAllianceProposal")
	proto.RegisterType((*MsgUpdateAllianceProposal)(nil), "alliance.alliance.MsgUpdateAllianceProposal")
	proto.RegisterType((*MsgDeleteAllianceProposal)(nil), "alliance.alliance.MsgDeleteAllianceProposal")
}

func init() { proto.RegisterFile("alliance/gov.proto", fileDescriptor_5518a6f5c90c8452) }

var fileDescriptor_5518a6f5c90c8452 = []byte{
	// 491 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x94, 0x31, 0x6f, 0xd3, 0x40,
	0x14, 0xc7, 0x6d, 0x9a, 0xa4, 0xe9, 0xb5, 0x48, 0xed, 0x11, 0x81, 0xdb, 0xc1, 0x8e, 0x22, 0x54,
	0x75, 0xe9, 0x19, 0xc1, 0xd6, 0x8d, 0x34, 0x0b, 0x20, 0x24, 0x64, 0x84, 0x10, 0x15, 0x52, 0x74,
	0xb1, 0x1f, 0x17, 0xab, 0x67, 0x9f, 0x75, 0xbe, 0xb4, 0xe4, 0x03, 0x20, 0x31, 0x32, 0x32, 0xf6,
	0x83, 0xf0, 0x01, 0x3a, 0x76, 0x44, 0x0c, 0xa1, 0x4a, 0x16, 0x66, 0x3e, 0x01, 0xba, 0xb3, 0x93,
	0x3a, 0xb0, 0x55, 0x94, 0x81, 0xc9, 0xef, 0xbd, 0xff, 0xf3, 0xcf, 0x4f, 0xef, 0xff, 0x64, 0x84,
	0x29, 0xe7, 0x31, 0x4d, 0x43, 0xf0, 0x99, 0x38, 0x21, 0x99, 0x14, 0x4a, 0xe0, 0xad, 0x79, 0x8d,
	0xcc, 0x83, 0x9d, 0x7b, 0x8b, 0xb6, 0x85, 0x66, 0x7a, 0x77, 0x5a, 0x4c, 0x30, 0x61, 0x42, 0x5f,
	0x47, 0x65, 0xd5, 0x65, 0x42, 0x30, 0x0e, 0xbe, 0xc9, 0x06, 0xa3, 0x77, 0x7e, 0x34, 0x92, 0x54,
	0xc5, 0x22, 0x2d, 0xf4, 0xce, 0x97, 0x1a, 0xda, 0x7e, 0x9e, 0xb3, 0x43, 0x09, 0x54, 0xc1, 0xe3,
	0x92, 0xf8, 0x42, 0x8a, 0x4c, 0xe4, 0x94, 0xe3, 0x16, 0xaa, 0xab, 0x58, 0x71, 0x70, 0xec, 0xb6,
	0xbd, 0xb7, 0x16, 0x14, 0x09, 0x6e, 0xa3, 0xf5, 0x08, 0xf2, 0x50, 0xc6, 0x99, 0x06, 0x39, 0xb7,
	0x8c, 0x56, 0x2d, 0xe1, 0x5d, 0x54, 0x8f, 0x20, 0x15, 0x89, 0xb3, 0xa2, 0xb5, 0xee, 0xe6, 0xcf,
	0x89, 0xb7, 0x31, 0xa6, 0x09, 0x3f, 0xe8, 0x98, 0x72, 0x27, 0x28, 0x64, 0xfc, 0x12, 0xdd, 0x96,
	0x70, 0x4a, 0x65, 0xd4, 0x3f, 0x85, 0x98, 0x0d, 0x95, 0x53, 0x33, 0xfd, 0xe4, 0x7c, 0xe2, 0x59,
	0xdf, 0x26, 0xde, 0x2e, 0x8b, 0xd5, 0x70, 0x34, 0x20, 0xa1, 0x48, 0xfc, 0x50, 0xe4, 0x89, 0xc8,
	0xcb, 0xc7, 0x7e, 0x1e, 0x1d, 0xfb, 0x6a, 0x9c, 0x41, 0x4e, 0x7a, 0x10, 0x06, 0x1b, 0x05, 0xe4,
	0xb5, 0x61, 0xe0, 0x23, 0x74, 0x67, 0x09, 0xda, 0x97, 0x34, 0x65, 0xe0, 0xd4, 0xdb, 0xf6, 0xde,
	0xfa, 0xc3, 0xfb, 0xe4, 0x8f, 0x95, 0x92, 0xa0, 0xf2, 0x76, 0xa0, 0x7b, 0xbb, 0x35, 0x3d, 0x40,
	0xb0, 0x25, 0x7f, 0x17, 0xf0, 0x33, 0xb4, 0xa6, 0xe8, 0x31, 0xf4, 0x25, 0x55, 0xe0, 0x34, 0xae,
	0x35, 0x6c, 0x53, 0x03, 0x02, 0xaa, 0x00, 0xbf, 0x45, 0xb8, 0x1c, 0x34, 0x1c, 0x6a, 0x7a, 0x41,
	0x5d, 0xbd, 0x16, 0x75, 0xb3, 0x20, 0x1d, 0x1a, 0x90, 0xa1, 0xbf, 0x41, 0x77, 0x97, 0xe9, 0x71,
	0xaa, 0x40, 0x9e, 0x50, 0xee, 0x34, 0xcd, 0x26, 0xb6, 0x49, 0x71, 0x1a, 0x64, 0x7e, 0x1a, 0xa4,
	0x57, 0x9e, 0x46, 0xb7, 0xa9, 0x3f, 0xfe, 0xf9, 0xbb, 0x67, 0x07, 0xad, 0x2a, 0xf6, 0x49, 0x09,
	0x38, 0x68, 0x7e, 0x3c, 0xf3, 0xac, 0x1f, 0x67, 0x9e, 0xd5, 0xb9, 0x5c, 0x31, 0xe7, 0xf3, 0x2a,
	0x8b, 0xfe, 0x9b, 0xf3, 0x59, 0xb2, 0xb8, 0x7e, 0x23, 0x16, 0x37, 0x6e, 0xdc, 0xe2, 0xd5, 0xbf,
	0x67, 0xf1, 0x07, 0xdb, 0x58, 0xdc, 0x03, 0x0e, 0xff, 0xde, 0xe2, 0xab, 0x39, 0xba, 0x4f, 0xcf,
	0xa7, 0xae, 0x7d, 0x31, 0x75, 0xed, 0xcb, 0xa9, 0x6b, 0x7f, 0x9a, 0xb9, 0xd6, 0xc5, 0xcc, 0xb5,
	0xbe, 0xce, 0x5c, 0xeb, 0xe8, 0x41, 0x65, 0x81, 0x0a, 0xa4, 0xa4, 0xfb, 0x89, 0x48, 0x61, 0xbc,
	0xf8, 0x41, 0xfa, 0xef, 0xaf, 0x42, 0xb3, 0xce, 0x41, 0xc3, 0x2c, 0xe4, 0xd1, 0xaf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xf3, 0xa6, 0x16, 0x4b, 0x74, 0x05, 0x00, 0x00,
}

func (m *MsgCreateAllianceProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateAllianceProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateAllianceProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.RewardChangeInterval, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.RewardChangeInterval):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintGov(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x42
	{
		size := m.RewardChangeRate.Size()
		i -= size
		if _, err := m.RewardChangeRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGov(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.TakeRate.Size()
		i -= size
		if _, err := m.TakeRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGov(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size, err := m.RewardWeightRange.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGov(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.RewardWeight.Size()
		i -= size
		if _, err := m.RewardWeight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGov(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintGov(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintGov(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintGov(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUpdateAllianceProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateAllianceProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateAllianceProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n3, err3 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.RewardChangeInterval, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.RewardChangeInterval):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintGov(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x3a
	{
		size := m.RewardChangeRate.Size()
		i -= size
		if _, err := m.RewardChangeRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGov(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.TakeRate.Size()
		i -= size
		if _, err := m.TakeRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGov(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.RewardWeight.Size()
		i -= size
		if _, err := m.RewardWeight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGov(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintGov(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintGov(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintGov(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgDeleteAllianceProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDeleteAllianceProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDeleteAllianceProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintGov(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintGov(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintGov(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGov(dAtA []byte, offset int, v uint64) int {
	offset -= sovGov(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgCreateAllianceProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovGov(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovGov(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovGov(uint64(l))
	}
	l = m.RewardWeight.Size()
	n += 1 + l + sovGov(uint64(l))
	l = m.RewardWeightRange.Size()
	n += 1 + l + sovGov(uint64(l))
	l = m.TakeRate.Size()
	n += 1 + l + sovGov(uint64(l))
	l = m.RewardChangeRate.Size()
	n += 1 + l + sovGov(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.RewardChangeInterval)
	n += 1 + l + sovGov(uint64(l))
	return n
}

func (m *MsgUpdateAllianceProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovGov(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovGov(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovGov(uint64(l))
	}
	l = m.RewardWeight.Size()
	n += 1 + l + sovGov(uint64(l))
	l = m.TakeRate.Size()
	n += 1 + l + sovGov(uint64(l))
	l = m.RewardChangeRate.Size()
	n += 1 + l + sovGov(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.RewardChangeInterval)
	n += 1 + l + sovGov(uint64(l))
	return n
}

func (m *MsgDeleteAllianceProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovGov(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovGov(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovGov(uint64(l))
	}
	return n
}

func sovGov(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGov(x uint64) (n int) {
	return sovGov(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgCreateAllianceProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGov
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgCreateAllianceProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateAllianceProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardWeight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardWeightRange", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardWeightRange.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TakeRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TakeRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardChangeRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardChangeRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardChangeInterval", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.RewardChangeInterval, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGov(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGov
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgUpdateAllianceProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGov
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgUpdateAllianceProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateAllianceProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardWeight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TakeRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TakeRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardChangeRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardChangeRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardChangeInterval", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.RewardChangeInterval, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGov(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGov
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgDeleteAllianceProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGov
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgDeleteAllianceProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDeleteAllianceProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGov
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGov
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGov
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGov(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGov
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGov(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGov
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGov
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGov
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGov
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGov
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGov
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGov        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGov          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGov = fmt.Errorf("proto: unexpected end of group")
)
