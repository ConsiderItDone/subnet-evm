package ibc

import (
	"fmt"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
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

type ibcBackend struct {
	db             database.Database
	snowCtx        *snow.Context
	signatureCache *cache.LRU[ids.ID, [bls.SignatureLen]byte]
}

// NewIbcBackend creates a new WarpBackend, and initializes the signature cache and message tracking database.
func NewIbcBackend(snowCtx *snow.Context, db database.Database, signatureCacheSize int) IbcBackend {
	return &ibcBackend{
		db:             db,
		snowCtx:        snowCtx,
		signatureCache: &cache.LRU[ids.ID, [bls.SignatureLen]byte]{Size: signatureCacheSize},
	}
}

func (ib *ibcBackend) AddMessage(payload []byte, key ids.ID) error {
	// Create a new unsigned message and add it to the warp backend.
	unsignedMessage, err := warp.NewUnsignedMessage(ib.snowCtx.ChainID, ib.snowCtx.ChainID, payload)
	if err != nil {
		return fmt.Errorf("failed to create unsigned message: %w", err)
	}

	if err := ib.db.Put(key[:], unsignedMessage.Bytes()); err != nil {
		return fmt.Errorf("failed to put in db: %w", err)
	}

	var signature [bls.SignatureLen]byte
	sig, err := ib.snowCtx.WarpSigner.Sign(unsignedMessage)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	copy(signature[:], sig)
	ib.signatureCache.Put(key, signature)
	return nil
}

func (ib *ibcBackend) GetProof(key ids.ID) ([bls.SignatureLen]byte, error) {
	// get proof from cache
	if sig, ok := ib.signatureCache.Get(key); ok {
		return sig, nil
	}

	unsignedMessageBytes, err := ib.db.Get(key[:])
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to get message %s from db: %w", key.String(), err)
	}

	unsignedMessage, err := warp.ParseUnsignedMessage(unsignedMessageBytes)
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to parse unsigned message %s: %w", key.String(), err)
	}

	var signature [bls.SignatureLen]byte
	sig, err := ib.snowCtx.WarpSigner.Sign(unsignedMessage)
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to sign message: %w", err)
	}

	copy(signature[:], sig)
	ib.signatureCache.Put(key, signature)
	return signature, nil
}
