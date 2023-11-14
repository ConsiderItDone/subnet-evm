package utils

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
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
	"github.com/ava-labs/subnet-evm/ethclient/subnetevmclient"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/precompile/contracts/ibc"
	"github.com/ava-labs/subnet-evm/precompile/contracts/ics20"
	"github.com/ava-labs/subnet-evm/rpc"
	contractBind "github.com/ava-labs/subnet-evm/tests/precompile/contract"
	"github.com/ava-labs/subnet-evm/tests/precompile/contract/ics20/ics20bank"
	"github.com/ava-labs/subnet-evm/tests/precompile/contract/ics20/ics20transferer"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
)

const (
	testChainID = "gaiahub-0"
	testPort    = ibctesting.TransferPort

	trustingPeriod time.Duration = time.Hour * 24 * 7 * 2
	ubdPeriod      time.Duration = time.Hour * 24 * 7 * 3
	maxClockDrift  time.Duration = time.Second * 10
)

var (
	testKey, _ = crypto.HexToECDSA("56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027")
	chainId    = big.NewInt(99999)

	ethClient           ethclient.Client
	subnetClient        *subnetevmclient.Client
	ibcContract         *contractBind.Contract
	ibcContractFilterer *contractBind.ContractFilterer
	auth                *bind.TransactOpts

	ics20BankAddr       common.Address
	ics20Bank           *ics20bank.Ics20bank
	ics20TransfererAddr common.Address
	ics20Transferer     *ics20transferer.Ics20transferer

	coordinator *ibctesting.Coordinator
	chainA      *ibctesting.TestChain
	chainB      *ibctesting.TestChain
	path        *ibctesting.Path

	clientIdA     = "07-tendermint-0"
	clientIdB     = "07-tendermint-1"
	connectionId0 = "connection-0"

	disabledTimeoutTimestamp = uint64(0)
	disabledTimeoutHeight    = clienttypes.ZeroHeight()
	defaultTimeoutHeight     = clienttypes.NewHeight(1, 100)

	marshaler *codec.ProtoCodec
)

func init() {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	connectiontypes.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler = codec.NewProtoCodec(interfaceRegistry)
}

func RunTestIbcInit(t *testing.T) {
	t.Log("executing new blockchain initialization")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Minute)
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
	t.Logf("subnet successfully created: %s", chainURI)

	rpcClient, err := rpc.DialContext(ctx, chainURI)
	require.NoError(t, err)

	ethClient = ethclient.NewClient(rpcClient)
	subnetClient = subnetevmclient.New(rpcClient)
	t.Log("eth client created")

	ibcContract, err = contractBind.NewContract(ibc.ContractAddress, ethClient)
	require.NoError(t, err)
	ibcContractFilterer, err = contractBind.NewContractFilterer(ibc.ContractAddress, ethClient)
	require.NoError(t, err)
	t.Log("contract binded")

	auth, err = bind.NewKeyedTransactorWithChainID(testKey, chainId)
	require.NoError(t, err)
	t.Log("transactor created")

	ics20bankAddr, ics20bankTx, ics20bank, err := ics20bank.DeployIcs20bank(auth, ethClient)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, ics20bankTx)
	require.NoError(t, err)
	ics20Bank = ics20bank
	ics20BankAddr = ics20bankAddr

	ics20transfererAddr, ics20transfererTx, ics20transferer, err := ics20transferer.DeployIcs20transferer(auth, ethClient, ibc.ContractAddress, ics20bankAddr)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, ics20transfererTx)
	require.NoError(t, err)
	ics20Transferer = ics20transferer
	ics20TransfererAddr = ics20transfererAddr

	setOperTx1, err := ics20bank.SetOperator(auth, auth.From)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, setOperTx1)
	require.NoError(t, err)

	setOperTx2, err := ics20bank.SetOperator(auth, ibc.ContractAddress)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, setOperTx2)
	require.NoError(t, err)

	setOperTx3, err := ics20bank.SetOperator(auth, ics20transfererAddr)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, setOperTx3)
	require.NoError(t, err)

	coordinator = ibctesting.NewCoordinator(t, 2)
	coordinator.CurrentTime = time.Now()
	chainA = coordinator.GetChain(ibctesting.GetChainID(1))
	coordinator.UpdateTimeForChain(chainA)
	chainB = coordinator.GetChain(ibctesting.GetChainID(2))
	coordinator.UpdateTimeForChain(chainB)
	path = ibctesting.NewPath(chainA, chainB)
	coordinator.SetupClients(path)
}

func RunTestIbcCreateClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	createIbcClient(t, ctx, path.EndpointA, clientIdA)
	createIbcClient(t, ctx, path.EndpointB, clientIdB)
}

func RunTestIbcConnectionOpenInit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	counterparty := connectiontypes.NewCounterparty(path.EndpointB.ClientID, path.EndpointB.ConnectionID, commitmenttypes.NewMerklePrefix([]byte("ibc")))
	counterpartybyte, err := counterparty.Marshal()
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

	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointA.ConnOpenInit)
	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, path.EndpointB.ConnOpenTry)
	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointA.UpdateClient)
}

func RunTestIbcConnectionOpenTry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointA.ConnOpenInit)
	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, path.EndpointB.UpdateClient)

	counterpartyClient := chainA.GetClientState(path.EndpointA.ClientID)
	counterparty := connectiontypes.NewCounterparty(path.EndpointA.ClientID, path.EndpointA.ConnectionID, chainA.GetPrefix())

	connectionKey := host.ConnectionKey(path.EndpointA.ConnectionID)
	proofInit, proofHeight := chainA.QueryProof(connectionKey)
	fmt.Printf("proofInit %#v\n", proofInit)
	fmt.Printf("proofHeight %#v\n", proofHeight)

	versions := connectiontypes.GetCompatibleVersions()
	consensusHeight := counterpartyClient.GetLatestHeight().(clienttypes.Height)

	consensusKey := host.FullConsensusStateKey(path.EndpointA.ClientID, consensusHeight)
	proofConsensus, _ := chainA.QueryProof(consensusKey)
	fmt.Printf("proofConsensus %#v\n", proofConsensus)

	// retrieve proof of counterparty clientstate on chainA
	clientKey := host.FullClientStateKey(path.EndpointA.ClientID)
	proofClient, _ := chainA.QueryProof(clientKey)
	fmt.Printf("proofClient %#v\n", proofClient)

	counterpartyByte, _ := counterparty.Marshal()
	fmt.Printf("counterparty %#v\n", counterpartyByte)

	clientStateByte, _ := clienttypes.MarshalClientState(marshaler, counterpartyClient)
	fmt.Printf("clientState %#v\n", clientStateByte)

	versionsByte, _ := json.Marshal(connectiontypes.ExportedVersionsToProto(versions))
	fmt.Printf("versions %#v\n", versionsByte)

	proofHeightByte, _ := proofHeight.Marshal()
	fmt.Printf("proofHeightByte %#v\n", proofHeightByte)

	consensusHeightByte, _ := marshaler.MarshalInterface(&consensusHeight)
	fmt.Printf("consensusHeightByte %#v\n", consensusHeightByte)

	tx, err := ibcContract.ConnOpenTry(
		auth,
		counterpartyByte,
		0,
		clientIdB,
		clientStateByte,
		versionsByte,
		proofInit,
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

func RunTestIbcConnectionOpenAck(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	counterpartyClient := chainB.GetClientState(path.EndpointB.ClientID)

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

	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointA.ConnOpenAck)
	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, path.EndpointB.ConnOpenConfirm)
	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointA.UpdateClient)
}

func RunTestIbcConnectionOpenConfirm(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, path.EndpointB.ConnOpenTry)
	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointB.ConnOpenAck)
	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, path.EndpointB.UpdateClient)

	connectionKey := host.ConnectionKey(path.EndpointA.ConnectionID)
	proofAck, proofHeight := chainA.QueryProof(connectionKey)

	proofHeightByte, err := marshaler.Marshal(&proofHeight)
	require.NoError(t, err)

	tx, err := ibcContract.ConnOpenConfirm(auth, path.EndpointB.ConnectionID, proofAck, proofHeightByte)
	require.NoError(t, err)
	re, err := waitForReceiptAndGet(ctx, ethClient, tx)
	require.NoError(t, err)
	t.Log(spew.Sdump(re.Logs))
}

