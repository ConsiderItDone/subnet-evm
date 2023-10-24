package ibc

import (
	"fmt"
	"math/big"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/ics20"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/log"
)

// Height is an auto generated low-level Go binding around an user-defined struct.
type Height struct {
	RevisionNumber *big.Int
	RevisionHeight *big.Int
}

// Packet is an auto generated low-level Go binding around an user-defined struct.
type Packet struct {
	Sequence           *big.Int
	SourcePort         string
	SourceChannel      string
	DestinationPort    string
	DestinationChannel string
	Data               []byte
	TimeoutHeight      Height
	TimeoutTimestamp   *big.Int
}

// IIBCMsgAcknowledgement is an auto generated low-level Go binding around an user-defined struct.
type IIBCMsgAcknowledgement struct {
	Packet          Packet
	Acknowledgement []byte
	ProofAcked      []byte
	ProofHeight     Height
	Signer          string
}

type OnAcknowledgementPacketInput struct {
	Packet          Packet
	Acknowledgement []byte
	Relayer         []byte
}

// IIBCMsgTimeout is an auto generated low-level Go binding around an user-defined struct.
type IIBCMsgTimeout struct {
	Packet           Packet
	ProofUnreceived  []byte
	ProofHeight      Height
	NextSequenceRecv *big.Int
	Signer           string
}

type OnTimeoutInput struct {
	Packet  Packet
	Relayer []byte
}

// IIBCMsgTimeoutOnClose is an auto generated low-level Go binding around an user-defined struct.
type IIBCMsgTimeoutOnClose struct {
	Packet           Packet
	ProofUnreceived  []byte
	ProofClose       []byte
	ProofHeight      Height
	NextSequenceRecv *big.Int
	Signer           string
}

type OnTimeoutOnCloseInput struct {
	Packet  Packet
	Relayer []byte
}

type MsgSendPacket struct {
	ChannelCapability *big.Int
	SourcePort        string
	SourceChannel     string
	TimeoutHeight     Height
	TimeoutTimestamp  *big.Int
	Data              []byte
}

// IIBCMsgRecvPacket is an auto generated low-level Go binding around an user-defined struct.
type IIBCMsgRecvPacket struct {
	Packet          Packet
	ProofCommitment []byte
	ProofHeight     Height
	Signer          string
}

type OnRecvPacketInput struct {
	Packet  Packet
	Relayer []byte
}

// PackOnRecvPacket packs [inputStruct] of type OnRecvPacketInput into the appropriate arguments for OnRecvPacket.
func PackOnRecvPacket(inputStruct OnRecvPacketInput) ([]byte, error) {
	return IBCABI.Pack("OnRecvPacket", inputStruct.Packet, inputStruct.Relayer)
}

// PackOnTimeoutOnCloseInput packs [inputStruct] of type OnTimeoutOnCloseInput into the appropriate arguments for OnTimeoutOnClose.
func PackOnTimeoutOnCloseInput(inputStruct OnTimeoutOnCloseInput) ([]byte, error) {
	return IBCABI.Pack("OnTimeoutOnClose", inputStruct.Packet, inputStruct.Relayer)
}

// PackOnTimeout packs [inputStruct] of type OnTimeoutInput into the appropriate arguments for OnTimeout.
func PackOnTimeout(inputStruct OnTimeoutInput) ([]byte, error) {
	return IBCABI.Pack("OnTimeout", inputStruct.Packet, inputStruct.Relayer)
}

// PackOnAcknowledgementPacket packs [inputStruct] of type OnAcknowledgementPacketInput into the appropriate arguments for OnAcknowledgementPacket.
func PackOnAcknowledgementPacket(inputStruct OnAcknowledgementPacketInput) ([]byte, error) {
	return IBCABI.Pack("OnAcknowledgementPacket", inputStruct.Packet, inputStruct.Acknowledgement, inputStruct.Relayer)
}

// UnpackSendPacketInput attempts to unpack [input] as SendPacketInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackSendPacketInput(input []byte) (MsgSendPacket, error) {
	inputStruct := MsgSendPacket{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "sendPacket", input)

	return inputStruct, err
}

