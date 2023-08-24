// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package codec

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

// FungibleTokenPacketData is an auto generated low-level Go binding around an user-defined struct.
type FungibleTokenPacketData struct {
	Denom    string
	Amount   *big.Int
	Sender   string
	Receiver common.Address
	Memo     string
}

// CodecMetaData contains all meta data concerning the Codec contract.
var CodecMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"rawdata\",\"type\":\"bytes\"}],\"name\":\"decode\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"denom\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"sender\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"internalType\":\"structFungibleTokenPacketData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"denom\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"sender\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"internalType\":\"structFungibleTokenPacketData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"encode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061085e806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806309aa1b821461003b578063e5c5e9a31461006b575b600080fd5b61005560048036038101906100509190610499565b61009b565b604051610062919061062f565b60405180910390f35b61008560048036038101906100809190610458565b6100c4565b6040516100929190610651565b60405180910390f35b6060816040516020016100ae9190610651565b6040516020818303038152906040529050919050565b6100cc6100e7565b818060200190518101906100e091906104da565b9050919050565b6040518060a00160405280606081526020016000815260200160608152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001606081525090565b600061013f61013a846106a4565b610673565b90508281526020810184848401111561015757600080fd5b610162848285610778565b509392505050565b600061017d610178846106d4565b610673565b90508281526020810184848401111561019557600080fd5b6101a0848285610778565b509392505050565b60006101bb6101b6846106d4565b610673565b9050828152602081018484840111156101d357600080fd5b6101de848285610787565b509392505050565b6000813590506101f5816107fa565b92915050565b60008151905061020a816107fa565b92915050565b600082601f83011261022157600080fd5b813561023184826020860161012c565b91505092915050565b600082601f83011261024b57600080fd5b813561025b84826020860161016a565b91505092915050565b600082601f83011261027557600080fd5b81516102858482602086016101a8565b91505092915050565b600060a082840312156102a057600080fd5b6102aa60a0610673565b9050600082013567ffffffffffffffff8111156102c657600080fd5b6102d28482850161023a565b60008301525060206102e68482850161042e565b602083015250604082013567ffffffffffffffff81111561030657600080fd5b6103128482850161023a565b6040830152506060610326848285016101e6565b606083015250608082013567ffffffffffffffff81111561034657600080fd5b6103528482850161023a565b60808301525092915050565b600060a0828403121561037057600080fd5b61037a60a0610673565b9050600082015167ffffffffffffffff81111561039657600080fd5b6103a284828501610264565b60008301525060206103b684828501610443565b602083015250604082015167ffffffffffffffff8111156103d657600080fd5b6103e284828501610264565b60408301525060606103f6848285016101fb565b606083015250608082015167ffffffffffffffff81111561041657600080fd5b61042284828501610264565b60808301525092915050565b60008135905061043d81610811565b92915050565b60008151905061045281610811565b92915050565b60006020828403121561046a57600080fd5b600082013567ffffffffffffffff81111561048457600080fd5b61049084828501610210565b91505092915050565b6000602082840312156104ab57600080fd5b600082013567ffffffffffffffff8111156104c557600080fd5b6104d18482850161028e565b91505092915050565b6000602082840312156104ec57600080fd5b600082015167ffffffffffffffff81111561050657600080fd5b6105128482850161035e565b91505092915050565b6105248161073c565b82525050565b600061053582610704565b61053f818561071a565b935061054f818560208601610787565b610558816107e9565b840191505092915050565b600061056e8261070f565b610578818561072b565b9350610588818560208601610787565b610591816107e9565b840191505092915050565b600060a08301600083015184820360008601526105b98282610563565b91505060208301516105ce6020860182610620565b50604083015184820360408601526105e68282610563565b91505060608301516105fb606086018261051b565b50608083015184820360808601526106138282610563565b9150508091505092915050565b6106298161076e565b82525050565b60006020820190508181036000830152610649818461052a565b905092915050565b6000602082019050818103600083015261066b818461059c565b905092915050565b6000604051905081810181811067ffffffffffffffff8211171561069a576106996107ba565b5b8060405250919050565b600067ffffffffffffffff8211156106bf576106be6107ba565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156106ef576106ee6107ba565b5b601f19601f8301169050602081019050919050565b600081519050919050565b600081519050919050565b600082825260208201905092915050565b600082825260208201905092915050565b60006107478261074e565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b838110156107a557808201518184015260208101905061078a565b838111156107b4576000848401525b50505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000601f19601f8301169050919050565b6108038161073c565b811461080e57600080fd5b50565b61081a8161076e565b811461082557600080fd5b5056fea2646970667358221220d58cb90bfa7801670055e928d24097a70ae060d08f7e5a5f610cbfaca87e062364736f6c63430008000033",
}

// CodecABI is the input ABI used to generate the binding from.
// Deprecated: Use CodecMetaData.ABI instead.
var CodecABI = CodecMetaData.ABI

// CodecBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CodecMetaData.Bin instead.
var CodecBin = CodecMetaData.Bin

// DeployCodec deploys a new Ethereum contract, binding an instance of Codec to it.
func DeployCodec(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Codec, error) {
	parsed, err := CodecMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CodecBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Codec{CodecCaller: CodecCaller{contract: contract}, CodecTransactor: CodecTransactor{contract: contract}, CodecFilterer: CodecFilterer{contract: contract}}, nil
}

