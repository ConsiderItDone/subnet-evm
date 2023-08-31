package ibc

import (
	"bytes"
	"fmt"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

func _sendPacket(opts *callOpts[MsgSendPacket]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	sequence, err := getNextSequenceSend(opts.accessibleState.GetStateDB(), opts.args.SourcePort, opts.args.SourceChannel)
	if err != nil {
		return err
	}

	channel, err := GetChannel(
		opts.accessibleState.GetStateDB(),
		opts.args.SourcePort,
		opts.args.SourceChannel,
	)
	if err != nil {
		return err
	}

	if channel.State != channeltypes.OPEN {
		return fmt.Errorf("%w, channel is not OPEN (got %s)", channeltypes.ErrInvalidChannelState, channel.State.String())
	}

	ok, err := GetCapability(opts.accessibleState.GetStateDB(), opts.args.SourcePort, opts.args.SourceChannel)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("caller does not own port capability for port ID %s, %w", opts.args.SourcePort, err)
	}

	height := clienttypes.Height{
		RevisionNumber: opts.args.TimeoutHeight.RevisionNumber.Uint64(),
		RevisionHeight: opts.args.TimeoutHeight.RevisionHeight.Uint64(),
	}

	packet := channeltypes.NewPacket(opts.args.Data, sequence, opts.args.SourcePort, opts.args.SourceChannel,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId, height, opts.args.TimeoutTimestamp.Uint64())

	connectionEnd, err := GetConnection(opts.accessibleState.GetStateDB(), channel.ConnectionHops[0])
	if err != nil {
		return err
	}

	clientState, err := GetClientState(opts.accessibleState.GetStateDB(), connectionEnd.GetClientID())
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}

	// check if packet is timed out on the receiving chain
	latestHeight := clientState.GetLatestHeight()
	if !height.IsZero() && latestHeight.GTE(height) {
		return fmt.Errorf("receiving chain block height >= packet timeout height (%s >= %s), err: %w", latestHeight, height, channeltypes.ErrPacketTimeout)
	}

	consensusState, err := GetConsensusState(opts.accessibleState.GetStateDB(), connectionEnd.ClientId, latestHeight)
	if err != nil {
		return fmt.Errorf("error loading consensus state, err: %w", err)
	}

	latestTimestamp := consensusState.GetTimestamp()
	if err != nil {
		return err
	}

	if packet.TimeoutTimestamp != 0 && latestTimestamp >= packet.TimeoutTimestamp {
		return fmt.Errorf("receiving chain block timestamp >= packet timeout timestamp (%s >= %s), err: %w", time.Unix(0, int64(latestTimestamp)), time.Unix(0, int64(packet.TimeoutTimestamp)), channeltypes.ErrPacketTimeout)
	}

	commitment := channeltypes.CommitPacket(marshaler, packet)

	topics, data, err := IBCABI.PackEvent(
		GeneratedPacketSentIdentifier.RawName,
		packet.Data,
		packet.TimeoutHeight.String(),
		big.NewInt(int64(packet.TimeoutTimestamp)),
		big.NewInt(int64(sequence)),
		packet.SourcePort,
		packet.SourceChannel,
		packet.DestinationPort,
		packet.DestinationChannel,
		channel.Ordering,
	)
	if err != nil {
		return fmt.Errorf("error packing event: %w", err)
	}
	blockNumber := opts.accessibleState.GetBlockContext().Number().Uint64()
	opts.accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

	setNextSequenceSend(opts.accessibleState.GetStateDB(), opts.args.SourcePort, opts.args.SourceChannel, sequence+1)
	setPacketCommitment(opts.accessibleState.GetStateDB(), opts.args.SourcePort, opts.args.SourceChannel, sequence, commitment)
	return nil
}

