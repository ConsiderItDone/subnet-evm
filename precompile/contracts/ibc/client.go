package ibc

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

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

func calculateKey(path []byte) string {
	return crypto.Keccak256Hash(crypto.Keccak256Hash(path).Bytes()).Hex()
}

func storeClientState(db contract.StateDB, clientId string, clientState *ibctm.ClientState) error {
	bz, err := clientState.Marshal()
	if err != nil {
		return err
	}
	key := calculateKey(host.FullClientStateKey(clientId))
	db.SetPrecompileState(common.BytesToAddress([]byte(key)), bz)

	return nil
}

func getClientState(db contract.StateDB, clientId string) (*ibctm.ClientState, bool, error) {
	key := calculateKey(host.FullClientStateKey(clientId))
	bz := db.GetPrecompileState(common.BytesToAddress([]byte(key)))

	if len(bz) == 0 {
		return nil, false, nil
	}

	clientState := &ibctm.ClientState{}
	err := clientState.Unmarshal(bz)
	if err != nil {
		return nil, false, err
	}

	return clientState, true, nil
}

func storeConsensusState(db contract.StateDB, clientId string, consensusState *ibctm.ConsensusState, height exported.Height) error {
	bz, err := consensusState.Marshal()
	if err != nil {
		return err
	}
	key := calculateKey(host.FullConsensusStateKey(clientId, height))
	db.SetPrecompileState(common.BytesToAddress([]byte(key)), bz)

	return nil
}

func getConsensusState(db contract.StateDB, clientId string, height exported.Height) (*ibctm.ConsensusState, bool, error) {
	key := calculateKey(host.FullConsensusStateKey(clientId, height))
	found := db.Exist(common.BytesToAddress([]byte(key)))
	if !found {
		return nil, false, nil
	}

	bz := db.GetPrecompileState(common.BytesToAddress([]byte(key)))
	consensusState := &ibctm.ConsensusState{}
	err := consensusState.Unmarshal(bz)
	if err != nil {
		return nil, false, err
	}

	return consensusState, true, nil
}

//func storeClientState(db contract.StateDB, clientState *ibctm.ClientState) error {
//	bz, err := clientState.Marshal()
//	if err != nil {
//		return err
//	}
//
//	// storage layout for bytes
//	// first 32 bytes - length of bytes
//	// next 32 byte slices - payload
//	if len(bz) <= 31 {
//
//	}
//	return nil
//}

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

	// generate clientID
	clientId, err := generateClientIdentifier(accessibleState.GetStateDB(), clientType)
	if err != nil {
		return nil, remainingGas, err
	}

	clientState := &ibctm.ClientState{}
	err = clientState.Unmarshal(inputStruct.ClientState)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("error unmarshalling client state: %w", err)
	}

	err = storeClientState(accessibleState.GetStateDB(), clientId, clientState)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("error storing client state: %w", err)
	}

	consensusState := &ibctm.ConsensusState{}
	err = consensusState.Unmarshal(inputStruct.ConsensusState)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling consensus state: %w", err)
	}

	err = storeConsensusState(accessibleState.GetStateDB(), clientId, consensusState, clientState.GetLatestHeight())
	if err != nil {
		return nil, remainingGas, fmt.Errorf("error storing consensus state: %w", err)
	}

	// emit event
	topics, data, err := IBCABI.PackEvent("ClientCreated", clientId)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("error packing event: %w", err)
	}
	blockNumber := accessibleState.GetBlockContext().Number().Uint64()
	accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

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

func _updateClient(opts *callOpts[UpdateClientInput]) error {
	stateDB := opts.accessibleState.GetStateDB()

	clientStatePath := fmt.Sprintf("clients/%s/clientState", opts.args.ClientID)
	found := stateDB.Exist(common.BytesToAddress([]byte(clientStatePath)))
	if !found {
		return fmt.Errorf("cannot update client with ID %s", opts.args.ClientID)
	}

	clientStateByte := stateDB.GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientState := &ibctm.ClientState{}
	if err := clientState.Unmarshal(clientStateByte); err != nil {
		return err
	}

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", opts.args.ClientID, clientState.GetLatestHeight())
	found = stateDB.Exist(common.BytesToAddress([]byte(consensusStatePath)))
	if !found {
		return fmt.Errorf("cannot update consensusState with ID %s", opts.args.ClientID)
	}

	clientMessage := &ibctm.Header{}
	if err := clientMessage.Unmarshal(opts.args.ClientMessage); err != nil {
		return fmt.Errorf("error unmarshalling client state file: %w", err)
	}

	consensusState := &ibctm.ConsensusState{
		Timestamp:          clientMessage.GetTime(),
		Root:               commitmenttypes.NewMerkleRoot(clientMessage.Header.GetAppHash()),
		NextValidatorsHash: clientMessage.Header.NextValidatorsHash,
	}
	// store ConsensusStateBytes
	consensusStateByte, err := consensusState.Marshal()
	if err != nil {
		return errors.New("consensusState marshaler error")
	}

	consensusStatePath = fmt.Sprintf("clients/%s/consensusStates/%s", opts.args.ClientID, clientMessage.GetHeight())
	stateDB.SetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)), consensusStateByte)
	return nil
}

