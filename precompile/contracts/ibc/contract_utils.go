package ibc

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/cometbft/cometbft/light"
	tmtypes "github.com/cometbft/cometbft/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrEmptyState    = errors.New("empty precompile state")
	ErrBadLength     = errors.New("bad length of data")
	ErrBadCapability = errors.New("capability had bad data")
	ErrAlreadyExist  = errors.New("already exist")
)

type (
	Marshaler interface {
		Marshal() ([]byte, error)
	}

	callOpts[T any] struct {
		accessibleState contract.AccessibleState
		caller          common.Address
		addr            common.Address
		suppliedGas     uint64
		readOnly        bool
		args            T
	}
)

func CalculateSlot(path []byte) common.Hash {
	return crypto.Keccak256Hash(path)
}

func ClientStateSlot(clientID string) common.Hash {
	return CalculateSlot(host.FullClientStateKey(clientID))
}

func ConsensusStateSlot(clientID string, height exported.Height) common.Hash {
	return CalculateSlot(host.FullConsensusStateKey(clientID, height))
}

func ConsensusStatePreviusSlot(clientID string, height exported.Height) common.Hash {
	return CalculateSlot(host.FullConsensusStateKey(clientID, height))
}

func ConnectionSlot(connectionID string) common.Hash {
	return CalculateSlot(host.ConnectionKey(connectionID))
}

func ChannelSlot(portID, channelID string) common.Hash {
	return CalculateSlot(host.ChannelKey(portID, channelID))
}

func PortSlot(portID string) common.Hash {
	return CalculateSlot([]byte(host.PortPath(portID)))
}

func ChannelCapabilitySlot(portID, channelID string) common.Hash {
	return CalculateSlot([]byte(host.ChannelCapabilityPath(portID, channelID)))
}

func ProcessedTimeSlot(height uint64) common.Hash {
	consensusStateKey := fmt.Sprintf("consensusStates/%d", height)
	return CalculateSlot([]byte(append([]byte(consensusStateKey), []byte("/processedTime")...)))
}

func ProcessedHeightSlot(height uint64) common.Hash {
	consensusStateKey := fmt.Sprintf("consensusStates/%d", height)
	return CalculateSlot([]byte(append([]byte(consensusStateKey), []byte("/processedHeight")...)))
}

func AddLog(as contract.AccessibleState, name string, args ...any) error {
	topics, data, err := IBCABI.PackEvent(name, args...)
	if err != nil {
		return err
	}
	blockNumber := as.GetBlockContext().Number().Uint64()
	as.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)
	return nil
}

func GetClientState(db contract.StateDB, clientId string) (*ibctm.ClientState, error) {
	state, err := GetState(db, ClientStateSlot(clientId))
	if err != nil {
		return nil, err
	}
	clientState := new(ibctm.ClientState)
	if err := clientState.Unmarshal(state); err != nil {
		return nil, err
	}
	return clientState, nil
}

func SetClientState(db contract.StateDB, clientId string, clientState *ibctm.ClientState) error {
	return SetState(db, ClientStateSlot(clientId), clientState)
}

func GetConsensusState(db contract.StateDB, clientId string, height exported.Height) (*ibctm.ConsensusState, error) {
	state, err := GetState(db, ConsensusStateSlot(clientId, height))
	if err != nil {
		return nil, err
	}
	consensusState := new(ibctm.ConsensusState)
	if err := consensusState.Unmarshal(state); err != nil {
		return nil, err
	}
	return consensusState, nil
}

type ConsensusStateStore struct {
	consState     *ibctm.ConsensusState
	previusHeight clienttypes.Height
}

func SetConsensusState(db contract.StateDB, clientId string, height exported.Height, consensusState *ibctm.ConsensusState) error {
	return SetState(db, ConsensusStateSlot(clientId, height), consensusState)
}

func GetPreviusHeight(db contract.StateDB, clientId string, height exported.Height) (*clienttypes.Height, error) {
	state, err := GetState(db, ConsensusStateSlot(clientId, height))
	if err != nil {
		return nil, err
	}
	previusHeight := new(clienttypes.Height)
	if err := previusHeight.Unmarshal(state); err != nil {
		return nil, err
	}
	return previusHeight, nil
}

func SetPreviusHeight(db contract.StateDB, clientId string, height exported.Height, previusHeight *clienttypes.Height) error {
	return SetState(db, ConsensusStatePreviusSlot(clientId, height), previusHeight)
}

