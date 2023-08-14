package ibc

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
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

func _connOpenInit(opts *callOpts[ConnOpenInitInput]) (string, error) {
	stateDB := opts.accessibleState.GetStateDB()

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

	// check ClientState exists in database
	_, found, err := getClientState(stateDB, opts.args.ClientID)
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("cannot update client with ID %s", opts.args.ClientID)
	}

	nextConnSeq := uint64(0)
	if stateDB.Exist(common.BytesToAddress([]byte("nextConnSeq"))) {
		b := stateDB.GetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")))
		nextConnSeq = binary.BigEndian.Uint64(b)
	}
	connectionID := fmt.Sprintf("%s%d", "connection-", nextConnSeq)
	nextConnSeq++
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, nextConnSeq)
	stateDB.SetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")), b)

	// connection defines chain A's ConnectionEnd
	connection := connectiontypes.NewConnectionEnd(connectiontypes.INIT, opts.args.ClientID, *counterparty, connectiontypes.ExportedVersionsToProto(versions), uint64(opts.args.DelayPeriod))

	connectionByte, err := marshaler.Marshal(&connection)
	if err != nil {
		return "", fmt.Errorf("connection marshaler error: %w", err)
	}
	connectionsPath := fmt.Sprintf("connections/%s", connectionID)
	stateDB.SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

	// emit event
	topics, data, err := IBCABI.PackEvent(GeneratedConnectionIdentifier.RawName, opts.args.ClientID, connectionID)
	if err != nil {
		return "", fmt.Errorf("error packing event: %w", err)
	}
	blockNumber := opts.accessibleState.GetBlockContext().Number().Uint64()
	opts.accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

	return connectionID, nil
}

func _connOpenTry(opts *callOpts[ConnOpenTryInput]) (string, error) {
	stateDB := opts.accessibleState.GetStateDB()

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
	if stateDB.Exist(common.BytesToAddress([]byte("nextConnSeq"))) {
		b := stateDB.GetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")))
		nextConnSeq = binary.BigEndian.Uint64(b)
	}
	connectionID := fmt.Sprintf("%s%d", "connection-", nextConnSeq)
	nextConnSeq++
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, nextConnSeq)
	stateDB.SetPrecompileState(common.BytesToAddress([]byte("nextConnSeq")), b)

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

	return connectionID, nil
}

