package utils

import (
	"context"
	"math/big"
	"time"

	"github.com/avast/retry-go"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/onsi/gomega"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/precompile/contracts/ibc"
	contractBind "github.com/ava-labs/subnet-evm/tests/precompile/contract"
)

var (
	testKey, _ = crypto.HexToECDSA("56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027")
	chainId    = big.NewInt(99999)
)

func RunIBCTests(ctx context.Context) {
	log.Info("Executing IBC tests on a new blockchain")

	genesisFilePath := "./tests/precompile/genesis/ibc.json"

	blockchainID := CreateNewSubnet(ctx, genesisFilePath)
	chainURI := GetDefaultChainURI(blockchainID)
	log.Info("Created subnet successfully", "ChainURI", chainURI)

	ethClient, err := ethclient.DialContext(ctx, chainURI)
	gomega.Expect(err).Should(gomega.BeNil())

	ibcContract, err := contractBind.NewContract(ibc.ContractAddress, ethClient)
	gomega.Expect(err).Should(gomega.BeNil())

	auth, err := bind.NewKeyedTransactorWithChainID(testKey, chainId)
	gomega.Expect(err).Should(gomega.BeNil())

	log.Info("Starting IBC tests")

	// create tendermint client
	tx, receipt, clientId := createClient(ctx, ethClient, ibcContract, auth)

	log.Info("Create Client tx", "hash", tx.Hash(), "blockNumber", receipt.BlockNumber.Uint64(), "clientId", clientId)
}

func createClient(ctx context.Context, ethClient ethclient.Client, contract *contractBind.Contract, auth *bind.TransactOpts) (*types.Transaction, *types.Receipt, string) {
	tx, err := contract.CreateClient(auth, exported.Tendermint, []byte{123}, []byte{21})
	gomega.Expect(err).Should(gomega.BeNil())

	receipt, err := waitForReceiptAndGet(ctx, ethClient, tx)
	gomega.Expect(err).Should(gomega.BeNil())

	contractFilterer, err := contractBind.NewContractFilterer(ibc.ContractAddress, ethClient)
	gomega.Expect(err).Should(gomega.BeNil())

	var clientCreatedLog types.Log
	if len(receipt.Logs) > 0 {
		clientCreatedLog = *receipt.Logs[0]
	}
	clientCreatedEvent, err := contractFilterer.ParseClientCreated(clientCreatedLog)
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(clientCreatedEvent.ClientId).Should(gomega.Equal("07-tendermint-0"))

	return tx, receipt, clientCreatedEvent.ClientId
}

func waitForReceiptAndGet(ctx context.Context, client ethclient.Client, tx *types.Transaction) (*types.Receipt, error) {
	var receipt *types.Receipt
	err := retry.Do(
		func() error {
			rc, err := client.TransactionReceipt(ctx, tx.Hash())
			if err != nil {
				return err
			}
			receipt = rc
			return nil
		},
		retry.Delay(1*time.Second),
		retry.Attempts(10),
	)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
