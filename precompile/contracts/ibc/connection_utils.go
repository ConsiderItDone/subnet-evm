package ibc

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	hosttypes "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

func makeConnectionID(db contract.StateDB) string {
	connSeq := uint64(0)
	if db.Exist(common.BytesToAddress([]byte("nextConnSeq"))) {
		b := db.GetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")))
		connSeq = binary.BigEndian.Uint64(b)
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, connSeq+1)
	db.SetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")), b)
	return fmt.Sprintf("%s%d", "connection-", connSeq)
}

func _connOpenInit(opts *callOpts[ConnOpenInitInput]) (string, error) {
	statedb := opts.accessibleState.GetStateDB()

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

	// check ClientState exists in database
	_, err := getClientState(statedb, opts.args.ClientID)
	if err != nil {
		return "", err
	}

	connectionID := makeConnectionID(statedb)

	// connection defines chain A's ConnectionEnd
	connection := connectiontypes.NewConnectionEnd(connectiontypes.INIT, opts.args.ClientID, *counterparty, connectiontypes.ExportedVersionsToProto(versions), uint64(opts.args.DelayPeriod))
	if err := setConnection(statedb, connectionID, &connection); err != nil {
		return "", fmt.Errorf("can't save connection: %w", err)
	}
	if err := addLog(opts.accessibleState, GeneratedConnectionIdentifier.RawName, opts.args.ClientID, connectionID); err != nil {
		return "", fmt.Errorf("error packing event: %w", err)
	}

	return connectionID, nil
}

func _connOpenTry(opts *callOpts[ConnOpenTryInput]) (string, error) {
	statedb := opts.accessibleState.GetStateDB()

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	clienttypes.RegisterInterfaces(interfaceRegistry)
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
	if err := proofHeight.Unmarshal(opts.args.ProofHeight); err != nil {
		return "", fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	consensusHeight := &clienttypes.Height{}
	if err = consensusHeight.Unmarshal(opts.args.ConsensusHeight); err != nil {
		return "", fmt.Errorf("error unmarshalling consensusHeight: %w", err)
	}

	connectionID := makeConnectionID(statedb)
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

	if err = clientVerification(connection, clientState, *proofHeight, opts.accessibleState, marshaler, opts.args.ProofClient); err != nil {
		return "", fmt.Errorf("error clientVerification: %w", err)
	}

	if err = connectionVerification(connection, expectedConnection, *proofHeight, opts.accessibleState, marshaler, connectionID, opts.args.ProofInit); err != nil {
		return "", fmt.Errorf("error connectionVerification: %w", err)
	}

	if err := setConnection(statedb, connectionID, &connection); err != nil {
		return "", fmt.Errorf("can't save connection: %w", err)
	}
	if err := addLog(opts.accessibleState, GeneratedConnectionIdentifier.RawName, opts.args.ClientID, connectionID); err != nil {
		return "", fmt.Errorf("error packing event: %w", err)
	}

	return connectionID, nil
}

func _connOpenAck(opts *callOpts[ConnOpenAckInput]) error {
	statedb := opts.accessibleState.GetStateDB()

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

	connection, err := getConnection(statedb, opts.args.ConnectionID)
	if err != nil {
		return fmt.Errorf("can't get connection: %w", err)
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

	if err := connectionVerification(*connection, expectedConnection, *proofHeight, opts.accessibleState, marshaler, string(opts.args.CounterpartyConnectionID), opts.args.ProofTry); err != nil {
		return fmt.Errorf("connection verification failed: %w", err)
	}

	// Check that ChainB stored the clientState provided in the msg
	if err := clientVerification(*connection, clientState, *proofHeight, opts.accessibleState, marshaler, opts.args.ProofClient); err != nil {
		return fmt.Errorf("client verification failed: %w", err)
	}

	// Update connection state to Open
	connection.State = connectiontypes.OPEN
	connection.Versions = []*connectiontypes.Version{&version}
	connection.Counterparty.ConnectionId = string(opts.args.CounterpartyConnectionID)

	if err := setConnection(statedb, opts.args.ConnectionID, connection); err != nil {
		return fmt.Errorf("can't save connection: %w", err)
	}

	return nil
}

func _connOpenConfirm(opts *callOpts[ConnOpenConfirmInput]) error {
	statedb := opts.accessibleState.GetStateDB()

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	proofHeight := &clienttypes.Height{}
	if err := marshaler.Unmarshal(opts.args.ProofHeight, proofHeight); err != nil {
		return fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	connection, err := getConnection(statedb, opts.args.ConnectionID)
	if err != nil {
		return fmt.Errorf("cannot find connection: %w", err)
	}

	// Check that connection state on ChainB is on state: TRYOPEN
	if connection.State != connectiontypes.TRYOPEN {
		return fmt.Errorf("connection state is not TRYOPEN (got %s), err: %w", connection.State.String(), connectiontypes.ErrInvalidConnectionState)
	}

	// prefix := k.GetCommitmentPrefix()
	expectedCounterparty := connectiontypes.NewCounterparty(connection.ClientId, opts.args.ConnectionID, commitmenttypes.NewMerklePrefix([]byte("ibc")))
	expectedConnection := connectiontypes.NewConnectionEnd(connectiontypes.OPEN, connection.Counterparty.ClientId, expectedCounterparty, connection.Versions, connection.DelayPeriod)

	if err := connectionVerification(*connection, expectedConnection, *proofHeight, opts.accessibleState, marshaler, opts.args.ConnectionID, opts.args.ProofAck); err != nil {
		return err
	}

	clientID := connection.GetClientID()

	clientState, err := getClientState(statedb, clientID)
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}

	consensusState, err := getConsensusState(statedb, clientID, clientState.GetLatestHeight())
	if err != nil {
		return fmt.Errorf("can't get consensus state: %w", err)
	}

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

	if err := setConnection(statedb, opts.args.ConnectionID, connection); err != nil {
		return fmt.Errorf("can't save connection: %w", err)
	}
	return nil
}

func clientVerification(
	connection connectiontypes.ConnectionEnd,
	clientState exported.ClientState,
	proofHeight exported.Height,
	accessibleState contract.AccessibleState,
	marshaler *codec.ProtoCodec,
	proofClientbyte []byte,
) error {
	clientID := connection.GetClientID()

	targetClientState, err := getClientState(accessibleState.GetStateDB(), clientID)
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}

	consensusState, err := getConsensusState(accessibleState.GetStateDB(), clientID, targetClientState.GetLatestHeight())
	if err != nil {
		return fmt.Errorf("error loading consensus state, err: %w", err)
	}

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

func connectionVerification(
	connection connectiontypes.ConnectionEnd,
	connectionEnd connectiontypes.ConnectionEnd,
	height exported.Height,
	accessibleState contract.AccessibleState,
	marshaler *codec.ProtoCodec,
	connectionID string,
	proof []byte,
) error {
	clientID := connection.GetClientID()

	clientState, err := getClientState(accessibleState.GetStateDB(), clientID)
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}

	consensusState, err := getConsensusState(accessibleState.GetStateDB(), clientID, clientState.GetLatestHeight())
	if err != nil {
		return fmt.Errorf("error loading consensus state, err: %w", err)
	}

	merklePath := commitmenttypes.NewMerklePath(hosttypes.ConnectionPath(connectionID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return fmt.Errorf("can't apply prefix %s: %w", connection.GetCounterparty().GetPrefix(), err)
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