func _connOpenAck(opts *callOpts[ConnOpenAckInput]) error {
	stateDB := opts.accessibleState.GetStateDB()

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
	connectionByte := stateDB.GetPrecompileState(common.BytesToAddress([]byte(connectionsPath)))
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

	if err := connectionVerification(connection, expectedConnection, *proofHeight, opts.accessibleState, marshaler, string(opts.args.CounterpartyConnectionID), opts.args.ProofTry); err != nil {
		return err
	}

	// Check that ChainB stored the clientState provided in the msg
	if err := clientVerification(connection, clientState, *proofHeight, opts.accessibleState, marshaler, opts.args.ProofClient); err != nil {
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
	stateDB.SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

	return nil
}

func _connOpenConfirm(opts *callOpts[ConnOpenConfirmInput]) error {
	stateDB := opts.accessibleState.GetStateDB()

	interfaceRegistry := cosmostypes.NewInterfaceRegistry()

	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	proofHeight := &clienttypes.Height{}
	if err := marshaler.Unmarshal(opts.args.ProofHeight, proofHeight); err != nil {
		return fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	connectionsPath := fmt.Sprintf("connections/%s", opts.args.ConnectionID)
	exist := stateDB.Exist(common.BytesToAddress([]byte(connectionsPath)))
	if !exist {
		return fmt.Errorf("cannot find connection with path: %s", connectionsPath)
	}

	connectionByte := stateDB.GetPrecompileState(common.BytesToAddress([]byte(connectionsPath)))
	connection := &connectiontypes.ConnectionEnd{}
	marshaler.MustUnmarshal(connectionByte, connection)

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

	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)
	clientStateByte := stateDB.GetPrecompileState(common.BytesToAddress([]byte(clientStatePath)))
	clientStateExp, err := clienttypes.UnmarshalClientState(marshaler, clientStateByte)
	if err != nil {
		return fmt.Errorf("error unmarshalling client state file, err: %w", err)
	}
	clientState := clientStateExp.(*ibctm.ClientState)

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", clientID, clientState.GetLatestHeight())
	consensusStateByte := stateDB.GetPrecompileState(common.BytesToAddress([]byte(consensusStatePath)))
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
	stateDB.SetPrecompileState(common.BytesToAddress([]byte(connectionsPath)), connectionByte)

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

	targetClientState, found, err := getClientState(accessibleState.GetStateDB(), clientID)
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}
	if !found {
		return fmt.Errorf("client state not found in database")
	}

	consensusState, found, err := getConsensusState(accessibleState.GetStateDB(), clientID, targetClientState.GetLatestHeight())
	if err != nil {
		return fmt.Errorf("error loading consensus state, err: %w", err)
	}
	if !found {
		return fmt.Errorf("consensus state not found in database")
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

	clientState, found, err := getClientState(accessibleState.GetStateDB(), clientID)
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}
	if !found {
		return fmt.Errorf("client state not found in database")
	}

	consensusState, found, err := getConsensusState(accessibleState.GetStateDB(), clientID, clientState.GetLatestHeight())
	if err != nil {
		return fmt.Errorf("error loading consensus state, err: %w", err)
	}
	if !found {
		return fmt.Errorf("consensus state not found in database")
	}

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

func getConnection(
	marshaler *codec.ProtoCodec,
	connectionsPath string,
	accessibleState contract.AccessibleState,
) (*connectiontypes.ConnectionEnd, error) {
	// connection hop length checked on msg.ValidateBasic()
	exist := accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(connectionsPath)))
	if !exist {
		return nil, fmt.Errorf("cannot find connection with path: %s", connectionsPath)
	}
	connection := &connectiontypes.ConnectionEnd{}
	connectionByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(connectionsPath)))

	marshaler.MustUnmarshal(connectionByte, connection)
	return connection, nil
}

func generateChannelIdentifier(accessibleState contract.AccessibleState) string {
	sequence := uint64(0)
	if accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte("nextChannelSeq"))) {
		b := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte("nextChannelSeq")))
		sequence = binary.BigEndian.Uint64(b)
	}
	sequence++
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, sequence)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte("nextChannelSeq")), b)

	return fmt.Sprintf("%s%d", "channel-", sequence)
}

// setNextSequenceSend sets a channel's next send sequence to the store
func setNextSequenceSend(accessibleState contract.AccessibleState, portID, channelID string, sequence uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, sequence)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(hosttypes.NextSequenceSendKey(portID, channelID))), b)
}

// setNextSequenceSend sets a channel's next send sequence to the store
func getNextSequenceSend(accessibleState contract.AccessibleState, portID, channelID string) uint64 {
	b := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(hosttypes.NextSequenceSendKey(portID, channelID))))
	return binary.BigEndian.Uint64(b)
}

// setNextSequenceRecv sets a channel's next receive sequence to the store
func setNextSequenceRecv(accessibleState contract.AccessibleState, portID, channelID string, sequence uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, sequence)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(hosttypes.NextSequenceRecvKey(portID, channelID))), b)
}

// setNextSequenceAck sets a channel's next ack sequence to the store
func setNextSequenceAck(accessibleState contract.AccessibleState, portID, channelID string, sequence uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, sequence)

	path := []byte(fmt.Sprintf("%s/%s/%s", portID, channelID, "nextSequenceAck"))
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress(path), b)
	// TODO
	// accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(hosttypes.NextSequenceAckKey(portID, channelID))), b)
}

func setPacketCommitment(accessibleState contract.AccessibleState, portID, channelID string, sequence uint64, commitmentHash []byte) {
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(hosttypes.PacketCommitmentKey(portID, channelID, sequence))), commitmentHash)
}

