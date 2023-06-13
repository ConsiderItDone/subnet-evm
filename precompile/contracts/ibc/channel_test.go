package ibc

import (
	"encoding/binary"
	"fmt"
	"math/big"

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

	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
)

const doesnotexist = "doesnotexist"

func (suite *KeeperTestSuite) TestChanOpenInit() {
	var path *ibctesting.Path

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	testCase := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{"success", func() {
			suite.coordinator.SetupConnections(path)
		}, true},
		{"connection doesn't exist", func() {
			// any non-empty values
			path.EndpointA.ConnectionID = "connection-0"
			path.EndpointB.ConnectionID = "connection-0"
		}, false},
	}

	for _, tc := range testCase {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			// run test for all types of ordering
			for _, order := range []channeltypes.Order{channeltypes.UNORDERED, channeltypes.ORDERED} {

				statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
				statedb.Finalise(true)
				vmctx := vm.BlockContext{
					CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
					Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
				}
				vmenv := vm.NewEVM(vmctx, vm.TxContext{}, statedb, params.TestChainConfig, vm.Config{ExtraEips: []int{2200}})

				suite.SetupTest() // reset
				path = ibctesting.NewPath(suite.chainA, suite.chainB)
				path.EndpointA.ChannelConfig.Order = order
				path.EndpointB.ChannelConfig.Order = order
				tc.malleate()

				counterparty := channeltypes.NewCounterparty(ibctesting.MockPort, ibctesting.FirstChannelID)

				var input []byte

				portIDbyte := []byte(path.EndpointA.ChannelConfig.PortID)
				portIDLen := make([]byte, 8)
				binary.BigEndian.PutUint64(portIDLen, uint64(len(portIDbyte)))

				input = append(input, portIDLen...)
				input = append(input, portIDbyte...)

				channel := channeltypes.NewChannel(channeltypes.INIT, order, counterparty, []string{path.EndpointB.ConnectionID}, path.EndpointA.ChannelConfig.Version)

				channelByte, _ := marshaler.Marshal(&channel)
				channelLen := make([]byte, 8)
				binary.BigEndian.PutUint64(channelLen, uint64(len(channelByte)))

				input = append(input, channelLen...)
				input = append(input, channelByte...)

				input = append(getChanOpenInitSignature, input...)
				admin := allowlist.TestAdminAddr
				enableds := allowlist.TestEnabledAddr

				allowlist.SetAllowListRole(vmenv.StateDB, ContractAddress, admin, allowlist.AdminRole)

				contract := createIbcGoPrecompile()

				suppliedGas := uint64(10000000)

				connection, _ := suite.chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(suite.chainA.GetContext(), path.EndpointA.ConnectionID)
				connectionByte := marshaler.MustMarshal(&connection)
				connectionsPath := fmt.Sprintf("connections/%s", path.EndpointA.ConnectionID)
				vmenv.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

				cs, _ := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID)
				cStore := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointA.ClientID)

				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					clientStateByte := clienttypes.MustMarshalClientState(marshaler, cs)

					clientStatePath := fmt.Sprintf("clients/%s/clientState", connection.GetClientID())
					vmenv.StateDB.SetPrecompileState(
						common.BytesToAddress([]byte(clientStatePath)),
						clientStateByte,
					)
					consensusStateByte := clienttypes.MustMarshalConsensusState(marshaler, consensusState)
					consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", connection.GetClientID(), clientState.GetLatestHeight())
					vmenv.StateDB.SetPrecompileState(
						common.BytesToAddress([]byte(consensusStatePath)),
						consensusStateByte,
					)
				}

				_, _, err := contract.Run(vmenv, admin, enableds, input, suppliedGas, false)

				// Testcase must have expectedPass = true AND channel order supported before
				// asserting the channel handshake initiation succeeded
				if tc.expPass {
					suite.Require().NoError(err)
				} else {
					suite.Require().Error(err)
				}
			}
		})
	}
}