func GetNextConsensusState(db contract.StateDB, clientState ibctm.ClientState, height exported.Height) (*ibctm.ConsensusState, error) {
	actualHeight := clientState.LatestHeight
	flag := true
	for flag {
		previusHeight, err := GetPreviusHeight(db, clientState.ChainId, actualHeight)
		if err != nil {
			return nil, err
		}
		if previusHeight.Compare(height) == 0 {
			return GetConsensusState(db, clientState.ChainId, actualHeight)
		} else if previusHeight.Compare(height) == 1 {
			actualHeight = *previusHeight
		} else {
			return nil, fmt.Errorf("Iterator out of range")
		}
	}
	return nil, fmt.Errorf("Iterator out of range")
}

func GetPreviousConsensusState(db contract.StateDB, clientState ibctm.ClientState, height exported.Height) (*ibctm.ConsensusState, error) {
	previusHeight, err := GetPreviusHeight(db, clientState.ChainId, height)
	if err != nil {
		return nil, err
	}
	return GetConsensusState(db, clientState.ChainId, previusHeight)

}

func GetConnection(db contract.StateDB, connectionID string) (*connectiontypes.ConnectionEnd, error) {
	state, err := GetState(db, ConnectionSlot(connectionID))
	if err != nil {
		return nil, err
	}
	connection := new(connectiontypes.ConnectionEnd)
	if err := connection.Unmarshal(state); err != nil {
		return nil, err
	}
	return connection, nil
}

func SetConnection(db contract.StateDB, connectionID string, conn *connectiontypes.ConnectionEnd) error {
	return SetState(db, ConnectionSlot(connectionID), conn)
}

func GetChannel(db contract.StateDB, portID string, channelID string) (*channeltypes.Channel, error) {
	state, err := GetState(db, ChannelSlot(portID, channelID))
	if err != nil {
		return nil, err
	}
	channel := new(channeltypes.Channel)
	if err := channel.Unmarshal(state); err != nil {
		return nil, err
	}
	return channel, nil
}

func SetChannel(db contract.StateDB, portID string, channelID string, channel *channeltypes.Channel) error {
	return SetState(db, ChannelSlot(portID, channelID), channel)
}

func GetPort(db contract.StateDB, portID string) (common.Address, error) {
	state, err := GetState(db, PortSlot(portID))
	if err != nil {
		return common.Address{}, err
	}
	if len(state) != common.AddressLength {
		return common.Address{}, ErrBadLength
	}
	return common.BytesToAddress(state), nil
}

func SetPort(db contract.StateDB, portID string, caller common.Address) error {
	if err := host.PortIdentifierValidator(portID); err != nil {
		return err
	}
	setState(db, ContractAddress, PortSlot(portID), caller[:])
	return nil
}

func GetCapability(db contract.StateDB, portID, channelID string) (bool, error) {
	state, err := GetState(db, ChannelCapabilitySlot(portID, channelID))
	if err != nil {
		return false, err
	}
	if !bytes.Equal(state, []byte{1}) {
		return false, ErrBadCapability
	}
	return true, nil
}

func SetCapability(db contract.StateDB, portID, channelID string) error {
	exist, err := GetCapability(db, portID, channelID)
	if err != nil && err != ErrEmptyState {
		fmt.Println("SetCapability")
		return err
	}
	if exist {
		return ErrAlreadyExist
	}
	setState(db, ContractAddress, ChannelCapabilitySlot(portID, channelID), []byte{1})
	return nil
}

func VerifyClientMessage(
	cs *ibctm.ClientState,
	clientID string,
	accessibleState contract.AccessibleState,
	clientMsg exported.ClientMessage,
) error {
	switch msg := clientMsg.(type) {
	case *ibctm.Header:
		return verifyHeader(cs, accessibleState, msg, clientID)
	case *ibctm.Misbehaviour:
		return verifyMisbehaviour(cs, accessibleState, msg, clientID)
	default:
		return clienttypes.ErrInvalidClientType
	}
}