func _chanOpenInit(opts *callOpts[ChanOpenInitInput]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	channel := &channeltypes.Channel{}
	if err := marshaler.Unmarshal(opts.args.Channel, channel); err != nil {
		return err
	}

	// connection hop length checked on msg.ValidateBasic()
	connectionsPath := fmt.Sprintf("connections/%s", channel.ConnectionHops[0])
	connectionEnd, err := getConnection(marshaler, connectionsPath, opts.accessibleState)
	if err != nil {
		return err
	}

	getVersions := connectionEnd.GetVersions()
	if len(getVersions) != 1 {
		return fmt.Errorf("single version must be negotiated on connection before opening channel, got: %v, err: %w",
			getVersions,
			connectiontypes.ErrInvalidVersion,
		)
	}

	if !connectiontypes.VerifySupportedFeature(getVersions[0], channel.Ordering.String()) {
		return fmt.Errorf("connection version %s does not support channel ordering: %s, err: %w",
			getVersions[0], channel.Ordering.String(),
			connectiontypes.ErrInvalidVersion,
		)
	}

	_, found, err := getClientState(opts.accessibleState.GetStateDB(), connectionEnd.ClientId)
	if err != nil {
		return fmt.Errorf("can't read client state: %w", err)
	}
	if !found {
		return fmt.Errorf("client state not found: %s", connectionEnd.ClientId)
	}

	channelID := generateChannelIdentifier(opts.accessibleState)

	err = makeCapability(opts.accessibleState.GetStateDB(), opts.args.PortID, channelID)
	if err != nil {
		return err
	}

	channelNew := channeltypes.NewChannel(channeltypes.INIT, channel.Ordering, channel.Counterparty, channel.ConnectionHops, channel.Version)
	bz := marshaler.MustMarshal(&channelNew)
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(hosttypes.ChannelKey(opts.args.PortID, channelID))), bz)

	setNextSequenceSend(opts.accessibleState, opts.args.PortID, channelID, 1)
	setNextSequenceRecv(opts.accessibleState, opts.args.PortID, channelID, 1)
	setNextSequenceAck(opts.accessibleState, opts.args.PortID, channelID, 1)

	return nil
}

func channelStateVerification(
	connection connectiontypes.ConnectionEnd,
	channel channeltypes.Channel,
	height exported.Height,
	accessibleState contract.AccessibleState,
	marshaler *codec.ProtoCodec,
	channelID string,
	proof []byte,
	portID string,
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

	merklePath := commitmenttypes.NewMerklePath(hosttypes.ChannelPath(portID, channelID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	bz, err := marshaler.Marshal(&channel)
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

func _chanOpenTry(opts *callOpts[ChanOpenTryInput]) (string, error) {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	channel := &channeltypes.Channel{}
	if err := marshaler.Unmarshal(opts.args.Channel, channel); err != nil {
		return "", fmt.Errorf("error unmarshalling channel: %w", err)
	}

	proofHeight := &clienttypes.Height{}
	if err := marshaler.Unmarshal(opts.args.ProofHeight, proofHeight); err != nil {
		return "", fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	if len(channel.ConnectionHops) != 1 {
		return "", fmt.Errorf("expected 1, got %d, err: %w", len(channel.ConnectionHops), channeltypes.ErrTooManyConnectionHops)
	}

	// generate a new channel
	channelID := generateChannelIdentifier(opts.accessibleState)

	// connection hop length checked on msg.ValidateBasic()
	connectionsPath := fmt.Sprintf("connections/%s", channel.ConnectionHops[0])
	connectionEnd, err := getConnection(marshaler, connectionsPath, opts.accessibleState)
	if err != nil {
		return "", err
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return "", fmt.Errorf("connection state is not OPEN (got %s), err: %w", connectiontypes.State(connectionEnd.GetState()).String(), connectiontypes.ErrInvalidConnectionState)
	}

	getVersions := connectionEnd.GetVersions()
	if len(getVersions) != 1 {
		return "", fmt.Errorf("single version must be negotiated on connection before opening channel, got: %v, err: %w", getVersions, connectiontypes.ErrInvalidVersion)
	}

	if !connectiontypes.VerifySupportedFeature(getVersions[0], channel.Ordering.String()) {
		return "", fmt.Errorf("connection version %s does not support channel ordering: %s, err: %w", getVersions[0], channel.Ordering.String(), connectiontypes.ErrInvalidVersion)
	}

	ok, err := getCapability(opts.accessibleState.GetStateDB(), opts.args.PortID, channelID)
	if !ok {
		return "", fmt.Errorf("caller does not own port capability for port ID %s, %w", opts.args.PortID, err)
	}

	counterpartyHops := []string{connectionEnd.GetCounterparty().GetConnectionID()}

	// expectedCounterpaty is the counterparty of the counterparty's channel end
	// (i.e self)
	expectedCounterparty := channeltypes.NewCounterparty(opts.args.PortID, "")
	expectedChannel := channeltypes.NewChannel(
		channeltypes.INIT, channel.Ordering, expectedCounterparty,
		counterpartyHops, opts.args.CounterpartyVersion,
	)

	if err := channelStateVerification(
		*connectionEnd,
		expectedChannel,
		*proofHeight,
		opts.accessibleState,
		marshaler,
		channel.Counterparty.ChannelId,
		opts.args.ProofInit,
		opts.args.PortID,
	); err != nil {
		return "", err
	}

	setNextSequenceSend(opts.accessibleState, opts.args.PortID, channelID, 1)
	setNextSequenceRecv(opts.accessibleState, opts.args.PortID, channelID, 1)
	setNextSequenceAck(opts.accessibleState, opts.args.PortID, channelID, 1)

	channelNew := channeltypes.NewChannel(channeltypes.TRYOPEN, channel.Ordering, channel.Counterparty, channel.ConnectionHops, channel.Version)
	bz := marshaler.MustMarshal(&channelNew)
	opts.accessibleState.GetStateDB().SetPrecompileState(
		common.BytesToAddress([]byte(hosttypes.ChannelKey(
			opts.args.PortID,
			channelID,
		))),
		bz,
	)

	return channelID, nil
}

func getChannelState(
	marshaler *codec.ProtoCodec,
	channelStatePath string,
	accessibleState contract.AccessibleState,
) (*channeltypes.Channel, error) {
	// connection hop length checked on msg.ValidateBasic()
	exist := accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte(channelStatePath)))
	if !exist {
		return nil, fmt.Errorf("cannot find channel state with path: %s", channelStatePath)
	}
	channelStateByte := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte(channelStatePath)))
	channelState := &channeltypes.Channel{}
	marshaler.MustUnmarshal(channelStateByte, channelState)
	return channelState, nil
}

