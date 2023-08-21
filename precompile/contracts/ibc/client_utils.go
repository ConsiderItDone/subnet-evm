package ibc

import (
	"errors"
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/ethereum/go-ethereum/common"
)

func makeClientID(db contract.StateDB, clientType string) string {
	clientSeq := db.GetState(ContractAddress, ClientSequenceSlot).Big()
	clientID := fmt.Sprintf("%s-%d", clientType, clientSeq.Int64())
	clientSeq.Add(clientSeq, common.Big1)
	db.SetState(ContractAddress, ClientSequenceSlot, common.BigToHash(clientSeq))
	return clientID
}

func _createClient(opts *callOpts[CreateClientInput]) (string, error) {
	statedb := opts.accessibleState.GetStateDB()
	clientType := opts.args.ClientType

	// supports only Tendermint for now
	if clientType != exported.Tendermint {
		return "", ErrWrongClientType
	}

	// generate clientID
	clientId := makeClientID(statedb, clientType)

	clientState := new(ibctm.ClientState)
	if err := clientState.Unmarshal(opts.args.ClientState); err != nil {
		return "", fmt.Errorf("error unmarshalling client state: %w", err)
	}
	if err := SetClientState(statedb, clientId, clientState); err != nil {
		return "", fmt.Errorf("error storing client state: %w", err)
	}

	consensusState := new(ibctm.ConsensusState)
	if err := consensusState.Unmarshal(opts.args.ConsensusState); err != nil {
		return "", fmt.Errorf("error unmarshalling consensus state: %w", err)
	}
	if err := SetConsensusState(statedb, clientId, clientState.GetLatestHeight(), consensusState); err != nil {
		return "", fmt.Errorf("error storing consensus state: %w", err)
	}

	if err := AddLog(opts.accessibleState, "ClientCreated", clientId); err != nil {
		return "", fmt.Errorf("error packing event: %w", err)
	}

	return clientId, nil
}

func _updateClient(opts *callOpts[UpdateClientInput]) error {
	statedb := opts.accessibleState.GetStateDB()

	clientState, err := GetClientState(statedb, opts.args.ClientID)
	if err != nil {
		return fmt.Errorf("can't get client state: %w", err)
	}

	consensusState, err := GetConsensusState(statedb, opts.args.ClientID, clientState.GetLatestHeight())
	if err != nil {
		return fmt.Errorf("can't get consensus state: %w", err)
	}

	clientMessage := &ibctm.Header{}
	if err := clientMessage.Unmarshal(opts.args.ClientMessage); err != nil {
		return fmt.Errorf("error unmarshalling client state file: %w", err)
	}

	clientState.LatestHeight = clientMessage.GetHeight().(clienttypes.Height)
	consensusState.Timestamp = clientMessage.GetTime()
	consensusState.Root = commitmenttypes.NewMerkleRoot(clientMessage.Header.GetAppHash())
	consensusState.NextValidatorsHash = clientMessage.Header.NextValidatorsHash

	if err := SetClientState(statedb, opts.args.ClientID, clientState); err != nil {
		return fmt.Errorf("can't update client state: %w", err)
	}

	if err := SetConsensusState(statedb, opts.args.ClientID, clientState.GetLatestHeight(), consensusState); err != nil {
		return fmt.Errorf("can't update consensus state: %w", err)
	}

	return nil
}

