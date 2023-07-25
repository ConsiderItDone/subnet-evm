//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IIBC {
  // Create IBC Client
  function createClient(
    string memory clientType,
    bytes memory clientState,
    bytes memory consensusState
  ) external returns (string memory clientID);

  function updateClient(string memory clientID, bytes memory clientMessage) external;

  function upgradeClient(
    string memory clientID,
    bytes memory upgradePath,
    bytes memory upgradedClien,
    bytes memory upgradedConsState,
    bytes memory proofUpgradeClient,
    bytes memory proofUpgradeConsState
  ) external;

  function connOpenInit(
    string memory clientID,
    bytes memory counterparty,
    bytes memory version,
    uint32 delayPeriod
  ) external returns (string memory connectionID);

  function connOpenTry(
    bytes memory counterparty,
    uint32 delayPeriod,
    string memory clientID,
    bytes memory clientState,
    bytes memory counterpartyVersions,
    bytes memory proofInit,
    bytes memory proofClient,
    bytes memory proofConsensus,
    bytes memory proofHeight,
    bytes memory consensusHeight
  ) external returns (string memory connectionID);

  function connOpenAck(
    string memory connectionID,
    bytes memory clientState,
    bytes memory version,
    bytes memory counterpartyConnectionID,
    bytes memory proofTry,
    bytes memory proofClient,
    bytes memory proofConsensus,
    bytes memory proofHeight,
    bytes memory consensusHeight
  ) external;

  function connOpenConfirm(string memory connectionID, bytes memory proofAck, bytes memory proofHeight) external;
}
