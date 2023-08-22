package ibc

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	hosttypes "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
)

const doesnotexist = "doesnotexist"

func malleateHeight(height clienttypes.Height, diff uint64) clienttypes.Height {
	return clienttypes.NewHeight(height.GetRevisionNumber(), height.GetRevisionHeight()+diff)
}

func TestChanOpenInit(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	var path *ibctesting.Path

	tests := map[string]testutils.PrecompileTest{
		"success": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
			},
			ExpectedRes: make([]byte, 0),
		},
		"connection doesn't exist": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				path.EndpointA.ConnectionID = "connection-0"
				path.EndpointB.ConnectionID = "connection-0"
			},
			ExpectedErr: "single version must be negotiated on connection before opening channel",
		},
	}
	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			coordinator = ibctesting.NewCoordinator(t, 2)
			chainA = coordinator.GetChain(ibctesting.GetChainID(1))
			chainB = coordinator.GetChain(ibctesting.GetChainID(2))

			statedb := state.NewTestStateDB(t)
			statedb.Finalise(true)

			path = ibctesting.NewPath(chainA, chainB)
			path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
			path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED

			if test.BeforeHook != nil {
				test.BeforeHook(t, statedb)
				test.BeforeHook = nil
			}

			counterparty := channeltypes.NewCounterparty(ibctesting.MockPort, ibctesting.FirstChannelID)
			channel := channeltypes.NewChannel(channeltypes.INIT, channeltypes.ORDERED, counterparty, []string{path.EndpointB.ConnectionID}, path.EndpointA.ChannelConfig.Version)
			channelByte, err := marshaler.Marshal(&channel)
			require.NoError(t, err)
			input, err := PackChanOpenInit(ChanOpenInitInput{
				PortID:  path.EndpointA.ChannelConfig.PortID,
				Channel: channelByte,
			})
			require.NoError(t, err)

			connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
			SetConnection(statedb, path.EndpointA.ConnectionID, &connection)

			cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
			cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)

				SetClientState(statedb, connection.GetClientID(), clientState)
				SetConsensusState(statedb, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
			}

			test.Input = input
			test.Caller = common.Address{1}
			test.SuppliedGas = ChanOpenInitGasCost
			test.ReadOnly = false
			test.Run(t, Module, statedb)
		})
	}
}

