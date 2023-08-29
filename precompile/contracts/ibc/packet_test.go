package ibc

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/cometbft/cometbft/libs/bytes"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	tendermint "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ics23 "github.com/cosmos/ics23/go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func TestRecvPacket(t *testing.T) {
	destinationChannel := "DestinationChannel"
	portID := "portName"
	data := []byte("00000000000000000000000000000000000000000000000000000000000000d5")
	sourcePort := "SourcePort"
	sourceChannel := "SourceChannel"
	connectionID := "connectionID"
	clientId := "clientId"
	prefix := types.MerklePrefix{}
	sequence := int64(1)

	tests := map[string]testutils.PrecompileTest{
		"sucsess case": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {

				input, err := PackRecvPacket(MsgRecvPacket{
					Packet: Packet{
						Sequence:           big.NewInt(sequence),
						SourcePort:         sourcePort,
						SourceChannel:      sourceChannel,
						DestinationPort:    portID,
						DestinationChannel: destinationChannel,
						Data:               data,
						TimeoutHeight: Height{
							RevisionNumber: big.NewInt(1),
							RevisionHeight: big.NewInt(1),
						},
						TimeoutTimestamp: big.NewInt(time.Now().UnixNano()),
					},
					ProofHeight: Height{
						RevisionNumber: big.NewInt(1),
						RevisionHeight: big.NewInt(1),
					},
					ProofCommitment: []byte("Proof"),
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

				SetPort(state, portID, contractAddress)

				err := SetCapability(state, portID, destinationChannel)
				if err != nil {
					t.Error(err)
				}

				SetChannel(state, portID, destinationChannel, &channeltypes.Channel{
					State:    0,
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
						RevisionHeight: uint64(0),
					},
					ProofSpecs:                   []*ics23.ProofSpec{},
					UpgradePath:                  []string{},
					AllowUpdateAfterExpiry:       false,
					AllowUpdateAfterMisbehaviour: false,
				},
				)
			},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				contractAddress := common.BytesToAddress([]byte("counter"))
				// common.Hash{} -> 0x0000000000000000000000000000000000000000000000000000
				newState := state.GetState(contractAddress, common.Hash{})
				// newState = 0x 000000 10 (dec)
				if !reflect.DeepEqual(newState.Bytes(), data) {
					t.Error("return value of test contract not equal data")
				}
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
				bytecode := "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80631ea0324b1461003b57806361bc221a14610057575b600080fd5b6100556004803603810190610050919061039b565b610075565b005b61005f6100dc565b60405161006c9190610416565b60405180910390f35b60008260a0015180602001905181019061008f9190610372565b9050806000808282829054906101000a900460070b6100ae91906104c2565b92506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550505050565b60008054906101000a900460070b81565b60006101006100fb84610462565b610431565b90508281526020810184848401111561011857600080fd5b61012384828561055f565b509392505050565b600061013e61013984610492565b610431565b90508281526020810184848401111561015657600080fd5b61016184828561055f565b509392505050565b600082601f83011261017a57600080fd5b813561018a8482602086016100ed565b91505092915050565b6000815190506101a2816105cc565b92915050565b600082601f8301126101b957600080fd5b81356101c984826020860161012b565b91505092915050565b6000604082840312156101e457600080fd5b6101ee6040610431565b905060006101fe8482850161035d565b60008301525060206102128482850161035d565b60208301525092915050565b6000610120828403121561023157600080fd5b61023c610100610431565b9050600061024c8482850161035d565b600083015250602082013567ffffffffffffffff81111561026c57600080fd5b610278848285016101a8565b602083015250604082013567ffffffffffffffff81111561029857600080fd5b6102a4848285016101a8565b604083015250606082013567ffffffffffffffff8111156102c457600080fd5b6102d0848285016101a8565b606083015250608082013567ffffffffffffffff8111156102f057600080fd5b6102fc848285016101a8565b60808301525060a082013567ffffffffffffffff81111561031c57600080fd5b61032884828501610169565b60a08301525060c061033c848285016101d2565b60c0830152506101006103518482850161035d565b60e08301525092915050565b60008135905061036c816105e3565b92915050565b60006020828403121561038457600080fd5b600061039284828501610193565b91505092915050565b600080604083850312156103ae57600080fd5b600083013567ffffffffffffffff8111156103c857600080fd5b6103d48582860161021e565b925050602083013567ffffffffffffffff8111156103f157600080fd5b6103fd85828601610169565b9150509250929050565b6104108161053e565b82525050565b600060208201905061042b6000830184610407565b92915050565b6000604051905081810181811067ffffffffffffffff821117156104585761045761059d565b5b8060405250919050565b600067ffffffffffffffff82111561047d5761047c61059d565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156104ad576104ac61059d565b5b601f19601f8301169050602081019050919050565b60006104cd8261053e565b91506104d88361053e565b925081677fffffffffffffff038313600083121516156104fb576104fa61056e565b5b817fffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000383126000831216156105335761053261056e565b5b828201905092915050565b60008160070b9050919050565b600067ffffffffffffffff82169050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105d58161053e565b81146105e057600080fd5b50565b6105ec8161054b565b81146105f757600080fd5b5056fea26469706673582212209e026e85ea91437372f076cae491d863bc86745ef88c421d038f2133152a6f4d64736f6c63430008000033"
				contractAddress := common.BytesToAddress([]byte("counter"))

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
	}
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}
