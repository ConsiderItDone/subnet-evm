package ibc

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/utils"

	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestConnOpenInit(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	var (
		path         *ibctesting.Path
		version      *connectiontypes.Version
		delayPeriod  uint64
		emptyConnBID bool
	)

	inputFn := func(t testing.TB) []byte {
		if emptyConnBID {
			path.EndpointB.ConnectionID = ""
		}
		counterparty := connectiontypes.NewCounterparty(path.EndpointB.ClientID, path.EndpointB.ConnectionID, chainB.GetPrefix())
		counterpartybyte, _ := marshaler.MarshalInterface(&counterparty)
		var versionbyte []byte
		if version == nil {
			versionbyte = []byte("")
		} else {
			versionbyte, _ = version.Marshal()
		}
		input, err := PackConnOpenInit(ConnOpenInitInput{
			ClientID:     path.EndpointA.ClientID,
			Counterparty: counterpartybyte,
			Version:      versionbyte,
			DelayPeriod:  uint32(delayPeriod),
		})
		require.NoError(t, err)
		return input
	}

	res, _ := PackConnOpenInitOutput("connection-0")

	tests := map[string]testutils.PrecompileTest{
		"success": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
			},
			ExpectedRes: res,
		},
		"success with empty counterparty identifier": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				emptyConnBID = true
			},
			ExpectedRes: res,
		},
		"success with non empty version": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				version = connectiontypes.ExportedVersionsToProto(connectiontypes.GetCompatibleVersions())[0]
			},
			ExpectedRes: res,
		},
		"success with non zero delayPeriod": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				delayPeriod = uint64(time.Hour.Nanoseconds())
			},
			ExpectedRes: res,
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
			emptyConnBID = false // must be explicitly changed
			version = nil        // must be explicitly

			path = ibctesting.NewPath(chainA, chainB)
			coordinator.SetupClients(path)

			clientState := path.EndpointA.GetClientState()
			require.NoError(t, SetClientState(statedb, path.EndpointA.ClientID, clientState.(*ibctm.ClientState)))

			consensusState := path.EndpointA.GetConsensusState(clientState.GetLatestHeight())

			consensusStateIbctm := consensusState.(*ibctm.ConsensusState)
			consensusStateIbctm.Timestamp = time.Now()
			require.NoError(t, SetConsensusState(statedb, path.EndpointA.ClientID, clientState.GetLatestHeight(), consensusStateIbctm))
			test.Config = NewConfig(utils.NewUint64(0))
			test.Caller = common.Address{1}
			test.SuppliedGas = ConnOpenInitGasCost
			test.ReadOnly = false
			test.InputFn = inputFn
			test.Run(t, Module, statedb)
		})
	}
}