// TestChanOpenTry tests the OpenTry handshake call for channels. It uses message passing
// to enter into the appropriate state and then calls ChanOpenTry directly. The channel
// is being created on chainB. The port capability must be created on chainB before
// ChanOpenTry can succeed.
func (suite *KeeperTestSuite) TestChanOpenTry() {
	var (
		path       *ibctesting.Path
		heightDiff uint64
	)

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{"success", func() {
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()
			err := path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			suite.chainB.CreatePortCapability(suite.chainB.GetSimApp().ScopedIBCMockKeeper, ibctesting.MockPort)
		}, true},
		{"connection is not OPEN", func() {
			suite.coordinator.SetupClients(path)
			// pass capability check
			suite.chainB.CreatePortCapability(suite.chainB.GetSimApp().ScopedIBCMockKeeper, ibctesting.MockPort)

			err := path.EndpointB.ConnOpenInit()
			suite.Require().NoError(err)
		}, false},
		{"consensus state not found", func() {
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()
			err := path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			suite.chainB.CreatePortCapability(suite.chainB.GetSimApp().ScopedIBCMockKeeper, ibctesting.MockPort)

			heightDiff = 3 // consensus state doesn't exist at this height
		}, false},
		{"channel verification failed", func() {
			// not creating a channel on chainA will result in an invalid proof of existence
			suite.coordinator.SetupConnections(path)
		}, false},
		{"connection version not negotiated", func() {
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()
			err := path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			// modify connB versions
			conn := path.EndpointB.GetConnection()

			version := connectiontypes.NewVersion("2", []string{"ORDER_ORDERED", "ORDER_UNORDERED"})
			conn.Versions = append(conn.Versions, version)

			suite.chainB.App.GetIBCKeeper().ConnectionKeeper.SetConnection(
				suite.chainB.GetContext(),
				path.EndpointB.ConnectionID, conn,
			)
			suite.chainB.CreatePortCapability(suite.chainB.GetSimApp().ScopedIBCMockKeeper, ibctesting.MockPort)
		}, false},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			heightDiff = 0    // must be explicitly changed in malleate
			path = ibctesting.NewPath(suite.chainA, suite.chainB)

			statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			statedb.Finalise(true)
			vmctx := vm.BlockContext{
				CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
				Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
			}
			vmenv := vm.NewEVM(vmctx, vm.TxContext{}, statedb, params.TestChainConfig, vm.Config{ExtraEips: []int{2200}})

			tc.malleate()

			if path.EndpointB.ClientID != "" {
				// ensure client is up to date
				err := path.EndpointB.UpdateClient()
				suite.Require().NoError(err)
			}

			counterparty := channeltypes.NewCounterparty(path.EndpointB.ChannelConfig.PortID, ibctesting.FirstChannelID)

			channelKey := hosttypes.ChannelKey(counterparty.PortId, counterparty.ChannelId)
			proof, proofHeight := suite.chainA.QueryProof(channelKey)

			var input []byte

			portIDbyte := []byte(path.EndpointB.ChannelConfig.PortID)
			portIDLen := make([]byte, 8)
			binary.BigEndian.PutUint64(portIDLen, uint64(len(portIDbyte)))

			input = append(input, portIDLen...)
			input = append(input, portIDbyte...)

			channel := channeltypes.NewChannel(channeltypes.INIT, channeltypes.ORDERED, counterparty, []string{path.EndpointB.ConnectionID}, path.EndpointA.ChannelConfig.Version)

			channelByte, _ := marshaler.Marshal(&channel)
			channelLen := make([]byte, 8)
			binary.BigEndian.PutUint64(channelLen, uint64(len(channelByte)))

			input = append(input, channelLen...)
			input = append(input, channelByte...)

			counterpartyVersionbyte := []byte(path.EndpointA.ChannelConfig.Version)
			counterpartyVersionLen := make([]byte, 8)
			binary.BigEndian.PutUint64(counterpartyVersionLen, uint64(len(counterpartyVersionbyte)))

			input = append(input, counterpartyVersionLen...)
			input = append(input, counterpartyVersionbyte...)

			proofInitLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofInitLen, uint64(len(proof)))

			input = append(input, proofInitLen...)
			input = append(input, proof...)

			height := malleateHeight(proofHeight, heightDiff)
			consensusHeightByte, _ := marshaler.Marshal(&height)
			consensusHeightByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(consensusHeightByteLen, uint64(len(consensusHeightByte)))

			input = append(input, consensusHeightByteLen...)
			input = append(input, consensusHeightByte...)

			input = append(getChanOpenTrySignature, input...)
			admin := allowlist.TestAdminAddr
			enableds := allowlist.TestEnabledAddr

			allowlist.SetAllowListRole(vmenv.StateDB, ContractAddress, admin, allowlist.AdminRole)

			contract := createIbcGoPrecompile()

			suppliedGas := uint64(10000000)

			connection, _ := suite.chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(suite.chainB.GetContext(), path.EndpointB.ConnectionID)
			connectionByte := marshaler.MustMarshal(&connection)
			connectionsPath := fmt.Sprintf("connections/%s", path.EndpointB.ConnectionID)
			vmenv.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

			cs, _ := suite.chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainB.GetContext(), path.EndpointB.ClientID)
			cStore := suite.chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainB.GetContext(), path.EndpointB.ClientID)

			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				clientStateByte := clienttypes.MustMarshalClientState(marshaler, cs)

				clientStatePath := fmt.Sprintf("clients/%s/clientState", connection.GetClientID())
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(clientStatePath)),
					clientStateByte,
				)
				consensusStateByte := clienttypes.MustMarshalConsensusState(marshaler, consensusState)
				consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", connection.GetClientID(), clientState.GetLatestHeight())
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(consensusStatePath)),
					consensusStateByte,
				)
			}
			_, _, err := contract.Run(vmenv, admin, enableds, input, suppliedGas, false)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func malleateHeight(height clienttypes.Height, diff uint64) clienttypes.Height {
	return clienttypes.NewHeight(height.GetRevisionNumber(), height.GetRevisionHeight()+diff)
}

