package ibc_old

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

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

	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
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

// TestConnOpenTry - chainB calls ConnOpenTry to verify the state of
// connection on chainA is INIT
func (suite *KeeperTestSuite) TestConnOpenTry() {
	var (
		path               *ibctesting.Path
		delayPeriod        uint64
		versions           []exported.Version
		consensusHeight    clienttypes.Height
		counterpartyClient exported.ClientState
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
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainA to pass as counterpartyClient
			counterpartyClient = suite.chainA.GetClientState(path.EndpointA.ClientID)
		}, true},
		{"success with delay period", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			delayPeriod = uint64(time.Hour.Nanoseconds())

			// set delay period on counterparty to non-zero value
			conn := path.EndpointA.GetConnection()
			conn.DelayPeriod = delayPeriod
			suite.chainA.App.GetIBCKeeper().ConnectionKeeper.SetConnection(suite.chainA.GetContext(), path.EndpointA.ConnectionID, conn)

			// commit in order for proof to return correct value
			suite.coordinator.CommitBlock(suite.chainA)
			err = path.EndpointB.UpdateClient()
			suite.Require().NoError(err)

			// retrieve client state of chainA to pass as counterpartyClient
			counterpartyClient = suite.chainA.GetClientState(path.EndpointA.ClientID)
		}, true},
		{"invalid counterparty client", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainA.GetClientState(path.EndpointA.ClientID)

			// Set an invalid client of chainA on chainB
			tmClient, ok := counterpartyClient.(*ibctm.ClientState)
			suite.Require().True(ok)
			tmClient.ChainId = "wrongchainid"

			suite.chainA.App.GetIBCKeeper().ClientKeeper.SetClientState(suite.chainA.GetContext(), path.EndpointA.ClientID, tmClient)
		}, false},
		{"counterparty versions is empty", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainA to pass as counterpartyClient
			counterpartyClient = suite.chainA.GetClientState(path.EndpointA.ClientID)

			versions = nil
		}, false},
		{"counterparty versions don't have a match", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainA to pass as counterpartyClient
			counterpartyClient = suite.chainA.GetClientState(path.EndpointA.ClientID)

			version := connectiontypes.NewVersion("0.0", nil)
			versions = []exported.Version{version}
		}, false},
		{"connection state verification failed", func() {
			// chainA connection not created

			// retrieve client state of chainA to pass as counterpartyClient
			counterpartyClient = suite.chainA.GetClientState(path.EndpointA.ClientID)
		}, false},
		{"client state verification failed", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainA to pass as counterpartyClient
			counterpartyClient = suite.chainA.GetClientState(path.EndpointA.ClientID)

			// modify counterparty client without setting in store so it still passes validate but fails proof verification
			tmClient, ok := counterpartyClient.(*ibctm.ClientState)
			suite.Require().True(ok)
			tmClient.LatestHeight = tmClient.LatestHeight.Increment().(clienttypes.Height)
		}, false},
		// {"consensus state verification failed", func() {
		// 	// retrieve client state of chainA to pass as counterpartyClient
		// 	counterpartyClient = suite.chainA.GetClientState(path.EndpointA.ClientID)

		// 	// give chainA wrong consensus state for chainB
		// 	consState, found := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetLatestClientConsensusState(suite.chainA.GetContext(), path.EndpointA.ClientID)
		// 	suite.Require().True(found)

		// 	tmConsState, ok := consState.(*ibctm.ConsensusState)
		// 	suite.Require().True(ok)

		// 	tmConsState.Timestamp = time.Now()
		// 	suite.chainA.App.GetIBCKeeper().ClientKeeper.SetClientConsensusState(suite.chainA.GetContext(), path.EndpointA.ClientID, counterpartyClient.GetLatestHeight(), tmConsState)

		// 	err := path.EndpointA.ConnOpenInit()
		// 	suite.Require().NoError(err)
		// }, false},
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

			suite.SetupTest()                                  // reset
			consensusHeight = clienttypes.ZeroHeight()         // may be changed in malleate
			versions = connectiontypes.GetCompatibleVersions() // may be changed in malleate
			delayPeriod = 0                                    // may be changed in malleate
			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			tc.malleate()

			counterparty := connectiontypes.NewCounterparty(path.EndpointA.ClientID, path.EndpointA.ConnectionID, suite.chainA.GetPrefix())

			// ensure client is up to date to receive proof
			err := path.EndpointB.UpdateClient()
			suite.Require().NoError(err)

			connectionKey := host.ConnectionKey(path.EndpointA.ConnectionID)
			proofInit, proofHeight := suite.chainA.QueryProof(connectionKey)
			fmt.Printf("proofInit %#v\n", proofInit)
			fmt.Printf("proofHeight %#v\n", proofHeight)
			if consensusHeight.IsZero() {
				// retrieve consensus state height to provide proof for
				consensusHeight = counterpartyClient.GetLatestHeight().(clienttypes.Height)
			}
			consensusKey := host.FullConsensusStateKey(path.EndpointA.ClientID, consensusHeight)
			proofConsensus, _ := suite.chainA.QueryProof(consensusKey)
			fmt.Printf("proofConsensus %#v\n", proofConsensus)

			// retrieve proof of counterparty clientstate on chainA
			clientKey := host.FullClientStateKey(path.EndpointA.ClientID)
			proofClient, _ := suite.chainA.QueryProof(clientKey)
			fmt.Printf("proofClient %#v\n", proofClient)

			var input []byte

			counterpartyByte, _ := counterparty.Marshal()
			fmt.Printf("counterparty %#v\n", counterpartyByte)

			counterpartyByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(counterpartyByteLen, uint64(len(counterpartyByte)))

			input = append(input, counterpartyByteLen...)
			input = append(input, counterpartyByte...)

			delayPeriodByte := make([]byte, 8)
			binary.BigEndian.PutUint64(delayPeriodByte, uint64(delayPeriod))

			input = append(input, delayPeriodByte...)

			clientIDByte := []byte(path.EndpointB.ClientID)
			fmt.Printf("ClientID %s\n", path.EndpointB.ClientID)

			clientIDByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(clientIDByteLen, uint64(len(clientIDByte)))

			input = append(input, clientIDByteLen...)
			input = append(input, clientIDByte...)

			// clientStateByte, _ := marshaler.MarshalInterface(counterpartyClient.(*ibctm.ClientState))
			clientStateByte, _ := clienttypes.MarshalClientState(marshaler, counterpartyClient)
			fmt.Printf("clientState %#v\n", clientStateByte)

			clientStateByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(clientStateByteLen, uint64(len(clientStateByte)))

			input = append(input, clientStateByteLen...)
			input = append(input, clientStateByte...)

			// versionsByte, _ := marshaler.Marshal()
			versionsByte, _ := json.Marshal(connectiontypes.ExportedVersionsToProto(versions))
			fmt.Printf("versions %#v\n", versionsByte)

			versionsByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(versionsByteLen, uint64(len(versionsByte)))

			input = append(input, versionsByteLen...)
			input = append(input, versionsByte...)

			proofInitLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofInitLen, uint64(len(proofInit)))

			input = append(input, proofInitLen...)
			input = append(input, proofInit...)

			proofClientLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofClientLen, uint64(len(proofClient)))

			input = append(input, proofClientLen...)
			input = append(input, proofClient...)

			proofConsensusLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofConsensusLen, uint64(len(proofConsensus)))

			input = append(input, proofConsensusLen...)
			input = append(input, proofConsensus...)

			proofHeightByte, _ := marshaler.MarshalInterface(&proofHeight)
			fmt.Printf("proofHeightByte %#v\n", proofHeightByte)

			proofHeightByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofHeightByteLen, uint64(len(proofHeightByte)))

			input = append(input, proofHeightByteLen...)
			input = append(input, proofHeightByte...)

			consensusHeightByte, _ := marshaler.MarshalInterface(&consensusHeight)
			fmt.Printf("consensusHeightByte %#v\n", consensusHeightByte)

			consensusHeightByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(consensusHeightByteLen, uint64(len(consensusHeightByte)))

			input = append(input, consensusHeightByteLen...)
			input = append(input, consensusHeightByte...)

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

				clientStatePath := fmt.Sprintf("clients/%s/clientState", path.EndpointB.ClientID)
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(clientStatePath)),
					clientStateByte,
				)
				consensusStateByte := clienttypes.MustMarshalConsensusState(marshaler, consensusState)
				consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", path.EndpointB.ClientID, clientState.GetLatestHeight())
				vmenv.StateDB.SetPrecompileState(
					common.BytesToAddress([]byte(consensusStatePath)),
					consensusStateByte,
				)
			}

			input = append(getConnOpenTrySignature, input...)
			admin := allowlist.TestAdminAddr
			enableds := allowlist.TestEnabledAddr

			allowlist.SetAllowListRole(vmenv.StateDB, ContractAddress, admin, allowlist.AdminRole)

			contract := createIbcGoPrecompile()

			suppliedGas := uint64(10000000)

			out, _, err := contract.Run(vmenv, admin, enableds, input, suppliedGas, false)
			connectionID := string(out)

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

