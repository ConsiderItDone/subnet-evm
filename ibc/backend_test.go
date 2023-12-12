package ibc

import (
	"testing"

	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/stretchr/testify/require"
)

var (
	networkID     uint32 = 54321
	sourceChainID        = ids.GenerateTestID()
	payload              = []byte("test")
	key                  = ids.GenerateTestID()
)

func TestProof(t *testing.T) {
	db := memdb.New()

	sk, err := bls.NewSecretKey()
	require.NoError(t, err)
	warpSigner := avalancheWarp.NewSigner(sk, networkID, sourceChainID)
	backend := NewBackend(networkID, sourceChainID, warpSigner, db, 500)

	// Add payload to backend
	err = backend.AddMessage(payload, key)
	require.NoError(t, err)

	// Verify that a proof is returned successfully, and compare to expected one.
	signature, err := backend.GetProof(key)
	require.NoError(t, err)

	unsignedMsg, err := avalancheWarp.NewUnsignedMessage(networkID, sourceChainID, payload)
	require.NoError(t, err)
	expectedSig, err := warpSigner.Sign(unsignedMsg)
	require.NoError(t, err)
	require.Equal(t, expectedSig, signature[:])
}
