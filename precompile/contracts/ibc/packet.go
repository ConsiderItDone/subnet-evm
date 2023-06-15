package ibc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	"github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	hosttypes "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
)

var (
	sendPacketSignature      = contract.CalculateFunctionSelector("sendPacket(uint64,bytes,uint64,bytes,bytes,uint64,bytes)")
	receivePacketSignature   = contract.CalculateFunctionSelector("receivePacket(uint64,bytes,uint64,bytes,uint64,bytes)")
	acknowledgementSignature = contract.CalculateFunctionSelector("receivePacket(uint64,bytes,uint64,bytes,uint64,bytes,uint64,bytes)")
)

func SendPacket(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
			input
			8 byte             - sourcePortLen
			clientIDLen byte   - sourcePortByte
			8 byte             - sourceChannelLen
			clientMessageLen   - sourceChannelByte
		    8 byte             - timeoutHeightLen
		    proofHeightbyte    - clienttypes.Height
		    8 byte             - timeoutTimestamp
			8 byte             - DataLen
		    data               - []byte

	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, sendPacketGas); err != nil {
		return nil, 0, fmt.Errorf("rrror DeductGas err: %w", err)
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

	// sourcePort
	carriage := uint64(0)
	sourcePortLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	sourcePort := string(getData(input, carriage, sourcePortLen))
	carriage = carriage + sourcePortLen

	// sourceChannel
	sourceChannelLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	sourceChannel := string(getData(input, carriage, sourceChannelLen))
	carriage = carriage + sourcePortLen

	// timeoutHeight
	timeoutHeightLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	timeoutHeightByte := getData(input, carriage, timeoutHeightLen)

	timeoutHeight := &clienttypes.Height{}
	err = marshaler.Unmarshal(timeoutHeightByte, timeoutHeight)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	// timeoutTimestamp
	timeoutTimestamp := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8

	// data
	DataLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	data := getData(input, carriage, DataLen)
	carriage = carriage + DataLen

	channel, err := getChannelState(marshaler, string(hosttypes.ChannelKey(sourcePort, sourceChannel)), accessibleState)
	if err != nil {
		return nil, 0, err
	}

	if channel.State != types.OPEN {
		return nil, 0, fmt.Errorf("%s: channel is not OPEN (got %s)", types.ErrInvalidChannelState, channel.State.String())
	}

	sequence, found := GetNextSequenceSend(accessibleState, sourcePort, sourceChannel)
	if !found {
		return nil, 0, fmt.Errorf("%s: source port: %s, source channel: %s", types.ErrSequenceSendNotFound, sourcePort, sourceChannel)
	}

	// construct packet from given fields and channel state
	packet := types.NewPacket(data, sequence, sourcePort, sourceChannel,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId, *timeoutHeight, timeoutTimestamp)

	if err := packet.ValidateBasic(); err != nil {
		return nil, 0, fmt.Errorf("%s: constructed packet failed basic validation", err)
	}

	commitment := types.CommitPacket(marshaler, packet)

	SetNextSequenceSend(accessibleState, sourcePort, sourceChannel, sequence+1)
	SetPacketCommitment(accessibleState, sourcePort, sourceChannel, packet.GetSequence(), commitment)

	var result []byte
	binary.BigEndian.PutUint64(result, sequence)

	return result, remainingGas, nil
}

func ReceivePacket(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
			input
			8 byte             - packetLen
			clientIDLen byte   - packetByte
			8 byte             - proofLen
			clientMessageLen   - proofByte
		    8 byte             - proofHeightLen
		    proofHeightbyte    - clienttypes.Height

	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, receivePacketGas); err != nil {
		return nil, 0, fmt.Errorf("rrror DeductGas err: %w", err)
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

	// sourcePort
	carriage := uint64(0)
	sourcePortLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	sourcePort := string(getData(input, carriage, sourcePortLen))
	carriage = carriage + sourcePortLen

	// sourceChannel
	sourceChannelLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	sourceChannel := string(getData(input, carriage, sourceChannelLen))
	carriage = carriage + sourcePortLen

	// timeoutHeight
	timeoutHeightLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	timeoutHeightByte := getData(input, carriage, timeoutHeightLen)

	timeoutHeight := &clienttypes.Height{}
	err = marshaler.Unmarshal(timeoutHeightByte, timeoutHeight)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	// timeoutTimestamp
	timeoutTimestamp := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8

	// data
	DataLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	data := getData(input, carriage, DataLen)
	carriage = carriage + DataLen

	channel, err := getChannelState(marshaler, string(hosttypes.ChannelKey(sourcePort, sourceChannel)), accessibleState)
	if err != nil {
		return nil, 0, err
	}

	if channel.State != types.OPEN {
		return nil, 0, fmt.Errorf("%s: channel is not OPEN (got %s)", types.ErrInvalidChannelState, channel.State.String())
	}

	sequence, found := GetNextSequenceSend(accessibleState, sourcePort, sourceChannel)
	if !found {
		return nil, 0, fmt.Errorf("%s: source port: %s, source channel: %s", types.ErrSequenceSendNotFound, sourcePort, sourceChannel)
	}

	// construct packet from given fields and channel state
	packet := types.NewPacket(data, sequence, sourcePort, sourceChannel,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId, *timeoutHeight, timeoutTimestamp)

	if err := packet.ValidateBasic(); err != nil {
		return nil, 0, fmt.Errorf("%s: constructed packet failed basic validation", err)
	}

	commitment := types.CommitPacket(marshaler, packet)

	SetNextSequenceSend(accessibleState, sourcePort, sourceChannel, sequence+1)
	SetPacketCommitment(accessibleState, sourcePort, sourceChannel, packet.GetSequence(), commitment)

	var result []byte
	binary.BigEndian.PutUint64(result, sequence)

	return result, remainingGas, nil
}