func TestConnOpenTry(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	var (
		path               *ibctesting.Path
		delayPeriod        uint64
		versions           []exported.Version
		consensusHeight    clienttypes.Height
		counterpartyClient exported.ClientState
	)

	res, err := PackConnOpenTryOutput("connection-0")
	require.NoError(t, err)

	tests := map[string]testutils.PrecompileTest{
		"success": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())
				// retrieve client state of chainA to pass as counterpartyClient
				counterpartyClient = chainA.GetClientState(path.EndpointA.ClientID)
			},
			ExpectedRes: res,
		},
		//"success with delay period": {
		//	BeforeHook: func(t testing.TB, state contract.StateDB) {
		//		require.NoError(t, path.EndpointA.ConnOpenInit())
		//
		//		delayPeriod = uint64(time.Hour.Nanoseconds())
		//
		//		// set delay period on counterparty to non-zero value
		//		conn := path.EndpointA.GetConnection()
		//		conn.DelayPeriod = delayPeriod
		//		chainA.App.GetIBCKeeper().ConnectionKeeper.SetConnection(chainA.GetContext(), path.EndpointA.ConnectionID, conn)
		//
		//		// commit in order for proof to return correct value
		//		coordinator.CommitBlock(chainA)
		//		require.NoError(t, path.EndpointB.UpdateClient())
		//
		//		// retrieve client state of chainA to pass as counterpartyClient
		//		counterpartyClient = chainA.GetClientState(path.EndpointA.ClientID)
		//	},
		//	ExpectedRes: res,
		//},
		"invalid counterparty client": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainA.GetClientState(path.EndpointA.ClientID)

				// Set an invalid client of chainA on chainB
				tmClient, ok := counterpartyClient.(*ibctm.ClientState)
				require.True(t, ok)
				tmClient.ChainId = "wrongchainid"

				chainA.App.GetIBCKeeper().ClientKeeper.SetClientState(chainA.GetContext(), path.EndpointA.ClientID, tmClient)
			},
			ExpectedErr: "error clientVerification: chained membership proof failed to verify membership of value",
		},
		"counterparty versions is empty": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainA to pass as counterpartyClient
				counterpartyClient = chainA.GetClientState(path.EndpointA.ClientID)

				versions = nil
			},
			ExpectedErr: "error PickVersion",
		},
		"counterparty versions don't have a match": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainA to pass as counterpartyClient
				counterpartyClient = chainA.GetClientState(path.EndpointA.ClientID)

				version := connectiontypes.NewVersion("0.0", nil)
				versions = []exported.Version{version}
			},
			ExpectedErr: "error PickVersion err: failed to find a matching counterparty version",
		},
		"connection state verification failed": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// chainA connection not created

				// retrieve client state of chainA to pass as counterpartyClient
				counterpartyClient = chainA.GetClientState(path.EndpointA.ClientID)
			},
			ExpectedErr: "error connectionVerification: chained membership proof contains nonexistence proof",
		},
		"client state verification failed": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainA to pass as counterpartyClient
				counterpartyClient = chainA.GetClientState(path.EndpointA.ClientID)

				// modify counterparty client without setting in store so it still passes validate but fails proof verification
				tmClient, ok := counterpartyClient.(*ibctm.ClientState)
				require.True(t, ok)
				tmClient.LatestHeight = tmClient.LatestHeight.Increment().(clienttypes.Height)
				counterpartyClient = tmClient
			},
			ExpectedErr: "error clientVerification: chained membership proof failed to verify membership",
		},
		//"consensus state verification failed": {
		//	BeforeHook: func(t testing.TB, state contract.StateDB) {
		//		// retrieve client state of chainA to pass as counterpartyClient
		//		counterpartyClient = chainA.GetClientState(path.EndpointA.ClientID)
		//
		//		// give chainA wrong consensus state for chainB
		//		consState, found := chainA.App.GetIBCKeeper().ClientKeeper.GetLatestClientConsensusState(chainA.GetContext(), path.EndpointA.ClientID)
		//		require.True(t, found)
		//
		//		tmConsState, ok := consState.(*ibctm.ConsensusState)
		//		require.True(t, ok)
		//
		//		tmConsState.Timestamp = time.Now()
		//		chainA.App.GetIBCKeeper().ClientKeeper.SetClientConsensusState(chainA.GetContext(), path.EndpointA.ClientID, counterpartyClient.GetLatestHeight(), tmConsState)
		//
		//		require.NoError(t, path.EndpointA.ConnOpenInit())
		//	},
		//	ExpectedErr: "error clientVerification",
		//},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			coordinator = ibctesting.NewCoordinator(t, 2)
			chainA = coordinator.GetChain(ibctesting.GetChainID(1))
			chainB = coordinator.GetChain(ibctesting.GetChainID(2))

			statedb := state.NewTestStateDB(t)
			statedb.Finalise(true)

			consensusHeight = clienttypes.ZeroHeight()         // may be changed in malleate
			versions = connectiontypes.GetCompatibleVersions() // may be changed in malleate
			delayPeriod = 0                                    // may be changed in malleate

			path = ibctesting.NewPath(chainA, chainB)
			coordinator.SetupClients(path)

			if test.BeforeHook != nil {
				test.BeforeHook(t, statedb)
				test.BeforeHook = nil
			}

			counterparty := connectiontypes.NewCounterparty(path.EndpointA.ClientID, path.EndpointA.ConnectionID, chainA.GetPrefix())

			// ensure client is up to date to receive proof
			require.NoError(t, path.EndpointB.UpdateClient())

			connectionKey := host.ConnectionKey(path.EndpointA.ConnectionID)
			proofInit, proofHeight := chainA.QueryProof(connectionKey)
			//fmt.Printf("proofInit %#v\n", proofInit)
			//fmt.Printf("proofHeight %#v\n", proofHeight)
			if consensusHeight.IsZero() {
				// retrieve consensus state height to provide proof for
				consensusHeight = counterpartyClient.GetLatestHeight().(clienttypes.Height)
			}
			consensusKey := host.FullConsensusStateKey(path.EndpointA.ClientID, consensusHeight)
			proofConsensus, _ := chainA.QueryProof(consensusKey)
			//fmt.Printf("proofConsensus %#v\n", proofConsensus)

			// retrieve proof of counterparty clientstate on chainA
			clientKey := host.FullClientStateKey(path.EndpointA.ClientID)
			proofClient, _ := chainA.QueryProof(clientKey)
			//fmt.Printf("proofClient %#v\n", proofClient)

			counterpartyByte, _ := counterparty.Marshal()
			//fmt.Printf("counterparty %#v\n", counterpartyByte)

			clientStateByte, _ := clienttypes.MarshalClientState(marshaler, counterpartyClient)
			//fmt.Printf("clientState %#v\n", clientStateByte)

			versionsByte, _ := json.Marshal(connectiontypes.ExportedVersionsToProto(versions))
			//fmt.Printf("versions %#v\n", versionsByte)

			proofHeightByte, _ := marshaler.MarshalInterface(&proofHeight)
			//fmt.Printf("proofHeightByte %#v\n", proofHeightByte)

			consensusHeightByte, _ := marshaler.MarshalInterface(&consensusHeight)
			//fmt.Printf("consensusHeightByte %#v\n", consensusHeightByte)

			input, err := PackConnOpenTry(ConnOpenTryInput{
				Counterparty:         counterpartyByte,
				DelayPeriod:          uint32(delayPeriod),
				ClientID:             path.EndpointB.ClientID,
				ClientState:          clientStateByte,
				CounterpartyVersions: versionsByte,
				ProofInit:            proofInit,
				ProofClient:          proofClient,
				ProofConsensus:       proofConsensus,
				ProofHeight:          proofHeightByte,
				ConsensusHeight:      consensusHeightByte,
			})
			require.NoError(t, err)

			connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
			SetConnection(statedb, path.EndpointB.ConnectionID, &connection)

			cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
			cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)
			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				require.NoError(t, SetClientState(statedb, path.EndpointB.ClientID, clientState))

				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				consensusStateIbctm := consensusState.(*ibctm.ConsensusState)
				consensusStateIbctm.Timestamp = time.Now()

				require.NoError(t, SetConsensusState(statedb, path.EndpointB.ClientID, clientState.GetLatestHeight(), consensusStateIbctm))
			}

			test.Caller = common.Address{1}
			test.SuppliedGas = ConnOpenTryGasCost
			test.ReadOnly = false
			test.Input = input
			test.Run(t, Module, statedb)
		})
	}
}

