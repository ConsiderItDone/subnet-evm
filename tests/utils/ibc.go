package utils

import (
	"context"
	"math/big"
	"time"

	"github.com/avast/retry-go"
	tmtypes "github.com/cometbft/cometbft/types"
	tmClientTypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	ibctestingmock "github.com/cosmos/ibc-go/v7/testing/mock"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/onsi/gomega"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/precompile/contracts/ibc"
	contractBind "github.com/ava-labs/subnet-evm/tests/precompile/contract"
)

const (
	testChainID = "gaiahub-0"

	trustingPeriod time.Duration = time.Hour * 24 * 7 * 2
	ubdPeriod      time.Duration = time.Hour * 24 * 7 * 3
	maxClockDrift  time.Duration = time.Second * 10
)

var (
	testKey, _ = crypto.HexToECDSA("56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027")
	chainId    = big.NewInt(99999)

	testClientHeight = tmClientTypes.NewHeight(0, 5)
)

func initIBCTest(ctx context.Context) (ethclient.Client, *contractBind.Contract, *bind.TransactOpts) {
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

	return ethClient, ibcContract, auth
}

func RunIBCCreateClientTests(ctx context.Context) {
	ethClient, ibcContract, auth := initIBCTest(ctx)

	log.Info("Starting IBC tests")

	// create tendermint client
	createClientTx, receipt, clientId := createClient(ctx, ethClient, ibcContract, auth)
	gomega.Expect(clientId).Should(gomega.Equal("07-tendermint-0"))

	log.Info("Create Client tx", "hash", createClientTx.Hash(), "blockNumber", receipt.BlockNumber.Uint64(), "clientId", clientId)

	createClientTx2, receipt2, clientId2 := createClient(ctx, ethClient, ibcContract, auth)
	gomega.Expect(clientId2).Should(gomega.Equal("07-tendermint-1"))

	log.Info("Create Client tx2", "hash", createClientTx2.Hash(), "blockNumber", receipt2.BlockNumber.Uint64(), "clientId", clientId2)
}

func RunIBCConnectionOpenInitTests(ctx context.Context) {
	ethClient, ibcContract, auth := initIBCTest(ctx)

	log.Info("Starting IBC tests")

	// create tendermint client
	createClientTx, receipt, clientId := createClient(ctx, ethClient, ibcContract, auth)
	gomega.Expect(clientId).Should(gomega.Equal("07-tendermint-0"))

	log.Info("Create Client tx", "hash", createClientTx.Hash(), "blockNumber", receipt.BlockNumber.Uint64(), "clientId", clientId)

	// connection open init
	connOpenInitTx, connOpenInitReceipt, connectionId, clientId2 := connectionOpenInit(ctx, ethClient, ibcContract, auth, clientId)
	gomega.Expect(clientId).Should(gomega.Equal(clientId2))

	log.Info("ConnOpenInit tx", "hash", connOpenInitTx.Hash(), "blockNumber", connOpenInitReceipt.BlockNumber.Uint64(), "connectionId", connectionId)
}

func RunIBCConnectionOpenTryTests(ctx context.Context) {
	ethClient, ibcContract, auth := initIBCTest(ctx)

	log.Info("Starting IBC tests")

	// create tendermint client
	createClientTx, receipt, clientId := createClient(ctx, ethClient, ibcContract, auth)
	gomega.Expect(clientId).Should(gomega.Equal("07-tendermint-0"))

	log.Info("Create Client tx", "hash", createClientTx.Hash(), "blockNumber", receipt.BlockNumber.Uint64(), "clientId", clientId)

	// connection open init
	connOpenInitTx, connOpenInitReceipt, connectionId, clientId2 := connectionOpenInit(ctx, ethClient, ibcContract, auth, clientId)
	gomega.Expect(clientId).Should(gomega.Equal(clientId2))

	log.Info("ConnOpenInit tx", "hash", connOpenInitTx.Hash(), "blockNumber", connOpenInitReceipt.BlockNumber.Uint64(), "connectionId", connectionId)

	// connection open try
	connectionOpenTry(ctx, ethClient, ibcContract, auth)
}