// PackSendPacket packs [inputStruct] of type SendPacketInput into the appropriate arguments for SendPacket.
func PackSendPacket(inputStruct MsgSendPacket) ([]byte, error) {
	return IBCABI.Pack("sendPacket", inputStruct.ChannelCapability, inputStruct.SourcePort, inputStruct.SourceChannel, inputStruct.TimeoutHeight, inputStruct.TimeoutTimestamp, inputStruct.Data)
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

	if err := _sendPacket(&callOpts[MsgSendPacket]{
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
func UnpackRecvPacketInput(input []byte) (*IIBCMsgRecvPacket, error) {
	res, err := IBCABI.UnpackInput("recvPacket", input)
	if err != nil {
		return nil, err
	}
	unpacked := *abi.ConvertType(res[0], new(IIBCMsgRecvPacket)).(*IIBCMsgRecvPacket)
	return &unpacked, nil
}

// PackRecvPacket packs [message] of type IIBCMsgRecvPacket into the appropriate arguments for RecvPacket.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackRecvPacket(message IIBCMsgRecvPacket) ([]byte, error) {
	return IBCABI.Pack("recvPacket", message)
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

	if err := _recvPacket(&callOpts[IIBCMsgRecvPacket]{
		accessibleState: accessibleState,
		caller:          caller,
		addr:            addr,
		suppliedGas:     suppliedGas,
		readOnly:        readOnly,
		args:            *inputStruct,
	}); err != nil {
		return nil, remainingGas, err
	}

	recvAddr, err := GetPort(accessibleState.GetStateDB(), inputStruct.Packet.DestinationPort)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("%w, port with portID: %s already bound", err, inputStruct.Packet.DestinationPort)
	}

	if inputStruct.Packet.DestinationPort == "transfer" {
		packetData, err := ics20.FungibleTokenPacketDataToABI(inputStruct.Packet.Data)
		if err != nil {
			return nil, remainingGas, err
		}
		inputStruct.Packet.Data = packetData
	}

	data, err := PackOnRecvPacket(OnRecvPacketInput{Packet: inputStruct.Packet, Relayer: []byte(inputStruct.Signer)})
	if err != nil {
		return nil, remainingGas, err
	}

	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	// ToDo: who mush gas do we need? now it is hardcoded as 100k
	_, remainingGas, err = accessibleState.CallFromPrecompile(ContractAddress, recvAddr, data, 100000, big.NewInt(0))
	if err != nil {
		topics, data, err := IBCABI.PackEvent(GeneratedAcknowledgementErrorIdentifier.RawName,
			inputStruct.Packet.Data,
			inputStruct.Packet.TimeoutHeight.String(),
			big.NewInt(inputStruct.Packet.TimeoutTimestamp.Int64()),
			big.NewInt(inputStruct.Packet.Sequence.Int64()),
			inputStruct.Packet.SourcePort,
			inputStruct.Packet.SourceChannel,
			inputStruct.Packet.DestinationPort,
			inputStruct.Packet.DestinationChannel,
			fmt.Sprintf("%s", err),
		)
		if err != nil {
			return nil, remainingGas, fmt.Errorf("error packing event: %w", err)
		}
		blockNumber := accessibleState.GetBlockContext().Number().Uint64()
		accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

		log.Warn("EXEPTION: recvPacket has bad request from CallFromPrecompile")

		ack = channeltypes.NewErrorAcknowledgement(err)
	}

	err = writeAcknowledgement(inputStruct.Packet, accessibleState, ack)
	if err != nil {
		return nil, remainingGas, err
	}
	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackTimeoutInput attempts to unpack [input] into the IIBCMsgTimeout type argument
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackTimeoutInput(input []byte) (IIBCMsgTimeout, error) {
	inputStruct := IIBCMsgTimeout{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "timeout", input)
	if err != nil {
		fmt.Println("UnpackInputIntoInterface")
	}
	return inputStruct, err
}

// PackTimeout packs [message] of type IIBCMsgTimeout into the appropriate arguments for Timeout.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackTimeout(message IIBCMsgTimeout) ([]byte, error) {
	return IBCABI.Pack("timeout", message.Packet, message.ProofUnreceived, message.ProofHeight, message.NextSequenceRecv, message.Signer)
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
	packedOutput := []byte{}

	err = _timeout(&callOpts[IIBCMsgTimeout]{
		accessibleState: accessibleState,
		caller:          caller,
		addr:            addr,
		suppliedGas:     suppliedGas,
		readOnly:        readOnly,
		args:            inputStruct,
	})
	switch err {
	case nil:
	case channeltypes.ErrNoOpMsg:
		return packedOutput, remainingGas, nil
	default:
		return nil, remainingGas, err
	}

	recvAddr, err := GetPort(accessibleState.GetStateDB(), inputStruct.Packet.DestinationPort)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("%w, port with portID: %s already bound", err, inputStruct.Packet.DestinationPort)
	}
	data, err := PackOnTimeout(OnTimeoutInput{Packet: inputStruct.Packet, Relayer: []byte(inputStruct.Signer)})
	if err != nil {
		return nil, remainingGas, err
	}
	_, remainingGas, err = accessibleState.CallFromPrecompile(ContractAddress, recvAddr, data, remainingGas, big.NewInt(0))
	if err != nil {
		return nil, remainingGas, err
	}
	// Delete packet commitment
	if err = TimeoutExecuted(accessibleState, inputStruct.Packet); err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackTimeoutOnCloseInput attempts to unpack [input] into the IIBCMsgTimeoutOnClose type argument
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackTimeoutOnCloseInput(input []byte) (IIBCMsgTimeoutOnClose, error) {
	inputStruct := IIBCMsgTimeoutOnClose{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "timeoutOnClose", input)
	if err != nil {
		fmt.Println("UnpackInputIntoInterface")
	}
	return inputStruct, err
}

// PackTimeoutOnClose packs [message] of type IIBCMsgTimeoutOnClose into the appropriate arguments for TimeoutOnClose.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackTimeoutOnClose(message IIBCMsgTimeoutOnClose) ([]byte, error) {
	return IBCABI.Pack("timeoutOnClose", message.Packet, message.ProofUnreceived, message.ProofClose, message.ProofHeight, message.NextSequenceRecv, message.Signer)
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
	packedOutput := []byte{}

	err = _timeoutOnClose(&callOpts[IIBCMsgTimeoutOnClose]{
		accessibleState: accessibleState,
		caller:          caller,
		addr:            addr,
		suppliedGas:     suppliedGas,
		readOnly:        readOnly,
		args:            inputStruct,
	})
	switch err {
	case nil:
	case channeltypes.ErrNoOpMsg:
		return packedOutput, remainingGas, nil
	default:
		return nil, remainingGas, err
	}

	recvAddr, err := GetPort(accessibleState.GetStateDB(), inputStruct.Packet.DestinationPort)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("%w, port with portID: %s already bound", err, inputStruct.Packet.DestinationPort)
	}
	data, err := PackOnTimeoutOnCloseInput(OnTimeoutOnCloseInput{Packet: inputStruct.Packet, Relayer: []byte(inputStruct.Signer)})
	if err != nil {
		return nil, remainingGas, err
	}
	_, remainingGas, err = accessibleState.CallFromPrecompile(ContractAddress, recvAddr, data, remainingGas, big.NewInt(0))
	if err != nil {
		return nil, remainingGas, err
	}
	// Delete packet commitment
	if err = TimeoutExecuted(accessibleState, inputStruct.Packet); err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

