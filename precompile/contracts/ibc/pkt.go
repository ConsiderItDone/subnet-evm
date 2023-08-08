package ibc

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
)

type Height struct {
	RevisionNumber uint64
	RevisionHeight uint64
}

type Packet struct {
	Sequence           uint64
	SourcePort         string
	SourceChannel      string
	DestinationPort    string
	DestinationChannel string
	Data               []byte
	TimeoutHeight      Height
	TimeoutTimestamp   uint64
}

type OnRecvPacketInput struct {
	Packet  Packet
	Relayer []byte
}

type SendPacketInput struct {
	ChannelCapability uint64
	SourcePort        string
	SourceChannel     string
	TimeoutHeight     Height
	TimeoutTimestamp  uint64
	Data              []byte
}

// PackOnRecvPacket packs [inputStruct] of type OnRecvPacketInput into the appropriate arguments for OnRecvPacket.
func PackOnRecvPacket(inputStruct OnRecvPacketInput) ([]byte, error) {
	return IBCABI.Pack("OnRecvPacket", inputStruct.Packet, inputStruct.Relayer)
}

// UnpackSendPacketInput attempts to unpack [input] as SendPacketInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackSendPacketInput(input []byte) (SendPacketInput, error) {
	inputStruct := SendPacketInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "SendPacket", input)

	return inputStruct, err
}

// PackSendPacket packs [inputStruct] of type SendPacketInput into the appropriate arguments for SendPacket.
func PackSendPacket(inputStruct SendPacketInput) ([]byte, error) {
	return IBCABI.Pack("SendPacket", inputStruct.ChannelCapability, inputStruct.SourcePort, inputStruct.SourceChannel, inputStruct.TimeoutHeight, inputStruct.TimeoutTimestamp, inputStruct.Data)
}

func sendPacket(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, SendPacketGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the SendPacketInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackSendPacketInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	_ = inputStruct // CUSTOM CODE OPERATES ON INPUT
	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}
