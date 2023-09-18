package ibc

import (
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	tendermint "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestQueryClientState(t *testing.T) {
	cs := &tendermint.ClientState{
		ChainId:         "chainID",
		TrustLevel:      tendermint.Fraction{Numerator: 1, Denominator: 2},
		TrustingPeriod:  100,
		UnbondingPeriod: 200,
		MaxClockDrift:   300,
		FrozenHeight:    types.Height{RevisionNumber: 5, RevisionHeight: 6},
		LatestHeight:    types.Height{RevisionNumber: 1, RevisionHeight: 1},
	}
	cdByte, err := cs.Marshal()
	if err != nil {
		t.Error(err)
	}
	packedOutput, err := PackQueryClientStateOutput(cdByte)
	if err != nil {
		t.Error(err)
	}

	tests := map[string]testutils.PrecompileTest{
		"invalid clientID": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				clientID := ""
				input, err := PackQueryClientStateInput(QueryClientStateInput{ClientID: clientID})
				require.NoError(t, err)
				return input
			},
			SuppliedGas: QueryClientStateGasCost,
			ReadOnly:    false,
			ExpectedErr: "error loading client state, err: empty precompile state",
		},
		"ClientState not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				input, err := PackQueryClientStateInput(QueryClientStateInput{ClientID: "clientID"})
				require.NoError(t, err)
				return input
			},
			SuppliedGas: QueryClientStateGasCost,
			ReadOnly:    false,
			ExpectedErr: "error loading client state, err: empty precompile state",
		},
		"success": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				input, err := PackQueryClientStateInput(QueryClientStateInput{ClientID: "clientID"})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				err := SetClientState(state, "clientID", cs)
				if err != nil {
					t.Error(err)
				}
			},
			SuppliedGas: QueryClientStateGasCost,
			ReadOnly:    false,
			ExpectedRes: packedOutput,
		},
	}
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}

func TestQueryConsensusState(t *testing.T) {
	cs := &tendermint.ClientState{
		ChainId:         "chainID",
		TrustLevel:      tendermint.Fraction{Numerator: 1, Denominator: 2},
		TrustingPeriod:  100,
		UnbondingPeriod: 200,
		MaxClockDrift:   300,
		FrozenHeight:    types.Height{RevisionNumber: 5, RevisionHeight: 6},
		LatestHeight:    types.Height{RevisionNumber: 1, RevisionHeight: 1},
	}

	consState := &tendermint.ConsensusState{
		Timestamp: time.Now(),
		Root:      commitmenttypes.MerkleRoot{},
	}
	cdByte, err := consState.Marshal()
	if err != nil {
		t.Error(err)
	}

	packedOutput, err := PackQueryConsensusStateOutput(cdByte)
	if err != nil {
		t.Error(err)
	}

	tests := map[string]testutils.PrecompileTest{
		"invalid clientID": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				clientID := ""
				input, err := PackQueryConsensusStateInput(QueryConsensusStateInput{ClientID: clientID})
				require.NoError(t, err)
				return input
			},
			SuppliedGas: QueryConsensusStateGasCost,
			ReadOnly:    false,
			ExpectedErr: "error loading client state, err: empty precompile state",
		},
		"ConsensusState not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				clientID := "clientID"
				input, err := PackQueryConsensusStateInput(QueryConsensusStateInput{ClientID: clientID})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				err := SetClientState(state, "clientID", cs)
				if err != nil {
					t.Error(err)
				}
			},
			SuppliedGas: QueryConsensusStateGasCost,
			ReadOnly:    false,
			ExpectedErr: "error loading consensus state, err: empty precompile state",
		},
		"success": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				input, err := PackQueryConsensusStateInput(QueryConsensusStateInput{ClientID: "clientID"})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				err := SetClientState(state, "clientID", cs)
				if err != nil {
					t.Error(err)
				}
				err = SetConsensusState(state, "clientID", cs.LatestHeight, consState)
				if err != nil {
					t.Error(err)
				}
			},
			SuppliedGas: QueryConsensusStateGasCost,
			ReadOnly:    false,
			ExpectedRes: packedOutput,
		},
	}
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}

func TestQueryConnection(t *testing.T) {
	connection := &connectiontypes.ConnectionEnd{
		ClientId:     "clientId",
		Versions:     []*connectiontypes.Version{},
		State:        connectiontypes.OPEN,
		Counterparty: connectiontypes.Counterparty{},
		DelayPeriod:  100,
	}
	cdByte, err := connection.Marshal()
	if err != nil {
		t.Error(err)
	}

	packedOutput, err := PackQueryConnectionOutput(cdByte)
	if err != nil {
		t.Error(err)
	}

	tests := map[string]testutils.PrecompileTest{
		"invalid ConnectionID": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				connectionID := ""
				input, err := PackQueryConnectionInput(QueryConnectionInput{ConnectionID: connectionID})
				require.NoError(t, err)
				return input
			},
			SuppliedGas: QueryConnectionGasCost,
			ReadOnly:    false,
			ExpectedErr: "error loading connection, err: empty precompile state",
		},
		"connection not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				ConnectionID := ""
				input, err := PackQueryConnectionInput(QueryConnectionInput{ConnectionID: ConnectionID})
				require.NoError(t, err)
				return input
			},
			SuppliedGas: QueryConnectionGasCost,
			ReadOnly:    false,
			ExpectedErr: "error loading connection, err: empty precompile state",
		},
		"success": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				input, err := PackQueryConnectionInput(QueryConnectionInput{ConnectionID: "ConnectionID"})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				err = SetConnection(state, "ConnectionID", connection)
				if err != nil {
					t.Error(err)
				}
			},
			SuppliedGas: QueryConnectionGasCost,
			ReadOnly:    false,
			ExpectedRes: packedOutput,
		},
	}
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}

func TestQueryChannel(t *testing.T) {
	channel := &channeltypes.Channel{
		State:          channeltypes.OPEN,
		Ordering:       channeltypes.NONE,
		Counterparty:   channeltypes.Counterparty{},
		ConnectionHops: []string{"ConnectionHops"},
		Version:        "version",
	}
	cdByte, err := channel.Marshal()
	if err != nil {
		t.Error(err)
	}

	packedOutput, err := PackQueryChannelOutput(cdByte)
	if err != nil {
		t.Error(err)
	}

	tests := map[string]testutils.PrecompileTest{
		"invalid ChannelID": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				channelID := ""
				portID := ""
				input, err := PackQueryChannelInput(QueryChannelInput{PortID: portID, ChannelID: channelID})
				require.NoError(t, err)
				return input
			},
			SuppliedGas: QueryChannelGasCost,
			ReadOnly:    false,
			ExpectedErr: "error loading channel, err: empty precompile state",
		},
		"channel not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				channelID := "channelID"
				portID := "portID"
				input, err := PackQueryChannelInput(QueryChannelInput{PortID: portID, ChannelID: channelID})
				require.NoError(t, err)
				return input
			},
			SuppliedGas: QueryChannelGasCost,
			ReadOnly:    false,
			ExpectedErr: "error loading channel, err: empty precompile state",
		},
		"success": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				channelID := "channelID"
				portID := "portID"
				input, err := PackQueryChannelInput(QueryChannelInput{PortID: portID, ChannelID: channelID})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				err = SetChannel(state, "portID", "channelID", channel)
				if err != nil {
					t.Error(err)
				}
			},
			SuppliedGas: QueryChannelGasCost,
			ReadOnly:    false,
			ExpectedRes: packedOutput,
		},
	}
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}
