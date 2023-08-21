package ibc

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ava-labs/subnet-evm/precompile/contract"
)

func getChunks(db contract.StateDB, addr common.Address, slot common.Hash) []common.Hash {
	firstChunk := db.GetState(addr, slot)
	firstChunkBn := firstChunk.Big()
	hasMultipleChunks := firstChunkBn.Bit(0) == 1

	if !hasMultipleChunks {
		return []common.Hash{firstChunk}
	}

	// get length
	firstChunkBn.Sub(firstChunkBn, big.NewInt(1))
	firstChunkBn.Div(firstChunkBn, big.NewInt(2))

	length := firstChunkBn.Int64()
	chunkAmount := chunkSize(int(length))

	chunks := make([]common.Hash, 0, chunkAmount+1)
	chunks = append(chunks, firstChunk)
	hash := crypto.Keccak256Hash(slot.Bytes())
	for i := 1; i < chunkAmount; i++ {
		slotBn := hash.Big()
		slotBn.Add(slotBn, big.NewInt(int64(i-1)))
		key := common.BytesToHash(slotBn.Bytes())

		data := db.GetState(addr, key)
		chunks = append(chunks, data)
	}

	return chunks
}

func joinChunks(chunks []common.Hash) []byte {
	var chunk common.Hash
	chunk = chunks[0]
	if len(chunks) == 1 {
		length := new(big.Int).SetBytes(chunk[31:])
		data := chunk[:length.Int64()/2]

		return data
	}

	length := chunk.Big()
	length.Sub(length, big.NewInt(1))
	length.Div(length, big.NewInt(2))

	result := make([]byte, 0, length.Int64())

	// skip first chunk because it contains only length bytes
	for i := 1; i < len(chunks); i++ {
		chunk = chunks[i]

		var lim int64
		if length.Int64() > common.HashLength {
			lim = common.HashLength
		} else {
			lim = length.Int64()
		}

		result = append(result, chunk.Bytes()[:lim]...)
		length.Sub(length, big.NewInt(common.HashLength))
	}

	return result
}

func getState(db contract.StateDB, addr common.Address, slot common.Hash) []byte {
	return joinChunks(getChunks(db, addr, slot))
}

func GetState(db contract.StateDB, slot common.Hash) ([]byte, error) {
	state := getState(db, ContractAddress, slot)
	if len(state) == 0 {
		return nil, ErrEmptyState
	}
	return state, nil
}

func chunkSize(length int) int {
	quotient := length / common.HashLength
	remainder := length % common.HashLength

	// 1 byte reserved for length of data encoded into payload
	if quotient == 0 && remainder <= common.HashLength-1 {
		return 1
	}

	if remainder > 0 {
		// 1 chunk for byte length, 1 chunk for remain bytes (less than 32)
		return quotient + 2
	}

	// we have exactly 32/64/96 bytes + 1 chunk for byte length
	return quotient + 1
}

func splitState(data []byte) []common.Hash {
	chunkSize := chunkSize(len(data))

	var chunk common.Hash
	chunks := make([]common.Hash, 0, chunkSize)

	lenBytes := math.U256Bytes(new(big.Int).SetUint64(uint64(len(data) * 2)))

	// edge case when we have len(data) <= 31 bytes
	if chunkSize == 1 {
		chunkBytes := common.RightPadBytes(data, 32)
		copy(chunkBytes[31:], lenBytes[31:])
		chunks = append(chunks, common.BytesToHash(chunkBytes))

		return chunks
	}

	lenBytes = math.U256Bytes(new(big.Int).SetUint64(uint64(len(data)*2 + 1)))
	chunks = append(chunks, common.BytesToHash(lenBytes))

	for len(data) >= common.HashLength {
		chunk, data = common.BytesToHash(data[:common.HashLength]), data[common.HashLength:]
		chunks = append(chunks, chunk)
	}
	if len(data) > 0 {
		chunks = append(chunks, common.BytesToHash(common.RightPadBytes(data, 32)))
	}

	return chunks
}

func setState(db contract.StateDB, addr common.Address, slot common.Hash, state []byte) error {
	hashes := splitState(state)
	hash := crypto.Keccak256Hash(slot.Bytes())
	for i, data := range hashes {
		var key common.Hash
		if i == 0 {
			key = slot
		} else {
			key = common.BytesToHash(new(big.Int).Add(hash.Big(), big.NewInt(int64(i-1))).Bytes())
		}
		db.SetState(addr, key, data)
	}
	return nil
}

func SetState(db contract.StateDB, slot common.Hash, obj Marshaler) error {
	state, err := obj.Marshal()
	if err != nil {
		return err
	}
	setState(db, ContractAddress, slot, state)
	return nil
}
