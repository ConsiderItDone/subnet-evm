// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = interfaces.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Height is an auto generated low-level Go binding around an user-defined struct.
type Height struct {
	RevisionNumber *big.Int
	RevisionHeight *big.Int
}

// IIBCMsgRecvPacket is an auto generated low-level Go binding around an user-defined struct.
type IIBCMsgRecvPacket struct {
	Packet          Packet
	ProofCommitment []byte
	ProofHeight     Height
	Signer          string
}

// Packet is an auto generated low-level Go binding around an user-defined struct.
type Packet struct {
	Sequence           *big.Int
	SourcePort         string
	SourceChannel      string
	DestinationPort    string
	DestinationChannel string
	Data               []byte
	TimeoutHeight      Height
	TimeoutTimestamp   *big.Int
}

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"timeoutHeight\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destPort\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destChannel\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int32\",\"name\":\"channelOrdering\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"connectionID\",\"type\":\"string\"}],\"name\":\"AcknowledgePacket\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"timeoutHeight\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destPort\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destChannel\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ack\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"connectionID\",\"type\":\"string\"}],\"name\":\"AcknowledgementWritten\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"connectionId\",\"type\":\"string\"}],\"name\":\"ChannelCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"}],\"name\":\"ClientCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"connectionId\",\"type\":\"string\"}],\"name\":\"ConnectionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"timeoutHeight\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destPort\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destChannel\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int32\",\"name\":\"channelOrdering\",\"type\":\"int32\"}],\"name\":\"PacketReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"timeoutHeight\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destPort\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destChannel\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int32\",\"name\":\"channelOrdering\",\"type\":\"int32\"}],\"name\":\"PacketSent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"timeoutHeight\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destPort\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destChannel\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int32\",\"name\":\"channelOrdering\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"connectionID\",\"type\":\"string\"}],\"name\":\"TimeoutPacket\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationPort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationChannel\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"timeoutHeight\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structPacket\",\"name\":\"packet\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"Relayer\",\"type\":\"bytes\"}],\"name\":\"OnRecvPacket\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationPort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationChannel\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"timeoutHeight\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structPacket\",\"name\":\"packet\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"Relayer\",\"type\":\"bytes\"}],\"name\":\"OnTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationPort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationChannel\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"timeoutHeight\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structPacket\",\"name\":\"packet\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"acknowledgement\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofAcked\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"proofHeight\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"signer\",\"type\":\"string\"}],\"name\":\"acknowledgement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"portID\",\"type\":\"string\"}],\"name\":\"bindPort\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"portID\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"channel\",\"type\":\"bytes\"}],\"name\":\"chanOpenInit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"portID\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"channel\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"counterpartyVersion\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"proofInit\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofHeight\",\"type\":\"bytes\"}],\"name\":\"chanOpenTry\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"channelID\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"portID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"channelID\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"proofInit\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofHeight\",\"type\":\"bytes\"}],\"name\":\"channelCloseConfirm\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"portID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"channelID\",\"type\":\"string\"}],\"name\":\"channelCloseInit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"portID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"channelID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"counterpartyChannelID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"counterpartyVersion\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"proofTry\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofHeight\",\"type\":\"bytes\"}],\"name\":\"channelOpenAck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"portID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"channelID\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"proofAck\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofHeight\",\"type\":\"bytes\"}],\"name\":\"channelOpenConfirm\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"connectionID\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"clientState\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"counterpartyConnectionID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofTry\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofClient\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofConsensus\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofHeight\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"consensusHeight\",\"type\":\"bytes\"}],\"name\":\"connOpenAck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"connectionID\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"proofAck\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofHeight\",\"type\":\"bytes\"}],\"name\":\"connOpenConfirm\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"clientID\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"counterparty\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"delayPeriod\",\"type\":\"uint32\"}],\"name\":\"connOpenInit\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"connectionID\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"counterparty\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"delayPeriod\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"clientID\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"clientState\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"counterpartyVersions\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofInit\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofClient\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofConsensus\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofHeight\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"consensusHeight\",\"type\":\"bytes\"}],\"name\":\"connOpenTry\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"connectionID\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"clientType\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"clientState\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"consensusState\",\"type\":\"bytes\"}],\"name\":\"createClient\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"clientID\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationPort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationChannel\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"timeoutHeight\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structPacket\",\"name\":\"packet\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"proofCommitment\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"proofHeight\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"signer\",\"type\":\"string\"}],\"internalType\":\"structIIBC.MsgRecvPacket\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"recvPacket\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"channelCapability\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"timeoutHeight\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"sendPacket\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationPort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationChannel\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"timeoutHeight\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structPacket\",\"name\":\"packet\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"proofUnreceived\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"proofHeight\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"nextSequenceRecv\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"signer\",\"type\":\"string\"}],\"name\":\"timeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationPort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationChannel\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"timeoutHeight\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structPacket\",\"name\":\"packet\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"proofUnreceived\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofClose\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"proofHeight\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"nextSequenceRecv\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"signer\",\"type\":\"string\"}],\"name\":\"timeoutOnClose\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"clientID\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"clientMessage\",\"type\":\"bytes\"}],\"name\":\"updateClient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"clientID\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"upgradePath\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"upgradedClien\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"upgradedConsState\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofUpgradeClient\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofUpgradeConsState\",\"type\":\"bytes\"}],\"name\":\"upgradeClient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// OnRecvPacket is a paid mutator transaction binding the contract method 0x85f7175c.
//
// Solidity: function OnRecvPacket((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes Relayer) returns()
func (_Contract *ContractTransactor) OnRecvPacket(opts *bind.TransactOpts, packet Packet, Relayer []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "OnRecvPacket", packet, Relayer)
}

// OnRecvPacket is a paid mutator transaction binding the contract method 0x85f7175c.
//
// Solidity: function OnRecvPacket((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes Relayer) returns()
func (_Contract *ContractSession) OnRecvPacket(packet Packet, Relayer []byte) (*types.Transaction, error) {
	return _Contract.Contract.OnRecvPacket(&_Contract.TransactOpts, packet, Relayer)
}

// OnRecvPacket is a paid mutator transaction binding the contract method 0x85f7175c.
//
// Solidity: function OnRecvPacket((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes Relayer) returns()
func (_Contract *ContractTransactorSession) OnRecvPacket(packet Packet, Relayer []byte) (*types.Transaction, error) {
	return _Contract.Contract.OnRecvPacket(&_Contract.TransactOpts, packet, Relayer)
}

// OnTimeout is a paid mutator transaction binding the contract method 0x36b8b913.
//
// Solidity: function OnTimeout((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes Relayer) returns()
func (_Contract *ContractTransactor) OnTimeout(opts *bind.TransactOpts, packet Packet, Relayer []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "OnTimeout", packet, Relayer)
}