// TestConnOpenAck - Chain A (ID #1) calls TestConnOpenAck to acknowledge (ACK state)
// the initialization (TRYINIT) of the connection on  Chain B (ID #2).
func (suite *KeeperTestSuite) TestConnOpenAck() {
	var (
		path               *ibctesting.Path
		consensusHeight    clienttypes.Height
		version            *connectiontypes.Version
		counterpartyClient exported.ClientState
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
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)
		}, true},
		{"invalid counterparty client", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)

			// Set an invalid client of chainA on chainB
			tmClient, ok := counterpartyClient.(*ibctm.ClientState)
			suite.Require().True(ok)
			tmClient.ChainId = "wrongchainid"

			suite.chainB.App.GetIBCKeeper().ClientKeeper.SetClientState(suite.chainB.GetContext(), path.EndpointB.ClientID, tmClient)
		}, false},
		{"consensus height >= latest height", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)

			consensusHeight = clienttypes.GetSelfHeight(suite.chainA.GetContext())
		}, false},
		{"connection not found", func() {
			// connections are never created

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)
		}, false},
		{"invalid counterparty connection ID", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)

			// modify connB to set counterparty connection identifier to wrong identifier
			connection, found := suite.chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(suite.chainA.GetContext(), path.EndpointA.ConnectionID)
			suite.Require().True(found)

			connection.Counterparty.ConnectionId = "badconnectionid"

			suite.chainA.App.GetIBCKeeper().ConnectionKeeper.SetConnection(suite.chainA.GetContext(), path.EndpointA.ConnectionID, connection)

			err = path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			err = path.EndpointB.UpdateClient()
			suite.Require().NoError(err)
		}, false},
		{"connection state is not INIT", func() {
			// connection state is already OPEN on chainA
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)

			err = path.EndpointA.ConnOpenAck()
			suite.Require().NoError(err)
		}, false},
		{"connection is in INIT but the proposed version is invalid", func() {
			// chainA is in INIT, chainB is in TRYOPEN
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)

			version = connectiontypes.NewVersion("2.0", nil)
		}, false},
		{"incompatible IBC versions", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)

			// set version to a non-compatible version
			version = connectiontypes.NewVersion("2.0", nil)
		}, false},
		{"empty version", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)

			version = &connectiontypes.Version{}
		}, false},
		{"feature set verification failed - unsupported feature", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)

			version = connectiontypes.NewVersion(connectiontypes.DefaultIBCVersionIdentifier, []string{"ORDER_ORDERED", "ORDER_UNORDERED", "ORDER_DAG"})
		}, false},
		{"self consensus state not found", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)

			consensusHeight = clienttypes.NewHeight(0, 1)
		}, false},
		{"connection state verification failed", func() {
			// chainB connection is not in INIT
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)
		}, false},
		{"client state verification failed", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)

			// modify counterparty client without setting in store so it still passes validate but fails proof verification
			tmClient, ok := counterpartyClient.(*ibctm.ClientState)
			suite.Require().True(ok)
			tmClient.LatestHeight = tmClient.LatestHeight.Increment().(clienttypes.Height)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)
		}, false},
		{"consensus state verification failed", func() {
			err := path.EndpointA.ConnOpenInit()
			suite.Require().NoError(err)

			// retrieve client state of chainB to pass as counterpartyClient
			counterpartyClient = suite.chainB.GetClientState(path.EndpointB.ClientID)

			// give chainB wrong consensus state for chainA
			consState, found := suite.chainB.App.GetIBCKeeper().ClientKeeper.GetLatestClientConsensusState(suite.chainB.GetContext(), path.EndpointB.ClientID)
			suite.Require().True(found)

			tmConsState, ok := consState.(*ibctm.ConsensusState)
			suite.Require().True(ok)

			tmConsState.Timestamp = tmConsState.Timestamp.Add(time.Second)
			suite.chainB.App.GetIBCKeeper().ClientKeeper.SetClientConsensusState(suite.chainB.GetContext(), path.EndpointB.ClientID, counterpartyClient.GetLatestHeight(), tmConsState)

			err = path.EndpointB.ConnOpenTry()
			suite.Require().NoError(err)
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

			suite.SetupTest()                                                                             // reset
			version = connectiontypes.ExportedVersionsToProto(connectiontypes.GetCompatibleVersions())[0] // must be explicitly changed in malleate
			consensusHeight = clienttypes.ZeroHeight()                                                    // must be explicitly changed in malleate
			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			tc.malleate()

			// ensure client is up to date to receive proof
			err := path.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			connectionKey := host.ConnectionKey(path.EndpointB.ConnectionID)
			proofTry, proofHeight := suite.chainB.QueryProof(connectionKey)

			if consensusHeight.IsZero() {
				// retrieve consensus state height to provide proof for
				clientState := suite.chainB.GetClientState(path.EndpointB.ClientID)
				consensusHeight = clientState.GetLatestHeight().(clienttypes.Height)
			}
			consensusKey := host.FullConsensusStateKey(path.EndpointB.ClientID, consensusHeight)
			proofConsensus, _ := suite.chainB.QueryProof(consensusKey)

			// retrieve proof of counterparty clientstate on chainA
			clientKey := host.FullClientStateKey(path.EndpointB.ClientID)
			proofClient, _ := suite.chainB.QueryProof(clientKey)

			var input []byte

			ConnectionIDByte := []byte(path.EndpointA.ConnectionID)
			ConnectionIDByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(ConnectionIDByteLen, uint64(len(ConnectionIDByte)))

			input = append(input, ConnectionIDByteLen...)
			input = append(input, ConnectionIDByte...)

			clientStateByte, _ := clienttypes.MarshalClientState(marshaler, counterpartyClient)
			clientStateByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(clientStateByteLen, uint64(len(clientStateByte)))

			input = append(input, clientStateByteLen...)
			input = append(input, clientStateByte...)

			versionByte, _ := marshaler.Marshal(version)
			versionByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(versionByteLen, uint64(len(versionByte)))

			input = append(input, versionByteLen...)
			input = append(input, versionByte...)

			counterpartyConnectionIDByte := []byte(path.EndpointB.ConnectionID)
			counterpartyConnectionIDByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(counterpartyConnectionIDByteLen, uint64(len(counterpartyConnectionIDByte)))

			input = append(input, counterpartyConnectionIDByteLen...)
			input = append(input, counterpartyConnectionIDByte...)

			proofTryLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofTryLen, uint64(len(proofTry)))

			input = append(input, proofTryLen...)
			input = append(input, proofTry...)

			proofClientLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofClientLen, uint64(len(proofClient)))

			input = append(input, proofClientLen...)
			input = append(input, proofClient...)

			proofConsensusLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofConsensusLen, uint64(len(proofConsensus)))

			input = append(input, proofConsensusLen...)
			input = append(input, proofConsensus...)

			proofHeightByte, _ := marshaler.Marshal(&proofHeight)
			proofHeightByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(proofHeightByteLen, uint64(len(proofHeightByte)))

			input = append(input, proofHeightByteLen...)
			input = append(input, proofHeightByte...)

			consensusHeightByte, _ := marshaler.Marshal(&consensusHeight)
			consensusHeightByteLen := make([]byte, 8)
			binary.BigEndian.PutUint64(consensusHeightByteLen, uint64(len(consensusHeightByte)))

			input = append(input, consensusHeightByteLen...)
			input = append(input, consensusHeightByte...)

			connection, _ := suite.chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(suite.chainA.GetContext(), path.EndpointA.ConnectionID)

			connectionByte, _ := marshaler.Marshal(&connection)

			connectionsPath := fmt.Sprintf("connections/%s", path.EndpointA.ConnectionID)
			vmenv.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

			cs, _ := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), connection.GetClientID())
			cStore := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), connection.GetClientID())

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

			input = append(getConnOpenAckSignature, input...)
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