func _channelOpenAck(opts *callOpts[ChannelOpenAckInput]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	proofHeight := &clienttypes.Height{}
	if err := marshaler.Unmarshal(opts.args.ProofHeight, proofHeight); err != nil {
		return fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	channel, err := getChannelState(
		marshaler,
		string(hosttypes.ChannelKey(
			opts.args.PortID,
			opts.args.ChannelID,
		)),
		opts.accessibleState,
	)
	if err != nil {
		return err
	}

	if channel.State != channeltypes.INIT {
		return fmt.Errorf("channel state should be INIT (got %s), err: %w", channel.State.String(), channeltypes.ErrInvalidChannelState)
	}

	connectionsPath := fmt.Sprintf("connections/%s", channel.ConnectionHops[0])
	connectionEnd, err := getConnection(marshaler, connectionsPath, opts.accessibleState)
	if err != nil {
		return err
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return fmt.Errorf("connection state is not OPEN (got %s), err: %w", connectiontypes.State(connectionEnd.GetState()).String(), connectiontypes.ErrInvalidConnectionState)
	}

	counterpartyHops := []string{connectionEnd.GetCounterparty().GetConnectionID()}

	// counterparty of the counterparty channel end (i.e self)
	expectedCounterparty := channeltypes.NewCounterparty(
		opts.args.PortID,
		opts.args.ChannelID,
	)

	ok, err := getCapability(opts.accessibleState.GetStateDB(), opts.args.PortID, opts.args.ChannelID)
	if !ok {
		return fmt.Errorf("caller does not own port capability for port ID %s, %w", opts.args.PortID, err)
	}

	expectedChannel := channeltypes.NewChannel(
		channeltypes.TRYOPEN, channel.Ordering, expectedCounterparty,
		counterpartyHops, opts.args.CounterpartyVersion,
	)

	err = channelStateVerification(*connectionEnd, expectedChannel, *proofHeight, opts.accessibleState, marshaler, opts.args.ChannelID, opts.args.ProofTry, channel.Counterparty.PortId)
	if err != nil {
		return fmt.Errorf("channel handshake open ack failed")
	}

	channel.State = channeltypes.OPEN
	channel.Version = opts.args.CounterpartyVersion
	channel.Counterparty.ChannelId = opts.args.CounterpartyChannelID

	bz := marshaler.MustMarshal(channel)
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(hosttypes.ChannelKey(
		opts.args.PortID,
		opts.args.ChannelID,
	))), bz)

	return nil
}