// OnTimeout is a paid mutator transaction binding the contract method 0x36b8b913.
//
// Solidity: function OnTimeout((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes Relayer) returns()
func (_Contract *ContractSession) OnTimeout(packet Packet, Relayer []byte) (*types.Transaction, error) {
	return _Contract.Contract.OnTimeout(&_Contract.TransactOpts, packet, Relayer)
}

// OnTimeout is a paid mutator transaction binding the contract method 0x36b8b913.
//
// Solidity: function OnTimeout((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes Relayer) returns()
func (_Contract *ContractTransactorSession) OnTimeout(packet Packet, Relayer []byte) (*types.Transaction, error) {
	return _Contract.Contract.OnTimeout(&_Contract.TransactOpts, packet, Relayer)
}

// Acknowledgement is a paid mutator transaction binding the contract method 0xf8831420.
//
// Solidity: function acknowledgement((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes acknowledgement, bytes proofAcked, (uint256,uint256) proofHeight, string signer) returns()
func (_Contract *ContractTransactor) Acknowledgement(opts *bind.TransactOpts, packet Packet, acknowledgement []byte, proofAcked []byte, proofHeight Height, signer string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "acknowledgement", packet, acknowledgement, proofAcked, proofHeight, signer)
}

// Acknowledgement is a paid mutator transaction binding the contract method 0xf8831420.
//
// Solidity: function acknowledgement((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes acknowledgement, bytes proofAcked, (uint256,uint256) proofHeight, string signer) returns()
func (_Contract *ContractSession) Acknowledgement(packet Packet, acknowledgement []byte, proofAcked []byte, proofHeight Height, signer string) (*types.Transaction, error) {
	return _Contract.Contract.Acknowledgement(&_Contract.TransactOpts, packet, acknowledgement, proofAcked, proofHeight, signer)
}

// Acknowledgement is a paid mutator transaction binding the contract method 0xf8831420.
//
// Solidity: function acknowledgement((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes acknowledgement, bytes proofAcked, (uint256,uint256) proofHeight, string signer) returns()
func (_Contract *ContractTransactorSession) Acknowledgement(packet Packet, acknowledgement []byte, proofAcked []byte, proofHeight Height, signer string) (*types.Transaction, error) {
	return _Contract.Contract.Acknowledgement(&_Contract.TransactOpts, packet, acknowledgement, proofAcked, proofHeight, signer)
}

// BindPort is a paid mutator transaction binding the contract method 0xc13b184f.
//
// Solidity: function bindPort(string portID) returns()
func (_Contract *ContractTransactor) BindPort(opts *bind.TransactOpts, portID string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "bindPort", portID)
}

// BindPort is a paid mutator transaction binding the contract method 0xc13b184f.
//
// Solidity: function bindPort(string portID) returns()
func (_Contract *ContractSession) BindPort(portID string) (*types.Transaction, error) {
	return _Contract.Contract.BindPort(&_Contract.TransactOpts, portID)
}

// BindPort is a paid mutator transaction binding the contract method 0xc13b184f.
//
// Solidity: function bindPort(string portID) returns()
func (_Contract *ContractTransactorSession) BindPort(portID string) (*types.Transaction, error) {
	return _Contract.Contract.BindPort(&_Contract.TransactOpts, portID)
}

// ChanOpenInit is a paid mutator transaction binding the contract method 0xa718c941.
//
// Solidity: function chanOpenInit(string portID, bytes channel) returns()
func (_Contract *ContractTransactor) ChanOpenInit(opts *bind.TransactOpts, portID string, channel []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "chanOpenInit", portID, channel)
}

// ChanOpenInit is a paid mutator transaction binding the contract method 0xa718c941.
//
// Solidity: function chanOpenInit(string portID, bytes channel) returns()
func (_Contract *ContractSession) ChanOpenInit(portID string, channel []byte) (*types.Transaction, error) {
	return _Contract.Contract.ChanOpenInit(&_Contract.TransactOpts, portID, channel)
}

// ChanOpenInit is a paid mutator transaction binding the contract method 0xa718c941.
//
// Solidity: function chanOpenInit(string portID, bytes channel) returns()
func (_Contract *ContractTransactorSession) ChanOpenInit(portID string, channel []byte) (*types.Transaction, error) {
	return _Contract.Contract.ChanOpenInit(&_Contract.TransactOpts, portID, channel)
}

// ChanOpenTry is a paid mutator transaction binding the contract method 0x0ce2b1f6.
//
// Solidity: function chanOpenTry(string portID, bytes channel, string counterpartyVersion, bytes proofInit, bytes proofHeight) returns(string channelID)
func (_Contract *ContractTransactor) ChanOpenTry(opts *bind.TransactOpts, portID string, channel []byte, counterpartyVersion string, proofInit []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "chanOpenTry", portID, channel, counterpartyVersion, proofInit, proofHeight)
}

// ChanOpenTry is a paid mutator transaction binding the contract method 0x0ce2b1f6.
//
// Solidity: function chanOpenTry(string portID, bytes channel, string counterpartyVersion, bytes proofInit, bytes proofHeight) returns(string channelID)
func (_Contract *ContractSession) ChanOpenTry(portID string, channel []byte, counterpartyVersion string, proofInit []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ChanOpenTry(&_Contract.TransactOpts, portID, channel, counterpartyVersion, proofInit, proofHeight)
}

// ChanOpenTry is a paid mutator transaction binding the contract method 0x0ce2b1f6.
//
// Solidity: function chanOpenTry(string portID, bytes channel, string counterpartyVersion, bytes proofInit, bytes proofHeight) returns(string channelID)
func (_Contract *ContractTransactorSession) ChanOpenTry(portID string, channel []byte, counterpartyVersion string, proofInit []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ChanOpenTry(&_Contract.TransactOpts, portID, channel, counterpartyVersion, proofInit, proofHeight)
}

// ChannelCloseConfirm is a paid mutator transaction binding the contract method 0x460d68fa.
//
// Solidity: function channelCloseConfirm(string portID, string channelID, bytes proofInit, bytes proofHeight) returns()
func (_Contract *ContractTransactor) ChannelCloseConfirm(opts *bind.TransactOpts, portID string, channelID string, proofInit []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "channelCloseConfirm", portID, channelID, proofInit, proofHeight)
}

// ChannelCloseConfirm is a paid mutator transaction binding the contract method 0x460d68fa.
//
// Solidity: function channelCloseConfirm(string portID, string channelID, bytes proofInit, bytes proofHeight) returns()
func (_Contract *ContractSession) ChannelCloseConfirm(portID string, channelID string, proofInit []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ChannelCloseConfirm(&_Contract.TransactOpts, portID, channelID, proofInit, proofHeight)
}