// TestChanOpenAck tests the OpenAck handshake call for channels. It uses message passing
// to enter into the appropriate state and then calls ChanOpenAck directly. The handshake
// call is occurring on chainA.
func (suite *KeeperTestSuite) TestChanOpenAck() {
	var (
		path                  *ibctesting.Path
		counterpartyChannelID string
		heightDiff            uint64
	)

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{"success", func() {
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()
			err := path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ChanOpenTry()
			suite.Require().NoError(err)
		}, true},
		{"success with empty stored counterparty channel ID", func() {
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()

			err := path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ChanOpenTry()
			suite.Require().NoError(err)

			// set the channel's counterparty channel identifier to empty string
			channel := path.EndpointA.GetChannel()
			channel.Counterparty.ChannelId = ""

			// use a different channel identifier
			counterpartyChannelID = path.EndpointB.ChannelID

			suite.chainA.App.GetIBCKeeper().ChannelKeeper.SetChannel(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, channel)
		}, true},
		{"channel state is not INIT", func() {
			// create fully open channels on both chains
			suite.coordinator.Setup(path)
		}, false},
		{"connection not found", func() {
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()
			err := path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ChanOpenTry()
			suite.Require().NoError(err)

			// set the channel's connection hops to wrong connection ID
			channel := path.EndpointA.GetChannel()
			channel.ConnectionHops[0] = doesnotexist
			suite.chainA.App.GetIBCKeeper().ChannelKeeper.SetChannel(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, channel)
		}, false},
		{"connection is not OPEN", func() {
			suite.coordinator.SetupClients(path)

			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// create channel in init
			path.SetChannelOrdered()

			err = path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			suite.chainA.CreateChannelCapability(suite.chainA.GetSimApp().ScopedIBCMockKeeper, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
		}, false},
		{"consensus state not found", func() {
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()

			err := path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ChanOpenTry()
			suite.Require().NoError(err)

			heightDiff = 3 // consensus state doesn't exist at this height
		}, false},
		{"channel verification failed", func() {
			// chainB is INIT, chainA in TRYOPEN
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()

			err := path.EndpointB.ChanOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointA.ChanOpenTry()
			suite.Require().NoError(err)
		}, false},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()          // reset
			counterpartyChannelID = "" // must be explicitly changed in malleate
			heightDiff = 0             // must be explicitly changed
			path = ibctesting.NewPath(suite.chainA, suite.chainB)

			statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			statedb.Finalise(true)
			vmctx := vm.BlockContext{
				CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
				Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
			}
			vmenv := vm.NewEVM(vmctx, vm.TxContext{}, statedb, params.TestChainConfig, vm.Config{ExtraEips: []int{2200}})

			tc.malleate()

			if counterpartyChannelID == "" {
				counterpartyChannelID = ibctesting.FirstChannelID
			}

			if path.EndpointA.ClientID != "" {
				// ensure client is up to date
				err := path.EndpointA.UpdateClient()
				suite.Require().NoError(err)
			}

			channelKey := hosttypes.ChannelKey(path.EndpointB.ChannelConfig.PortID, ibctesting.FirstChannelID)
			proof, proofHeight := suite.chainB.QueryProof(channelKey)

			var input []byte

			portIDbyte := []byte(path.EndpointA.ChannelConfig.PortID)
			portIDLen := make([]byte, 8)
			binary.BigEndian.PutUint64(portIDLen, uint64(len(portIDbyte)))

			input = append(input, portIDLen...)
			input = append(input, portIDbyte...)

			channelIdbyte := []byte(path.EndpointA.ChannelID)
			channelIdLen := make([]byte, 8)
			binary.BigEndian.PutUint64(channelIdLen, uint64(len(channelIdbyte)))

			input = append(input, channelIdLen...)
			input = append(input, channelIdbyte...)

			counterpartyChannelIdbyte := []byte(counterpartyChannelID)
			counterpartyChannelIdLen := make([]byte, 8)
			binary.BigEndian.PutUint64(counterpartyChannelIdLen, uint64(len(counterpartyChannelIdbyte)))

			input = append(input, counterpartyChannelIdLen...)
			input = append(input, counterpartyChannelIdbyte...)

			counterpartyVersionbyte := []byte(path.EndpointB.ChannelConfig.Version)
			counterpartyVersionLen := make([]byte, 8)
			binary.BigEndian.PutUint64(counterpartyVersionLen, uint64(len(counterpartyVersionbyte)))

			input = append(input, counterpartyVersionLen...)
			input = append(input, counterpartyVersionbyte...)

			ProofTryLen := make([]byte, 8)
			binary.BigEndian.PutUint64(ProofTryLen, uint64(len(proof)))

			input = append(input, ProofTryLen...)
			input = append(input, proof...)

			height := malleateHeight(proofHeight, heightDiff)
			consensusHeightByte, _ := marshaler.Marshal(&height)
			consensusHeightByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(consensusHeightByteLen, uint64(len(consensusHeightByte)))

			input = append(input, consensusHeightByteLen...)
			input = append(input, consensusHeightByte...)

			input = append(getChanOpenAckSignature, input...)
			admin := allowlist.TestAdminAddr
			enableds := allowlist.TestEnabledAddr

			allowlist.SetAllowListRole(vmenv.StateDB, ContractAddress, admin, allowlist.AdminRole)

			contract := createIbcGoPrecompile()

			suppliedGas := uint64(10000000)

			channel, _ := suite.chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			chanByte := marshaler.MustMarshal(&channel)
			chanPath := hosttypes.ChannelKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			vmenv.GetStateDB().SetPrecompileState(common.BytesToAddress(chanPath), chanByte)

			connection, _ := suite.chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(suite.chainA.GetContext(), path.EndpointA.ConnectionID)
			connectionByte := marshaler.MustMarshal(&connection)
			connectionsPath := fmt.Sprintf("connections/%s", path.EndpointA.ConnectionID)
			vmenv.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

			cs, _ := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID)
			cStore := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointA.ClientID)

			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				clientStateByte := clienttypes.MustMarshalClientState(marshaler, cs)

				clientStatePath := fmt.Sprintf("clients/%s/clientState", connection.GetClientID())
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(clientStatePath)),
					clientStateByte,
				)
				consensusStateByte := clienttypes.MustMarshalConsensusState(marshaler, consensusState)
				consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", connection.GetClientID(), clientState.GetLatestHeight())
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(consensusStatePath)),
					consensusStateByte,
				)
			}
			_, _, err := contract.Run(vmenv, admin, enableds, input, suppliedGas, false)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// TestChanOpenConfirm tests the OpenAck handshake call for channels. It uses message passing
