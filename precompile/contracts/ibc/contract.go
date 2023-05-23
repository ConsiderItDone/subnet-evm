package ibc

import (
	_ "embed"
	"encoding/binary"
	"encoding/json"
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
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	hosttypes "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
)

const (
	upgradeClientGas = uint64(1)
	updateClientGas  = uint64(1)
	createClientGas  = uint64(1)

	connOpenInitGas    = uint64(1)
	connOpenTryGas     = uint64(1)
	connOpenConfirmGas = uint64(1)
)

// Singleton StatefulPrecompiledContract and signatures.
var (
	IbcGoPrecompile = createIbcGoPrecompile() // will be initialized by init function

	getCreateClientSignature    = contract.CalculateFunctionSelector("createClient(uint64,bytes,uint64,bytes)")
	getUpdateClientSignature    = contract.CalculateFunctionSelector("updateClient(uint64,bytes,uint64,bytes)")
	getUpgradeClientSignature   = contract.CalculateFunctionSelector("upgradeClient(uint64,bytes,uint64,bytes,uint64,bytes,uint64,bytes,uint64,bytes,uint64,bytes)")
	getConnOpenInitSignature    = contract.CalculateFunctionSelector("connOpenInit(uint64,bytes,uint64,bytes,uint64,bytes)")
	getConnOpenConfirmSignature = contract.CalculateFunctionSelector("connOpenConfirm(uint64,bytes,uint64,bytes,uint64,bytes)")
	getConnOpenTrySignature     = contract.CalculateFunctionSelector("connOpenTry(uint64,bytes,uint64,uint64,bytes,uint64,bytes,uint64,bytes,uint64,bytes,uint64,bytes,uint64,bytes,uint64,bytes,uint64,bytes)")
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

func ConnOpenInit(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
		8 byte                       - clientIDLen
		clientIDbyte                 - clientID
		8 byte                       - counterpartyLen
		counterpartybyte             - counterparty
		8 byte                       - versionLen
		versionbyte                  - Version
		8 byte                       - delayPeriod
	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, upgradeClientGas); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

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

	carriage := uint64(0)
	// clientId
	clientIDLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	clientID := string(getData(input, carriage, clientIDLen))
	carriage = carriage + clientIDLen

	//counterparty
	counterpartyLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	counterpartyByte := getData(input, carriage, counterpartyLen)
	carriage = carriage + counterpartyLen

	counterparty := &connectiontypes.Counterparty{}
	err = counterparty.Unmarshal(counterpartyByte)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling counterparty: %w", err)
	}

	//version
	versionLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	versionByte := getData(input, carriage, versionLen)
	carriage = carriage + versionLen

	version := &connectiontypes.Version{}
	err = version.Unmarshal(versionByte)

	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling version: %w", err)
	}

	//delayPeriod
	delayPeriod := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()

	versions := connectiontypes.GetCompatibleVersions()
	if versionLen != 0 {
		if !connectiontypes.IsSupportedVersion(connectiontypes.GetCompatibleVersions(), version) {
			return nil, 0, fmt.Errorf("%w : version is not supported", connectiontypes.ErrInvalidVersion)
		}

		versions = []exported.Version{version}
	}

	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)
	found := accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(clientStatePath)))
	if !found {
		return nil, 0, fmt.Errorf("cannot update client with ID %s", clientID)
	}

	nextConnSeq := uint64(0)
	if accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte("nextConnSeq"))) {
		b := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")))
		nextConnSeq = binary.BigEndian.Uint64(b)
	}
	connectionID := fmt.Sprintf("%s%d", "connection-", nextConnSeq)
	nextConnSeq++
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, nextConnSeq)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")), b)

	// connection defines chain A's ConnectionEnd
	connection := connectiontypes.NewConnectionEnd(connectiontypes.INIT, clientID, *counterparty, connectiontypes.ExportedVersionsToProto(versions), delayPeriod)

	connectionByte, err := marshaler.Marshal(&connection)
	if err != nil {
		return nil, 0, fmt.Errorf("connection marshaler error: %w", err)
	}
	connectionsPath := fmt.Sprintf("connections/%s", connectionID)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

	return []byte(connectionID), connOpenInitGas, nil

}