func RunTestIbcChannelOpenInit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	path.EndpointA.ChannelConfig.PortID = testPort
	path.EndpointA.ChannelConfig.Order = channeltypes.UNORDERED
	path.EndpointA.ChannelConfig.Version = transfertypes.Version

	path.EndpointB.ChannelConfig.PortID = testPort
	path.EndpointB.ChannelConfig.Order = channeltypes.UNORDERED
	path.EndpointB.ChannelConfig.Version = transfertypes.Version

	counterparty := channeltypes.NewCounterparty(testPort, ibctesting.FirstChannelID)
	channel := channeltypes.NewChannel(channeltypes.INIT, channeltypes.UNORDERED, counterparty, []string{path.EndpointB.ConnectionID}, path.EndpointA.ChannelConfig.Version)
	channelByte, err := marshaler.Marshal(&channel)
	require.NoError(t, err)

	tx, err := ibcContract.ChanOpenInit(auth, testPort, channelByte)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, tx)
	require.NoError(t, err)

	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, func() error {
		if err := path.EndpointA.ChanOpenInit(); err != nil {
			return err
		}
		return path.EndpointA.UpdateClient()
	})
	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, path.EndpointB.ChanOpenTry)
	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointA.UpdateClient)
}

func RunTestIbcChannelOpenTry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, path.EndpointB.ConnOpenConfirm)
	path.SetChannelOrdered()
	require.NoError(t, path.EndpointA.ChanOpenInit())
	chainB.CreatePortCapability(chainB.GetSimApp().ScopedIBCMockKeeper, ibctesting.MockPort)

	counterparty := channeltypes.NewCounterparty(path.EndpointB.ChannelConfig.PortID, ibctesting.FirstChannelID)
	channel := channeltypes.NewChannel(channeltypes.INIT, channeltypes.ORDERED, counterparty, []string{path.EndpointB.ConnectionID}, path.EndpointA.ChannelConfig.Version)
	channelByte, err := marshaler.Marshal(&channel)
	require.NoError(t, err)

	channelKey := host.ChannelKey(counterparty.PortId, counterparty.ChannelId)
	proof, proofHeight := chainA.QueryProof(channelKey)

	consensusHeightByte, err := proofHeight.Marshal()
	require.NoError(t, err)
	tx, err := ibcContract.ChanOpenTry(
		auth,
		path.EndpointB.ChannelConfig.PortID,
		channelByte,
		path.EndpointA.ChannelConfig.Version,
		proof,
		consensusHeightByte,
	)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, tx)
	require.NoError(t, err)
}

func RunTestIbcChannelOpenAck(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	channelKey := host.ChannelKey(path.EndpointB.ChannelConfig.PortID, ibctesting.FirstChannelID)
	proof, proofHeight := chainB.QueryProof(channelKey)
	proofHeightByte, err := marshaler.Marshal(&proofHeight)
	require.NoError(t, err)

	tx, err := ibcContract.ChannelOpenAck(
		auth,
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		ibctesting.FirstChannelID,
		path.EndpointB.ChannelConfig.Version,
		proof,
		proofHeightByte,
	)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, tx)
	require.NoError(t, err)

	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointA.ChanOpenAck)
	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, path.EndpointB.ChanOpenConfirm)
}

func RunTestIbcChannelOpenConfirm(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, path.EndpointB.ChanOpenTry)
	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointA.ChanOpenAck)
	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, path.EndpointB.UpdateClient)

	channelKey := host.ChannelKey(path.EndpointA.ChannelConfig.PortID, ibctesting.FirstChannelID)
	proof, proofHeight := chainA.QueryProof(channelKey)

	consensusHeightByte, err := marshaler.Marshal(&proofHeight)
	require.NoError(t, err)

	tx, err := ibcContract.ChannelOpenConfirm(
		auth,
		path.EndpointB.ChannelConfig.PortID,
		ibctesting.FirstChannelID,
		proof,
		consensusHeightByte,
	)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, tx)
	require.NoError(t, err)
}