// to enter into the appropriate state and then calls ChanOpenConfirm directly. The handshake
// call is occurring on chainB.
func (suite *KeeperTestSuite) TestChanOpenConfirm() {
	var (
		path       *ibctesting.Path
		heightDiff uint64
	)

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{"success", func() {
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()

			err := path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ChanOpenTry()
			suite.Require().NoError(err)

			err = path.EndpointA.ChanOpenAck()
			suite.Require().NoError(err)

		}, true},
		{"channel state is not TRYOPEN", func() {
			// create fully open channels on both cahins
			suite.coordinator.Setup(path)
		}, false},
		{"connection not found", func() {
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()

			err := path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ChanOpenTry()
			suite.Require().NoError(err)

			err = path.EndpointA.ChanOpenAck()
			suite.Require().NoError(err)

			// set the channel's connection hops to wrong connection ID
			channel := path.EndpointB.GetChannel()
			channel.ConnectionHops[0] = doesnotexist
			suite.chainB.App.GetIBCKeeper().ChannelKeeper.SetChannel(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, channel)
		}, false},
		{"connection is not OPEN", func() {
			suite.coordinator.SetupClients(path)

			err := path.EndpointB.ConnOpenInit()
			suite.Require().NoError(err)

			suite.chainB.CreateChannelCapability(suite.chainB.GetSimApp().ScopedIBCMockKeeper, path.EndpointB.ChannelConfig.PortID, ibctesting.FirstChannelID)
		}, false},
		{"consensus state not found", func() {
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()

			err := path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ChanOpenTry()
			suite.Require().NoError(err)

			err = path.EndpointA.ChanOpenAck()
			suite.Require().NoError(err)

			heightDiff = 3
		}, false},
		{"channel verification failed", func() {
			// chainA is INIT, chainB in TRYOPEN
			suite.coordinator.SetupConnections(path)
			path.SetChannelOrdered()

			err := path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ChanOpenTry()
			suite.Require().NoError(err)
		}, false},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			heightDiff = 0    // must be explicitly changed
			path = ibctesting.NewPath(suite.chainA, suite.chainB)

			statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			statedb.Finalise(true)
			vmctx := vm.BlockContext{
				CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
				Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
			}
			vmenv := vm.NewEVM(vmctx, vm.TxContext{}, statedb, params.TestChainConfig, vm.Config{ExtraEips: []int{2200}})

			tc.malleate()

			if path.EndpointB.ClientID != "" {
				// ensure client is up to date
				err := path.EndpointB.UpdateClient()
				suite.Require().NoError(err)

			}

			channelKey := hosttypes.ChannelKey(path.EndpointA.ChannelConfig.PortID, ibctesting.FirstChannelID)
			proof, proofHeight := suite.chainA.QueryProof(channelKey)

			var input []byte

			portIDbyte := []byte(path.EndpointB.ChannelConfig.PortID)
			portIDLen := make([]byte, 8)
			binary.BigEndian.PutUint64(portIDLen, uint64(len(portIDbyte)))

			input = append(input, portIDLen...)
			input = append(input, portIDbyte...)

			channelIdbyte := []byte(ibctesting.FirstChannelID)
			channelIdLen := make([]byte, 8)
			binary.BigEndian.PutUint64(channelIdLen, uint64(len(channelIdbyte)))

			input = append(input, channelIdLen...)
			input = append(input, channelIdbyte...)

			ProofTryLen := make([]byte, 8)
			binary.BigEndian.PutUint64(ProofTryLen, uint64(len(proof)))

			input = append(input, ProofTryLen...)
			input = append(input, proof...)

			height := malleateHeight(proofHeight, heightDiff)
			consensusHeightByte, _ := marshaler.Marshal(&height)
			consensusHeightByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(consensusHeightByteLen, uint64(len(consensusHeightByte)))

			input = append(input, consensusHeightByteLen...)
			input = append(input, consensusHeightByte...)

			input = append(getChanOpenConfirmSignature, input...)
			admin := allowlist.TestAdminAddr
			enableds := allowlist.TestEnabledAddr

			allowlist.SetAllowListRole(vmenv.StateDB, ContractAddress, admin, allowlist.AdminRole)

			contract := createIbcGoPrecompile()

			suppliedGas := uint64(10000000)

			channel, _ := suite.chainB.App.GetIBCKeeper().ChannelKeeper.GetChannel(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			chanByte := marshaler.MustMarshal(&channel)
			chanPath := hosttypes.ChannelKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			vmenv.GetStateDB().SetPrecompileState(common.BytesToAddress(chanPath), chanByte)

			connection, _ := suite.chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(suite.chainB.GetContext(), path.EndpointB.ConnectionID)
			connectionByte := marshaler.MustMarshal(&connection)
			connectionsPath := fmt.Sprintf("connections/%s", path.EndpointB.ConnectionID)
			vmenv.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

			cs, _ := suite.chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainB.GetContext(), path.EndpointB.ClientID)
			cStore := suite.chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainB.GetContext(), path.EndpointB.ClientID)

			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				clientStateByte := clienttypes.MustMarshalClientState(marshaler, cs)

				clientStatePath := fmt.Sprintf("clients/%s/clientState", connection.GetClientID())
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(clientStatePath)),
					clientStateByte,
				)
				consensusStateByte := clienttypes.MustMarshalConsensusState(marshaler, consensusState)
				consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", connection.GetClientID(), clientState.GetLatestHeight())
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(consensusStatePath)),
					consensusStateByte,
				)
			}

			_, _, err := contract.Run(vmenv, admin, enableds, input, suppliedGas, false)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// TestChanCloseInit tests the initial closing of a handshake on chainA by calling