func TestChanOpenTry(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	var (
		path       *ibctesting.Path
		heightDiff uint64
	)

	res, err := PackChanOpenTryOutput("channel-0")
	require.NoError(t, err)

	tests := map[string]testutils.PrecompileTest{
		"success": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointA.ChanOpenInit())
				chainB.CreatePortCapability(chainB.GetSimApp().ScopedIBCMockKeeper, ibctesting.MockPort)
			},
			ExpectedRes: res,
		},
		"connection is not OPEN": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupClients(path)
				// pass capability check
				chainB.CreatePortCapability(chainB.GetSimApp().ScopedIBCMockKeeper, ibctesting.MockPort)
				require.NoError(t, path.EndpointB.ConnOpenInit())
			},
			ExpectedErr: "connection state is not OPEN",
		},
		"consensus state not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointA.ChanOpenInit())
				chainB.CreatePortCapability(chainB.GetSimApp().ScopedIBCMockKeeper, ibctesting.MockPort)
				heightDiff = 3 // consensus state doesn't exist at this height
			},
			ExpectedErr: "client state height < proof height ",
		},
		"channel verification failed": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// not creating a channel on chainA will result in an invalid proof of existence
				coordinator.SetupConnections(path)
			},
			ExpectedErr: "chained membership proof contains nonexistence proof at index",
		},
		"connection version not negotiated": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointA.ChanOpenInit())

				// modify connB versions
				conn := path.EndpointB.GetConnection()

				version := connectiontypes.NewVersion("2", []string{"ORDER_ORDERED", "ORDER_UNORDERED"})
				conn.Versions = append(conn.Versions, version)

				chainB.App.GetIBCKeeper().ConnectionKeeper.SetConnection(
					chainB.GetContext(),
					path.EndpointB.ConnectionID, conn,
				)
				chainB.CreatePortCapability(chainB.GetSimApp().ScopedIBCMockKeeper, ibctesting.MockPort)
			},
			ExpectedErr: "single version must be negotiated on connection before opening channel",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			coordinator = ibctesting.NewCoordinator(t, 2)
			chainA = coordinator.GetChain(ibctesting.GetChainID(1))
			chainB = coordinator.GetChain(ibctesting.GetChainID(2))

			statedb := state.NewTestStateDB(t)
			statedb.Finalise(true)

			heightDiff = 0 // must be explicitly changed in malleate
			path = ibctesting.NewPath(chainA, chainB)

			if test.BeforeHook != nil {
				test.BeforeHook(t, statedb)
				test.BeforeHook = nil
			}

			if path.EndpointB.ClientID != "" {
				// ensure client is up to date
				require.NoError(t, path.EndpointB.UpdateClient())
			}

			counterparty := channeltypes.NewCounterparty(path.EndpointB.ChannelConfig.PortID, ibctesting.FirstChannelID)
			channel := channeltypes.NewChannel(channeltypes.INIT, channeltypes.ORDERED, counterparty, []string{path.EndpointB.ConnectionID}, path.EndpointA.ChannelConfig.Version)
			channelByte, err := marshaler.Marshal(&channel)
			require.NoError(t, err)

			channelKey := hosttypes.ChannelKey(counterparty.PortId, counterparty.ChannelId)
			proof, proofHeight := chainA.QueryProof(channelKey)

			height := malleateHeight(proofHeight, heightDiff)
			consensusHeightByte, err := marshaler.Marshal(&height)
			require.NoError(t, err)

			input, err := PackChanOpenTry(ChanOpenTryInput{
				PortID:              path.EndpointB.ChannelConfig.PortID,
				Channel:             channelByte,
				CounterpartyVersion: path.EndpointA.ChannelConfig.Version,
				ProofInit:           proof,
				ProofHeight:         consensusHeightByte,
			})
			require.NoError(t, err)

			connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
			SetConnection(statedb, path.EndpointB.ConnectionID, &connection)

			cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
			cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)
			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				SetClientState(statedb, connection.GetClientID(), clientState)
				SetConsensusState(statedb, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
			}

			test.Input = input
			test.Caller = common.Address{1}
			test.SuppliedGas = ChanOpenInitGasCost
			test.ReadOnly = false
			test.Run(t, Module, statedb)
		})
	}
}

