package ibc

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/cometbft/cometbft/libs/bytes"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	tendermint "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	ibcmock "github.com/cosmos/ibc-go/v7/testing/mock"
	ics23 "github.com/cosmos/ics23/go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

var (
	disabledTimeoutTimestamp = uint64(0)
	disabledTimeoutHeight    = clienttypes.ZeroHeight()
	defaultTimeoutHeight     = clienttypes.NewHeight(1, 100)

	// for when the testing package cannot be used
	connIDA = "connA"
	connIDB = "connB"
)

func TestRecvPacket(t *testing.T) {
	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	path := ibctesting.NewPath(chainA, chainB)
	coordinator.Setup(path)
	data := common.FromHex("000000000000000000000000000000000000000000000000000000000000002d")

	sequence, _ := path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, data)

	tests := map[string]testutils.PrecompileTest{
		"success UNORDERED channel": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				packetKey := host.PacketCommitmentKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				proof, proofHeight := path.EndpointA.QueryProof(packetKey)

				input, err := PackRecvPacket(MsgRecvPacket{
					Packet: Packet{
						Sequence:           big.NewInt(int64(sequence)),
						SourcePort:         path.EndpointA.ChannelConfig.PortID,
						SourceChannel:      path.EndpointA.ChannelID,
						DestinationPort:    path.EndpointB.ChannelConfig.PortID,
						DestinationChannel: path.EndpointB.ChannelID,
						Data:               data,
						TimeoutHeight: Height{
							RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
							RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
						},
						TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
					},
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(proofHeight.RevisionHeight)),
					},
					ProofCommitment: proof,
					Signer:          "Signer",
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))

				SetPort(state, path.EndpointB.ChannelConfig.PortID, contractAddress)

				err := SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				if err != nil {
					t.Error(err)
				}

				channel, _ := chainB.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				SetChannel(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, &channel)

				connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
				SetConnection(state, path.EndpointB.ConnectionID, &connection)

				cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
				cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
		},
		"success: ORDERED channel": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				path.SetChannelOrdered()
				packetKey := host.PacketCommitmentKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				proof, proofHeight := path.EndpointA.QueryProof(packetKey)

				input, err := PackRecvPacket(MsgRecvPacket{
					Packet: Packet{
						Sequence:           big.NewInt(int64(sequence)),
						SourcePort:         path.EndpointA.ChannelConfig.PortID,
						SourceChannel:      path.EndpointA.ChannelID,
						DestinationPort:    path.EndpointB.ChannelConfig.PortID,
						DestinationChannel: path.EndpointB.ChannelID,
						Data:               data,
						TimeoutHeight: Height{
							RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
							RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
						},
						TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
					},
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(proofHeight.RevisionHeight)),
					},
					ProofCommitment: proof,
					Signer:          "Signer",
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))

				SetPort(state, path.EndpointB.ChannelConfig.PortID, contractAddress)

				err := SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				if err != nil {
					t.Error(err)
				}

				channel, _ := chainB.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				SetChannel(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, &channel)

				connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
				SetConnection(state, path.EndpointB.ConnectionID, &connection)

				cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
				cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
		},
		"channel not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				path.SetChannelOrdered()
				packetKey := host.PacketCommitmentKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				proof, proofHeight := path.EndpointA.QueryProof(packetKey)

				input, err := PackRecvPacket(MsgRecvPacket{
					Packet: Packet{
						Sequence:           big.NewInt(int64(sequence)),
						SourcePort:         path.EndpointA.ChannelConfig.PortID,
						SourceChannel:      path.EndpointA.ChannelID,
						DestinationPort:    path.EndpointB.ChannelConfig.PortID,
						DestinationChannel: path.EndpointB.ChannelID,
						Data:               data,
						TimeoutHeight: Height{
							RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
							RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
						},
						TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
					},
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(proofHeight.RevisionHeight)),
					},
					ProofCommitment: proof,
					Signer:          "Signer",
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))

				SetPort(state, path.EndpointB.ChannelConfig.PortID, contractAddress)

				err := SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				if err != nil {
					t.Error(err)
				}

				connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
				SetConnection(state, path.EndpointB.ConnectionID, &connection)

				cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
				cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "empty precompile state",
		},
		"channel not open": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				path.SetChannelOrdered()
				packetKey := host.PacketCommitmentKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				proof, proofHeight := path.EndpointA.QueryProof(packetKey)

				input, err := PackRecvPacket(MsgRecvPacket{
					Packet: Packet{
						Sequence:           big.NewInt(int64(sequence)),
						SourcePort:         path.EndpointA.ChannelConfig.PortID,
						SourceChannel:      path.EndpointA.ChannelID,
						DestinationPort:    path.EndpointB.ChannelConfig.PortID,
						DestinationChannel: path.EndpointB.ChannelID,
						Data:               data,
						TimeoutHeight: Height{
							RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
							RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
						},
						TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
					},
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(proofHeight.RevisionHeight)),
					},
					ProofCommitment: proof,
					Signer:          "Signer",
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))

				SetPort(state, path.EndpointB.ChannelConfig.PortID, contractAddress)

				err := SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				if err != nil {
					t.Error(err)
				}

				channel, _ := chainB.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				channel.State = channeltypes.CLOSED
				SetChannel(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, &channel)

				connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
				SetConnection(state, path.EndpointB.ConnectionID, &connection)

				cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
				cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "invalid channel state, channel state is not OPEN (got STATE_CLOSED)",
		},
		"capability not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				path.SetChannelOrdered()
				packetKey := host.PacketCommitmentKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				proof, proofHeight := path.EndpointA.QueryProof(packetKey)

				input, err := PackRecvPacket(MsgRecvPacket{
					Packet: Packet{
						Sequence:           big.NewInt(int64(sequence)),
						SourcePort:         path.EndpointA.ChannelConfig.PortID,
						SourceChannel:      path.EndpointA.ChannelID,
						DestinationPort:    path.EndpointB.ChannelConfig.PortID,
						DestinationChannel: path.EndpointB.ChannelID,
						Data:               data,
						TimeoutHeight: Height{
							RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
							RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
						},
						TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
					},
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(proofHeight.RevisionHeight)),
					},
					ProofCommitment: proof,
					Signer:          "Signer",
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))

				SetPort(state, path.EndpointB.ChannelConfig.PortID, contractAddress)

				channel, _ := chainB.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				SetChannel(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, &channel)

				connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
				SetConnection(state, path.EndpointB.ConnectionID, &connection)

				cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
				cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "empty precompile state, could not retrieve module from port-id",
		},
		"connection not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				path.SetChannelOrdered()
				packetKey := host.PacketCommitmentKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				proof, proofHeight := path.EndpointA.QueryProof(packetKey)

				input, err := PackRecvPacket(MsgRecvPacket{
					Packet: Packet{
						Sequence:           big.NewInt(int64(sequence)),
						SourcePort:         path.EndpointA.ChannelConfig.PortID,
						SourceChannel:      path.EndpointA.ChannelID,
						DestinationPort:    path.EndpointB.ChannelConfig.PortID,
						DestinationChannel: path.EndpointB.ChannelID,
						Data:               data,
						TimeoutHeight: Height{
							RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
							RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
						},
						TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
					},
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(proofHeight.RevisionHeight)),
					},
					ProofCommitment: proof,
					Signer:          "Signer",
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))

				SetPort(state, path.EndpointB.ChannelConfig.PortID, contractAddress)

				err := SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				if err != nil {
					t.Error(err)
				}

				channel, _ := chainB.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				SetChannel(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, &channel)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "empty precompile state",
		},
		"connection not OPEN": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				path.SetChannelOrdered()
				packetKey := host.PacketCommitmentKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				proof, proofHeight := path.EndpointA.QueryProof(packetKey)

				input, err := PackRecvPacket(MsgRecvPacket{
					Packet: Packet{
						Sequence:           big.NewInt(int64(sequence)),
						SourcePort:         path.EndpointA.ChannelConfig.PortID,
						SourceChannel:      path.EndpointA.ChannelID,
						DestinationPort:    path.EndpointB.ChannelConfig.PortID,
						DestinationChannel: path.EndpointB.ChannelID,
						Data:               data,
						TimeoutHeight: Height{
							RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
							RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
						},
						TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
					},
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(proofHeight.RevisionHeight)),
					},
					ProofCommitment: proof,
					Signer:          "Signer",
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))

				SetPort(state, path.EndpointB.ChannelConfig.PortID, contractAddress)

				err := SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				if err != nil {
					t.Error(err)
				}

				channel, _ := chainB.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				SetChannel(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, &channel)

				connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
				connection.State = connectiontypes.INIT
				SetConnection(state, path.EndpointB.ConnectionID, &connection)

				cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
				cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "invalid connection state,connection state is not OPEN (got STATE_INIT)",
		},
		"validation failed": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				path.SetChannelOrdered()
				packetKey := host.PacketCommitmentKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				proof, proofHeight := path.EndpointA.QueryProof(packetKey)

				input, err := PackRecvPacket(MsgRecvPacket{
					Packet: Packet{
						Sequence:           big.NewInt(int64(sequence)),
						SourcePort:         path.EndpointA.ChannelConfig.PortID,
						SourceChannel:      path.EndpointA.ChannelID,
						DestinationPort:    path.EndpointB.ChannelConfig.PortID,
						DestinationChannel: path.EndpointB.ChannelID,
						Data:               data,
						TimeoutHeight: Height{
							RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
							RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
						},
						TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
					},
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(proofHeight.RevisionHeight)),
					},
					ProofCommitment: proof,
					Signer:          "Signer",
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))

				SetPort(state, path.EndpointB.ChannelConfig.PortID, contractAddress)

				err := SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				if err != nil {
					t.Error(err)
				}

				channel, _ := chainB.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				SetChannel(state, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, &channel)

				connection, _ := chainB.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainB.GetContext(), path.EndpointB.ConnectionID)
				SetConnection(state, path.EndpointB.ConnectionID, &connection)

				cs, _ := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
				cStore := chainB.App.GetIBCKeeper().ClientKeeper.ClientStore(chainB.GetContext(), path.EndpointB.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					clientState.ProofSpecs = []*ics23.ProofSpec{}
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "length of specs: 0 not equal to length of proof: 2: invalid merkle proof [cosmos/ibc-go/v7@v7.2.0/modules/core/23-commitment/types/merkle.go:304], failed packet commitment verification for client (07-tendermint-0), couldn't verify counterparty packet commitment",
		},
	}
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}

