package ibc

import (
	"fmt"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
)

// IbcBackend tracks BLS signature
type IbcBackend interface {
	// AddMessage adds payload to the backend database
	AddMessage(payload []byte, key ids.ID) error

	// GetProof returns the proof of the requested key.
	GetProof(key ids.ID) ([bls.SignatureLen]byte, error)
}

type backend struct {
	networkID     uint32
	sourceChainID ids.ID
	db            database.Database
	warpSigner    warp.Signer

	signatureCache *cache.LRU[ids.ID, [bls.SignatureLen]byte]
}

// NewBackend creates a new WarpBackend, and initializes the signature cache and message tracking database.
func NewBackend(networkID uint32, sourceChainID ids.ID, warpSigner warp.Signer, db database.Database, cacheSize int) IbcBackend {
	return &backend{
		networkID:      networkID,
		sourceChainID:  sourceChainID,
		db:             db,
		warpSigner:     warpSigner,
		signatureCache: &cache.LRU[ids.ID, [bls.SignatureLen]byte]{Size: cacheSize},
	}
}

func (b *backend) AddMessage(payload []byte, key ids.ID) error {
	// Create a new unsigned message and add it to the warp backend.
	unsignedMessage, err := warp.NewUnsignedMessage(b.networkID, b.sourceChainID, payload)
	if err != nil {
		return fmt.Errorf("failed to create unsigned message: %w", err)
	}

	if err := b.db.Put(key[:], unsignedMessage.Bytes()); err != nil {
		return fmt.Errorf("failed to put in db: %w", err)
	}

	var signature [bls.SignatureLen]byte
	sig, err := b.warpSigner.Sign(unsignedMessage)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	copy(signature[:], sig)
	b.signatureCache.Put(key, signature)
	return nil
}

func (b *backend) GetProof(key ids.ID) ([bls.SignatureLen]byte, error) {
	// get proof from cache
	if sig, ok := b.signatureCache.Get(key); ok {
		return sig, nil
	}

	unsignedMessageBytes, err := b.db.Get(key[:])
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to get message %s from db: %w", key.String(), err)
	}

	unsignedMessage, err := warp.ParseUnsignedMessage(unsignedMessageBytes)
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to parse unsigned message %s: %w", key.String(), err)
	}

	var signature [bls.SignatureLen]byte
	sig, err := b.warpSigner.Sign(unsignedMessage)
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to sign message: %w", err)
	}

	copy(signature[:], sig)
	b.signatureCache.Put(key, signature)
	return signature, nil
}