func TestConnOpenAck(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))

	var (
		path               *ibctesting.Path
		consensusHeight    clienttypes.Height
		version            *connectiontypes.Version
		counterpartyClient exported.ClientState

		proofTry       []byte
		proofHeight    clienttypes.Height
		proofConsensus []byte
		proofClient    []byte
	)

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	tests := map[string]testutils.PrecompileTest{
		"success": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())
				require.NoError(t, path.EndpointB.ConnOpenTry())
				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)
			},
			ExpectedRes: make([]byte, 0),
		},
		"invalid counterparty client": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())
				require.NoError(t, path.EndpointB.ConnOpenTry())

				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)

				// Set an invalid client of chainA on chainB
				tmClient, ok := counterpartyClient.(*ibctm.ClientState)
				require.True(t, ok)
				tmClient.ChainId = "wrongchainid"

				chainB.App.GetIBCKeeper().ClientKeeper.SetClientState(chainB.GetContext(), path.EndpointB.ClientID, tmClient)
			},
			ExpectedErr: "chained membership proof failed",
		},
		"consensus height >= latest height": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)

				require.NoError(t, path.EndpointB.ConnOpenTry())
				consensusHeight = clienttypes.GetSelfHeight(chainA.GetContext())
			},
			ExpectedErr: "chained membership proof failed",
		},
		"connection not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)
			},
			ExpectedErr: "connection state is not INIT",
		},
		"invalid counterparty connection ID": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())
				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)
				require.NoError(t, path.EndpointB.ConnOpenTry())

				// modify connB to set counterparty connection identifier to wrong identifier
				connection, found := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				require.True(t, found)

				connection.Counterparty.ConnectionId = "badconnectionid"

				chainA.App.GetIBCKeeper().ConnectionKeeper.SetConnection(chainA.GetContext(), path.EndpointA.ConnectionID, connection)

				require.NoError(t, path.EndpointA.UpdateClient())
				require.NoError(t, path.EndpointB.UpdateClient())
			},
			ExpectedErr: "chained membership proof failed to verify membership",
		},
		"connection state is not INIT": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// connection state is already OPEN on chainA
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)

				require.NoError(t, path.EndpointB.ConnOpenTry())
				require.NoError(t, path.EndpointA.ConnOpenAck())
			},
			ExpectedErr: "connection state is not INIT",
		},
		"connection is in INIT but the proposed version is invalid": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// chainA is in INIT, chainB is in TRYOPEN
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)

				require.NoError(t, path.EndpointB.ConnOpenTry())

				version = connectiontypes.NewVersion("2.0", nil)
			},
			ExpectedErr: "the counterparty selected version",
		},
		"incompatible IBC versions": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)

				require.NoError(t, path.EndpointB.ConnOpenTry())

				// set version to a non-compatible version
				version = connectiontypes.NewVersion("2.0", nil)
			},
			ExpectedErr: "the counterparty selected version",
		},
		"empty version": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)

				require.NoError(t, path.EndpointB.ConnOpenTry())

				version = &connectiontypes.Version{}
			},
			ExpectedErr: "the counterparty selected version",
		},
		"feature set verification failed - unsupported feature": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)

				require.NoError(t, path.EndpointB.ConnOpenTry())

				version = connectiontypes.NewVersion(connectiontypes.DefaultIBCVersionIdentifier, []string{"ORDER_ORDERED", "ORDER_UNORDERED", "ORDER_DAG"})
			},
			ExpectedErr: "the counterparty selected version",
		},
		"self consensus state not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)

				require.NoError(t, path.EndpointB.ConnOpenTry())

				consensusHeight = clienttypes.NewHeight(0, 1)
			},
			ExpectedErr: "chained membership proof failed to verify membership of value",
		},
		"connection state verification failed": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())
				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)
			},
			ExpectedErr: "chained membership proof contains nonexistence proof",
		},
		"client state verification failed": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)

				// modify counterparty client without setting in store so it still passes validate but fails proof verification
				tmClient, ok := counterpartyClient.(*ibctm.ClientState)
				require.True(t, ok)
				tmClient.LatestHeight = tmClient.LatestHeight.Increment().(clienttypes.Height)

				require.NoError(t, path.EndpointB.ConnOpenTry())
			},
			ExpectedErr: "chained membership proof failed to verify membership",
		},
		"consensus state verification failed": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				require.NoError(t, path.EndpointA.ConnOpenInit())

				// retrieve client state of chainB to pass as counterpartyClient
				counterpartyClient = chainB.GetClientState(path.EndpointB.ClientID)

				// give chainB wrong consensus state for chainA
				consState, found := chainB.App.GetIBCKeeper().ClientKeeper.GetLatestClientConsensusState(chainB.GetContext(), path.EndpointB.ClientID)
				require.True(t, found)

				tmConsState, ok := consState.(*ibctm.ConsensusState)
				require.True(t, ok)

				tmConsState.Timestamp = tmConsState.Timestamp.Add(time.Second)
				chainB.App.GetIBCKeeper().ClientKeeper.SetClientConsensusState(chainB.GetContext(), path.EndpointB.ClientID, counterpartyClient.GetLatestHeight(), tmConsState)

				require.NoError(t, path.EndpointB.ConnOpenTry())
			},
			ExpectedErr: "chained membership proof failed to verify membership of value",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			statedb := state.NewTestStateDB(t)
			statedb.Finalise(true)

			coordinator := ibctesting.NewCoordinator(t, 2)
			chainA = coordinator.GetChain(ibctesting.GetChainID(1))
			chainB = coordinator.GetChain(ibctesting.GetChainID(2))

			version = connectiontypes.ExportedVersionsToProto(connectiontypes.GetCompatibleVersions())[0] // must be explicitly changed in malleate
			consensusHeight = clienttypes.ZeroHeight()
			path = ibctesting.NewPath(chainA, chainB)
			coordinator.SetupClients(path)

			orgBeforeHook := test.BeforeHook
			test.Caller = common.Address{1}
			test.SuppliedGas = ConnOpenAckGasCost
			test.ReadOnly = false
			test.BeforeHook = func(t testing.TB, state contract.StateDB) {
				orgBeforeHook(t, state)
				require.NoError(t, path.EndpointA.UpdateClient())

				connectionKey := host.ConnectionKey(path.EndpointB.ConnectionID)
				proofTry, proofHeight = chainB.QueryProof(connectionKey)

				if consensusHeight.IsZero() {
					// retrieve consensus state height to provide proof for
					clientState := chainB.GetClientState(path.EndpointB.ClientID)
					consensusHeight = clientState.GetLatestHeight().(clienttypes.Height)
				}
				consensusKey := host.FullConsensusStateKey(path.EndpointB.ClientID, consensusHeight)
				proofConsensus, _ = chainB.QueryProof(consensusKey)

				// retrieve proof of counterparty clientstate on chainA
				clientKey := host.FullClientStateKey(path.EndpointB.ClientID)
				proofClient, _ = chainB.QueryProof(clientKey)

				connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				SetConnection(state, path.EndpointA.ConnectionID, &connection)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), connection.GetClientID())
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), connection.GetClientID())

				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					SetClientState(state, connection.GetClientID(), clientState)

					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					consensusStateIbctm := consensusState.(*ibctm.ConsensusState)
					consensusStateIbctm.Timestamp = time.Now()

					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusStateIbctm)
				}
			}
			test.InputFn = func(t testing.TB) []byte {
				clientStateByte, _ := clienttypes.MarshalClientState(marshaler, counterpartyClient)
				versionByte, _ := marshaler.Marshal(version)
				proofHeightByte, _ := marshaler.Marshal(&proofHeight)
				consensusHeightByte, _ := marshaler.Marshal(&consensusHeight)
				input, err := PackConnOpenAck(ConnOpenAckInput{
					ConnectionID:             path.EndpointA.ConnectionID,
					ClientState:              clientStateByte,
					Version:                  versionByte,
					CounterpartyConnectionID: []byte(path.EndpointB.ConnectionID),
					ProofTry:                 proofTry,
					ProofClient:              proofClient,
					ProofConsensus:           proofConsensus,
					ProofHeight:              proofHeightByte,
					ConsensusHeight:          consensusHeightByte,
				})
				require.NoError(t, err)

				return input
			}
			test.Run(t, Module, statedb)
		})
	}
}

