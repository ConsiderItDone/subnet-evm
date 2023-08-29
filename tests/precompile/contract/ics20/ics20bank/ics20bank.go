// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ics20bank

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

// Ics20bankMetaData contains all meta data concerning the Ics20bank contract.
var Ics20bankMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OPERATOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50620000537fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775620000476200005960201b60201c565b6200006160201b60201c565b620001d2565b600033905090565b6200007382826200007760201b60201c565b5050565b6200008982826200016860201b60201c565b6200016457600160008084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550620001096200005960201b60201c565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b600080600084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b61232a80620001e26000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c8063b3ab15fb11610097578063d547741f11610066578063d547741f146102a9578063f24dc1da146102c5578063f45346dc146102e1578063f5b541a6146102fd57610100565b8063b3ab15fb14610225578063b9b092c814610241578063ba7aef4314610271578063c45b71de1461028d57610100565b806369328dec116100d357806369328dec1461019d57806375b238fc146101b957806391d14854146101d7578063a217fddf1461020757610100565b806301ffc9a714610105578063248a9ca3146101355780632f2ff15d1461016557806336568abe14610181575b600080fd5b61011f600480360381019061011a9190611887565b61031b565b60405161012c9190611bc6565b60405180910390f35b61014f600480360381019061014a9190611822565b610395565b60405161015c9190611be1565b60405180910390f35b61017f600480360381019061017a919061184b565b6103b4565b005b61019b6004803603810190610196919061184b565b6103d5565b005b6101b760048036038101906101b291906117aa565b610458565b005b6101c1610535565b6040516101ce9190611be1565b60405180910390f35b6101f160048036038101906101ec919061184b565b610559565b6040516101fe9190611bc6565b60405180910390f35b61020f6105c3565b60405161021c9190611be1565b60405180910390f35b61023f600480360381019061023a919061163d565b6105ca565b005b61025b600480360381019061025691906116e6565b610667565b6040516102689190611d3e565b60405180910390f35b61028b6004803603810190610286919061173e565b61073f565b005b6102a760048036038101906102a2919061173e565b610804565b005b6102c360048036038101906102be919061184b565b6108c9565b005b6102df60048036038101906102da9190611666565b6108ea565b005b6102fb60048036038101906102f691906117aa565b610bfa565b005b610305610cd9565b6040516103129190611be1565b60405180910390f35b60007f7965db0b000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916148061038e575061038d82610cfd565b5b9050919050565b6000806000838152602001908152602001600020600101549050919050565b6103bd82610395565b6103c681610d67565b6103d08383610d7b565b505050565b6103dd610e5b565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461044a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161044190611cfe565b60405180910390fd5b6104548282610e63565b5050565b6104778373ffffffffffffffffffffffffffffffffffffffff16610f44565b61048057600080fd5b61049a61048b610e5b565b61049485610f67565b84610f79565b8273ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82846040518363ffffffff1660e01b81526004016104d5929190611b9d565b602060405180830381600087803b1580156104ef57600080fd5b505af1158015610503573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061052791906117f9565b61053057600080fd5b505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b600080600084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b6000801b81565b6105fb7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756105f6610e5b565b610559565b61063a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161063190611cbe565b60405180910390fd5b6106647f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b929826110f3565b50565b60008073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1614156106d8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106cf90611cde565b60405180910390fd5b600183836040516106ea929190611afc565b908152602001604051809103902060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490509392505050565b6107707f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b92961076b610e5b565b610559565b6107af576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107a690611d1e565b60405180910390fd5b6107fe8484848080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505083611101565b50505050565b6108357f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b929610830610e5b565b610559565b610874576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161086b90611d1e565b60405180910390fd5b6108c38484848080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505083610f79565b50505050565b6108d282610395565b6108db81610d67565b6108e58383610e63565b505050565b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16141561095a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161095190611c3e565b60405180910390fd5b610962610e5b565b73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1614806109c857506109c77f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b9296109c2610e5b565b610559565b5b610a07576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109fe90611c7e565b60405180910390fd5b600060018484604051610a1b929190611afc565b908152602001604051809103902060008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610aac576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610aa390611c9e565b60405180910390fd5b8181610ab89190611e30565b60018585604051610aca929190611afc565b908152602001604051809103902060008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508160018585604051610b2c929190611afc565b908152602001604051809103902060008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610b869190611d80565b925050819055508473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610bea9190611d3e565b60405180910390a3505050505050565b610c198373ffffffffffffffffffffffffffffffffffffffff16610f44565b610c2257600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd610c46610e5b565b30856040518463ffffffff1660e01b8152600401610c6693929190611b66565b602060405180830381600087803b158015610c8057600080fd5b505af1158015610c94573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cb891906117f9565b610cc157600080fd5b610cd481610cce85610f67565b84611101565b505050565b7f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b92981565b60007f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b610d7881610d73610e5b565b6111de565b50565b610d858282610559565b610e5757600160008084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550610dfc610e5b565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b600033905090565b610e6d8282610559565b15610f4057600080600084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550610ee5610e5b565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b60405160405180910390a45b5050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b6060610f7282611263565b9050919050565b6000600183604051610f8b9190611b15565b908152602001604051809103902060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508181101561101c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161101390611c5e565b60405180910390fd5b81816110289190611e30565b6001846040516110389190611b15565b908152602001604051809103902060008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516110e59190611d3e565b60405180910390a350505050565b6110fd8282610d7b565b5050565b806001836040516111129190611b15565b908152602001604051809103902060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461116c9190611d80565b925050819055508273ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516111d19190611d3e565b60405180910390a3505050565b6111e88282610559565b61125f576111f581611263565b6112038360001c6020611290565b604051602001611214929190611b2c565b6040516020818303038152906040526040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112569190611bfc565b60405180910390fd5b5050565b60606112898273ffffffffffffffffffffffffffffffffffffffff16601460ff16611290565b9050919050565b6060600060028360026112a39190611dd6565b6112ad9190611d80565b67ffffffffffffffff8111156112ec577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519080825280601f01601f19166020018201604052801561131e5781602001600182028036833780820191505090505b5090507f30000000000000000000000000000000000000000000000000000000000000008160008151811061137c577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053507f780000000000000000000000000000000000000000000000000000000000000081600181518110611406577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600060018460026114469190611dd6565b6114509190611d80565b90505b600181111561153c577f3031323334353637383961626364656600000000000000000000000000000000600f8616601081106114b8577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b1a60f81b8282815181106114f5577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600485901c94508061153590611f24565b9050611453565b5060008414611580576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161157790611c1e565b60405180910390fd5b8091505092915050565b60008135905061159981612281565b92915050565b6000815190506115ae81612298565b92915050565b6000813590506115c3816122af565b92915050565b6000813590506115d8816122c6565b92915050565b60008083601f8401126115f057600080fd5b8235905067ffffffffffffffff81111561160957600080fd5b60208301915083600182028301111561162157600080fd5b9250929050565b600081359050611637816122dd565b92915050565b60006020828403121561164f57600080fd5b600061165d8482850161158a565b91505092915050565b60008060008060006080868803121561167e57600080fd5b600061168c8882890161158a565b955050602061169d8882890161158a565b945050604086013567ffffffffffffffff8111156116ba57600080fd5b6116c6888289016115de565b935093505060606116d988828901611628565b9150509295509295909350565b6000806000604084860312156116fb57600080fd5b60006117098682870161158a565b935050602084013567ffffffffffffffff81111561172657600080fd5b611732868287016115de565b92509250509250925092565b6000806000806060858703121561175457600080fd5b60006117628782880161158a565b945050602085013567ffffffffffffffff81111561177f57600080fd5b61178b878288016115de565b9350935050604061179e87828801611628565b91505092959194509250565b6000806000606084860312156117bf57600080fd5b60006117cd8682870161158a565b93505060206117de86828701611628565b92505060406117ef8682870161158a565b9150509250925092565b60006020828403121561180b57600080fd5b60006118198482850161159f565b91505092915050565b60006020828403121561183457600080fd5b6000611842848285016115b4565b91505092915050565b6000806040838503121561185e57600080fd5b600061186c858286016115b4565b925050602061187d8582860161158a565b9150509250929050565b60006020828403121561189957600080fd5b60006118a7848285016115c9565b91505092915050565b6118b981611e64565b82525050565b6118c881611e76565b82525050565b6118d781611e82565b82525050565b60006118e98385611d75565b93506118f6838584611ee2565b82840190509392505050565b600061190d82611d59565b6119178185611d64565b9350611927818560208601611ef1565b61193081611f7d565b840191505092915050565b600061194682611d59565b6119508185611d75565b9350611960818560208601611ef1565b80840191505092915050565b6000611979602083611d64565b915061198482611f8e565b602082019050919050565b600061199c602783611d64565b91506119a782611fb7565b604082019050919050565b60006119bf602683611d64565b91506119ca82612006565b604082019050919050565b60006119e2602b83611d64565b91506119ed82612055565b604082019050919050565b6000611a05602c83611d64565b9150611a10826120a4565b604082019050919050565b6000611a28602883611d64565b9150611a33826120f3565b604082019050919050565b6000611a4b601783611d75565b9150611a5682612142565b601782019050919050565b6000611a6e602d83611d64565b9150611a798261216b565b604082019050919050565b6000611a91601183611d75565b9150611a9c826121ba565b601182019050919050565b6000611ab4602f83611d64565b9150611abf826121e3565b604082019050919050565b6000611ad7602883611d64565b9150611ae282612232565b604082019050919050565b611af681611ed8565b82525050565b6000611b098284866118dd565b91508190509392505050565b6000611b21828461193b565b915081905092915050565b6000611b3782611a3e565b9150611b43828561193b565b9150611b4e82611a84565b9150611b5a828461193b565b91508190509392505050565b6000606082019050611b7b60008301866118b0565b611b8860208301856118b0565b611b956040830184611aed565b949350505050565b6000604082019050611bb260008301856118b0565b611bbf6020830184611aed565b9392505050565b6000602082019050611bdb60008301846118bf565b92915050565b6000602082019050611bf660008301846118ce565b92915050565b60006020820190508181036000830152611c168184611902565b905092915050565b60006020820190508181036000830152611c378161196c565b9050919050565b60006020820190508181036000830152611c578161198f565b9050919050565b60006020820190508181036000830152611c77816119b2565b9050919050565b60006020820190508181036000830152611c97816119d5565b9050919050565b60006020820190508181036000830152611cb7816119f8565b9050919050565b60006020820190508181036000830152611cd781611a1b565b9050919050565b60006020820190508181036000830152611cf781611a61565b9050919050565b60006020820190508181036000830152611d1781611aa7565b9050919050565b60006020820190508181036000830152611d3781611aca565b9050919050565b6000602082019050611d536000830184611aed565b92915050565b600081519050919050565b600082825260208201905092915050565b600081905092915050565b6000611d8b82611ed8565b9150611d9683611ed8565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611dcb57611dca611f4e565b5b828201905092915050565b6000611de182611ed8565b9150611dec83611ed8565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611e2557611e24611f4e565b5b828202905092915050565b6000611e3b82611ed8565b9150611e4683611ed8565b925082821015611e5957611e58611f4e565b5b828203905092915050565b6000611e6f82611eb8565b9050919050565b60008115159050919050565b6000819050919050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b83811015611f0f578082015181840152602081019050611ef4565b83811115611f1e576000848401525b50505050565b6000611f2f82611ed8565b91506000821415611f4357611f42611f4e565b5b600182039050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000601f19601f8301169050919050565b7f537472696e67733a20686578206c656e67746820696e73756666696369656e74600082015250565b7f494353323042616e6b3a207472616e7366657220746f20746865207a65726f2060008201527f6164647265737300000000000000000000000000000000000000000000000000602082015250565b7f494353323042616e6b3a206275726e20616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b7f494353323042616e6b3a2063616c6c6572206973206e6f74206f776e6572206e60008201527f6f7220617070726f766564000000000000000000000000000000000000000000602082015250565b7f494353323042616e6b3a20696e73756666696369656e742062616c616e63652060008201527f666f72207472616e736665720000000000000000000000000000000000000000602082015250565b7f6d75737420686176652061646d696e20726f6c6520746f20736574206e65772060008201527f6f70657261746f72000000000000000000000000000000000000000000000000602082015250565b7f416363657373436f6e74726f6c3a206163636f756e7420000000000000000000600082015250565b7f494353323042616e6b3a2062616c616e636520717565727920666f722074686560008201527f207a65726f206164647265737300000000000000000000000000000000000000602082015250565b7f206973206d697373696e6720726f6c6520000000000000000000000000000000600082015250565b7f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560008201527f20726f6c657320666f722073656c660000000000000000000000000000000000602082015250565b7f494353323042616e6b3a206d7573742068617665206d696e74657220726f6c6560008201527f20746f206d696e74000000000000000000000000000000000000000000000000602082015250565b61228a81611e64565b811461229557600080fd5b50565b6122a181611e76565b81146122ac57600080fd5b50565b6122b881611e82565b81146122c357600080fd5b50565b6122cf81611e8c565b81146122da57600080fd5b50565b6122e681611ed8565b81146122f157600080fd5b5056fea2646970667358221220ba80ae06b565776ef83792046ffdb01459e6a4efd7e41e8e70ff48821e5ebfb764736f6c63430008010033",
}

