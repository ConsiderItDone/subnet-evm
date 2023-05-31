package ibc

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/ethereum/go-ethereum/common"
)

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
