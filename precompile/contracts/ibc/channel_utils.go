package ibc

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	hosttypes "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

func generateChannelIdentifier(accessibleState contract.AccessibleState) string {
	sequence := uint64(0)
	if accessibleState.GetStateDB().Exist(common.BytesToAddress([]byte("nextChannelSeq"))) {
		b := accessibleState.GetStateDB().GetPrecompileState(common.BytesToAddress([]byte("nextChannelSeq")))
		sequence = binary.BigEndian.Uint64(b)
	}
	channelId := fmt.Sprintf("%s%d", "channel-", sequence)
	sequence++
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, sequence)
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte("nextChannelSeq")), b)

	return channelId
}

// setNextSequenceSend sets a channel's next send sequence to the store
func setNextSequenceSend(accessibleState contract.AccessibleState, portID, channelID string, sequence uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, sequence)
	key := calculateKey([]byte(hosttypes.NextSequenceSendKey(portID, channelID)))
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(key)), b)
}

// setNextSequenceRecv sets a channel's next receive sequence to the store
func setNextSequenceRecv(accessibleState contract.AccessibleState, portID, channelID string, sequence uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, sequence)
	key := calculateKey([]byte(hosttypes.NextSequenceRecvKey(portID, channelID)))
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(key)), b)
}

// setNextSequenceAck sets a channel's next ack sequence to the store
func setNextSequenceAck(accessibleState contract.AccessibleState, portID, channelID string, sequence uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, sequence)

	key := calculateKey([]byte(hosttypes.NextSequenceAckKey(portID, channelID)))
	accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(key)), b)
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
	bz, err := channelNew.Marshal()
	if err != nil {
		return fmt.Errorf("can't serialize channel state: %w", err)
	}
	path := calculateKey(hosttypes.ChannelKey(opts.args.PortID, channelID))
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(path)), bz)

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

	clientState, clientStateFound, err := getClientState(accessibleState.GetStateDB(), clientID)
	if err != nil {
		return fmt.Errorf("can't read client state: %w", err)
	}
	if !clientStateFound {
		return fmt.Errorf("client state not found: %s", clientID)
	}

	consensusState, consensusStateFound, err := getConsensusState(accessibleState.GetStateDB(), clientID, clientState.GetLatestHeight())
	if err != nil {
		return fmt.Errorf("can't read consensus state: %w", err)
	}
	if !consensusStateFound {
		return fmt.Errorf("consensus state not found: %s", clientID)
	}

	merklePath := commitmenttypes.NewMerklePath(hosttypes.ChannelPath(portID, channelID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	bz, err := channel.Marshal()
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

	return merkleProof.VerifyMembership(clientState.ProofSpecs, consensusState.GetRoot(), merklePath, bz)
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

	err = makeCapability(opts.accessibleState.GetStateDB(), opts.args.PortID, channelID)
	if err != nil {
		return "", fmt.Errorf("can't make capability [%s,%s]: %w", opts.args.PortID, channelID, err)
	}

	channelNew := channeltypes.NewChannel(channeltypes.TRYOPEN, channel.Ordering, channel.Counterparty, channel.ConnectionHops, channel.Version)
	bz := marshaler.MustMarshal(&channelNew)
	path := calculateKey(hosttypes.ChannelKey(
		opts.args.PortID,
		channelID,
	))
	opts.accessibleState.GetStateDB().SetPrecompileState(
		common.BytesToAddress([]byte(path)),
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
	key := common.BytesToAddress([]byte(calculateKey([]byte(channelStatePath))))
	exist := accessibleState.GetStateDB().Exist(key)
	if !exist {
		return nil, fmt.Errorf("cannot find channel state with path: %s", channelStatePath)
	}
	channelStateByte := accessibleState.GetStateDB().GetPrecompileState(key)
	channelState := &channeltypes.Channel{}
	marshaler.MustUnmarshal(channelStateByte, channelState)
	return channelState, nil
}

func _channelOpenAck(opts *callOpts[ChannelOpenAckInput]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	channeltypes.RegisterInterfaces(interfaceRegistry)
	clienttypes.RegisterInterfaces(interfaceRegistry)
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
		return fmt.Errorf("channel handshake open ack failed: %w", err)
	}

	channel.State = channeltypes.OPEN
	channel.Version = opts.args.CounterpartyVersion
	channel.Counterparty.ChannelId = opts.args.CounterpartyChannelID

	bz := marshaler.MustMarshal(channel)
	key := calculateKey(hosttypes.ChannelKey(
		opts.args.PortID,
		opts.args.ChannelID,
	))
	opts.accessibleState.GetStateDB().SetPrecompileState(common.BytesToAddress([]byte(key)), bz)

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
