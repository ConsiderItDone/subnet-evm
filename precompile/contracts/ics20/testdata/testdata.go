package testdata

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind/backends"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/precompile/contracts/ics20/testdata/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var PrivatKey *ecdsa.PrivateKey

func init() {
	pkey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	PrivatKey = pkey
}

func NewCodecEnv() (*backends.SimulatedBackend, *bind.TransactOpts, *codec.Codec, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(PrivatKey, big.NewInt(1337))
	if err != nil {
		return nil, nil, nil, err
	}

	balance := new(big.Int)
	balance.SetString("10000000000000000000000", 10)
	allocations := map[common.Address]core.GenesisAccount{
		auth.From: {
			Balance: balance,
		},
	}

	backend := backends.NewSimulatedBackend(allocations, 4712388)
	_, _, contract, err := codec.DeployCodec(auth, backend)
	if err != nil {
		return nil, nil, nil, err
	}
	backend.Commit(true)

	return backend, auth, contract, nil
}