func verifyHeader(
	cs *ibctm.ClientState,
	accessibleState contract.AccessibleState,
	header *ibctm.Header,
	clientID string,
) error {
	timeInt := accessibleState.GetBlockContext().Timestamp().Int64()
	currentTimestamp := time.Unix(timeInt, 0)

	consState, err := GetConsensusState(accessibleState.GetStateDB(), clientID, header.TrustedHeight)
	if err != nil {
		return fmt.Errorf("can't get consensus state, err: %w", err)
	}

	if err := checkTrustedHeader(header, consState); err != nil {
		return err
	}

	// UpdateClient only accepts updates with a header at the same revision
	// as the trusted consensus state
	if header.GetHeight().GetRevisionNumber() != header.TrustedHeight.RevisionNumber {
		return fmt.Errorf("header height revision %d does not match trusted header revision %d, err: %w", header.GetHeight().GetRevisionNumber(), header.TrustedHeight.RevisionNumber, ibctm.ErrInvalidHeaderHeight)
	}

	tmTrustedValidators, err := tmtypes.ValidatorSetFromProto(header.TrustedValidators)
	if err != nil {
		return fmt.Errorf("trusted validator set in not tendermint validator set type, err: %w", err)
	}

	tmSignedHeader, err := tmtypes.SignedHeaderFromProto(header.SignedHeader)
	if err != nil {
		return fmt.Errorf("signed header in not tendermint signed header type, err: %w", err)
	}

	tmValidatorSet, err := tmtypes.ValidatorSetFromProto(header.ValidatorSet)
	if err != nil {
		return fmt.Errorf("validator set in not tendermint validator set type, err: %w", err)
	}

	// assert header height is newer than consensus state
	if header.GetHeight().LTE(header.TrustedHeight) {
		return fmt.Errorf("header height ≤ consensus state height (%s ≤ %s), err: %w", header.GetHeight(), header.TrustedHeight, clienttypes.ErrInvalidHeader)
	}

	// Construct a trusted header using the fields in consensus state
	// Only Height, Time, and NextValidatorsHash are necessary for verification
	// NOTE: updates must be within the same revision
	trustedHeader := tmtypes.Header{
		ChainID:            cs.GetChainID(),
		Height:             int64(header.TrustedHeight.RevisionHeight),
		Time:               consState.Timestamp,
		NextValidatorsHash: consState.NextValidatorsHash,
	}
	signedHeader := tmtypes.SignedHeader{
		Header: &trustedHeader,
	}

	// Verify next header with the passed-in trustedVals
	// - asserts trusting period not passed
	// - assert header timestamp is not past the trusting period
	// - assert header timestamp is past latest stored consensus state timestamp
	// - assert that a TrustLevel proportion of TrustedValidators signed new Commit
	err = light.Verify(
		&signedHeader,
		tmTrustedValidators, tmSignedHeader, tmValidatorSet,
		cs.TrustingPeriod, currentTimestamp, cs.MaxClockDrift, cs.TrustLevel.ToTendermint(),
	)
	if err != nil {
		return fmt.Errorf("failed to verify header, err: %w", err)
	}

	return nil
}

func verifyMisbehaviour(
	cs *ibctm.ClientState,
	accessibleState contract.AccessibleState,
	misbehaviour *ibctm.Misbehaviour,
	clientID string,
) error {
	// Regardless of the type of misbehaviour, ensure that both headers are valid and would have been accepted by light-client

	// Retrieve trusted consensus states for each Header in misbehaviour
	tmConsensusState1, err := GetConsensusState(accessibleState.GetStateDB(), clientID, misbehaviour.Header1.TrustedHeight)
	if err != nil {
		return fmt.Errorf("can't get consensus state, could not get trusted consensus state from clientStore for Header1 at TrustedHeight: %s, err: %w", misbehaviour.Header1.TrustedHeight, err)
	}

	tmConsensusState2, err := GetConsensusState(accessibleState.GetStateDB(), clientID, misbehaviour.Header2.TrustedHeight)
	if err != nil {
		return fmt.Errorf("can't get consensus state, could not get trusted consensus state from clientStore for Header2 at TrustedHeight: %s, err: %w", misbehaviour.Header2.TrustedHeight, err)
	}

	// Check the validity of the two conflicting headers against their respective
	// trusted consensus states
	// NOTE: header height and commitment root assertions are checked in
	// misbehaviour.ValidateBasic by the client keeper and msg.ValidateBasic
	// by the base application.
	time := time.Unix(accessibleState.GetBlockContext().Timestamp().Int64(), 0)

	if err := checkMisbehaviourHeader(
		cs, tmConsensusState1, misbehaviour.Header1, time,
	); err != nil {
		return fmt.Errorf("verifying Header1 in Misbehaviour failed, err: %w", err)
	}
	if err := checkMisbehaviourHeader(
		cs, tmConsensusState2, misbehaviour.Header2, time,
	); err != nil {
		return fmt.Errorf("verifying Header2 in Misbehaviour failed, err: %w", err)
	}

	return nil
}