// ChannelCloseConfirm is a paid mutator transaction binding the contract method 0x460d68fa.
//
// Solidity: function channelCloseConfirm(string portID, string channelID, bytes proofInit, bytes proofHeight) returns()
func (_Contract *ContractTransactorSession) ChannelCloseConfirm(portID string, channelID string, proofInit []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ChannelCloseConfirm(&_Contract.TransactOpts, portID, channelID, proofInit, proofHeight)
}

// ChannelCloseInit is a paid mutator transaction binding the contract method 0x7eb320da.
//
// Solidity: function channelCloseInit(string portID, string channelID) returns()
func (_Contract *ContractTransactor) ChannelCloseInit(opts *bind.TransactOpts, portID string, channelID string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "channelCloseInit", portID, channelID)
}

// ChannelCloseInit is a paid mutator transaction binding the contract method 0x7eb320da.
//
// Solidity: function channelCloseInit(string portID, string channelID) returns()
func (_Contract *ContractSession) ChannelCloseInit(portID string, channelID string) (*types.Transaction, error) {
	return _Contract.Contract.ChannelCloseInit(&_Contract.TransactOpts, portID, channelID)
}

// ChannelCloseInit is a paid mutator transaction binding the contract method 0x7eb320da.
//
// Solidity: function channelCloseInit(string portID, string channelID) returns()
func (_Contract *ContractTransactorSession) ChannelCloseInit(portID string, channelID string) (*types.Transaction, error) {
	return _Contract.Contract.ChannelCloseInit(&_Contract.TransactOpts, portID, channelID)
}

// ChannelOpenAck is a paid mutator transaction binding the contract method 0xbd6f4bde.
//
// Solidity: function channelOpenAck(string portID, string channelID, string counterpartyChannelID, string counterpartyVersion, bytes proofTry, bytes proofHeight) returns()
func (_Contract *ContractTransactor) ChannelOpenAck(opts *bind.TransactOpts, portID string, channelID string, counterpartyChannelID string, counterpartyVersion string, proofTry []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "channelOpenAck", portID, channelID, counterpartyChannelID, counterpartyVersion, proofTry, proofHeight)
}

// ChannelOpenAck is a paid mutator transaction binding the contract method 0xbd6f4bde.
//
// Solidity: function channelOpenAck(string portID, string channelID, string counterpartyChannelID, string counterpartyVersion, bytes proofTry, bytes proofHeight) returns()
func (_Contract *ContractSession) ChannelOpenAck(portID string, channelID string, counterpartyChannelID string, counterpartyVersion string, proofTry []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ChannelOpenAck(&_Contract.TransactOpts, portID, channelID, counterpartyChannelID, counterpartyVersion, proofTry, proofHeight)
}

// ChannelOpenAck is a paid mutator transaction binding the contract method 0xbd6f4bde.
//
// Solidity: function channelOpenAck(string portID, string channelID, string counterpartyChannelID, string counterpartyVersion, bytes proofTry, bytes proofHeight) returns()
func (_Contract *ContractTransactorSession) ChannelOpenAck(portID string, channelID string, counterpartyChannelID string, counterpartyVersion string, proofTry []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ChannelOpenAck(&_Contract.TransactOpts, portID, channelID, counterpartyChannelID, counterpartyVersion, proofTry, proofHeight)
}

// ChannelOpenConfirm is a paid mutator transaction binding the contract method 0x9c110621.
//
// Solidity: function channelOpenConfirm(string portID, string channelID, bytes proofAck, bytes proofHeight) returns()
func (_Contract *ContractTransactor) ChannelOpenConfirm(opts *bind.TransactOpts, portID string, channelID string, proofAck []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "channelOpenConfirm", portID, channelID, proofAck, proofHeight)
}

// ChannelOpenConfirm is a paid mutator transaction binding the contract method 0x9c110621.
//
// Solidity: function channelOpenConfirm(string portID, string channelID, bytes proofAck, bytes proofHeight) returns()
func (_Contract *ContractSession) ChannelOpenConfirm(portID string, channelID string, proofAck []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ChannelOpenConfirm(&_Contract.TransactOpts, portID, channelID, proofAck, proofHeight)
}

// ChannelOpenConfirm is a paid mutator transaction binding the contract method 0x9c110621.
//
// Solidity: function channelOpenConfirm(string portID, string channelID, bytes proofAck, bytes proofHeight) returns()
func (_Contract *ContractTransactorSession) ChannelOpenConfirm(portID string, channelID string, proofAck []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ChannelOpenConfirm(&_Contract.TransactOpts, portID, channelID, proofAck, proofHeight)
}

// ConnOpenAck is a paid mutator transaction binding the contract method 0x528d5cd3.
//
// Solidity: function connOpenAck(string connectionID, bytes clientState, bytes version, bytes counterpartyConnectionID, bytes proofTry, bytes proofClient, bytes proofConsensus, bytes proofHeight, bytes consensusHeight) returns()
func (_Contract *ContractTransactor) ConnOpenAck(opts *bind.TransactOpts, connectionID string, clientState []byte, version []byte, counterpartyConnectionID []byte, proofTry []byte, proofClient []byte, proofConsensus []byte, proofHeight []byte, consensusHeight []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "connOpenAck", connectionID, clientState, version, counterpartyConnectionID, proofTry, proofClient, proofConsensus, proofHeight, consensusHeight)
}

// ConnOpenAck is a paid mutator transaction binding the contract method 0x528d5cd3.
//
// Solidity: function connOpenAck(string connectionID, bytes clientState, bytes version, bytes counterpartyConnectionID, bytes proofTry, bytes proofClient, bytes proofConsensus, bytes proofHeight, bytes consensusHeight) returns()
func (_Contract *ContractSession) ConnOpenAck(connectionID string, clientState []byte, version []byte, counterpartyConnectionID []byte, proofTry []byte, proofClient []byte, proofConsensus []byte, proofHeight []byte, consensusHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ConnOpenAck(&_Contract.TransactOpts, connectionID, clientState, version, counterpartyConnectionID, proofTry, proofClient, proofConsensus, proofHeight, consensusHeight)
}

// ConnOpenAck is a paid mutator transaction binding the contract method 0x528d5cd3.
//
// Solidity: function connOpenAck(string connectionID, bytes clientState, bytes version, bytes counterpartyConnectionID, bytes proofTry, bytes proofClient, bytes proofConsensus, bytes proofHeight, bytes consensusHeight) returns()
func (_Contract *ContractTransactorSession) ConnOpenAck(connectionID string, clientState []byte, version []byte, counterpartyConnectionID []byte, proofTry []byte, proofClient []byte, proofConsensus []byte, proofHeight []byte, consensusHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ConnOpenAck(&_Contract.TransactOpts, connectionID, clientState, version, counterpartyConnectionID, proofTry, proofClient, proofConsensus, proofHeight, consensusHeight)
}

