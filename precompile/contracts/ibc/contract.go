// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package ibc

import (
	_ "embed"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/precompile/contract"
)

const (
	CreateClientGasCost  uint64 = 1
	UpdateClientGasCost  uint64 = 1
	UpgradeClientGasCost uint64 = 1

	ConnOpenAckGasCost     uint64 = 1
	ConnOpenConfirmGasCost uint64 = 1
	ConnOpenInitGasCost    uint64 = 1
	ConnOpenTryGasCost     uint64 = 1

	ChanOpenInitGasCost        uint64 = 1
	ChanOpenTryGasCost         uint64 = 1
	ChannelCloseConfirmGasCost uint64 = 1
	ChannelCloseInitGasCost    uint64 = 1
	ChannelOpenAckGasCost      uint64 = 1
	ChannelOpenConfirmGasCost  uint64 = 1

	BindPortGasCost uint64 = 1

	RecvPacketGasCost      uint64 = 1
	SendPacketGasCost      uint64 = 1
	AcknowledgementGasCost uint64 = 1
	TimeoutGasCost         uint64 = 1
	TimeoutOnCloseGasCost  uint64 = 1

	QueryClientStateGasCost    uint64 = 1
	QueryConsensusStateGasCost uint64 = 1
	QueryConnectionGasCost     uint64 = 1
	QueryChannelGasCost        uint64 = 1

	QueryPacketCommitmentGasCost      uint64 = 1
	QueryPacketAcknowledgementGasCost uint64 = 1
)

// CUSTOM CODE STARTS HERE
// Reference imports to suppress errors from unused imports. This code and any unnecessary imports can be removed.
var (
	_ = abi.JSON
	_ = errors.New
	_ = big.NewInt
)

// Singleton StatefulPrecompiledContract and signatures.
var (

	// IBCRawABI contains the raw ABI of IBC contract.
	//go:embed contract.abi
	IBCRawABI string

	IBCABI                                    = contract.ParseABI(IBCRawABI)
	IBCPrecompile                             = createIBCPrecompile()
	GeneratedClientIdentifier                 = IBCABI.Events["ClientCreated"]
	GeneratedConnectionIdentifier             = IBCABI.Events["ConnectionCreated"]
	GeneratedPacketSentIdentifier             = IBCABI.Events["PacketSent"]
	GeneratedPacketReceivedIdentifier         = IBCABI.Events["PacketReceived"]
	GeneratedTimeoutPacketIdentifier          = IBCABI.Events["TimeoutPacket"]
	GeneratedAcknowledgePacketIdentifier      = IBCABI.Events["AcknowledgePacket"]
	GeneratedAcknowledgementWrittenIdentifier = IBCABI.Events["AcknowledgementWritten"]
	GeneratedAcknowledgementErrorIdentifier   = IBCABI.Events["AcknowledgementError"]

	ClientSequenceSlot     = common.BytesToHash([]byte("client-sequence"))
	ConnectionSequenceSlot = common.BytesToHash([]byte("connection-sequence"))
	ChannelSequenceSlot    = common.BytesToHash([]byte("channel-sequence"))

	ErrWrongClientType = errors.New("wrong client type. Only Tendermint supported")
)

// createIBCPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.
func createIBCPrecompile() contract.StatefulPrecompiledContract {
	var functions []*contract.StatefulPrecompileFunction

	abiFunctionMap := map[string]contract.RunStatefulPrecompileFunc{
		//"bindPort":            bindPort,
		"chanOpenInit":        chanOpenInit,
		"chanOpenTry":         chanOpenTry,
		"channelCloseConfirm": channelCloseConfirm,
		"channelCloseInit":    channelCloseInit,
		"channelOpenAck":      channelOpenAck,
		"channelOpenConfirm":  channelOpenConfirm,
		"connOpenAck":         connOpenAck,
		"connOpenConfirm":     connOpenConfirm,
		"connOpenInit":        connOpenInit,
		"connOpenTry":         connOpenTry,
		"createClient":        createClient,
		"updateClient":        updateClient,
		"upgradeClient":       upgradeClient,
		"recvPacket":          recvPacket,
		"sendPacket":          sendPacket,
		"timeout":             timeout,
		"timeoutOnClose":      timeoutOnClose,
		"acknowledgement":     acknowledgement,

		"queryClientState":           queryClientState,
		"queryConsensusState":        queryConsensusState,
		"queryConnection":            queryConnection,
		"queryChannel":               queryChannel,
		"queryPacketCommitment":      queryPacketCommitment,
		"queryPacketAcknowledgement": queryPacketAcknowledgement,
	}

	for name, function := range abiFunctionMap {
		method, ok := IBCABI.Methods[name]
		if !ok {
			panic(fmt.Errorf("given method (%s) does not exist in the ABI", name))
		}
		functions = append(functions, contract.NewStatefulPrecompileFunction(method.ID, function))
	}
	// Construct the contract with no fallback function.
	statefulContract, err := contract.NewStatefulPrecompileContract(nil, functions)
	if err != nil {
		panic(err)
	}
	return statefulContract
}
