package ibc

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// IbcAPI introduces snowman specific functionality to the evm
type IbcAPI struct {
	Backend IbcBackend
}

// GetProof returns the BLS signature by key
func (api *IbcAPI) GetProof(ctx context.Context, key ids.ID) (hexutil.Bytes, error) {
	signature, err := api.Backend.GetProof(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get proof %w", err)
	}
	return signature[:], nil
}