func _channelOpenConfirm(opts *callOpts[ChannelOpenConfirmInput]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	proofHeight := &clienttypes.Height{}
	if err := marshaler.Unmarshal(opts.args.ProofHeight, proofHeight); err != nil {
		return fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	channel, err := getChannelState(
		marshaler,
		string(hosttypes.ChannelKey(
			opts.args.PortID,
			opts.args.ChannelID,
		)),
		opts.accessibleState,
	)
	if err != nil {
		return err
	}

	ok, err := getCapability(opts.accessibleState.GetStateDB(), opts.args.PortID, opts.args.ChannelID)
	if !ok {
		return fmt.Errorf("caller does not own port capability for port ID %s, %w", opts.args.PortID, err)
	}

	if channel.State != channeltypes.TRYOPEN {
		return fmt.Errorf("channel state is not TRYOPEN (got %s), err: %w", channel.State.String(), channeltypes.ErrInvalidChannelState)
	}

	connectionsPath := fmt.Sprintf("connections/%s", channel.ConnectionHops[0])
	connectionEnd, err := getConnection(marshaler, connectionsPath, opts.accessibleState)
	if err != nil {
		return err
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return fmt.Errorf("connection state is not OPEN (got %s), err: %w", connectiontypes.State(connectionEnd.GetState()).String(), connectiontypes.ErrInvalidConnectionState)
	}

	counterpartyHops := []string{connectionEnd.GetCounterparty().GetConnectionID()}

	counterparty := channeltypes.NewCounterparty(
		opts.args.PortID,
		opts.args.ChannelID,
	)
	expectedChannel := channeltypes.NewChannel(
		channeltypes.OPEN, channel.Ordering, counterparty,
		counterpartyHops, channel.Version,
	)

	err = channelStateVerification(*connectionEnd, expectedChannel, *proofHeight, opts.accessibleState, marshaler, channel.Counterparty.ChannelId, opts.args.ProofAck, channel.Counterparty.PortId)
	if err != nil {
		return fmt.Errorf("channel handshake open ack failed")
	}
	channel.State = channeltypes.OPEN

	bz := marshaler.MustMarshal(channel)
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(hosttypes.ChannelKey(opts.args.PortID, opts.args.ChannelID))), bz)

	return nil
}

func _channelCloseInit(opts *callOpts[ChannelCloseInitInput]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	channel, err := getChannelState(
		marshaler,
		string(hosttypes.ChannelKey(
			opts.args.PortID,
			opts.args.ChannelID,
		)),
		opts.accessibleState,
	)
	if err != nil {
		return err
	}

	if channel.State == channeltypes.CLOSED {
		return fmt.Errorf("channel is already CLOSED: %w", channeltypes.ErrInvalidChannelState)
	}

	if len(channel.ConnectionHops) == 0 {
		return fmt.Errorf("length channel.ConnectionHops == 0")
	}
	connectionsPath := fmt.Sprintf("connections/%s", channel.ConnectionHops[0])
	connectionEnd, err := getConnection(marshaler, connectionsPath, opts.accessibleState)
	if err != nil {
		return err
	}

	_, found, err := getClientState(opts.accessibleState.GetStateDB(), connectionEnd.ClientId)
	if err != nil {
		return fmt.Errorf("can't read client state: %w", err)
	}
	if !found {
		return fmt.Errorf("client state not found: %s", connectionEnd.ClientId)
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return fmt.Errorf("connection state is not OPEN (got %s), err: %w", connectiontypes.State(connectionEnd.GetState()).String(), connectiontypes.ErrInvalidConnectionState)
	}

	channel.State = channeltypes.CLOSED

	bz := marshaler.MustMarshal(channel)
	opts.accessibleState.GetStateDB().SetPrecompileState(
		common.BytesToAddress([]byte(hosttypes.ChannelKey(
			opts.args.PortID,
			opts.args.ChannelID,
		))),
		bz,
	)

	return nil
}

