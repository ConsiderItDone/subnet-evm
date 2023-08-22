package utils

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/ethclient"
)

type EVMStorageReader struct {
	ethClient ethclient.Client
}

func NewEVMStorageReader(client ethclient.Client) *EVMStorageReader {
	return &EVMStorageReader{
		ethClient: client,
	}
}

func (r *EVMStorageReader) GetState(addr common.Address, hash common.Hash) common.Hash {
	data, err := ethClient.StorageAt(context.Background(), addr, hash, nil)
	if err != nil {
		return common.Hash{}
	}

	return common.BytesToHash(data)
}
