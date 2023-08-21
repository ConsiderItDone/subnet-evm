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
	// Gas costs for each function. These are set to 1 by default.
	// You should set a gas cost for each function in your contract.
	// Generally, you should not set gas costs very low as this may cause your network to be vulnerable to DoS attacks.
	// There are some predefined gas costs in contract/utils.go that you can use.
	ChanOpenInitGasCost        uint64 = 1 /* SET A GAS COST HERE */
	ChanOpenTryGasCost         uint64 = 1 /* SET A GAS COST HERE */
	ChannelCloseConfirmGasCost uint64 = 1 /* SET A GAS COST HERE */
	ChannelCloseInitGasCost    uint64 = 1 /* SET A GAS COST HERE */
	ChannelOpenAckGasCost      uint64 = 1 /* SET A GAS COST HERE */
	ChannelOpenConfirmGasCost  uint64 = 1 /* SET A GAS COST HERE */
	ConnOpenAckGasCost         uint64 = 1 /* SET A GAS COST HERE */
	ConnOpenConfirmGasCost     uint64 = 1 /* SET A GAS COST HERE */
	ConnOpenInitGasCost        uint64 = 1 /* SET A GAS COST HERE */
	ConnOpenTryGasCost         uint64 = 1 /* SET A GAS COST HERE */
	CreateClientGasCost        uint64 = 1 /* SET A GAS COST HERE */
	UpdateClientGasCost        uint64 = 1 /* SET A GAS COST HERE */
	UpgradeClientGasCost       uint64 = 1 /* SET A GAS COST HERE */
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

	IBCABI                        = contract.ParseABI(IBCRawABI)
	IBCPrecompile                 = createIBCPrecompile()
	GeneratedClientIdentifier     = IBCABI.Events["ClientCreated"]
	GeneratedConnectionIdentifier = IBCABI.Events["ConnectionCreated"]

	NextClientSeqStorageKey = common.Hash{'n', 'c', 's', 'e', 'q', 's', 'k'}

	ErrWrongClientType = errors.New("wrong client type. Only Tendermint supported")
)

// createIBCPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.
func createIBCPrecompile() contract.StatefulPrecompiledContract {
	var functions []*contract.StatefulPrecompileFunction

	abiFunctionMap := map[string]contract.RunStatefulPrecompileFunc{
		"chanOpenInit":        chanOpenInit,
		"chanOpenTry":         chanOpenTry,
		"channelCloseConfirm": channelCloseConfirm,
		"channelCloseInit":    channelCloseInit,
		"channelOpenAck":      channelOpenAck,
		"channelOpenConfirm":  channelOpenConfirm,
		"connOpenAck":         connOpenAck,
		"connOpenConfirm":     connOpenConfirm,
		"connOpenInit":        connOpenInit,
		"connOpenTry":         connOpenTry,
		"createClient":        createClient,
		"updateClient":        updateClient,
		"upgradeClient":       upgradeClient,
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