func _recvPacket(opts *callOpts[MsgRecvPacket]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	found, err := GetCapability(opts.accessibleState.GetStateDB(), opts.args.Packet.DestinationPort, opts.args.Packet.DestinationChannel)
	if err != nil {
		return fmt.Errorf("%w, could not retrieve module from port-id", err)
	}
	if !found {
		return fmt.Errorf("module with port-ID: %s and channel-ID: %s, does not exist", opts.args.Packet.DestinationPort, opts.args.Packet.DestinationChannel)
	}

	channel, err := GetChannel(
		opts.accessibleState.GetStateDB(),
		opts.args.Packet.DestinationPort,
		opts.args.Packet.DestinationChannel,
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

	connectionEnd, err := GetConnection(opts.accessibleState.GetStateDB(), channel.ConnectionHops[0])
	if err != nil {
		return err
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return fmt.Errorf("%w,connection state is not OPEN (got %s)", connectiontypes.ErrInvalidConnectionState, connectiontypes.State(connectionEnd.GetState()).String())
	}

	// check if packet timeouted by comparing it with the latest height of the chain
	selfHeight := clienttypes.Height{
		RevisionNumber: opts.args.Packet.TimeoutHeight.RevisionNumber.Uint64(), // TODO
		RevisionHeight: opts.accessibleState.GetBlockContext().Number().Uint64(),
	}
	timeoutHeight := clienttypes.Height{
		RevisionNumber: opts.args.Packet.TimeoutHeight.RevisionNumber.Uint64(),
		RevisionHeight: opts.args.Packet.TimeoutHeight.RevisionHeight.Uint64(),
	}

	if !timeoutHeight.IsZero() && selfHeight.GTE(timeoutHeight) {
		return fmt.Errorf("%w, block height >= packet timeout height (%s >= %s)", channeltypes.ErrPacketTimeout, selfHeight, timeoutHeight)
	}

	// check if packet timeouted by comparing it with the latest timestamp of the chain
	if opts.args.Packet.TimeoutTimestamp.Uint64() != 0 && opts.accessibleState.GetBlockContext().Timestamp().Uint64() >= opts.args.Packet.TimeoutTimestamp.Uint64() {
		return fmt.Errorf("block timestamp >= packet timeout timestamp (%d >= %d)", opts.accessibleState.GetBlockContext().Timestamp().Uint64(), int64(opts.args.Packet.TimeoutTimestamp.Uint64()))
	}

	packet := channeltypes.NewPacket(opts.args.Packet.Data, opts.args.Packet.Sequence.Uint64(), opts.args.Packet.SourcePort, opts.args.Packet.SourceChannel,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId, clienttypes.Height{
			RevisionNumber: opts.args.Packet.TimeoutHeight.RevisionNumber.Uint64(),
			RevisionHeight: opts.args.Packet.TimeoutHeight.RevisionHeight.Uint64(),
		}, opts.args.Packet.TimeoutTimestamp.Uint64())

	commitment := channeltypes.CommitPacket(marshaler, packet)

	height := clienttypes.Height{
		RevisionNumber: opts.args.ProofHeight.RevisionNumber.Uint64(),
		RevisionHeight: opts.args.ProofHeight.RevisionHeight.Uint64(),
	}

	// verify that the counterparty did commit to sending this packet
	if err := VerifyPacketCommitment(
		marshaler, connectionEnd, height, opts.args.ProofCommitment,
		packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence(),
		commitment, opts.accessibleState,
	); err != nil {
		return fmt.Errorf("%w, couldn't verify counterparty packet commitment", err)
	}

	switch channel.Ordering {
	case channeltypes.UNORDERED:

		// check if the packet receipt has been received already for unordered channels
		_, found := getPacketReceipt(opts.accessibleState.GetStateDB(), packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())
		if found {
			topics, data, err := IBCABI.PackEvent(GeneratedPacketReceivedIdentifier.RawName,
				packet.Data,
				packet.TimeoutHeight.String(),
				packet.TimeoutTimestamp,
				packet.Sequence,
				packet.SourcePort,
				packet.SourceChannel,
				packet.DestinationPort,
				packet.DestinationChannel,
				channel.Ordering.String(),
			)
			if err != nil {
				return fmt.Errorf("error packing event: %w", err)
			}
			blockNumber := opts.accessibleState.GetBlockContext().Number().Uint64()
			opts.accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

			// This error indicates that the packet has already been relayed. Core IBC will
			// treat this error as a no-op in order to prevent an entire relay transaction
			// from failing and consuming unnecessary fees.
			return channeltypes.ErrNoOpMsg
		}

		// All verification complete, update state
		// For unordered channels we must set the receipt so it can be verified on the other side.
		// This receipt does not contain any data, since the packet has not yet been processed,
		// it's just a single store key set to an empty string to indicate that the packet has been received
		setPacketReceipt(opts.accessibleState.GetStateDB(), packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())

	case channeltypes.ORDERED:
		// check if the packet is being received in order
		nextSequenceRecv, err := getNextSequenceRecv(opts.accessibleState.GetStateDB(), packet.GetDestPort(), packet.GetDestChannel())
		if err != nil {
			return err
		}

		if packet.GetSequence() < nextSequenceRecv {
			topics, data, err := IBCABI.PackEvent(GeneratedPacketReceivedIdentifier.RawName,
				packet.Data,
				packet.TimeoutHeight.String(),
				packet.TimeoutTimestamp,
				packet.Sequence,
				packet.SourcePort,
				packet.SourceChannel,
				packet.DestinationPort,
				packet.DestinationChannel,
				channel.Ordering.String(),
			)
			if err != nil {
				return fmt.Errorf("error packing event: %w", err)
			}
			blockNumber := opts.accessibleState.GetBlockContext().Number().Uint64()
			opts.accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

			// This error indicates that the packet has already been relayed. Core IBC will
			// treat this error as a no-op in order to prevent an entire relay transaction
			// from failing and consuming unnecessary fees.
			return channeltypes.ErrNoOpMsg
		}

		if packet.GetSequence() != nextSequenceRecv {
			return fmt.Errorf("%w, packet sequence ≠ next receive sequence (%d ≠ %d)", channeltypes.ErrPacketSequenceOutOfOrder, packet.GetSequence(), nextSequenceRecv)
		}

		// All verification complete, update state
		// In ordered case, we must increment nextSequenceRecv
		nextSequenceRecv++

		// incrementing nextSequenceRecv and storing under this chain's channelEnd identifiers
		// Since this is the receiving chain, our channelEnd is packet's destination port and channel
		setNextSequenceRecv(opts.accessibleState.GetStateDB(), packet.GetDestPort(), packet.GetDestChannel(), nextSequenceRecv)
	}
	// emit an event that the relayer can query for
	topics, data, err := IBCABI.PackEvent(GeneratedPacketReceivedIdentifier.RawName,
		packet.Data,
		packet.TimeoutHeight.String(),
		packet.TimeoutTimestamp,
		packet.Sequence,
		packet.SourcePort,
		packet.SourceChannel,
		packet.DestinationPort,
		packet.DestinationChannel,
		channel.Ordering.String(),
	)
	if err != nil {
		return fmt.Errorf("error packing event: %w", err)
	}
	blockNumber := opts.accessibleState.GetBlockContext().Number().Uint64()
	opts.accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

	return nil
}

func writeAcknowledgement(
	packet Packet,
	accessibleState contract.AccessibleState,
) error {
	channel, err := GetChannel(
		accessibleState.GetStateDB(),
		packet.DestinationPort,
		packet.DestinationChannel,
	)
	if err != nil {
		return err
	}

	if channel.State != channeltypes.OPEN {
		fmt.Errorf("%w, channel state is not OPEN (got %s)", channeltypes.ErrInvalidChannelState, channel.State.String())
	}

	if hasPacketAcknowledgement(accessibleState.GetStateDB(), packet.DestinationPort, packet.DestinationChannel, packet.Sequence.Uint64()) {
		return channeltypes.ErrAcknowledgementExists
	}

	// set the acknowledgement so that it can be verified on the other side
	setPacketAcknowledgement(
		accessibleState.GetStateDB(), packet.DestinationPort, packet.DestinationChannel, packet.Sequence.Uint64(),
		channeltypes.CommitAcknowledgement([]byte{1}),
	)

	return nil
}

func _timeout(opts *callOpts[MsgTimeout]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	found, err := GetCapability(opts.accessibleState.GetStateDB(), opts.args.Packet.SourcePort, opts.args.Packet.SourceChannel)
	if err != nil {
		return fmt.Errorf("%w, could not retrieve module from port-id", err)
	}
	if !found {
		return fmt.Errorf("module with port-ID: %s and channel-ID: %s, does not exist", opts.args.Packet.DestinationPort, opts.args.Packet.DestinationChannel)
	}

	height := clienttypes.Height{
		RevisionNumber: opts.args.Packet.TimeoutHeight.RevisionNumber.Uint64(),
		RevisionHeight: opts.args.Packet.TimeoutHeight.RevisionHeight.Uint64(),
	}

	channel, err := GetChannel(
		opts.accessibleState.GetStateDB(),
		opts.args.Packet.SourcePort,
		opts.args.Packet.SourceChannel,
	)
	if err != nil {
		return err
	}

	packet := channeltypes.NewPacket(opts.args.Packet.Data, opts.args.Packet.Sequence.Uint64(), opts.args.Packet.SourcePort, opts.args.Packet.SourceChannel,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId, height, opts.args.Packet.TimeoutTimestamp.Uint64())

	proofHeight := clienttypes.Height{
		RevisionNumber: opts.args.ProofHeight.RevisionNumber.Uint64(),
		RevisionHeight: opts.args.ProofHeight.RevisionHeight.Uint64(),
	}

	err = TimeoutPacket(marshaler, packet, opts.args.ProofUnreceived, proofHeight, opts.args.NextSequenceRecv.Uint64(), opts.accessibleState)
	return err
}

func _timeoutOnClose(opts *callOpts[MsgTimeoutOnClose]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	found, err := GetCapability(opts.accessibleState.GetStateDB(), opts.args.Packet.SourcePort, opts.args.Packet.SourceChannel)
	if err != nil {
		return fmt.Errorf("%w, could not retrieve module from port-id", err)
	}
	if !found {
		return fmt.Errorf("module with port-ID: %s and channel-ID: %s, does not exist", opts.args.Packet.DestinationPort, opts.args.Packet.DestinationChannel)
	}

	height := clienttypes.Height{
		RevisionNumber: opts.args.Packet.TimeoutHeight.RevisionNumber.Uint64(),
		RevisionHeight: opts.args.Packet.TimeoutHeight.RevisionHeight.Uint64(),
	}

	channel, err := GetChannel(
		opts.accessibleState.GetStateDB(),
		opts.args.Packet.SourcePort,
		opts.args.Packet.SourceChannel,
	)
	if err != nil {
		return err
	}

	packet := channeltypes.NewPacket(opts.args.Packet.Data, opts.args.Packet.Sequence.Uint64(), opts.args.Packet.SourcePort, opts.args.Packet.SourceChannel,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId, height, opts.args.Packet.TimeoutTimestamp.Uint64())

	if packet.GetDestPort() != channel.Counterparty.PortId {
		return fmt.Errorf("%w, packet destination port doesn't match the counterparty's port (%s ≠ %s)", channeltypes.ErrInvalidPacket, packet.GetDestPort(), channel.Counterparty.PortId)
	}

	if packet.GetDestChannel() != channel.Counterparty.ChannelId {
		return fmt.Errorf("%w, packet destination channel doesn't match the counterparty's channel (%s ≠ %s)", channeltypes.ErrInvalidPacket, packet.GetDestChannel(), channel.Counterparty.ChannelId)
	}

	connectionEnd, err := GetConnection(opts.accessibleState.GetStateDB(), channel.ConnectionHops[0])
	if err != nil {
		return err
	}

	commitment, err := getPacketCommitment(opts.accessibleState.GetStateDB(), packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	if err != nil {
		return err
	}

	if len(commitment) == 0 {
		// emit an event marking that we have processed the timeout
		topics, data, err := IBCABI.PackEvent(GeneratedTimeoutPacketIdentifier.RawName,
			opts.args.Packet.TimeoutTimestamp,
			opts.args.Packet.Sequence,
			opts.args.Packet.SourcePort,
			opts.args.Packet.SourceChannel,
			opts.args.Packet.DestinationPort,
			opts.args.Packet.DestinationChannel,
			channel.Ordering.String(),
			channel.ConnectionHops[0],
		)
		if err != nil {
			return fmt.Errorf("error packing event: %w", err)
		}
		blockNumber := opts.accessibleState.GetBlockContext().Number().Uint64()
		opts.accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

		// This error indicates that the timeout has already been relayed
		// or there is a misconfigured relayer attempting to prove a timeout
		// for a packet never sent. Core IBC will treat this error as a no-op in order to
		// prevent an entire relay transaction from failing and consuming unnecessary fees.
		return channeltypes.ErrNoOpMsg
	}

	packetCommitment := channeltypes.CommitPacket(marshaler, packet)

	// verify we sent the packet and haven't cleared it out yet
	if !bytes.Equal(commitment, packetCommitment) {
		return fmt.Errorf("%w, packet commitment bytes are not equal: got (%v), expected (%v)", channeltypes.ErrInvalidPacket, commitment, packetCommitment)
	}

	counterpartyHops := []string{connectionEnd.GetCounterparty().GetConnectionID()}

	counterparty := channeltypes.NewCounterparty(packet.GetSourcePort(), packet.GetSourceChannel())
	expectedChannel := channeltypes.NewChannel(
		channeltypes.CLOSED, channel.Ordering, counterparty, counterpartyHops, channel.Version,
	)

	proofHeight := clienttypes.Height{
		RevisionNumber: opts.args.ProofHeight.RevisionNumber.Uint64(),
		RevisionHeight: opts.args.ProofHeight.RevisionHeight.Uint64(),
	}

	// check that the opposing channel end has closed
	if err := VerifyChannelState(
		marshaler, connectionEnd, proofHeight, opts.args.ProofClose,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId,
		expectedChannel, opts.accessibleState,
	); err != nil {
		return err
	}

	switch channel.Ordering {
	case channeltypes.ORDERED:
		// check that packet has not been received
		if opts.args.NextSequenceRecv.Uint64() > packet.GetSequence() {
			return fmt.Errorf("%w, packet already received, next sequence receive > packet sequence (%d > %d", channeltypes.ErrInvalidPacket, opts.args.NextSequenceRecv.Uint64(), packet.GetSequence())
		}

		// check that the recv sequence is as claimed
		err = VerifyNextSequenceRecv(
			marshaler, connectionEnd, proofHeight, opts.args.ProofUnreceived,
			packet.GetDestPort(), packet.GetDestChannel(), opts.args.NextSequenceRecv.Uint64(), opts.accessibleState,
		)
	case channeltypes.UNORDERED:
		err = VerifyPacketReceiptAbsence(
			marshaler, connectionEnd, proofHeight, opts.args.ProofUnreceived,
			packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence(), opts.accessibleState,
		)
	default:
		panic(fmt.Errorf("%w, %s", channeltypes.ErrInvalidChannelOrdering, channel.Ordering.String()))
	}

	if err != nil {
		return err
	}

	return nil
}

func _acknowledgement(opts *callOpts[MsgAcknowledgement]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	found, err := GetCapability(opts.accessibleState.GetStateDB(), opts.args.Packet.SourcePort, opts.args.Packet.SourceChannel)
	if err != nil {
		return fmt.Errorf("%w, could not retrieve module from port-id", err)
	}
	if !found {
		return fmt.Errorf("module with port-ID: %s and channel-ID: %s, does not exist", opts.args.Packet.DestinationPort, opts.args.Packet.DestinationChannel)
	}

	channel, err := GetChannel(
		opts.accessibleState.GetStateDB(),
		opts.args.Packet.SourcePort,
		opts.args.Packet.SourceChannel,
	)
	if err != nil {
		return err
	}

	height := clienttypes.Height{
		RevisionNumber: opts.args.Packet.TimeoutHeight.RevisionNumber.Uint64(),
		RevisionHeight: opts.args.Packet.TimeoutHeight.RevisionHeight.Uint64(),
	}

	packet := channeltypes.NewPacket(opts.args.Packet.Data, opts.args.Packet.Sequence.Uint64(), opts.args.Packet.SourcePort, opts.args.Packet.SourceChannel,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId, height, opts.args.Packet.TimeoutTimestamp.Uint64())

	if channel.State != channeltypes.OPEN {
		return fmt.Errorf("%w, channel state is not OPEN (got %s)", channeltypes.ErrInvalidChannelState, channel.State.String())
	}

	found, err = GetCapability(opts.accessibleState.GetStateDB(), opts.args.Packet.SourcePort, opts.args.Packet.SourceChannel)
	if err != nil {
		return fmt.Errorf("%w, could not retrieve module from port-id", err)
	}
	if !found {
		return fmt.Errorf("module with port-ID: %s and channel-ID: %s, does not exist", opts.args.Packet.DestinationPort, opts.args.Packet.DestinationChannel)
	}

	// packet must have been sent to the channel's counterparty
	if packet.GetDestPort() != channel.Counterparty.PortId {
		return fmt.Errorf("%w, packet destination port doesn't match the counterparty's port (%s ≠ %s)", channeltypes.ErrInvalidPacket, packet.GetDestPort(), channel.Counterparty.PortId)
	}

	if packet.GetDestChannel() != channel.Counterparty.ChannelId {
		return fmt.Errorf("%w, packet destination channel doesn't match the counterparty's channel (%s ≠ %s)", channeltypes.ErrInvalidPacket, packet.GetDestChannel(), channel.Counterparty.ChannelId)
	}

	connectionEnd, err := GetConnection(opts.accessibleState.GetStateDB(), channel.ConnectionHops[0])
	if err != nil {
		return err
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return fmt.Errorf("%w, connection state is not OPEN (got %s)", connectiontypes.ErrInvalidConnectionState, connectiontypes.State(connectionEnd.GetState()).String())
	}

	commitment, err := getPacketCommitment(opts.accessibleState.GetStateDB(), packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	if err != nil {
		return err
	}

	if len(commitment) == 0 {
		topics, data, err := IBCABI.PackEvent(GeneratedAcknowledgePacketIdentifier.RawName,
			packet.TimeoutHeight.String(),
			opts.args.Packet.TimeoutTimestamp,
			opts.args.Packet.Sequence,
			opts.args.Packet.SourcePort,
			opts.args.Packet.SourceChannel,
			opts.args.Packet.DestinationPort,
			opts.args.Packet.DestinationChannel,
			channel.Ordering.String(),
			channel.ConnectionHops[0],
		)
		if err != nil {
			return fmt.Errorf("error packing event: %w", err)
		}
		blockNumber := opts.accessibleState.GetBlockContext().Number().Uint64()
		opts.accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

		// This error indicates that the acknowledgement has already been relayed
		// or there is a misconfigured relayer attempting to prove an acknowledgement
		// for a packet never sent. Core IBC will treat this error as a no-op in order to
		// prevent an entire relay transaction from failing and consuming unnecessary fees.
		return channeltypes.ErrNoOpMsg
	}

	packetCommitment := channeltypes.CommitPacket(marshaler, packet)

	// verify we sent the packet and haven't cleared it out yet
	if !bytes.Equal(commitment, packetCommitment) {
		return fmt.Errorf("%w, commitment bytes are not equal: got (%v), expected (%v)", channeltypes.ErrInvalidPacket, packetCommitment, commitment)
	}

	height = clienttypes.Height{
		RevisionNumber: opts.args.ProofHeight.RevisionNumber.Uint64(),
		RevisionHeight: opts.args.ProofHeight.RevisionHeight.Uint64(),
	}

	if err := VerifyPacketAcknowledgement(
		marshaler,
		connectionEnd,
		height,
		opts.args.ProofAcked,
		packet.GetDestPort(),
		packet.GetDestChannel(),
		packet.GetSequence(),
		opts.args.Acknowledgement,
		opts.accessibleState,
	); err != nil {
		return err
	}

	// assert packets acknowledged in order
	if channel.Ordering == channeltypes.ORDERED {

		nextSequenceAck, err := getNextSequenceAck(opts.accessibleState.GetStateDB(), opts.args.Packet.SourcePort, opts.args.Packet.SourceChannel)
		if err != nil {
			return err
		}

		if opts.args.Packet.Sequence.Uint64() != nextSequenceAck {
			return fmt.Errorf("%w, packet sequence ≠ next ack sequence (%d ≠ %d)", channeltypes.ErrPacketSequenceOutOfOrder, opts.args.Packet.Sequence.Uint64(), nextSequenceAck)
		}

		// All verification complete, in the case of ORDERED channels we must increment nextSequenceAck
		nextSequenceAck++

		setNextSequenceAck(opts.accessibleState.GetStateDB(), opts.args.Packet.SourcePort, opts.args.Packet.SourceChannel, opts.args.Packet.Sequence.Uint64())
	}

	// Delete packet commitment, since the packet has been acknowledged, the commitement is no longer necessary
	deletePacketCommitment(opts.accessibleState.GetStateDB(), opts.args.Packet.SourcePort, opts.args.Packet.SourceChannel, opts.args.Packet.Sequence.Uint64())

	// emit an event marking that we have processed the acknowledgement
	topics, data, err := IBCABI.PackEvent(GeneratedAcknowledgePacketIdentifier.RawName,
		packet.TimeoutHeight.String(),
		opts.args.Packet.TimeoutTimestamp,
		opts.args.Packet.Sequence,
		opts.args.Packet.SourcePort,
		opts.args.Packet.SourceChannel,
		opts.args.Packet.DestinationPort,
		opts.args.Packet.DestinationChannel,
		channel.Ordering.String(),
		channel.ConnectionHops[0],
	)
	if err != nil {
		return fmt.Errorf("error packing event: %w", err)
	}
	blockNumber := opts.accessibleState.GetBlockContext().Number().Uint64()
	opts.accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

	return nil
}

func TimeoutPacket(
	marshaler *codec.ProtoCodec,
	packet exported.PacketI,
	proof []byte,
	proofHeight exported.Height,
	nextSequenceRecv uint64,
	accessibleState contract.AccessibleState,
) error {

	channel, err := GetChannel(
		accessibleState.GetStateDB(),
		packet.GetSourcePort(),
		packet.GetSourceChannel(),
	)
	if err != nil {
		return err
	}
	// NOTE: TimeoutPacket is called by the AnteHandler which acts upon the packet.Route(),
	// so the capability authentication can be omitted here

	if packet.GetDestPort() != channel.Counterparty.PortId {
		return fmt.Errorf("%w, packet destination port doesn't match the counterparty's port (%s ≠ %s)", channeltypes.ErrInvalidPacket, packet.GetDestPort(), channel.Counterparty.PortId)
	}

	if packet.GetDestChannel() != channel.Counterparty.ChannelId {
		return fmt.Errorf("%w, packet destination channel doesn't match the counterparty's channel (%s ≠ %s)", channeltypes.ErrInvalidPacket, packet.GetDestChannel(), channel.Counterparty.ChannelId)
	}

	connectionEnd, err := GetConnection(accessibleState.GetStateDB(), channel.ConnectionHops[0])
	if err != nil {
		return err
	}

	// check that timeout height or timeout timestamp has passed on the other end
	consensusState, err := GetConsensusState(accessibleState.GetStateDB(), connectionEnd.ClientId, proofHeight)
	if err != nil {
		return fmt.Errorf("%w, %w, please ensure the proof was constructed against a height that exists on the client", clienttypes.ErrConsensusStateNotFound, err)
	}
	proofTimestamp := consensusState.GetTimestamp()

	timeoutHeight := packet.GetTimeoutHeight()
	if (timeoutHeight.IsZero() || proofHeight.LT(timeoutHeight)) &&
		(packet.GetTimeoutTimestamp() == 0 || proofTimestamp < packet.GetTimeoutTimestamp()) {
		return fmt.Errorf("%w, packet timeout has not been reached for height or timestamp", channeltypes.ErrPacketTimeout)
	}

	commitment, err := getPacketCommitment(accessibleState.GetStateDB(), packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	if err != nil {
		return err
	}

	if len(commitment) == 0 {
		// emit an event marking that we have processed the timeout
		topics, data, err := IBCABI.PackEvent(GeneratedTimeoutPacketIdentifier.RawName,
			packet.GetTimeoutTimestamp(),
			packet.GetSequence(),
			packet.GetSourcePort(),
			packet.GetSourceChannel(),
			packet.GetDestPort(),
			packet.GetDestChannel(),
			channel.Ordering.String(),
			channel.ConnectionHops[0],
		)
		if err != nil {
			return fmt.Errorf("error packing event: %w", err)
		}
		blockNumber := accessibleState.GetBlockContext().Number().Uint64()
		accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

		// This error indicates that the timeout has already been relayed
		// or there is a misconfigured relayer attempting to prove a timeout
		// for a packet never sent. Core IBC will treat this error as a no-op in order to
		// prevent an entire relay transaction from failing and consuming unnecessary fees.
		return channeltypes.ErrNoOpMsg
	}

	if channel.State != channeltypes.OPEN {
		return fmt.Errorf("%w, channel state is not OPEN (got %s)", channeltypes.ErrInvalidChannelState, channel.State.String())
	}

	packetCommitment := channeltypes.CommitPacket(marshaler, packet)

	// verify we sent the packet and haven't cleared it out yet
	if !bytes.Equal(commitment, packetCommitment) {
		return fmt.Errorf("%w, packet commitment bytes are not equal: got (%v), expected (%v)", channeltypes.ErrInvalidPacket, commitment, packetCommitment)
	}

	switch channel.Ordering {
	case channeltypes.ORDERED:
		// check that packet has not been received
		if nextSequenceRecv > packet.GetSequence() {
			return fmt.Errorf("%w, packet already received, next sequence receive > packet sequence (%d > %d)", channeltypes.ErrPacketReceived, nextSequenceRecv, packet.GetSequence())
		}

		// check that the recv sequence is as claimed
		err = VerifyNextSequenceRecv(
			marshaler, connectionEnd, proofHeight, proof,
			packet.GetDestPort(), packet.GetDestChannel(), nextSequenceRecv, accessibleState,
		)
	case channeltypes.UNORDERED:
		err = VerifyPacketReceiptAbsence(
			marshaler, connectionEnd, proofHeight, proof,
			packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence(), accessibleState,
		)
	default:
		panic(fmt.Errorf("%w, %s", channeltypes.ErrInvalidChannelOrdering, channel.Ordering.String()))
	}

	if err != nil {
		return err
	}

	return nil
}

func TimeoutExecuted(
	accessibleState contract.AccessibleState,
	packet Packet,
) error {

	channel, err := GetChannel(
		accessibleState.GetStateDB(),
		packet.SourcePort,
		packet.SourceChannel,
	)
	if err != nil {
		return err
	}

	ok, err := GetCapability(accessibleState.GetStateDB(), packet.SourcePort, packet.SourceChannel)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("caller does not own port capability for port ID %s, %w", packet.SourcePort, err)
	}

	deletePacketCommitment(accessibleState.GetStateDB(), packet.SourcePort, packet.SourceChannel, packet.Sequence.Uint64())

	if channel.Ordering == channeltypes.ORDERED {
		channel.State = channeltypes.CLOSED
		SetChannel(accessibleState.GetStateDB(),
			packet.SourcePort,
			packet.SourceChannel,
			channel,
		)
	}

	// emit an event marking that we have processed the timeout
	topics, data, err := IBCABI.PackEvent(GeneratedTimeoutPacketIdentifier.RawName,
		packet.TimeoutTimestamp,
		packet.Sequence,
		packet.SourcePort,
		packet.SourceChannel,
		packet.DestinationPort,
		packet.DestinationChannel,
		channel.Ordering.String(),
		channel.ConnectionHops[0],
	)
	if err != nil {
		return fmt.Errorf("error packing event: %w", err)
	}
	blockNumber := accessibleState.GetBlockContext().Number().Uint64()
	accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)

	if channel.Ordering == channeltypes.ORDERED && channel.State == channeltypes.CLOSED {
		// TODO
		// emitChannelClosedEvent(ctx, packet, channel)
	}

	return nil
}
