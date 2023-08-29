// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ics20transferer

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

// Ics20transfererMetaData contains all meta data concerning the Ics20transferer contract.
var Ics20transfererMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_ibcAddr\",\"type\":\"address\"},{\"internalType\":\"contractIICS20Bank\",\"name\":\"_bank\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"sourcePort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"sourceChannel\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationPort\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"destinationChannel\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"revisionNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revisionHeight\",\"type\":\"uint256\"}],\"internalType\":\"structHeight\",\"name\":\"timeoutHeight\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"timeoutTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structPacket\",\"name\":\"packet\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"OnRecvPacket\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"someAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"portId\",\"type\":\"string\"}],\"name\":\"bindPort\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ibcAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chan\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"chanAddr\",\"type\":\"address\"}],\"name\":\"setChannelEscrowAddresses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001c8f38038062001c8f8339818101604052810190620000379190620001db565b620000576200004b620000e160201b60201c565b620000e960201b60201c565b81600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505062000298565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600081519050620001be8162000264565b92915050565b600081519050620001d5816200027e565b92915050565b60008060408385031215620001ef57600080fd5b6000620001ff85828601620001ad565b92505060206200021285828601620001c4565b9150509250929050565b6000620002298262000244565b9050919050565b60006200023d826200021c565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6200026f816200021c565b81146200027b57600080fd5b50565b620002898162000230565b81146200029557600080fd5b50565b6119e780620002a86000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806385f7175c1161005b57806385f7175c146100c65780638da5cb5b146100f65780639d19876514610114578063f2fde38b146101305761007d565b80635849f2df14610082578063696a9bf41461009e578063715018a6146100bc575b600080fd5b61009c6004803603810190610097919061108a565b61014c565b005b6100a66101b5565b6040516100b391906112a2565b60405180910390f35b6100c46101df565b005b6100e060048036038101906100db919061111f565b6101f3565b6040516100ed9190611347565b60405180910390f35b6100fe61035f565b60405161010b91906112a2565b60405180910390f35b61012e60048036038101906101299190611036565b610388565b005b61014a6004803603810190610145919061100d565b6103f7565b005b61015461047b565b80600183604051610165919061128b565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6101e761047b565b6101f160006104f9565b565b60006101fd6105bd565b73ffffffffffffffffffffffffffffffffffffffff1661021b6101b5565b73ffffffffffffffffffffffffffffffffffffffff1614610271576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610268906113a4565b60405180910390fd5b60008360a0015180602001905181019061028b91906110de565b9050600061029c82600001516105c5565b905060006102d06102b5876020015188604001516105f3565b6102c285600001516105c5565b6106dc90919063ffffffff16565b90506102e5818361077690919063ffffffff16565b6103175761030d6102f9876040015161078c565b846060015185600001518660200151610813565b9350505050610359565b600061033d8361032f89606001518a608001516105f3565b6108b490919063ffffffff16565b90506103528460600151828660200151610988565b9450505050505b92915050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b8173ffffffffffffffffffffffffffffffffffffffff1663c13b184f826040518263ffffffff1660e01b81526004016103c19190611362565b600060405180830381600087803b1580156103db57600080fd5b505af11580156103ef573d6000803e3d6000fd5b505050505050565b6103ff61047b565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561046f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161046690611384565b60405180910390fd5b610478816104f9565b50565b6104836105bd565b73ffffffffffffffffffffffffffffffffffffffff166104a161035f565b73ffffffffffffffffffffffffffffffffffffffff16146104f7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104ee906113c4565b60405180910390fd5b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600033905090565b6105cd610c0c565b600060208301905060405180604001604052808451815260200182815250915050919050565b6105fb610c0c565b6106d46106cf61063f6040518060400160405280600181526020017f2f000000000000000000000000000000000000000000000000000000000000008152506105c5565b6106c16106bc61064e876105c5565b6106ae6106a96106926040518060400160405280600181526020017f2f000000000000000000000000000000000000000000000000000000000000008152506105c5565b61069b8c6105c5565b6108b490919063ffffffff16565b6105c5565b6108b490919063ffffffff16565b6105c5565b6108b490919063ffffffff16565b6105c5565b905092915050565b6106e4610c0c565b8160000151836000015110156106fc57829050610770565b600060019050826020015184602001511461072a578251602085015160208501518281208383201493505050505b801561076b578260000151846000018181516107469190611747565b915081815250508260000151846020018181516107639190611492565b915081815250505b839150505b92915050565b6000806107838484610a26565b14905092915050565b60008060018360405161079f919061128b565b908152602001604051809103902060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561080a57600080fd5b80915050919050565b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663f24dc1da868686866040518563ffffffff1660e01b815260040161087694939291906112bd565b600060405180830381600087803b15801561089057600080fd5b505af11580156108a4573d6000803e3d6000fd5b5050505060019050949350505050565b60606000826000015184600001516108cc9190611492565b67ffffffffffffffff81111561090b577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519080825280601f01601f19166020018201604052801561093d5781602001600182028036833780820191505090505b509050600060208201905061095b8186602001518760000151610b60565b61097d85600001518261096e9190611492565b85602001518660000151610b60565b819250505092915050565b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ba7aef438585856040518463ffffffff1660e01b81526004016109e993929190611309565b600060405180830381600087803b158015610a0357600080fd5b505af1158015610a17573d6000803e3d6000fd5b50505050600190509392505050565b60008083600001519050836000015183600001511015610a4857826000015190505b60008460200151905060008460200151905060005b83811015610b3f576000808451915083519050808214610b0b5760007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90506020871015610ae557600184886020610ab59190611747565b610abf9190611492565b6008610acb9190611659565b6002610ad7919061153b565b610ae19190611747565b1990505b600081831682851603905060008114610b08578098505050505050505050610b5a565b50505b602085610b189190611492565b9450602084610b279190611492565b93505050602081610b389190611492565b9050610a5d565b5084600001518660000151610b5491906116b3565b93505050505b92915050565b5b60208110610b9f5781518352602083610b7a9190611492565b9250602082610b899190611492565b9150602081610b989190611747565b9050610b61565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90506000821115610bf6576001826020610bdc9190611747565b610100610be9919061153b565b610bf39190611747565b90505b8019835116818551168181178652505050505050565b604051806040016040528060008152602001600081525090565b6000610c39610c3484611409565b6113e4565b905082815260208101848484011115610c5157600080fd5b610c5c8482856117cd565b509392505050565b6000610c77610c728461143a565b6113e4565b905082815260208101848484011115610c8f57600080fd5b610c9a8482856117cd565b509392505050565b6000610cb5610cb08461143a565b6113e4565b905082815260208101848484011115610ccd57600080fd5b610cd88482856117dc565b509392505050565b600081359050610cef81611983565b92915050565b600081519050610d0481611983565b92915050565b600082601f830112610d1b57600080fd5b8135610d2b848260208601610c26565b91505092915050565b600082601f830112610d4557600080fd5b8135610d55848260208601610c64565b91505092915050565b600082601f830112610d6f57600080fd5b8151610d7f848260208601610ca2565b91505092915050565b600060a08284031215610d9a57600080fd5b610da460a06113e4565b9050600082015167ffffffffffffffff811115610dc057600080fd5b610dcc84828501610d5e565b6000830152506020610de084828501610ff8565b602083015250604082015167ffffffffffffffff811115610e0057600080fd5b610e0c84828501610d5e565b6040830152506060610e2084828501610cf5565b606083015250608082015167ffffffffffffffff811115610e4057600080fd5b610e4c84828501610d5e565b60808301525092915050565b600060408284031215610e6a57600080fd5b610e7460406113e4565b90506000610e8484828501610fe3565b6000830152506020610e9884828501610fe3565b60208301525092915050565b60006101208284031215610eb757600080fd5b610ec26101006113e4565b90506000610ed284828501610fe3565b600083015250602082013567ffffffffffffffff811115610ef257600080fd5b610efe84828501610d34565b602083015250604082013567ffffffffffffffff811115610f1e57600080fd5b610f2a84828501610d34565b604083015250606082013567ffffffffffffffff811115610f4a57600080fd5b610f5684828501610d34565b606083015250608082013567ffffffffffffffff811115610f7657600080fd5b610f8284828501610d34565b60808301525060a082013567ffffffffffffffff811115610fa257600080fd5b610fae84828501610d0a565b60a08301525060c0610fc284828501610e58565b60c083015250610100610fd784828501610fe3565b60e08301525092915050565b600081359050610ff28161199a565b92915050565b6000815190506110078161199a565b92915050565b60006020828403121561101f57600080fd5b600061102d84828501610ce0565b91505092915050565b6000806040838503121561104957600080fd5b600061105785828601610ce0565b925050602083013567ffffffffffffffff81111561107457600080fd5b61108085828601610d34565b9150509250929050565b6000806040838503121561109d57600080fd5b600083013567ffffffffffffffff8111156110b757600080fd5b6110c385828601610d34565b92505060206110d485828601610ce0565b9150509250929050565b6000602082840312156110f057600080fd5b600082015167ffffffffffffffff81111561110a57600080fd5b61111684828501610d88565b91505092915050565b6000806040838503121561113257600080fd5b600083013567ffffffffffffffff81111561114c57600080fd5b61115885828601610ea4565b925050602083013567ffffffffffffffff81111561117557600080fd5b61118185828601610d0a565b9150509250929050565b6111948161177b565b82525050565b6111a38161178d565b82525050565b60006111b48261146b565b6111be8185611476565b93506111ce8185602086016117dc565b6111d78161189e565b840191505092915050565b60006111ed8261146b565b6111f78185611487565b93506112078185602086016117dc565b80840191505092915050565b6000611220602683611476565b915061122b826118bc565b604082019050919050565b6000611243603483611476565b915061124e8261190b565b604082019050919050565b6000611266602083611476565b91506112718261195a565b602082019050919050565b611285816117c3565b82525050565b600061129782846111e2565b915081905092915050565b60006020820190506112b7600083018461118b565b92915050565b60006080820190506112d2600083018761118b565b6112df602083018661118b565b81810360408301526112f181856111a9565b9050611300606083018461127c565b95945050505050565b600060608201905061131e600083018661118b565b818103602083015261133081856111a9565b905061133f604083018461127c565b949350505050565b600060208201905061135c600083018461119a565b92915050565b6000602082019050818103600083015261137c81846111a9565b905092915050565b6000602082019050818103600083015261139d81611213565b9050919050565b600060208201905081810360008301526113bd81611236565b9050919050565b600060208201905081810360008301526113dd81611259565b9050919050565b60006113ee6113ff565b90506113fa828261180f565b919050565b6000604051905090565b600067ffffffffffffffff8211156114245761142361186f565b5b61142d8261189e565b9050602081019050919050565b600067ffffffffffffffff8211156114555761145461186f565b5b61145e8261189e565b9050602081019050919050565b600081519050919050565b600082825260208201905092915050565b600081905092915050565b600061149d826117c3565b91506114a8836117c3565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156114dd576114dc611840565b5b828201905092915050565b6000808291508390505b60018511156115325780860481111561150e5761150d611840565b5b600185161561151d5780820291505b808102905061152b856118af565b94506114f2565b94509492505050565b6000611546826117c3565b9150611551836117c3565b925061157e7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8484611586565b905092915050565b6000826115965760019050611652565b816115a45760009050611652565b81600181146115ba57600281146115c4576115f3565b6001915050611652565b60ff8411156115d6576115d5611840565b5b8360020a9150848211156115ed576115ec611840565b5b50611652565b5060208310610133831016604e8410600b84101617156116285782820a90508381111561162357611622611840565b5b611652565b61163584848460016114e8565b9250905081840481111561164c5761164b611840565b5b81810290505b9392505050565b6000611664826117c3565b915061166f836117c3565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156116a8576116a7611840565b5b828202905092915050565b60006116be82611799565b91506116c983611799565b9250827f80000000000000000000000000000000000000000000000000000000000000000182126000841215161561170457611703611840565b5b827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01821360008412161561173c5761173b611840565b5b828203905092915050565b6000611752826117c3565b915061175d836117c3565b9250828210156117705761176f611840565b5b828203905092915050565b6000611786826117a3565b9050919050565b60008115159050919050565b6000819050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b838110156117fa5780820151818401526020810190506117df565b83811115611809576000848401525b50505050565b6118188261189e565b810181811067ffffffffffffffff821117156118375761183661186f565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000601f19601f8301169050919050565b60008160011c9050919050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b7f4261736546756e6769626c65546f6b656e4170703a2063616c6c65722069732060008201527f6e6f74207468652049424320636f6e7472616374000000000000000000000000602082015250565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b61198c8161177b565b811461199757600080fd5b50565b6119a3816117c3565b81146119ae57600080fd5b5056fea26469706673582212204743feb0f22eba399406da06351d34a51867d752ebe3a0edc15f6a5afe1acdfc64736f6c63430008010033",
}