func RunTestIbcRecvPacket(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	mintFungibleTokenPacketData, err := json.Marshal(ics20.FungibleTokenPacketData{
		Denom:    "USDT",
		Amount:   "1000",
		Sender:   common.Address{}.Hex(),
		Receiver: auth.From.Hex(),
		Memo:     "some memo",
	})
	require.NoError(t, err)

	sequence, err := path.EndpointB.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, mintFungibleTokenPacketData)
	require.NoError(t, err)
	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointA.UpdateClient)
	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, nil)

	packetKey := host.PacketCommitmentKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence)
	proof, proofHeight := path.EndpointB.QueryProof(packetKey)

	bindTx, err := ics20Transferer.BindPort(auth, ibc.ContractAddress, testPort)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, bindTx)
	require.NoError(t, err)

	setEscrowAddrTx, err := ics20Transferer.SetChannelEscrowAddresses(auth, path.EndpointA.ChannelID, auth.From)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, setEscrowAddrTx)
	require.NoError(t, err)

	auth.GasLimit = 200000
	recvTx, err := ibcContract.RecvPacket(auth, contractBind.IIBCMsgRecvPacket{
		Packet: contractBind.Packet{
			Sequence:           big.NewInt(int64(sequence)),
			SourcePort:         path.EndpointA.ChannelConfig.PortID,
			SourceChannel:      path.EndpointA.ChannelID,
			DestinationPort:    testPort,
			DestinationChannel: path.EndpointB.ChannelID,
			Data:               mintFungibleTokenPacketData,
			TimeoutHeight: contractBind.Height{
				RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
				RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
			},
			TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
		},
		ProofCommitment: proof,
		ProofHeight: contractBind.Height{
			RevisionNumber: new(big.Int).SetUint64(proofHeight.RevisionNumber),
			RevisionHeight: new(big.Int).SetUint64(proofHeight.RevisionHeight),
		},
		Signer: "",
	})
	auth.GasLimit = 0
	require.NoError(t, err)

	re, err := waitForReceiptAndGet(ctx, ethClient, recvTx)
	require.NoError(t, err)
	require.Equal(t, len(re.Logs), 2)

	transferlog, err := ics20Bank.Ics20bankFilterer.ParseTransfer(*re.Logs[1])
	require.NoError(t, err)
	assert.Equal(t, common.Address{}, transferlog.From)
	assert.Equal(t, auth.From, transferlog.To)
	assert.Equal(t, "transfer/channel-0/USDT", transferlog.Path)
	assert.Equal(t, big.NewInt(1000), transferlog.Value)
}

func RunTestIbcSendPacket(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	mintFungibleTokenPacketData, err := json.Marshal(ics20.FungibleTokenPacketData{
		Denom:    "USDT",
		Amount:   "1000",
		Sender:   common.Address{}.Hex(),
		Receiver: auth.From.Hex(),
		Memo:     "some memo",
	})
	require.NoError(t, err)

	_, err = path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, mintFungibleTokenPacketData)
	require.NoError(t, err)
	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, nil)
	updateIbcClientAfterFunc(t, clientIdB, path.EndpointB, path.EndpointB.UpdateClient)

	bindTx, err := ics20Transferer.BindPort(auth, ibc.ContractAddress, testPort)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, bindTx)
	require.NoError(t, err)

	setEscrowAddrTx, err := ics20Transferer.SetChannelEscrowAddresses(auth, path.EndpointB.ChannelID, auth.From)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, setEscrowAddrTx)
	require.NoError(t, err)

	auth.GasLimit = 200000

	recvTx, err := ibcContract.SendPacket(auth,
		big.NewInt(int64(0)),
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		contractBind.Height{
			RevisionNumber: new(big.Int).SetUint64(defaultTimeoutHeight.RevisionNumber),
			RevisionHeight: new(big.Int).SetUint64(defaultTimeoutHeight.RevisionHeight),
		},
		big.NewInt(int64(disabledTimeoutTimestamp)),
		mintFungibleTokenPacketData)

	auth.GasLimit = 0
	require.NoError(t, err)

	re, err := waitForReceiptAndGet(ctx, ethClient, recvTx)
	require.NoError(t, err)
	require.Equal(t, len(re.Logs), 1)
}

