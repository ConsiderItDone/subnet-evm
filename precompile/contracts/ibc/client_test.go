package ibc

import (
	"fmt"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	ibctestingmock "github.com/cosmos/ibc-go/v7/testing/mock"
)

const (
	testChainID          = "gaiahub-0"
	testChainIDRevision1 = "gaiahub-1"

	testClientID1 = "07-tendermint-0"
	testClientID2 = "07-tendermint-1"
	testClientID3 = "07-tendermint-2"

	trustingPeriod time.Duration = time.Hour * 24 * 7 * 2
	ubdPeriod      time.Duration = time.Hour * 24 * 7 * 3
	maxClockDrift  time.Duration = time.Second * 10
)

var (
	privVal          = ibctestingmock.NewPV()
	testClientHeight = clienttypes.NewHeight(0, 5)

	validator  *tmtypes.Validator
	validators *tmtypes.ValidatorSet
)

func init() {
	pubKey, _ := privVal.GetPubKey()
	validator = tmtypes.NewValidator(pubKey, 1)
	validators = tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
}

func TestCreateClient(t *testing.T) {
	createClientOuput, err := PackCreateClientOutput(testClientID1)
	require.NoError(t, err)

	tests := map[string]testutils.PrecompileTest{
		"success: 07-tendermint client type supported": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				clientState, err := ibctm.
					NewClientState(
						testChainID,
						ibctm.DefaultTrustLevel,
						trustingPeriod,
						ubdPeriod,
						maxClockDrift,
						testClientHeight,
						commitmenttypes.GetSDKSpecs(),
						ibctesting.UpgradePath,
					).
					Marshal()
				require.NoError(t, err)

				consensusState, err := ibctm.
					NewConsensusState(
						time.Now(),
						commitmenttypes.NewMerkleRoot([]byte("hash")),
						validators.Hash(),
					).
					Marshal()
				require.NoError(t, err)

				input, err := PackCreateClient(CreateClientInput{
					ClientType:     exported.Tendermint,
					ClientState:    clientState,
					ConsensusState: consensusState,
				})
				require.NoError(t, err)
				return input
			},
			SuppliedGas: CreateClientGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
			ExpectedRes: createClientOuput,
		},
	}

	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}

