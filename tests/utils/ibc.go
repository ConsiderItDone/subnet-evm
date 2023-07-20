package utils

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/onsi/gomega"

	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/precompile/contracts/ibc"
	"github.com/ava-labs/subnet-evm/tests/precompile/contract"
)

func RunIBCTests(ctx context.Context) {
	log.Info("Executing IBC tests on a new blockchain")

	genesisFilePath := fmt.Sprintf("./tests/precompile/genesis/ibc.json")

	blockchainID := CreateNewSubnet(ctx, genesisFilePath)
	chainURI := GetDefaultChainURI(blockchainID)
	log.Info("Created subnet successfully", "ChainURI", chainURI)

	client, err := ethclient.Dial(chainURI)
	gomega.Expect(err).Should(gomega.BeNil())

	ibcContract, err := contract.NewContract(ibc.ContractAddress, client)
	gomega.Expect(err).Should(gomega.BeNil())

	ibcContract.CreateClient()

	//cmdPath := "./contracts"
	//// test path is relative to the cmd path
	//testPath := fmt.Sprintf("./test/%s.ts", test)
	//cmd := exec.Command("npx", "hardhat", "test", testPath, "--network", "local")
	//cmd.Dir = cmdPath
	//
	//RunTestCMD(cmd, chainURI)
}
