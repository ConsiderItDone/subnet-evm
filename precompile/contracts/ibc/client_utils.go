package ibc

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
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
	//interfaceRegistry := cosmostypes.NewInterfaceRegistry()
	//std.RegisterInterfaces(interfaceRegistry)
	//ibctm.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	//marshaler := codec.NewProtoCodec(interfaceRegistry)

	statedb := opts.accessibleState.GetStateDB()

	clientState, err := GetClientState(statedb, opts.args.ClientID)
	if err != nil {
		return fmt.Errorf("can't get client state: %w", err)
	}

	if Status(opts.accessibleState, *clientState, opts.args.ClientID) != exported.Active {
		return fmt.Errorf("client is not active")
	}

	// ToDo: fix type checking and using
	//clientMsg, err := clienttypes.UnmarshalClientMessage(marshaler, opts.args.ClientMessage)
	//if err != nil {
	//	return fmt.Errorf("can't unmarshal client message: %w", err)
	//}
	//
	//if err := VerifyClientMessage(clientState, opts.args.ClientID, opts.accessibleState, clientMsg); err != nil {
	//	return err
	//}

	foundMisbehaviour := checkForMisbehaviour(*clientState, marshaler, clientMsg, opts.args.ClientID, opts.accessibleState)
	if foundMisbehaviour {
		clientState.FrozenHeight = ibctm.FrozenHeight
		if err := SetClientState(statedb, opts.args.ClientID, clientState); err != nil {
			return fmt.Errorf("can't update client state: %w", err)
		}

		topics, data, err := IBCABI.PackEvent(GeneratedTypeSubmitMisbehaviourIdentifier.RawName,
			opts.args.ClientID,
			clientState.ClientType(),
		)
		if err != nil {
			return fmt.Errorf("error packing event: %w", err)
		}
		blockNumber := opts.accessibleState.GetBlockContext().Number().Uint64()
		opts.accessibleState.GetStateDB().AddLog(ContractAddress, topics, data, blockNumber)
		return nil
	}

	consensusState, err := GetConsensusState(statedb, opts.args.ClientID, clientState.GetLatestHeight())
	if err != nil {
		return fmt.Errorf("can't get consensus state: %w", err)
	}

	header := clientMsg.(*ibctm.Header)
	clientState.LatestHeight = header.GetHeight().(clienttypes.Height)
	consensusState.Timestamp = header.GetTime()
	consensusState.Root = commitmenttypes.NewMerkleRoot(header.Header.GetAppHash())
	consensusState.NextValidatorsHash = header.Header.NextValidatorsHash

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

	if Status(opts.accessibleState, *clientState, opts.args.ClientID) != exported.Active {
		return fmt.Errorf("client is not active")
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

// CheckForMisbehaviour detects duplicate height misbehaviour and BFT time violation misbehaviour
// in a submitted Header message and verifies the correctness of a submitted Misbehaviour ClientMessage
func checkForMisbehaviour(
	cs ibctm.ClientState,
	cdc codec.BinaryCodec,
	msg exported.ClientMessage,
	clientID string,
	accessibleState contract.AccessibleState) bool {
	switch msg := msg.(type) {
	case *ibctm.Header:
		tmHeader := msg
		// consState := tmHeader.ConsensusState()

		// Check if the Client store already has a consensus state for the header's height
		// If the consensus state exists, and it matches the header then we return early
		// since header has already been submitted in a previous UpdateClient.
		existingConsState, err := GetConsensusState(accessibleState.GetStateDB(), clientID, tmHeader.GetHeight())
		if err != nil {
			return true
		}

		if existingConsState != nil {
			// This header has already been submitted and the necessary state is already stored
			// in client store, thus we can return early without further validation.
			if reflect.DeepEqual(existingConsState, tmHeader.ConsensusState()) { //nolint:gosimple
				return false
			}

			// A consensus state already exists for this height, but it does not match the provided header.
			// The assumption is that Header has already been validated. Thus we can return true as misbehaviour is present
			return true
		}

		// TODO
		// prevCons, err := GetPreviousConsensusState(clientStore, cdc, tmHeader.GetHeight())
		// if err != nil && !prevCons.Timestamp.Before(consState.Timestamp) {
		// 	return true
		// }
		// nextCons, err := GetNextConsensusState(clientStore, cdc, tmHeader.GetHeight())
		// if err != nil && !nextCons.Timestamp.After(consState.Timestamp) {
		// 	return true
		// }

	case *ibctm.Misbehaviour:
		// if heights are equal check that this is valid misbehaviour of a fork
		// otherwise if heights are unequal check that this is valid misbehavior of BFT time violation
		if msg.Header1.GetHeight().EQ(msg.Header2.GetHeight()) {
			blockID1, err := tmtypes.BlockIDFromProto(&msg.Header1.SignedHeader.Commit.BlockID)
			if err != nil {
				return false
			}

			blockID2, err := tmtypes.BlockIDFromProto(&msg.Header2.SignedHeader.Commit.BlockID)
			if err != nil {
				return false
			}

			// Ensure that Commit Hashes are different
			if !bytes.Equal(blockID1.Hash, blockID2.Hash) {
				return true
			}

		} else if !msg.Header1.SignedHeader.Header.Time.After(msg.Header2.SignedHeader.Header.Time) {
			// Header1 is at greater height than Header2, therefore Header1 time must be less than or equal to
			// Header2 time in order to be valid misbehaviour (violation of monotonic time).
			return true
		}
	}

	return false
}
