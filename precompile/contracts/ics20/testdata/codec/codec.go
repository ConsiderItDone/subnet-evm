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
	Sender   []byte
	Receiver []byte
	Memo     []byte
}

// CodecMetaData contains all meta data concerning the Codec contract.
var CodecMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"rawdata\",\"type\":\"bytes\"}],\"name\":\"decode\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"denom\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"memo\",\"type\":\"bytes\"}],\"internalType\":\"structFungibleTokenPacketData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"denom\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"memo\",\"type\":\"bytes\"}],\"internalType\":\"structFungibleTokenPacketData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"encode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506108d6806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063486f3da81461003b578063e5c5e9a31461006b575b600080fd5b610055600480360381019061005091906104f1565b61009b565b60405161006291906106b8565b60405180910390f35b610085600480360381019061008091906104b0565b6100c4565b60405161009291906106da565b60405180910390f35b6060816040516020016100ae91906106da565b6040516020818303038152906040529050919050565b6100cc6100e7565b818060200190518101906100e09190610532565b9050919050565b6040518060a0016040528060608152602001600081526020016060815260200160608152602001606081525090565b600061012961012484610721565b6106fc565b90508281526020810184848401111561014157600080fd5b61014c8482856107d6565b509392505050565b600061016761016284610721565b6106fc565b90508281526020810184848401111561017f57600080fd5b61018a8482856107e5565b509392505050565b60006101a56101a084610752565b6106fc565b9050828152602081018484840111156101bd57600080fd5b6101c88482856107d6565b509392505050565b60006101e36101de84610752565b6106fc565b9050828152602081018484840111156101fb57600080fd5b6102068482856107e5565b509392505050565b600082601f83011261021f57600080fd5b813561022f848260208601610116565b91505092915050565b600082601f83011261024957600080fd5b8151610259848260208601610154565b91505092915050565b600082601f83011261027357600080fd5b8135610283848260208601610192565b91505092915050565b600082601f83011261029d57600080fd5b81516102ad8482602086016101d0565b91505092915050565b600060a082840312156102c857600080fd5b6102d260a06106fc565b9050600082013567ffffffffffffffff8111156102ee57600080fd5b6102fa84828501610262565b600083015250602061030e84828501610486565b602083015250604082013567ffffffffffffffff81111561032e57600080fd5b61033a8482850161020e565b604083015250606082013567ffffffffffffffff81111561035a57600080fd5b6103668482850161020e565b606083015250608082013567ffffffffffffffff81111561038657600080fd5b6103928482850161020e565b60808301525092915050565b600060a082840312156103b057600080fd5b6103ba60a06106fc565b9050600082015167ffffffffffffffff8111156103d657600080fd5b6103e28482850161028c565b60008301525060206103f68482850161049b565b602083015250604082015167ffffffffffffffff81111561041657600080fd5b61042284828501610238565b604083015250606082015167ffffffffffffffff81111561044257600080fd5b61044e84828501610238565b606083015250608082015167ffffffffffffffff81111561046e57600080fd5b61047a84828501610238565b60808301525092915050565b60008135905061049581610889565b92915050565b6000815190506104aa81610889565b92915050565b6000602082840312156104c257600080fd5b600082013567ffffffffffffffff8111156104dc57600080fd5b6104e88482850161020e565b91505092915050565b60006020828403121561050357600080fd5b600082013567ffffffffffffffff81111561051d57600080fd5b610529848285016102b6565b91505092915050565b60006020828403121561054457600080fd5b600082015167ffffffffffffffff81111561055e57600080fd5b61056a8482850161039e565b91505092915050565b600061057e82610783565b6105888185610799565b93506105988185602086016107e5565b6105a181610878565b840191505092915050565b60006105b782610783565b6105c181856107aa565b93506105d18185602086016107e5565b6105da81610878565b840191505092915050565b60006105f08261078e565b6105fa81856107bb565b935061060a8185602086016107e5565b61061381610878565b840191505092915050565b600060a083016000830151848203600086015261063b82826105e5565b915050602083015161065060208601826106a9565b50604083015184820360408601526106688282610573565b915050606083015184820360608601526106828282610573565b9150506080830151848203608086015261069c8282610573565b9150508091505092915050565b6106b2816107cc565b82525050565b600060208201905081810360008301526106d281846105ac565b905092915050565b600060208201905081810360008301526106f4818461061e565b905092915050565b6000610706610717565b90506107128282610818565b919050565b6000604051905090565b600067ffffffffffffffff82111561073c5761073b610849565b5b61074582610878565b9050602081019050919050565b600067ffffffffffffffff82111561076d5761076c610849565b5b61077682610878565b9050602081019050919050565b600081519050919050565b600081519050919050565b600082825260208201905092915050565b600082825260208201905092915050565b600082825260208201905092915050565b6000819050919050565b82818337600083830152505050565b60005b838110156108035780820151818401526020810190506107e8565b83811115610812576000848401525b50505050565b61082182610878565b810181811067ffffffffffffffff821117156108405761083f610849565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000601f19601f8301169050919050565b610892816107cc565b811461089d57600080fd5b5056fea26469706673582212202fb0938ade95f4015ce19a5830cc242efe50ce375391794f24d3885cc424d11564736f6c63430008010033",
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
// Solidity: function decode(bytes rawdata) pure returns((string,uint256,bytes,bytes,bytes))
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
// Solidity: function decode(bytes rawdata) pure returns((string,uint256,bytes,bytes,bytes))
func (_Codec *CodecSession) Decode(rawdata []byte) (FungibleTokenPacketData, error) {
	return _Codec.Contract.Decode(&_Codec.CallOpts, rawdata)
}

// Decode is a free data retrieval call binding the contract method 0xe5c5e9a3.
//
// Solidity: function decode(bytes rawdata) pure returns((string,uint256,bytes,bytes,bytes))
func (_Codec *CodecCallerSession) Decode(rawdata []byte) (FungibleTokenPacketData, error) {
	return _Codec.Contract.Decode(&_Codec.CallOpts, rawdata)
}

// Encode is a free data retrieval call binding the contract method 0x486f3da8.
//
// Solidity: function encode((string,uint256,bytes,bytes,bytes) data) pure returns(bytes)
func (_Codec *CodecCaller) Encode(opts *bind.CallOpts, data FungibleTokenPacketData) ([]byte, error) {
	var out []interface{}
	err := _Codec.contract.Call(opts, &out, "encode", data)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Encode is a free data retrieval call binding the contract method 0x486f3da8.
//
// Solidity: function encode((string,uint256,bytes,bytes,bytes) data) pure returns(bytes)
func (_Codec *CodecSession) Encode(data FungibleTokenPacketData) ([]byte, error) {
	return _Codec.Contract.Encode(&_Codec.CallOpts, data)
}

// Encode is a free data retrieval call binding the contract method 0x486f3da8.
//
// Solidity: function encode((string,uint256,bytes,bytes,bytes) data) pure returns(bytes)
func (_Codec *CodecCallerSession) Encode(data FungibleTokenPacketData) ([]byte, error) {
	return _Codec.Contract.Encode(&_Codec.CallOpts, data)
}
