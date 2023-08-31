package ibc

import (
	"fmt"
	"math"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	hosttypes "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

func VerifyPacketCommitment(
	cdc codec.BinaryCodec,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
	commitmentBytes []byte,
	accessibleState contract.AccessibleState,
) error {
	clientID := connection.GetClientID()

	clientState, err := GetClientState(accessibleState.GetStateDB(), clientID)
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}

	// get time and block delays
	timeDelay := connection.GetDelayPeriod()
	expectedTimePerBlock := 2
	blockDelay := uint64(math.Ceil(float64(timeDelay) / float64(expectedTimePerBlock)))

	merklePath := commitmenttypes.NewMerklePath(hosttypes.PacketCommitmentPath(portID, channelID, sequence))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	if err := VerifyMembership(
		*clientState, cdc, height,
		timeDelay, blockDelay, accessibleState,
		proof, merklePath, commitmentBytes, clientID,
	); err != nil {
		return fmt.Errorf("%w, failed packet commitment verification for client (%s)", err, clientID)
	}
	return nil
}
func VerifyPacketAcknowledgement(
	cdc codec.BinaryCodec,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
	acknowledgement []byte,
	accessibleState contract.AccessibleState,
) error {
	clientID := connection.GetClientID()

	clientState, err := GetClientState(accessibleState.GetStateDB(), clientID)
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}

	// get time and block delays
	timeDelay := connection.GetDelayPeriod()
	expectedTimePerBlock := 2
	blockDelay := uint64(math.Ceil(float64(timeDelay) / float64(expectedTimePerBlock)))

	merklePath := commitmenttypes.NewMerklePath(hosttypes.PacketAcknowledgementPath(portID, channelID, sequence))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	if err := VerifyMembership(
		*clientState, cdc, height,
		timeDelay, blockDelay, accessibleState,
		proof, merklePath, channeltypes.CommitAcknowledgement(acknowledgement), clientID,
	); err != nil {
		return fmt.Errorf("%w, failed packet acknowledgement verification for client (%s)", err, clientID)
	}
	return nil
}

