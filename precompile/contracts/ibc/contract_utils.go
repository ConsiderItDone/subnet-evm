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

func calculateKey(path []byte) common.Address {
	return common.BytesToAddress(crypto.Keccak256Hash(path).Bytes())
}

func clientStateKey(clientID string) common.Address {
	return calculateKey(host.FullClientStateKey(clientID))
}

func consensusStateKey(clientID string, height exported.Height) common.Address {
	return calculateKey(host.FullConsensusStateKey(clientID, height))
}

func connectionKey(connectionID string) common.Address {
	return calculateKey(host.ConnectionKey(connectionID))
}

func channelKey(portID, channelID string) common.Address {
	return calculateKey(host.ChannelKey(portID, channelID))
}

func portKey(portID string) common.Address {
	return calculateKey([]byte(host.PortPath(portID)))
}

func channelCapabilityKey(portID, channelID string) common.Address {
	return calculateKey([]byte(host.ChannelCapabilityPath(portID, channelID)))
}

func getPrecompileState(db contract.StateDB, addr common.Address) ([]byte, error) {
	state := db.GetPrecompileState(addr)
	if len(state) == 0 {
		return nil, ErrEmptyState
	}
	return state, nil
}

func setPrecompileState(db contract.StateDB, addr common.Address, obj Marshaler) error {
	state, err := obj.Marshal()
	if err != nil {
		return err
	}
	db.SetPrecompileState(addr, state)
	return nil
}

func addLog(as contract.AccessibleState, name string, args ...any) error {
	topics, data, err := IBCABI.PackEvent(name, args...)
	if err != nil {
		return err
	}
	blockNumber := as.GetBlockContext().Number().Uint64()
	as.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)
	return nil
}

func getClientState(db contract.StateDB, clientId string) (*ibctm.ClientState, error) {
	state, err := getPrecompileState(db, clientStateKey(clientId))
	if err != nil {
		return nil, err
	}
	clientState := new(ibctm.ClientState)
	if err := clientState.Unmarshal(state); err != nil {
		return nil, err
	}
	return clientState, nil
}

func setClientState(db contract.StateDB, clientId string, clientState *ibctm.ClientState) error {
	return setPrecompileState(db, clientStateKey(clientId), clientState)
}

func getConsensusState(db contract.StateDB, clientId string, height exported.Height) (*ibctm.ConsensusState, error) {
	state, err := getPrecompileState(db, consensusStateKey(clientId, height))
	if err != nil {
		return nil, err
	}
	consensusState := new(ibctm.ConsensusState)
	if err := consensusState.Unmarshal(state); err != nil {
		return nil, err
	}
	return consensusState, nil
}

func setConsensusState(db contract.StateDB, clientId string, height exported.Height, consensusState *ibctm.ConsensusState) error {
	return setPrecompileState(db, consensusStateKey(clientId, height), consensusState)
}

func getConnection(db contract.StateDB, connectionID string) (*connectiontypes.ConnectionEnd, error) {
	state, err := getPrecompileState(db, connectionKey(connectionID))
	if err != nil {
		return nil, err
	}
	connection := new(connectiontypes.ConnectionEnd)
	if err := connection.Unmarshal(state); err != nil {
		return nil, err
	}
	return connection, nil
}

func setConnection(db contract.StateDB, connectionID string, conn *connectiontypes.ConnectionEnd) error {
	return setPrecompileState(db, connectionKey(connectionID), conn)
}

func getChannel(db contract.StateDB, portID string, channelID string) (*channeltypes.Channel, error) {
	state, err := getPrecompileState(db, channelKey(portID, channelID))
	if err != nil {
		return nil, err
	}
	channel := new(channeltypes.Channel)
	if err := channel.Unmarshal(state); err != nil {
		return nil, err
	}
	return channel, nil
}

func setChannel(db contract.StateDB, portID string, channelID string, channel *channeltypes.Channel) error {
	return setPrecompileState(db, channelKey(portID, channelID), channel)
}

func getPort(db contract.StateDB, portID string) (common.Address, error) {
	state, err := getPrecompileState(db, portKey(portID))
	if err != nil {
		return common.Address{}, err
	}
	if len(state) != common.AddressLength {
		return common.Address{}, ErrBadLength
	}
	return common.BytesToAddress(state), nil
}

func setPort(db contract.StateDB, portID string, caller common.Address) error {
	if err := host.PortIdentifierValidator(portID); err != nil {
		return err
	}
	db.SetPrecompileState(portKey(portID), caller[:])
	return nil
}

func getCapability(db contract.StateDB, portID, channelID string) (bool, error) {
	state, err := getPrecompileState(db, channelCapabilityKey(portID, channelID))
	if err != nil {
		return false, err
	}
	if !bytes.Equal(state, []byte{1}) {
		return false, ErrBadCapability
	}
	return true, nil
}

func setCapability(db contract.StateDB, portID, channelID string) error {
	exist, err := getCapability(db, portID, channelID)
	if err != nil && err != ErrEmptyState {
		return err
	}
	if exist {
		return ErrAlreadyExist
	}
	db.SetPrecompileState(channelCapabilityKey(portID, channelID), []byte{1})
	return nil
}
