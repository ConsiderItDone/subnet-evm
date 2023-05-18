package ibc

import (
	_ "embed"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/ethereum/go-ethereum/common"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

const (
	upgradeClientGas = uint64(1)
	updateClientGas  = uint64(1)
	createClientGas  = uint64(1)
)

// Singleton StatefulPrecompiledContract and signatures.
var (
	IbcGoPrecompile = createIbcGoPrecompile() // will be initialized by init function

	getCreateClientSignature  = contract.CalculateFunctionSelector("createClient(uint64,byte[],uint64,byte[])")
	getUpdateClientSignature  = contract.CalculateFunctionSelector("updateClient(uint64,byte[],uint64,byte[])")
	getUpgradeClientSignature = contract.CalculateFunctionSelector("upgradeClient(uint64,byte[],uint64,byte[],uint64,byte[],uint64,byte[],uint64,byte[],uint64,byte[])")
)


// createClient generates a new client identifier and isolated prefix store for the provided client state.
// The client state is responsible for setting any client-specific data in the store via the Initialize method.
// This includes the client state, initial consensus state and any associated metadata.
func createClient(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
		input
		8 byte             - clientStateLen
		clientStateLen     - clientState
		8 byte             - consensusStateLen
		consensusStateLen  - consensusState
	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, createClientGas); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	// Allow list is enabled and AllowFeeRecipients is a state-changer function.
	// This part of the code restricts the function to be called only by enabled/admin addresses in the allow list.
	// You can modify/delete this code if you don't want this function to be restricted by the allow list.
	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("non-enabled cannot call createClient: %s", caller)
	}

	nextClientSeq := uint64(0)
	if accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte("nextClientSeq"))) {
		b := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte("nextClientSeq")))
		nextClientSeq = binary.BigEndian.Uint64(b)
	}

	
	// ClientStateBytes
	carriage := uint64(0)
	clientStateLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8

	clientStateByte := getData(input, carriage, clientStateLen)

	carriage = carriage + clientStateLen

	clientState := &ibctm.ClientState{}
	err = clientState.Unmarshal(clientStateByte)
	if err != nil {
		return nil, createClientGas, fmt.Errorf("error unmarshalling client state file: %w", err)
	}

	clientID := fmt.Sprintf("%s-%d", clientState.ClientType(), nextClientSeq)
	nextClientSeq++

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, nextClientSeq)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte("nextClientSeq")), b)

	// ConsensusStateBytes
	consensusStateLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	consensusStateByte := getData(input, carriage, consensusStateLen)

	consensusState := &ibctm.ConsensusState{}
	err = consensusState.Unmarshal(consensusStateByte)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling client state file: %w", err)
	}
	// store ClientStateBytes
	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(clientStatePath)), clientStateByte)

	// store ConsensusStateBytes
	consensusStateByte, err = consensusState.Marshal()
	if err != nil {
		return nil, 0, errors.New("consensusState marshaler error")
	}

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", clientID, clientState.GetLatestHeight())
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)), consensusStateByte)
	return []byte(clientID), createClientGas, nil
}

func updateClient(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
		input
		8 byte             - clientIDLen
		clientIDLen byte   - clientID
		8 byte             - clientMessageLen
		clientMessageLen   - clientMessageByte
	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, updateClientGas); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	// Allow list is enabled and AllowFeeRecipients is a state-changer function.
	// This part of the code restricts the function to be called only by enabled/admin addresses in the allow list.
	// You can modify/delete this code if you don't want this function to be restricted by the allow list.
	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("non-enabled cannot call updateClient: %s", caller)
	}

	// clientId
	carriage := uint64(0)
	clientIDLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	clientID := string(getData(input, carriage, clientIDLen))
	carriage = carriage + clientIDLen

	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)
	found := accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(clientStatePath)))
	if !found {
		return nil, 0, fmt.Errorf("cannot update client with ID %s", clientID)
	}

	clientStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientState := &ibctm.ClientState{}
	err = clientState.Unmarshal(clientStateByte)
	if err != nil {
		return nil, 0, err
	}

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", clientID, clientState.GetLatestHeight())
	found = accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(consensusStatePath)))
	if !found {
		return nil, 0, fmt.Errorf("cannot update consensusState with ID %s", clientID)
	}

	// bytes clientMessage;
	clientMessageLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	clientMessageByte := getData(input, carriage, clientMessageLen)

	clientMessage := &ibctm.Header{}
	err = clientMessage.Unmarshal(clientMessageByte)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling client state file: %w", err)
	}

	consensusState := &ibctm.ConsensusState{
		Timestamp:          clientMessage.GetTime(),
		Root:               commitmenttypes.NewMerkleRoot(clientMessage.Header.GetAppHash()),
		NextValidatorsHash: clientMessage.Header.NextValidatorsHash,
	}
	// store ConsensusStateBytes
	consensusStateByte, err := consensusState.Marshal()
	if err != nil {
		return nil, 0, errors.New("consensusState marshaler error")
	}

	consensusStatePath = fmt.Sprintf("clients/%s/consensusStates/%s", clientID, clientMessage.GetHeight())
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)), consensusStateByte)
	return nil, updateClientGas, nil
}

