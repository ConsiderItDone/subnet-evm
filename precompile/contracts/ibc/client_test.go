package ibc

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/core/state"
)

func TestClientSeqStorage(t *testing.T) {
	stateDB := state.NewTestStateDB(t)

	seq := getStoredNextClientSeq(stateDB)
	require.Equal(t, seq.String(), "0")

	err := storeNextClientSeq(stateDB, big.NewInt(5))
	require.NoError(t, err)

	seq = getStoredNextClientSeq(stateDB)
	require.Equal(t, seq.String(), "5")
}