func TestChanOpenAck(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	var (
		path                  *ibctesting.Path
		counterpartyChannelID string
		heightDiff            uint64
	)

	tests := map[string]testutils.PrecompileTest{
		"success": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointA.ChanOpenInit())
				require.NoError(t, path.EndpointB.ChanOpenTry())
			},
			ExpectedRes: make([]byte, 0),
		},
		"success with empty stored counterparty channel ID": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointA.ChanOpenInit())
				require.NoError(t, path.EndpointB.ChanOpenTry())

				// set the channel's counterparty channel identifier to empty string
				channel := path.EndpointA.GetChannel()
				channel.Counterparty.ChannelId = ""

				// use a different channel identifier
				counterpartyChannelID = path.EndpointB.ChannelID

				chainA.App.GetIBCKeeper().ChannelKeeper.SetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, channel)
			},
			ExpectedRes: make([]byte, 0),
		},
		"channel state is not INIT": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// create fully open channels on both chains
				coordinator.Setup(path)
			},
			ExpectedErr: "channel state should be INIT",
		},
		"connection not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointA.ChanOpenInit())
				require.NoError(t, path.EndpointB.ChanOpenTry())

				// set the channel's connection hops to wrong connection ID
				channel := path.EndpointA.GetChannel()
				channel.ConnectionHops[0] = doesnotexist
				chainA.App.GetIBCKeeper().ChannelKeeper.SetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, channel)
			},
			ExpectedErr: "can't read connection",
		},
		"connection is not OPEN": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
				require.NoError(t, path.EndpointA.ChanOpenInit())
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointA.ChanOpenInit())
				chainA.CreateChannelCapability(chainA.GetSimApp().ScopedIBCMockKeeper, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			},
			ExpectedErr: "channel handshake open ack failed",
		},
		"consensus state not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()

				require.NoError(t, path.EndpointA.ChanOpenInit())
				require.NoError(t, path.EndpointB.ChanOpenTry())

				heightDiff = 3 // consensus state doesn't exist at this height
			},
			ExpectedErr: "channel handshake open ack failed",
		},
		"channel verification failed": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// chainB is INIT, chainA in TRYOPEN
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()

				require.NoError(t, path.EndpointB.ChanOpenInit())
				require.NoError(t, path.EndpointA.ChanOpenTry())
			},
			ExpectedErr: "channel state should be INIT",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			coordinator = ibctesting.NewCoordinator(t, 2)
			chainA = coordinator.GetChain(ibctesting.GetChainID(1))
			chainB = coordinator.GetChain(ibctesting.GetChainID(2))

			statedb := state.NewTestStateDB(t)
			statedb.Finalise(true)

			counterpartyChannelID = "" // must be explicitly changed in malleate
			heightDiff = 0             // must be explicitly changed
			path = ibctesting.NewPath(chainA, chainB)

			if test.BeforeHook != nil {
				test.BeforeHook(t, statedb)
				test.BeforeHook = nil
			}

			if counterpartyChannelID == "" {
				counterpartyChannelID = ibctesting.FirstChannelID
			}

			if path.EndpointA.ClientID != "" {
				require.NoError(t, path.EndpointA.UpdateClient())
			}

			proof, proofHeight := chainB.QueryProof(hosttypes.ChannelKey(path.EndpointB.ChannelConfig.PortID, ibctesting.FirstChannelID))

			height := malleateHeight(proofHeight, heightDiff)
			consensusHeightByte, err := marshaler.Marshal(&height)
			require.NoError(t, err)

			input, err := PackChannelOpenAck(ChannelOpenAckInput{
				PortID:                path.EndpointA.ChannelConfig.PortID,
				ChannelID:             path.EndpointA.ChannelID,
				CounterpartyChannelID: counterpartyChannelID,
				CounterpartyVersion:   path.EndpointB.ChannelConfig.Version,
				ProofTry:              proof,
				ProofHeight:           consensusHeightByte,
			})
			require.NoError(t, err)

			SetCapability(statedb, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)

			channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			SetChannel(statedb, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)

			connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
			SetConnection(statedb, path.EndpointA.ConnectionID, &connection)

			cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
			cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				SetClientState(statedb, connection.GetClientID(), clientState)
				SetConsensusState(statedb, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
			}

			test.Input = input
			test.Caller = common.Address{1}
			test.SuppliedGas = ChannelOpenAckGasCost
			test.ReadOnly = false
			test.Run(t, Module, statedb)
		})
	}
}