func ConnOpenTry(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
		input
		8 byte                       - counterpartyLen
		counterpartybyte             - counterparty
		8 byte                       - delayPeriod
		8 byte                       - clientIDLen
		clientIDbyte                 - clientID
		8 byte                       - clientStateLen
		clientStatebyte              - clientState
		8 byte                       - counterpartyVersionsLen
		counterpartyVersionsbyte     - []exported.Version
		8 byte                       - proofInitLen
		proofInitbyte                - []byte
		8 byte                       - proofClientLen
		proofClientbyte              - []byte
		8 byte                       - proofConsensusLen
		proofConsensusbyte           - []byte
		8 byte                       - proofHeightLen
		proofHeightbyte              - exported.Height
		8 byte                       - consensusHeightLen
		consensusHeightbyte          - exported.Height
	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, upgradeClientGas); err != nil {
		return nil, 0, fmt.Errorf("Error DeductGas err: %w", err)
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("non-enabled cannot call upgradeClient: %s", caller)
	}

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	carriage := uint64(0)

	//counterparty
	counterpartyLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	counterpartyByte := getData(input, carriage, counterpartyLen)
	carriage = carriage + counterpartyLen

	counterparty := &connectiontypes.Counterparty{}
	err = counterparty.Unmarshal(counterpartyByte)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling counterparty: %w", err)
	}

	//delayPeriod
	delayPeriod := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8

	// clientId
	clientIDLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	clientID := string(getData(input, carriage, clientIDLen))
	carriage = carriage + clientIDLen

	//clientState
	clientStateLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	clientStateByte := getData(input, carriage, clientStateLen)
	carriage = carriage + clientStateLen

	clientStateExp, err := clienttypes.UnmarshalClientState(marshaler, clientStateByte)
	clientState := clientStateExp.(*ibctm.ClientState)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling clientState: %w", err)
	}

	//counterpartyVersions
	counterpartyVersionsLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	counterpartyVersionsByte := getData(input, carriage, counterpartyVersionsLen)
	carriage = carriage + counterpartyVersionsLen

	counterpartyVersions := []*connectiontypes.Version{}
	err = json.Unmarshal(counterpartyVersionsByte, &counterpartyVersions)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling counterpartyVersions: %w", err)
	}

	//proofInitbyte
	proofInitLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofInitbyte := getData(input, carriage, proofInitLen)
	carriage = carriage + proofInitLen

	//proofClientbyte
	proofClientLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofClientbyte := getData(input, carriage, proofClientLen)
	carriage = carriage + proofClientLen

	//proofConsensusbyte
	proofConsensusLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	/* proofConsensusbyte := */ getData(input, carriage, proofConsensusLen)
	carriage = carriage + proofConsensusLen

	//proofHeight
	proofHeightLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofHeightbyte := getData(input, carriage, proofHeightLen)
	carriage = carriage + proofHeightLen

	proofHeight := &clienttypes.Height{}
	err = marshaler.UnmarshalInterface(proofHeightbyte, proofHeight)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	//consensusHeight
	consensusHeightLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	consensusHeightbyte := getData(input, carriage, consensusHeightLen)

	consensusHeight := &clienttypes.Height{}
	err = marshaler.UnmarshalInterface(consensusHeightbyte, consensusHeight)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling consensusHeight: %w", err)
	}

	nextConnSeq := uint64(0)
	if accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte("nextConnSeq"))) {
		b := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")))
		nextConnSeq = binary.BigEndian.Uint64(b)
	}
	connectionID := fmt.Sprintf("%s%d", "connection-", nextConnSeq)
	nextConnSeq++
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, nextConnSeq)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")), b)

	expectedCounterparty := connectiontypes.NewCounterparty(clientID, "", commitmenttypes.NewMerklePrefix([]byte("ibc")))
	expectedConnection := connectiontypes.NewConnectionEnd(connectiontypes.INIT, counterparty.ClientId, expectedCounterparty, counterpartyVersions, delayPeriod)

	// chain B picks a version from Chain A's available versions that is compatible
	// with Chain B's supported IBC versions. PickVersion will select the intersection
	// of the supported versions and the counterparty versions.
	version, err := connectiontypes.PickVersion(connectiontypes.GetCompatibleVersions(), connectiontypes.ProtoVersionsToExported(counterpartyVersions))
	if err != nil {
		return nil, 0, fmt.Errorf("Error PickVersion err: %w", err)
	}

	// connection defines chain B's ConnectionEnd
	connection := connectiontypes.NewConnectionEnd(connectiontypes.TRYOPEN, clientID, *counterparty, []*connectiontypes.Version{version}, delayPeriod)

	err = clientVerefication(connection, clientState, *proofHeight, accessibleState, marshaler, proofClientbyte)
	if err != nil {
		return nil, 0, fmt.Errorf("Error clientVerefication err: %w", err)
	}
	err = connectionVerefication(connection, expectedConnection, *proofHeight, accessibleState, marshaler, connectionID, proofInitbyte)
	if err != nil {
		return nil, 0, fmt.Errorf("Error connectionVerefication err: %w", err)
	}

	// err = consensusStateVerefication(connection, expectedConsensusState, proofHeight, accessibleState, marshaler, connectionID, proofClientbyte)
	// if err != nil {
	// 	return nil, 0, err
	// }

	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)
	found := accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(clientStatePath)))
	if !found {
		return nil, 0, clienttypes.ErrClientNotFound
	}

	conns, found := getClientConnectionPaths(marshaler, clientID, accessibleState)
	if !found {
		conns = []string{}
	}
	conns = append(conns, connectionID)
	clientPaths := connectiontypes.ClientPaths{Paths: conns}
	bz := marshaler.MustMarshal(&clientPaths)

	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress(hosttypes.ClientConnectionsKey(clientID)), bz)
	return []byte(connectionID), connOpenTryGas, nil
}

