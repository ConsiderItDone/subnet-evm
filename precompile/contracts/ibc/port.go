package ibc

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

// UnpackBindPortInput attempts to unpack [input] into the string type argument
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackBindPortInput(input []byte) (string, error) {
	res, err := IBCABI.UnpackInput("bindPort", input)
	if err != nil {
		return "", err
	}
	unpacked := *abi.ConvertType(res[0], new(string)).(*string)
	return unpacked, nil
}

// PackBindPort packs [portID] of type string into the appropriate arguments for bindPort.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackBindPort(portID string) ([]byte, error) {
	return IBCABI.Pack("bindPort", portID)
}

func bindPort(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, BindPortGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// TODO CHECK caller PERMISSION

	// attempts to unpack [input] into the arguments to the BindPortInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackBindPortInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	portID := inputStruct

	if _, err := GetPort(accessibleState.GetStateDB(), portID); err == nil {
		return nil, remainingGas, fmt.Errorf("port with portID: %s already bound", portID)
	}

	if err = SetPort(accessibleState.GetStateDB(), portID, caller); err != nil {
		return nil, remainingGas, err
	}

	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}