// Ics20bankABI is the input ABI used to generate the binding from.
// Deprecated: Use Ics20bankMetaData.ABI instead.
var Ics20bankABI = Ics20bankMetaData.ABI

// Ics20bankBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Ics20bankMetaData.Bin instead.
var Ics20bankBin = Ics20bankMetaData.Bin

// DeployIcs20bank deploys a new Ethereum contract, binding an instance of Ics20bank to it.
func DeployIcs20bank(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Ics20bank, error) {
	parsed, err := Ics20bankMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Ics20bankBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ics20bank{Ics20bankCaller: Ics20bankCaller{contract: contract}, Ics20bankTransactor: Ics20bankTransactor{contract: contract}, Ics20bankFilterer: Ics20bankFilterer{contract: contract}}, nil
}

// Ics20bank is an auto generated Go binding around an Ethereum contract.
type Ics20bank struct {
	Ics20bankCaller     // Read-only binding to the contract
	Ics20bankTransactor // Write-only binding to the contract
	Ics20bankFilterer   // Log filterer for contract events
}

// Ics20bankCaller is an auto generated read-only Go binding around an Ethereum contract.
type Ics20bankCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Ics20bankTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Ics20bankTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Ics20bankFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Ics20bankFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Ics20bankSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Ics20bankSession struct {
	Contract     *Ics20bank        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Ics20bankCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Ics20bankCallerSession struct {
	Contract *Ics20bankCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// Ics20bankTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Ics20bankTransactorSession struct {
	Contract     *Ics20bankTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// Ics20bankRaw is an auto generated low-level Go binding around an Ethereum contract.
type Ics20bankRaw struct {
	Contract *Ics20bank // Generic contract binding to access the raw methods on
}

// Ics20bankCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Ics20bankCallerRaw struct {
	Contract *Ics20bankCaller // Generic read-only contract binding to access the raw methods on
}

// Ics20bankTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Ics20bankTransactorRaw struct {
	Contract *Ics20bankTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIcs20bank creates a new instance of Ics20bank, bound to a specific deployed contract.
func NewIcs20bank(address common.Address, backend bind.ContractBackend) (*Ics20bank, error) {
	contract, err := bindIcs20bank(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ics20bank{Ics20bankCaller: Ics20bankCaller{contract: contract}, Ics20bankTransactor: Ics20bankTransactor{contract: contract}, Ics20bankFilterer: Ics20bankFilterer{contract: contract}}, nil
}

// NewIcs20bankCaller creates a new read-only instance of Ics20bank, bound to a specific deployed contract.
func NewIcs20bankCaller(address common.Address, caller bind.ContractCaller) (*Ics20bankCaller, error) {
	contract, err := bindIcs20bank(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Ics20bankCaller{contract: contract}, nil
}

// NewIcs20bankTransactor creates a new write-only instance of Ics20bank, bound to a specific deployed contract.
func NewIcs20bankTransactor(address common.Address, transactor bind.ContractTransactor) (*Ics20bankTransactor, error) {
	contract, err := bindIcs20bank(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Ics20bankTransactor{contract: contract}, nil
}

// NewIcs20bankFilterer creates a new log filterer instance of Ics20bank, bound to a specific deployed contract.
func NewIcs20bankFilterer(address common.Address, filterer bind.ContractFilterer) (*Ics20bankFilterer, error) {
	contract, err := bindIcs20bank(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Ics20bankFilterer{contract: contract}, nil
}

// bindIcs20bank binds a generic wrapper to an already deployed contract.
func bindIcs20bank(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Ics20bankABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ics20bank *Ics20bankRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ics20bank.Contract.Ics20bankCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ics20bank *Ics20bankRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ics20bank.Contract.Ics20bankTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ics20bank *Ics20bankRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ics20bank.Contract.Ics20bankTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ics20bank *Ics20bankCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ics20bank.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ics20bank *Ics20bankTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ics20bank.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ics20bank *Ics20bankTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ics20bank.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Ics20bank *Ics20bankCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Ics20bank.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Ics20bank *Ics20bankSession) ADMINROLE() ([32]byte, error) {
	return _Ics20bank.Contract.ADMINROLE(&_Ics20bank.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Ics20bank *Ics20bankCallerSession) ADMINROLE() ([32]byte, error) {
	return _Ics20bank.Contract.ADMINROLE(&_Ics20bank.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Ics20bank *Ics20bankCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Ics20bank.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Ics20bank *Ics20bankSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Ics20bank.Contract.DEFAULTADMINROLE(&_Ics20bank.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Ics20bank *Ics20bankCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Ics20bank.Contract.DEFAULTADMINROLE(&_Ics20bank.CallOpts)
}

// OPERATORROLE is a free data retrieval call binding the contract method 0xf5b541a6.
//
// Solidity: function OPERATOR_ROLE() view returns(bytes32)
func (_Ics20bank *Ics20bankCaller) OPERATORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Ics20bank.contract.Call(opts, &out, "OPERATOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OPERATORROLE is a free data retrieval call binding the contract method 0xf5b541a6.
//
// Solidity: function OPERATOR_ROLE() view returns(bytes32)
func (_Ics20bank *Ics20bankSession) OPERATORROLE() ([32]byte, error) {
	return _Ics20bank.Contract.OPERATORROLE(&_Ics20bank.CallOpts)
}

// OPERATORROLE is a free data retrieval call binding the contract method 0xf5b541a6.
//
// Solidity: function OPERATOR_ROLE() view returns(bytes32)
func (_Ics20bank *Ics20bankCallerSession) OPERATORROLE() ([32]byte, error) {
	return _Ics20bank.Contract.OPERATORROLE(&_Ics20bank.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0xb9b092c8.
//
// Solidity: function balanceOf(address account, string id) view returns(uint256)
func (_Ics20bank *Ics20bankCaller) BalanceOf(opts *bind.CallOpts, account common.Address, id string) (*big.Int, error) {
	var out []interface{}
	err := _Ics20bank.contract.Call(opts, &out, "balanceOf", account, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0xb9b092c8.
//
// Solidity: function balanceOf(address account, string id) view returns(uint256)
func (_Ics20bank *Ics20bankSession) BalanceOf(account common.Address, id string) (*big.Int, error) {
	return _Ics20bank.Contract.BalanceOf(&_Ics20bank.CallOpts, account, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0xb9b092c8.
//
// Solidity: function balanceOf(address account, string id) view returns(uint256)
func (_Ics20bank *Ics20bankCallerSession) BalanceOf(account common.Address, id string) (*big.Int, error) {
	return _Ics20bank.Contract.BalanceOf(&_Ics20bank.CallOpts, account, id)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Ics20bank *Ics20bankCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Ics20bank.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Ics20bank *Ics20bankSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Ics20bank.Contract.GetRoleAdmin(&_Ics20bank.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Ics20bank *Ics20bankCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Ics20bank.Contract.GetRoleAdmin(&_Ics20bank.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Ics20bank *Ics20bankCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Ics20bank.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Ics20bank *Ics20bankSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Ics20bank.Contract.HasRole(&_Ics20bank.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Ics20bank *Ics20bankCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Ics20bank.Contract.HasRole(&_Ics20bank.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Ics20bank *Ics20bankCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Ics20bank.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Ics20bank *Ics20bankSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Ics20bank.Contract.SupportsInterface(&_Ics20bank.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Ics20bank *Ics20bankCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Ics20bank.Contract.SupportsInterface(&_Ics20bank.CallOpts, interfaceId)
}

// Burn is a paid mutator transaction binding the contract method 0xc45b71de.
//
// Solidity: function burn(address account, string id, uint256 amount) returns()
func (_Ics20bank *Ics20bankTransactor) Burn(opts *bind.TransactOpts, account common.Address, id string, amount *big.Int) (*types.Transaction, error) {
	return _Ics20bank.contract.Transact(opts, "burn", account, id, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xc45b71de.
//
// Solidity: function burn(address account, string id, uint256 amount) returns()
func (_Ics20bank *Ics20bankSession) Burn(account common.Address, id string, amount *big.Int) (*types.Transaction, error) {
	return _Ics20bank.Contract.Burn(&_Ics20bank.TransactOpts, account, id, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xc45b71de.
//
// Solidity: function burn(address account, string id, uint256 amount) returns()
func (_Ics20bank *Ics20bankTransactorSession) Burn(account common.Address, id string, amount *big.Int) (*types.Transaction, error) {
	return _Ics20bank.Contract.Burn(&_Ics20bank.TransactOpts, account, id, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xf45346dc.
//
// Solidity: function deposit(address tokenContract, uint256 amount, address receiver) returns()
func (_Ics20bank *Ics20bankTransactor) Deposit(opts *bind.TransactOpts, tokenContract common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _Ics20bank.contract.Transact(opts, "deposit", tokenContract, amount, receiver)
}

// Deposit is a paid mutator transaction binding the contract method 0xf45346dc.
//
// Solidity: function deposit(address tokenContract, uint256 amount, address receiver) returns()
func (_Ics20bank *Ics20bankSession) Deposit(tokenContract common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.Deposit(&_Ics20bank.TransactOpts, tokenContract, amount, receiver)
}

// Deposit is a paid mutator transaction binding the contract method 0xf45346dc.
//
// Solidity: function deposit(address tokenContract, uint256 amount, address receiver) returns()
func (_Ics20bank *Ics20bankTransactorSession) Deposit(tokenContract common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.Deposit(&_Ics20bank.TransactOpts, tokenContract, amount, receiver)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Ics20bank *Ics20bankTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ics20bank.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Ics20bank *Ics20bankSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.GrantRole(&_Ics20bank.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Ics20bank *Ics20bankTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.GrantRole(&_Ics20bank.TransactOpts, role, account)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address account, string id, uint256 amount) returns()
func (_Ics20bank *Ics20bankTransactor) Mint(opts *bind.TransactOpts, account common.Address, id string, amount *big.Int) (*types.Transaction, error) {
	return _Ics20bank.contract.Transact(opts, "mint", account, id, amount)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address account, string id, uint256 amount) returns()
func (_Ics20bank *Ics20bankSession) Mint(account common.Address, id string, amount *big.Int) (*types.Transaction, error) {
	return _Ics20bank.Contract.Mint(&_Ics20bank.TransactOpts, account, id, amount)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address account, string id, uint256 amount) returns()
func (_Ics20bank *Ics20bankTransactorSession) Mint(account common.Address, id string, amount *big.Int) (*types.Transaction, error) {
	return _Ics20bank.Contract.Mint(&_Ics20bank.TransactOpts, account, id, amount)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Ics20bank *Ics20bankTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ics20bank.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Ics20bank *Ics20bankSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.RenounceRole(&_Ics20bank.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Ics20bank *Ics20bankTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.RenounceRole(&_Ics20bank.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Ics20bank *Ics20bankTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ics20bank.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Ics20bank *Ics20bankSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.RevokeRole(&_Ics20bank.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Ics20bank *Ics20bankTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.RevokeRole(&_Ics20bank.TransactOpts, role, account)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_Ics20bank *Ics20bankTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _Ics20bank.contract.Transact(opts, "setOperator", operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_Ics20bank *Ics20bankSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.SetOperator(&_Ics20bank.TransactOpts, operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_Ics20bank *Ics20bankTransactorSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.SetOperator(&_Ics20bank.TransactOpts, operator)
}

// TransferFrom is a paid mutator transaction binding the contract method 0xf24dc1da.
//
// Solidity: function transferFrom(address from, address to, string id, uint256 amount) returns()
func (_Ics20bank *Ics20bankTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, id string, amount *big.Int) (*types.Transaction, error) {
	return _Ics20bank.contract.Transact(opts, "transferFrom", from, to, id, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0xf24dc1da.
//
// Solidity: function transferFrom(address from, address to, string id, uint256 amount) returns()
func (_Ics20bank *Ics20bankSession) TransferFrom(from common.Address, to common.Address, id string, amount *big.Int) (*types.Transaction, error) {
	return _Ics20bank.Contract.TransferFrom(&_Ics20bank.TransactOpts, from, to, id, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0xf24dc1da.
//
// Solidity: function transferFrom(address from, address to, string id, uint256 amount) returns()
func (_Ics20bank *Ics20bankTransactorSession) TransferFrom(from common.Address, to common.Address, id string, amount *big.Int) (*types.Transaction, error) {
	return _Ics20bank.Contract.TransferFrom(&_Ics20bank.TransactOpts, from, to, id, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address tokenContract, uint256 amount, address receiver) returns()
func (_Ics20bank *Ics20bankTransactor) Withdraw(opts *bind.TransactOpts, tokenContract common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _Ics20bank.contract.Transact(opts, "withdraw", tokenContract, amount, receiver)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address tokenContract, uint256 amount, address receiver) returns()
func (_Ics20bank *Ics20bankSession) Withdraw(tokenContract common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.Withdraw(&_Ics20bank.TransactOpts, tokenContract, amount, receiver)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address tokenContract, uint256 amount, address receiver) returns()
func (_Ics20bank *Ics20bankTransactorSession) Withdraw(tokenContract common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _Ics20bank.Contract.Withdraw(&_Ics20bank.TransactOpts, tokenContract, amount, receiver)
}

// Ics20bankRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Ics20bank contract.
type Ics20bankRoleAdminChangedIterator struct {
	Event *Ics20bankRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *Ics20bankRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Ics20bankRoleAdminChanged)
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
		it.Event = new(Ics20bankRoleAdminChanged)
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
func (it *Ics20bankRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Ics20bankRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Ics20bankRoleAdminChanged represents a RoleAdminChanged event raised by the Ics20bank contract.
type Ics20bankRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Ics20bank *Ics20bankFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*Ics20bankRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Ics20bank.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &Ics20bankRoleAdminChangedIterator{contract: _Ics20bank.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Ics20bank *Ics20bankFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *Ics20bankRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Ics20bank.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Ics20bankRoleAdminChanged)
				if err := _Ics20bank.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Ics20bank *Ics20bankFilterer) ParseRoleAdminChanged(log types.Log) (*Ics20bankRoleAdminChanged, error) {
	event := new(Ics20bankRoleAdminChanged)
	if err := _Ics20bank.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Ics20bankRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Ics20bank contract.
type Ics20bankRoleGrantedIterator struct {
	Event *Ics20bankRoleGranted // Event containing the contract specifics and raw log

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
func (it *Ics20bankRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Ics20bankRoleGranted)
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
		it.Event = new(Ics20bankRoleGranted)
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
func (it *Ics20bankRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Ics20bankRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Ics20bankRoleGranted represents a RoleGranted event raised by the Ics20bank contract.
type Ics20bankRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ics20bank *Ics20bankFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*Ics20bankRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Ics20bank.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &Ics20bankRoleGrantedIterator{contract: _Ics20bank.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ics20bank *Ics20bankFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *Ics20bankRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Ics20bank.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Ics20bankRoleGranted)
				if err := _Ics20bank.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ics20bank *Ics20bankFilterer) ParseRoleGranted(log types.Log) (*Ics20bankRoleGranted, error) {
	event := new(Ics20bankRoleGranted)
	if err := _Ics20bank.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Ics20bankRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Ics20bank contract.
type Ics20bankRoleRevokedIterator struct {
	Event *Ics20bankRoleRevoked // Event containing the contract specifics and raw log

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
func (it *Ics20bankRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Ics20bankRoleRevoked)
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
		it.Event = new(Ics20bankRoleRevoked)
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
func (it *Ics20bankRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Ics20bankRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Ics20bankRoleRevoked represents a RoleRevoked event raised by the Ics20bank contract.
type Ics20bankRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ics20bank *Ics20bankFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*Ics20bankRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Ics20bank.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &Ics20bankRoleRevokedIterator{contract: _Ics20bank.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ics20bank *Ics20bankFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *Ics20bankRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Ics20bank.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Ics20bankRoleRevoked)
				if err := _Ics20bank.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ics20bank *Ics20bankFilterer) ParseRoleRevoked(log types.Log) (*Ics20bankRoleRevoked, error) {
	event := new(Ics20bankRoleRevoked)
	if err := _Ics20bank.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Ics20bankTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Ics20bank contract.
type Ics20bankTransferIterator struct {
	Event *Ics20bankTransfer // Event containing the contract specifics and raw log

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
func (it *Ics20bankTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Ics20bankTransfer)
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
		it.Event = new(Ics20bankTransfer)
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
func (it *Ics20bankTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Ics20bankTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Ics20bankTransfer represents a Transfer event raised by the Ics20bank contract.
type Ics20bankTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Ics20bank *Ics20bankFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*Ics20bankTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Ics20bank.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Ics20bankTransferIterator{contract: _Ics20bank.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Ics20bank *Ics20bankFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *Ics20bankTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Ics20bank.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Ics20bankTransfer)
				if err := _Ics20bank.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Ics20bank *Ics20bankFilterer) ParseTransfer(log types.Log) (*Ics20bankTransfer, error) {
	event := new(Ics20bankTransfer)
	if err := _Ics20bank.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