func _upgradeClient(opts *callOpts[UpgradeClientInput]) error {
	stateDB := opts.accessibleState.GetStateDB()

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	upgradedClient := &ibctm.ClientState{}
	if err := upgradedClient.Unmarshal(opts.args.UpgradedClien); err != nil {
		return fmt.Errorf("error unmarshalling upgraded client: %w", err)
	}

	upgradedConsState := &ibctm.ConsensusState{}
	if err := upgradedConsState.Unmarshal(opts.args.UpgradedConsState); err != nil {
		return fmt.Errorf("error unmarshalling upgraded ConsensusState: %w", err)
	}

	clientStatePath := fmt.Sprintf("clients/%s/clientState", opts.args.ClientID)
	clientStateByte := stateDB.GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientStateExp := clienttypes.MustUnmarshalClientState(marshaler, clientStateByte)
	clientState, ok := clientStateExp.(*ibctm.ClientState)
	if !ok {
		return fmt.Errorf("error unmarshalling client state file")
	}

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", opts.args.ClientID, clientState.GetLatestHeight())
	consensusStateByte := stateDB.GetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)))
	consensusStateExp := clienttypes.MustUnmarshalConsensusState(marshaler, consensusStateByte)
	consensusState, ok := consensusStateExp.(*ibctm.ConsensusState)
	if !ok {
		return fmt.Errorf("error unmarshalling consensus state file")
	}

	if len(clientState.UpgradePath) == 0 {
		return errors.New("cannot upgrade client, no upgrade path set")
	}

	// last height of current counterparty chain must be client's latest height
	lastHeight := clientState.GetLatestHeight()

	if !upgradedClient.GetLatestHeight().GT(lastHeight) {
		return fmt.Errorf("upgraded client height %s must be at greater than current client height %s", upgradedClient.GetLatestHeight(), lastHeight)
	}

	// unmarshal proofs
	var merkleProofClient, merkleProofConsState commitmenttypes.MerkleProof
	if err := marshaler.Unmarshal(opts.args.ProofUpgradeClient, &merkleProofClient); err != nil {
		return fmt.Errorf("could not unmarshal client merkle proof: %v", err)
	}
	if err := marshaler.Unmarshal(opts.args.ProofUpgradeConsState, &merkleProofConsState); err != nil {
		return fmt.Errorf("could not unmarshal consensus state merkle proof: %v", err)
	}

	// Verify client proof
	client := upgradedClient.ZeroCustomFields()
	bz, err := marshaler.MarshalInterface(client)
	if err != nil {
		return fmt.Errorf("could not marshal client state: %v", err)
	}
	// copy all elements from upgradePath except final element
	clientPath := make([]string, len(clientState.UpgradePath)-1)
	copy(clientPath, clientState.UpgradePath)

	// append lastHeight and `upgradedClient` to last key of upgradePath and use as lastKey of clientPath
	// this will create the IAVL key that is used to store client in upgrade store
	lastKey := clientState.UpgradePath[len(clientState.UpgradePath)-1]
	appendedKey := fmt.Sprintf("%s/%d/%s", lastKey, lastHeight.GetRevisionHeight(), upgradetypes.KeyUpgradedClient)

	clientPath = append(clientPath, appendedKey)

	// construct clientState Merkle path
	upgradeClientPath := commitmenttypes.NewMerklePath(clientPath...)

	if err := merkleProofClient.VerifyMembership(clientState.ProofSpecs, consensusState.GetRoot(), upgradeClientPath, bz); err != nil {
		return fmt.Errorf("client state proof failed. Path: %s, err: %v", upgradeClientPath.Pretty(), err)
	}

	// Verify consensus state proof
	bz, err = marshaler.MarshalInterface(upgradedConsState)
	if err != nil {
		return fmt.Errorf("could not marshal consensus state: %v", err)
	}

	// copy all elements from upgradePath except final element
	consPath := make([]string, len(clientState.UpgradePath)-1)
	copy(consPath, clientState.UpgradePath)

	// append lastHeight and `upgradedClient` to last key of upgradePath and use as lastKey of clientPath
	// this will create the IAVL key that is used to store client in upgrade store
	lastKey = clientState.UpgradePath[len(clientState.UpgradePath)-1]
	appendedKey = fmt.Sprintf("%s/%d/%s", lastKey, lastHeight.GetRevisionHeight(), upgradetypes.KeyUpgradedConsState)

	consPath = append(consPath, appendedKey)
	// construct consensus state Merkle path
	upgradeConsStatePath := commitmenttypes.NewMerklePath(consPath...)

	if err := merkleProofConsState.VerifyMembership(clientState.ProofSpecs, consensusState.GetRoot(), upgradeConsStatePath, bz); err != nil {
		return fmt.Errorf("consensus state proof failed. Path: %s", upgradeConsStatePath.Pretty())
	}

	newClientState := ibctm.NewClientState(
		upgradedClient.ChainId, clientState.TrustLevel, clientState.TrustingPeriod, upgradedClient.UnbondingPeriod,
		clientState.MaxClockDrift, upgradedClient.LatestHeight, upgradedClient.ProofSpecs, upgradedClient.UpgradePath,
	)

	if err := newClientState.Validate(); err != nil {
		return fmt.Errorf("updated client state failed basic validation")
	}

	newConsState := ibctm.NewConsensusState(
		upgradedConsState.Timestamp, commitmenttypes.NewMerkleRoot([]byte(ibctm.SentinelRoot)), upgradedConsState.NextValidatorsHash,
	)

	consensusStateByte, err = marshaler.MarshalInterface(newConsState)
	if err != nil {
		return errors.New("consensusState marshaler error")
	}
	consensusStatePath = fmt.Sprintf("clients/%s/consensusStates/%s", opts.args.ClientID, lastHeight)
	stateDB.SetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)), consensusStateByte)

	clientStateByte, err = marshaler.MarshalInterface(newClientState)
	if err != nil {
		return errors.New("clientState marshaler error")
	}
	clientStatePath = fmt.Sprintf("clients/%s/clientState", opts.args.ClientID)
	stateDB.SetPrecompileState(common.BytesToAddress([]byte(clientStatePath)), clientStateByte)
	return nil
}