// Ics20transfererABI is the input ABI used to generate the binding from.
// Deprecated: Use Ics20transfererMetaData.ABI instead.
var Ics20transfererABI = Ics20transfererMetaData.ABI

// Ics20transfererBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Ics20transfererMetaData.Bin instead.
var Ics20transfererBin = Ics20transfererMetaData.Bin

// DeployIcs20transferer deploys a new Ethereum contract, binding an instance of Ics20transferer to it.
func DeployIcs20transferer(auth *bind.TransactOpts, backend bind.ContractBackend, _ibcAddr common.Address, _bank common.Address) (common.Address, *types.Transaction, *Ics20transferer, error) {
	parsed, err := Ics20transfererMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Ics20transfererBin), backend, _ibcAddr, _bank)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ics20transferer{Ics20transfererCaller: Ics20transfererCaller{contract: contract}, Ics20transfererTransactor: Ics20transfererTransactor{contract: contract}, Ics20transfererFilterer: Ics20transfererFilterer{contract: contract}}, nil
}

// Ics20transferer is an auto generated Go binding around an Ethereum contract.
type Ics20transferer struct {
	Ics20transfererCaller     // Read-only binding to the contract
	Ics20transfererTransactor // Write-only binding to the contract
	Ics20transfererFilterer   // Log filterer for contract events
}

