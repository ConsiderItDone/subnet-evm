package utils

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	wallet "github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/plugin/evm"
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
	testKey, _       = crypto.HexToECDSA("56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027")
	chainId          = big.NewInt(99999)
	testClientHeight = clienttypes.NewHeight(0, 5)

	ethClient           ethclient.Client
	ibcContract         *contractBind.Contract
	ibcContractFilterer *contractBind.ContractFilterer
	auth                *bind.TransactOpts

	coordinator *ibctesting.Coordinator
	chainA      *ibctesting.TestChain
	chainB      *ibctesting.TestChain
	path        *ibctesting.Path

	clientIdA     = "07-tendermint-0"
	clientIdB     = "07-tendermint-1"
	connectionId0 = "connection-0"

	marshaler *codec.ProtoCodec
)

func init() {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	connectiontypes.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler = codec.NewProtoCodec(interfaceRegistry)
}

func RunTestIbcInit(t *testing.T) {
	t.Log("executing new blockchain initialization")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	kc := secp256k1fx.NewKeychain(genesis.EWOQKey)

	// NewWalletFromURI fetches the available UTXOs owned by [kc] on the network
	// that [LocalAPIURI] is hosting.
	wallet, err := wallet.NewWalletFromURI(ctx, DefaultLocalNodeURI, kc)
	require.NoError(t, err)

	pWallet := wallet.P()

	owner := &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs: []ids.ShortID{
			genesis.EWOQKey.PublicKey().Address(),
		},
	}

	genesisBytes, err := os.ReadFile("../precompile/genesis/ibc.json")
	require.NoError(t, err)

	createSubnetTxID, err := pWallet.IssueCreateSubnetTx(owner)
	require.NoError(t, err)
	t.Logf("new subnet id: %s", createSubnetTxID)

	genesis := new(core.Genesis)
	require.NoError(t, genesis.UnmarshalJSON(genesisBytes))

	createChainTxID, err := pWallet.IssueCreateChainTx(
		createSubnetTxID,
		genesisBytes,
		evm.ID,
		nil,
		"testChain",
	)
	require.NoError(t, err)
	t.Logf("new chain id: %s", createSubnetTxID)

	// Confirm the new blockchain is ready by waiting for the readiness endpoint
	infoClient := info.NewClient(DefaultLocalNodeURI)
	bootstrapped, err := info.AwaitBootstrapped(ctx, infoClient, createChainTxID.String(), 2*time.Second)
	require.NoError(t, err)
	require.True(t, bootstrapped, "network isn't bootstaped")

	chainURI := GetDefaultChainURI(createChainTxID.String())
	t.Log("subnet successfully created: %s", chainURI)

	ethClient, err = ethclient.DialContext(ctx, chainURI)
	require.NoError(t, err)
	t.Log("eth client created")

	ibcContract, err = contractBind.NewContract(ibc.ContractAddress, ethClient)
	require.NoError(t, err)
	ibcContractFilterer, err = contractBind.NewContractFilterer(ibc.ContractAddress, ethClient)
	require.NoError(t, err)
	t.Log("contract binded")

	auth, err = bind.NewKeyedTransactorWithChainID(testKey, chainId)
	require.NoError(t, err)
	t.Log("transactor created")

	coordinator = ibctesting.NewCoordinator(t, 2)
	chainA = coordinator.GetChain(ibctesting.GetChainID(1))
	chainB = coordinator.GetChain(ibctesting.GetChainID(2))
	path = ibctesting.NewPath(chainA, chainB)
	coordinator.SetupClients(path)
}

func RunTestIbcCreateClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	clientStateA, existA := chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(chainA.GetContext(), path.EndpointA.ClientID)
	require.True(t, existA)
	clientStateByteA, err := clientStateA.(*ibctm.ClientState).Marshal()
	require.NoError(t, err)

	consensusStateA, foundA := chainA.App.GetIBCKeeper().ClientKeeper.GetLatestClientConsensusState(chainA.GetContext(), path.EndpointB.ClientID)
	require.True(t, foundA)
	consensusStateByteA, err := consensusStateA.(*ibctm.ConsensusState).Marshal()
	require.NoError(t, err)

	txA, err := ibcContract.CreateClient(auth, exported.Tendermint, clientStateByteA, consensusStateByteA)
	require.NoError(t, err)
	reA, err := waitForReceiptAndGet(ctx, ethClient, txA)
	require.NoError(t, err)
	require.True(t, len(reA.Logs) > 0)
	eventA, err := ibcContractFilterer.ParseClientCreated(*reA.Logs[0])
	require.NoError(t, err)
	assert.Equal(t, clientIdA, eventA.ClientId)

	clientStateB, existB := chainB.App.GetIBCKeeper().ClientKeeper.GetClientState(chainB.GetContext(), path.EndpointB.ClientID)
	require.True(t, existB)
	clientStateByteB, err := clientStateB.(*ibctm.ClientState).Marshal()
	require.NoError(t, err)

	consensusStateB, foundB := chainB.App.GetIBCKeeper().ClientKeeper.GetLatestClientConsensusState(chainB.GetContext(), path.EndpointB.ClientID)
	require.True(t, foundB)
	consensusStateByteB, err := consensusStateB.(*ibctm.ConsensusState).Marshal()
	require.NoError(t, err)

	txB, err := ibcContract.CreateClient(auth, exported.Tendermint, clientStateByteB, consensusStateByteB)
	require.NoError(t, err)
	reB, err := waitForReceiptAndGet(ctx, ethClient, txB)
	require.NoError(t, err)
	require.True(t, len(reA.Logs) > 0)
	eventB, err := ibcContractFilterer.ParseClientCreated(*reB.Logs[0])
	require.NoError(t, err)
	assert.Equal(t, clientIdB, eventB.ClientId)
}

