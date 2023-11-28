// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// TODO: replace with gomock
package ibc

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
)

var _ precompileconfig.Config = &noopStatefulPrecompileConfig{}

type noopStatefulPrecompileConfig struct {
	timestamp *big.Int
}

func NewNoopStatefulPrecompileConfig() *noopStatefulPrecompileConfig {
	return &noopStatefulPrecompileConfig{}
}

func (n *noopStatefulPrecompileConfig) Key() string {
	return ""
}

func (n *noopStatefulPrecompileConfig) Address() common.Address {
	return common.Address{}
}

func (n *noopStatefulPrecompileConfig) Timestamp() *uint64 {
	t := n.timestamp.Uint64()

	return &t
}

func (n *noopStatefulPrecompileConfig) IsDisabled() bool {
	return false
}

func (n *noopStatefulPrecompileConfig) Equal(precompileconfig.Config) bool {
	return false
}

func (n *noopStatefulPrecompileConfig) Verify(chainConfig precompileconfig.ChainConfig) error {
	return nil
}