func upgradeClient(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
		input
		8 byte                   - clientIDLen
		clientIDLenbyte          - clientID
		8 byte        			 - UpgradePathLen
		UpgradePathLen           - UpgradePath
		8 byte                   - upgradedClientLen
		upgradedClientLen byte   - upgradedClientByte
		8 byte                   - upgradedConsStateLen
		upgradedConsStateLen     - upgradedConsStateByte
		8 byte                   - proofUpgradeClientLen
		proofUpgradeClientLen    - proofUpgradeClientByte
		8 byte                   - proofUpgradeConsStateLen
		proofUpgradeConsStateLen - proofUpgradeConsStateByte
	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, upgradeClientGas); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	// Allow list is enabled and AllowFeeRecipients is a state-changer function.
	// This part of the code restricts the function to be called only by enabled/admin addresses in the allow list.
	// You can modify/delete this code if you don't want this function to be restricted by the allow list.
	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("non-enabled cannot call upgradeClient: %s", caller)
	}

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	// clientId
	carriage := uint64(0)
	clientIDLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	clientID := string(getData(input, carriage, clientIDLen))
	carriage = carriage + clientIDLen
	//upgradedClientByte
	upgradedClientLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	upgradedClientByte := getData(input, carriage, upgradedClientLen)
	carriage = carriage + upgradedClientLen

	upgradedClient := &ibctm.ClientState{}
	err = upgradedClient.Unmarshal(upgradedClientByte)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling upgraded client: %w", err)
	}

	//upgradedConsStateByte
	upgradedConsStateLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	upgradedConsStateByte := getData(input, carriage, upgradedConsStateLen)
	carriage = carriage + upgradedConsStateLen

	upgradedConsState := &ibctm.ConsensusState{}
	err = upgradedConsState.Unmarshal(upgradedConsStateByte)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling upgraded ConsensusState: %w", err)
	}

	//proofUpgradeClientByte
	proofUpgradeClientLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofUpgradeClientByte := getData(input, carriage, proofUpgradeClientLen)
	carriage = carriage + proofUpgradeClientLen

	//proofUpgradeConsStateByte
	proofUpgradeConsStateLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofUpgradeConsStateByte := getData(input, carriage, proofUpgradeConsStateLen)

	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)
	clientStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientStateExp := clienttypes.MustUnmarshalClientState(marshaler, clientStateByte)
	clientState, ok := clientStateExp.(*ibctm.ClientState)
	if !ok {
		return nil, 0, fmt.Errorf("error unmarshalling client state file")
	}
	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", clientID, clientState.GetLatestHeight())
	consensusStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)))
	consensusStateExp := clienttypes.MustUnmarshalConsensusState(marshaler, consensusStateByte)
	consensusState, ok := consensusStateExp.(*ibctm.ConsensusState)
	if !ok {
		return nil, 0, fmt.Errorf("error unmarshalling consensus state file")
	}

	if len(clientState.UpgradePath) == 0 {
		return nil, 0, errors.New("cannot upgrade client, no upgrade path set")
	}

	// last height of current counterparty chain must be client's latest height
	lastHeight := clientState.GetLatestHeight()

	if !upgradedClient.GetLatestHeight().GT(lastHeight) {
		return nil, 0, fmt.Errorf("upgraded client height %s must be at greater than current client height %s", upgradedClient.GetLatestHeight(), lastHeight)
	}

	// unmarshal proofs
	var merkleProofClient, merkleProofConsState commitmenttypes.MerkleProof
	if err := marshaler.Unmarshal(proofUpgradeClientByte, &merkleProofClient); err != nil {
		return nil, 0, fmt.Errorf("could not unmarshal client merkle proof: %v", err)
	}
	if err := marshaler.Unmarshal(proofUpgradeConsStateByte, &merkleProofConsState); err != nil {
		return nil, 0, fmt.Errorf("could not unmarshal consensus state merkle proof: %v", err)
	}

	// Verify client proof
	client := upgradedClient.ZeroCustomFields()
	bz, err := marshaler.MarshalInterface(client)
	if err != nil {
		return nil, 0, fmt.Errorf("could not marshal client state: %v", err)
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
		return nil, 0, fmt.Errorf("client state proof failed. Path: %s, err: %v", upgradeClientPath.Pretty(), err)
	}

	// Verify consensus state proof
	bz, err = marshaler.MarshalInterface(upgradedConsState)
	if err != nil {
		return nil, 0, fmt.Errorf("could not marshal consensus state: %v", err)
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
		return nil, 0, fmt.Errorf("consensus state proof failed. Path: %s", upgradeConsStatePath.Pretty())
	}

	newClientState := ibctm.NewClientState(
		upgradedClient.ChainId, clientState.TrustLevel, clientState.TrustingPeriod, upgradedClient.UnbondingPeriod,
		clientState.MaxClockDrift, upgradedClient.LatestHeight, upgradedClient.ProofSpecs, upgradedClient.UpgradePath,
	)

	if err := newClientState.Validate(); err != nil {
		return nil, 0, fmt.Errorf("updated client state failed basic validation")
	}

	newConsState := ibctm.NewConsensusState(
		upgradedConsState.Timestamp, commitmenttypes.NewMerkleRoot([]byte(ibctm.SentinelRoot)), upgradedConsState.NextValidatorsHash,
	)

	consensusStateByte, err = marshaler.MarshalInterface(newConsState)
	if err != nil {
		return nil, 0, errors.New("consensusState marshaler error")
	}
	consensusStatePath = fmt.Sprintf("clients/%s/consensusStates/%s", clientID, lastHeight)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)), consensusStateByte)

	clientStateByte, err = marshaler.MarshalInterface(newClientState)
	if err != nil {
		return nil, 0, errors.New("clientState marshaler error")
	}
	clientStatePath = fmt.Sprintf("clients/%s/clientState", clientID)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(clientStatePath)), clientStateByte)
	return nil, upgradeClientGas, nil
}

// getData returns a slice from the data based on the start and size and pads
// up to size with zero's. This function is overflow safe.
func getData(data []byte, start uint64, size uint64) []byte {
	length := uint64(len(data))
	if start > length {
		start = length
	}
	end := start + size
	if end > length {
		end = length
	}
	return common.RightPadBytes(data[start:end], int(size))
}

// createIbcGoPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.
// Access to the getters/setters is controlled by an allow list for [precompileAddr].
func createIbcGoPrecompile() contract.StatefulPrecompiledContract {
	enabledFuncs := allowlist.CreateAllowListFunctions(ContractAddress)

	enabledFuncs = append(enabledFuncs, 
		contract.NewStatefulPrecompileFunction(getCreateClientSignature, createClient), 
		contract.NewStatefulPrecompileFunction(getUpdateClientSignature, updateClient),
		contract.NewStatefulPrecompileFunction(getUpgradeClientSignature, upgradeClient),
	)

	// Construct the contract with no fallback function.
	contract, err := contract.NewStatefulPrecompileContract(nil, enabledFuncs)
	// TODO: Change this to be returned as an error after refactoring this precompile
	// to use the new precompile template.
	if err != nil {
		panic(err)
	}
	return contract
}