func ConnOpenAck(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
		input
		8 byte                       - connectionIDLen
		connectionIDbyte             - connectionID
		8 byte                       - clientStateLen
		clientStatebyte              - clientState
		8 byte                       - versionLen
		versionbyte                  - Version
		8 byte                       - counterpartyConnectionIDLen
		counterpartyConnectionIDbyte - counterpartyConnectionID
		8 byte                       - proofTryLen
		proofTrybyte                 - []byte
		8 byte                       - proofClientLen
		proofClientbyte              - []byte
		8 byte                       - proofConsensusLen
		proofConsensusbyte           - []byte
		8 byte                       - proofHeightLen
		proofHeightbyte              - exported.Height
		8 byte                       - consensusHeightLen
		consensusHeightbyte          - exported.Height
	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, upgradeClientGas); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("non-enabled cannot call upgradeClient: %s", caller)
	}

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	carriage := uint64(0)

	//proofClientbyte
	connectionIDLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	connectionID := string(getData(input, carriage, connectionIDLen))
	carriage = carriage + connectionIDLen

	//clientState
	clientStateLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	clientStateByte := getData(input, carriage, clientStateLen)
	carriage = carriage + clientStateLen

	clientStateExp, err := clienttypes.UnmarshalClientState(marshaler, clientStateByte)
	clientState := clientStateExp.(*ibctm.ClientState)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling clientState: %w", err)
	}

	//version
	versionLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	versionByte := getData(input, carriage, versionLen)
	carriage = carriage + versionLen

	version := connectiontypes.Version{}
	err = json.Unmarshal(versionByte, &version)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling counterpartyVersions: %w", err)
	}
	
	// counterpartyConnectionID
	counterpartyConnectionIDLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	counterpartyConnectionID := string(getData(input, carriage, counterpartyConnectionIDLen))
	carriage = carriage + counterpartyConnectionIDLen


	//proofTrybyte
	proofTryLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofTry := getData(input, carriage, proofTryLen)
	carriage = carriage + proofTryLen


	//proofClientbyte
	proofClientLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofClient := getData(input, carriage, proofClientLen)
	carriage = carriage + proofClientLen

	//proofClientbyte
	proofConsensusLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	/* proofConsensus := */ getData(input, carriage, proofConsensusLen)
	carriage = carriage + proofConsensusLen

	//proofHeight
	proofHeightLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofHeightbyte := getData(input, carriage, proofHeightLen)
	carriage = carriage + proofHeightLen

	proofHeight := &clienttypes.Height{}
	err = marshaler.UnmarshalInterface(proofHeightbyte, proofHeight)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	//consensusHeight
	consensusHeightLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	consensusHeightbyte := getData(input, carriage, consensusHeightLen)

	consensusHeight := &clienttypes.Height{}
	err = marshaler.UnmarshalInterface(consensusHeightbyte, consensusHeight)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling consensusHeight: %w", err)
	}


	connectionByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(connectionID)))
	connection := connectiontypes.ConnectionEnd{}
	err = marshaler.UnmarshalInterface(connectionByte, connection)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling connection id: %s, error: %w", connectionID, err)
	}

	// verify the previously set connection state
	if connection.State != connectiontypes.INIT {
		return nil, 0, fmt.Errorf("connection state is not INIT (got %s), error: %w", connection.State.String(), connectiontypes.ErrInvalidConnectionState)
	}

	// ensure selected version is supported
	if !connectiontypes.IsSupportedVersion(connectiontypes.ProtoVersionsToExported(connection.Versions), &version) {
		return nil, 0, fmt.Errorf("the counterparty selected version %s is not supported by versions selected on INIT, error: %w", version, connectiontypes.ErrInvalidConnectionState)
	}

	expectedCounterparty := connectiontypes.NewCounterparty(connection.ClientId, connectionID, commitmenttypes.NewMerklePrefix([]byte("ibc")))
	expectedConnection := connectiontypes.NewConnectionEnd(connectiontypes.TRYOPEN, connection.Counterparty.ClientId, expectedCounterparty, []*connectiontypes.Version{&version}, connection.DelayPeriod)

	if err := connectionVerefication(connection, expectedConnection, proofHeight, accessibleState, marshaler, counterpartyConnectionID, proofTry); err != nil {
		return nil, 0, err
	}

	// Check that ChainB stored the clientState provided in the msg
	if err := clientVerefication(connection, clientState, proofHeight, accessibleState, marshaler, proofClient); err != nil {
		return nil, 0, err
	}

	// // Ensure that ChainB has stored the correct ConsensusState for chainA at the consensusHeight
	// if err := k.VerifyClientConsensusState(
	// 	ctx, connection, proofHeight, consensusHeight, proofConsensus, expectedConsensusState,
	// ); err != nil {
	// 	return err
	// }

	// Update connection state to Open
	connection.State = connectiontypes.OPEN
	connection.Versions = []*connectiontypes.Version{&version}
	connection.Counterparty.ConnectionId = counterpartyConnectionID

	connectionByte, err = marshaler.Marshal(&connection)
	if err != nil {
		return nil, 0, errors.New("connection marshaler error")
	}
	connectionsPath := fmt.Sprintf("connections/%s", connectionID)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

	return nil, 0, nil
}

