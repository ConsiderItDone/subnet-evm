package ibc

import (
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

type ChanOpenInitInput struct {
	PortID  string
	Channel []byte
}

type ChanOpenTryInput struct {
	PortID              string
	Channel             []byte
	CounterpartyVersion string
	ProofInit           []byte
	ProofHeight         []byte
}

type ChannelOpenAckInput struct {
	PortID                string
	ChannelID             string
	CounterpartyChannelID string
	CounterpartyVersion   string
	ProofTry              []byte
	ProofHeight           []byte
}

type ChannelCloseConfirmInput struct {
	PortID      string
	ChannelID   string
	ProofInit   []byte
	ProofHeight []byte
}

type ChannelCloseInitInput struct {
	PortID    string
	ChannelID string
}

type ChannelOpenConfirmInput struct {
	PortID      string
	ChannelID   string
	ProofAck    []byte
	ProofHeight []byte
}

// UnpackChanOpenInitInput attempts to unpack [input] as ChanOpenInitInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackChanOpenInitInput(input []byte) (ChanOpenInitInput, error) {
	inputStruct := ChanOpenInitInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "chanOpenInit", input)

	return inputStruct, err
}

// PackChanOpenInit packs [inputStruct] of type ChanOpenInitInput into the appropriate arguments for chanOpenInit.
func PackChanOpenInit(inputStruct ChanOpenInitInput) ([]byte, error) {
	return IBCABI.Pack("chanOpenInit", inputStruct.PortID, inputStruct.Channel)
}

func chanOpenInit(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ChanOpenInitGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the ChanOpenInitInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackChanOpenInitInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	if err := _chanOpenInit(&callOpts[ChanOpenInitInput]{
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

// UnpackChanOpenTryInput attempts to unpack [input] as ChanOpenTryInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackChanOpenTryInput(input []byte) (ChanOpenTryInput, error) {
	inputStruct := ChanOpenTryInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "chanOpenTry", input)

	return inputStruct, err
}

// PackChanOpenTry packs [inputStruct] of type ChanOpenTryInput into the appropriate arguments for chanOpenTry.
func PackChanOpenTry(inputStruct ChanOpenTryInput) ([]byte, error) {
	return IBCABI.Pack("chanOpenTry", inputStruct.PortID, inputStruct.Channel, inputStruct.CounterpartyVersion, inputStruct.ProofInit, inputStruct.ProofHeight)
}

// PackChanOpenTryOutput attempts to pack given channelID of type string
// to conform the ABI outputs.
func PackChanOpenTryOutput(channelID string) ([]byte, error) {
	return IBCABI.PackOutput("chanOpenTry", channelID)
}

func chanOpenTry(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ChanOpenTryGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the ChanOpenTryInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackChanOpenTryInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	output, err := _chanOpenTry(&callOpts[ChanOpenTryInput]{
		accessibleState: accessibleState,
		caller:          caller,
		addr:            addr,
		suppliedGas:     suppliedGas,
		readOnly:        readOnly,
		args:            inputStruct,
	})
	if err != nil {
		return nil, remainingGas, err
	}

	packedOutput, err := PackChanOpenTryOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackChannelOpenAckInput attempts to unpack [input] as ChannelOpenAckInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackChannelOpenAckInput(input []byte) (ChannelOpenAckInput, error) {
	inputStruct := ChannelOpenAckInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "channelOpenAck", input)

	return inputStruct, err
}

// PackChannelOpenAck packs [inputStruct] of type ChannelOpenAckInput into the appropriate arguments for channelOpenAck.
func PackChannelOpenAck(inputStruct ChannelOpenAckInput) ([]byte, error) {
	return IBCABI.Pack("channelOpenAck", inputStruct.PortID, inputStruct.ChannelID, inputStruct.CounterpartyChannelID, inputStruct.CounterpartyVersion, inputStruct.ProofTry, inputStruct.ProofHeight)
}

func channelOpenAck(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ChannelOpenAckGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the ChannelOpenAckInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackChannelOpenAckInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	if err := _channelOpenAck(&callOpts[ChannelOpenAckInput]{
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

// UnpackChannelOpenConfirmInput attempts to unpack [input] as ChannelOpenConfirmInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackChannelOpenConfirmInput(input []byte) (ChannelOpenConfirmInput, error) {
	inputStruct := ChannelOpenConfirmInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "channelOpenConfirm", input)

	return inputStruct, err
}

// PackChannelOpenConfirm packs [inputStruct] of type ChannelOpenConfirmInput into the appropriate arguments for channelOpenConfirm.
func PackChannelOpenConfirm(inputStruct ChannelOpenConfirmInput) ([]byte, error) {
	return IBCABI.Pack("channelOpenConfirm", inputStruct.PortID, inputStruct.ChannelID, inputStruct.ProofAck, inputStruct.ProofHeight)
}

func channelOpenConfirm(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ChannelOpenConfirmGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the ChannelOpenConfirmInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackChannelOpenConfirmInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	if err := _channelOpenConfirm(&callOpts[ChannelOpenConfirmInput]{
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

// UnpackChannelCloseInitInput attempts to unpack [input] as ChannelCloseInitInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackChannelCloseInitInput(input []byte) (ChannelCloseInitInput, error) {
	inputStruct := ChannelCloseInitInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "channelCloseInit", input)

	return inputStruct, err
}

// PackChannelCloseInit packs [inputStruct] of type ChannelCloseInitInput into the appropriate arguments for channelCloseInit.
func PackChannelCloseInit(inputStruct ChannelCloseInitInput) ([]byte, error) {
	return IBCABI.Pack("channelCloseInit", inputStruct.PortID, inputStruct.ChannelID)
}

func channelCloseInit(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ChannelCloseInitGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the ChannelCloseInitInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackChannelCloseInitInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	if err := _channelCloseInit(&callOpts[ChannelCloseInitInput]{
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

// UnpackChannelCloseConfirmInput attempts to unpack [input] as ChannelCloseConfirmInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackChannelCloseConfirmInput(input []byte) (ChannelCloseConfirmInput, error) {
	inputStruct := ChannelCloseConfirmInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "channelCloseConfirm", input)

	return inputStruct, err
}

// PackChannelCloseConfirm packs [inputStruct] of type ChannelCloseConfirmInput into the appropriate arguments for channelCloseConfirm.
func PackChannelCloseConfirm(inputStruct ChannelCloseConfirmInput) ([]byte, error) {
	return IBCABI.Pack("channelCloseConfirm", inputStruct.PortID, inputStruct.ChannelID, inputStruct.ProofInit, inputStruct.ProofHeight)
}

func channelCloseConfirm(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ChannelCloseConfirmGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the ChannelCloseConfirmInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackChannelCloseConfirmInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	if err := _channelCloseConfirm(&callOpts[ChannelCloseConfirmInput]{
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