// checkMisbehaviourHeader checks that a Header in Misbehaviour is valid misbehaviour given
// a trusted ConsensusState
func checkMisbehaviourHeader(
	clientState *ibctm.ClientState, consState *ibctm.ConsensusState, header *ibctm.Header, currentTimestamp time.Time,
) error {
	tmTrustedValset, err := tmtypes.ValidatorSetFromProto(header.TrustedValidators)
	if err != nil {
		return fmt.Errorf("trusted validator set is not tendermint validator set type, err: %w", err)
	}

	tmCommit, err := tmtypes.CommitFromProto(header.Commit)
	if err != nil {
		return fmt.Errorf("commit is not tendermint commit type, err: %w", err)
	}

	// check the trusted fields for the header against ConsensusState
	if err := checkTrustedHeader(header, consState); err != nil {
		return err
	}

	// assert that the age of the trusted consensus state is not older than the trusting period
	if currentTimestamp.Sub(consState.Timestamp) >= clientState.TrustingPeriod {
		return fmt.Errorf("current timestamp minus the latest consensus state timestamp is greater than or equal to the trusting period (%d >= %d), err: %w", currentTimestamp.Sub(consState.Timestamp), clientState.TrustingPeriod, ibctm.ErrTrustingPeriodExpired)
	}

	chainID := clientState.GetChainID()
	// If chainID is in revision format, then set revision number of chainID with the revision number
	// of the misbehaviour header
	// NOTE: misbehaviour verification is not supported for chains which upgrade to a new chainID without
	// strictly following the chainID revision format
	if clienttypes.IsRevisionFormat(chainID) {
		chainID, _ = clienttypes.SetRevisionNumber(chainID, header.GetHeight().GetRevisionNumber())
	}

	// - ValidatorSet must have TrustLevel similarity with trusted FromValidatorSet
	// - ValidatorSets on both headers are valid given the last trusted ValidatorSet
	if err := tmTrustedValset.VerifyCommitLightTrusting(
		chainID, tmCommit, clientState.TrustLevel.ToTendermint(),
	); err != nil {
		return fmt.Errorf("validator set in header has too much change from trusted validator set: %v: %w", err, clienttypes.ErrInvalidMisbehaviour)
	}
	return nil
}

// checkTrustedHeader checks that consensus state matches trusted fields of Header
func checkTrustedHeader(header *ibctm.Header, consState *ibctm.ConsensusState) error {
	tmTrustedValidators, err := tmtypes.ValidatorSetFromProto(header.TrustedValidators)
	if err != nil {
		return fmt.Errorf("trusted validator set in not tendermint validator set type, err: %w", err)
	}

	// assert that trustedVals is NextValidators of last trusted header
	// to do this, we check that trustedVals.Hash() == consState.NextValidatorsHash
	tvalHash := tmTrustedValidators.Hash()
	if !bytes.Equal(consState.NextValidatorsHash, tvalHash) {
		return fmt.Errorf("trusted validators %s, does not hash to latest trusted validators. Expected: %X, got: %X, err: %w",
			header.TrustedValidators, consState.NextValidatorsHash, tvalHash, ibctm.ErrInvalidValidatorSet)
	}
	return nil
}

func Status(
	accessibleState contract.AccessibleState,
	cs ibctm.ClientState,
	clientId string,
) exported.Status {
	if !cs.FrozenHeight.IsZero() {
		return exported.Frozen
	}

	// get latest consensus state from clientStore to check for expiry
	consState, err := GetConsensusState(accessibleState.GetStateDB(), clientId, cs.GetLatestHeight())
	if err != nil {
		// if the client state does not have an associated consensus state for its latest height
		// then it must be expired
		return exported.Expired
	}

	now := time.Unix(accessibleState.GetBlockContext().Timestamp().Int64(), 0)
	if cs.IsExpired(consState.Timestamp, now) {
		return exported.Expired
	}

	return exported.Active
}
func IsExpired(latestTimestamp time.Time, cs ibctm.ClientState, now time.Time) bool {
	expirationTime := latestTimestamp.Add(cs.TrustingPeriod)
	return !expirationTime.After(now)
}