func TestUpdateClient(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))
	//now := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	past := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	var (
		path         *ibctesting.Path
		updateHeader *ibctm.Header
	)

	// Must create header creation functions since suite.header gets recreated on each test case
	createFutureUpdateFn := func(trustedHeight clienttypes.Height) *ibctm.Header {
		header, err := chainA.ConstructUpdateTMClientHeaderWithTrustedHeight(path.EndpointB.Chain, path.EndpointA.ClientID, trustedHeight)
		require.NoError(t, err)
		return header
	}
	createPastUpdateFn := func(fillHeight, trustedHeight clienttypes.Height) *ibctm.Header {
		consState, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientConsensusState(chainA.GetContext(), path.EndpointA.ClientID, trustedHeight)
		//require.True(t, found)

		return chainB.CreateTMClientHeader(chainB.ChainID, int64(fillHeight.RevisionHeight), trustedHeight, consState.(*ibctm.ConsensusState).Timestamp.Add(time.Second*5),
			chainB.Vals, chainB.Vals, chainB.Vals, chainB.Signers)
	}
	inputFn := func(t testing.TB) []byte {
		clientMessage, err := updateHeader.Marshal()
		require.NoError(t, err)
		input, err := PackUpdateClient(UpdateClientInput{
			ClientID:      path.EndpointA.ClientID,
			ClientMessage: clientMessage,
		})
		require.NoError(t, err)
		return input
	}

	tests := map[string]testutils.PrecompileTest{
		"valid past update": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				clientState := path.EndpointA.GetClientState()
				trustedHeight := clientState.GetLatestHeight().(clienttypes.Height)

				currHeight := chainB.CurrentHeader.Height
				fillHeight := clienttypes.NewHeight(clientState.GetLatestHeight().GetRevisionNumber(), uint64(currHeight))

				// commit a couple blocks to allow client to fill in gaps
				coordinator.CommitBlock(chainB) // this height is not filled in yet
				coordinator.CommitBlock(chainB) // this height is filled in by the update below

				require.NoError(t, path.EndpointA.UpdateClient())
				require.NoError(t, setClientState(
					state,
					path.EndpointA.ClientID,
					clientState.(*ibctm.ClientState),
				))

				// store previous consensus state
				prevConsState := &ibctm.ConsensusState{
					Timestamp:          time.Now(),
					NextValidatorsHash: chainB.Vals.Hash(),
				}
				require.NoError(t, setConsensusState(
					state,
					path.EndpointA.ClientID,
					clientState.GetLatestHeight(),
					prevConsState,
				))

				// ensure fill height not set
				chainA.App.GetIBCKeeper().ClientKeeper.GetClientConsensusState(chainA.GetContext(), path.EndpointA.ClientID, fillHeight)
				//require.True(t, found)

				// updateHeader will fill in consensus state between prevConsState and suite.consState
				// clientState should not be updated
				updateHeader = createPastUpdateFn(fillHeight, trustedHeight)
			},
			ExpectedRes: make([]byte, 0),
		},
		"misbehaviour detection: conflicting header": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				clientID := path.EndpointA.ClientID

				height1 := clienttypes.NewHeight(1, 1)
				// store previous consensus state
				prevConsState := &ibctm.ConsensusState{
					Timestamp:          past,
					NextValidatorsHash: chainB.Vals.Hash(),
				}
				chainA.App.GetIBCKeeper().ClientKeeper.SetClientConsensusState(chainA.GetContext(), clientID, height1, prevConsState)

				height5 := clienttypes.NewHeight(1, 5)
				// store next consensus state to check that trustedHeight does not need to be hightest consensus state before header height
				nextConsState := &ibctm.ConsensusState{
					Timestamp:          past.Add(time.Minute),
					NextValidatorsHash: chainB.Vals.Hash(),
				}
				chainA.App.GetIBCKeeper().ClientKeeper.SetClientConsensusState(chainA.GetContext(), clientID, height5, nextConsState)

				height3 := clienttypes.NewHeight(1, 3)
				// updateHeader will fill in consensus state between prevConsState and suite.consState
				// clientState should not be updated
				updateHeader = createPastUpdateFn(height3, height1)
				// set conflicting consensus state in store to create misbehaviour scenario
				conflictConsState := updateHeader.ConsensusState()
				conflictConsState.Root = commitmenttypes.NewMerkleRoot([]byte("conflicting apphash"))
				chainA.App.GetIBCKeeper().ClientKeeper.SetClientConsensusState(chainA.GetContext(), clientID, updateHeader.GetHeight(), conflictConsState)
			},
			ExpectedErr: "can't get client state",
		},
		"client state not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				updateHeader = createFutureUpdateFn(path.EndpointA.GetClientState().GetLatestHeight().(clienttypes.Height))
				path.EndpointA.ClientID = ibctesting.InvalidID
			},
			ExpectedErr: "can't get client state",
		},
		"consensus state not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				clientState := path.EndpointA.GetClientState()
				tmClient, ok := clientState.(*ibctm.ClientState)
				require.True(t, ok)

				tmClient.LatestHeight = tmClient.LatestHeight.Increment().(clienttypes.Height)
				require.NoError(t, setClientState(
					state,
					path.EndpointA.ClientID,
					clientState.(*ibctm.ClientState),
				))
				updateHeader = createFutureUpdateFn(clientState.GetLatestHeight().(clienttypes.Height))
			},
			ExpectedErr: "can't get consensus state",
		},
		"client is not active": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				clientState := path.EndpointA.GetClientState().(*ibctm.ClientState)
				clientState.FrozenHeight = clienttypes.NewHeight(1, 1)
				require.NoError(t, setClientState(
					state,
					path.EndpointA.ClientID,
					clientState,
				))
				updateHeader = createFutureUpdateFn(clientState.GetLatestHeight().(clienttypes.Height))
			},
			ExpectedErr: "can't get consensus state",
		},
		"invalid header": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				updateHeader = createFutureUpdateFn(path.EndpointA.GetClientState().GetLatestHeight().(clienttypes.Height))
				updateHeader.TrustedHeight = updateHeader.TrustedHeight.Increment().(clienttypes.Height)
			},
			ExpectedErr: "can't get client state",
		},
	}

	statedb := state.NewTestStateDB(t)
	statedb.Finalise(true)

	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if name != "misbehaviour detection: conflicting header" {
				statedb = state.NewTestStateDB(t)
				statedb.Finalise(true)
			}

			path = ibctesting.NewPath(chainA, chainB)
			coordinator.SetupClients(path)

			test.Caller = common.Address{1}
			test.SuppliedGas = UpgradeClientGasCost
			test.ReadOnly = false
			test.InputFn = inputFn
			test.Run(t, Module, statedb)
		})
	}
}