func createClient(ctx context.Context, ethClient ethclient.Client, contract *contractBind.Contract, auth *bind.TransactOpts) (*types.Transaction, *types.Receipt, string) {
	clientStateBytes := getClientStateBytes()
	consensusStateBytes := getConsensusStateBytes()

	tx, err := contract.CreateClient(auth, exported.Tendermint, clientStateBytes, consensusStateBytes)
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

	return tx, receipt, clientCreatedEvent.ClientId
}

func connectionOpenInit(ctx context.Context, ethClient ethclient.Client, contract *contractBind.Contract, auth *bind.TransactOpts, clientId string) (*types.Transaction, *types.Receipt, string, string) {
	counterpartyBytes := getCounterpartyBytes()

	version := connectiontypes.ExportedVersionsToProto(connectiontypes.GetCompatibleVersions())[0]
	versionBytes, err := version.Marshal()
	gomega.Expect(err).Should(gomega.BeNil())

	delayPeriod := uint32(time.Hour.Nanoseconds())

	tx, err := contract.ConnOpenInit(auth, clientId, counterpartyBytes, versionBytes, delayPeriod)
	gomega.Expect(err).Should(gomega.BeNil())

	receipt, err := waitForReceiptAndGet(ctx, ethClient, tx)
	gomega.Expect(err).Should(gomega.BeNil())

	contractFilterer, err := contractBind.NewContractFilterer(ibc.ContractAddress, ethClient)
	gomega.Expect(err).Should(gomega.BeNil())

	var connectionCreatedLog types.Log
	if len(receipt.Logs) > 0 {
		connectionCreatedLog = *receipt.Logs[0]
	}
	connectionCreatedEvent, err := contractFilterer.ParseConnectionCreated(connectionCreatedLog)
	gomega.Expect(err).Should(gomega.BeNil())

	// TODO: query proofs for client state, consensus state and connection

	return tx, receipt, connectionCreatedEvent.ConnectionId, connectionCreatedEvent.ClientId
}

