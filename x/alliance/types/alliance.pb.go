// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: alliance/alliance.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

type RewardWeightRange struct {
	Min cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=min,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"min"`
	Max cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=max,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"max"`
}

func (m *RewardWeightRange) Reset()         { *m = RewardWeightRange{} }
func (m *RewardWeightRange) String() string { return proto.CompactTextString(m) }
func (*RewardWeightRange) ProtoMessage()    {}
func (*RewardWeightRange) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7dbf17f28cd0f90, []int{0}
}
func (m *RewardWeightRange) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RewardWeightRange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RewardWeightRange.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RewardWeightRange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RewardWeightRange.Merge(m, src)
}
func (m *RewardWeightRange) XXX_Size() int {
	return m.Size()
}
func (m *RewardWeightRange) XXX_DiscardUnknown() {
	xxx_messageInfo_RewardWeightRange.DiscardUnknown(m)
}

var xxx_messageInfo_RewardWeightRange proto.InternalMessageInfo

// key: denom value: AllianceAsset
type AllianceAsset struct {
	// Denom of the asset. It could either be a native token or an IBC token
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty" yaml:"denom"`
	// The reward weight specifies the ratio of rewards that will be given to each alliance asset
	// It does not need to sum to 1. rate = weight / total_weight
	// Native asset is always assumed to have a weight of 1.s
	RewardWeight cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=reward_weight,json=rewardWeight,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"reward_weight"`
	// A positive take rate is used for liquid staking derivatives. It defines an rate that is applied per take_rate_interval
	// that will be redirected to the distribution rewards pool
	TakeRate             cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=take_rate,json=takeRate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"take_rate"`
	TotalTokens          cosmossdk_io_math.Int       `protobuf:"bytes,4,opt,name=total_tokens,json=totalTokens,proto3,customtype=cosmossdk.io/math.Int" json:"total_tokens"`
	TotalValidatorShares cosmossdk_io_math.LegacyDec `protobuf:"bytes,5,opt,name=total_validator_shares,json=totalValidatorShares,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"total_validator_shares"`
	RewardStartTime      time.Time                   `protobuf:"bytes,6,opt,name=reward_start_time,json=rewardStartTime,proto3,stdtime" json:"reward_start_time"`
	RewardChangeRate     cosmossdk_io_math.LegacyDec `protobuf:"bytes,7,opt,name=reward_change_rate,json=rewardChangeRate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"reward_change_rate"`
	RewardChangeInterval time.Duration               `protobuf:"bytes,8,opt,name=reward_change_interval,json=rewardChangeInterval,proto3,stdduration" json:"reward_change_interval"`
	LastRewardChangeTime time.Time                   `protobuf:"bytes,9,opt,name=last_reward_change_time,json=lastRewardChangeTime,proto3,stdtime" json:"last_reward_change_time"`
	// set a bound of weight range to limit how much reward weights can scale.
	RewardWeightRange RewardWeightRange `protobuf:"bytes,10,opt,name=reward_weight_range,json=rewardWeightRange,proto3" json:"reward_weight_range"`
	// flag to check if an asset has completed the initialization process after the reward delay
	IsInitialized bool `protobuf:"varint,11,opt,name=is_initialized,json=isInitialized,proto3" json:"is_initialized,omitempty"`
}

func (m *AllianceAsset) Reset()         { *m = AllianceAsset{} }
func (m *AllianceAsset) String() string { return proto.CompactTextString(m) }
func (*AllianceAsset) ProtoMessage()    {}
func (*AllianceAsset) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7dbf17f28cd0f90, []int{1}
}
func (m *AllianceAsset) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AllianceAsset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AllianceAsset.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AllianceAsset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AllianceAsset.Merge(m, src)
}
func (m *AllianceAsset) XXX_Size() int {
	return m.Size()
}
func (m *AllianceAsset) XXX_DiscardUnknown() {
	xxx_messageInfo_AllianceAsset.DiscardUnknown(m)
}

var xxx_messageInfo_AllianceAsset proto.InternalMessageInfo