func (h Height) String() string {
	return fmt.Sprintf("%d-%d", h.RevisionNumber.Uint64(), h.RevisionHeight.Uint64())
}

// UnpackAcknowledgementInput attempts to unpack [input] into the IIBCMsgAcknowledgement type argument
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackAcknowledgementInput(input []byte) (IIBCMsgAcknowledgement, error) {
	inputStruct := IIBCMsgAcknowledgement{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "acknowledgement", input)

	return inputStruct, err
}

// PackAcknowledgement packs [message] of type IIBCMsgAcknowledgement into the appropriate arguments for Acknowledgement.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackAcknowledgement(inputStruct IIBCMsgAcknowledgement) ([]byte, error) {
	return IBCABI.Pack("acknowledgement", inputStruct.Packet, inputStruct.Acknowledgement, inputStruct.ProofAcked, inputStruct.ProofHeight, inputStruct.Signer)
}

func acknowledgement(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, AcknowledgementGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the TimeoutOnCloseInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackAcknowledgementInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	err = _acknowledgement(&callOpts[IIBCMsgAcknowledgement]{
		accessibleState: accessibleState,
		caller:          caller,
		addr:            addr,
		suppliedGas:     suppliedGas,
		readOnly:        readOnly,
		args:            inputStruct,
	})
	switch err {
	case nil:
	case channeltypes.ErrNoOpMsg:
		return []byte{}, remainingGas, nil
	default:
		return nil, remainingGas, err
	}

	recvAddr, err := GetPort(accessibleState.GetStateDB(), inputStruct.Packet.DestinationPort)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("%w, port with portID: %s already bound", err, inputStruct.Packet.DestinationPort)
	}
	data, err := PackOnAcknowledgementPacket(OnAcknowledgementPacketInput{Packet: inputStruct.Packet, Acknowledgement: inputStruct.Acknowledgement, Relayer: []byte(inputStruct.Signer)})
	if err != nil {
		return nil, remainingGas, err
	}
	_, remainingGas, err = accessibleState.CallFromPrecompile(ContractAddress, recvAddr, data, remainingGas, big.NewInt(0))
	if err != nil {
		return nil, remainingGas, fmt.Errorf("can't call fuction via CallFromPrecompile: %w", err)
	}

	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}