func TestSendPacket(t *testing.T) {
	// destinationChannel := "DestinationChannel"
	portID := "portName"
	data := []byte("00000000000000000000000000000000000000000000000000000000000000d5")
	sourcePort := "SourcePort"
	sourceChannel := "SourceChannel"
	connectionID := "connectionID"
	clientId := "clientId"
	prefix := types.MerklePrefix{}
	sequence := uint64(1)
	heith := int64(10)
	bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
	contractAddress := common.BytesToAddress([]byte("counter"))

	tests := map[string]testutils.PrecompileTest{
		"sucsess case": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				input, err := PackSendPacket(MsgSendPacket{
					ChannelCapability: big.NewInt(0),
					SourcePort:        sourcePort,
					SourceChannel:     sourceChannel,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(1),
						RevisionHeight: big.NewInt(heith),
					},
					TimeoutTimestamp: big.NewInt(time.Now().UnixNano()),
					Data:             data,
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, portID, contractAddress)

				err := SetCapability(state, sourcePort, sourceChannel)
				if err != nil {
					t.Error(err)
				}

				setNextSequenceSend(state, sourcePort, sourceChannel, sequence)
				SetChannel(state, sourcePort, sourceChannel, &channeltypes.Channel{
					State:    channeltypes.OPEN,
					Ordering: channeltypes.UNORDERED,
					Counterparty: channeltypes.Counterparty{
						PortId:    sourcePort,
						ChannelId: sourceChannel,
					},
					ConnectionHops: []string{connectionID},
					Version:        "version",
				})
				SetConnection(state, connectionID, &connectiontypes.ConnectionEnd{
					ClientId: clientId,
					Versions: []*connectiontypes.Version{},
					State:    connectiontypes.OPEN,
					Counterparty: connectiontypes.Counterparty{
						ClientId:     clientId,
						ConnectionId: connectionID,
						Prefix:       prefix,
					},
					DelayPeriod: 2,
				})

				SetClientState(state, clientId, &tendermint.ClientState{
					ChainId: clientId,
					TrustLevel: tendermint.Fraction{
						Numerator:   1,
						Denominator: 3,
					},
					TrustingPeriod:  100,
					UnbondingPeriod: 100,
					MaxClockDrift:   100,
					FrozenHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: 0,
					},
					LatestHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: uint64(heith - 1),
					},
					ProofSpecs:                   []*ics23.ProofSpec{},
					UpgradePath:                  []string{},
					AllowUpdateAfterExpiry:       false,
					AllowUpdateAfterMisbehaviour: false,
				},
				)

				SetConsensusState(state, clientId, clienttypes.Height{
					RevisionNumber: 0,
					RevisionHeight: uint64(heith - 1),
				},
					&tendermint.ConsensusState{
						Timestamp:          time.Now(),
						Root:               types.MerkleRoot{},
						NextValidatorsHash: bytes.HexBytes{},
					},
				)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
		},
		"success: ORDERED channel": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				input, err := PackSendPacket(MsgSendPacket{
					ChannelCapability: big.NewInt(0),
					SourcePort:        sourcePort,
					SourceChannel:     sourceChannel,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(1),
						RevisionHeight: big.NewInt(heith),
					},
					TimeoutTimestamp: big.NewInt(time.Now().UnixNano()),
					Data:             data,
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, portID, contractAddress)

				err := SetCapability(state, sourcePort, sourceChannel)
				if err != nil {
					t.Error(err)
				}

				setNextSequenceSend(state, sourcePort, sourceChannel, sequence)
				SetChannel(state, sourcePort, sourceChannel, &channeltypes.Channel{
					State:    channeltypes.OPEN,
					Ordering: channeltypes.UNORDERED,
					Counterparty: channeltypes.Counterparty{
						PortId:    sourcePort,
						ChannelId: sourceChannel,
					},
					ConnectionHops: []string{connectionID},
					Version:        "version",
				})
				SetConnection(state, connectionID, &connectiontypes.ConnectionEnd{
					ClientId: clientId,
					Versions: []*connectiontypes.Version{},
					State:    connectiontypes.OPEN,
					Counterparty: connectiontypes.Counterparty{
						ClientId:     clientId,
						ConnectionId: connectionID,
						Prefix:       prefix,
					},
					DelayPeriod: 2,
				})

				SetClientState(state, clientId, &tendermint.ClientState{
					ChainId: clientId,
					TrustLevel: tendermint.Fraction{
						Numerator:   1,
						Denominator: 3,
					},
					TrustingPeriod:  100,
					UnbondingPeriod: 100,
					MaxClockDrift:   100,
					FrozenHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: 0,
					},
					LatestHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: uint64(heith - 1),
					},
					ProofSpecs:                   []*ics23.ProofSpec{},
					UpgradePath:                  []string{},
					AllowUpdateAfterExpiry:       false,
					AllowUpdateAfterMisbehaviour: false,
				},
				)

				SetConsensusState(state, clientId, clienttypes.Height{
					RevisionNumber: 0,
					RevisionHeight: uint64(heith - 1),
				},
					&tendermint.ConsensusState{
						Timestamp:          time.Now(),
						Root:               types.MerkleRoot{},
						NextValidatorsHash: bytes.HexBytes{},
					},
				)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
		},
		"failed: channel not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				input, err := PackSendPacket(MsgSendPacket{
					ChannelCapability: big.NewInt(0),
					SourcePort:        sourcePort,
					SourceChannel:     sourceChannel,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(1),
						RevisionHeight: big.NewInt(heith),
					},
					TimeoutTimestamp: big.NewInt(time.Now().UnixNano()),
					Data:             data,
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, portID, contractAddress)

				err := SetCapability(state, sourcePort, sourceChannel)
				if err != nil {
					t.Error(err)
				}

				setNextSequenceSend(state, sourcePort, sourceChannel, sequence)

				SetConnection(state, connectionID, &connectiontypes.ConnectionEnd{
					ClientId: clientId,
					Versions: []*connectiontypes.Version{},
					State:    connectiontypes.OPEN,
					Counterparty: connectiontypes.Counterparty{
						ClientId:     clientId,
						ConnectionId: connectionID,
						Prefix:       prefix,
					},
					DelayPeriod: 2,
				})

				SetClientState(state, clientId, &tendermint.ClientState{
					ChainId: clientId,
					TrustLevel: tendermint.Fraction{
						Numerator:   1,
						Denominator: 3,
					},
					TrustingPeriod:  100,
					UnbondingPeriod: 100,
					MaxClockDrift:   100,
					FrozenHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: 0,
					},
					LatestHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: uint64(heith - 1),
					},
					ProofSpecs:                   []*ics23.ProofSpec{},
					UpgradePath:                  []string{},
					AllowUpdateAfterExpiry:       false,
					AllowUpdateAfterMisbehaviour: false,
				},
				)

				SetConsensusState(state, clientId, clienttypes.Height{
					RevisionNumber: 0,
					RevisionHeight: uint64(heith - 1),
				},
					&tendermint.ConsensusState{
						Timestamp:          time.Now(),
						Root:               types.MerkleRoot{},
						NextValidatorsHash: bytes.HexBytes{},
					},
				)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "empty precompile state",
		},
		"failed: connection not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				input, err := PackSendPacket(MsgSendPacket{
					ChannelCapability: big.NewInt(0),
					SourcePort:        sourcePort,
					SourceChannel:     sourceChannel,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(1),
						RevisionHeight: big.NewInt(heith),
					},
					TimeoutTimestamp: big.NewInt(time.Now().UnixNano()),
					Data:             data,
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, portID, contractAddress)

				err := SetCapability(state, sourcePort, sourceChannel)
				if err != nil {
					t.Error(err)
				}

				setNextSequenceSend(state, sourcePort, sourceChannel, sequence)
				SetChannel(state, sourcePort, sourceChannel, &channeltypes.Channel{
					State:    channeltypes.OPEN,
					Ordering: channeltypes.UNORDERED,
					Counterparty: channeltypes.Counterparty{
						PortId:    sourcePort,
						ChannelId: sourceChannel,
					},
					ConnectionHops: []string{connectionID},
					Version:        "version",
				})

				SetClientState(state, clientId, &tendermint.ClientState{
					ChainId: clientId,
					TrustLevel: tendermint.Fraction{
						Numerator:   1,
						Denominator: 3,
					},
					TrustingPeriod:  100,
					UnbondingPeriod: 100,
					MaxClockDrift:   100,
					FrozenHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: 0,
					},
					LatestHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: uint64(heith - 1),
					},
					ProofSpecs:                   []*ics23.ProofSpec{},
					UpgradePath:                  []string{},
					AllowUpdateAfterExpiry:       false,
					AllowUpdateAfterMisbehaviour: false,
				},
				)

				SetConsensusState(state, clientId, clienttypes.Height{
					RevisionNumber: 0,
					RevisionHeight: uint64(heith - 1),
				},
					&tendermint.ConsensusState{
						Timestamp:          time.Now(),
						Root:               types.MerkleRoot{},
						NextValidatorsHash: bytes.HexBytes{},
					},
				)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "empty precompile state",
		},
		"failed: client state not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				input, err := PackSendPacket(MsgSendPacket{
					ChannelCapability: big.NewInt(0),
					SourcePort:        sourcePort,
					SourceChannel:     sourceChannel,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(1),
						RevisionHeight: big.NewInt(heith),
					},
					TimeoutTimestamp: big.NewInt(time.Now().UnixNano()),
					Data:             data,
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, portID, contractAddress)

				err := SetCapability(state, sourcePort, sourceChannel)
				if err != nil {
					t.Error(err)
				}

				setNextSequenceSend(state, sourcePort, sourceChannel, sequence)
				SetChannel(state, sourcePort, sourceChannel, &channeltypes.Channel{
					State:    channeltypes.OPEN,
					Ordering: channeltypes.UNORDERED,
					Counterparty: channeltypes.Counterparty{
						PortId:    sourcePort,
						ChannelId: sourceChannel,
					},
					ConnectionHops: []string{connectionID},
					Version:        "version",
				})
				SetConnection(state, connectionID, &connectiontypes.ConnectionEnd{
					ClientId: clientId,
					Versions: []*connectiontypes.Version{},
					State:    connectiontypes.OPEN,
					Counterparty: connectiontypes.Counterparty{
						ClientId:     clientId,
						ConnectionId: connectionID,
						Prefix:       prefix,
					},
					DelayPeriod: 2,
				})

				SetConsensusState(state, clientId, clienttypes.Height{
					RevisionNumber: 0,
					RevisionHeight: uint64(heith - 1),
				},
					&tendermint.ConsensusState{
						Timestamp:          time.Now(),
						Root:               types.MerkleRoot{},
						NextValidatorsHash: bytes.HexBytes{},
					},
				)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "error loading client state, err: empty precompile state",
		},
		"failed: consensus not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				input, err := PackSendPacket(MsgSendPacket{
					ChannelCapability: big.NewInt(0),
					SourcePort:        sourcePort,
					SourceChannel:     sourceChannel,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(1),
						RevisionHeight: big.NewInt(heith),
					},
					TimeoutTimestamp: big.NewInt(time.Now().UnixNano()),
					Data:             data,
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, portID, contractAddress)

				err := SetCapability(state, sourcePort, sourceChannel)
				if err != nil {
					t.Error(err)
				}

				setNextSequenceSend(state, sourcePort, sourceChannel, sequence)
				SetChannel(state, sourcePort, sourceChannel, &channeltypes.Channel{
					State:    channeltypes.OPEN,
					Ordering: channeltypes.UNORDERED,
					Counterparty: channeltypes.Counterparty{
						PortId:    sourcePort,
						ChannelId: sourceChannel,
					},
					ConnectionHops: []string{connectionID},
					Version:        "version",
				})
				SetConnection(state, connectionID, &connectiontypes.ConnectionEnd{
					ClientId: clientId,
					Versions: []*connectiontypes.Version{},
					State:    connectiontypes.OPEN,
					Counterparty: connectiontypes.Counterparty{
						ClientId:     clientId,
						ConnectionId: connectionID,
						Prefix:       prefix,
					},
					DelayPeriod: 2,
				})

				SetClientState(state, clientId, &tendermint.ClientState{
					ChainId: clientId,
					TrustLevel: tendermint.Fraction{
						Numerator:   1,
						Denominator: 3,
					},
					TrustingPeriod:  100,
					UnbondingPeriod: 100,
					MaxClockDrift:   100,
					FrozenHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: 0,
					},
					LatestHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: uint64(heith - 1),
					},
					ProofSpecs:                   []*ics23.ProofSpec{},
					UpgradePath:                  []string{},
					AllowUpdateAfterExpiry:       false,
					AllowUpdateAfterMisbehaviour: false,
				},
				)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "error loading consensus state, err: empty precompile state",
		},
		"failed: capability not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				input, err := PackSendPacket(MsgSendPacket{
					ChannelCapability: big.NewInt(0),
					SourcePort:        sourcePort,
					SourceChannel:     sourceChannel,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(1),
						RevisionHeight: big.NewInt(heith),
					},
					TimeoutTimestamp: big.NewInt(time.Now().UnixNano()),
					Data:             data,
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, portID, contractAddress)

				setNextSequenceSend(state, sourcePort, sourceChannel, sequence)
				SetChannel(state, sourcePort, sourceChannel, &channeltypes.Channel{
					State:    channeltypes.OPEN,
					Ordering: channeltypes.UNORDERED,
					Counterparty: channeltypes.Counterparty{
						PortId:    sourcePort,
						ChannelId: sourceChannel,
					},
					ConnectionHops: []string{connectionID},
					Version:        "version",
				})
				SetConnection(state, connectionID, &connectiontypes.ConnectionEnd{
					ClientId: clientId,
					Versions: []*connectiontypes.Version{},
					State:    connectiontypes.OPEN,
					Counterparty: connectiontypes.Counterparty{
						ClientId:     clientId,
						ConnectionId: connectionID,
						Prefix:       prefix,
					},
					DelayPeriod: 2,
				})

				SetClientState(state, clientId, &tendermint.ClientState{
					ChainId: clientId,
					TrustLevel: tendermint.Fraction{
						Numerator:   1,
						Denominator: 3,
					},
					TrustingPeriod:  100,
					UnbondingPeriod: 100,
					MaxClockDrift:   100,
					FrozenHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: 0,
					},
					LatestHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: uint64(heith - 1),
					},
					ProofSpecs:                   []*ics23.ProofSpec{},
					UpgradePath:                  []string{},
					AllowUpdateAfterExpiry:       false,
					AllowUpdateAfterMisbehaviour: false,
				},
				)

				SetConsensusState(state, clientId, clienttypes.Height{
					RevisionNumber: 0,
					RevisionHeight: uint64(heith - 1),
				},
					&tendermint.ConsensusState{
						Timestamp:          time.Now(),
						Root:               types.MerkleRoot{},
						NextValidatorsHash: bytes.HexBytes{},
					},
				)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "empty precompile state",
		},
		"failed: channel is not open": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				input, err := PackSendPacket(MsgSendPacket{
					ChannelCapability: big.NewInt(0),
					SourcePort:        sourcePort,
					SourceChannel:     sourceChannel,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(1),
						RevisionHeight: big.NewInt(heith),
					},
					TimeoutTimestamp: big.NewInt(time.Now().UnixNano()),
					Data:             data,
				})
				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, portID, contractAddress)

				err := SetCapability(state, sourcePort, sourceChannel)
				if err != nil {
					t.Error(err)
				}

				setNextSequenceSend(state, sourcePort, sourceChannel, sequence)
				SetChannel(state, sourcePort, sourceChannel, &channeltypes.Channel{
					State:    channeltypes.CLOSED,
					Ordering: channeltypes.UNORDERED,
					Counterparty: channeltypes.Counterparty{
						PortId:    sourcePort,
						ChannelId: sourceChannel,
					},
					ConnectionHops: []string{connectionID},
					Version:        "version",
				})
				SetConnection(state, connectionID, &connectiontypes.ConnectionEnd{
					ClientId: clientId,
					Versions: []*connectiontypes.Version{},
					State:    connectiontypes.OPEN,
					Counterparty: connectiontypes.Counterparty{
						ClientId:     clientId,
						ConnectionId: connectionID,
						Prefix:       prefix,
					},
					DelayPeriod: 2,
				})

				SetClientState(state, clientId, &tendermint.ClientState{
					ChainId: clientId,
					TrustLevel: tendermint.Fraction{
						Numerator:   1,
						Denominator: 3,
					},
					TrustingPeriod:  100,
					UnbondingPeriod: 100,
					MaxClockDrift:   100,
					FrozenHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: 0,
					},
					LatestHeight: clienttypes.Height{
						RevisionNumber: 0,
						RevisionHeight: uint64(heith - 1),
					},
					ProofSpecs:                   []*ics23.ProofSpec{},
					UpgradePath:                  []string{},
					AllowUpdateAfterExpiry:       false,
					AllowUpdateAfterMisbehaviour: false,
				},
				)

				SetConsensusState(state, clientId, clienttypes.Height{
					RevisionNumber: 0,
					RevisionHeight: uint64(heith - 1),
				},
					&tendermint.ConsensusState{
						Timestamp:          time.Now(),
						Root:               types.MerkleRoot{},
						NextValidatorsHash: bytes.HexBytes{},
					},
				)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "invalid channel state, channel is not OPEN (got STATE_CLOSED)",
		},
	}
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}

