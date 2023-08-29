package ics20

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

type (
	FungibleTokenPacketData struct {
		// the token denomination to be transferred
		Denom string `json:"denom"`
		// the token amount to be transferred
		Amount string `json:"amount"`
		// the sender address
		Sender string `json:"sender"`
		// the recipient address on the destination chain
		Receiver string `json:"receiver"`
		// optional memo
		Memo string `json:"memo,omitempty"`
	}
)

var (
	FungibleTokenAbiArgument abi.Arguments

	ErrDenomNotFound    = errors.New("denom not found")
	ErrAmountNotFound   = errors.New("amount not found")
	ErrAmountCantParse  = errors.New("amount has unknown format")
	ErrSenderNotFound   = errors.New("sender not found")
	ErrReceiverNotFound = errors.New("receiver not found")
	ErrBadReceiverAddr  = errors.New("bad receiver address")
)

func init() {
	var fungibleTokenPacketDataArg abi.Argument
	if err := json.Unmarshal([]byte(`
		{
			"components": [
				{
					"internalType": "string",
					"name": "denom",
					"type": "string"
				},
				{
					"internalType": "uint256",
					"name": "amount",
					"type": "uint256"
				},
				{
					"internalType": "string",
					"name": "sender",
					"type": "string"
				},
				{
					"internalType": "address",
					"name": "receiver",
					"type": "address"
				},
				{
					"internalType": "string",
					"name": "memo",
					"type": "string"
				}
			],
			"internalType": "struct FungibleTokenPacketData",
			"name": "data",
			"type": "tuple"
		}
	`), &fungibleTokenPacketDataArg); err != nil {
		panic(err)
	}
	FungibleTokenAbiArgument = abi.Arguments{fungibleTokenPacketDataArg}
}

func FungibleTokenPacketDataToABI(rawdata []byte) ([]byte, error) {
	var data FungibleTokenPacketData
	if err := json.Unmarshal(rawdata, &data); err != nil {
		return nil, fmt.Errorf("bad json data: %w", err)
	}

	if len(data.Denom) == 0 {
		return nil, ErrDenomNotFound
	}

	if len(data.Amount) == 0 {
		return nil, ErrAmountNotFound
	}

	if len(data.Sender) == 0 {
		return nil, ErrSenderNotFound
	}

	if len(data.Receiver) == 0 {
		return nil, ErrReceiverNotFound
	}

	if !common.IsHexAddress(data.Receiver) {
		return nil, ErrBadReceiverAddr
	}

	amount, ok := math.ParseBig256(data.Amount)
	if !ok {
		return nil, ErrAmountCantParse
	}

	abidata, err := FungibleTokenAbiArgument.Pack(&struct {
		Denom    string
		Amount   *big.Int
		Sender   string
		Receiver common.Address
		Memo     string
	}{
		Denom:    data.Denom,
		Amount:   amount,
		Sender:   data.Sender,
		Receiver: common.HexToAddress(data.Receiver),
		Memo:     data.Memo,
	})
	if err != nil {
		return nil, fmt.Errorf("can't abi encode: %w", err)
	}

	return abidata, nil
}