func _channelCloseConfirm(opts *callOpts[ChannelCloseConfirmInput]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	proofHeight := &clienttypes.Height{}
	if err := marshaler.Unmarshal(opts.args.ProofHeight, proofHeight); err != nil {
		return fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	channel, err := getChannelState(
		marshaler,
		string(hosttypes.ChannelKey(
			opts.args.PortID,
			opts.args.ChannelID,
		)),
		opts.accessibleState,
	)
	if err != nil {
		return err
	}

	if channel.State == channeltypes.CLOSED {
		return fmt.Errorf("channel is already CLOSED: %w", channeltypes.ErrInvalidChannelState)
	}

	connectionsPath := fmt.Sprintf("connections/%s", channel.ConnectionHops[0])
	connectionEnd, err := getConnection(marshaler, connectionsPath, opts.accessibleState)
	if err != nil {
		return err
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return fmt.Errorf("connection state is not OPEN (got %s), err: %w", connectiontypes.State(connectionEnd.GetState()).String(), channeltypes.ErrInvalidChannelState)
	}

	counterpartyHops := []string{connectionEnd.GetCounterparty().GetConnectionID()}
	counterparty := channeltypes.NewCounterparty(
		opts.args.PortID,
		opts.args.ChannelID,
	)
	expectedChannel := channeltypes.NewChannel(
		channeltypes.CLOSED, channel.Ordering, counterparty,
		counterpartyHops, channel.Version,
	)

	err = channelStateVerification(*connectionEnd, expectedChannel, *proofHeight, opts.accessibleState, marshaler, channel.Counterparty.ChannelId, opts.args.ProofInit, channel.Counterparty.PortId)
	if err != nil {
		return err
	}

	channel.State = channeltypes.CLOSED

	bz := marshaler.MustMarshal(channel)
	opts.accessibleState.GetStateDB().SetPrecompileState(
		common.BytesToAddress([]byte(hosttypes.ChannelKey(
			opts.args.PortID,
			opts.args.ChannelID,
		))),
		bz,
	)

	return nil
}

func _sendPacket(opts *callOpts[SendPacketInput]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	sequence := getNextSequenceSend(opts.accessibleState, opts.args.SourcePort, opts.args.SourceChannel)

	channel, err := getChannelState(
		marshaler,
		string(hosttypes.ChannelKey(
			opts.args.SourcePort,
			opts.args.SourceChannel,
		)),
		opts.accessibleState,
	)
	if err != nil {
		return err
	}

	if channel.State != channeltypes.OPEN {
		return fmt.Errorf("%w, channel is not OPEN (got %s)", channeltypes.ErrInvalidChannelState, channel.State.String())
	}

	getCapability(opts.accessibleState.GetStateDB(), opts.args.SourcePort, opts.args.SourceChannel)
	if !k.scopedKeeper.AuthenticateCapability(ctx, channelCap, host.ChannelCapabilityPath(opts.args.SourcePort, opts.args.SourceChannel)) {
		return fmt.Errorf("%w, caller does not own capability for channel, port ID (%s) channel ID (%s)", opts.args.SourcePort, opts.args.SourceChannel, channeltypes.ErrChannelCapabilityNotFound)
	}

	height := clienttypes.Height{
		RevisionNumber: opts.args.TimeoutHeight.RevisionNumber,
		RevisionHeight: opts.args.TimeoutHeight.RevisionHeight,
	}

	packet := channeltypes.NewPacket(opts.args.Data, sequence, opts.args.SourcePort, opts.args.SourceChannel,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId, height, opts.args.TimeoutTimestamp)

	connectionsPath := fmt.Sprintf("connections/%s", channel.ConnectionHops[0])
	connectionEnd, err := getConnection(marshaler, connectionsPath, opts.accessibleState)
	if err != nil {
		return err
	}

	clientState, found, err := getClientState(opts.accessibleState.GetStateDB(), connectionEnd.GetClientID())
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}
	if !found {
		return fmt.Errorf("client state not found in database")
	}

	// check if packet is timed out on the receiving chain
	latestHeight := clientState.GetLatestHeight()
	if !height.IsZero() && latestHeight.GTE(height) {
		return fmt.Errorf("receiving chain block height >= packet timeout height (%s >= %s), err: %w", latestHeight, height, channeltypes.ErrPacketTimeout)
	}

	consensusState, found, err := getConsensusState(opts.accessibleState.GetStateDB(), connectionEnd.ClientId, latestHeight)
	if err != nil {
		return fmt.Errorf("error loading consensus state, err: %w", err)
	}
	if !found {
		return fmt.Errorf("consensus state not found in database")
	}
	latestTimestamp := consensusState.GetTimestamp()
	if err != nil {
		return err
	}

	if packet.TimeoutTimestamp != 0 && latestTimestamp >= packet.TimeoutTimestamp {
		return fmt.Errorf("receiving chain block timestamp >= packet timeout timestamp (%s >= %s), err: %w", time.Unix(0, int64(latestTimestamp)), time.Unix(0, int64(packet.TimeoutTimestamp)), channeltypes.ErrPacketTimeout)
	}

	commitment := channeltypes.CommitPacket(marshaler, packet)

	setNextSequenceSend(opts.accessibleState, opts.args.SourcePort, opts.args.SourceChannel, sequence+1)
	setPacketCommitment(opts.accessibleState, opts.args.SourcePort, opts.args.SourceChannel, sequence, commitment)
	return nil
}