// ChanCloseInit. Both chains will use message passing to setup OPEN channels.
func (suite *KeeperTestSuite) TestChanCloseInit() {
	var (
		path *ibctesting.Path
	)

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{"success", func() {
			suite.coordinator.Setup(path)
		}, true},
		{"channel doesn't exist", func() {
			// any non-nil values work for connections
			path.EndpointA.ConnectionID = ibctesting.FirstConnectionID
			path.EndpointB.ConnectionID = ibctesting.FirstConnectionID

			path.EndpointA.ChannelID = ibctesting.FirstChannelID
			path.EndpointB.ChannelID = ibctesting.FirstChannelID

			// ensure channel capability check passes
			suite.chainA.CreateChannelCapability(suite.chainA.GetSimApp().ScopedIBCMockKeeper, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
		}, false},
		{"channel state is CLOSED", func() {
			suite.coordinator.Setup(path)

			// close channel
			err := path.EndpointA.SetChannelState(channeltypes.CLOSED)
			suite.Require().NoError(err)
		}, false},
		{"connection not found", func() {
			suite.coordinator.Setup(path)

			// set the channel's connection hops to wrong connection ID
			channel := path.EndpointA.GetChannel()
			channel.ConnectionHops[0] = doesnotexist
			suite.chainA.App.GetIBCKeeper().ChannelKeeper.SetChannel(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, channel)
		}, false},
		{"connection is not OPEN", func() {
			suite.coordinator.SetupClients(path)

			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// create channel in init
			path.SetChannelOrdered()
			err = path.EndpointA.ChanOpenInit()
			suite.Require().NoError(err)

			// ensure channel capability check passes
			suite.chainA.CreateChannelCapability(suite.chainA.GetSimApp().ScopedIBCMockKeeper, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
		}, false},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			path = ibctesting.NewPath(suite.chainA, suite.chainB)

			statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			statedb.Finalise(true)
			vmctx := vm.BlockContext{
				CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
				Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
			}
			vmenv := vm.NewEVM(vmctx, vm.TxContext{}, statedb, params.TestChainConfig, vm.Config{ExtraEips: []int{2200}})

			tc.malleate()

			var input []byte

			portIDbyte := []byte(path.EndpointA.ChannelConfig.PortID)
			portIDLen := make([]byte, 8)
			binary.BigEndian.PutUint64(portIDLen, uint64(len(portIDbyte)))

			input = append(input, portIDLen...)
			input = append(input, portIDbyte...)

			channelIdbyte := []byte(ibctesting.FirstChannelID)
			channelIdLen := make([]byte, 8)
			binary.BigEndian.PutUint64(channelIdLen, uint64(len(channelIdbyte)))

			input = append(input, channelIdLen...)
			input = append(input, channelIdbyte...)

			input = append(getChanCloseInitSignature, input...)
			admin := allowlist.TestAdminAddr
			enableds := allowlist.TestEnabledAddr

			allowlist.SetAllowListRole(vmenv.StateDB, ContractAddress, admin, allowlist.AdminRole)

			contract := createIbcGoPrecompile()

			suppliedGas := uint64(10000000)

			channel, _ := suite.chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			chanByte := marshaler.MustMarshal(&channel)
			chanPath := hosttypes.ChannelKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			vmenv.GetStateDB().SetPrecompileState(common.BytesToAddress(chanPath), chanByte)

			connection, _ := suite.chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(suite.chainA.GetContext(), path.EndpointA.ConnectionID)
			connectionByte := marshaler.MustMarshal(&connection)
			connectionsPath := fmt.Sprintf("connections/%s", path.EndpointA.ConnectionID)
			vmenv.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

			cs, _ := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID)
			cStore := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointA.ClientID)

			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				clientStateByte := clienttypes.MustMarshalClientState(marshaler, cs)

				clientStatePath := fmt.Sprintf("clients/%s/clientState", connection.GetClientID())
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(clientStatePath)),
					clientStateByte,
				)
				consensusStateByte := clienttypes.MustMarshalConsensusState(marshaler, consensusState)
				consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", connection.GetClientID(), clientState.GetLatestHeight())
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(consensusStatePath)),
					consensusStateByte,
				)
			}

			_, _, err := contract.Run(vmenv, admin, enableds, input, suppliedGas, false)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// TestChanCloseConfirm tests the confirming closing channel ends by calling ChanCloseConfirm
