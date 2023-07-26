package ibc

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
)

type ConnOpenAckInput struct {
	ConnectionID             string
	ClientState              []byte
	Version                  []byte
	CounterpartyConnectionID []byte
	ProofTry                 []byte
	ProofClient              []byte
	ProofConsensus           []byte
	ProofHeight              []byte
	ConsensusHeight          []byte
}

type ConnOpenConfirmInput struct {
	ConnectionID string
	ProofAck     []byte
	ProofHeight  []byte
}

type ConnOpenInitInput struct {
	ClientID     string
	Counterparty []byte
	Version      []byte
	DelayPeriod  uint32
}

type ConnOpenTryInput struct {
	Counterparty         []byte
	DelayPeriod          uint32
	ClientID             string
	ClientState          []byte
	CounterpartyVersions []byte
	ProofInit            []byte
	ProofClient          []byte
	ProofConsensus       []byte
	ProofHeight          []byte
	ConsensusHeight      []byte
}

// UnpackConnOpenAckInput attempts to unpack [input] as ConnOpenAckInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackConnOpenAckInput(input []byte) (ConnOpenAckInput, error) {
	inputStruct := ConnOpenAckInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "connOpenAck", input)

	return inputStruct, err
}

// PackConnOpenAck packs [inputStruct] of type ConnOpenAckInput into the appropriate arguments for connOpenAck.
func PackConnOpenAck(inputStruct ConnOpenAckInput) ([]byte, error) {
	return IBCABI.Pack("connOpenAck", inputStruct.ConnectionID, inputStruct.ClientState, inputStruct.Version, inputStruct.CounterpartyConnectionID, inputStruct.ProofTry, inputStruct.ProofClient, inputStruct.ProofConsensus, inputStruct.ProofHeight, inputStruct.ConsensusHeight)
}

func connOpenAck(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ConnOpenAckGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the ConnOpenAckInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackConnOpenAckInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	if err := _connOpenAck(&callOpts[ConnOpenAckInput]{
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

// UnpackConnOpenConfirmInput attempts to unpack [input] as ConnOpenConfirmInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackConnOpenConfirmInput(input []byte) (ConnOpenConfirmInput, error) {
	inputStruct := ConnOpenConfirmInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "connOpenConfirm", input)

	return inputStruct, err
}

// PackConnOpenConfirm packs [inputStruct] of type ConnOpenConfirmInput into the appropriate arguments for connOpenConfirm.
func PackConnOpenConfirm(inputStruct ConnOpenConfirmInput) ([]byte, error) {
	return IBCABI.Pack("connOpenConfirm", inputStruct.ConnectionID, inputStruct.ProofAck, inputStruct.ProofHeight)
}

func connOpenConfirm(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ConnOpenConfirmGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the ConnOpenConfirmInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackConnOpenConfirmInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	if err := _connOpenConfirm(&callOpts[ConnOpenConfirmInput]{
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

// UnpackConnOpenInitInput attempts to unpack [input] as ConnOpenInitInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackConnOpenInitInput(input []byte) (ConnOpenInitInput, error) {
	inputStruct := ConnOpenInitInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "connOpenInit", input)

	return inputStruct, err
}

// PackConnOpenInit packs [inputStruct] of type ConnOpenInitInput into the appropriate arguments for connOpenInit.
func PackConnOpenInit(inputStruct ConnOpenInitInput) ([]byte, error) {
	return IBCABI.Pack("connOpenInit", inputStruct.ClientID, inputStruct.Counterparty, inputStruct.Version, inputStruct.DelayPeriod)
}

// PackConnOpenInitOutput attempts to pack given connectionID of type string
// to conform the ABI outputs.
func PackConnOpenInitOutput(connectionID string) ([]byte, error) {
	return IBCABI.PackOutput("connOpenInit", connectionID)
}

func connOpenInit(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ConnOpenInitGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the ConnOpenInitInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackConnOpenInitInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	connectionID, err := _connOpenInit(&callOpts[ConnOpenInitInput]{
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

	packedOutput, err := PackConnOpenInitOutput(connectionID)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackConnOpenTryInput attempts to unpack [input] as ConnOpenTryInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackConnOpenTryInput(input []byte) (ConnOpenTryInput, error) {
	inputStruct := ConnOpenTryInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "connOpenTry", input)

	return inputStruct, err
}

// PackConnOpenTry packs [inputStruct] of type ConnOpenTryInput into the appropriate arguments for connOpenTry.
func PackConnOpenTry(inputStruct ConnOpenTryInput) ([]byte, error) {
	return IBCABI.Pack("connOpenTry", inputStruct.Counterparty, inputStruct.DelayPeriod, inputStruct.ClientID, inputStruct.ClientState, inputStruct.CounterpartyVersions, inputStruct.ProofInit, inputStruct.ProofClient, inputStruct.ProofConsensus, inputStruct.ProofHeight, inputStruct.ConsensusHeight)
}

// PackConnOpenTryOutput attempts to pack given connectionID of type string
// to conform the ABI outputs.
func PackConnOpenTryOutput(connectionID string) ([]byte, error) {
	return IBCABI.PackOutput("connOpenTry", connectionID)
}

func connOpenTry(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ConnOpenTryGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the ConnOpenTryInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackConnOpenTryInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	connectionID, err := _connOpenTry(&callOpts[ConnOpenTryInput]{
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

	packedOutput, err := PackConnOpenTryOutput(connectionID)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}
