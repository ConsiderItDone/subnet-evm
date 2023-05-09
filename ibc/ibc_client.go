package ibc

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ava-labs/subnet-evm/rpc"
)

var _ IbcClient = (*ibcClient)(nil)

type IbcClient interface {
	GetProof(ctx context.Context, key ids.ID) ([]byte, error)
}

// ibcClient implementation for interacting with EVM [chain]
type ibcClient struct {
	client *rpc.Client
}

func NewIbcClient(uri, chain string) (IbcClient, error) {
	client, err := rpc.Dial(fmt.Sprintf("%s/ext/bc/%s/rpc", uri, chain))
	if err != nil {
		return nil, fmt.Errorf("failed to dial client. err: %w", err)
	}
	return &ibcClient{
		client: client,
	}, nil
}

func (c *ibcClient) GetProof(ctx context.Context, key ids.ID) ([]byte, error) {
	var res hexutil.Bytes
	err := c.client.CallContext(ctx, &res, "ibc_getProof", key)
	if err != nil {
		return nil, fmt.Errorf("call to ibc_getProof failed. err: %w", err)
	}
	return res, err
}
