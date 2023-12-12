package ics20

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile/contracts/ics20/testdata"
	"github.com/ava-labs/subnet-evm/precompile/contracts/ics20/testdata/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/fungible_token_cases.json
	rawFungibleTokenData   []byte
	fungibleTokenTestCases []FungibleTokenTestCase
)

type FungibleTokenTestCase struct {
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input,omitempty"`
	Error *string         `json:"error,omitempty"`
}

func init() {
	if err := json.Unmarshal(rawFungibleTokenData, &fungibleTokenTestCases); err != nil {
		panic(err)
	}
}

func TestFungibleTokenPacketDataToABI(t *testing.T) {
	_, _, cdc, err := testdata.NewCodecEnv()
	require.NoError(t, err, "can't create enviroment")

	for i := range fungibleTokenTestCases {
		testcase := fungibleTokenTestCases[i]
		t.Run(testcase.Name, func(t *testing.T) {
			actual, err := FungibleTokenPacketDataToABI(testcase.Input)
			if testcase.Error != nil {
				assert.ErrorContains(t, err, *testcase.Error)
				return
			}
			require.NoError(t, err)

			var jsondata FungibleTokenPacketData
			require.NoError(t, json.Unmarshal(testcase.Input, &jsondata), "can't parse json input")
			amount, ok := math.ParseBig256(jsondata.Amount)
			require.True(t, ok)

			sender, err := common.ParseHexOrString(jsondata.Sender)
			require.NoError(t, err)

			receiver, err := common.ParseHexOrString(jsondata.Receiver)
			require.NoError(t, err)

			expected, err := cdc.Encode(nil, codec.FungibleTokenPacketData{
				Denom:    jsondata.Denom,
				Amount:   amount,
				Sender:   sender,
				Receiver: receiver,
				Memo:     []byte(jsondata.Memo),
			})
			require.NoError(t, err, "can't call test contract 'codec'")

			assert.Equal(t, expected, actual)
		})
	}
}