// Ics20transfererCaller is an auto generated read-only Go binding around an Ethereum contract.
type Ics20transfererCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Ics20transfererTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Ics20transfererTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Ics20transfererFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Ics20transfererFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Ics20transfererSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Ics20transfererSession struct {
	Contract     *Ics20transferer  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Ics20transfererCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Ics20transfererCallerSession struct {
	Contract *Ics20transfererCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// Ics20transfererTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Ics20transfererTransactorSession struct {
	Contract     *Ics20transfererTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// Ics20transfererRaw is an auto generated low-level Go binding around an Ethereum contract.
type Ics20transfererRaw struct {
	Contract *Ics20transferer // Generic contract binding to access the raw methods on
}

// Ics20transfererCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Ics20transfererCallerRaw struct {
	Contract *Ics20transfererCaller // Generic read-only contract binding to access the raw methods on
}

// Ics20transfererTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Ics20transfererTransactorRaw struct {
	Contract *Ics20transfererTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIcs20transferer creates a new instance of Ics20transferer, bound to a specific deployed contract.
func NewIcs20transferer(address common.Address, backend bind.ContractBackend) (*Ics20transferer, error) {
	contract, err := bindIcs20transferer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ics20transferer{Ics20transfererCaller: Ics20transfererCaller{contract: contract}, Ics20transfererTransactor: Ics20transfererTransactor{contract: contract}, Ics20transfererFilterer: Ics20transfererFilterer{contract: contract}}, nil
}

// NewIcs20transfererCaller creates a new read-only instance of Ics20transferer, bound to a specific deployed contract.
func NewIcs20transfererCaller(address common.Address, caller bind.ContractCaller) (*Ics20transfererCaller, error) {
	contract, err := bindIcs20transferer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Ics20transfererCaller{contract: contract}, nil
}

// NewIcs20transfererTransactor creates a new write-only instance of Ics20transferer, bound to a specific deployed contract.
func NewIcs20transfererTransactor(address common.Address, transactor bind.ContractTransactor) (*Ics20transfererTransactor, error) {
	contract, err := bindIcs20transferer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Ics20transfererTransactor{contract: contract}, nil
}

// NewIcs20transfererFilterer creates a new log filterer instance of Ics20transferer, bound to a specific deployed contract.
func NewIcs20transfererFilterer(address common.Address, filterer bind.ContractFilterer) (*Ics20transfererFilterer, error) {
	contract, err := bindIcs20transferer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Ics20transfererFilterer{contract: contract}, nil
}

// bindIcs20transferer binds a generic wrapper to an already deployed contract.
func bindIcs20transferer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Ics20transfererABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ics20transferer *Ics20transfererRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ics20transferer.Contract.Ics20transfererCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ics20transferer *Ics20transfererRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ics20transferer.Contract.Ics20transfererTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ics20transferer *Ics20transfererRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ics20transferer.Contract.Ics20transfererTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ics20transferer *Ics20transfererCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ics20transferer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ics20transferer *Ics20transfererTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ics20transferer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ics20transferer *Ics20transfererTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ics20transferer.Contract.contract.Transact(opts, method, params...)
}

