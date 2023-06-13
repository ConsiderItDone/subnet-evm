package ibc

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	hosttypes "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
)

var (
	sendPacketSignature    = contract.CalculateFunctionSelector("sendPacket(uint64,bytes,uint64,bytes,bytes,uint64,bytes)")
	receivePacketSignature = contract.CalculateFunctionSelector("receivePacket(uint64,bytes,uint64,bytes,uint64,bytes)")
)

func SendPacket(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
			input
			8 byte             - sourcePortLen
			clientIDLen byte   - sourcePortByte
			8 byte             - sourceChannelLen
			clientMessageLen   - sourceChannelByte
		    8 byte             - timeoutHeightLen
		    proofHeightbyte    - clienttypes.Height
		    8 byte             - timeoutTimestamp
			8 byte             - DataLen
		    data               - []byte

	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, sendPacketGas); err != nil {
		return nil, 0, fmt.Errorf("rrror DeductGas err: %w", err)
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("non-enabled cannot call upgradeClient: %s", caller)
	}

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	// sourcePort
	carriage := uint64(0)
	sourcePortLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	sourcePort := string(getData(input, carriage, sourcePortLen))
	carriage = carriage + sourcePortLen

	// sourceChannel
	sourceChannelLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	sourceChannel := string(getData(input, carriage, sourceChannelLen))
	carriage = carriage + sourcePortLen

	// timeoutHeight
	timeoutHeightLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	timeoutHeightByte := getData(input, carriage, timeoutHeightLen)

	timeoutHeight := &clienttypes.Height{}
	err = marshaler.Unmarshal(timeoutHeightByte, timeoutHeight)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	// timeoutTimestamp
	timeoutTimestamp := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8

	// data
	DataLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	data := getData(input, carriage, DataLen)
	carriage = carriage + DataLen

	channel, err := getChannelState(marshaler, string(hosttypes.ChannelKey(sourcePort, sourceChannel)), accessibleState)
	if err != nil {
		return nil, 0, err
	}

	if channel.State != types.OPEN {
		return nil, 0, fmt.Errorf("%s: channel is not OPEN (got %s)", types.ErrInvalidChannelState, channel.State.String())
	}

	sequence, found := GetNextSequenceSend(accessibleState, sourcePort, sourceChannel)
	if !found {
		return nil, 0, fmt.Errorf("%s: source port: %s, source channel: %s", types.ErrSequenceSendNotFound, sourcePort, sourceChannel)
	}

	// construct packet from given fields and channel state
	packet := types.NewPacket(data, sequence, sourcePort, sourceChannel,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId, *timeoutHeight, timeoutTimestamp)

	if err := packet.ValidateBasic(); err != nil {
		return nil, 0, fmt.Errorf("%s: constructed packet failed basic validation", err)
	}

	commitment := types.CommitPacket(marshaler, packet)

	SetNextSequenceSend(accessibleState, sourcePort, sourceChannel, sequence+1)
	SetPacketCommitment(accessibleState, sourcePort, sourceChannel, packet.GetSequence(), commitment)

	var result []byte
	binary.BigEndian.PutUint64(result, sequence)

	return result, remainingGas, nil
}

func ReceivePacket(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
			input
			8 byte             - packetLen
			clientIDLen byte   - packetByte
			8 byte             - proofLen
			clientMessageLen   - proofByte
		    8 byte             - proofHeightLen
		    proofHeightbyte    - clienttypes.Height

	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, sendPacketGas); err != nil {
		return nil, 0, fmt.Errorf("rrror DeductGas err: %w", err)
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("non-enabled cannot call upgradeClient: %s", caller)
	}

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	// sourcePort
	carriage := uint64(0)
	sourcePortLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	sourcePort := string(getData(input, carriage, sourcePortLen))
	carriage = carriage + sourcePortLen

	// sourceChannel
	sourceChannelLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	sourceChannel := string(getData(input, carriage, sourceChannelLen))
	carriage = carriage + sourcePortLen

	// timeoutHeight
	timeoutHeightLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	timeoutHeightByte := getData(input, carriage, timeoutHeightLen)

	timeoutHeight := &clienttypes.Height{}
	err = marshaler.Unmarshal(timeoutHeightByte, timeoutHeight)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	// timeoutTimestamp
	timeoutTimestamp := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8

	// data
	DataLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	data := getData(input, carriage, DataLen)
	carriage = carriage + DataLen

	channel, err := getChannelState(marshaler, string(hosttypes.ChannelKey(sourcePort, sourceChannel)), accessibleState)
	if err != nil {
		return nil, 0, err
	}

	if channel.State != types.OPEN {
		return nil, 0, fmt.Errorf("%s: channel is not OPEN (got %s)", types.ErrInvalidChannelState, channel.State.String())
	}

	sequence, found := GetNextSequenceSend(accessibleState, sourcePort, sourceChannel)
	if !found {
		return nil, 0, fmt.Errorf("%s: source port: %s, source channel: %s", types.ErrSequenceSendNotFound, sourcePort, sourceChannel)
	}

	// construct packet from given fields and channel state
	packet := types.NewPacket(data, sequence, sourcePort, sourceChannel,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId, *timeoutHeight, timeoutTimestamp)

	if err := packet.ValidateBasic(); err != nil {
		return nil, 0, fmt.Errorf("%s: constructed packet failed basic validation", err)
	}

	commitment := types.CommitPacket(marshaler, packet)

	SetNextSequenceSend(accessibleState, sourcePort, sourceChannel, sequence+1)
	SetPacketCommitment(accessibleState, sourcePort, sourceChannel, packet.GetSequence(), commitment)

	var result []byte
	binary.BigEndian.PutUint64(result, sequence)

	return result, remainingGas, nil
}
