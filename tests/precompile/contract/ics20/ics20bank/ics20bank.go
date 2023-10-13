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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OPERATOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50620000537fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775620000476200005960201b60201c565b6200006160201b60201c565b620001d2565b600033905090565b6200007382826200007760201b60201c565b5050565b6200008982826200016860201b60201c565b6200016457600160008084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550620001096200005960201b60201c565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b600080600084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b6123c180620001e26000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c8063b3ab15fb11610097578063d547741f11610066578063d547741f146102a9578063f24dc1da146102c5578063f45346dc146102e1578063f5b541a6146102fd57610100565b8063b3ab15fb14610225578063b9b092c814610241578063ba7aef4314610271578063c45b71de1461028d57610100565b806369328dec116100d357806369328dec1461019d57806375b238fc146101b957806391d14854146101d7578063a217fddf1461020757610100565b806301ffc9a714610105578063248a9ca3146101355780632f2ff15d1461016557806336568abe14610181575b600080fd5b61011f600480360381019061011a919061188f565b61031b565b60405161012c9190611bfb565b60405180910390f35b61014f600480360381019061014a919061182a565b610395565b60405161015c9190611c16565b60405180910390f35b61017f600480360381019061017a9190611853565b6103b4565b005b61019b60048036038101906101969190611853565b6103d5565b005b6101b760048036038101906101b291906117b2565b610458565b005b6101c1610535565b6040516101ce9190611c16565b60405180910390f35b6101f160048036038101906101ec9190611853565b610559565b6040516101fe9190611bfb565b60405180910390f35b61020f6105c3565b60405161021c9190611c16565b60405180910390f35b61023f600480360381019061023a9190611645565b6105ca565b005b61025b600480360381019061025691906116ee565b610667565b6040516102689190611dd5565b60405180910390f35b61028b60048036038101906102869190611746565b61073f565b005b6102a760048036038101906102a29190611746565b610804565b005b6102c360048036038101906102be9190611853565b6108c9565b005b6102df60048036038101906102da919061166e565b6108ea565b005b6102fb60048036038101906102f691906117b2565b610bfe565b005b610305610cdd565b6040516103129190611c16565b60405180910390f35b60007f7965db0b000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916148061038e575061038d82610d01565b5b9050919050565b6000806000838152602001908152602001600020600101549050919050565b6103bd82610395565b6103c681610d6b565b6103d08383610d7f565b505050565b6103dd610e5f565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461044a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161044190611d95565b60405180910390fd5b6104548282610e67565b5050565b6104778373ffffffffffffffffffffffffffffffffffffffff16610f48565b61048057600080fd5b61049a61048b610e5f565b61049485610f6b565b84610f7d565b8273ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82846040518363ffffffff1660e01b81526004016104d5929190611bd2565b602060405180830381600087803b1580156104ef57600080fd5b505af1158015610503573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105279190611801565b61053057600080fd5b505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b600080600084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b6000801b81565b6105fb7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756105f6610e5f565b610559565b61063a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161063190611d55565b60405180910390fd5b6106647f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b929826110f9565b50565b60008073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1614156106d8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106cf90611d75565b60405180910390fd5b600183836040516106ea929190611b31565b908152602001604051809103902060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490509392505050565b6107707f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b92961076b610e5f565b610559565b6107af576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107a690611db5565b60405180910390fd5b6107fe8484848080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505083611107565b50505050565b6108357f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b929610830610e5f565b610559565b610874576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161086b90611db5565b60405180910390fd5b6108c38484848080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505083610f7d565b50505050565b6108d282610395565b6108db81610d6b565b6108e58383610e67565b505050565b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16141561095a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161095190611cd5565b60405180910390fd5b610962610e5f565b73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1614806109c857506109c77f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b9296109c2610e5f565b610559565b5b610a07576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109fe90611d15565b60405180910390fd5b600060018484604051610a1b929190611b31565b908152602001604051809103902060008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610aac576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610aa390611d35565b60405180910390fd5b8181610ab89190611ec7565b60018585604051610aca929190611b31565b908152602001604051809103902060008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508160018585604051610b2c929190611b31565b908152602001604051809103902060008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610b869190611e17565b925050819055508473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff167f1d30d3db8e01fa0d5626c471596f822f597e720c26a2930ef20d3387313c3d78868686604051610bee93929190611c31565b60405180910390a3505050505050565b610c1d8373ffffffffffffffffffffffffffffffffffffffff16610f48565b610c2657600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd610c4a610e5f565b30856040518463ffffffff1660e01b8152600401610c6a93929190611b9b565b602060405180830381600087803b158015610c8457600080fd5b505af1158015610c98573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cbc9190611801565b610cc557600080fd5b610cd881610cd285610f6b565b84611107565b505050565b7f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b92981565b60007f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b610d7c81610d77610e5f565b6111e6565b50565b610d898282610559565b610e5b57600160008084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550610e00610e5f565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b600033905090565b610e718282610559565b15610f4457600080600084815260200190815260200160002060000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550610ee9610e5f565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b60405160405180910390a45b5050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b6060610f768261126b565b9050919050565b6000600183604051610f8f9190611b4a565b908152602001604051809103902060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015611020576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161101790611cf5565b60405180910390fd5b818161102c9190611ec7565b60018460405161103c9190611b4a565b908152602001604051809103902060008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f1d30d3db8e01fa0d5626c471596f822f597e720c26a2930ef20d3387313c3d7885856040516110eb929190611c85565b60405180910390a350505050565b6111038282610d7f565b5050565b806001836040516111189190611b4a565b908152602001604051809103902060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546111729190611e17565b925050819055508273ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f1d30d3db8e01fa0d5626c471596f822f597e720c26a2930ef20d3387313c3d7884846040516111d9929190611c85565b60405180910390a3505050565b6111f08282610559565b611267576111fd8161126b565b61120b8360001c6020611298565b60405160200161121c929190611b61565b6040516020818303038152906040526040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161125e9190611c63565b60405180910390fd5b5050565b60606112918273ffffffffffffffffffffffffffffffffffffffff16601460ff16611298565b9050919050565b6060600060028360026112ab9190611e6d565b6112b59190611e17565b67ffffffffffffffff8111156112f4577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519080825280601f01601f1916602001820160405280156113265781602001600182028036833780820191505090505b5090507f300000000000000000000000000000000000000000000000000000000000000081600081518110611384577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053507f78000000000000000000000000000000000000000000000000000000000000008160018151811061140e577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506000600184600261144e9190611e6d565b6114589190611e17565b90505b6001811115611544577f3031323334353637383961626364656600000000000000000000000000000000600f8616601081106114c0577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b1a60f81b8282815181106114fd577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600485901c94508061153d90611fbb565b905061145b565b5060008414611588576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161157f90611cb5565b60405180910390fd5b8091505092915050565b6000813590506115a181612318565b92915050565b6000815190506115b68161232f565b92915050565b6000813590506115cb81612346565b92915050565b6000813590506115e08161235d565b92915050565b60008083601f8401126115f857600080fd5b8235905067ffffffffffffffff81111561161157600080fd5b60208301915083600182028301111561162957600080fd5b9250929050565b60008135905061163f81612374565b92915050565b60006020828403121561165757600080fd5b600061166584828501611592565b91505092915050565b60008060008060006080868803121561168657600080fd5b600061169488828901611592565b95505060206116a588828901611592565b945050604086013567ffffffffffffffff8111156116c257600080fd5b6116ce888289016115e6565b935093505060606116e188828901611630565b9150509295509295909350565b60008060006040848603121561170357600080fd5b600061171186828701611592565b935050602084013567ffffffffffffffff81111561172e57600080fd5b61173a868287016115e6565b92509250509250925092565b6000806000806060858703121561175c57600080fd5b600061176a87828801611592565b945050602085013567ffffffffffffffff81111561178757600080fd5b611793878288016115e6565b935093505060406117a687828801611630565b91505092959194509250565b6000806000606084860312156117c757600080fd5b60006117d586828701611592565b93505060206117e686828701611630565b92505060406117f786828701611592565b9150509250925092565b60006020828403121561181357600080fd5b6000611821848285016115a7565b91505092915050565b60006020828403121561183c57600080fd5b600061184a848285016115bc565b91505092915050565b6000806040838503121561186657600080fd5b6000611874858286016115bc565b925050602061188585828601611592565b9150509250929050565b6000602082840312156118a157600080fd5b60006118af848285016115d1565b91505092915050565b6118c181611efb565b82525050565b6118d081611f0d565b82525050565b6118df81611f19565b82525050565b60006118f18385611dfb565b93506118fe838584611f79565b61190783612014565b840190509392505050565b600061191e8385611e0c565b935061192b838584611f79565b82840190509392505050565b600061194282611df0565b61194c8185611dfb565b935061195c818560208601611f88565b61196581612014565b840191505092915050565b600061197b82611df0565b6119858185611e0c565b9350611995818560208601611f88565b80840191505092915050565b60006119ae602083611dfb565b91506119b982612025565b602082019050919050565b60006119d1602783611dfb565b91506119dc8261204e565b604082019050919050565b60006119f4602683611dfb565b91506119ff8261209d565b604082019050919050565b6000611a17602b83611dfb565b9150611a22826120ec565b604082019050919050565b6000611a3a602c83611dfb565b9150611a458261213b565b604082019050919050565b6000611a5d602883611dfb565b9150611a688261218a565b604082019050919050565b6000611a80601783611e0c565b9150611a8b826121d9565b601782019050919050565b6000611aa3602d83611dfb565b9150611aae82612202565b604082019050919050565b6000611ac6601183611e0c565b9150611ad182612251565b601182019050919050565b6000611ae9602f83611dfb565b9150611af48261227a565b604082019050919050565b6000611b0c602883611dfb565b9150611b17826122c9565b604082019050919050565b611b2b81611f6f565b82525050565b6000611b3e828486611912565b91508190509392505050565b6000611b568284611970565b915081905092915050565b6000611b6c82611a73565b9150611b788285611970565b9150611b8382611ab9565b9150611b8f8284611970565b91508190509392505050565b6000606082019050611bb060008301866118b8565b611bbd60208301856118b8565b611bca6040830184611b22565b949350505050565b6000604082019050611be760008301856118b8565b611bf46020830184611b22565b9392505050565b6000602082019050611c1060008301846118c7565b92915050565b6000602082019050611c2b60008301846118d6565b92915050565b60006040820190508181036000830152611c4c8185876118e5565b9050611c5b6020830184611b22565b949350505050565b60006020820190508181036000830152611c7d8184611937565b905092915050565b60006040820190508181036000830152611c9f8185611937565b9050611cae6020830184611b22565b9392505050565b60006020820190508181036000830152611cce816119a1565b9050919050565b60006020820190508181036000830152611cee816119c4565b9050919050565b60006020820190508181036000830152611d0e816119e7565b9050919050565b60006020820190508181036000830152611d2e81611a0a565b9050919050565b60006020820190508181036000830152611d4e81611a2d565b9050919050565b60006020820190508181036000830152611d6e81611a50565b9050919050565b60006020820190508181036000830152611d8e81611a96565b9050919050565b60006020820190508181036000830152611dae81611adc565b9050919050565b60006020820190508181036000830152611dce81611aff565b9050919050565b6000602082019050611dea6000830184611b22565b92915050565b600081519050919050565b600082825260208201905092915050565b600081905092915050565b6000611e2282611f6f565b9150611e2d83611f6f565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611e6257611e61611fe5565b5b828201905092915050565b6000611e7882611f6f565b9150611e8383611f6f565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611ebc57611ebb611fe5565b5b828202905092915050565b6000611ed282611f6f565b9150611edd83611f6f565b925082821015611ef057611eef611fe5565b5b828203905092915050565b6000611f0682611f4f565b9050919050565b60008115159050919050565b6000819050919050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b83811015611fa6578082015181840152602081019050611f8b565b83811115611fb5576000848401525b50505050565b6000611fc682611f6f565b91506000821415611fda57611fd9611fe5565b5b600182039050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000601f19601f8301169050919050565b7f537472696e67733a20686578206c656e67746820696e73756666696369656e74600082015250565b7f494353323042616e6b3a207472616e7366657220746f20746865207a65726f2060008201527f6164647265737300000000000000000000000000000000000000000000000000602082015250565b7f494353323042616e6b3a206275726e20616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b7f494353323042616e6b3a2063616c6c6572206973206e6f74206f776e6572206e60008201527f6f7220617070726f766564000000000000000000000000000000000000000000602082015250565b7f494353323042616e6b3a20696e73756666696369656e742062616c616e63652060008201527f666f72207472616e736665720000000000000000000000000000000000000000602082015250565b7f6d75737420686176652061646d696e20726f6c6520746f20736574206e65772060008201527f6f70657261746f72000000000000000000000000000000000000000000000000602082015250565b7f416363657373436f6e74726f6c3a206163636f756e7420000000000000000000600082015250565b7f494353323042616e6b3a2062616c616e636520717565727920666f722074686560008201527f207a65726f206164647265737300000000000000000000000000000000000000602082015250565b7f206973206d697373696e6720726f6c6520000000000000000000000000000000600082015250565b7f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560008201527f20726f6c657320666f722073656c660000000000000000000000000000000000602082015250565b7f494353323042616e6b3a206d7573742068617665206d696e74657220726f6c6560008201527f20746f206d696e74000000000000000000000000000000000000000000000000602082015250565b61232181611efb565b811461232c57600080fd5b50565b61233881611f0d565b811461234357600080fd5b50565b61234f81611f19565b811461235a57600080fd5b50565b61236681611f23565b811461237157600080fd5b50565b61237d81611f6f565b811461238857600080fd5b5056fea264697066735822122008154672f506e59a032c577d7fb44eede5e11dc06b33f174acc319b88b59945c64736f6c63430008010033",
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
	Path  string
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0x1d30d3db8e01fa0d5626c471596f822f597e720c26a2930ef20d3387313c3d78.
//
// Solidity: event Transfer(address indexed from, address indexed to, string path, uint256 value)
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

// WatchTransfer is a free log subscription operation binding the contract event 0x1d30d3db8e01fa0d5626c471596f822f597e720c26a2930ef20d3387313c3d78.
//
// Solidity: event Transfer(address indexed from, address indexed to, string path, uint256 value)
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

// ParseTransfer is a log parse operation binding the contract event 0x1d30d3db8e01fa0d5626c471596f822f597e720c26a2930ef20d3387313c3d78.
//
// Solidity: event Transfer(address indexed from, address indexed to, string path, uint256 value)
func (_Ics20bank *Ics20bankFilterer) ParseTransfer(log types.Log) (*Ics20bankTransfer, error) {
	event := new(Ics20bankTransfer)
	if err := _Ics20bank.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
