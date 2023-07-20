// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ibc_old

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
)

func TestEqual(t *testing.T) {
	admins := []common.Address{allowlist.TestAdminAddr}
	enableds := []common.Address{allowlist.TestEnabledAddr}
	tests := map[string]testutils.ConfigEqualTest{
		"non-nil config and nil other": {
			Config:   NewConfig(big.NewInt(3), admins, enableds),
			Other:    nil,
			Expected: false,
		},
		"different type": {
			Config:   NewConfig(big.NewInt(3), admins, enableds),
			Other:    precompileconfig.NewNoopStatefulPrecompileConfig(),
			Expected: false,
		},
		"different timestamp": {
			Config:   NewConfig(big.NewInt(3), admins, nil),
			Other:    NewConfig(big.NewInt(4), admins, nil),
			Expected: false,
		},
		"same config": {
			Config:   NewConfig(big.NewInt(3), admins, nil),
			Other:    NewConfig(big.NewInt(3), admins, nil),
			Expected: true,
		},
	}
	allowlist.EqualPrecompileWithAllowListTests(t, Module, tests)
}