func TestTimeoutPacket(t *testing.T) {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	sequence := uint64(1)

	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))
	coordinator.CommitNBlocks(chainA, 2)
	coordinator.CommitNBlocks(chainB, 2)
	path := ibctesting.NewPath(chainA, chainB)

	path.SetChannelOrdered()
	coordinator.Setup(path)

	data := common.FromHex("000000000000000000000000000000000000000000000000000000000000002d")

	var (
		err         error
		proof       []byte
		proofHeight exported.Height
	)

	timeoutHeight := clienttypes.GetSelfHeight(chainB.GetContext())
	timeoutTimestamp := uint64(chainB.GetContext().BlockTime().UnixNano())

	sequence, err = path.EndpointA.SendPacket(timeoutHeight, timeoutTimestamp, data)
	if err != nil {
		t.Error(err)
	}

	packet := Packet{
		Sequence:           big.NewInt(int64(sequence)),
		SourcePort:         path.EndpointA.ChannelConfig.PortID,
		SourceChannel:      path.EndpointA.ChannelID,
		DestinationPort:    path.EndpointB.ChannelConfig.PortID,
		DestinationChannel: path.EndpointB.ChannelID,
		Data:               data,
		TimeoutHeight: Height{
			RevisionNumber: big.NewInt(int64(timeoutHeight.RevisionNumber)),
			RevisionHeight: big.NewInt(int64(timeoutHeight.RevisionHeight)),
		},
		TimeoutTimestamp: big.NewInt(int64(timeoutTimestamp)),
	}

	err = path.EndpointA.UpdateClient()
	if err != nil {
		t.Error(err)
	}

	tests := map[string]testutils.PrecompileTest{
		"sucsess case": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				orderedPacketKey := host.NextSequenceRecvKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				// unorderedPacketKey := host.PacketReceiptKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sequence)

				proof, proofHeight = path.EndpointB.QueryProof(orderedPacketKey)

				input, err := PackTimeout(MsgTimeout{
					Packet: packet,
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.GetRevisionNumber())),
						RevisionHeight: big.NewInt(int64(proofHeight.GetRevisionHeight())),
					},
					ProofUnreceived:  proof,
					Signer:           "Signer",
					NextSequenceRecv: big.NewInt(1),
				})
				if err != nil {
					t.Error(err)
				}

				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))

				SetPort(state, path.EndpointA.ChannelConfig.PortID, contractAddress)

				err := SetCapability(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if err != nil {
					t.Error(err)
				}

				SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointA.ChannelID)

				channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				SetChannel(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)

				connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				SetConnection(state, path.EndpointA.ConnectionID, &connection)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
				commitment := chainA.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				setPacketCommitment(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence, commitment)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
		},
	}
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}

