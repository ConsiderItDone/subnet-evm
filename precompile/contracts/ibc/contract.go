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
	CreateClientGasCost uint64 = 1
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

	IBCABI                    = contract.ParseABI(IBCRawABI)
	IBCPrecompile             = createIBCPrecompile()
	GeneratedClientIdentifier = IBCABI.Events["ClientCreated"]

	nextClientSeqStorageKey = common.Hash{'n', 'c', 's', 'e', 'q', 's', 'k'}
	clientStateStorageKey   = common.Hash{'c', 's', 't', 's', 'k'}

	ErrWrongClientType = errors.New("wrong client type. Only Tendermint supported")
)

// createIBCPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.
func createIBCPrecompile() contract.StatefulPrecompiledContract {
	var functions []*contract.StatefulPrecompileFunction

	abiFunctionMap := map[string]contract.RunStatefulPrecompileFunc{
		"createClient": createClient,
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
