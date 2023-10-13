// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// TODO: replace with gomock
package ibc

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
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

func (n *noopStatefulPrecompileConfig) Timestamp() *big.Int {
	return n.timestamp
}

func (n *noopStatefulPrecompileConfig) IsDisabled() bool {
	return false
}

func (n *noopStatefulPrecompileConfig) Equal(precompileconfig.Config) bool {
	return false
}

func (n *noopStatefulPrecompileConfig) Verify() error {
	return nil
}
