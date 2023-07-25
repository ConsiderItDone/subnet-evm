package ibc

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/ethereum/go-ethereum/common"
)

var (
	keyClientSeq = common.BytesToHash([]byte("client-seq"))

	ErrWrongClientType = errors.New("wrong client type")
)

type callOpts[T any] struct {
	accessibleState contract.AccessibleState
	caller          common.Address
	addr            common.Address
	suppliedGas     uint64
	readOnly        bool
	args            T
}

func _createClient(opts *callOpts[CreateClientInput]) (string, error) {
	if opts.args.ClientType != exported.Tendermint {
		return "", ErrWrongClientType
	}

	db := opts.accessibleState.GetStateDB()
	clientSeq := db.GetState(ContractAddress, keyClientSeq)
	newClientSeq := common.BigToHash(
		new(big.Int).Add(
			clientSeq.Big(),
			big.NewInt(1),
		),
	)
	db.SetState(ContractAddress, keyClientSeq, newClientSeq)

	return fmt.Sprintf("%s-%d", opts.args.ClientType, clientSeq.Big().Int64()), nil
}

func _updateClient(opts *callOpts[UpdateClientInput]) error {
	panic("not implemented")
}

func _upgradeClient(opts *callOpts[UpgradeClientInput]) error {
	panic("not implemented")
}

func _connOpenInit(opts *callOpts[ConnOpenInitInput]) (string, error) {
	panic("not implemented")
}

func _connOpenTry(opts *callOpts[ConnOpenTryInput]) (string, error) {
	panic("not implemented")
}

func _connOpenAck(opts *callOpts[ConnOpenAckInput]) error {
	panic("not implemented")
}

func _connOpenConfirm(opts *callOpts[ConnOpenConfirmInput]) error {
	panic("not implemented")
}
