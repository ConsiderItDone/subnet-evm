package ibc

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	hosttypes "github.com/cosmos/ibc-go/v7/modules/core/24-host"
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

	_, ok, _ := getPort(accessibleState.GetStateDB(), portID)
	if !ok {
		return nil, remainingGas, fmt.Errorf("port with portID: %s already bound", portID)
	}
	err = storePortID(accessibleState.GetStateDB(), portID, caller)
	if err != nil {
		return nil, remainingGas, err
	}
	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

func storePortID(db contract.StateDB, portID string, caller common.Address) error {
	if err := hosttypes.PortIdentifierValidator(portID); err != nil {
		panic(err.Error())
	}

	bz := caller[:]

	key := calculateKey([]byte(hosttypes.PortPath(portID)))
	db.SetPrecompileState(common.BytesToAddress([]byte(key)), bz)
	return nil
}

func getPort(db contract.StateDB, portID string) (common.Address, bool, error) {
	key := calculateKey([]byte(hosttypes.PortPath(portID)))
	bz := db.GetPrecompileState(common.BytesToAddress([]byte(key)))

	if len(bz) == 0 {
		return common.Address{}, false, fmt.Errorf("Bind port with this portID: %s not exist", portID)
	}

	if len(bz) != common.AddressLength {
		return common.Address{}, false, fmt.Errorf("Lenght of data by this portID: %s, not equal AddressLength", portID)
	}

	return common.BytesToAddress([]byte(key)), true, nil
}

func makeCapability(db contract.StateDB, portID, channelID string) error {
	name := hosttypes.ChannelCapabilityPath(portID, channelID)

	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("capability name cannot be empty")
	}

	if ok, _ := getCapability(db, portID, channelID); ok {
		return fmt.Errorf("Capability with portID: %s, name: %s already exist", portID, name)
	}

	bz := []byte{1}

	key := calculateKey([]byte(name))
	db.SetPrecompileState(common.BytesToAddress([]byte(key)), bz)
	return nil
}

func getCapability(db contract.StateDB, portID, channelID string) (bool, error) {
	name := hosttypes.ChannelCapabilityPath(portID, channelID)
	key := calculateKey([]byte(name))
	bz := db.GetPrecompileState(common.BytesToAddress([]byte(key)))

	if len(bz) == 0 {
		return false, fmt.Errorf("Capability with this name: %s not exist", name)
	}

	// TODO make deleteCapability func to operate all condition of variable
	if !reflect.DeepEqual(bz, []byte{1}) {
		return false, fmt.Errorf("Capability with this name: %s has bad data", name)
	}

	return true, nil
}
