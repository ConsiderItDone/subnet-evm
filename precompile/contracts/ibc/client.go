package ibc

import (
	"fmt"
	"math/big"

	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
)

type CreateClientInput struct {
	ClientType     string
	ClientState    []byte
	ConsensusState []byte
}

type UpdateClientInput struct {
	ClientID      string
	ClientMessage []byte
}

type UpgradeClientInput struct {
	ClientID              string
	UpgradePath           []byte
	UpgradedClien         []byte
	UpgradedConsState     []byte
	ProofUpgradeClient    []byte
	ProofUpgradeConsState []byte
}

// UnpackCreateClientInput attempts to unpack [input] as CreateClientInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackCreateClientInput(input []byte) (CreateClientInput, error) {
	inputStruct := CreateClientInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "createClient", input)

	return inputStruct, err
}

// PackCreateClient packs [inputStruct] of type CreateClientInput into the appropriate arguments for createClient.
func PackCreateClient(inputStruct CreateClientInput) ([]byte, error) {
	return IBCABI.Pack("createClient", inputStruct.ClientType, inputStruct.ClientState, inputStruct.ConsensusState)
}

// PackCreateClientOutput attempts to pack given clientID of type string
// to conform the ABI outputs.
func PackCreateClientOutput(clientID string) ([]byte, error) {
	return IBCABI.PackOutput("createClient", clientID)
}

func getStoredNextClientSeq(db contract.StateDB) *big.Int {
	val := db.GetState(ContractAddress, nextClientSeqStorageKey)
	return val.Big()
}

func storeNextClientSeq(db contract.StateDB, value *big.Int) error {
	if value == nil {
		return fmt.Errorf("client seq cannot be nil")
	}
	db.SetState(ContractAddress, nextClientSeqStorageKey, common.BigToHash(value))

	return nil
}

func generateClientIdentifier(db contract.StateDB, clientType string) (string, error) {
	nextClientId := getStoredNextClientSeq(db)
	clientId := fmt.Sprintf("%s-%d", clientType, nextClientId.Int64())

	nextClientId.Add(nextClientId, big.NewInt(1))
	err := storeNextClientSeq(db, nextClientId)
	if err != nil {
		return "", err
	}

	return clientId, nil
}

func storeClientState(db contract.StateDB, clientState *ibctm.ClientState) error {
	bz, err := clientState.Marshal()
	if err != nil {
		return err
	}

	// storage layout for bytes
	// first 32 bytes - length of bytes
	// next 32 byte slices - payload
	if len(bz) <= 31 {

	}
	return nil
}

func createClient(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, CreateClientGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the CreateClientInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackCreateClientInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	clientType := inputStruct.ClientType

	// supports only Tendermint for now
	if clientType != exported.Tendermint {
		return nil, remainingGas, ErrWrongClientType
	}

	clientId, err := generateClientIdentifier(accessibleState.GetStateDB(), clientType)
	if err != nil {
		return nil, remainingGas, err
	}

	// emit event
	topics := make([]common.Hash, 1)
	topics[0] = GeneratedClientIdentifier.ID
	data, err := GeneratedClientIdentifier.Inputs.Pack(clientId)
	if err != nil {
		return nil, remainingGas, ErrWrongClientType
	}
	blockNumber := accessibleState.GetBlockContext().Number().Uint64()
	accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

	//clientState := &ibctm.ClientState{}
	//err = clientState.Unmarshal(inputStruct.ClientState)
	//if err != nil {
	//	return nil, remainingGas, err
	//}
	//
	//consensusState := &ibctm.ConsensusState{}
	//err = consensusState.Unmarshal(inputStruct.ConsensusState)
	//if err != nil {
	//	return nil, remainingGas, err
	//}

	packedOutput, err := PackCreateClientOutput(clientId)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackUpdateClientInput attempts to unpack [input] as UpdateClientInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackUpdateClientInput(input []byte) (UpdateClientInput, error) {
	inputStruct := UpdateClientInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "updateClient", input)

	return inputStruct, err
}

// PackUpdateClient packs [inputStruct] of type UpdateClientInput into the appropriate arguments for updateClient.
func PackUpdateClient(inputStruct UpdateClientInput) ([]byte, error) {
	return IBCABI.Pack("updateClient", inputStruct.ClientID, inputStruct.ClientMessage)
}

func updateClient(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, UpdateClientGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the UpdateClientInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackUpdateClientInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	if err := _updateClient(&callOpts[UpdateClientInput]{
		accessibleState: accessibleState,
		caller:          caller,
		addr:            addr,
		suppliedGas:     suppliedGas,
		readOnly:        readOnly,
		args:            inputStruct,
	}); err != nil {
		return nil, remainingGas, err
	}

	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackUpgradeClientInput attempts to unpack [input] as UpgradeClientInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackUpgradeClientInput(input []byte) (UpgradeClientInput, error) {
	inputStruct := UpgradeClientInput{}
	err := IBCABI.UnpackInputIntoInterface(&inputStruct, "upgradeClient", input)

	return inputStruct, err
}

// PackUpgradeClient packs [inputStruct] of type UpgradeClientInput into the appropriate arguments for upgradeClient.
func PackUpgradeClient(inputStruct UpgradeClientInput) ([]byte, error) {
	return IBCABI.Pack("upgradeClient", inputStruct.ClientID, inputStruct.UpgradePath, inputStruct.UpgradedClien, inputStruct.UpgradedConsState, inputStruct.ProofUpgradeClient, inputStruct.ProofUpgradeConsState)
}

func upgradeClient(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, UpgradeClientGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the UpgradeClientInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackUpgradeClientInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	if err := _upgradeClient(&callOpts[UpgradeClientInput]{
		accessibleState: accessibleState,
		caller:          caller,
		addr:            addr,
		suppliedGas:     suppliedGas,
		readOnly:        readOnly,
		args:            inputStruct,
	}); err != nil {
		return nil, remainingGas, err
	}

	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}
