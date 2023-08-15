// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"context"
	"testing"
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	// Import the solidity package, so that ginkgo maps out the tests declared within the package
	"github.com/ava-labs/avalanchego/api/health"
	_ "github.com/ava-labs/subnet-evm/tests/precompile/solidity"
	"github.com/ava-labs/subnet-evm/tests/utils"
)

func TestE2E(t *testing.T) {
	utils.RegisterNodeRun()
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "subnet-evm precompile ginkgo test suite")
}

func TestIBC(t *testing.T) {
	t.Log("starting avalanche node")
	cmd, err := utils.RunCommand("../../scripts/run.sh")
	require.NoError(t, err)
	defer cmd.Stop()

	healthClient := health.NewClient(utils.DefaultLocalNodeURI)
	healthy, err := health.AwaitReady(context.Background(), healthClient, 5*time.Second, nil)
	require.NoError(t, err)
	require.True(t, healthy)
	t.Log("avalanche node started")

	t.Run("part a", func(t *testing.T) {
		t.Run("create chain", utils.RunTestIbcInit)
		t.Run("create clients", utils.RunTestIbcCreateClient)
		t.Run("connection open init", utils.RunTestIbcConnectionOpenInit)
		t.Run("connection open ack", utils.RunTestIbcConnectionOpenAck)
		t.Run("channel open init", utils.RunTestIncChannelOpenInit)
		t.Run("channel open ack", utils.RunTestIncChannelOpenAck)
	})

	t.Run("part b", func(t *testing.T) {
		t.Run("create chain", utils.RunTestIbcInit)
		t.Run("create clients", utils.RunTestIbcCreateClient)
		t.Run("connection open try", utils.RunTestIbcConnectionOpenTry)
		t.Run("connection open confirm", utils.RunTestIbcConnectionOpenConfirm)
	})
}