// on chainB. Both chains will use message passing to setup OPEN channels. ChanCloseInit is
// bypassed on chainA by setting the channel state in the ChannelKeeper.
func (suite *KeeperTestSuite) TestChanCloseConfirm() {
	var (
		path       *ibctesting.Path
		heightDiff uint64
	)

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{"success", func() {
			suite.coordinator.Setup(path)
			err := path.EndpointA.SetChannelState(channeltypes.CLOSED)
			suite.Require().NoError(err)
		}, true},
		{"channel state is CLOSED", func() {
			suite.coordinator.Setup(path)

			err := path.EndpointB.SetChannelState(channeltypes.CLOSED)
			suite.Require().NoError(err)
		}, false},
		{"connection not found", func() {
			suite.coordinator.Setup(path)

			// set the channel's connection hops to wrong connection ID
			channel := path.EndpointB.GetChannel()
			channel.ConnectionHops[0] = doesnotexist
			suite.chainB.App.GetIBCKeeper().ChannelKeeper.SetChannel(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, channel)
		}, false},
		{"connection is not OPEN", func() {
			suite.coordinator.SetupClients(path)

			err := path.EndpointB.ConnOpenInit()
			suite.Require().NoError(err)

			// create channel in init
			path.SetChannelOrdered()
			err = path.EndpointB.ChanOpenInit()
			suite.Require().NoError(err)

			// ensure channel capability check passes
			suite.chainB.CreateChannelCapability(suite.chainB.GetSimApp().ScopedIBCMockKeeper, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
		}, false},
		{"consensus state not found", func() {
			suite.coordinator.Setup(path)

			err := path.EndpointA.SetChannelState(channeltypes.CLOSED)
			suite.Require().NoError(err)

			heightDiff = 3
		}, false},
		{"channel verification failed", func() {
			// channel not closed
			suite.coordinator.Setup(path)
		}, false},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			heightDiff = 0    // must explicitly be changed
			path = ibctesting.NewPath(suite.chainA, suite.chainB)

			statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			statedb.Finalise(true)
			vmctx := vm.BlockContext{
				CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
				Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
			}
			vmenv := vm.NewEVM(vmctx, vm.TxContext{}, statedb, params.TestChainConfig, vm.Config{ExtraEips: []int{2200}})

			tc.malleate()

			channelKey := hosttypes.ChannelKey(path.EndpointA.ChannelConfig.PortID, ibctesting.FirstChannelID)
			proof, proofHeight := suite.chainA.QueryProof(channelKey)

			var input []byte

			portIDbyte := []byte(path.EndpointB.ChannelConfig.PortID)
			portIDLen := make([]byte, 8)
			binary.BigEndian.PutUint64(portIDLen, uint64(len(portIDbyte)))

			input = append(input, portIDLen...)
			input = append(input, portIDbyte...)

			channelIdbyte := []byte(ibctesting.FirstChannelID)
			channelIdLen := make([]byte, 8)
			binary.BigEndian.PutUint64(channelIdLen, uint64(len(channelIdbyte)))

			input = append(input, channelIdLen...)
			input = append(input, channelIdbyte...)

			proofInitLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofInitLen, uint64(len(proof)))

			input = append(input, proofInitLen...)
			input = append(input, proof...)

			height := malleateHeight(proofHeight, heightDiff)
			consensusHeightByte, _ := marshaler.Marshal(&height)
			consensusHeightByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(consensusHeightByteLen, uint64(len(consensusHeightByte)))

			input = append(input, consensusHeightByteLen...)
			input = append(input, consensusHeightByte...)

			input = append(getChanCloseConfirmSignature, input...)
			admin := allowlist.TestAdminAddr
			enableds := allowlist.TestEnabledAddr

			allowlist.SetAllowListRole(vmenv.StateDB, ContractAddress, admin, allowlist.AdminRole)

			contract := createIbcGoPrecompile()

			suppliedGas := uint64(10000000)

			channel, _ := suite.chainB.App.GetIBCKeeper().ChannelKeeper.GetChannel(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			chanByte := marshaler.MustMarshal(&channel)
			chanPath := hosttypes.ChannelKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			vmenv.GetStateDB().SetPrecompileState(common.BytesToAddress(chanPath), chanByte)

			connection, _ := suite.chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(suite.chainB.GetContext(), path.EndpointB.ConnectionID)
			connectionByte := marshaler.MustMarshal(&connection)
			connectionsPath := fmt.Sprintf("connections/%s", path.EndpointB.ConnectionID)
			vmenv.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

			cs, _ := suite.chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainB.GetContext(), path.EndpointB.ClientID)
			cStore := suite.chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainB.GetContext(), path.EndpointB.ClientID)

			if cs != nil {
				clientState := cs.(*ibctm.ClientState)
				bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
				consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
				clientStateByte := clienttypes.MustMarshalClientState(marshaler, cs)

				clientStatePath := fmt.Sprintf("clients/%s/clientState", connection.GetClientID())
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(clientStatePath)),
					clientStateByte,
				)
				consensusStateByte := clienttypes.MustMarshalConsensusState(marshaler, consensusState)
				consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", connection.GetClientID(), clientState.GetLatestHeight())
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(consensusStatePath)),
					consensusStateByte,
				)
			}

			_, _, err := contract.Run(vmenv, admin, enableds, input, suppliedGas, false)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
