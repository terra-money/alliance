package types_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
			types.NewMsgDelegate("delegator", "validator", sdk.NewCoin("Alliance", sdk.NewInt(1000000000000000000))),
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
			p:     types.NewMsgCreateAllianceProposal("Alliance1", "Alliance with 1", "ibc/denom1", sdk.NewDec(1), sdk.NewDec(1), sdk.NewDec(1), time.Second),
			title: "Alliance1",
			desc:  "Alliance with 1",
			typ:   "msg_create_alliance_proposal",
			str:   "title:\"Alliance1\" description:\"Alliance with 1\" denom:\"ibc/denom1\" reward_weight:\"1000000000000000000\" take_rate:\"1000000000000000000\" reward_change_rate:\"1000000000000000000\" reward_change_interval:<seconds:1 > ",
		},
		"msg_update_alliance_proposal": {
			p:     types.NewMsgUpdateAllianceProposal("Alliance2", "Alliance with 2", "ibc/denom2", sdk.NewDec(2), sdk.NewDec(2), sdk.NewDec(2), time.Hour),
			title: "Alliance2",
			desc:  "Alliance with 2",
			typ:   "msg_update_alliance_proposal",
			str:   "title:\"Alliance2\" description:\"Alliance with 2\" denom:\"ibc/denom2\" reward_weight:\"2000000000000000000\" take_rate:\"2000000000000000000\" reward_change_rate:\"2000000000000000000\" reward_change_interval:<seconds:3600 > ",
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