func TestTimeoutOnClosePacket(t *testing.T) {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	sequence := uint64(1)

	coordinator := ibctesting.NewCoordinator(t, 2)
	chainA := coordinator.GetChain(ibctesting.GetChainID(1))
	chainB := coordinator.GetChain(ibctesting.GetChainID(2))
	coordinator.CommitNBlocks(chainA, 2)
	coordinator.CommitNBlocks(chainB, 2)
	path := ibctesting.NewPath(chainA, chainB)

	path.SetChannelOrdered()
	coordinator.Setup(path)

	data := common.FromHex("000000000000000000000000000000000000000000000000000000000000002d")

	var (
		err   error
		proof []byte
	)

	timeoutHeight := clienttypes.GetSelfHeight(chainB.GetContext())
	timeoutTimestamp := uint64(chainB.GetContext().BlockTime().UnixNano())

	sequence, err = path.EndpointA.SendPacket(timeoutHeight, timeoutTimestamp, data)
	if err != nil {
		t.Error(err)
	}

	err = path.EndpointB.SetChannelState(channeltypes.CLOSED)
	if err != nil {
		t.Error(err)
	}

	err = path.EndpointA.UpdateClient()
	if err != nil {
		t.Error(err)
	}

	packet := Packet{
		Sequence:           big.NewInt(int64(sequence)),
		SourcePort:         path.EndpointA.ChannelConfig.PortID,
		SourceChannel:      path.EndpointA.ChannelID,
		DestinationPort:    path.EndpointB.ChannelConfig.PortID,
		DestinationChannel: path.EndpointB.ChannelID,
		Data:               data,
		TimeoutHeight: Height{
			RevisionNumber: big.NewInt(int64(timeoutHeight.RevisionNumber)),
			RevisionHeight: big.NewInt(int64(timeoutHeight.RevisionHeight)),
		},
		TimeoutTimestamp: big.NewInt(int64(timeoutTimestamp)),
	}

	tests := map[string]testutils.PrecompileTest{
		"sucsess case": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				channelKey := host.ChannelKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
				orderedPacketKey := host.NextSequenceRecvKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)

				proofClosed, proofHeight := chainB.QueryProof(channelKey)

				proof, _ = chainB.QueryProof(orderedPacketKey)

				input, err := PackTimeoutOnClose(MsgTimeoutOnClose{
					Packet: packet,
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.GetRevisionNumber())),
						RevisionHeight: big.NewInt(int64(proofHeight.GetRevisionHeight())),
					},
					ProofUnreceived:  proof,
					ProofClose:       proofClosed,
					Signer:           "Signer",
					NextSequenceRecv: big.NewInt(1),
				})
				if err != nil {
					t.Error(err)
				}

				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))

				SetPort(state, path.EndpointA.ChannelConfig.PortID, contractAddress)

				err := SetCapability(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if err != nil {
					t.Error(err)
				}

				SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointA.ChannelID)

				channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				SetChannel(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)

				connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				SetConnection(state, path.EndpointA.ConnectionID, &connection)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
				commitment := chainA.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				setPacketCommitment(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence, commitment)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
		},
	}
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}

