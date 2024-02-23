package tests_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/terra-money/alliance/x/alliance/types"
)

type ProposalWrapper struct {
	Prop govtypes.Content
}

func TestMarshalJSONMsgs(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)

	testCases := []struct {
		name      string
		input     sdk.Msg
		strOutput string
	}{
		{
			"Msg Delegate",
			types.NewMsgDelegate("delegator", "validator", sdk.NewCoin("Alliance", math.NewInt(1000000000000000000))),
			`{"delegator_address":"delegator","validator_address":"validator","amount":{"denom":"Alliance","amount":"1000000000000000000"}}`,
		},
	}

	for _, tc := range testCases {
		bz, err := cdc.MarshalJSON(tc.input)
		assert.NoError(t, err)
		assert.Equal(t, tc.strOutput, string(bz))

		var msgDelegate types.MsgDelegate
		assert.NoError(t, cdc.UnmarshalJSON(bz, &msgDelegate))

		assert.Equal(t, tc.input, &msgDelegate)
	}
}

func TestProposalsContent(t *testing.T) {
	cases := map[string]struct {
		p     govtypes.Content
		title string
		desc  string
		typ   string
		str   string
	}{
		"msg_create_alliance_proposal": {
			p:     types.NewMsgCreateAllianceProposal("Alliance1", "Alliance with 1", "ibc/denom1", math.LegacyNewDec(1), types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)}, math.LegacyNewDec(1), math.LegacyNewDec(1), time.Second),
			title: "Alliance1",
			desc:  "Alliance with 1",
			typ:   "msg_create_alliance_proposal",
			str:   "title:\"Alliance1\" description:\"Alliance with 1\" denom:\"ibc/denom1\" reward_weight:\"1000000000000000000\" take_rate:\"1000000000000000000\" reward_change_rate:\"1000000000000000000\" reward_change_interval:<seconds:1 > reward_weight_range:<min:\"0\" max:\"5000000000000000000\" > ",
		},
		"msg_update_alliance_proposal": {
			p:     types.NewMsgUpdateAllianceProposal("Alliance2", "Alliance with 2", "ibc/denom2", math.LegacyNewDec(2), types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)}, math.LegacyNewDec(2), math.LegacyNewDec(2), time.Hour),
			title: "Alliance2",
			desc:  "Alliance with 2",
			typ:   "msg_update_alliance_proposal",
			str:   "title:\"Alliance2\" description:\"Alliance with 2\" denom:\"ibc/denom2\" reward_weight:\"2000000000000000000\" take_rate:\"2000000000000000000\" reward_change_rate:\"2000000000000000000\" reward_change_interval:<seconds:3600 > reward_weight_range:<min:\"0\" max:\"5000000000000000000\" > ",
		},
		"msg_delete_alliance_proposal": {
			p:     types.NewMsgDeleteAllianceProposal("test", "abcd", "ibc/denom"),
			title: "test",
			desc:  "abcd",
			typ:   "msg_delete_alliance_proposal",
			str:   "title:\"test\" description:\"abcd\" denom:\"ibc/denom\" ",
		},
	}

	cdc := codec.NewLegacyAmino()
	govtypes.RegisterLegacyAminoCodec(cdc)
	types.RegisterLegacyAminoCodec(cdc)

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.title, tc.p.GetTitle())
			assert.Equal(t, tc.desc, tc.p.GetDescription())
			assert.Equal(t, tc.typ, tc.p.ProposalType())
			assert.Equal(t, "alliance", tc.p.ProposalRoute())
			assert.Equal(t, tc.str, tc.p.String())

			// try to encode and decode type to ensure codec works
			wrap := ProposalWrapper{tc.p}
			bz, err := cdc.Marshal(&wrap)
			require.NoError(t, err)
			unwrap := ProposalWrapper{}
			err = cdc.Unmarshal(bz, &unwrap)
			require.NoError(t, err)

			// all methods should look the same
			assert.Equal(t, tc.title, unwrap.Prop.GetTitle())
			assert.Equal(t, tc.desc, unwrap.Prop.GetDescription())
			assert.Equal(t, tc.typ, unwrap.Prop.ProposalType())
			assert.Equal(t, "alliance", unwrap.Prop.ProposalRoute())
			assert.Equal(t, tc.str, unwrap.Prop.String())
		})
	}
}

