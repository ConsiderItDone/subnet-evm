package ics20

import (
	"encoding/json"
	"errors"
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

	ErrBadDemonValue    = errors.New("denom has bad json data")
	ErrBadJsonValue     = errors.New("isc20: has bad json data")
	ErrBadAmountValue   = errors.New("amount has bad value")
	ErrBadSenderValue   = errors.New("sender has bad value")
	ErrBadReceiverValue = errors.New("receiver has bad value")
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
					"internalType": "bytes",
					"name": "sender",
					"type": "bytes"
				},
				{
					"internalType": "bytes",
					"name": "receiver",
					"type": "bytes"
				},
				{
					"internalType": "bytes",
					"name": "memo",
					"type": "bytes"
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
		return nil, ErrBadJsonValue
	}

	if len(data.Denom) == 0 {
		return nil, ErrBadDemonValue
	}

	if len(data.Amount) == 0 {
		return nil, ErrBadAmountValue
	}

	amount, ok := math.ParseBig256(data.Amount)
	if !ok {
		return nil, ErrBadAmountValue
	}

	sender, err := common.ParseHexOrString(data.Sender)
	if err != nil {
		return nil, ErrBadSenderValue
	}

	receiver, err := common.ParseHexOrString(data.Receiver)
	if err != nil {
		return nil, ErrBadReceiverValue
	}

	return FungibleTokenAbiArgument.Pack(&struct {
		Denom    string
		Amount   *big.Int
		Sender   []byte
		Receiver []byte
		Memo     []byte
	}{
		Denom:    data.Denom,
		Amount:   amount,
		Sender:   sender,
		Receiver: receiver,
		Memo:     []byte(data.Memo),
	})
}