func RunTestIbcConnectionOpenInit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	counterparty := connectiontypes.NewCounterparty(path.EndpointB.ClientID, path.EndpointB.ConnectionID, chainB.GetPrefix())
	counterpartybyte, err := marshaler.MarshalInterface(&counterparty)
	require.NoError(t, err)

	version := connectiontypes.ExportedVersionsToProto(connectiontypes.GetCompatibleVersions())[0]
	versionByte, err := marshaler.Marshal(version)
	require.NoError(t, err)

	tx, err := ibcContract.ConnOpenInit(auth, clientIdA, counterpartybyte, versionByte, 0)
	require.NoError(t, err)
	re, err := waitForReceiptAndGet(ctx, ethClient, tx)
	require.NoError(t, err)
	require.True(t, len(re.Logs) > 0)

	ev, err := ibcContractFilterer.ParseConnectionCreated(*re.Logs[0])
	require.NoError(t, err)
	assert.Equal(t, clientIdA, ev.ClientId)
	assert.Equal(t, connectionId0, ev.ConnectionId)

}

func RunTestIbcConnectionOpenAck(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	require.NoError(t, path.EndpointA.ConnOpenInit())
	require.NoError(t, path.EndpointB.ConnOpenTry())

	counterpartyClient := chainB.GetClientState(path.EndpointB.ClientID)
	require.NoError(t, path.EndpointA.UpdateClient())

	connectionKey := host.ConnectionKey(path.EndpointB.ConnectionID)
	proofTry, proofHeight := chainB.QueryProof(connectionKey)

	clientState := chainB.GetClientState(path.EndpointB.ClientID)
	consensusHeight := clientState.GetLatestHeight().(clienttypes.Height)

	consensusKey := host.FullConsensusStateKey(path.EndpointB.ClientID, consensusHeight)
	proofConsensus, _ := chainB.QueryProof(consensusKey)

	clientKey := host.FullClientStateKey(path.EndpointB.ClientID)
	proofClient, _ := chainB.QueryProof(clientKey)

	clientStateByte, err := clienttypes.MarshalClientState(marshaler, counterpartyClient)
	require.NoError(t, err)

	version := connectiontypes.ExportedVersionsToProto(connectiontypes.GetCompatibleVersions())[0]
	versionByte, err := marshaler.Marshal(version)
	require.NoError(t, err)

	proofHeightByte, err := marshaler.Marshal(&proofHeight)
	require.NoError(t, err)

	consensusHeightByte, err := marshaler.Marshal(&consensusHeight)
	require.NoError(t, err)

	tx, err := ibcContract.ConnOpenAck(
		auth,
		path.EndpointA.ConnectionID,
		clientStateByte,
		versionByte,
		[]byte(path.EndpointB.ConnectionID),
		proofTry,
		proofClient,
		proofConsensus,
		proofHeightByte,
		consensusHeightByte,
	)
	require.NoError(t, err)
	re, err := waitForReceiptAndGet(ctx, ethClient, tx)
	require.NoError(t, err)
	t.Log(spew.Sdump(re.Logs))
}

func RunTestIncChannelOpenInit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	require.NoError(t, path.EndpointA.ConnOpenAck())
	require.NoError(t, path.EndpointA.UpdateClient())
	require.NoError(t, path.EndpointB.ConnOpenConfirm())
	require.NoError(t, path.EndpointB.UpdateClient())

	counterparty := channeltypes.NewCounterparty(ibctesting.MockPort, ibctesting.FirstChannelID)
	channel := channeltypes.NewChannel(channeltypes.INIT, channeltypes.ORDERED, counterparty, []string{path.EndpointB.ConnectionID}, path.EndpointA.ChannelConfig.Version)
	channelByte, err := marshaler.Marshal(&channel)
	require.NoError(t, err)

	tx, err := ibcContract.ChanOpenInit(auth, path.EndpointA.ChannelConfig.PortID, channelByte)
	require.NoError(t, err)
	re, err := waitForReceiptAndGet(ctx, ethClient, tx)
	require.NoError(t, err)
	t.Log(spew.Sdump(re.Logs))
}

func RunTestIncChannelOpenAck(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	require.NoError(t, path.EndpointA.ChanOpenInit())
	require.NoError(t, path.EndpointA.UpdateClient())
	require.NoError(t, path.EndpointB.ChanOpenTry())
	require.NoError(t, path.EndpointB.UpdateClient())

	channelKey := host.ChannelKey(path.EndpointB.ChannelConfig.PortID, ibctesting.FirstChannelID)
	proof, proofHeight := chainB.QueryProof(channelKey)

	proofHeightByte, err := marshaler.Marshal(&proofHeight)
	require.NoError(t, err)

	tx, err := ibcContract.ChannelOpenAck(
		auth,
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		path.EndpointB.ChannelID,
		path.EndpointB.ChannelConfig.Version,
		proof,
		proofHeightByte,
	)
	require.NoError(t, err)
	re, err := waitForReceiptAndGet(ctx, ethClient, tx)
	require.NoError(t, err)
	t.Log(spew.Sdump(re.Logs))
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