func RunTestIbcAckPacket(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	addr, err := cosmostypes.Bech32ifyAddressBytes("cosmos", auth.From.Bytes())
	require.NoError(t, err)
	amount := 1000

	//prefix := fmt.Sprintf("%s/%s/", path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
	//prefix := "ibc"
	prefix := ""
	mintFungibleTokenPacket := transfertypes.FungibleTokenPacketData{
		Denom:    prefix + "USDT",
		Amount:   strconv.Itoa(amount),
		Sender:   auth.From.Hex(),
		Receiver: addr,
		Memo:     "some memo",
	}

	mintFungibleTokenPacketData, err := transfertypes.ModuleCdc.MarshalJSON(&mintFungibleTokenPacket)
	require.NoError(t, err)

	bindTx, err := ics20Transferer.BindPort(auth, ibc.ContractAddress, testPort)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, bindTx)
	require.NoError(t, err)

	setEscrowAddrTx, err := ics20Transferer.SetChannelEscrowAddresses(auth, path.EndpointB.ChannelID, auth.From)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, setEscrowAddrTx)
	require.NoError(t, err)

	auth.GasLimit = 200000
	sendTx, err := ibcContract.SendPacket(auth,
		big.NewInt(int64(0)),
		testPort,
		path.EndpointA.ChannelID,
		contractBind.Height{
			RevisionNumber: new(big.Int).SetUint64(defaultTimeoutHeight.RevisionNumber),
			RevisionHeight: new(big.Int).SetUint64(defaultTimeoutHeight.RevisionHeight),
		},
		big.NewInt(int64(disabledTimeoutTimestamp)),
		mintFungibleTokenPacketData)

	auth.GasLimit = 0
	require.NoError(t, err)
	re, err := waitForReceiptAndGet(ctx, ethClient, sendTx)
	require.NoError(t, err)
	require.Equal(t, len(re.Logs), 1)

	sequence, err := path.EndpointA.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, mintFungibleTokenPacketData)
	require.NoError(t, err)

	err = path.EndpointB.RecvPacket(channeltypes.Packet{
		Sequence:           sequence,
		SourcePort:         testPort,
		SourceChannel:      path.EndpointA.ChannelID,
		DestinationPort:    testPort,
		DestinationChannel: path.EndpointB.ChannelID,
		TimeoutHeight: clienttypes.Height{
			RevisionNumber: defaultTimeoutHeight.RevisionNumber,
			RevisionHeight: defaultTimeoutHeight.RevisionHeight,
		},
		TimeoutTimestamp: disabledTimeoutTimestamp,
		Data:             mintFungibleTokenPacketData,
	})
	require.NoError(t, err)

	respAllBalances, err := path.EndpointB.Chain.GetSimApp().BankKeeper.AllBalances(path.EndpointB.Chain.GetContext(), &banktypes.QueryAllBalancesRequest{
		Address: addr,
	})
	require.NoError(t, err)
	denom := respAllBalances.Balances[0].Denom
	resp, err := path.EndpointB.Chain.GetSimApp().BankKeeper.Balance(path.EndpointB.Chain.GetContext(), &banktypes.QueryBalanceRequest{
		Address: addr,
		Denom:   denom,
	})
	require.NoError(t, err)
	require.Equal(t, resp.Balance.Amount.Int64(), int64(amount))

	hexHash := denom[len(transfertypes.DenomPrefix+"/"):]
	hash, err := transfertypes.ParseHexHash(hexHash)
	require.NoError(t, err)
	denometrace, found := path.EndpointB.Chain.GetSimApp().TransferKeeper.GetDenomTrace(path.EndpointB.Chain.GetContext(), hash)
	if !found {
		require.NoError(t, transfertypes.ErrTraceNotFound)
	}
	require.Equal(t, denometrace.Path, "transfer/channel-0")
	require.Equal(t, denometrace.BaseDenom, "USDT")

	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointA.UpdateClient)

	acknowledgement := channeltypes.NewResultAcknowledgement([]byte{byte(1)})

	packetKey := host.PacketAcknowledgementKey(testPort, path.EndpointB.ChannelID, sequence)
	proof, proofHeight := path.EndpointB.QueryProof(packetKey)

	auth.GasLimit = 200000
	recvTx, err := ibcContract.SendPacket(auth,
		big.NewInt(int64(0)),
		ibctesting.TransferPort,
		path.EndpointA.ChannelID,
		contractBind.Height{
			RevisionNumber: new(big.Int).SetUint64(defaultTimeoutHeight.RevisionNumber),
			RevisionHeight: new(big.Int).SetUint64(defaultTimeoutHeight.RevisionHeight),
		},
		big.NewInt(int64(disabledTimeoutTimestamp)),
		mintFungibleTokenPacketData)

	auth.GasLimit = 0
	require.NoError(t, err)
	re, err := waitForReceiptAndGet(ctx, ethClient, recvTx)
	require.NoError(t, err)
	require.Equal(t, len(re.Logs), 1)

	err = path.EndpointB.RecvPacket(channeltypes.Packet{
		Sequence:           sequence,
		SourcePort:         ibctesting.TransferPort,
		SourceChannel:      path.EndpointA.ChannelID,
		DestinationPort:    ibctesting.TransferPort,
		DestinationChannel: path.EndpointB.ChannelID,
		TimeoutHeight: clienttypes.Height{
			RevisionNumber: defaultTimeoutHeight.RevisionNumber,
			RevisionHeight: defaultTimeoutHeight.RevisionHeight,
		},
		TimeoutTimestamp: disabledTimeoutTimestamp,
		Data:             mintFungibleTokenPacketData,
	})
	require.NoError(t, err)

	updateIbcClientAfterFunc(t, clientIdA, path.EndpointA, path.EndpointA.UpdateClient)

	packetKey := host.PacketAcknowledgementKey(ibctesting.TransferPort, path.EndpointB.ChannelID, sequence)
	proof, proofHeight := path.EndpointB.QueryProof(packetKey)

	acknowledgement := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	ack := sha256.Sum256(acknowledgement.Acknowledgement())

	packetAckTx, err := ibcContract.Acknowledgement(
		auth,
		contractBind.Packet{
			Sequence:           big.NewInt(int64(sequence)),
			SourcePort:         testPort,
			SourceChannel:      path.EndpointA.ChannelID,
			DestinationPort:    testPort,
			DestinationChannel: path.EndpointB.ChannelID,
			Data:               mintFungibleTokenPacketData,
			TimeoutHeight: contractBind.Height{
				RevisionNumber: new(big.Int).SetUint64(defaultTimeoutHeight.RevisionNumber),
				RevisionHeight: new(big.Int).SetUint64(defaultTimeoutHeight.RevisionHeight),
			},
			TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
		},
		acknowledgement.Acknowledgement(),
		proof,
		contractBind.Height{
			RevisionNumber: big.NewInt(int64(proofHeight.RevisionNumber)),
			RevisionHeight: big.NewInt(int64(proofHeight.RevisionHeight)),
		},
		"",
	)
	require.NoError(t, err)
	_, err = waitForReceiptAndGet(ctx, ethClient, packetAckTx)
	require.NoError(t, err)

	// 	require.NoError(t, err)
	// 	mintFungibleTokenPacket = transfertypes.FungibleTokenPacketData{
	// 		Denom:    prefix + "USDT",
	// 		Amount:   strconv.Itoa(amount),
	// 		Sender:   addr,
	// 		Receiver: auth.From.Hex(),
	// 		Memo:     "some memo",
	// 	}
	// 	mintFungibleTokenPacketData, err = transfertypes.ModuleCdc.MarshalJSON(&mintFungibleTokenPacket)
	// 	require.NoError(t, err)

	// 	sequence, err = path.EndpointB.SendPacket(defaultTimeoutHeight, disabledTimeoutTimestamp, mintFungibleTokenPacketData)
	// 	require.NoError(t, err)

	// 	auth.GasLimit = 200000
	// 	recvTx, err := ibcContract.RecvPacket(auth, contractBind.IIBCMsgRecvPacket{
	// 		Packet: contractBind.Packet{
	// 			Sequence:           big.NewInt(int64(sequence)),
	// 			SourcePort:         testPort,
	// 			SourceChannel:      path.EndpointB.ChannelID,
	// 			DestinationPort:    testPort,
	// 			DestinationChannel: path.EndpointA.ChannelID,
	// 			Data:               mintFungibleTokenPacketData,
	// 			TimeoutHeight: contractBind.Height{
	// 				RevisionNumber: big.NewInt(int64(defaultTimeoutHeight.RevisionNumber)),
	// 				RevisionHeight: big.NewInt(int64(defaultTimeoutHeight.RevisionHeight)),
	// 			},
	// 			TimeoutTimestamp: big.NewInt(int64(disabledTimeoutTimestamp)),
	// 		},
	// 		ProofCommitment: proof,
	// 		ProofHeight: contractBind.Height{
	// 			RevisionNumber: new(big.Int).SetUint64(proofHeight.RevisionNumber),
	// 			RevisionHeight: new(big.Int).SetUint64(proofHeight.RevisionHeight),
	// 		},
	// 		Signer: "",
	// 	})
	// 	auth.GasLimit = 0
	// 	require.NoError(t, err)

	// 	re, err = waitForReceiptAndGet(ctx, ethClient, recvTx)
	// 	require.NoError(t, err)
	// 	require.Equal(t, len(re.Logs), 2)

	// 	transferlog, err := ics20Bank.Ics20bankFilterer.ParseTransfer(*re.Logs[1])
	// 	require.NoError(t, err)
	// 	assert.Equal(t, common.Address{}, transferlog.From)
	// 	assert.Equal(t, auth.From, transferlog.To)
	// 	assert.Equal(t, "transfer/channel-0/USDT", transferlog.Path)
	// 	assert.Equal(t, big.NewInt(1000), transferlog.Value)

	// 	acknowledgement = channeltypes.NewResultAcknowledgement([]byte{byte(1)})

	//	path.EndpointB.AcknowledgePacket(channeltypes.Packet{
	//		Sequence:           sequence,
	//		SourcePort:         testPort,
	//		SourceChannel:      path.EndpointB.ChannelID,
	//		DestinationPort:    testPort,
	//		DestinationChannel: path.EndpointA.ChannelID,
	//		TimeoutHeight: clienttypes.Height{
	//			RevisionNumber: defaultTimeoutHeight.RevisionNumber,
	//			RevisionHeight: defaultTimeoutHeight.RevisionHeight,
	//		},
	//		TimeoutTimestamp: disabledTimeoutTimestamp,
	//		Data:             mintFungibleTokenPacketData,
	//	},
	//
	//	acknowledgement.Acknowledgement())
}