func ConnOpenConfirm(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
		input
		8 byte                       - connectionIDLen
		connectionIDbyte             - connectionID
		8 byte                       - proofAckLen
		proofAckbyte                 - []byte
		8 byte                       - proofHeightLen
		proofHeightbyte              - exported.Height
	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, upgradeClientGas); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("non-enabled cannot call upgradeClient: %s", caller)
	}

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	carriage := uint64(0)
	// clientId
	connectionIDLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	connectionID := string(getData(input, carriage, connectionIDLen))
	carriage = carriage + connectionIDLen

	//proofClientbyte
	proofAckLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofAckbyte := getData(input, carriage, proofAckLen)
	carriage = carriage + proofAckLen

	//counterparty
	proofHeightLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofHeightbyte := getData(input, carriage, proofHeightLen)

	proofHeight := &clienttypes.Height{}
	err = marshaler.Unmarshal(proofHeightbyte, proofHeight)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	connectionsPath := fmt.Sprintf("connections/%s", connectionID)

	exist := accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(connectionsPath)))
	if !exist {
		return nil, 0, fmt.Errorf("cannot find connection with path: %s", connectionsPath)
	}

	connectionByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(connectionsPath)))
	connection := &connectiontypes.ConnectionEnd{}
	marshaler.MustUnmarshal(connectionByte, connection)

	// Check that connection state on ChainB is on state: TRYOPEN
	if connection.State != connectiontypes.TRYOPEN {
		return nil, 0, fmt.Errorf("connection state is not TRYOPEN (got %s), err: %w", connection.State.String(), connectiontypes.ErrInvalidConnectionState)
	}

	// prefix := k.GetCommitmentPrefix()
	expectedCounterparty := connectiontypes.NewCounterparty(connection.ClientId, connectionID, commitmenttypes.NewMerklePrefix([]byte("ibc")))
	expectedConnection := connectiontypes.NewConnectionEnd(connectiontypes.OPEN, connection.Counterparty.ClientId, expectedCounterparty, connection.Versions, connection.DelayPeriod)

	err = connectionVerefication(*connection, expectedConnection, *proofHeight, accessibleState, marshaler, connectionID, proofAckbyte)
	if err != nil {
		return nil, 0, err
	}

	clientID := connection.GetClientID()

	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)
	clientStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientStateExp, err := clienttypes.UnmarshalClientState(marshaler, clientStateByte)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling client state file, err: %w", err)
	}
	clientState := clientStateExp.(*ibctm.ClientState)

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", clientID, clientState.GetLatestHeight())
	consensusStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)))
	consensusStateExp, err := clienttypes.UnmarshalConsensusState(marshaler, consensusStateByte)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling consensus state file, err: %w", err)
	}
	consensusState := consensusStateExp.(*ibctm.ConsensusState)

	merklePath := commitmenttypes.NewMerklePath(hosttypes.ConnectionPath(connectionID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return nil, 0, err
	}

	bz, err := marshaler.Marshal(&expectedConnection)
	if err != nil {
		return nil, 0, err
	}

	if clientState.GetLatestHeight().LT(*proofHeight) {
		return nil, 0, fmt.Errorf("client state height < proof height (%d < %d), please ensure the client has been updated", clientState.GetLatestHeight(), proofHeight)
	}

	var merkleProof commitmenttypes.MerkleProof
	if err := marshaler.Unmarshal(proofAckbyte, &merkleProof); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal proof into ICS 23 commitment merkle proof")
	}
	merkleProof.VerifyMembership(clientState.ProofSpecs, consensusState.GetRoot(), merklePath, bz)

	// Update ChainB's connection to Open
	connection.State = connectiontypes.OPEN

	connectionByte, err = marshaler.Marshal(connection)
	if err != nil {
		return nil, 0, errors.New("connection marshaler error")
	}
	connectionsPath = fmt.Sprintf("connections/%s", connectionID)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

	return nil, connOpenConfirmGas, err
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
		contract.NewStatefulPrecompileFunction(getConnOpenInitSignature, ConnOpenInit),
		contract.NewStatefulPrecompileFunction(getConnOpenTrySignature, ConnOpenTry),
		contract.NewStatefulPrecompileFunction(getConnOpenConfirmSignature, ConnOpenConfirm),
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