func TestInvalidProposalsContent(t *testing.T) {
	byteArray := []byte{'a', 'l', 'l', 'i', 'a', 'n', 'c', 'e', 0, '2'}
	invalidDenom := string(byteArray)
	cases := map[string]struct {
		p     govtypes.Content
		title string
		desc  string
		typ   string
		str   string
	}{
		"msg_create_alliance_proposal": {
			p:     types.NewMsgCreateAllianceProposal("Alliance1", "Alliance with 1", "ibc/denom1", math.LegacyNewDec(1), types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)}, math.LegacyNewDec(1), math.LegacyNewDec(1), -time.Second),
			title: "Alliance1",
			desc:  "Alliance with 1",
			typ:   "msg_create_alliance_proposal",
		},
		"msg_create_alliance_proposal_invalid_denom": {
			p:     types.NewMsgCreateAllianceProposal("Alliance1", "Alliance with 1", invalidDenom, math.LegacyNewDec(1), types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)}, math.LegacyNewDec(1), math.LegacyNewDec(1), time.Second),
			title: "Alliance1",
			desc:  "Alliance with 1",
			typ:   "msg_create_alliance_proposal",
		},
		"msg_update_alliance_proposal": {
			p:     types.NewMsgUpdateAllianceProposal("Alliance2", "Alliance with 2", "ibc/denom2", math.LegacyNewDec(2), types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)}, math.LegacyNewDec(2), math.LegacyNewDec(2), -time.Hour),
			title: "Alliance2",
			desc:  "Alliance with 2",
			typ:   "msg_update_alliance_proposal",
		},
	}

	cdc := codec.NewLegacyAmino()
	govtypes.RegisterLegacyAminoCodec(cdc)
	types.RegisterLegacyAminoCodec(cdc)

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := tc.p.ValidateBasic()
			require.Error(t, err)
		})
	}
}

func TestAminoJSON(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	govtypes.RegisterLegacyAminoCodec(cdc)
	types.RegisterLegacyAminoCodec(cdc)
	legacytx.RegressionTestingAminoCodec = cdc

	msgDelegate := types.NewMsgDelegate("delegator", "validator", sdk.NewCoin("Alliance", math.NewInt(1000000000000000000)))
	require.Equal(t,
		`{"account_number":"1","chain_id":"foo","fee":{"amount":[],"gas":"0"},"memo":"memo","msgs":[{"type":"alliance/MsgDelegate","value":{"amount":{"amount":"1000000000000000000","denom":"Alliance"},"delegator_address":"delegator","validator_address":"validator"}}],"sequence":"1","timeout_height":"1"}`,
		string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{msgDelegate}, "memo")),
	)

	msgUndelegate := types.NewMsgUndelegate("delegator", "validator", sdk.NewCoin("Alliance", math.NewInt(1000000000000000000)))
	require.Equal(t,
		`{"account_number":"1","chain_id":"foo","fee":{"amount":[],"gas":"0"},"memo":"memo","msgs":[{"type":"alliance/MsgUndelegate","value":{"amount":{"amount":"1000000000000000000","denom":"Alliance"},"delegator_address":"delegator","validator_address":"validator"}}],"sequence":"1","timeout_height":"1"}`,
		string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{msgUndelegate}, "memo")),
	)

	msgRedelegate := types.NewMsgRedelegate("delegator", "validator", "validator1", sdk.NewCoin("Alliance", math.NewInt(1000000000000000000)))
	require.Equal(t,
		`{"account_number":"1","chain_id":"foo","fee":{"amount":[],"gas":"0"},"memo":"memo","msgs":[{"type":"alliance/MsgRedelegate","value":{"amount":{"amount":"1000000000000000000","denom":"Alliance"},"delegator_address":"delegator","validator_dst_address":"validator1","validator_src_address":"validator"}}],"sequence":"1","timeout_height":"1"}`,
		string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{msgRedelegate}, "memo")),
	)

	msgClaimDelegationRewards := types.NewMsgClaimDelegationRewards("delegator", "validator", "Alliance")
	require.Equal(t,
		`{"account_number":"1","chain_id":"foo","fee":{"amount":[],"gas":"0"},"memo":"memo","msgs":[{"type":"alliance/MsgClaimDelegationRewards","value":{"delegator_address":"delegator","denom":"Alliance","validator_address":"validator"}}],"sequence":"1","timeout_height":"1"}`,
		string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{msgClaimDelegationRewards}, "memo")),
	)
}
