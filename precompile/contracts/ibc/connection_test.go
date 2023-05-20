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
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
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
