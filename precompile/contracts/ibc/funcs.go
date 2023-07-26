package ibc

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"

	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	hosttypes "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

var (
	keyClientSeq = common.BytesToHash([]byte("client-seq"))
)

type callOpts[T any] struct {
	accessibleState contract.AccessibleState
	caller          common.Address
	addr            common.Address
	suppliedGas     uint64
	readOnly        bool
	args            T
}

func _createClient(opts *callOpts[CreateClientInput]) (string, error) {
	if opts.args.ClientType != exported.Tendermint {
		return "", ErrWrongClientType
	}

	db := opts.accessibleState.GetStateDB()
	clientSeq := db.GetState(ContractAddress, keyClientSeq)
	newClientSeq := common.BigToHash(
		new(big.Int).Add(
			clientSeq.Big(),
			big.NewInt(1),
		),
	)
	db.SetState(ContractAddress, keyClientSeq, newClientSeq)

	return fmt.Sprintf("%s-%d", opts.args.ClientType, clientSeq.Big().Int64()), nil
}

func _updateClient(opts *callOpts[UpdateClientInput]) error {
	stateDB := opts.accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, opts.caller)
	if !callerStatus.IsEnabled() {
		return fmt.Errorf("non-enabled cannot call updateClient: %s", opts.caller)
	}

	clientStatePath := fmt.Sprintf("clients/%s/clientState", opts.args.ClientID)
	found := opts.accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(clientStatePath)))
	if !found {
		return fmt.Errorf("cannot update client with ID %s", opts.args.ClientID)
	}

	clientStateByte := opts.accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientState := &ibctm.ClientState{}
	if err := clientState.Unmarshal(clientStateByte); err != nil {
		return err
	}

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", opts.args.ClientID, clientState.GetLatestHeight())
	found = opts.accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(consensusStatePath)))
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
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)), consensusStateByte)
	return nil
}

func _upgradeClient(opts *callOpts[UpgradeClientInput]) error {
	stateDB := opts.accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, opts.caller)
	if !callerStatus.IsEnabled() {
		return fmt.Errorf("non-enabled cannot call upgradeClient: %s", opts.caller)
	}

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
	clientStateByte := opts.accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientStateExp := clienttypes.MustUnmarshalClientState(marshaler, clientStateByte)
	clientState, ok := clientStateExp.(*ibctm.ClientState)
	if !ok {
		return fmt.Errorf("error unmarshalling client state file")
	}

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", opts.args.ClientID, clientState.GetLatestHeight())
	consensusStateByte := opts.accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)))
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
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)), consensusStateByte)

	clientStateByte, err = marshaler.MarshalInterface(newClientState)
	if err != nil {
		return errors.New("clientState marshaler error")
	}
	clientStatePath = fmt.Sprintf("clients/%s/clientState", opts.args.ClientID)
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(clientStatePath)), clientStateByte)
	return nil
}

func _connOpenInit(opts *callOpts[ConnOpenInitInput]) (string, error) {
	stateDB := opts.accessibleState.GetStateDB()

	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, opts.caller)
	if !callerStatus.IsEnabled() {
		return "", fmt.Errorf("non-enabled cannot call upgradeClient: %s", opts.caller)
	}

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	counterparty := &connectiontypes.Counterparty{}
	if err := counterparty.Unmarshal(opts.args.Counterparty); err != nil {
		return "", fmt.Errorf("error unmarshalling counterparty: %w", err)
	}

	version := &connectiontypes.Version{}
	if err := version.Unmarshal(opts.args.Version); err != nil {
		return "", fmt.Errorf("error unmarshalling version: %w", err)
	}

	versions := connectiontypes.GetCompatibleVersions()
	if len(opts.args.Version) != 0 {
		if !connectiontypes.IsSupportedVersion(connectiontypes.GetCompatibleVersions(), version) {
			return "", fmt.Errorf("%w : version is not supported", connectiontypes.ErrInvalidVersion)
		}
		versions = []exported.Version{version}
	}

	clientStatePath := fmt.Sprintf("clients/%s/clientState", opts.args.ClientID)
	found := opts.accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(clientStatePath)))
	if !found {
		return "", fmt.Errorf("cannot update client with ID %s", opts.args.ClientID)
	}

	nextConnSeq := uint64(0)
	if opts.accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte("nextConnSeq"))) {
		b := opts.accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")))
		nextConnSeq = binary.BigEndian.Uint64(b)
	}
	connectionID := fmt.Sprintf("%s%d", "connection-", nextConnSeq)
	nextConnSeq++
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, nextConnSeq)
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")), b)

	// connection defines chain A's ConnectionEnd
	connection := connectiontypes.NewConnectionEnd(connectiontypes.INIT, opts.args.ClientID, *counterparty, connectiontypes.ExportedVersionsToProto(versions), uint64(opts.args.DelayPeriod))

	connectionByte, err := marshaler.Marshal(&connection)
	if err != nil {
		return "", fmt.Errorf("connection marshaler error: %w", err)
	}
	connectionsPath := fmt.Sprintf("connections/%s", connectionID)
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

	return connectionID, nil
}