func Acknowledgement(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	/*
			input
			8 byte          - packetLen
			packet          - Packet
			8 byte          - acknowledgementLen
			acknowledgement - []byte
			8 byte          - proofAckedLen
			proofAcked      - []byte
		    8 byte          - proofHeightLen
			proofHeight     - types.Height
	*/
	if remainingGas, err = contract.DeductGas(suppliedGas, acknowledgementGas); err != nil {
		return nil, 0, fmt.Errorf("rrror DeductGas err: %w", err)
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
	// Packet
	packetLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	packetByte := getData(input, carriage, packetLen)
	carriage = carriage + packetLen

	packet := &types.Packet{}
	err = marshaler.Unmarshal(packetByte, packet)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling packet: %w", err)
	}

	// acknowledgement
	acknowledgementLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	acknowledgement := getData(input, carriage, acknowledgementLen)
	carriage = carriage + acknowledgementLen

	// acknowledgement
	proofAckedLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofAcked := getData(input, carriage, proofAckedLen)
	carriage = carriage + proofAckedLen

	// timeoutHeight
	proofHeightLen := new(big.Int).SetBytes(getData(input, carriage, 8)).Uint64()
	carriage = carriage + 8
	proofHeightByte := getData(input, carriage, proofHeightLen)

	proofHeight := &clienttypes.Height{}
	err = marshaler.Unmarshal(proofHeightByte, proofHeight)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling proofHeight: %w", err)
	}

	channel, err := getChannelState(marshaler, string(hosttypes.ChannelKey(packet.GetSourcePort(), packet.GetSourceChannel())), accessibleState)
	if err != nil {
		return nil, 0, err
	}

	if channel.State != types.OPEN {
		return nil, 0, fmt.Errorf("channel state is not OPEN (got %s): %w", channel.State.String(), err)
	}

	// packet must have been sent to the channel's counterparty
	if packet.GetDestPort() != channel.Counterparty.PortId {
		return nil, 0, fmt.Errorf("packet destination port doesn't match the counterparty's port (%s ≠ %s)", packet.GetDestPort(), channel.Counterparty.PortId)
	}

	if packet.GetDestChannel() != channel.Counterparty.ChannelId {
		return nil, 0, fmt.Errorf("packet destination channel doesn't match the counterparty's channel (%s ≠ %s)", packet.GetDestChannel(), channel.Counterparty.ChannelId)
	}

	connectionsPath := fmt.Sprintf("connections/%s", channel.ConnectionHops[0])
	connectionEnd, err := getConnection(marshaler, connectionsPath, accessibleState)
	if err != nil {
		return nil, 0, err
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return nil, 0, fmt.Errorf("connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String())
	}

	commitment, err := GetPacketCommitment(accessibleState, packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	if err != nil {
		return nil, 0, err
	}

	if len(commitment) == 0 {
		return nil, 0, fmt.Errorf("message is redundant, no-op will be performed")
	}

	packetCommitment := types.CommitPacket(marshaler, packet)

	// verify we sent the packet and haven't cleared it out yet
	if !bytes.Equal(commitment, packetCommitment) {
		return nil, 0, fmt.Errorf("commitment bytes are not equal: got (%v), expected (%v)", packetCommitment, commitment)
	}

	clientID := connectionEnd.GetClientID()

	clientStatePath := fmt.Sprintf("clients/%s/clientState", clientID)
	clientState, err := getClientState(marshaler, clientStatePath, accessibleState)
	if err != nil {
		return nil, 0, err
	}

	// get time and block delays
	// timeDelay := connectionEnd.GetDelayPeriod()

	// TODO
	// expectedTimePerBlock := GetMaxExpectedTimePerBlock()
	// if expectedTimePerBlock == 0 {
	// 	return nil, 0, err
	// }
	// blockDelay := uint64(math.Ceil(float64(timeDelay) / float64(expectedTimePerBlock)))

	merklePath := commitmenttypes.NewMerklePath(hosttypes.PacketAcknowledgementPath(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence()))
	merklePath, err = commitmenttypes.ApplyPrefix(connectionEnd.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return nil, 0, err
	}

	if clientState.GetLatestHeight().LT(*proofHeight) {
		return nil, 0, fmt.Errorf("client state height < proof height (%d < %d), please ensure the client has been updated", clientState.GetLatestHeight(), proofHeight)
	}

	// if timeDelay != 0 {
	// TODO
	// // check that executing chain's timestamp has passed consensusState's processed time + delay time period
	// processedTime, ok := GetProcessedTime(store, proofHeight)
	// if !ok {
	// 	return nil, 0, fmt.Errorf("processed time not found for height: %s", proofHeight)
	// }

	// currentTimestamp := uint64(ctx.BlockTime().UnixNano())
	// validTime := processedTime + timeDelay

	// // NOTE: delay time period is inclusive, so if currentTimestamp is validTime, then we return no error
	// if currentTimestamp < validTime {
	// 	return nil, 0, fmt.Errorf("cannot verify packet until time: %d, current time: %d", validTime, currentTimestamp)
	// }
	// }

	// if blockDelay != 0 {
	// TODO
	// // check that executing chain's height has passed consensusState's processed height + delay block period
	// processedHeight, ok := GetProcessedHeight(store, proofHeight)
	// if !ok {
	// 	return nil, 0, fmt.Errorf("processed height not found for height: %s", proofHeight)
	// }

	// currentHeight := clienttypes.GetSelfHeight(ctx)
	// validHeight := clienttypes.NewHeight(processedHeight.GetRevisionNumber(), processedHeight.GetRevisionHeight()+blockDelay)

	// // NOTE: delay block period is inclusive, so if currentHeight is validHeight, then we return no error
	// if currentHeight.LT(validHeight) {
	// 	return nil, 0, fmt.Errorf("cannot verify packet until height: %s, current height: %s",
	// 	validHeight, currentHeight)
	// }
	// }

	var merkleProof commitmenttypes.MerkleProof
	if err := marshaler.Unmarshal(proofAcked, &merkleProof); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal proof into ICS 23 commitment merkle proof")
	}

	consensusStatePath := fmt.Sprintf("clients/%s/consensusStates/%s", clientID, proofHeight)
	consensusState, err := getConsensusState(marshaler, consensusStatePath, accessibleState)

	if err != nil {
		return nil, 0, fmt.Errorf("please ensure the proof was constructed against a height that exists on the client")
	}

	err = merkleProof.VerifyMembership(clientState.ProofSpecs, consensusState.GetRoot(), merklePath, types.CommitAcknowledgement(acknowledgement))
	if err != nil {
		return nil, 0, fmt.Errorf("failed packet acknowledgement verification for client (%s)", clientID)
	}

	// assert packets acknowledged in order
	if channel.Ordering == types.ORDERED {
		nextSequenceAck, found := GetNextSequenceAck(accessibleState, packet.GetSourcePort(), packet.GetSourceChannel())
		if !found {
			return nil, 0, fmt.Errorf("source port: %s, source channel: %s", packet.GetSourcePort(), packet.GetSourceChannel())
		}

		if packet.GetSequence() != nextSequenceAck {
			return nil, 0, fmt.Errorf("packet sequence ≠ next ack sequence (%d ≠ %d)", packet.GetSequence(), nextSequenceAck)
		}

		// All verification complete, in the case of ORDERED channels we must increment nextSequenceAck
		nextSequenceAck++

		// incrementing NextSequenceAck and storing under this chain's channelEnd identifiers
		// Since this is the original sending chain, our channelEnd is packet's source port and channel
		SetNextSequenceAck(accessibleState, packet.GetSourcePort(), packet.GetSourceChannel(), nextSequenceAck)

	}

	// Delete packet commitment, since the packet has been acknowledged, the commitement is no longer necessary
	DeletePacketCommitment(accessibleState, packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())

	return nil, remainingGas, err
}
