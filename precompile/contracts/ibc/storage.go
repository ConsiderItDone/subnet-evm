package ibc

import (
	"encoding/hex"
	"fmt"
	"math/big"

	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ava-labs/subnet-evm/precompile/contract"
)

func split(data []byte) []common.Hash {
	chunkSize := chunkSize(data)

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
		chunks = append(chunks, common.BytesToHash(data[:len(data)]))
	}

	return chunks
}

func join(chunks []common.Hash) []byte {
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

func chunkSize(data []byte) int {
	dataLen := len(data)
	quotient := dataLen / common.HashLength
	remainder := dataLen % common.HashLength

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

func storeBytes(db contract.StateDB, slot common.Hash, input []byte) {
	hashes := split(input)
	hash := crypto.Keccak256Hash(slot.Bytes())
	for i, data := range hashes {
		var key common.Hash
		if i == 0 {
			key = slot
		} else {
			key = common.BytesToHash(new(big.Int).Add(hash.Big(), big.NewInt(int64(i-1))).Bytes())
		}
		db.SetState(ContractAddress, key, data)
	}
}

func getBytesChunks(db contract.StateDB, slot common.Hash) []common.Hash {
	firstChunk := db.GetState(ContractAddress, slot)
	firstChunkBn := firstChunk.Big()
	hasMultipleChunks := firstChunkBn.Bit(0) == 1

	if !hasMultipleChunks {
		return []common.Hash{firstChunk}
	}

	// get length
	firstChunkBn.Sub(firstChunkBn, big.NewInt(1))
	firstChunkBn.Div(firstChunkBn, big.NewInt(2))

	length := firstChunkBn.Int64()
	chunkAmount := length / common.HashLength

	chunks := make([]common.Hash, 0, chunkAmount)
	hash := crypto.Keccak256Hash(slot.Bytes())
	for i := 1; i <= int(chunkAmount); i++ {
		slotBn := hash.Big()
		slotBn.Add(slotBn, big.NewInt(int64(i-1)))
		key := common.BytesToHash(slotBn.Bytes())

		data := db.GetState(ContractAddress, key)
		chunks = append(chunks, data)
	}

	return chunks
}
func getBytes(db contract.StateDB, slot common.Hash) []byte {
	chunks := getBytesChunks(db, slot)
	return join(chunks)
}

func ClientSlot(clientID string) common.Hash {
	return CalculateKey(host.FullClientStateKey(clientID))
}

func packBytesToEVMStorage(input []byte) ([]byte, error) {
	byteLength := len(input)                   // Get the byte length of the input slice
	packedBytes := make([]byte, byteLength+32) // Allocate a slice to hold the packed bytes

	// Convert the byte length to a big-endian hexadecimal string
	lengthHex := fmt.Sprintf("%064x", byteLength)
	lengthBytes, err := hex.DecodeString(lengthHex)
	if err != nil {
		return nil, err
	}

	// Copy the length bytes to the packed bytes slice starting from index 0
	copy(packedBytes, lengthBytes)

	// Copy the input bytes to the packed bytes slice starting from index 32
	copy(packedBytes[32:], input)

	return packedBytes, nil
}

func main() {
	input := []byte("Hello, World!") // The byte slice to be packed
	packedBytes, err := packBytesToEVMStorage(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Packed Bytes: %x\n", packedBytes)
}