// Codec is an auto generated Go binding around an Ethereum contract.
type Codec struct {
	CodecCaller     // Read-only binding to the contract
	CodecTransactor // Write-only binding to the contract
	CodecFilterer   // Log filterer for contract events
}

// CodecCaller is an auto generated read-only Go binding around an Ethereum contract.
type CodecCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CodecTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CodecTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CodecFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CodecFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CodecSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CodecSession struct {
	Contract     *Codec            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CodecCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CodecCallerSession struct {
	Contract *CodecCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// CodecTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CodecTransactorSession struct {
	Contract     *CodecTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CodecRaw is an auto generated low-level Go binding around an Ethereum contract.
type CodecRaw struct {
	Contract *Codec // Generic contract binding to access the raw methods on
}

// CodecCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CodecCallerRaw struct {
	Contract *CodecCaller // Generic read-only contract binding to access the raw methods on
}

// CodecTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CodecTransactorRaw struct {
	Contract *CodecTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCodec creates a new instance of Codec, bound to a specific deployed contract.
func NewCodec(address common.Address, backend bind.ContractBackend) (*Codec, error) {
	contract, err := bindCodec(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Codec{CodecCaller: CodecCaller{contract: contract}, CodecTransactor: CodecTransactor{contract: contract}, CodecFilterer: CodecFilterer{contract: contract}}, nil
}

// NewCodecCaller creates a new read-only instance of Codec, bound to a specific deployed contract.
func NewCodecCaller(address common.Address, caller bind.ContractCaller) (*CodecCaller, error) {
	contract, err := bindCodec(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CodecCaller{contract: contract}, nil
}

// NewCodecTransactor creates a new write-only instance of Codec, bound to a specific deployed contract.
func NewCodecTransactor(address common.Address, transactor bind.ContractTransactor) (*CodecTransactor, error) {
	contract, err := bindCodec(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CodecTransactor{contract: contract}, nil
}

// NewCodecFilterer creates a new log filterer instance of Codec, bound to a specific deployed contract.
func NewCodecFilterer(address common.Address, filterer bind.ContractFilterer) (*CodecFilterer, error) {
	contract, err := bindCodec(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CodecFilterer{contract: contract}, nil
}

// bindCodec binds a generic wrapper to an already deployed contract.
func bindCodec(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CodecABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Codec *CodecRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Codec.Contract.CodecCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Codec *CodecRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Codec.Contract.CodecTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Codec *CodecRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Codec.Contract.CodecTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Codec *CodecCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Codec.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Codec *CodecTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Codec.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Codec *CodecTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Codec.Contract.contract.Transact(opts, method, params...)
}

// Decode is a free data retrieval call binding the contract method 0xe5c5e9a3.
//
// Solidity: function decode(bytes rawdata) pure returns((string,uint256,string,address,string))
func (_Codec *CodecCaller) Decode(opts *bind.CallOpts, rawdata []byte) (FungibleTokenPacketData, error) {
	var out []interface{}
	err := _Codec.contract.Call(opts, &out, "decode", rawdata)

	if err != nil {
		return *new(FungibleTokenPacketData), err
	}

	out0 := *abi.ConvertType(out[0], new(FungibleTokenPacketData)).(*FungibleTokenPacketData)

	return out0, err

}

// Decode is a free data retrieval call binding the contract method 0xe5c5e9a3.
//
// Solidity: function decode(bytes rawdata) pure returns((string,uint256,string,address,string))
func (_Codec *CodecSession) Decode(rawdata []byte) (FungibleTokenPacketData, error) {
	return _Codec.Contract.Decode(&_Codec.CallOpts, rawdata)
}

// Decode is a free data retrieval call binding the contract method 0xe5c5e9a3.
//
// Solidity: function decode(bytes rawdata) pure returns((string,uint256,string,address,string))
func (_Codec *CodecCallerSession) Decode(rawdata []byte) (FungibleTokenPacketData, error) {
	return _Codec.Contract.Decode(&_Codec.CallOpts, rawdata)
}

// Encode is a free data retrieval call binding the contract method 0x09aa1b82.
//
// Solidity: function encode((string,uint256,string,address,string) data) pure returns(bytes)
func (_Codec *CodecCaller) Encode(opts *bind.CallOpts, data FungibleTokenPacketData) ([]byte, error) {
	var out []interface{}
	err := _Codec.contract.Call(opts, &out, "encode", data)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Encode is a free data retrieval call binding the contract method 0x09aa1b82.
//
// Solidity: function encode((string,uint256,string,address,string) data) pure returns(bytes)
func (_Codec *CodecSession) Encode(data FungibleTokenPacketData) ([]byte, error) {
	return _Codec.Contract.Encode(&_Codec.CallOpts, data)
}

// Encode is a free data retrieval call binding the contract method 0x09aa1b82.
//
// Solidity: function encode((string,uint256,string,address,string) data) pure returns(bytes)
func (_Codec *CodecCallerSession) Encode(data FungibleTokenPacketData) ([]byte, error) {
	return _Codec.Contract.Encode(&_Codec.CallOpts, data)
}
