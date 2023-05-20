package ibc

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/ethereum/go-ethereum/common"
)

// TestConnOpenInit - chainA initializes (INIT state) a connection with
// chainB which is yet UNINITIALIZED
func (suite *KeeperTestSuite) TestConnOpenInit() {
	var (
		path         *ibctesting.Path
		version      *connectiontypes.Version
		delayPeriod  uint64
		emptyConnBID bool
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{"success", func() {
		}, true},
		{"success with empty counterparty identifier", func() {
			emptyConnBID = true
		}, true},
		{"success with non empty version", func() {
			version = connectiontypes.ExportedVersionsToProto(connectiontypes.GetCompatibleVersions())[0]
		}, true},
		{"success with non zero delayPeriod", func() {
			delayPeriod = uint64(time.Hour.Nanoseconds())
		}, true},
	}

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.msg, func() {

			statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			statedb.Finalise(true)
			vmctx := vm.BlockContext{
				CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
				Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
			}
			vmenv := vm.NewEVM(vmctx, vm.TxContext{}, statedb, params.TestChainConfig, vm.Config{ExtraEips: []int{2200}})

			suite.SetupTest()    // reset
			emptyConnBID = false // must be explicitly changed
			version = nil        // must be explicitly

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			tc.malleate()

			if emptyConnBID {
				path.EndpointB.ConnectionID = ""
			}
			counterparty := connectiontypes.NewCounterparty(path.EndpointB.ClientID, path.EndpointB.ConnectionID, suite.chainB.GetPrefix())

			var input []byte

			clientIDByte := []byte(path.EndpointA.ClientID)
			clientIDByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(clientIDByteLen, uint64(len(clientIDByte)))

			input = append(input, clientIDByteLen...)
			input = append(input, clientIDByte...)

			counterpartybyte, _ := marshaler.MarshalInterface(&counterparty)
			counterpartyLen := make([]byte, 8)
			binary.BigEndian.PutUint64(counterpartyLen, uint64(len(counterpartybyte)))

			input = append(input, counterpartyLen...)
			input = append(input, counterpartybyte...)

			var versionbyte []byte
			var versionLen []byte
			if version == nil {
				versionbyte = []byte("")
				versionLen = make([]byte, 8)
				binary.BigEndian.PutUint64(versionLen, uint64(len(versionbyte)))
			} else {
				versionbyte, _ = version.Marshal()
				versionLen = make([]byte, 8)
				binary.BigEndian.PutUint64(versionLen, uint64(len(versionbyte)))
			}

			input = append(input, versionLen...)
			input = append(input, versionbyte...)

			delayPeriodByte := make([]byte, 8)
			binary.BigEndian.PutUint64(delayPeriodByte, uint64(delayPeriod))

			input = append(input, delayPeriodByte...)

			input = append(getConnOpenInitSignature, input...)
			admin := allowlist.TestAdminAddr
			enableds := allowlist.TestEnabledAddr

			allowlist.SetAllowListRole(vmenv.StateDB, ContractAddress, admin, allowlist.AdminRole)

			contract := createIbcGoPrecompile()

			suppliedGas := uint64(10000000)

			cs, _ := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID)
			clientStateByte, _ := marshaler.MarshalInterface(cs.(*ibctm.ClientState))
			clientStatePath := fmt.Sprintf("clients/%s/clientState", path.EndpointA.ClientID)
			vmenv.StateDB.SetPrecompileState(
				common.BytesToAddress([]byte(clientStatePath)),
				clientStateByte,
			)

			output, _, err := contract.Run(vmenv, admin, enableds, input, suppliedGas, false)

			connectionID := string(output)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(connectiontypes.FormatConnectionIdentifier(0), connectionID)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal("", connectionID)
			}
		})
	}
}

// TestConnOpenConfirm - chainB calls ConnOpenConfirm to confirm that
// chainA state is now OPEN.
func (suite *KeeperTestSuite) TestConnOpenConfirm() {
	var path *ibctesting.Path

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
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)

			err = path.EndpointA.ConnOpenAck()
			suite.Require().NoError(err)
		}, true},
		{"connection not found", func() {
			// connections are never created
		}, false},
		{"chain B's connection state is not TRYOPEN", func() {
			// connections are OPEN
			suite.coordinator.CreateConnections(path)
		}, false},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.msg, func() {
			statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			statedb.Finalise(true)
			vmctx := vm.BlockContext{
				CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
				Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
			}
			vmenv := vm.NewEVM(vmctx, vm.TxContext{}, statedb, params.TestChainConfig, vm.Config{ExtraEips: []int{2200}})

			suite.SetupTest() // reset
			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			tc.malleate()

			// ensure client is up to date to receive proof
			err := path.EndpointB.UpdateClient()
			suite.Require().NoError(err)

			connectionKey := host.ConnectionKey(path.EndpointA.ConnectionID)
			proofAck, proofHeight := suite.chainA.QueryProof(connectionKey)

			var input []byte

			ConnectionIDByte := []byte(path.EndpointB.ConnectionID)
			ConnectionIDByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(ConnectionIDByteLen, uint64(len(ConnectionIDByte)))

			input = append(input, ConnectionIDByteLen...)
			input = append(input, ConnectionIDByte...)

			proofAckLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofAckLen, uint64(len(proofAck)))

			input = append(input, proofAckLen...)
			input = append(input, proofAck...)

			proofHeightByte, _ := marshaler.Marshal(&proofHeight)
			proofHeightByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofHeightByteLen, uint64(len(proofHeightByte)))

			input = append(input, proofHeightByteLen...)
			input = append(input, proofHeightByte...)

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

			input = append(getConnOpenConfirmSignature, input...)
			admin := allowlist.TestAdminAddr
			enableds := allowlist.TestEnabledAddr

			allowlist.SetAllowListRole(vmenv.StateDB, ContractAddress, admin, allowlist.AdminRole)

			contract := createIbcGoPrecompile()

			suppliedGas := uint64(10000000)

			_, _, err = contract.Run(vmenv, admin, enableds, input, suppliedGas, false)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