func _upgradeClient(opts *callOpts[UpgradeClientInput]) error {
	interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	statedb := opts.accessibleState.GetStateDB()

	upgradedClient := new(ibctm.ClientState)
	if err := upgradedClient.Unmarshal(opts.args.UpgradedClien); err != nil {
		return fmt.Errorf("error unmarshalling upgraded client: %w", err)
	}

	upgradedConsState := new(ibctm.ConsensusState)
	if err := upgradedConsState.Unmarshal(opts.args.UpgradedConsState); err != nil {
		return fmt.Errorf("error unmarshalling upgraded ConsensusState: %w", err)
	}

	clientState, err := GetClientState(statedb, opts.args.ClientID)
	if err != nil {
		return fmt.Errorf("can't get client state: %w", err)
	}

	consensusState, err := GetConsensusState(statedb, opts.args.ClientID, clientState.GetLatestHeight())
	if err != nil {
		return fmt.Errorf("can't get consensus state: %w", err)
	}

	if len(clientState.UpgradePath) == 0 {
		return errors.New("cannot upgrade client, no upgrade path set")
	}

	// last height of current counterparty chain must be client's latest height
	lastHeight := clientState.GetLatestHeight()

	if !upgradedClient.GetLatestHeight().GT(lastHeight) {
		return fmt.Errorf("upgraded client height %s must be at greater than current client height %s", upgradedClient.GetLatestHeight(), lastHeight)
	}

	// unmarshal proofs
	var merkleProofClient, merkleProofConsState commitmenttypes.MerkleProof
	if err := marshaler.Unmarshal(opts.args.ProofUpgradeClient, &merkleProofClient); err != nil {
		return fmt.Errorf("could not unmarshal client merkle proof: %v", err)
	}
	if err := marshaler.Unmarshal(opts.args.ProofUpgradeConsState, &merkleProofConsState); err != nil {
		return fmt.Errorf("could not unmarshal consensus state merkle proof: %v", err)
	}

	// Verify client proof
	client := upgradedClient.ZeroCustomFields()
	bz, err := marshaler.MarshalInterface(client)
	if err != nil {
		return fmt.Errorf("could not marshal client state: %v", err)
	}
	// copy all elements from upgradePath except final element
	clientPath := make([]string, len(clientState.UpgradePath)-1)
	copy(clientPath, clientState.UpgradePath)

	// append lastHeight and `upgradedClient` to last key of upgradePath and use as lastKey of clientPath
	// this will create the IAVL key that is used to store client in upgrade store
	lastKey := clientState.UpgradePath[len(clientState.UpgradePath)-1]
	appendedKey := fmt.Sprintf("%s/%d/%s", lastKey, lastHeight.GetRevisionHeight(), upgradetypes.KeyUpgradedClient)

	clientPath = append(clientPath, appendedKey)

	// construct clientState Merkle path
	upgradeClientPath := commitmenttypes.NewMerklePath(clientPath...)

	if err := merkleProofClient.VerifyMembership(clientState.ProofSpecs, consensusState.GetRoot(), upgradeClientPath, bz); err != nil {
		return fmt.Errorf("client state proof failed. Path: %s, err: %v", upgradeClientPath.Pretty(), err)
	}

	// Verify consensus state proof
	bz, err = marshaler.MarshalInterface(upgradedConsState)
	if err != nil {
		return fmt.Errorf("could not marshal consensus state: %v", err)
	}

	// copy all elements from upgradePath except final element
	consPath := make([]string, len(clientState.UpgradePath)-1)
	copy(consPath, clientState.UpgradePath)

	// append lastHeight and `upgradedClient` to last key of upgradePath and use as lastKey of clientPath
	// this will create the IAVL key that is used to store client in upgrade store
	lastKey = clientState.UpgradePath[len(clientState.UpgradePath)-1]
	appendedKey = fmt.Sprintf("%s/%d/%s", lastKey, lastHeight.GetRevisionHeight(), upgradetypes.KeyUpgradedConsState)

	consPath = append(consPath, appendedKey)
	// construct consensus state Merkle path
	upgradeConsStatePath := commitmenttypes.NewMerklePath(consPath...)
	if err := merkleProofConsState.VerifyMembership(clientState.ProofSpecs, consensusState.GetRoot(), upgradeConsStatePath, bz); err != nil {
		return fmt.Errorf("consensus state proof failed. Path: %s", upgradeConsStatePath.Pretty())
	}

	newClientState := ibctm.NewClientState(
		upgradedClient.ChainId, clientState.TrustLevel, clientState.TrustingPeriod, upgradedClient.UnbondingPeriod,
		clientState.MaxClockDrift, upgradedClient.LatestHeight, upgradedClient.ProofSpecs, upgradedClient.UpgradePath,
	)
	if err := newClientState.Validate(); err != nil {
		return fmt.Errorf("updated client state failed basic validation")
	}
	newConsState := ibctm.NewConsensusState(
		upgradedConsState.Timestamp, commitmenttypes.NewMerkleRoot([]byte(ibctm.SentinelRoot)), upgradedConsState.NextValidatorsHash,
	)

	if err := SetClientState(statedb, opts.args.ClientID, newClientState); err != nil {
		return fmt.Errorf("error storing client state: %w", err)
	}
	if err := SetConsensusState(statedb, opts.args.ClientID, newClientState.GetLatestHeight(), newConsState); err != nil {
		return fmt.Errorf("error storing consensus state: %w", err)
	}
	return nil
}