func VerifyChannelState(
	cdc codec.BinaryCodec,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID,
	channelID string,
	channel exported.ChannelI,
	accessibleState contract.AccessibleState,
) error {
	clientID := connection.GetClientID()

	clientState, err := GetClientState(accessibleState.GetStateDB(), clientID)
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}

	merklePath := commitmenttypes.NewMerklePath(hosttypes.ChannelPath(portID, channelID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	channelEnd, ok := channel.(channeltypes.Channel)
	if !ok {
		return fmt.Errorf("invalid channel type %T", channel)
	}

	bz, err := cdc.Marshal(&channelEnd)
	if err != nil {
		return err
	}

	if err := VerifyMembership(
		*clientState, cdc, height,
		0, 0, // skip delay period checks for non-packet processing verification
		accessibleState, proof, merklePath, bz, clientID,
	); err != nil {
		return fmt.Errorf("%w, failed channel state verification for client (%s)", err, clientID)
	}
	return nil
}

func VerifyNextSequenceRecv(
	cdc codec.BinaryCodec,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID,
	channelID string,
	nextSequenceRecv uint64,
	accessibleState contract.AccessibleState,
) error {
	clientID := connection.GetClientID()

	clientState, err := GetClientState(accessibleState.GetStateDB(), clientID)
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}

	// get time and block delays
	timeDelay := connection.GetDelayPeriod()
	expectedTimePerBlock := 2
	blockDelay := uint64(math.Ceil(float64(timeDelay) / float64(expectedTimePerBlock)))

	merklePath := commitmenttypes.NewMerklePath(hosttypes.NextSequenceRecvPath(portID, channelID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	if err := VerifyMembership(
		*clientState, cdc, height,
		timeDelay, blockDelay,
		accessibleState, proof, merklePath, sdk.Uint64ToBigEndian(nextSequenceRecv),
		clientID,
	); err != nil {
		return fmt.Errorf("%w, failed next sequence receive verification for client (%s)", err, clientID)
	}
	return nil
}

func VerifyPacketReceiptAbsence(
	cdc codec.BinaryCodec,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
	accessibleState contract.AccessibleState,
) error {
	clientID := connection.GetClientID()

	clientState, err := GetClientState(accessibleState.GetStateDB(), clientID)
	if err != nil {
		return fmt.Errorf("error loading client state, err: %w", err)
	}

	// get time and block delays
	timeDelay := connection.GetDelayPeriod()
	expectedTimePerBlock := 2
	blockDelay := uint64(math.Ceil(float64(timeDelay) / float64(expectedTimePerBlock)))

	merklePath := commitmenttypes.NewMerklePath(hosttypes.PacketReceiptPath(portID, channelID, sequence))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	if err := VerifyNonMembership(
		*clientState, cdc, height,
		timeDelay, blockDelay,
		accessibleState, proof, merklePath,
	); err != nil {
		return fmt.Errorf("%w, failed packet receipt absence verification for client (%s)", err, clientID)
	}
	return nil
}

func VerifyMembership(
	cs ibctm.ClientState,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	accessibleState contract.AccessibleState,
	proof []byte,
	path exported.Path,
	value []byte,
	clientID string,
) error {

	if cs.GetLatestHeight().LT(height) {
		return fmt.Errorf("client state height < proof height (%d < %d), please ensure the client has been updated", cs.GetLatestHeight(), height)
	}

	if err := verifyDelayPeriodPassed(height, delayTimePeriod, delayBlockPeriod, accessibleState); err != nil {
		return err
	}

	var merkleProof commitmenttypes.MerkleProof
	if err := cdc.Unmarshal(proof, &merkleProof); err != nil {
		return fmt.Errorf("%w, failed to unmarshal proof into ICS 23 commitment merkle proof", commitmenttypes.ErrInvalidProof)
	}

	merklePath, ok := path.(commitmenttypes.MerklePath)
	if !ok {
		return fmt.Errorf(", expected %T, got %T", commitmenttypes.MerklePath{}, path)
	}

	consensusState, err := GetConsensusState(accessibleState.GetStateDB(), clientID, height)
	if err != nil {
		return fmt.Errorf("%w, %w, please ensure the proof was constructed against a height that exists on the client", clienttypes.ErrConsensusStateNotFound, err)
	}

	return merkleProof.VerifyMembership(cs.ProofSpecs, consensusState.GetRoot(), merklePath, value)
}

func VerifyNonMembership(
	cs ibctm.ClientState,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	accessibleState contract.AccessibleState,
	proof []byte,
	path exported.Path,
) error {

	if cs.GetLatestHeight().LT(height) {
		return fmt.Errorf("client state height < proof height (%d < %d), please ensure the client has been updated", cs.GetLatestHeight(), height)
	}

	if err := verifyDelayPeriodPassed(height, delayTimePeriod, delayBlockPeriod, accessibleState); err != nil {
		return err
	}

	var merkleProof commitmenttypes.MerkleProof
	if err := cdc.Unmarshal(proof, &merkleProof); err != nil {
		return fmt.Errorf("%w, failed to unmarshal proof into ICS 23 commitment merkle proof", commitmenttypes.ErrInvalidProof)
	}

	merklePath, ok := path.(commitmenttypes.MerklePath)
	if !ok {
		return fmt.Errorf(", expected %T, got %T", commitmenttypes.MerklePath{}, path)
	}

	consensusState, err := GetConsensusState(accessibleState.GetStateDB(), cs.ChainId, height)
	if err != nil {
		return fmt.Errorf("%w, %w, please ensure the proof was constructed against a height that exists on the client", clienttypes.ErrConsensusStateNotFound, err)
	}

	return merkleProof.VerifyNonMembership(cs.ProofSpecs, consensusState.GetRoot(), merklePath)
}

func verifyDelayPeriodPassed(
	proofHeight exported.Height,
	delayTimePeriod, delayBlockPeriod uint64,
	accessibleState contract.AccessibleState,
) error {
	if delayTimePeriod != 0 {
		// check that executing chain's timestamp has passed consensusState's processed time + delay time period
		processedTime, err := GetProcessedTime(accessibleState.GetStateDB(), proofHeight.GetRevisionHeight())
		if err != nil {
			return fmt.Errorf("%w, processed time not found for height: %s", err, proofHeight)
		}

		currentTimestamp := accessibleState.GetBlockContext().Timestamp().Uint64()
		validTime := processedTime + delayTimePeriod

		// NOTE: delay time period is inclusive, so if currentTimestamp is validTime, then we return no error
		if currentTimestamp < validTime {
			return fmt.Errorf("cannot verify packet until time: %d, current time: %d",
				validTime, currentTimestamp)
		}

	}

	if delayBlockPeriod != 0 {
		// check that executing chain's height has passed consensusState's processed height + delay block period
		processedHeight, err := GetProcessedHeight(accessibleState.GetStateDB(), proofHeight.GetRevisionHeight())
		if err != nil {
			return fmt.Errorf("%w, processed height not found for height: %s", err, proofHeight)
		}

		currentHeight := accessibleState.GetBlockContext().Number().Uint64()
		validHeight := clienttypes.NewHeight(processedHeight.GetRevisionNumber(), processedHeight.GetRevisionHeight()+delayBlockPeriod)

		// NOTE: delay block period is inclusive, so if currentHeight is validHeight, then we return no error
		if currentHeight < validHeight.RevisionHeight {
			return fmt.Errorf("cannot verify packet until height: %s, current height: %d",
				validHeight.String(), currentHeight)
		}
	}

	return nil
}