func TestConnOpenConfirm(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))

	var path *ibctesting.Path

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	tests := map[string]testutils.PrecompileTest{
		"success": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				path.EndpointA.ConnOpenInit()
				path.EndpointB.ConnOpenTry()
				path.EndpointA.ConnOpenAck()
			},
			ExpectedRes: make([]byte, 0),
		},
		"connection not found": {
			ExpectedErr: "connection state is not TRYOPEN",
		},
		"chain B's connection state is not TRYOPEN": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.CreateConnections(path)
			},
			ExpectedErr: "connection state is not TRYOPEN",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			statedb := state.NewTestStateDB(t)
			statedb.Finalise(true)

			coordinator = ibctesting.NewCoordinator(t, 2)
			chainA = coordinator.GetChain(ibctesting.GetChainID(1))
			chainB = coordinator.GetChain(ibctesting.GetChainID(2))
			path = ibctesting.NewPath(chainA, chainB)
			coordinator.SetupClients(path)

			if test.BeforeHook != nil {
				test.BeforeHook(t, statedb)
				test.BeforeHook = nil
			}

			require.NoError(t, path.EndpointB.UpdateClient())

			connectionKey := host.ConnectionKey(path.EndpointA.ConnectionID)
			proofAck, proofHeight := chainA.QueryProof(connectionKey)

			proofHeightByte, err := marshaler.Marshal(&proofHeight)
			require.NoError(t, err)

			input, err := PackConnOpenConfirm(ConnOpenConfirmInput{
				ConnectionID: path.EndpointB.ConnectionID,
				ProofAck:     proofAck,
				ProofHeight:  proofHeightByte,
			})
			require.NoError(t, err)

			connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
			SetConnection(statedb, path.EndpointB.ConnectionID, &connection)

			cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
			cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)

			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				SetClientState(statedb, connection.GetClientID(), clientState)

				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				consensusStateIbctm := consensusState.(*ibctm.ConsensusState)
				consensusStateIbctm.Timestamp = time.Now()
				SetConsensusState(statedb, connection.GetClientID(), clientState.GetLatestHeight(), consensusStateIbctm)
			}

			test.Caller = common.Address{1}
			test.SuppliedGas = ConnOpenConfirmGasCost
			test.ReadOnly = false
			test.Input = input
			test.Run(t, Module, statedb)
		})
	}
}