func _connOpenTry(opts *callOpts[ConnOpenTryInput]) (string, error) {
	stateDB := opts.accessibleState.GetStateDB()

	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, opts.caller)
	if !callerStatus.IsEnabled() {
		return "", fmt.Errorf("non-enabled cannot call upgradeClient: %s", opts.caller)
	}

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	counterparty := &connectiontypes.Counterparty{}
	if err := counterparty.Unmarshal(opts.args.Counterparty); err != nil {
		return "", fmt.Errorf("error unmarshalling counterparty: %w", err)
	}

	clientStateExp, err := clienttypes.UnmarshalClientState(marshaler, opts.args.ClientState)
	clientState := clientStateExp.(*ibctm.ClientState)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling clientState: %w", err)
	}

	counterpartyVersions := []*connectiontypes.Version{}
	if err := json.Unmarshal(opts.args.CounterpartyVersions, &counterpartyVersions); err != nil {
		return "", fmt.Errorf("error unmarshalling counterpartyVersions: %w", err)
	}

	proofHeight := &clienttypes.Height{}
	if err := marshaler.UnmarshalInterface(opts.args.ProofHeight, proofHeight); err != nil {
		return "", fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	consensusHeight := &clienttypes.Height{}
	if err = marshaler.UnmarshalInterface(opts.args.ConsensusHeight, consensusHeight); err != nil {
		return "", fmt.Errorf("error unmarshalling consensusHeight: %w", err)
	}

	nextConnSeq := uint64(0)
	if opts.accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte("nextConnSeq"))) {
		b := opts.accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")))
		nextConnSeq = binary.BigEndian.Uint64(b)
	}
	connectionID := fmt.Sprintf("%s%d", "connection-", nextConnSeq)
	nextConnSeq++
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, nextConnSeq)
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")), b)

	expectedCounterparty := connectiontypes.NewCounterparty(opts.args.ClientID, "", commitmenttypes.NewMerklePrefix([]byte("ibc")))
	expectedConnection := connectiontypes.NewConnectionEnd(connectiontypes.INIT, counterparty.ClientId, expectedCounterparty, counterpartyVersions, uint64(opts.args.DelayPeriod))

	// chain B picks a version from Chain A's available versions that is compatible
	// with Chain B's supported IBC versions. PickVersion will select the intersection
	// of the supported versions and the counterparty versions.
	version, err := connectiontypes.PickVersion(connectiontypes.GetCompatibleVersions(), connectiontypes.ProtoVersionsToExported(counterpartyVersions))
	if err != nil {
		return "", fmt.Errorf("error PickVersion err: %w", err)
	}

	// connection defines chain B's ConnectionEnd
	connection := connectiontypes.NewConnectionEnd(connectiontypes.TRYOPEN, opts.args.ClientID, *counterparty, []*connectiontypes.Version{version}, uint64(opts.args.DelayPeriod))

	if err = clientVerefication(connection, clientState, *proofHeight, opts.accessibleState, marshaler, opts.args.ProofClient); err != nil {
		return "", fmt.Errorf("error clientVerefication err: %w", err)
	}

	if err = connectionVerefication(connection, expectedConnection, *proofHeight, opts.accessibleState, marshaler, connectionID, opts.args.ProofInit); err != nil {
		return "", fmt.Errorf("error connectionVerefication err: %w", err)
	}

	clientStatePath := fmt.Sprintf("clients/%s/clientState", opts.args.ClientID)
	found := opts.accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(clientStatePath)))
	if !found {
		return "", clienttypes.ErrClientNotFound
	}

	conns, found := getClientConnectionPaths(marshaler, opts.args.ClientID, opts.accessibleState)
	if !found {
		conns = []string{}
	}
	conns = append(conns, connectionID)
	clientPaths := connectiontypes.ClientPaths{Paths: conns}
	bz := marshaler.MustMarshal(&clientPaths)

	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress(hosttypes.ClientConnectionsKey(opts.args.ClientID)), bz)
	return connectionID, nil
}

