package ics20

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type (
	FungibleTokenPacketData struct {
		// the token denomination to be transferred
		Denom string `json:"denom"`
		// the token amount to be transferred
		Amount *big.Int `json:"amount"`
		// the sender address
		Sender string `json:"sender"`
		// the recipient address on the destination chain
		Receiver common.Address `json:"receiver"`
		// optional memo
		Memo string `json:"memo,omitempty"`
	}
)

var (
	FungibleTokenAbiStruct   abi.Type
	FungibleTokenAbiArgument abi.Arguments

	ErrDenomNotFound    = errors.New("denom not found")
	ErrSenderNotFound   = errors.New("sender not found")
	ErrReceiverNotFound = errors.New("receiver not found")
)

func init() {
	fungibleTokenAbiStruct, err := abi.NewType("tuple", "struct thing", []abi.ArgumentMarshaling{
		{Name: "denom", Type: "string", InternalType: "string"},
		{Name: "amount", Type: "uint256", InternalType: "uint256"},
		{Name: "sender", Type: "string", InternalType: "string"},
		{Name: "receiver", Type: "address", InternalType: "address"},
		{Name: "memo", Type: "string", InternalType: "string"},
	})
	if err != nil {
		panic(err)
	}

	FungibleTokenAbiStruct = fungibleTokenAbiStruct
	FungibleTokenAbiArgument = abi.Arguments{{
		Name: "rawdata",
		Type: FungibleTokenAbiStruct,
	}}
}

func FungibleTokenPacketDataToABI(rawdata []byte) ([]byte, error) {
	var data FungibleTokenPacketData
	if err := json.Unmarshal(rawdata, &data); err != nil {
		return nil, fmt.Errorf("bad json data: %w", err)
	}

	if len(data.Denom) == 0 {
		return nil, ErrDenomNotFound
	}

	if len(data.Sender) == 0 {
		return nil, ErrSenderNotFound
	}

	abidata, err := FungibleTokenAbiArgument.Pack(&data)
	if err != nil {
		return nil, fmt.Errorf("can't abi encode: %w", err)
	}

	return abidata, nil
}