// ConnOpenConfirm is a paid mutator transaction binding the contract method 0x45870d5e.
//
// Solidity: function connOpenConfirm(string connectionID, bytes proofAck, bytes proofHeight) returns()
func (_Contract *ContractTransactor) ConnOpenConfirm(opts *bind.TransactOpts, connectionID string, proofAck []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "connOpenConfirm", connectionID, proofAck, proofHeight)
}

// ConnOpenConfirm is a paid mutator transaction binding the contract method 0x45870d5e.
//
// Solidity: function connOpenConfirm(string connectionID, bytes proofAck, bytes proofHeight) returns()
func (_Contract *ContractSession) ConnOpenConfirm(connectionID string, proofAck []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ConnOpenConfirm(&_Contract.TransactOpts, connectionID, proofAck, proofHeight)
}

// ConnOpenConfirm is a paid mutator transaction binding the contract method 0x45870d5e.
//
// Solidity: function connOpenConfirm(string connectionID, bytes proofAck, bytes proofHeight) returns()
func (_Contract *ContractTransactorSession) ConnOpenConfirm(connectionID string, proofAck []byte, proofHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ConnOpenConfirm(&_Contract.TransactOpts, connectionID, proofAck, proofHeight)
}

// ConnOpenInit is a paid mutator transaction binding the contract method 0xd198062b.
//
// Solidity: function connOpenInit(string clientID, bytes counterparty, bytes version, uint32 delayPeriod) returns(string connectionID)
func (_Contract *ContractTransactor) ConnOpenInit(opts *bind.TransactOpts, clientID string, counterparty []byte, version []byte, delayPeriod uint32) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "connOpenInit", clientID, counterparty, version, delayPeriod)
}

// ConnOpenInit is a paid mutator transaction binding the contract method 0xd198062b.
//
// Solidity: function connOpenInit(string clientID, bytes counterparty, bytes version, uint32 delayPeriod) returns(string connectionID)
func (_Contract *ContractSession) ConnOpenInit(clientID string, counterparty []byte, version []byte, delayPeriod uint32) (*types.Transaction, error) {
	return _Contract.Contract.ConnOpenInit(&_Contract.TransactOpts, clientID, counterparty, version, delayPeriod)
}

// ConnOpenInit is a paid mutator transaction binding the contract method 0xd198062b.
//
// Solidity: function connOpenInit(string clientID, bytes counterparty, bytes version, uint32 delayPeriod) returns(string connectionID)
func (_Contract *ContractTransactorSession) ConnOpenInit(clientID string, counterparty []byte, version []byte, delayPeriod uint32) (*types.Transaction, error) {
	return _Contract.Contract.ConnOpenInit(&_Contract.TransactOpts, clientID, counterparty, version, delayPeriod)
}

// ConnOpenTry is a paid mutator transaction binding the contract method 0x6954535e.
//
// Solidity: function connOpenTry(bytes counterparty, uint32 delayPeriod, string clientID, bytes clientState, bytes counterpartyVersions, bytes proofInit, bytes proofClient, bytes proofConsensus, bytes proofHeight, bytes consensusHeight) returns(string connectionID)
func (_Contract *ContractTransactor) ConnOpenTry(opts *bind.TransactOpts, counterparty []byte, delayPeriod uint32, clientID string, clientState []byte, counterpartyVersions []byte, proofInit []byte, proofClient []byte, proofConsensus []byte, proofHeight []byte, consensusHeight []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "connOpenTry", counterparty, delayPeriod, clientID, clientState, counterpartyVersions, proofInit, proofClient, proofConsensus, proofHeight, consensusHeight)
}

// ConnOpenTry is a paid mutator transaction binding the contract method 0x6954535e.
//
// Solidity: function connOpenTry(bytes counterparty, uint32 delayPeriod, string clientID, bytes clientState, bytes counterpartyVersions, bytes proofInit, bytes proofClient, bytes proofConsensus, bytes proofHeight, bytes consensusHeight) returns(string connectionID)
func (_Contract *ContractSession) ConnOpenTry(counterparty []byte, delayPeriod uint32, clientID string, clientState []byte, counterpartyVersions []byte, proofInit []byte, proofClient []byte, proofConsensus []byte, proofHeight []byte, consensusHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ConnOpenTry(&_Contract.TransactOpts, counterparty, delayPeriod, clientID, clientState, counterpartyVersions, proofInit, proofClient, proofConsensus, proofHeight, consensusHeight)
}

// ConnOpenTry is a paid mutator transaction binding the contract method 0x6954535e.
//
// Solidity: function connOpenTry(bytes counterparty, uint32 delayPeriod, string clientID, bytes clientState, bytes counterpartyVersions, bytes proofInit, bytes proofClient, bytes proofConsensus, bytes proofHeight, bytes consensusHeight) returns(string connectionID)
func (_Contract *ContractTransactorSession) ConnOpenTry(counterparty []byte, delayPeriod uint32, clientID string, clientState []byte, counterpartyVersions []byte, proofInit []byte, proofClient []byte, proofConsensus []byte, proofHeight []byte, consensusHeight []byte) (*types.Transaction, error) {
	return _Contract.Contract.ConnOpenTry(&_Contract.TransactOpts, counterparty, delayPeriod, clientID, clientState, counterpartyVersions, proofInit, proofClient, proofConsensus, proofHeight, consensusHeight)
}

// CreateClient is a paid mutator transaction binding the contract method 0x2629636b.
//
// Solidity: function createClient(string clientType, bytes clientState, bytes consensusState) returns(string clientID)
func (_Contract *ContractTransactor) CreateClient(opts *bind.TransactOpts, clientType string, clientState []byte, consensusState []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "createClient", clientType, clientState, consensusState)
}

// CreateClient is a paid mutator transaction binding the contract method 0x2629636b.
//
// Solidity: function createClient(string clientType, bytes clientState, bytes consensusState) returns(string clientID)
func (_Contract *ContractSession) CreateClient(clientType string, clientState []byte, consensusState []byte) (*types.Transaction, error) {
	return _Contract.Contract.CreateClient(&_Contract.TransactOpts, clientType, clientState, consensusState)
}

// CreateClient is a paid mutator transaction binding the contract method 0x2629636b.
//
// Solidity: function createClient(string clientType, bytes clientState, bytes consensusState) returns(string clientID)
func (_Contract *ContractTransactorSession) CreateClient(clientType string, clientState []byte, consensusState []byte) (*types.Transaction, error) {
	return _Contract.Contract.CreateClient(&_Contract.TransactOpts, clientType, clientState, consensusState)
}

// RecvPacket is a paid mutator transaction binding the contract method 0x15d1edac.
//
// Solidity: function recvPacket(((uint256,string,string,string,string,bytes,(uint256,uint256),uint256),bytes,(uint256,uint256),string) message) returns()
func (_Contract *ContractTransactor) RecvPacket(opts *bind.TransactOpts, message IIBCMsgRecvPacket) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "recvPacket", message)
}

