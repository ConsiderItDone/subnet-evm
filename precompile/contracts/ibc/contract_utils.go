package ibc

import (
	"bytes"
	"errors"

	"github.com/ava-labs/subnet-evm/precompile/contract"
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

func SetConsensusState(db contract.StateDB, clientId string, height exported.Height, consensusState *ibctm.ConsensusState) error {
	return SetState(db, ConsensusStateSlot(clientId, height), consensusState)
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
		return err
	}
	if exist {
		return ErrAlreadyExist
	}
	setState(db, ContractAddress, ChannelCapabilitySlot(portID, channelID), []byte{1})
	return nil
}
