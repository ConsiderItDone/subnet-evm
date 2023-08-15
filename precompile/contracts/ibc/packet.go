package ibc

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/accounts/abi"
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

type MsgAcknowledgement struct {
	Packet          Packet
	Acknowledgement []byte
	ProofAcked      []byte
	ProofHeight     Height
	Signer          string
}

type MsgRecvPacket struct {
	Packet          Packet
	ProofCommitment []byte
	ProofHeight     Height
	Signer          string
}

type MsgTimeout struct {
	Packet           Packet
	ProofUnreceived  []byte
	ProofHeight      Height
	NextSequenceRecv uint64
	Signer           string
}

type MsgTimeoutOnClose struct {
	Packet           Packet
	ProofUnreceived  []byte
	ProofClose       []byte
	ProofHeight      Height
	NextSequenceRecv uint64
	Signer           string
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

type RecvPacketInput struct {
	Packet          Packet
	ProofCommitment []byte
	ProofHeight     Height
	Signer          string
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

	if err := _sendPacket(&callOpts[SendPacketInput]{
		accessibleState: accessibleState,
		caller:          caller,
		addr:            addr,
		suppliedGas:     suppliedGas,
		readOnly:        readOnly,
		args:            inputStruct,
	}); err != nil {
		return nil, remainingGas, err
	}

	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackRecvPacketInput attempts to unpack [input] into the IIBCMsgRecvPacket type argument
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackRecvPacketInput(input []byte) (*MsgRecvPacket, error) {
	res, err := IBCABI.UnpackInput("RecvPacket", input)
	if err != nil {
		return nil, err
	}
	unpacked := abi.ConvertType(res[0], new(MsgRecvPacket)).(*MsgRecvPacket)
	return unpacked, nil
}

// PackRecvPacket packs [message] of type IIBCMsgRecvPacket into the appropriate arguments for RecvPacket.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackRecvPacket(message MsgRecvPacket) ([]byte, error) {
	return IBCABI.Pack("RecvPacket", message)
}

func recvPacket(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, RecvPacketGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the RecvPacketInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackRecvPacketInput(input)
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

// UnpackTimeoutInput attempts to unpack [input] into the IIBCMsgTimeout type argument
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackTimeoutInput(input []byte) (*MsgTimeout, error) {
	res, err := IBCABI.UnpackInput("Timeout", input)
	if err != nil {
		return nil, err
	}
	unpacked := abi.ConvertType(res[0], new(MsgTimeout)).(*MsgTimeout)
	return unpacked, nil
}

// PackTimeout packs [message] of type IIBCMsgTimeout into the appropriate arguments for Timeout.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackTimeout(message MsgTimeout) ([]byte, error) {
	return IBCABI.Pack("Timeout", message)
}

func timeout(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, TimeoutGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the TimeoutInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackTimeoutInput(input)
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

// UnpackTimeoutOnCloseInput attempts to unpack [input] into the IIBCMsgTimeoutOnClose type argument
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackTimeoutOnCloseInput(input []byte) (*MsgTimeoutOnClose, error) {
	res, err := IBCABI.UnpackInput("TimeoutOnClose", input)
	if err != nil {
		return nil, err
	}
	unpacked := abi.ConvertType(res[0], new(MsgTimeoutOnClose)).(*MsgTimeoutOnClose)
	return unpacked, nil
}

// PackTimeoutOnClose packs [message] of type IIBCMsgTimeoutOnClose into the appropriate arguments for TimeoutOnClose.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackTimeoutOnClose(message MsgTimeoutOnClose) ([]byte, error) {
	return IBCABI.Pack("TimeoutOnClose", message)
}

func timeoutOnClose(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, TimeoutOnCloseGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the TimeoutOnCloseInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackTimeoutOnCloseInput(input)
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