// RecvPacket is a paid mutator transaction binding the contract method 0x15d1edac.
//
// Solidity: function recvPacket(((uint256,string,string,string,string,bytes,(uint256,uint256),uint256),bytes,(uint256,uint256),string) message) returns()
func (_Contract *ContractSession) RecvPacket(message IIBCMsgRecvPacket) (*types.Transaction, error) {
	return _Contract.Contract.RecvPacket(&_Contract.TransactOpts, message)
}

// RecvPacket is a paid mutator transaction binding the contract method 0x15d1edac.
//
// Solidity: function recvPacket(((uint256,string,string,string,string,bytes,(uint256,uint256),uint256),bytes,(uint256,uint256),string) message) returns()
func (_Contract *ContractTransactorSession) RecvPacket(message IIBCMsgRecvPacket) (*types.Transaction, error) {
	return _Contract.Contract.RecvPacket(&_Contract.TransactOpts, message)
}

// SendPacket is a paid mutator transaction binding the contract method 0x6052bf6f.
//
// Solidity: function sendPacket(uint256 channelCapability, string sourcePort, string sourceChannel, (uint256,uint256) timeoutHeight, uint256 timeoutTimestamp, bytes data) returns()
func (_Contract *ContractTransactor) SendPacket(opts *bind.TransactOpts, channelCapability *big.Int, sourcePort string, sourceChannel string, timeoutHeight Height, timeoutTimestamp *big.Int, data []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "sendPacket", channelCapability, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

// SendPacket is a paid mutator transaction binding the contract method 0x6052bf6f.
//
// Solidity: function sendPacket(uint256 channelCapability, string sourcePort, string sourceChannel, (uint256,uint256) timeoutHeight, uint256 timeoutTimestamp, bytes data) returns()
func (_Contract *ContractSession) SendPacket(channelCapability *big.Int, sourcePort string, sourceChannel string, timeoutHeight Height, timeoutTimestamp *big.Int, data []byte) (*types.Transaction, error) {
	return _Contract.Contract.SendPacket(&_Contract.TransactOpts, channelCapability, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

// SendPacket is a paid mutator transaction binding the contract method 0x6052bf6f.
//
// Solidity: function sendPacket(uint256 channelCapability, string sourcePort, string sourceChannel, (uint256,uint256) timeoutHeight, uint256 timeoutTimestamp, bytes data) returns()
func (_Contract *ContractTransactorSession) SendPacket(channelCapability *big.Int, sourcePort string, sourceChannel string, timeoutHeight Height, timeoutTimestamp *big.Int, data []byte) (*types.Transaction, error) {
	return _Contract.Contract.SendPacket(&_Contract.TransactOpts, channelCapability, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

// Timeout is a paid mutator transaction binding the contract method 0x8883ff39.
//
// Solidity: function timeout((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes proofUnreceived, (uint256,uint256) proofHeight, uint256 nextSequenceRecv, string signer) returns()
func (_Contract *ContractTransactor) Timeout(opts *bind.TransactOpts, packet Packet, proofUnreceived []byte, proofHeight Height, nextSequenceRecv *big.Int, signer string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "timeout", packet, proofUnreceived, proofHeight, nextSequenceRecv, signer)
}

// Timeout is a paid mutator transaction binding the contract method 0x8883ff39.
//
// Solidity: function timeout((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes proofUnreceived, (uint256,uint256) proofHeight, uint256 nextSequenceRecv, string signer) returns()
func (_Contract *ContractSession) Timeout(packet Packet, proofUnreceived []byte, proofHeight Height, nextSequenceRecv *big.Int, signer string) (*types.Transaction, error) {
	return _Contract.Contract.Timeout(&_Contract.TransactOpts, packet, proofUnreceived, proofHeight, nextSequenceRecv, signer)
}

// Timeout is a paid mutator transaction binding the contract method 0x8883ff39.
//
// Solidity: function timeout((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes proofUnreceived, (uint256,uint256) proofHeight, uint256 nextSequenceRecv, string signer) returns()
func (_Contract *ContractTransactorSession) Timeout(packet Packet, proofUnreceived []byte, proofHeight Height, nextSequenceRecv *big.Int, signer string) (*types.Transaction, error) {
	return _Contract.Contract.Timeout(&_Contract.TransactOpts, packet, proofUnreceived, proofHeight, nextSequenceRecv, signer)
}

// TimeoutOnClose is a paid mutator transaction binding the contract method 0xc519baa9.
//
// Solidity: function timeoutOnClose((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes proofUnreceived, bytes proofClose, (uint256,uint256) proofHeight, uint256 nextSequenceRecv, string signer) returns()
func (_Contract *ContractTransactor) TimeoutOnClose(opts *bind.TransactOpts, packet Packet, proofUnreceived []byte, proofClose []byte, proofHeight Height, nextSequenceRecv *big.Int, signer string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "timeoutOnClose", packet, proofUnreceived, proofClose, proofHeight, nextSequenceRecv, signer)
}

// TimeoutOnClose is a paid mutator transaction binding the contract method 0xc519baa9.
//
// Solidity: function timeoutOnClose((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes proofUnreceived, bytes proofClose, (uint256,uint256) proofHeight, uint256 nextSequenceRecv, string signer) returns()
func (_Contract *ContractSession) TimeoutOnClose(packet Packet, proofUnreceived []byte, proofClose []byte, proofHeight Height, nextSequenceRecv *big.Int, signer string) (*types.Transaction, error) {
	return _Contract.Contract.TimeoutOnClose(&_Contract.TransactOpts, packet, proofUnreceived, proofClose, proofHeight, nextSequenceRecv, signer)
}

// TimeoutOnClose is a paid mutator transaction binding the contract method 0xc519baa9.
//
// Solidity: function timeoutOnClose((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes proofUnreceived, bytes proofClose, (uint256,uint256) proofHeight, uint256 nextSequenceRecv, string signer) returns()
func (_Contract *ContractTransactorSession) TimeoutOnClose(packet Packet, proofUnreceived []byte, proofClose []byte, proofHeight Height, nextSequenceRecv *big.Int, signer string) (*types.Transaction, error) {
	return _Contract.Contract.TimeoutOnClose(&_Contract.TransactOpts, packet, proofUnreceived, proofClose, proofHeight, nextSequenceRecv, signer)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x6fbf8079.
//
// Solidity: function updateClient(string clientID, bytes clientMessage) returns()
func (_Contract *ContractTransactor) UpdateClient(opts *bind.TransactOpts, clientID string, clientMessage []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateClient", clientID, clientMessage)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x6fbf8079.
//
// Solidity: function updateClient(string clientID, bytes clientMessage) returns()
func (_Contract *ContractSession) UpdateClient(clientID string, clientMessage []byte) (*types.Transaction, error) {
	return _Contract.Contract.UpdateClient(&_Contract.TransactOpts, clientID, clientMessage)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x6fbf8079.
//
// Solidity: function updateClient(string clientID, bytes clientMessage) returns()
func (_Contract *ContractTransactorSession) UpdateClient(clientID string, clientMessage []byte) (*types.Transaction, error) {
	return _Contract.Contract.UpdateClient(&_Contract.TransactOpts, clientID, clientMessage)
}

// UpgradeClient is a paid mutator transaction binding the contract method 0xe5fd74d3.
//
// Solidity: function upgradeClient(string clientID, bytes upgradePath, bytes upgradedClien, bytes upgradedConsState, bytes proofUpgradeClient, bytes proofUpgradeConsState) returns()
func (_Contract *ContractTransactor) UpgradeClient(opts *bind.TransactOpts, clientID string, upgradePath []byte, upgradedClien []byte, upgradedConsState []byte, proofUpgradeClient []byte, proofUpgradeConsState []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "upgradeClient", clientID, upgradePath, upgradedClien, upgradedConsState, proofUpgradeClient, proofUpgradeConsState)
}

// UpgradeClient is a paid mutator transaction binding the contract method 0xe5fd74d3.
//
// Solidity: function upgradeClient(string clientID, bytes upgradePath, bytes upgradedClien, bytes upgradedConsState, bytes proofUpgradeClient, bytes proofUpgradeConsState) returns()
func (_Contract *ContractSession) UpgradeClient(clientID string, upgradePath []byte, upgradedClien []byte, upgradedConsState []byte, proofUpgradeClient []byte, proofUpgradeConsState []byte) (*types.Transaction, error) {
	return _Contract.Contract.UpgradeClient(&_Contract.TransactOpts, clientID, upgradePath, upgradedClien, upgradedConsState, proofUpgradeClient, proofUpgradeConsState)
}

// UpgradeClient is a paid mutator transaction binding the contract method 0xe5fd74d3.
//
// Solidity: function upgradeClient(string clientID, bytes upgradePath, bytes upgradedClien, bytes upgradedConsState, bytes proofUpgradeClient, bytes proofUpgradeConsState) returns()
func (_Contract *ContractTransactorSession) UpgradeClient(clientID string, upgradePath []byte, upgradedClien []byte, upgradedConsState []byte, proofUpgradeClient []byte, proofUpgradeConsState []byte) (*types.Transaction, error) {
	return _Contract.Contract.UpgradeClient(&_Contract.TransactOpts, clientID, upgradePath, upgradedClien, upgradedConsState, proofUpgradeClient, proofUpgradeConsState)
}

// ContractAcknowledgePacketIterator is returned from FilterAcknowledgePacket and is used to iterate over the raw logs and unpacked data for AcknowledgePacket events raised by the Contract contract.
type ContractAcknowledgePacketIterator struct {
	Event *ContractAcknowledgePacket // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAcknowledgePacketIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAcknowledgePacket)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAcknowledgePacket)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAcknowledgePacketIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAcknowledgePacketIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAcknowledgePacket represents a AcknowledgePacket event raised by the Contract contract.
type ContractAcknowledgePacket struct {
	TimeoutHeight    string
	TimeoutTimestamp *big.Int
	Sequence         *big.Int
	SourcePort       string
	SourceChannel    string
	DestPort         string
	DestChannel      string
	ChannelOrdering  int32
	ConnectionID     string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterAcknowledgePacket is a free log retrieval operation binding the contract event 0x643d36ddde0bd3af37ec1d67f146b0f353d1f5b01eaa8a3879d3890ab9cc224d.
//
// Solidity: event AcknowledgePacket(string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering, string connectionID)
func (_Contract *ContractFilterer) FilterAcknowledgePacket(opts *bind.FilterOpts) (*ContractAcknowledgePacketIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "AcknowledgePacket")
	if err != nil {
		return nil, err
	}
	return &ContractAcknowledgePacketIterator{contract: _Contract.contract, event: "AcknowledgePacket", logs: logs, sub: sub}, nil
}

// WatchAcknowledgePacket is a free log subscription operation binding the contract event 0x643d36ddde0bd3af37ec1d67f146b0f353d1f5b01eaa8a3879d3890ab9cc224d.
//
// Solidity: event AcknowledgePacket(string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering, string connectionID)
func (_Contract *ContractFilterer) WatchAcknowledgePacket(opts *bind.WatchOpts, sink chan<- *ContractAcknowledgePacket) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "AcknowledgePacket")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAcknowledgePacket)
				if err := _Contract.contract.UnpackLog(event, "AcknowledgePacket", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAcknowledgePacket is a log parse operation binding the contract event 0x643d36ddde0bd3af37ec1d67f146b0f353d1f5b01eaa8a3879d3890ab9cc224d.
//
// Solidity: event AcknowledgePacket(string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering, string connectionID)
func (_Contract *ContractFilterer) ParseAcknowledgePacket(log types.Log) (*ContractAcknowledgePacket, error) {
	event := new(ContractAcknowledgePacket)
	if err := _Contract.contract.UnpackLog(event, "AcknowledgePacket", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAcknowledgementWrittenIterator is returned from FilterAcknowledgementWritten and is used to iterate over the raw logs and unpacked data for AcknowledgementWritten events raised by the Contract contract.
type ContractAcknowledgementWrittenIterator struct {
	Event *ContractAcknowledgementWritten // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAcknowledgementWrittenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAcknowledgementWritten)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAcknowledgementWritten)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAcknowledgementWrittenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAcknowledgementWrittenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAcknowledgementWritten represents a AcknowledgementWritten event raised by the Contract contract.
type ContractAcknowledgementWritten struct {
	Data             []byte
	TimeoutHeight    string
	TimeoutTimestamp *big.Int
	Sequence         *big.Int
	SourcePort       string
	SourceChannel    string
	DestPort         string
	DestChannel      string
	Ack              []byte
	ConnectionID     string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterAcknowledgementWritten is a free log retrieval operation binding the contract event 0x24fcffd9284a3995cbe257809fac5494cecf20a31833139c768aacadc0def5cd.
//
// Solidity: event AcknowledgementWritten(bytes data, string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, bytes ack, string connectionID)
func (_Contract *ContractFilterer) FilterAcknowledgementWritten(opts *bind.FilterOpts) (*ContractAcknowledgementWrittenIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "AcknowledgementWritten")
	if err != nil {
		return nil, err
	}
	return &ContractAcknowledgementWrittenIterator{contract: _Contract.contract, event: "AcknowledgementWritten", logs: logs, sub: sub}, nil
}

// WatchAcknowledgementWritten is a free log subscription operation binding the contract event 0x24fcffd9284a3995cbe257809fac5494cecf20a31833139c768aacadc0def5cd.
//
// Solidity: event AcknowledgementWritten(bytes data, string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, bytes ack, string connectionID)
func (_Contract *ContractFilterer) WatchAcknowledgementWritten(opts *bind.WatchOpts, sink chan<- *ContractAcknowledgementWritten) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "AcknowledgementWritten")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAcknowledgementWritten)
				if err := _Contract.contract.UnpackLog(event, "AcknowledgementWritten", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAcknowledgementWritten is a log parse operation binding the contract event 0x24fcffd9284a3995cbe257809fac5494cecf20a31833139c768aacadc0def5cd.
//
// Solidity: event AcknowledgementWritten(bytes data, string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, bytes ack, string connectionID)
func (_Contract *ContractFilterer) ParseAcknowledgementWritten(log types.Log) (*ContractAcknowledgementWritten, error) {
	event := new(ContractAcknowledgementWritten)
	if err := _Contract.contract.UnpackLog(event, "AcknowledgementWritten", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractChannelCreatedIterator is returned from FilterChannelCreated and is used to iterate over the raw logs and unpacked data for ChannelCreated events raised by the Contract contract.
type ContractChannelCreatedIterator struct {
	Event *ContractChannelCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractChannelCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractChannelCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractChannelCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractChannelCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractChannelCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractChannelCreated represents a ChannelCreated event raised by the Contract contract.
type ContractChannelCreated struct {
	ClientId     string
	ConnectionId string
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterChannelCreated is a free log retrieval operation binding the contract event 0xb403db35509133349144332a6f2aece4ef7d5989aee88f9e4a62b8f24fd57b46.
//
// Solidity: event ChannelCreated(string clientId, string connectionId)
func (_Contract *ContractFilterer) FilterChannelCreated(opts *bind.FilterOpts) (*ContractChannelCreatedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ChannelCreated")
	if err != nil {
		return nil, err
	}
	return &ContractChannelCreatedIterator{contract: _Contract.contract, event: "ChannelCreated", logs: logs, sub: sub}, nil
}

// WatchChannelCreated is a free log subscription operation binding the contract event 0xb403db35509133349144332a6f2aece4ef7d5989aee88f9e4a62b8f24fd57b46.
//
// Solidity: event ChannelCreated(string clientId, string connectionId)
func (_Contract *ContractFilterer) WatchChannelCreated(opts *bind.WatchOpts, sink chan<- *ContractChannelCreated) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ChannelCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractChannelCreated)
				if err := _Contract.contract.UnpackLog(event, "ChannelCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelCreated is a log parse operation binding the contract event 0xb403db35509133349144332a6f2aece4ef7d5989aee88f9e4a62b8f24fd57b46.
//
// Solidity: event ChannelCreated(string clientId, string connectionId)
func (_Contract *ContractFilterer) ParseChannelCreated(log types.Log) (*ContractChannelCreated, error) {
	event := new(ContractChannelCreated)
	if err := _Contract.contract.UnpackLog(event, "ChannelCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractClientCreatedIterator is returned from FilterClientCreated and is used to iterate over the raw logs and unpacked data for ClientCreated events raised by the Contract contract.
type ContractClientCreatedIterator struct {
	Event *ContractClientCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractClientCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractClientCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractClientCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractClientCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractClientCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractClientCreated represents a ClientCreated event raised by the Contract contract.
type ContractClientCreated struct {
	ClientId string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterClientCreated is a free log retrieval operation binding the contract event 0xeb98df470d17266538e4ee034952206621fad8d86ca38b090e92f64589108482.
//
// Solidity: event ClientCreated(string clientId)
func (_Contract *ContractFilterer) FilterClientCreated(opts *bind.FilterOpts) (*ContractClientCreatedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ClientCreated")
	if err != nil {
		return nil, err
	}
	return &ContractClientCreatedIterator{contract: _Contract.contract, event: "ClientCreated", logs: logs, sub: sub}, nil
}

// WatchClientCreated is a free log subscription operation binding the contract event 0xeb98df470d17266538e4ee034952206621fad8d86ca38b090e92f64589108482.
//
// Solidity: event ClientCreated(string clientId)
func (_Contract *ContractFilterer) WatchClientCreated(opts *bind.WatchOpts, sink chan<- *ContractClientCreated) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ClientCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractClientCreated)
				if err := _Contract.contract.UnpackLog(event, "ClientCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClientCreated is a log parse operation binding the contract event 0xeb98df470d17266538e4ee034952206621fad8d86ca38b090e92f64589108482.
//
// Solidity: event ClientCreated(string clientId)
func (_Contract *ContractFilterer) ParseClientCreated(log types.Log) (*ContractClientCreated, error) {
	event := new(ContractClientCreated)
	if err := _Contract.contract.UnpackLog(event, "ClientCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractConnectionCreatedIterator is returned from FilterConnectionCreated and is used to iterate over the raw logs and unpacked data for ConnectionCreated events raised by the Contract contract.
type ContractConnectionCreatedIterator struct {
	Event *ContractConnectionCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractConnectionCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractConnectionCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractConnectionCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractConnectionCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractConnectionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractConnectionCreated represents a ConnectionCreated event raised by the Contract contract.
type ContractConnectionCreated struct {
	ClientId     string
	ConnectionId string
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterConnectionCreated is a free log retrieval operation binding the contract event 0x4d26138feedd7dd05ad1f2cb1a4bef7b74ac9612fbadc01ab5488c46b91e7e77.
//
// Solidity: event ConnectionCreated(string clientId, string connectionId)
func (_Contract *ContractFilterer) FilterConnectionCreated(opts *bind.FilterOpts) (*ContractConnectionCreatedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ConnectionCreated")
	if err != nil {
		return nil, err
	}
	return &ContractConnectionCreatedIterator{contract: _Contract.contract, event: "ConnectionCreated", logs: logs, sub: sub}, nil
}

// WatchConnectionCreated is a free log subscription operation binding the contract event 0x4d26138feedd7dd05ad1f2cb1a4bef7b74ac9612fbadc01ab5488c46b91e7e77.
//
// Solidity: event ConnectionCreated(string clientId, string connectionId)
func (_Contract *ContractFilterer) WatchConnectionCreated(opts *bind.WatchOpts, sink chan<- *ContractConnectionCreated) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ConnectionCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractConnectionCreated)
				if err := _Contract.contract.UnpackLog(event, "ConnectionCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConnectionCreated is a log parse operation binding the contract event 0x4d26138feedd7dd05ad1f2cb1a4bef7b74ac9612fbadc01ab5488c46b91e7e77.
//
// Solidity: event ConnectionCreated(string clientId, string connectionId)
func (_Contract *ContractFilterer) ParseConnectionCreated(log types.Log) (*ContractConnectionCreated, error) {
	event := new(ContractConnectionCreated)
	if err := _Contract.contract.UnpackLog(event, "ConnectionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPacketReceivedIterator is returned from FilterPacketReceived and is used to iterate over the raw logs and unpacked data for PacketReceived events raised by the Contract contract.
type ContractPacketReceivedIterator struct {
	Event *ContractPacketReceived // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractPacketReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPacketReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractPacketReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractPacketReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPacketReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPacketReceived represents a PacketReceived event raised by the Contract contract.
type ContractPacketReceived struct {
	Data             []byte
	TimeoutHeight    string
	TimeoutTimestamp *big.Int
	Sequence         *big.Int
	SourcePort       string
	SourceChannel    string
	DestPort         string
	DestChannel      string
	ChannelOrdering  int32
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPacketReceived is a free log retrieval operation binding the contract event 0xa4e7a49b834c9f544209bb6332320ed9fe1587c227f023f0b9facc7a3cccff40.
//
// Solidity: event PacketReceived(bytes data, string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering)
func (_Contract *ContractFilterer) FilterPacketReceived(opts *bind.FilterOpts) (*ContractPacketReceivedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "PacketReceived")
	if err != nil {
		return nil, err
	}
	return &ContractPacketReceivedIterator{contract: _Contract.contract, event: "PacketReceived", logs: logs, sub: sub}, nil
}

// WatchPacketReceived is a free log subscription operation binding the contract event 0xa4e7a49b834c9f544209bb6332320ed9fe1587c227f023f0b9facc7a3cccff40.
//
// Solidity: event PacketReceived(bytes data, string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering)
func (_Contract *ContractFilterer) WatchPacketReceived(opts *bind.WatchOpts, sink chan<- *ContractPacketReceived) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "PacketReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPacketReceived)
				if err := _Contract.contract.UnpackLog(event, "PacketReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePacketReceived is a log parse operation binding the contract event 0xa4e7a49b834c9f544209bb6332320ed9fe1587c227f023f0b9facc7a3cccff40.
//
// Solidity: event PacketReceived(bytes data, string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering)
func (_Contract *ContractFilterer) ParsePacketReceived(log types.Log) (*ContractPacketReceived, error) {
	event := new(ContractPacketReceived)
	if err := _Contract.contract.UnpackLog(event, "PacketReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPacketSentIterator is returned from FilterPacketSent and is used to iterate over the raw logs and unpacked data for PacketSent events raised by the Contract contract.
type ContractPacketSentIterator struct {
	Event *ContractPacketSent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractPacketSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPacketSent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractPacketSent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractPacketSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPacketSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPacketSent represents a PacketSent event raised by the Contract contract.
type ContractPacketSent struct {
	Data             []byte
	TimeoutHeight    string
	TimeoutTimestamp *big.Int
	Sequence         *big.Int
	SourcePort       string
	SourceChannel    string
	DestPort         string
	DestChannel      string
	ChannelOrdering  int32
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPacketSent is a free log retrieval operation binding the contract event 0xb1d37c3162423067c25ffe8e6b0f1ccb90c2ac2717f532ae3479cbdc9b822201.
//
// Solidity: event PacketSent(bytes data, string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering)
func (_Contract *ContractFilterer) FilterPacketSent(opts *bind.FilterOpts) (*ContractPacketSentIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "PacketSent")
	if err != nil {
		return nil, err
	}
	return &ContractPacketSentIterator{contract: _Contract.contract, event: "PacketSent", logs: logs, sub: sub}, nil
}

// WatchPacketSent is a free log subscription operation binding the contract event 0xb1d37c3162423067c25ffe8e6b0f1ccb90c2ac2717f532ae3479cbdc9b822201.
//
// Solidity: event PacketSent(bytes data, string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering)
func (_Contract *ContractFilterer) WatchPacketSent(opts *bind.WatchOpts, sink chan<- *ContractPacketSent) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "PacketSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPacketSent)
				if err := _Contract.contract.UnpackLog(event, "PacketSent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePacketSent is a log parse operation binding the contract event 0xb1d37c3162423067c25ffe8e6b0f1ccb90c2ac2717f532ae3479cbdc9b822201.
//
// Solidity: event PacketSent(bytes data, string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering)
func (_Contract *ContractFilterer) ParsePacketSent(log types.Log) (*ContractPacketSent, error) {
	event := new(ContractPacketSent)
	if err := _Contract.contract.UnpackLog(event, "PacketSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractTimeoutPacketIterator is returned from FilterTimeoutPacket and is used to iterate over the raw logs and unpacked data for TimeoutPacket events raised by the Contract contract.
type ContractTimeoutPacketIterator struct {
	Event *ContractTimeoutPacket // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractTimeoutPacketIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractTimeoutPacket)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractTimeoutPacket)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractTimeoutPacketIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractTimeoutPacketIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractTimeoutPacket represents a TimeoutPacket event raised by the Contract contract.
type ContractTimeoutPacket struct {
	TimeoutHeight    string
	TimeoutTimestamp *big.Int
	Sequence         *big.Int
	SourcePort       string
	SourceChannel    string
	DestPort         string
	DestChannel      string
	ChannelOrdering  int32
	ConnectionID     string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterTimeoutPacket is a free log retrieval operation binding the contract event 0x280b5c88e7ecdaacc40ca0de13e47206493bdee68e9656ef49e359cb36aa4c12.
//
// Solidity: event TimeoutPacket(string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering, string connectionID)
func (_Contract *ContractFilterer) FilterTimeoutPacket(opts *bind.FilterOpts) (*ContractTimeoutPacketIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "TimeoutPacket")
	if err != nil {
		return nil, err
	}
	return &ContractTimeoutPacketIterator{contract: _Contract.contract, event: "TimeoutPacket", logs: logs, sub: sub}, nil
}

// WatchTimeoutPacket is a free log subscription operation binding the contract event 0x280b5c88e7ecdaacc40ca0de13e47206493bdee68e9656ef49e359cb36aa4c12.
//
// Solidity: event TimeoutPacket(string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering, string connectionID)
func (_Contract *ContractFilterer) WatchTimeoutPacket(opts *bind.WatchOpts, sink chan<- *ContractTimeoutPacket) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "TimeoutPacket")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractTimeoutPacket)
				if err := _Contract.contract.UnpackLog(event, "TimeoutPacket", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTimeoutPacket is a log parse operation binding the contract event 0x280b5c88e7ecdaacc40ca0de13e47206493bdee68e9656ef49e359cb36aa4c12.
//
// Solidity: event TimeoutPacket(string timeoutHeight, uint256 timeoutTimestamp, uint256 sequence, string sourcePort, string sourceChannel, string destPort, string destChannel, int32 channelOrdering, string connectionID)
func (_Contract *ContractFilterer) ParseTimeoutPacket(log types.Log) (*ContractTimeoutPacket, error) {
	event := new(ContractTimeoutPacket)
	if err := _Contract.contract.UnpackLog(event, "TimeoutPacket", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