func connectionOpenTry(ctx context.Context, ethClient ethclient.Client, contract *contractBind.Contract, auth *bind.TransactOpts) {
	counterpartyBytes := []byte{0xa, 0xf, 0x30, 0x37, 0x2d, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x74, 0x2d, 0x30, 0x12, 0xc, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x30, 0x1a, 0x5, 0xa, 0x3, 0x69, 0x62, 0x63}
	delayPeriod := uint32(time.Hour.Nanoseconds())
	clientStateBytes := []byte{0xa, 0x2b, 0x2f, 0x69, 0x62, 0x63, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x7d, 0xa, 0xc, 0x74, 0x65, 0x73, 0x74, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x32, 0x2d, 0x31, 0x12, 0x4, 0x8, 0x1, 0x10, 0x3, 0x1a, 0x4, 0x8, 0x80, 0xea, 0x49, 0x22, 0x4, 0x8, 0x80, 0xdf, 0x6e, 0x2a, 0x2, 0x8, 0xa, 0x32, 0x0, 0x3a, 0x4, 0x8, 0x1, 0x10, 0x3, 0x42, 0x19, 0xa, 0x9, 0x8, 0x1, 0x18, 0x1, 0x20, 0x1, 0x2a, 0x1, 0x0, 0x12, 0xc, 0xa, 0x2, 0x0, 0x1, 0x10, 0x21, 0x18, 0x4, 0x20, 0xc, 0x30, 0x1, 0x42, 0x19, 0xa, 0x9, 0x8, 0x1, 0x18, 0x1, 0x20, 0x1, 0x2a, 0x1, 0x0, 0x12, 0xc, 0xa, 0x2, 0x0, 0x1, 0x10, 0x20, 0x18, 0x1, 0x20, 0x1, 0x30, 0x1, 0x4a, 0x7, 0x75, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x4a, 0x10, 0x75, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x64, 0x49, 0x42, 0x43, 0x53, 0x74, 0x61, 0x74, 0x65}
	counterpartyVersions := []byte{0x5b, 0x7b, 0x22, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x2c, 0x22, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x22, 0x3a, 0x5b, 0x22, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x45, 0x44, 0x22, 0x2c, 0x22, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x55, 0x4e, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x45, 0x44, 0x22, 0x5d, 0x7d, 0x5d}

	proofInit := []byte{0xa, 0xac, 0x2, 0xa, 0xa9, 0x2, 0xa, 0x18, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x30, 0x12, 0x52, 0xa, 0xf, 0x30, 0x37, 0x2d, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x74, 0x2d, 0x30, 0x12, 0x23, 0xa, 0x1, 0x31, 0x12, 0xd, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x45, 0x44, 0x12, 0xf, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x55, 0x4e, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x45, 0x44, 0x18, 0x1, 0x22, 0x18, 0xa, 0xf, 0x30, 0x37, 0x2d, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x74, 0x2d, 0x30, 0x1a, 0x5, 0xa, 0x3, 0x69, 0x62, 0x63, 0x1a, 0xb, 0x8, 0x1, 0x18, 0x1, 0x20, 0x1, 0x2a, 0x3, 0x0, 0x2, 0xa, 0x22, 0x29, 0x8, 0x1, 0x12, 0x25, 0x2, 0x4, 0xa, 0x20, 0xb9, 0xa, 0x68, 0x2b, 0x88, 0x3c, 0xe7, 0xc2, 0xdc, 0x16, 0xa9, 0x14, 0x3e, 0xe4, 0xd5, 0xe4, 0xea, 0x69, 0xe9, 0x83, 0xbb, 0xd2, 0x8a, 0x1a, 0xab, 0x9, 0x7e, 0x93, 0x42, 0x57, 0xd0, 0x6d, 0x20, 0x22, 0x29, 0x8, 0x1, 0x12, 0x25, 0x4, 0x8, 0xa, 0x20, 0x42, 0x37, 0xc3, 0xbd, 0xc0, 0x20, 0x80, 0x78, 0xf3, 0x6b, 0xe5, 0xd9, 0xc0, 0xf, 0x47, 0x8c, 0x2d, 0xa9, 0x50, 0xe9, 0xca, 0x7e, 0x9f, 0x19, 0xc9, 0xc5, 0x2a, 0xfe, 0xc, 0x4a, 0xfb, 0x66, 0x20, 0x22, 0x29, 0x8, 0x1, 0x12, 0x25, 0x6, 0xe, 0xa, 0x20, 0xa2, 0x64, 0x6d, 0x25, 0x3f, 0x3a, 0xc5, 0xaf, 0xdc, 0x8c, 0x41, 0xb8, 0x62, 0x4e, 0x6b, 0xc0, 0x18, 0x2e, 0x43, 0x69, 0x3f, 0x42, 0x96, 0x68, 0xec, 0x6e, 0x66, 0x98, 0xc9, 0x60, 0x6b, 0x27, 0x20, 0x22, 0x2b, 0x8, 0x1, 0x12, 0x4, 0x8, 0x14, 0xa, 0x20, 0x1a, 0x21, 0x20, 0x9c, 0x2a, 0xf9, 0xb2, 0xb4, 0xe4, 0x4, 0xdb, 0x69, 0x90, 0x5d, 0xef, 0x9a, 0xd6, 0xc, 0xda, 0xfe, 0x86, 0x4c, 0x42, 0xcc, 0xb, 0x4, 0x9a, 0x16, 0x2e, 0xe5, 0x1f, 0x37, 0xa6, 0xe9, 0x31, 0xa, 0xfe, 0x1, 0xa, 0xfb, 0x1, 0xa, 0x3, 0x69, 0x62, 0x63, 0x12, 0x20, 0x6, 0x7, 0x6b, 0x34, 0xdd, 0x5f, 0x86, 0x82, 0x6f, 0xe8, 0x5e, 0x64, 0x43, 0x81, 0x63, 0x36, 0x63, 0x43, 0xba, 0x3f, 0x7b, 0x5e, 0x52, 0xf0, 0x26, 0x3e, 0x23, 0x23, 0xf7, 0xdb, 0x63, 0x79, 0x1a, 0x9, 0x8, 0x1, 0x18, 0x1, 0x20, 0x1, 0x2a, 0x1, 0x0, 0x22, 0x27, 0x8, 0x1, 0x12, 0x1, 0x1, 0x1a, 0x20, 0xc5, 0x5, 0xc0, 0xfd, 0x48, 0xb1, 0xcf, 0x2b, 0x65, 0x61, 0x9f, 0x12, 0xb3, 0x14, 0x4e, 0x19, 0xc6, 0xb, 0x4e, 0x6b, 0x62, 0x52, 0x5e, 0xe9, 0x36, 0xbe, 0x93, 0xd0, 0x41, 0x62, 0x3e, 0xbf, 0x22, 0x27, 0x8, 0x1, 0x12, 0x1, 0x1, 0x1a, 0x20, 0xb5, 0x94, 0xb4, 0xb9, 0x92, 0x43, 0xce, 0xda, 0xa5, 0x46, 0x37, 0x7d, 0xe1, 0xe, 0xbb, 0xc9, 0xfd, 0x57, 0xe3, 0xba, 0x3b, 0x76, 0x4, 0x82, 0x89, 0x2f, 0xa6, 0xc6, 0xd3, 0x4c, 0x82, 0x6a, 0x22, 0x25, 0x8, 0x1, 0x12, 0x21, 0x1, 0x47, 0xe, 0x2a, 0x73, 0x10, 0x14, 0xa6, 0x24, 0xf2, 0x68, 0xdc, 0x7b, 0xc9, 0x84, 0x2d, 0xc5, 0xb1, 0x65, 0x5, 0x9e, 0xb8, 0xc7, 0xe0, 0xb1, 0x5c, 0x73, 0x4e, 0x1d, 0x1b, 0x42, 0x53, 0xe4, 0x22, 0x25, 0x8, 0x1, 0x12, 0x21, 0x1, 0xb6, 0x5f, 0x9a, 0xae, 0x6c, 0x59, 0x3, 0xff, 0xf1, 0x9, 0x87, 0xcd, 0x1, 0xcf, 0x9, 0xd9, 0xe5, 0xbb, 0x9, 0x20, 0x80, 0x9c, 0xe2, 0x83, 0xee, 0x94, 0xc4, 0x8d, 0x5e, 0x11, 0xf7, 0xc4, 0x22, 0x27, 0x8, 0x1, 0x12, 0x1, 0x1, 0x1a, 0x20, 0xf8, 0xa3, 0xd7, 0x24, 0x32, 0x64, 0x99, 0x96, 0x6a, 0xdd, 0xc4, 0xc5, 0x58, 0x57, 0x51, 0xae, 0x96, 0x40, 0x2a, 0xca, 0xb6, 0x1a, 0xeb, 0x74, 0xd3, 0x54, 0x2a, 0xb7, 0x82, 0x7, 0x50, 0x38}
	proofConsensus := []byte{0xa, 0xc5, 0x2, 0xa, 0xc2, 0x2, 0xa, 0x2b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x30, 0x37, 0x2d, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x74, 0x2d, 0x30, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x2f, 0x31, 0x2d, 0x33, 0x12, 0x80, 0x1, 0xa, 0x2e, 0x2f, 0x69, 0x62, 0x63, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x4e, 0xa, 0x6, 0x8, 0x85, 0xe5, 0xb4, 0xf0, 0x5, 0x12, 0x22, 0xa, 0x20, 0x8b, 0xf5, 0x39, 0xdc, 0xca, 0x6a, 0x3a, 0xaa, 0xfc, 0x9b, 0x12, 0x6c, 0xde, 0x59, 0x28, 0xe4, 0xdb, 0x8c, 0xff, 0xc2, 0x3c, 0x49, 0xce, 0x55, 0xb0, 0x6d, 0x29, 0x3f, 0xf9, 0x8f, 0x44, 0xa4, 0x1a, 0x20, 0xda, 0x9a, 0x13, 0x57, 0x71, 0xd7, 0x1f, 0xf9, 0x68, 0x85, 0x62, 0x78, 0xd6, 0xe0, 0x70, 0xf6, 0x27, 0xc7, 0x37, 0xcb, 0x5a, 0xbc, 0x4b, 0x1c, 0xd2, 0x39, 0xbc, 0x68, 0x9d, 0xb5, 0xcf, 0x32, 0x1a, 0xb, 0x8, 0x1, 0x18, 0x1, 0x20, 0x1, 0x2a, 0x3, 0x0, 0x2, 0x6, 0x22, 0x29, 0x8, 0x1, 0x12, 0x25, 0x4, 0x6, 0xa, 0x20, 0x3e, 0xd9, 0x94, 0x3e, 0xc9, 0x98, 0x80, 0x4b, 0x36, 0xad, 0xe3, 0x46, 0x34, 0xfd, 0xe6, 0xa4, 0xea, 0xb4, 0xb, 0xc0, 0xd1, 0xd1, 0x98, 0x3f, 0x8, 0xf3, 0xc7, 0x91, 0xf5, 0x90, 0x37, 0x49, 0x20, 0x22, 0x2b, 0x8, 0x1, 0x12, 0x4, 0x6, 0xe, 0xa, 0x20, 0x1a, 0x21, 0x20, 0x43, 0x7f, 0x23, 0x36, 0x96, 0xb5, 0xe5, 0x75, 0x22, 0xf, 0xce, 0x90, 0x94, 0xae, 0xd6, 0x1d, 0xb5, 0x8f, 0xd4, 0x91, 0xbb, 0x95, 0x85, 0xa0, 0x1e, 0x5d, 0x3d, 0x78, 0xcf, 0xa2, 0xcd, 0x6b, 0x22, 0x2b, 0x8, 0x1, 0x12, 0x4, 0x8, 0x14, 0xa, 0x20, 0x1a, 0x21, 0x20, 0x9c, 0x2a, 0xf9, 0xb2, 0xb4, 0xe4, 0x4, 0xdb, 0x69, 0x90, 0x5d, 0xef, 0x9a, 0xd6, 0xc, 0xda, 0xfe, 0x86, 0x4c, 0x42, 0xcc, 0xb, 0x4, 0x9a, 0x16, 0x2e, 0xe5, 0x1f, 0x37, 0xa6, 0xe9, 0x31, 0xa, 0xfe, 0x1, 0xa, 0xfb, 0x1, 0xa, 0x3, 0x69, 0x62, 0x63, 0x12, 0x20, 0x6, 0x7, 0x6b, 0x34, 0xdd, 0x5f, 0x86, 0x82, 0x6f, 0xe8, 0x5e, 0x64, 0x43, 0x81, 0x63, 0x36, 0x63, 0x43, 0xba, 0x3f, 0x7b, 0x5e, 0x52, 0xf0, 0x26, 0x3e, 0x23, 0x23, 0xf7, 0xdb, 0x63, 0x79, 0x1a, 0x9, 0x8, 0x1, 0x18, 0x1, 0x20, 0x1, 0x2a, 0x1, 0x0, 0x22, 0x27, 0x8, 0x1, 0x12, 0x1, 0x1, 0x1a, 0x20, 0xc5, 0x5, 0xc0, 0xfd, 0x48, 0xb1, 0xcf, 0x2b, 0x65, 0x61, 0x9f, 0x12, 0xb3, 0x14, 0x4e, 0x19, 0xc6, 0xb, 0x4e, 0x6b, 0x62, 0x52, 0x5e, 0xe9, 0x36, 0xbe, 0x93, 0xd0, 0x41, 0x62, 0x3e, 0xbf, 0x22, 0x27, 0x8, 0x1, 0x12, 0x1, 0x1, 0x1a, 0x20, 0xb5, 0x94, 0xb4, 0xb9, 0x92, 0x43, 0xce, 0xda, 0xa5, 0x46, 0x37, 0x7d, 0xe1, 0xe, 0xbb, 0xc9, 0xfd, 0x57, 0xe3, 0xba, 0x3b, 0x76, 0x4, 0x82, 0x89, 0x2f, 0xa6, 0xc6, 0xd3, 0x4c, 0x82, 0x6a, 0x22, 0x25, 0x8, 0x1, 0x12, 0x21, 0x1, 0x47, 0xe, 0x2a, 0x73, 0x10, 0x14, 0xa6, 0x24, 0xf2, 0x68, 0xdc, 0x7b, 0xc9, 0x84, 0x2d, 0xc5, 0xb1, 0x65, 0x5, 0x9e, 0xb8, 0xc7, 0xe0, 0xb1, 0x5c, 0x73, 0x4e, 0x1d, 0x1b, 0x42, 0x53, 0xe4, 0x22, 0x25, 0x8, 0x1, 0x12, 0x21, 0x1, 0xb6, 0x5f, 0x9a, 0xae, 0x6c, 0x59, 0x3, 0xff, 0xf1, 0x9, 0x87, 0xcd, 0x1, 0xcf, 0x9, 0xd9, 0xe5, 0xbb, 0x9, 0x20, 0x80, 0x9c, 0xe2, 0x83, 0xee, 0x94, 0xc4, 0x8d, 0x5e, 0x11, 0xf7, 0xc4, 0x22, 0x27, 0x8, 0x1, 0x12, 0x1, 0x1, 0x1a, 0x20, 0xf8, 0xa3, 0xd7, 0x24, 0x32, 0x64, 0x99, 0x96, 0x6a, 0xdd, 0xc4, 0xc5, 0x58, 0x57, 0x51, 0xae, 0x96, 0x40, 0x2a, 0xca, 0xb6, 0x1a, 0xeb, 0x74, 0xd3, 0x54, 0x2a, 0xb7, 0x82, 0x7, 0x50, 0x38}
	proofClient := []byte{0xa, 0x98, 0x3, 0xa, 0x95, 0x3, 0xa, 0x23, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x30, 0x37, 0x2d, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x74, 0x2d, 0x30, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0xac, 0x1, 0xa, 0x2b, 0x2f, 0x69, 0x62, 0x63, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x7d, 0xa, 0xc, 0x74, 0x65, 0x73, 0x74, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x32, 0x2d, 0x31, 0x12, 0x4, 0x8, 0x1, 0x10, 0x3, 0x1a, 0x4, 0x8, 0x80, 0xea, 0x49, 0x22, 0x4, 0x8, 0x80, 0xdf, 0x6e, 0x2a, 0x2, 0x8, 0xa, 0x32, 0x0, 0x3a, 0x4, 0x8, 0x1, 0x10, 0x3, 0x42, 0x19, 0xa, 0x9, 0x8, 0x1, 0x18, 0x1, 0x20, 0x1, 0x2a, 0x1, 0x0, 0x12, 0xc, 0xa, 0x2, 0x0, 0x1, 0x10, 0x21, 0x18, 0x4, 0x20, 0xc, 0x30, 0x1, 0x42, 0x19, 0xa, 0x9, 0x8, 0x1, 0x18, 0x1, 0x20, 0x1, 0x2a, 0x1, 0x0, 0x12, 0xc, 0xa, 0x2, 0x0, 0x1, 0x10, 0x20, 0x18, 0x1, 0x20, 0x1, 0x30, 0x1, 0x4a, 0x7, 0x75, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x4a, 0x10, 0x75, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x64, 0x49, 0x42, 0x43, 0x53, 0x74, 0x61, 0x74, 0x65, 0x1a, 0xb, 0x8, 0x1, 0x18, 0x1, 0x20, 0x1, 0x2a, 0x3, 0x0, 0x2, 0x6, 0x22, 0x2b, 0x8, 0x1, 0x12, 0x4, 0x2, 0x4, 0xa, 0x20, 0x1a, 0x21, 0x20, 0x80, 0x19, 0x39, 0x94, 0x7e, 0x6e, 0xb1, 0x5, 0x20, 0xab, 0xb4, 0x9c, 0xa5, 0xca, 0x68, 0x82, 0x13, 0xf0, 0x96, 0x77, 0x2a, 0xeb, 0xae, 0x77, 0x3b, 0x5e, 0xf7, 0x1a, 0xd, 0xdc, 0xfe, 0x3, 0x22, 0x2b, 0x8, 0x1, 0x12, 0x4, 0x4, 0x6, 0xa, 0x20, 0x1a, 0x21, 0x20, 0x8a, 0xe9, 0xc2, 0x4f, 0x93, 0x86, 0xfb, 0x21, 0x6, 0x8c, 0x7, 0x73, 0x7b, 0xbc, 0x68, 0xdf, 0xfe, 0x29, 0xe3, 0xe0, 0xae, 0x70, 0x64, 0xf6, 0xa1, 0x31, 0x29, 0x88, 0x75, 0x73, 0x17, 0x40, 0x22, 0x2b, 0x8, 0x1, 0x12, 0x4, 0x6, 0xe, 0xa, 0x20, 0x1a, 0x21, 0x20, 0x43, 0x7f, 0x23, 0x36, 0x96, 0xb5, 0xe5, 0x75, 0x22, 0xf, 0xce, 0x90, 0x94, 0xae, 0xd6, 0x1d, 0xb5, 0x8f, 0xd4, 0x91, 0xbb, 0x95, 0x85, 0xa0, 0x1e, 0x5d, 0x3d, 0x78, 0xcf, 0xa2, 0xcd, 0x6b, 0x22, 0x2b, 0x8, 0x1, 0x12, 0x4, 0x8, 0x14, 0xa, 0x20, 0x1a, 0x21, 0x20, 0x9c, 0x2a, 0xf9, 0xb2, 0xb4, 0xe4, 0x4, 0xdb, 0x69, 0x90, 0x5d, 0xef, 0x9a, 0xd6, 0xc, 0xda, 0xfe, 0x86, 0x4c, 0x42, 0xcc, 0xb, 0x4, 0x9a, 0x16, 0x2e, 0xe5, 0x1f, 0x37, 0xa6, 0xe9, 0x31, 0xa, 0xfe, 0x1, 0xa, 0xfb, 0x1, 0xa, 0x3, 0x69, 0x62, 0x63, 0x12, 0x20, 0x6, 0x7, 0x6b, 0x34, 0xdd, 0x5f, 0x86, 0x82, 0x6f, 0xe8, 0x5e, 0x64, 0x43, 0x81, 0x63, 0x36, 0x63, 0x43, 0xba, 0x3f, 0x7b, 0x5e, 0x52, 0xf0, 0x26, 0x3e, 0x23, 0x23, 0xf7, 0xdb, 0x63, 0x79, 0x1a, 0x9, 0x8, 0x1, 0x18, 0x1, 0x20, 0x1, 0x2a, 0x1, 0x0, 0x22, 0x27, 0x8, 0x1, 0x12, 0x1, 0x1, 0x1a, 0x20, 0xc5, 0x5, 0xc0, 0xfd, 0x48, 0xb1, 0xcf, 0x2b, 0x65, 0x61, 0x9f, 0x12, 0xb3, 0x14, 0x4e, 0x19, 0xc6, 0xb, 0x4e, 0x6b, 0x62, 0x52, 0x5e, 0xe9, 0x36, 0xbe, 0x93, 0xd0, 0x41, 0x62, 0x3e, 0xbf, 0x22, 0x27, 0x8, 0x1, 0x12, 0x1, 0x1, 0x1a, 0x20, 0xb5, 0x94, 0xb4, 0xb9, 0x92, 0x43, 0xce, 0xda, 0xa5, 0x46, 0x37, 0x7d, 0xe1, 0xe, 0xbb, 0xc9, 0xfd, 0x57, 0xe3, 0xba, 0x3b, 0x76, 0x4, 0x82, 0x89, 0x2f, 0xa6, 0xc6, 0xd3, 0x4c, 0x82, 0x6a, 0x22, 0x25, 0x8, 0x1, 0x12, 0x21, 0x1, 0x47, 0xe, 0x2a, 0x73, 0x10, 0x14, 0xa6, 0x24, 0xf2, 0x68, 0xdc, 0x7b, 0xc9, 0x84, 0x2d, 0xc5, 0xb1, 0x65, 0x5, 0x9e, 0xb8, 0xc7, 0xe0, 0xb1, 0x5c, 0x73, 0x4e, 0x1d, 0x1b, 0x42, 0x53, 0xe4, 0x22, 0x25, 0x8, 0x1, 0x12, 0x21, 0x1, 0xb6, 0x5f, 0x9a, 0xae, 0x6c, 0x59, 0x3, 0xff, 0xf1, 0x9, 0x87, 0xcd, 0x1, 0xcf, 0x9, 0xd9, 0xe5, 0xbb, 0x9, 0x20, 0x80, 0x9c, 0xe2, 0x83, 0xee, 0x94, 0xc4, 0x8d, 0x5e, 0x11, 0xf7, 0xc4, 0x22, 0x27, 0x8, 0x1, 0x12, 0x1, 0x1, 0x1a, 0x20, 0xf8, 0xa3, 0xd7, 0x24, 0x32, 0x64, 0x99, 0x96, 0x6a, 0xdd, 0xc4, 0xc5, 0x58, 0x57, 0x51, 0xae, 0x96, 0x40, 0x2a, 0xca, 0xb6, 0x1a, 0xeb, 0x74, 0xd3, 0x54, 0x2a, 0xb7, 0x82, 0x7, 0x50, 0x38}
	proofHeightByte := []byte(nil)
	consensusHeightByte := []byte(nil)

	tx, err := contract.ConnOpenTry(auth, counterpartyBytes, delayPeriod, "07-tendermint-0", clientStateBytes, counterpartyVersions, proofInit, proofClient, proofConsensus, proofHeightByte, consensusHeightByte)
	gomega.Expect(err).Should(gomega.BeNil())

	receipt, err := waitForReceiptAndGet(ctx, ethClient, tx)
	gomega.Expect(err).Should(gomega.BeNil())

	log.Info("ConnOpenTry tx", "hash", tx.Hash(), "blockNumber", receipt.BlockNumber.Uint64())
}

func getClientStateBytes() []byte {
	clientState := ibctm.NewClientState(testChainID, ibctm.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath)
	clientStateBytes, err := clientState.Marshal()
	gomega.Expect(err).Should(gomega.BeNil())

	return clientStateBytes
}

func getConsensusStateBytes() []byte {
	now := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	privVal := ibctestingmock.NewPV()
	pubKey, err := privVal.GetPubKey()
	gomega.Expect(err).Should(gomega.BeNil())

	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
	consensusState := ibctm.NewConsensusState(now, commitmenttypes.NewMerkleRoot([]byte("hash")), valSet.Hash())
	consensusStateBytes, err := consensusState.Marshal()
	gomega.Expect(err).Should(gomega.BeNil())

	return consensusStateBytes
}

func getCounterpartyBytes() []byte {
	prefix := commitmenttypes.NewMerklePrefix([]byte("storePrefixKey"))
	counterparty := connectiontypes.NewCounterparty("connectiontotest", "clienttotest", prefix)
	counterpartyBytes, err := counterparty.Marshal()
	gomega.Expect(err).Should(gomega.BeNil())

	return counterpartyBytes
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