func _recvPacket(opts *callOpts[RecvPacketInput]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)



	relayer, err := sdk.AccAddressFromBech32(opts.args.Signer)
	if err != nil {
		ctx.Logger().Error("receive packet failed", "error", errorsmod.Wrap(err, "Invalid address for msg Signer"))
		return nil, errorsmod.Wrap(err, "Invalid address for msg Signer")
	}


	found, err := getCapability(opts.accessibleState.GetStateDB(), opts.args.Packet.DestinationPort, opts.args.Packet.DestinationChannel)
	if err != nil {
		return fmt.Errorf("%w, could not retrieve module from port-id", err)
	}
	if !found {
		return fmt.Errorf("module with port-ID: %s and channel-ID: %s, does not exist",  opts.args.Packet.DestinationPort, opts.args.Packet.DestinationChannel)
	}

	caller_from_port_id, found, err := getPort(opts.accessibleState.GetStateDB(), opts.args.Packet.DestinationPort)
	if err != nil {
		return fmt.Errorf("%w, could not retrieve module from port-id", err)
	}
	if !found {
		return fmt.Errorf("module with port-ID: %s and channel-ID: %s, does not exist",  opts.args.Packet.DestinationPort, opts.args.Packet.DestinationChannel)
	}




	err = k.ChannelKeeper.RecvPacket(cacheCtx, capability, msg.Packet, msg.ProofCommitment, msg.ProofHeight)


	channel, err := getChannelState(
		marshaler,
		string(hosttypes.ChannelKey(
			opts.args.Packet.DestinationPort,
			opts.args.Packet.DestinationChannel,
		)),
		opts.accessibleState,
	)
	if err != nil {
		return err
	}

	if channel.State != channeltypes.OPEN {
		fmt.Errorf("%w, channel state is not OPEN (got %s)", channeltypes.ErrInvalidChannelState, channel.State.String())
	}

	// packet must come from the channel's counterparty
	if opts.args.Packet.SourcePort != channel.Counterparty.PortId {
		return fmt.Errorf("%w, packet source port doesn't match the counterparty's port (%s ≠ %s)", channeltypes.ErrInvalidPacket, opts.args.Packet.SourcePort, channel.Counterparty.PortId)
	}

	if opts.args.Packet.SourceChannel != channel.Counterparty.ChannelId {
		return fmt.Errorf("%w, packet source channel doesn't match the counterparty's channel (%s ≠ %s)", channeltypes.ErrInvalidPacket, opts.args.Packet.SourceChannel, channel.Counterparty.ChannelId)
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return errorsmod.Wrap(connectiontypes.ErrConnectionNotFound, channel.ConnectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return errorsmod.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	// check if packet timeouted by comparing it with the latest height of the chain
	selfHeight := clienttypes.GetSelfHeight(ctx)
	timeoutHeight := packet.GetTimeoutHeight()
	if !timeoutHeight.IsZero() && selfHeight.GTE(timeoutHeight) {
		return errorsmod.Wrapf(
			types.ErrPacketTimeout,
			"block height >= packet timeout height (%s >= %s)", selfHeight, timeoutHeight,
		)
	}

	// check if packet timeouted by comparing it with the latest timestamp of the chain
	if packet.GetTimeoutTimestamp() != 0 && uint64(ctx.BlockTime().UnixNano()) >= packet.GetTimeoutTimestamp() {
		return errorsmod.Wrapf(
			types.ErrPacketTimeout,
			"block timestamp >= packet timeout timestamp (%s >= %s)", ctx.BlockTime(), time.Unix(0, int64(packet.GetTimeoutTimestamp())),
		)
	}

	commitment := types.CommitPacket(k.cdc, packet)

	// verify that the counterparty did commit to sending this packet
	if err := k.connectionKeeper.VerifyPacketCommitment(
		ctx, connectionEnd, proofHeight, proof,
		packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence(),
		commitment,
	); err != nil {
		return errorsmod.Wrap(err, "couldn't verify counterparty packet commitment")
	}

	switch channel.Ordering {
	case types.UNORDERED:
		// check if the packet receipt has been received already for unordered channels
		_, found := k.GetPacketReceipt(ctx, packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())
		if found {
			emitRecvPacketEvent(ctx, packet, channel)
			// This error indicates that the packet has already been relayed. Core IBC will
			// treat this error as a no-op in order to prevent an entire relay transaction
			// from failing and consuming unnecessary fees.
			return types.ErrNoOpMsg
		}

		// All verification complete, update state
		// For unordered channels we must set the receipt so it can be verified on the other side.
		// This receipt does not contain any data, since the packet has not yet been processed,
		// it's just a single store key set to an empty string to indicate that the packet has been received
		k.SetPacketReceipt(ctx, packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())

	case types.ORDERED:
		// check if the packet is being received in order
		nextSequenceRecv, found := k.GetNextSequenceRecv(ctx, packet.GetDestPort(), packet.GetDestChannel())
		if !found {
			return errorsmod.Wrapf(
				types.ErrSequenceReceiveNotFound,
				"destination port: %s, destination channel: %s", packet.GetDestPort(), packet.GetDestChannel(),
			)
		}

		if packet.GetSequence() < nextSequenceRecv {
			emitRecvPacketEvent(ctx, packet, channel)
			// This error indicates that the packet has already been relayed. Core IBC will
			// treat this error as a no-op in order to prevent an entire relay transaction
			// from failing and consuming unnecessary fees.
			return types.ErrNoOpMsg
		}

		if packet.GetSequence() != nextSequenceRecv {
			return errorsmod.Wrapf(
				types.ErrPacketSequenceOutOfOrder,
				"packet sequence ≠ next receive sequence (%d ≠ %d)", packet.GetSequence(), nextSequenceRecv,
			)
		}

		// All verification complete, update state
		// In ordered case, we must increment nextSequenceRecv
		nextSequenceRecv++

		// incrementing nextSequenceRecv and storing under this chain's channelEnd identifiers
		// Since this is the receiving chain, our channelEnd is packet's destination port and channel
		k.SetNextSequenceRecv(ctx, packet.GetDestPort(), packet.GetDestChannel(), nextSequenceRecv)
	}
	// emit an event that the relayer can query for
	emitRecvPacketEvent(ctx, packet, channel)






	ret, remainingGas, err := opts.accessibleState.CallFromPrecompile(ContractAddress, caller_from_port_id, PackIBCPacketMessage(), remainingGas, 0)

	err := k.ChannelKeeper.WriteAcknowledgement(ctx, capability, msg.Packet, ack)


	return &channeltypes.MsgRecvPacketResponse{Result: channeltypes.SUCCESS}, nil


	return nil
}