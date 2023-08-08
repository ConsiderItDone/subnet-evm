//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IIBC {
  event ClientCreated(string clientId);
  event ConnectionCreated(string clientId, string connectionId);
  event ChannelCreated(string clientId, string connectionId);

  struct Height {
      uint64 revisionNumber;
      uint64 revisionHeight;
  }

  struct Packet {
      uint64 sequence;
      string sourcePort;
      string sourceChannel;
      string destinationPort;
      string destinationChannel;
      bytes data;
      Height timeoutHeight;
      uint64 timeoutTimestamp;
  }

  function OnRecvPacket(Packet memory packet, bytes memory relayer) external;

  function SendPacket(
      uint64 channelCapability,
      string memory sourcePort,
      string memory sourceChannel,
      Height memory timeoutHeight,
      uint64 timeoutTimestamp,
      bytes memory data
  ) external;

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

  function chanOpenInit(string memory portID, bytes memory channel) external;

  function chanOpenTry(
    string memory portID,
    bytes memory channel,
    string memory counterpartyVersion,
    bytes memory proofInit,
    bytes memory proofHeight
  ) external returns (string memory channelID);

  function channelOpenAck(
    string memory portID,
    string memory channelID,
    string memory counterpartyChannelID,
    string memory counterpartyVersion,
    bytes memory proofTry,
    bytes memory proofHeight
  ) external;

  function channelOpenConfirm(
    string memory portID,
    string memory channelID,
    bytes memory proofAck,
    bytes memory proofHeight
  ) external;

  function channelCloseInit(string memory portID, string memory channelID) external;

  function channelCloseConfirm(
    string memory portID,
    string memory channelID,
    bytes memory proofInit,
    bytes memory proofHeight
  ) external;
}