// IbcAddress is a free data retrieval call binding the contract method 0x696a9bf4.
//
// Solidity: function ibcAddress() view returns(address)
func (_Ics20transferer *Ics20transfererCaller) IbcAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ics20transferer.contract.Call(opts, &out, "ibcAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IbcAddress is a free data retrieval call binding the contract method 0x696a9bf4.
//
// Solidity: function ibcAddress() view returns(address)
func (_Ics20transferer *Ics20transfererSession) IbcAddress() (common.Address, error) {
	return _Ics20transferer.Contract.IbcAddress(&_Ics20transferer.CallOpts)
}

// IbcAddress is a free data retrieval call binding the contract method 0x696a9bf4.
//
// Solidity: function ibcAddress() view returns(address)
func (_Ics20transferer *Ics20transfererCallerSession) IbcAddress() (common.Address, error) {
	return _Ics20transferer.Contract.IbcAddress(&_Ics20transferer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ics20transferer *Ics20transfererCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ics20transferer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ics20transferer *Ics20transfererSession) Owner() (common.Address, error) {
	return _Ics20transferer.Contract.Owner(&_Ics20transferer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ics20transferer *Ics20transfererCallerSession) Owner() (common.Address, error) {
	return _Ics20transferer.Contract.Owner(&_Ics20transferer.CallOpts)
}

// OnRecvPacket is a paid mutator transaction binding the contract method 0x85f7175c.
//
// Solidity: function OnRecvPacket((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes ) returns(bool)
func (_Ics20transferer *Ics20transfererTransactor) OnRecvPacket(opts *bind.TransactOpts, packet Packet, arg1 []byte) (*types.Transaction, error) {
	return _Ics20transferer.contract.Transact(opts, "OnRecvPacket", packet, arg1)
}

// OnRecvPacket is a paid mutator transaction binding the contract method 0x85f7175c.
//
// Solidity: function OnRecvPacket((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes ) returns(bool)
func (_Ics20transferer *Ics20transfererSession) OnRecvPacket(packet Packet, arg1 []byte) (*types.Transaction, error) {
	return _Ics20transferer.Contract.OnRecvPacket(&_Ics20transferer.TransactOpts, packet, arg1)
}

// OnRecvPacket is a paid mutator transaction binding the contract method 0x85f7175c.
//
// Solidity: function OnRecvPacket((uint256,string,string,string,string,bytes,(uint256,uint256),uint256) packet, bytes ) returns(bool)
func (_Ics20transferer *Ics20transfererTransactorSession) OnRecvPacket(packet Packet, arg1 []byte) (*types.Transaction, error) {
	return _Ics20transferer.Contract.OnRecvPacket(&_Ics20transferer.TransactOpts, packet, arg1)
}

// BindPort is a paid mutator transaction binding the contract method 0x9d198765.
//
// Solidity: function bindPort(address someAddr, string portId) returns()
func (_Ics20transferer *Ics20transfererTransactor) BindPort(opts *bind.TransactOpts, someAddr common.Address, portId string) (*types.Transaction, error) {
	return _Ics20transferer.contract.Transact(opts, "bindPort", someAddr, portId)
}

// BindPort is a paid mutator transaction binding the contract method 0x9d198765.
//
// Solidity: function bindPort(address someAddr, string portId) returns()
func (_Ics20transferer *Ics20transfererSession) BindPort(someAddr common.Address, portId string) (*types.Transaction, error) {
	return _Ics20transferer.Contract.BindPort(&_Ics20transferer.TransactOpts, someAddr, portId)
}

// BindPort is a paid mutator transaction binding the contract method 0x9d198765.
//
// Solidity: function bindPort(address someAddr, string portId) returns()
func (_Ics20transferer *Ics20transfererTransactorSession) BindPort(someAddr common.Address, portId string) (*types.Transaction, error) {
	return _Ics20transferer.Contract.BindPort(&_Ics20transferer.TransactOpts, someAddr, portId)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ics20transferer *Ics20transfererTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ics20transferer.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ics20transferer *Ics20transfererSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ics20transferer.Contract.RenounceOwnership(&_Ics20transferer.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ics20transferer *Ics20transfererTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ics20transferer.Contract.RenounceOwnership(&_Ics20transferer.TransactOpts)
}

// SetChannelEscrowAddresses is a paid mutator transaction binding the contract method 0x5849f2df.
//
// Solidity: function setChannelEscrowAddresses(string chan, address chanAddr) returns()
func (_Ics20transferer *Ics20transfererTransactor) SetChannelEscrowAddresses(opts *bind.TransactOpts, arg0 string, chanAddr common.Address) (*types.Transaction, error) {
	return _Ics20transferer.contract.Transact(opts, "setChannelEscrowAddresses", arg0, chanAddr)
}

// SetChannelEscrowAddresses is a paid mutator transaction binding the contract method 0x5849f2df.
//
// Solidity: function setChannelEscrowAddresses(string chan, address chanAddr) returns()
func (_Ics20transferer *Ics20transfererSession) SetChannelEscrowAddresses(arg0 string, chanAddr common.Address) (*types.Transaction, error) {
	return _Ics20transferer.Contract.SetChannelEscrowAddresses(&_Ics20transferer.TransactOpts, arg0, chanAddr)
}

// SetChannelEscrowAddresses is a paid mutator transaction binding the contract method 0x5849f2df.
//
// Solidity: function setChannelEscrowAddresses(string chan, address chanAddr) returns()
func (_Ics20transferer *Ics20transfererTransactorSession) SetChannelEscrowAddresses(arg0 string, chanAddr common.Address) (*types.Transaction, error) {
	return _Ics20transferer.Contract.SetChannelEscrowAddresses(&_Ics20transferer.TransactOpts, arg0, chanAddr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ics20transferer *Ics20transfererTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ics20transferer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ics20transferer *Ics20transfererSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ics20transferer.Contract.TransferOwnership(&_Ics20transferer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ics20transferer *Ics20transfererTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ics20transferer.Contract.TransferOwnership(&_Ics20transferer.TransactOpts, newOwner)
}

// Ics20transfererOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ics20transferer contract.
type Ics20transfererOwnershipTransferredIterator struct {
	Event *Ics20transfererOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *Ics20transfererOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Ics20transfererOwnershipTransferred)
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
		it.Event = new(Ics20transfererOwnershipTransferred)
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
func (it *Ics20transfererOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Ics20transfererOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Ics20transfererOwnershipTransferred represents a OwnershipTransferred event raised by the Ics20transferer contract.
type Ics20transfererOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ics20transferer *Ics20transfererFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Ics20transfererOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ics20transferer.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Ics20transfererOwnershipTransferredIterator{contract: _Ics20transferer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ics20transferer *Ics20transfererFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *Ics20transfererOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ics20transferer.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Ics20transfererOwnershipTransferred)
				if err := _Ics20transferer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ics20transferer *Ics20transfererFilterer) ParseOwnershipTransferred(log types.Log) (*Ics20transfererOwnershipTransferred, error) {
	event := new(Ics20transfererOwnershipTransferred)
	if err := _Ics20transferer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