func TestAcknowledgement(t *testing.T) {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	ack := ibcmock.MockAcknowledgement

	var (
		packet      Packet
		sequence    = uint64(1)
		err         error
		chainA      *ibctesting.TestChain
		chainB      *ibctesting.TestChain
		path        *ibctesting.Path
		coordinator *ibctesting.Coordinator
	)

	tests := map[string]testutils.PrecompileTest{
		"success on ordered channel": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				packetKey := host.PacketAcknowledgementKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sequence)
				proof, proofHeight := path.EndpointB.QueryProof(packetKey)

				input, err := PackAcknowledgement(MsgAcknowledgement{
					Packet:          packet,
					Acknowledgement: ack.Acknowledgement(),
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.GetRevisionNumber())),
						RevisionHeight: big.NewInt(int64(proofHeight.GetRevisionHeight())),
					},
					ProofAcked: proof,
					Signer:     "Signer",
				})
				if err != nil {
					t.Error(err)
				}

				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				path.SetChannelOrdered()
				coordinator.Setup(path)
				data := ibctesting.MockPacketData
				sequence, err = path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, data)
				if err != nil {
					t.Error(err)
				}

				channelPacket := channeltypes.NewPacket(data, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, defaultTimeoutHeight, disabledTimeoutTimestamp)
				packet = Packet{
					Sequence:           big.NewInt(int64(sequence)),
					SourcePort:         path.EndpointA.ChannelConfig.PortID,
					SourceChannel:      path.EndpointA.ChannelID,
					DestinationPort:    path.EndpointB.ChannelConfig.PortID,
					DestinationChannel: path.EndpointB.ChannelID,
					Data:               data,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
					},
					TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
				}
				err = path.EndpointB.RecvPacket(channelPacket)
				if err != nil {
					t.Error(err)
				}

				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))
				nextSequenceAck, ok := chainA.App.GetIBCKeeper().ChannelKeeper.GetNextSequenceAck(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if !ok {
					t.Error("nextSequenceAck is not exist")
				}

				setNextSequenceAck(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, nextSequenceAck)
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, path.EndpointA.ChannelConfig.PortID, contractAddress)
				err = SetCapability(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if err != nil {
					t.Error(err)
				}

				SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointA.ChannelID)
				channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				SetChannel(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)
				connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				SetConnection(state, path.EndpointA.ConnectionID, &connection)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
				commitment := chainA.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				setPacketCommitment(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence, commitment)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
		},
		"success on unordered channel": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				packetKey := host.PacketAcknowledgementKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sequence)
				proof, proofHeight := path.EndpointB.QueryProof(packetKey)

				input, err := PackAcknowledgement(MsgAcknowledgement{
					Packet:          packet,
					Acknowledgement: ack.Acknowledgement(),
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.GetRevisionNumber())),
						RevisionHeight: big.NewInt(int64(proofHeight.GetRevisionHeight())),
					},
					ProofAcked: proof,
					Signer:     "Signer",
				})
				if err != nil {
					t.Error(err)
				}

				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				coordinator.Setup(path)
				data := ibctesting.MockPacketData
				sequence, err = path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, data)
				if err != nil {
					t.Error(err)
				}

				channelPacket := channeltypes.NewPacket(data, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, defaultTimeoutHeight, disabledTimeoutTimestamp)
				packet = Packet{
					Sequence:           big.NewInt(int64(sequence)),
					SourcePort:         path.EndpointA.ChannelConfig.PortID,
					SourceChannel:      path.EndpointA.ChannelID,
					DestinationPort:    path.EndpointB.ChannelConfig.PortID,
					DestinationChannel: path.EndpointB.ChannelID,
					Data:               data,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
					},
					TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
				}
				err = path.EndpointB.RecvPacket(channelPacket)
				if err != nil {
					t.Error(err)
				}

				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))
				nextSequenceAck, ok := chainA.App.GetIBCKeeper().ChannelKeeper.GetNextSequenceAck(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if !ok {
					t.Error("nextSequenceAck is not exist")
				}

				setNextSequenceAck(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, nextSequenceAck)
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, path.EndpointA.ChannelConfig.PortID, contractAddress)
				err = SetCapability(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if err != nil {
					t.Error(err)
				}

				SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointA.ChannelID)
				channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				SetChannel(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)
				connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				SetConnection(state, path.EndpointA.ConnectionID, &connection)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
				commitment := chainA.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				setPacketCommitment(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence, commitment)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
		},
		"channel not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				packetKey := host.PacketAcknowledgementKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sequence)
				proof, proofHeight := path.EndpointB.QueryProof(packetKey)

				input, err := PackAcknowledgement(MsgAcknowledgement{
					Packet:          packet,
					Acknowledgement: ack.Acknowledgement(),
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.GetRevisionNumber())),
						RevisionHeight: big.NewInt(int64(proofHeight.GetRevisionHeight())),
					},
					ProofAcked: proof,
					Signer:     "Signer",
				})
				if err != nil {
					t.Error(err)
				}

				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				path.SetChannelOrdered()
				coordinator.Setup(path)
				data := ibctesting.MockPacketData
				sequence, err = path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, data)
				if err != nil {
					t.Error(err)
				}

				channelPacket := channeltypes.NewPacket(data, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, defaultTimeoutHeight, disabledTimeoutTimestamp)
				packet = Packet{
					Sequence:           big.NewInt(int64(sequence)),
					SourcePort:         path.EndpointA.ChannelConfig.PortID,
					SourceChannel:      path.EndpointA.ChannelID,
					DestinationPort:    path.EndpointB.ChannelConfig.PortID,
					DestinationChannel: path.EndpointB.ChannelID,
					Data:               data,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
					},
					TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
				}
				err = path.EndpointB.RecvPacket(channelPacket)
				if err != nil {
					t.Error(err)
				}

				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))
				nextSequenceAck, ok := chainA.App.GetIBCKeeper().ChannelKeeper.GetNextSequenceAck(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if !ok {
					t.Error("nextSequenceAck is not exist")
				}

				setNextSequenceAck(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, nextSequenceAck)
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, path.EndpointA.ChannelConfig.PortID, contractAddress)
				err = SetCapability(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if err != nil {
					t.Error(err)
				}

				SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointA.ChannelID)
				connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				SetConnection(state, path.EndpointA.ConnectionID, &connection)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
				commitment := chainA.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				setPacketCommitment(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence, commitment)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "empty precompile state",
		},
		"channel not open": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				packetKey := host.PacketAcknowledgementKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sequence)
				proof, proofHeight := path.EndpointB.QueryProof(packetKey)

				input, err := PackAcknowledgement(MsgAcknowledgement{
					Packet:          packet,
					Acknowledgement: ack.Acknowledgement(),
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.GetRevisionNumber())),
						RevisionHeight: big.NewInt(int64(proofHeight.GetRevisionHeight())),
					},
					ProofAcked: proof,
					Signer:     "Signer",
				})
				if err != nil {
					t.Error(err)
				}

				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				path.SetChannelOrdered()
				coordinator.Setup(path)
				data := ibctesting.MockPacketData
				sequence, err = path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, data)
				if err != nil {
					t.Error(err)
				}

				channelPacket := channeltypes.NewPacket(data, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, defaultTimeoutHeight, disabledTimeoutTimestamp)
				packet = Packet{
					Sequence:           big.NewInt(int64(sequence)),
					SourcePort:         path.EndpointA.ChannelConfig.PortID,
					SourceChannel:      path.EndpointA.ChannelID,
					DestinationPort:    path.EndpointB.ChannelConfig.PortID,
					DestinationChannel: path.EndpointB.ChannelID,
					Data:               data,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
					},
					TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
				}
				err = path.EndpointB.RecvPacket(channelPacket)
				if err != nil {
					t.Error(err)
				}

				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))
				nextSequenceAck, ok := chainA.App.GetIBCKeeper().ChannelKeeper.GetNextSequenceAck(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if !ok {
					t.Error("nextSequenceAck is not exist")
				}

				setNextSequenceAck(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, nextSequenceAck)
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, path.EndpointA.ChannelConfig.PortID, contractAddress)
				err = SetCapability(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if err != nil {
					t.Error(err)
				}

				SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointA.ChannelID)
				channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				channel.State = channeltypes.CLOSED
				SetChannel(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)
				connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				SetConnection(state, path.EndpointA.ConnectionID, &connection)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
				commitment := chainA.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				setPacketCommitment(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence, commitment)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "invalid channel state, channel state is not OPEN (got STATE_CLOSED)",
		},
		"capability not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				packetKey := host.PacketAcknowledgementKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sequence)
				proof, proofHeight := path.EndpointB.QueryProof(packetKey)

				input, err := PackAcknowledgement(MsgAcknowledgement{
					Packet:          packet,
					Acknowledgement: ack.Acknowledgement(),
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.GetRevisionNumber())),
						RevisionHeight: big.NewInt(int64(proofHeight.GetRevisionHeight())),
					},
					ProofAcked: proof,
					Signer:     "Signer",
				})
				if err != nil {
					t.Error(err)
				}

				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				path.SetChannelOrdered()
				coordinator.Setup(path)
				data := ibctesting.MockPacketData
				sequence, err = path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, data)
				if err != nil {
					t.Error(err)
				}

				channelPacket := channeltypes.NewPacket(data, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, defaultTimeoutHeight, disabledTimeoutTimestamp)
				packet = Packet{
					Sequence:           big.NewInt(int64(sequence)),
					SourcePort:         path.EndpointA.ChannelConfig.PortID,
					SourceChannel:      path.EndpointA.ChannelID,
					DestinationPort:    path.EndpointB.ChannelConfig.PortID,
					DestinationChannel: path.EndpointB.ChannelID,
					Data:               data,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
					},
					TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
				}
				err = path.EndpointB.RecvPacket(channelPacket)
				if err != nil {
					t.Error(err)
				}

				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))
				nextSequenceAck, ok := chainA.App.GetIBCKeeper().ChannelKeeper.GetNextSequenceAck(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if !ok {
					t.Error("nextSequenceAck is not exist")
				}

				setNextSequenceAck(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, nextSequenceAck)
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, path.EndpointA.ChannelConfig.PortID, contractAddress)
				channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				SetChannel(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)
				connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				SetConnection(state, path.EndpointA.ConnectionID, &connection)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
				commitment := chainA.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				setPacketCommitment(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence, commitment)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "empty precompile state, could not retrieve module from port-id",
		},
		"connection not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				packetKey := host.PacketAcknowledgementKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sequence)
				proof, proofHeight := path.EndpointB.QueryProof(packetKey)

				input, err := PackAcknowledgement(MsgAcknowledgement{
					Packet:          packet,
					Acknowledgement: ack.Acknowledgement(),
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.GetRevisionNumber())),
						RevisionHeight: big.NewInt(int64(proofHeight.GetRevisionHeight())),
					},
					ProofAcked: proof,
					Signer:     "Signer",
				})
				if err != nil {
					t.Error(err)
				}

				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				path.SetChannelOrdered()
				coordinator.Setup(path)
				data := ibctesting.MockPacketData
				sequence, err = path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, data)
				if err != nil {
					t.Error(err)
				}

				channelPacket := channeltypes.NewPacket(data, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, defaultTimeoutHeight, disabledTimeoutTimestamp)
				packet = Packet{
					Sequence:           big.NewInt(int64(sequence)),
					SourcePort:         path.EndpointA.ChannelConfig.PortID,
					SourceChannel:      path.EndpointA.ChannelID,
					DestinationPort:    path.EndpointB.ChannelConfig.PortID,
					DestinationChannel: path.EndpointB.ChannelID,
					Data:               data,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
					},
					TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
				}
				err = path.EndpointB.RecvPacket(channelPacket)
				if err != nil {
					t.Error(err)
				}

				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))
				nextSequenceAck, ok := chainA.App.GetIBCKeeper().ChannelKeeper.GetNextSequenceAck(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if !ok {
					t.Error("nextSequenceAck is not exist")
				}
				setNextSequenceAck(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, nextSequenceAck)

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, path.EndpointA.ChannelConfig.PortID, contractAddress)
				err = SetCapability(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if err != nil {
					t.Error(err)
				}

				SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointA.ChannelID)
				channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				SetChannel(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)

				commitment := chainA.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				setPacketCommitment(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence, commitment)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "empty precompile state",
		},
		"connection not OPEN": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				packetKey := host.PacketAcknowledgementKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sequence)
				proof, proofHeight := path.EndpointB.QueryProof(packetKey)

				input, err := PackAcknowledgement(MsgAcknowledgement{
					Packet:          packet,
					Acknowledgement: ack.Acknowledgement(),
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.GetRevisionNumber())),
						RevisionHeight: big.NewInt(int64(proofHeight.GetRevisionHeight())),
					},
					ProofAcked: proof,
					Signer:     "Signer",
				})
				if err != nil {
					t.Error(err)
				}

				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				path.SetChannelOrdered()
				coordinator.Setup(path)
				data := ibctesting.MockPacketData
				sequence, err = path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, data)
				if err != nil {
					t.Error(err)
				}

				channelPacket := channeltypes.NewPacket(data, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, defaultTimeoutHeight, disabledTimeoutTimestamp)
				packet = Packet{
					Sequence:           big.NewInt(int64(sequence)),
					SourcePort:         path.EndpointA.ChannelConfig.PortID,
					SourceChannel:      path.EndpointA.ChannelID,
					DestinationPort:    path.EndpointB.ChannelConfig.PortID,
					DestinationChannel: path.EndpointB.ChannelID,
					Data:               data,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
					},
					TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
				}
				err = path.EndpointB.RecvPacket(channelPacket)
				if err != nil {
					t.Error(err)
				}

				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))
				nextSequenceAck, ok := chainA.App.GetIBCKeeper().ChannelKeeper.GetNextSequenceAck(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if !ok {
					t.Error("nextSequenceAck is not exist")
				}

				setNextSequenceAck(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, nextSequenceAck)
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, path.EndpointA.ChannelConfig.PortID, contractAddress)
				err = SetCapability(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if err != nil {
					t.Error(err)
				}

				SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointA.ChannelID)
				channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				SetChannel(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)
				connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				connection.State = connectiontypes.INIT
				SetConnection(state, path.EndpointA.ConnectionID, &connection)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
				commitment := chainA.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				setPacketCommitment(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence, commitment)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "invalid connection state, connection state is not OPEN (got STATE_INIT)",
		},
		"packet commitment bytes do not match": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				packetKey := host.PacketAcknowledgementKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sequence)
				proof, proofHeight := path.EndpointB.QueryProof(packetKey)

				input, err := PackAcknowledgement(MsgAcknowledgement{
					Packet:          packet,
					Acknowledgement: ack.Acknowledgement(),
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.GetRevisionNumber())),
						RevisionHeight: big.NewInt(int64(proofHeight.GetRevisionHeight())),
					},
					ProofAcked: proof,
					Signer:     "Signer",
				})
				if err != nil {
					t.Error(err)
				}

				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				path.SetChannelOrdered()
				coordinator.Setup(path)
				data := []byte("some data")
				sequence, err = path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, data)
				if err != nil {
					t.Error(err)
				}

				channelPacket := channeltypes.NewPacket(data, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, defaultTimeoutHeight, disabledTimeoutTimestamp)
				packet = Packet{
					Sequence:           big.NewInt(int64(sequence)),
					SourcePort:         path.EndpointA.ChannelConfig.PortID,
					SourceChannel:      path.EndpointA.ChannelID,
					DestinationPort:    path.EndpointB.ChannelConfig.PortID,
					DestinationChannel: path.EndpointB.ChannelID,
					Data:               data,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
					},
					TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
				}
				err = path.EndpointB.RecvPacket(channelPacket)
				if err != nil {
					t.Error(err)
				}

				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))
				nextSequenceAck, ok := chainA.App.GetIBCKeeper().ChannelKeeper.GetNextSequenceAck(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if !ok {
					t.Error("nextSequenceAck is not exist")
				}

				setNextSequenceAck(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, nextSequenceAck)
				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, path.EndpointA.ChannelConfig.PortID, contractAddress)
				err = SetCapability(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if err != nil {
					t.Error(err)
				}

				SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointA.ChannelID)
				channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				SetChannel(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)
				connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				SetConnection(state, path.EndpointA.ConnectionID, &connection)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "empty precompile state",
		},
		"next ack sequence not found": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				packetKey := host.PacketAcknowledgementKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sequence)
				proof, proofHeight := path.EndpointB.QueryProof(packetKey)

				input, err := PackAcknowledgement(MsgAcknowledgement{
					Packet:          packet,
					Acknowledgement: ack.Acknowledgement(),
					ProofHeight: Height{
						RevisionNumber: big.NewInt(int64(proofHeight.GetRevisionNumber())),
						RevisionHeight: big.NewInt(int64(proofHeight.GetRevisionHeight())),
					},
					ProofAcked: proof,
					Signer:     "Signer",
				})
				if err != nil {
					t.Error(err)
				}

				require.NoError(t, err)
				return input
			},
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				path.SetChannelOrdered()
				coordinator.Setup(path)
				data := ibctesting.MockPacketData
				sequence, err = path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, data)
				if err != nil {
					t.Error(err)
				}

				channelPacket := channeltypes.NewPacket(data, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, defaultTimeoutHeight, disabledTimeoutTimestamp)
				packet = Packet{
					Sequence:           big.NewInt(int64(sequence)),
					SourcePort:         path.EndpointA.ChannelConfig.PortID,
					SourceChannel:      path.EndpointA.ChannelID,
					DestinationPort:    path.EndpointB.ChannelConfig.PortID,
					DestinationChannel: path.EndpointB.ChannelID,
					Data:               data,
					TimeoutHeight: Height{
						RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
						RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
					},
					TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
				}
				err = path.EndpointB.RecvPacket(channelPacket)
				if err != nil {
					t.Error(err)
				}

				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

				state.CreateAccount(contractAddress)
				state.SetCode(contractAddress, hexutil.MustDecode(bytecode))
				SetPort(state, path.EndpointA.ChannelConfig.PortID, contractAddress)
				err = SetCapability(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				if err != nil {
					t.Error(err)
				}

				SetCapability(state, path.EndpointB.ChannelConfig.PortID, path.EndpointA.ChannelID)
				channel, _ := chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				SetChannel(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, &channel)
				connection, _ := chainA.App.GetIBCKeeper().ConnectionKeeper.GetConnection(chainA.GetContext(), path.EndpointA.ConnectionID)
				SetConnection(state, path.EndpointA.ConnectionID, &connection)

				cs, _ := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
				cStore := chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(chainA.GetContext(), path.EndpointA.ClientID)
				if cs != nil {
					clientState := cs.(*ibctm.ClientState)
					bz := cStore.Get([]byte(fmt.Sprintf("consensusStates/%s", cs.GetLatestHeight())))
					consensusState := clienttypes.MustUnmarshalConsensusState(marshaler, bz)
					SetClientState(state, connection.GetClientID(), clientState)
					SetConsensusState(state, connection.GetClientID(), clientState.GetLatestHeight(), consensusState.(*ibctm.ConsensusState))
				}
				commitment := chainA.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
				setPacketCommitment(state, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence, commitment)
			},
			SuppliedGas: BindPortGasCost,
			ReadOnly:    false,
			ExpectedErr: "empty precompile state",
		},
	}
	// Run tests.
	for name, test := range tests {
		coordinator = ibctesting.NewCoordinator(t, 2)
		chainA = coordinator.GetChain(ibctesting.GetChainID(1))
		chainB = coordinator.GetChain(ibctesting.GetChainID(2))
		coordinator.CommitNBlocks(chainA, 2)
		coordinator.CommitNBlocks(chainB, 2)
		path = ibctesting.NewPath(chainA, chainB)

		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}
