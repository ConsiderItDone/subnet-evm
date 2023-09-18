package ibc

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

type QueryClientStateInput struct {
	ClientID string
}

type QueryConsensusStateInput struct {
	ClientID string
}

type QueryConnectionInput struct {
	ConnectionID string
}

type QueryChannelInput struct {
	PortID    string
	ChannelID string
}

// PackQueryClientStateInput packs [inputStruct] of type QueryClientStateInput into the appropriate arguments for queryClientState.
func PackQueryClientStateInput(inputStruct QueryClientStateInput) ([]byte, error) {
	return IBCABI.Pack("queryClientState", inputStruct.ClientID)
}

// UnpackQueryClientStateInput attempts to unpack [input] as QueryClientStateInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackQueryClientStateInput(input []byte) (QueryClientStateInput, error) {
	inputStruct := QueryClientStateInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "queryClientState", input)

	return inputStruct, err
}

// PackQueryClientStateOutput attempts to pack given clientState of type []byte
// to conform the ABI outputs.
func PackQueryClientStateOutput(clientState []byte) ([]byte, error) {
	return IBCABI.PackOutput("queryClientState", clientState)
}

func queryClientState(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, QueryClientStateGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the QueryClientState.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackQueryClientStateInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	clientState, err := GetClientState(accessibleState.GetStateDB(), inputStruct.ClientID)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("error loading client state, err: %w", err)
	}

	out, err := clientState.Marshal()
	if err != nil {
		return nil, remainingGas, err
	}

	packedOutput, err := PackQueryClientStateOutput(out)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// PackQueryConsensusStateInput packs [inputStruct] of type QueryConsensusStateInput into the appropriate arguments for queryConsensusState.
func PackQueryConsensusStateInput(inputStruct QueryConsensusStateInput) ([]byte, error) {
	return IBCABI.Pack("queryConsensusState", inputStruct.ClientID)
}

// UnpackQueryConsensusStateInput attempts to unpack [input] as QueryConsensusStateInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackQueryConsensusStateInput(input []byte) (QueryConsensusStateInput, error) {
	inputStruct := QueryConsensusStateInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "queryConsensusState", input)

	return inputStruct, err
}

// PackQueryConsensusState attempts to pack given consensusState of type []byte
// to conform the ABI outputs.
func PackQueryConsensusStateOutput(consensusState []byte) ([]byte, error) {
	return IBCABI.PackOutput("queryConsensusState", consensusState)
}

func queryConsensusState(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, QueryConsensusStateGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the QueryConsensusState.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackQueryConsensusStateInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	clientState, err := GetClientState(accessibleState.GetStateDB(), inputStruct.ClientID)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("error loading client state, err: %w", err)
	}

	consensusState, err := GetConsensusState(accessibleState.GetStateDB(), inputStruct.ClientID, clientState.LatestHeight)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("error loading consensus state, err: %w", err)
	}

	out, err := consensusState.Marshal()
	if err != nil {
		return nil, remainingGas, err
	}

	packedOutput, err := PackQueryConsensusStateOutput(out)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// PackQueryConnectionInput packs [inputStruct] of type QueryConnectionInput into the appropriate arguments for queryConnection.
func PackQueryConnectionInput(inputStruct QueryConnectionInput) ([]byte, error) {
	return IBCABI.Pack("queryConnection", inputStruct.ConnectionID)
}

// UnpackQueryConnectionInput attempts to unpack [input] as QueryConnectionInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackQueryConnectionInput(input []byte) (QueryConnectionInput, error) {
	inputStruct := QueryConnectionInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "queryConnection", input)

	return inputStruct, err
}

// PackQueryConnection attempts to pack given connection of type []byte
// to conform the ABI outputs.
func PackQueryConnectionOutput(connection []byte) ([]byte, error) {
	return IBCABI.PackOutput("queryConnection", connection)
}

func queryConnection(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, QueryConnectionGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the QueryConnection.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackQueryConnectionInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	connection, err := GetConnection(accessibleState.GetStateDB(), inputStruct.ConnectionID)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("error loading connection, err: %w", err)
	}

	out, err := connection.Marshal()
	if err != nil {
		return nil, remainingGas, err
	}

	packedOutput, err := PackQueryConnectionOutput(out)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// PackQueryChannelInput packs [inputStruct] of type QueryChannelInput into the appropriate arguments for queryChannel.
func PackQueryChannelInput(inputStruct QueryChannelInput) ([]byte, error) {
	return IBCABI.Pack("queryChannel", inputStruct.PortID, inputStruct.ChannelID)
}

// UnpackQueryChannelInput attempts to unpack [input] as QueryChannelInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackQueryChannelInput(input []byte) (QueryChannelInput, error) {
	inputStruct := QueryChannelInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "queryChannel", input)

	return inputStruct, err
}

// PackQueryChannel attempts to pack given channel of type []byte
// to conform the ABI outputs.
func PackQueryChannelOutput(channel []byte) ([]byte, error) {
	return IBCABI.PackOutput("queryChannel", channel)
}

func queryChannel(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, QueryChannelGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the QueryChannel.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackQueryChannelInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	channel, err := GetChannel(accessibleState.GetStateDB(), inputStruct.PortID, inputStruct.ChannelID)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("error loading channel, err: %w", err)
	}

	out, err := channel.Marshal()
	if err != nil {
		return nil, remainingGas, err
	}

	packedOutput, err := PackQueryChannelOutput(out)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}