func _connOpenAck(opts *callOpts[ConnOpenAckInput]) error {
	stateDB := opts.accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, opts.caller)
	if !callerStatus.IsEnabled() {
		return fmt.Errorf("non-enabled cannot call upgradeClient: %s", opts.caller)
	}

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	clientStateExp, err := clienttypes.UnmarshalClientState(marshaler, opts.args.ClientState)
	clientState := clientStateExp.(*ibctm.ClientState)
	if err != nil {
		return fmt.Errorf("error unmarshalling clientState: %w", err)
	}

	version := connectiontypes.Version{}
	if err = marshaler.Unmarshal(opts.args.Version, &version); err != nil {
		return fmt.Errorf("error unmarshalling counterpartyVersions: %w", err)
	}

	proofHeight := &clienttypes.Height{}
	if err = marshaler.Unmarshal(opts.args.ProofHeight, proofHeight); err != nil {
		return fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	consensusHeight := &clienttypes.Height{}
	if err = marshaler.Unmarshal(opts.args.ConsensusHeight, consensusHeight); err != nil {
		return fmt.Errorf("error unmarshalling consensusHeight: %w", err)
	}

	connectionsPath := fmt.Sprintf("connections/%s", opts.args.ConnectionID)
	connectionByte := opts.accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(connectionsPath)))
	connection := connectiontypes.ConnectionEnd{}
	if err = marshaler.Unmarshal(connectionByte, &connection); err != nil {
		return fmt.Errorf("error unmarshalling connection id: %s, error: %w", opts.args.ConnectionID, err)
	}

	// verify the previously set connection state
	if connection.State != connectiontypes.INIT {
		return fmt.Errorf("connection state is not INIT (got %s), error: %w", connection.State.String(), connectiontypes.ErrInvalidConnectionState)
	}

	// ensure selected version is supported
	if !connectiontypes.IsSupportedVersion(connectiontypes.ProtoVersionsToExported(connection.Versions), &version) {
		return fmt.Errorf("the counterparty selected version %s is not supported by versions selected on INIT, error: %w", version, connectiontypes.ErrInvalidConnectionState)
	}

	expectedCounterparty := connectiontypes.NewCounterparty(connection.ClientId, opts.args.ConnectionID, commitmenttypes.NewMerklePrefix([]byte("ibc")))
	expectedConnection := connectiontypes.NewConnectionEnd(connectiontypes.TRYOPEN, connection.Counterparty.ClientId, expectedCounterparty, []*connectiontypes.Version{&version}, connection.DelayPeriod)

	if err := connectionVerefication(connection, expectedConnection, *proofHeight, opts.accessibleState, marshaler, string(opts.args.CounterpartyConnectionID), opts.args.ProofTry); err != nil {
		return err
	}

	// Check that ChainB stored the clientState provided in the msg
	if err := clientVerefication(connection, clientState, *proofHeight, opts.accessibleState, marshaler, opts.args.ProofClient); err != nil {
		return err
	}

	// Update connection state to Open
	connection.State = connectiontypes.OPEN
	connection.Versions = []*connectiontypes.Version{&version}
	connection.Counterparty.ConnectionId = string(opts.args.CounterpartyConnectionID)

	connectionByte, err = marshaler.Marshal(&connection)
	if err != nil {
		return errors.New("connection marshaler error")
	}
	connectionsPath = fmt.Sprintf("connections/%s", opts.args.ConnectionID)
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

	return nil
}