func clientVerefication(
	connection connectiontypes.ConnectionEnd,
	clientState exported.ClientState,
	proofHeight exported.Height,
	accessibleState contract.AccessibleState,
	marshaler *codec.ProtoCodec,
	proofClientbyte []byte,
) error {

	clientID := connection.GetClientID()

	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)

	clientStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientStateExp, err := clienttypes.UnmarshalClientState(marshaler, clientStateByte)
	if err != nil {
		return fmt.Errorf("error unmarshalling client state file, err: %w", err)
	}
	targetClientState := clientStateExp.(*ibctm.ClientState)

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", clientID, targetClientState.GetLatestHeight())
	consensusStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)))
	consensusStateExp, err := clienttypes.UnmarshalConsensusState(marshaler, consensusStateByte)
	if err != nil {
		return fmt.Errorf("error unmarshalling consensus state file, err: %w", err)
	}
	consensusState := consensusStateExp.(*ibctm.ConsensusState)

	merklePath := commitmenttypes.NewMerklePath(hosttypes.FullClientStatePath(connection.GetCounterparty().GetClientID()))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	bz, err := marshaler.MarshalInterface(clientState)
	if err != nil {
		return err
	}

	if targetClientState.GetLatestHeight().LT(proofHeight) {
		return fmt.Errorf("client state height < proof height (%d < %d), please ensure the client has been updated", targetClientState.GetLatestHeight(), proofHeight)
	}

	var merkleProof commitmenttypes.MerkleProof
	if err := marshaler.Unmarshal(proofClientbyte, &merkleProof); err != nil {
		return fmt.Errorf("failed to unmarshal proof into ICS 23 commitment merkle proof")
	}
	err = merkleProof.VerifyMembership(targetClientState.ProofSpecs, consensusState.GetRoot(), merklePath, bz)
	if err != nil {
		return err
	}
	return err
}