func TestChanOpenConfirm(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	var (
		path       *ibctesting.Path
		heightDiff uint64
	)

	tests := map[string]testutils.PrecompileTest{
		"success": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointA.ChanOpenInit())
				require.NoError(t, path.EndpointB.ChanOpenTry())
				require.NoError(t, path.EndpointA.ChanOpenAck())
			},
			ExpectedRes: make([]byte, 0),
		},
		"channel state is not TRYOPEN": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// create fully open channels on both cahins
				coordinator.Setup(path)
			},
			ExpectedErr: "channel state is not TRYOPEN ",
		},
		"connection not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointA.ChanOpenInit())
				require.NoError(t, path.EndpointB.ChanOpenTry())
				require.NoError(t, path.EndpointA.ChanOpenAck())

				// set the channel's connection hops to wrong connection ID
				channel := path.EndpointB.GetChannel()
				channel.ConnectionHops[0] = doesnotexist
				chainB.App.GetIBCKeeper().ChannelKeeper.SetChannel(chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, channel)
			},
			ExpectedErr: "can't read connection: empty precompile state",
		},
		"connection is not OPEN": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupClients(path)
				require.NoError(t, path.EndpointB.ConnOpenInit())
				chainB.CreateChannelCapability(chainB.GetSimApp().ScopedIBCMockKeeper, path.EndpointB.ChannelConfig.PortID, ibctesting.FirstChannelID)
			},
			ExpectedErr: "can't read channel",
		},
		"consensus state not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointA.ChanOpenInit())
				require.NoError(t, path.EndpointB.ChanOpenTry())
				require.NoError(t, path.EndpointA.ChanOpenAck())

				heightDiff = 3
			},
			ExpectedErr: "channel handshake open ack failed",
		},
		"channel verification failed": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// chainA is INIT, chainB in TRYOPEN
				coordinator.SetupConnections(path)
				path.SetChannelOrdered()

				require.NoError(t, path.EndpointA.ChanOpenInit())
				require.NoError(t, path.EndpointB.ChanOpenTry())
			},
			ExpectedErr: "channel handshake open ack failed",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			coordinator = ibctesting.NewCoordinator(t, 2)
			chainA = coordinator.GetChain(ibctesting.GetChainID(1))
			chainB = coordinator.GetChain(ibctesting.GetChainID(2))

			statedb := state.NewTestStateDB(t)
			statedb.Finalise(true)

			heightDiff = 0 // must be explicitly changed
			path = ibctesting.NewPath(chainA, chainB)

			if test.BeforeHook != nil {
				test.BeforeHook(t, statedb)
				test.BeforeHook = nil
			}

			if path.EndpointB.ClientID != "" {
				require.NoError(t, path.EndpointB.UpdateClient())
			}

			proof, proofHeight := chainA.QueryProof(hosttypes.ChannelKey(path.EndpointA.ChannelConfig.PortID, ibctesting.FirstChannelID))

			height := malleateHeight(proofHeight, heightDiff)
			consensusHeightByte, err := marshaler.Marshal(&height)
			require.NoError(t, err)

			input, err := PackChannelOpenConfirm(ChannelOpenConfirmInput{
				PortID:      path.EndpointB.ChannelConfig.PortID,
				ChannelID:   ibctesting.FirstChannelID,
				ProofAck:    proof,
				ProofHeight: consensusHeightByte,
			})
			require.NoError(t, err)

			SetCapability(statedb, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)

			channel, _ := chainB.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			SetChannel(statedb, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, &channel)

			connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
			SetConnection(statedb, path.EndpointB.ConnectionID, &connection)

			cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
			cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)
			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				SetClientState(statedb, connection.GetClientID(), clientState)
				SetConsensusState(statedb, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
			}

			test.Input = input
			test.Caller = common.Address{1}
			test.SuppliedGas = ChannelOpenConfirmGasCost
			test.ReadOnly = false
			test.Run(t, Module, statedb)
		})
	}
}

func TestChanCloseInit(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	var path *ibctesting.Path

	tests := map[string]testutils.PrecompileTest{
		"success": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.Setup(path)
			},
			ExpectedRes: make([]byte, 0),
		},
		"channel doesn't exist": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// any non-nil values work for connections
				path.EndpointA.ConnectionID = ibctesting.FirstConnectionID
				path.EndpointB.ConnectionID = ibctesting.FirstConnectionID

				path.EndpointA.ChannelID = ibctesting.FirstChannelID
				path.EndpointB.ChannelID = ibctesting.FirstChannelID

				// ensure channel capability check passes
				chainA.CreateChannelCapability(chainA.GetSimApp().ScopedIBCMockKeeper, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			},
			ExpectedErr: "length channel.ConnectionHops == 0",
		},
		"channel state is CLOSED": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.Setup(path)

				require.NoError(t, path.EndpointA.SetChannelState(channeltypes.CLOSED))
			},
			ExpectedErr: "channel is already CLOSED",
		},
		"connection not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.Setup(path)

				// set the channel's connection hops to wrong connection ID
				channel := path.EndpointA.GetChannel()
				channel.ConnectionHops[0] = doesnotexist
				chainA.App.GetIBCKeeper().ChannelKeeper.SetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, channel)

			},
			ExpectedErr: "can't read connection",
		},
		"connection is not OPEN": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.SetupClients(path)

				require.NoError(t, path.EndpointA.ConnOpenInit())

				// create channel in init
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// ensure channel capability check passes
				chainA.CreateChannelCapability(chainA.GetSimApp().ScopedIBCMockKeeper, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			},
			ExpectedErr: "can't read channel",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			coordinator = ibctesting.NewCoordinator(t, 2)
			chainA = coordinator.GetChain(ibctesting.GetChainID(1))
			chainB = coordinator.GetChain(ibctesting.GetChainID(2))

			statedb := state.NewTestStateDB(t)
			statedb.Finalise(true)

			path = ibctesting.NewPath(chainA, chainB)

			if test.BeforeHook != nil {
				test.BeforeHook(t, statedb)
				test.BeforeHook = nil
			}

			input, err := PackChannelCloseInit(ChannelCloseInitInput{
				PortID:    path.EndpointA.ChannelConfig.PortID,
				ChannelID: ibctesting.FirstChannelID,
			})
			require.NoError(t, err)

			channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			SetChannel(statedb, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)

			connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
			SetConnection(statedb, path.EndpointA.ConnectionID, &connection)

			cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
			cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				SetClientState(statedb, connection.GetClientID(), clientState)
				SetConsensusState(statedb, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
			}

			test.Input = input
			test.Caller = common.Address{1}
			test.SuppliedGas = ChannelCloseInitGasCost
			test.ReadOnly = false
			test.Run(t, Module, statedb)
		})
	}
}