func _connOpenConfirm(opts *callOpts[ConnOpenConfirmInput]) error {
	stateDB := opts.accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, opts.caller)
	if !callerStatus.IsEnabled() {
		return fmt.Errorf("non-enabled cannot call upgradeClient: %s", opts.caller)
	}

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	proofHeight := &clienttypes.Height{}
	if err := marshaler.Unmarshal(opts.args.ProofHeight, proofHeight); err != nil {
		return fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	connectionsPath := fmt.Sprintf("connections/%s", opts.args.ConnectionID)
	exist := opts.accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(connectionsPath)))
	if !exist {
		return fmt.Errorf("cannot find connection with path: %s", connectionsPath)
	}

	connectionByte := opts.accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(connectionsPath)))
	connection := &connectiontypes.ConnectionEnd{}
	marshaler.MustUnmarshal(connectionByte, connection)

	// Check that connection state on ChainB is on state: TRYOPEN
	if connection.State != connectiontypes.TRYOPEN {
		return fmt.Errorf("connection state is not TRYOPEN (got %s), err: %w", connection.State.String(), connectiontypes.ErrInvalidConnectionState)
	}

	// prefix := k.GetCommitmentPrefix()
	expectedCounterparty := connectiontypes.NewCounterparty(connection.ClientId, opts.args.ConnectionID, commitmenttypes.NewMerklePrefix([]byte("ibc")))
	expectedConnection := connectiontypes.NewConnectionEnd(connectiontypes.OPEN, connection.Counterparty.ClientId, expectedCounterparty, connection.Versions, connection.DelayPeriod)

	if err := connectionVerefication(*connection, expectedConnection, *proofHeight, opts.accessibleState, marshaler, opts.args.ConnectionID, opts.args.ProofAck); err != nil {
		return err
	}

	clientID := connection.GetClientID()

	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)
	clientStateByte := opts.accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientStateExp, err := clienttypes.UnmarshalClientState(marshaler, clientStateByte)
	if err != nil {
		return fmt.Errorf("error unmarshalling client state file, err: %w", err)
	}
	clientState := clientStateExp.(*ibctm.ClientState)

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", clientID, clientState.GetLatestHeight())
	consensusStateByte := opts.accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)))
	consensusStateExp, err := clienttypes.UnmarshalConsensusState(marshaler, consensusStateByte)
	if err != nil {
		return fmt.Errorf("error unmarshalling consensus state file, err: %w", err)
	}
	consensusState := consensusStateExp.(*ibctm.ConsensusState)

	merklePath := commitmenttypes.NewMerklePath(hosttypes.ConnectionPath(opts.args.ConnectionID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	bz, err := marshaler.Marshal(&expectedConnection)
	if err != nil {
		return err
	}

	if clientState.GetLatestHeight().LT(*proofHeight) {
		return fmt.Errorf("client state height < proof height (%d < %d), please ensure the client has been updated", clientState.GetLatestHeight(), proofHeight)
	}

	var merkleProof commitmenttypes.MerkleProof
	if err := marshaler.Unmarshal(opts.args.ProofAck, &merkleProof); err != nil {
		return fmt.Errorf("failed to unmarshal proof into ICS 23 commitment merkle proof")
	}
	merkleProof.VerifyMembership(clientState.ProofSpecs, consensusState.GetRoot(), merklePath, bz)

	// Update ChainB's connection to Open
	connection.State = connectiontypes.OPEN

	connectionByte, err = marshaler.Marshal(connection)
	if err != nil {
		return errors.New("connection marshaler error")
	}
	connectionsPath = fmt.Sprintf("connections/%s", opts.args.ConnectionID)
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

	return nil
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