func consensusStateVerefication(
	connection connectiontypes.ConnectionEnd,
	consensusState exported.ConsensusState,
	height exported.Height,
	accessibleState contract.AccessibleState,
	marshaler *codec.ProtoCodec,
	connectionID string,
	proof []byte,
) error {
	clientID := connection.GetClientID()

	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)
	clientStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientStateExp, err := clienttypes.UnmarshalClientState(marshaler, clientStateByte)
	if err != nil {
		return fmt.Errorf("error unmarshalling client state file, err: %w", err)
	}
	clientState := clientStateExp.(*ibctm.ClientState)

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", clientID, clientState.GetLatestHeight())
	consensusStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)))
	consensusStateExp, err := clienttypes.UnmarshalConsensusState(marshaler, consensusStateByte)
	if err != nil {
		return fmt.Errorf("error unmarshalling consensus state file, err: %w", err)
	}
	targetConsensusState := consensusStateExp.(*ibctm.ConsensusState)

	merklePath := commitmenttypes.NewMerklePath(hosttypes.ConnectionPath(connectionID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	bz, err := marshaler.MarshalInterface(consensusState)
	if err != nil {
		return err
	}

	if clientState.GetLatestHeight().LT(height) {
		return fmt.Errorf("client state height < proof height (%d < %d), please ensure the client has been updated", clientState.GetLatestHeight(), height)
	}

	var merkleProof commitmenttypes.MerkleProof
	if err := marshaler.Unmarshal(proof, &merkleProof); err != nil {
		return fmt.Errorf("failed to unmarshal proof into ICS 23 commitment merkle proof")
	}
	err = merkleProof.VerifyMembership(clientState.ProofSpecs, targetConsensusState.GetRoot(), merklePath, bz)
	if err != nil {
		return err
	}
	return nil
}

func connectionVerefication(
	connection connectiontypes.ConnectionEnd,
	connectionEnd connectiontypes.ConnectionEnd,
	height exported.Height,
	accessibleState contract.AccessibleState,
	marshaler *codec.ProtoCodec,
	connectionID string,
	proof []byte,
) error {
	clientID := connection.GetClientID()

	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)
	clientStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientStateExp, err := clienttypes.UnmarshalClientState(marshaler, clientStateByte)
	if err != nil {
		return fmt.Errorf("error unmarshalling client state file, err: %w", err)
	}
	clientState := clientStateExp.(*ibctm.ClientState)

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", clientID, clientState.GetLatestHeight())
	consensusStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)))
	consensusStateExp, err := clienttypes.UnmarshalConsensusState(marshaler, consensusStateByte)
	if err != nil {
		return fmt.Errorf("error unmarshalling consensus state file, err: %w", err)
	}
	consensusState := consensusStateExp.(*ibctm.ConsensusState)

	merklePath := commitmenttypes.NewMerklePath(hosttypes.ConnectionPath(connectionID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	bz, err := marshaler.Marshal(&connectionEnd)
	if err != nil {
		return err
	}

	if clientState.GetLatestHeight().LT(height) {
		return fmt.Errorf("client state height < proof height (%d < %d), please ensure the client has been updated", clientState.GetLatestHeight(), height)
	}

	var merkleProof commitmenttypes.MerkleProof
	if err := marshaler.Unmarshal(proof, &merkleProof); err != nil {
		return fmt.Errorf("failed to unmarshal proof into ICS 23 commitment merkle proof")
	}
	err = merkleProof.VerifyMembership(clientState.ProofSpecs, consensusState.GetRoot(), merklePath, bz)
	if err != nil {
		return err
	}
	return nil
}

func getClientConnectionPaths(
	marshaler *codec.ProtoCodec,
	clientID string,
	accessibleState contract.AccessibleState,
) ([]string, bool) {

	bz := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress(hosttypes.ClientConnectionsKey(clientID)))
	if len(bz) == 0 {
		return nil, false
	}
	var clientPaths connectiontypes.ClientPaths
	marshaler.MustUnmarshal(bz, &clientPaths)
	return clientPaths.Paths, true
}