func TestUpgradeClient(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))

	var (
		path                                        *ibctesting.Path
		upgradedClient                              exported.ClientState
		lastHeight                                  exported.Height
		proofUpgradedClient, proofUpgradedConsState []byte
		upgradedClientBz, upgradedConsStateBz       []byte
		//clientStatePath, consensusStatePath         string
	)

	inputFn := func(
		clientId string,
		upgradedClient exported.ClientState,
		upgradedConsState *ibctm.ConsensusState,
	) func(t testing.TB) []byte {
		return func(t testing.TB) []byte {
			upgradedClientByte, err := upgradedClient.(*ibctm.ClientState).Marshal()
			require.NoError(t, err)

			upgradedConsStateByte, err := upgradedConsState.Marshal()
			require.NoError(t, err)

			output, err := PackUpgradeClient(UpgradeClientInput{
				ClientID:              clientId,
				UpgradePath:           []byte("path"),
				UpgradedClien:         upgradedClientByte,
				UpgradedConsState:     upgradedConsStateByte,
				ProofUpgradeClient:    proofUpgradedClient,
				ProofUpgradeConsState: proofUpgradedConsState,
			})
			require.NoError(t, err)

			return output
		}
	}

	tests := map[string]testutils.PrecompileTest{
		"successful upgrade": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				lastHeight = clienttypes.NewHeight(1, uint64(chainB.GetContext().BlockHeight()+1))

				require.NoError(t, chainB.GetSimApp().UpgradeKeeper.SetUpgradedClient(chainB.GetContext(), int64(lastHeight.GetRevisionHeight()), upgradedClientBz))
				require.NoError(t, chainB.GetSimApp().UpgradeKeeper.SetUpgradedConsensusState(chainB.GetContext(), int64(lastHeight.GetRevisionHeight()), upgradedConsStateBz))

				// commit upgrade store changes and update clients
				coordinator.CommitBlock(chainB)
				require.NoError(t, path.EndpointA.UpdateClient())

				cs, found := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				require.True(t, found)

				proofUpgradedClient, _ = chainB.QueryUpgradeProof(upgradetypes.UpgradedClientKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
				proofUpgradedConsState, _ = chainB.QueryUpgradeProof(upgradetypes.UpgradedConsStateKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
			},
			ExpectedRes: make([]byte, 0),
		},
		"client state not found": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// last Height is at next block
				lastHeight = clienttypes.NewHeight(1, uint64(chainB.GetContext().BlockHeight()+1))

				// zero custom fields and store in upgrade store
				require.NoError(t, chainB.GetSimApp().UpgradeKeeper.SetUpgradedClient(chainB.GetContext(), int64(lastHeight.GetRevisionHeight()), upgradedClientBz))
				require.NoError(t, chainB.GetSimApp().UpgradeKeeper.SetUpgradedConsensusState(chainB.GetContext(), int64(lastHeight.GetRevisionHeight()), upgradedConsStateBz))

				// commit upgrade store changes and update clients
				coordinator.CommitBlock(chainB)
				require.NoError(t, path.EndpointA.UpdateClient())

				cs, found := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				require.True(t, found)

				proofUpgradedClient, _ = chainB.QueryUpgradeProof(upgradetypes.UpgradedClientKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
				proofUpgradedConsState, _ = chainB.QueryUpgradeProof(upgradetypes.UpgradedConsStateKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())

				path.EndpointA.ClientID = "wrongclientid"
			},
			ExpectedErr: "can't get client state",
		},
		"tendermint client VerifyUpgrade fails": {
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				// last Height is at next block
				lastHeight = clienttypes.NewHeight(1, uint64(chainB.GetContext().BlockHeight()+1))

				// zero custom fields and store in upgrade store
				require.NoError(t, chainB.GetSimApp().UpgradeKeeper.SetUpgradedClient(chainB.GetContext(), int64(lastHeight.GetRevisionHeight()), upgradedClientBz))
				require.NoError(t, chainB.GetSimApp().UpgradeKeeper.SetUpgradedConsensusState(chainB.GetContext(), int64(lastHeight.GetRevisionHeight()), upgradedConsStateBz))

				// change upgradedClient client-specified parameters
				tmClient := upgradedClient.(*ibctm.ClientState)
				tmClient.ChainId = "wrongchainID"
				upgradedClient = tmClient

				coordinator.CommitBlock(chainB)
				require.NoError(t, path.EndpointA.UpdateClient())

				cs, found := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				require.True(t, found)

				proofUpgradedClient, _ = chainB.QueryUpgradeProof(upgradetypes.UpgradedClientKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
				proofUpgradedConsState, _ = chainB.QueryUpgradeProof(upgradetypes.UpgradedConsStateKey(int64(lastHeight.GetRevisionHeight())), cs.GetLatestHeight().GetRevisionHeight())
			},
			ExpectedErr: "client state proof failed.",
		},
	}

	// Run tests.
	for name, test := range tests {
		statedb := state.NewTestStateDB(t)
		statedb.Finalise(true)

		path = ibctesting.NewPath(chainA, chainB)
		coordinator.SetupClients(path)

		clientState := path.EndpointA.GetClientState().(*ibctm.ClientState)
		revisionNumber := clienttypes.ParseChainID(clientState.ChainId)

		newChainID, err := clienttypes.SetRevisionNumber(clientState.ChainId, revisionNumber+1)
		require.NoError(t, err)

		upgradedClient = ibctm.NewClientState(newChainID, ibctm.DefaultTrustLevel, trustingPeriod, ubdPeriod+trustingPeriod, maxClockDrift, clienttypes.NewHeight(revisionNumber+1, clientState.GetLatestHeight().GetRevisionHeight()+1), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath)
		upgradedClient = upgradedClient.ZeroCustomFields()
		upgradedClientBz, err = clienttypes.MarshalClientState(chainA.App.AppCodec(), upgradedClient)
		require.NoError(t, err)

		upgradedConsState := &ibctm.ConsensusState{
			NextValidatorsHash: []byte("nextValsHash"),
		}
		upgradedConsStateBz, err = clienttypes.MarshalConsensusState(chainA.App.AppCodec(), upgradedConsState)
		require.NoError(t, err)

		t.Run(name, func(t *testing.T) {
			originalBeforeHook := test.BeforeHook
			test.BeforeHook = func(t testing.TB, state contract.StateDB) {
				originalBeforeHook(t, state)
				interfaceRegistry := cosmostypes.NewInterfaceRegistry()
				marshaler := codec.NewProtoCodec(interfaceRegistry)

				std.RegisterInterfaces(interfaceRegistry)
				ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)

				if cs != nil {
					clientState = cs.(*ibctm.ClientState)
					setClientState(statedb, clientState.ChainId, clientState)

					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					setConsensusState(statedb, clientState.ChainId, clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
			}
			test.Caller = common.Address{1}
			test.SuppliedGas = UpgradeClientGasCost
			test.ReadOnly = false
			test.InputFn = inputFn(clientState.ChainId, upgradedClient, upgradedConsState)
			test.Run(t, Module, statedb)
		})
	}
}