func TestChanCloseConfirm(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	var (
		path       *ibctesting.Path
		heightDiff uint64
	)

	tests := map[string]testutils.PrecompileTest{
		"success": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.Setup(path)
				require.NoError(t, path.EndpointA.SetChannelState(channeltypes.CLOSED))
			},
			ExpectedRes: make([]byte, 0),
		},
		"channel state is CLOSED": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.Setup(path)
				require.NoError(t, path.EndpointB.SetChannelState(channeltypes.CLOSED))
			},
			ExpectedErr: "channel is already CLOSED",
		},
		"connection not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.Setup(path)

				// set the channel's connection hops to wrong connection ID
				channel := path.EndpointB.GetChannel()
				channel.ConnectionHops[0] = doesnotexist
				chainB.App.GetIBCKeeper().ChannelKeeper.SetChannel(chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, channel)
			},
			ExpectedErr: "can't read connection",
		},
		"connection is not OPEN": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.Setup(path)

				require.NoError(t, path.EndpointB.ConnOpenInit())

				// create channel in init
				path.SetChannelOrdered()
				require.NoError(t, path.EndpointB.ChanOpenInit())

				// ensure channel capability check passes
				chainB.CreateChannelCapability(chainB.GetSimApp().ScopedIBCMockKeeper, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			},
			ExpectedErr: "can't read channel",
		},
		"consensus state not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.Setup(path)

				require.NoError(t, path.EndpointA.SetChannelState(channeltypes.CLOSED))
				heightDiff = 3
			},
			ExpectedErr: "client state height < proof height",
		},
		"channel verification failed": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// channel not closed
				coordinator.Setup(path)
			},
			ExpectedErr: "client state height < proof height",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			coordinator = ibctesting.NewCoordinator(t, 2)
			chainA = coordinator.GetChain(ibctesting.GetChainID(1))
			chainB = coordinator.GetChain(ibctesting.GetChainID(2))

			statedb := state.NewTestStateDB(t)
			statedb.Finalise(true)

			path = ibctesting.NewPath(chainA, chainB)

			if test.BeforeHook != nil {
				test.BeforeHook(t, statedb)
				test.BeforeHook = nil
			}

			channelKey := hosttypes.ChannelKey(path.EndpointA.ChannelConfig.PortID, ibctesting.FirstChannelID)
			proof, proofHeight := chainA.QueryProof(channelKey)

			height := malleateHeight(proofHeight, heightDiff)
			consensusHeightByte, _ := marshaler.Marshal(&height)

			input, err := PackChannelCloseConfirm(ChannelCloseConfirmInput{
				PortID:      path.EndpointB.ChannelConfig.PortID,
				ChannelID:   ibctesting.FirstChannelID,
				ProofInit:   proof,
				ProofHeight: consensusHeightByte,
			})
			require.NoError(t, err)

			channel, _ := chainB.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			SetChannel(statedb, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, &channel)

			connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
			SetConnection(statedb, path.EndpointB.ConnectionID, &connection)

			cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
			cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)
			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				SetClientState(statedb, connection.GetClientID(), clientState)
				SetConsensusState(statedb, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
			}

			test.Input = input
			test.Caller = common.Address{1}
			test.SuppliedGas = ChannelCloseConfirmGasCost
			test.ReadOnly = false
			test.Run(t, Module, statedb)
		})
	}
}