func QueryProofs(t *testing.T) {
	clientId := clientIdA

	data, err := ethClient.StorageAt(context.Background(), ibc.ContractAddress, ibc.ClientSequenceSlot, nil)
	require.NoError(t, err)
	t.Logf("Client seq storage data: %x\n", data)

	proof, err := subnetClient.GetProof(context.Background(), ibc.ContractAddress, []string{ibc.ClientSequenceSlot.Hex()}, nil)
	require.NoError(t, err)
	t.Logf("Client seq storage merkle tree proof: %+v\n", proof)

	clientStateBz, err := ethClient.StorageAt(context.Background(), ibc.ContractAddress, ibc.ClientStateSlot(clientId), nil)
	require.NoError(t, err)
	t.Logf("Client state storage data: %x\n", clientStateBz)
}

func createIbcClient(t *testing.T, ctx context.Context, enpoint *ibctesting.Endpoint, clientId string) {
	clientState, ok1 := enpoint.GetClientState().(*ibctm.ClientState)
	clientState.MaxClockDrift = 5 * time.Minute
	require.True(t, ok1)
	clientStateByte, err := clientState.Marshal()
	require.NoError(t, err)

	consensusState, ok2 := enpoint.GetConsensusState(clientState.LatestHeight).(*ibctm.ConsensusState)
	require.True(t, ok2)
	consensusStateByte, err := consensusState.Marshal()
	require.NoError(t, err)

	tx, err := ibcContract.CreateClient(auth, exported.Tendermint, clientStateByte, consensusStateByte)
	require.NoError(t, err)

	re, err := waitForReceiptAndGet(ctx, ethClient, tx)
	require.NoError(t, err)
	require.True(t, len(re.Logs) > 0)

	event, err := ibcContractFilterer.ParseClientCreated(*re.Logs[0])
	require.NoError(t, err)

	assert.Equal(t, clientId, event.ClientId)
}

func queryClientStateFromContract(t *testing.T, cliendId string) *ibctm.ClientState {
	clientStateByte, err := ibcContract.QueryClientState(nil, cliendId)
	require.NoError(t, err)

	var clientState ibctm.ClientState
	require.NoError(t, clientState.Unmarshal(clientStateByte))

	return &clientState
}

func updateIbcClientAfterFunc(t *testing.T, cliendId string, endpoint *ibctesting.Endpoint, fn func() error) {
	if fn != nil {
		require.NoError(t, fn())
	}

	clientState := queryClientStateFromContract(t, cliendId)

	header, err := endpoint.Chain.ConstructUpdateTMClientHeaderWithTrustedHeight(endpoint.Counterparty.Chain, cliendId, clientState.LatestHeight)
	require.NoError(t, err)

	msg, err := clienttypes.MarshalClientMessage(marshaler, exported.ClientMessage(header))
	require.NoError(t, err)

	tx, err := ibcContract.UpdateClient(auth, cliendId, msg)
	require.NoError(t, err)

	re, err := waitForReceiptAndGet(context.Background(), ethClient, tx)
	require.NoError(t, err)

	t.Logf("'%s' updated: %#v", cliendId, re.Logs)
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

func getRandomAddr() (common.Address, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}