type RewardWeightChangeSnapshot struct {
	PrevRewardWeight cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=prev_reward_weight,json=prevRewardWeight,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"prev_reward_weight"`
	RewardHistories  []RewardHistory             `protobuf:"bytes,2,rep,name=reward_histories,json=rewardHistories,proto3" json:"reward_histories"`
}

func (m *RewardWeightChangeSnapshot) Reset()         { *m = RewardWeightChangeSnapshot{} }
func (m *RewardWeightChangeSnapshot) String() string { return proto.CompactTextString(m) }
func (*RewardWeightChangeSnapshot) ProtoMessage()    {}
func (*RewardWeightChangeSnapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7dbf17f28cd0f90, []int{2}
}
func (m *RewardWeightChangeSnapshot) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RewardWeightChangeSnapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RewardWeightChangeSnapshot.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RewardWeightChangeSnapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RewardWeightChangeSnapshot.Merge(m, src)
}
func (m *RewardWeightChangeSnapshot) XXX_Size() int {
	return m.Size()
}
func (m *RewardWeightChangeSnapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_RewardWeightChangeSnapshot.DiscardUnknown(m)
}

var xxx_messageInfo_RewardWeightChangeSnapshot proto.InternalMessageInfo

func init() {
	proto.RegisterType((*RewardWeightRange)(nil), "alliance.RewardWeightRange")
	proto.RegisterType((*AllianceAsset)(nil), "alliance.AllianceAsset")
	proto.RegisterType((*RewardWeightChangeSnapshot)(nil), "alliance.RewardWeightChangeSnapshot")
}

func init() { proto.RegisterFile("alliance/alliance.proto", fileDescriptor_f7dbf17f28cd0f90) }

var fileDescriptor_f7dbf17f28cd0f90 = []byte{
	// 682 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0xcd, 0x6e, 0xd3, 0x4a,
	0x14, 0xc7, 0xe3, 0x7e, 0xdd, 0x74, 0xd2, 0xde, 0xdb, 0xfa, 0xa6, 0xad, 0x9b, 0x4a, 0x76, 0x14,
	0x04, 0x8a, 0x84, 0x6a, 0xa3, 0x22, 0x36, 0x5d, 0xd1, 0xd0, 0x45, 0x82, 0x10, 0x02, 0xb7, 0x02,
	0x15, 0x16, 0xd6, 0x34, 0x19, 0xec, 0x51, 0x6d, 0x4f, 0x34, 0x73, 0xfa, 0x11, 0x9e, 0x00, 0x89,
	0x4d, 0x97, 0x2c, 0xfb, 0x10, 0xbc, 0x01, 0x9b, 0xee, 0xa8, 0x58, 0x21, 0x16, 0x05, 0xb5, 0x1b,
	0xd6, 0x3c, 0x01, 0x9a, 0xf1, 0xb8, 0x4d, 0x9a, 0x4d, 0xb3, 0xf3, 0xcc, 0x99, 0xff, 0x6f, 0xfe,
	0xe7, 0xcc, 0x39, 0x46, 0x4b, 0x38, 0x8e, 0x29, 0x4e, 0xdb, 0xc4, 0xcb, 0x3f, 0xdc, 0x2e, 0x67,
	0xc0, 0xcc, 0x62, 0xbe, 0xae, 0x94, 0x43, 0x16, 0x32, 0xb5, 0xe9, 0xc9, 0xaf, 0x2c, 0x5e, 0x59,
	0x6e, 0x33, 0x91, 0x30, 0x11, 0x64, 0x81, 0x6c, 0xa1, 0x43, 0x0b, 0x57, 0xcc, 0x2e, 0xe6, 0x38,
	0xc9, 0xb7, 0xed, 0x90, 0xb1, 0x30, 0x26, 0x9e, 0x5a, 0xed, 0xee, 0xbf, 0xf3, 0x3a, 0xfb, 0x1c,
	0x03, 0x65, 0xa9, 0x8e, 0x3b, 0x37, 0xe3, 0x40, 0x13, 0x22, 0x00, 0x27, 0xdd, 0xec, 0x40, 0xed,
	0xa3, 0x81, 0xe6, 0x7d, 0x72, 0x88, 0x79, 0xe7, 0x35, 0xa1, 0x61, 0x04, 0x3e, 0x4e, 0x43, 0x62,
	0x3e, 0x42, 0xe3, 0x09, 0x4d, 0x2d, 0xa3, 0x6a, 0xd4, 0xa7, 0x1b, 0x77, 0x4e, 0xcf, 0x9d, 0xc2,
	0x8f, 0x73, 0x67, 0x25, 0x33, 0x24, 0x3a, 0x7b, 0x2e, 0x65, 0x5e, 0x82, 0x21, 0x72, 0x9f, 0x91,
	0x10, 0xb7, 0x7b, 0x9b, 0xa4, 0xed, 0xcb, 0xf3, 0x4a, 0x86, 0x8f, 0xac, 0xb1, 0x51, 0x64, 0xf8,
	0x68, 0xbd, 0xf8, 0xe1, 0xc4, 0x29, 0xfc, 0x3e, 0x71, 0x0a, 0xb5, 0xaf, 0x53, 0x68, 0x76, 0x43,
	0x27, 0xba, 0x21, 0x04, 0x01, 0xf3, 0x1e, 0x9a, 0xec, 0x90, 0x94, 0x25, 0xda, 0xcb, 0xdc, 0x9f,
	0x73, 0x67, 0xa6, 0x87, 0x93, 0x78, 0xbd, 0xa6, 0xb6, 0x6b, 0x7e, 0x16, 0x36, 0x9b, 0x68, 0x96,
	0xab, 0x34, 0x82, 0x43, 0x95, 0xc7, 0x28, 0x26, 0x66, 0x78, 0x5f, 0x01, 0xcc, 0xc7, 0x68, 0x1a,
	0xf0, 0x1e, 0x09, 0x38, 0x06, 0x62, 0x8d, 0xdf, 0x9e, 0x52, 0x94, 0x2a, 0x1f, 0x03, 0x31, 0x9f,
	0xa3, 0x19, 0x60, 0x80, 0xe3, 0x00, 0xd8, 0x1e, 0x49, 0x85, 0x35, 0xa1, 0x20, 0xf7, 0x35, 0x64,
	0x61, 0x18, 0xd2, 0x4a, 0xe1, 0xdb, 0xe7, 0x55, 0xa4, 0x1f, 0xbc, 0x95, 0x82, 0x5f, 0x52, 0x80,
	0x6d, 0xa5, 0x37, 0x77, 0xd0, 0x62, 0xc6, 0x3b, 0xc0, 0x31, 0xed, 0x60, 0x60, 0x3c, 0x10, 0x11,
	0xe6, 0x44, 0x58, 0x93, 0xb7, 0xb7, 0x57, 0x56, 0x88, 0x57, 0x39, 0x61, 0x4b, 0x01, 0xcc, 0x17,
	0x68, 0x5e, 0x97, 0x4d, 0x00, 0xe6, 0x10, 0xc8, 0xf6, 0xb0, 0xa6, 0xaa, 0x46, 0xbd, 0xb4, 0x56,
	0x71, 0xb3, 0xde, 0x71, 0xf3, 0xde, 0x71, 0xb7, 0xf3, 0xde, 0x69, 0x14, 0xe5, 0x8d, 0xc7, 0x3f,
	0x1d, 0xc3, 0xff, 0x2f, 0x93, 0x6f, 0x49, 0xb5, 0x8c, 0x9b, 0x2f, 0x91, 0xa9, 0x89, 0xed, 0x48,
	0xf6, 0x52, 0x56, 0xc7, 0x7f, 0x6e, 0x6f, 0x74, 0x2e, 0x93, 0x3f, 0x51, 0x6a, 0x55, 0xcf, 0x1d,
	0xb4, 0x38, 0x88, 0xa4, 0x29, 0x10, 0x7e, 0x80, 0x63, 0xab, 0xa8, 0x9c, 0x2e, 0x0f, 0x39, 0xdd,
	0xd4, 0x53, 0x90, 0x19, 0xfd, 0x24, 0x8d, 0x96, 0xfb, 0xb1, 0x2d, 0x0d, 0x30, 0xdf, 0xa2, 0xa5,
	0x18, 0x0b, 0x08, 0x06, 0xf9, 0xaa, 0x0a, 0xd3, 0x23, 0x54, 0xa1, 0x2c, 0x21, 0x7e, 0xdf, 0x05,
	0xba, 0x14, 0xff, 0x0f, 0xf4, 0x64, 0xc0, 0x65, 0xc8, 0x42, 0x0a, 0xbc, 0xe2, 0x5e, 0xfd, 0x1c,
	0x86, 0xe6, 0xaf, 0x31, 0x21, 0xc9, 0xfe, 0x3c, 0x1f, 0x1a, 0xcc, 0xbb, 0xe8, 0x5f, 0x2a, 0x02,
	0x9a, 0x52, 0xa0, 0x38, 0xa6, 0xef, 0x49, 0xc7, 0x2a, 0x55, 0x8d, 0x7a, 0xd1, 0x9f, 0xa5, 0xa2,
	0x75, 0xbd, 0xd9, 0x37, 0x51, 0x5f, 0x0c, 0x54, 0xe9, 0xe7, 0x67, 0xf6, 0xb6, 0x52, 0xdc, 0x15,
	0x11, 0x03, 0xf9, 0x5a, 0x5d, 0x4e, 0x0e, 0x82, 0xc1, 0xd9, 0x19, 0x61, 0xee, 0xe7, 0xa4, 0xbc,
	0xff, 0x02, 0xb3, 0x89, 0xf4, 0x0b, 0x06, 0x11, 0x15, 0xc0, 0x38, 0x25, 0xc2, 0x1a, 0xab, 0x8e,
	0xd7, 0x4b, 0x6b, 0x4b, 0x37, 0x53, 0x6e, 0xaa, 0x03, 0x3d, 0x9d, 0xae, 0x6e, 0xa5, 0x66, 0xae,
	0xba, 0xce, 0xa2, 0xf1, 0xf4, 0xf4, 0xc2, 0x36, 0xce, 0x2e, 0x6c, 0xe3, 0xd7, 0x85, 0x6d, 0x1c,
	0x5f, 0xda, 0x85, 0xb3, 0x4b, 0xbb, 0xf0, 0xfd, 0xd2, 0x2e, 0xbc, 0x79, 0x10, 0x52, 0x88, 0xf6,
	0x77, 0xdd, 0x36, 0x4b, 0x3c, 0x20, 0x9c, 0xe3, 0xd5, 0x84, 0xa5, 0xa4, 0x77, 0xf5, 0xe7, 0xf5,
	0x8e, 0xae, 0x3f, 0xa1, 0xd7, 0x25, 0x62, 0x77, 0x4a, 0xbd, 0xe4, 0xc3, 0xbf, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x83, 0xfb, 0x94, 0x6d, 0xa6, 0x05, 0x00, 0x00,
}

func (m *RewardWeightRange) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RewardWeightRange) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RewardWeightRange) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Max.Size()
		i -= size
		if _, err := m.Max.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintAlliance(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.Min.Size()
		i -= size
		if _, err := m.Min.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintAlliance(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *AllianceAsset) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AllianceAsset) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AllianceAsset) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IsInitialized {
		i--
		if m.IsInitialized {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x58
	}
	{
		size, err := m.RewardWeightRange.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintAlliance(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.LastRewardChangeTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.LastRewardChangeTime):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintAlliance(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x4a
	n3, err3 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.RewardChangeInterval, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.RewardChangeInterval):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintAlliance(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x42
	{
		size := m.RewardChangeRate.Size()
		i -= size
		if _, err := m.RewardChangeRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintAlliance(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	n4, err4 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.RewardStartTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.RewardStartTime):])
	if err4 != nil {
		return 0, err4
	}
	i -= n4
	i = encodeVarintAlliance(dAtA, i, uint64(n4))
	i--
	dAtA[i] = 0x32
	{
		size := m.TotalValidatorShares.Size()
		i -= size
		if _, err := m.TotalValidatorShares.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintAlliance(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.TotalTokens.Size()
		i -= size
		if _, err := m.TotalTokens.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintAlliance(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.TakeRate.Size()
		i -= size
		if _, err := m.TakeRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintAlliance(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.RewardWeight.Size()
		i -= size
		if _, err := m.RewardWeight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintAlliance(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintAlliance(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RewardWeightChangeSnapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RewardWeightChangeSnapshot) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RewardWeightChangeSnapshot) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.RewardHistories) > 0 {
		for iNdEx := len(m.RewardHistories) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RewardHistories[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintAlliance(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size := m.PrevRewardWeight.Size()
		i -= size
		if _, err := m.PrevRewardWeight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintAlliance(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintAlliance(dAtA []byte, offset int, v uint64) int {
	offset -= sovAlliance(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RewardWeightRange) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Min.Size()
	n += 1 + l + sovAlliance(uint64(l))
	l = m.Max.Size()
	n += 1 + l + sovAlliance(uint64(l))
	return n
}

func (m *AllianceAsset) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovAlliance(uint64(l))
	}
	l = m.RewardWeight.Size()
	n += 1 + l + sovAlliance(uint64(l))
	l = m.TakeRate.Size()
	n += 1 + l + sovAlliance(uint64(l))
	l = m.TotalTokens.Size()
	n += 1 + l + sovAlliance(uint64(l))
	l = m.TotalValidatorShares.Size()
	n += 1 + l + sovAlliance(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.RewardStartTime)
	n += 1 + l + sovAlliance(uint64(l))
	l = m.RewardChangeRate.Size()
	n += 1 + l + sovAlliance(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.RewardChangeInterval)
	n += 1 + l + sovAlliance(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.LastRewardChangeTime)
	n += 1 + l + sovAlliance(uint64(l))
	l = m.RewardWeightRange.Size()
	n += 1 + l + sovAlliance(uint64(l))
	if m.IsInitialized {
		n += 2
	}
	return n
}

func (m *RewardWeightChangeSnapshot) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.PrevRewardWeight.Size()
	n += 1 + l + sovAlliance(uint64(l))
	if len(m.RewardHistories) > 0 {
		for _, e := range m.RewardHistories {
			l = e.Size()
			n += 1 + l + sovAlliance(uint64(l))
		}
	}
	return n
}

func sovAlliance(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAlliance(x uint64) (n int) {
	return sovAlliance(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RewardWeightRange) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAlliance
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
			return fmt.Errorf("proto: RewardWeightRange: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RewardWeightRange: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Min", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Min.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Max", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Max.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAlliance(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAlliance
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
func (m *AllianceAsset) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAlliance
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
			return fmt.Errorf("proto: AllianceAsset: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AllianceAsset: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardWeight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TakeRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TakeRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalTokens", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalTokens.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalValidatorShares", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalValidatorShares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardStartTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.RewardStartTime, dAtA[iNdEx:postIndex]); err != nil {
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
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
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
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.RewardChangeInterval, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastRewardChangeTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.LastRewardChangeTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardWeightRange", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardWeightRange.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsInitialized", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsInitialized = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipAlliance(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAlliance
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
func (m *RewardWeightChangeSnapshot) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAlliance
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
			return fmt.Errorf("proto: RewardWeightChangeSnapshot: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RewardWeightChangeSnapshot: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PrevRewardWeight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PrevRewardWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardHistories", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlliance
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
				return ErrInvalidLengthAlliance
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAlliance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardHistories = append(m.RewardHistories, RewardHistory{})
			if err := m.RewardHistories[len(m.RewardHistories)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAlliance(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAlliance
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
func skipAlliance(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAlliance
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
					return 0, ErrIntOverflowAlliance
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
					return 0, ErrIntOverflowAlliance
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
				return 0, ErrInvalidLengthAlliance
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAlliance
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAlliance
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAlliance        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAlliance          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAlliance = fmt.Errorf("proto: unexpected end of group")
)
