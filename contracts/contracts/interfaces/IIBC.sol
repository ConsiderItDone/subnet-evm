//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

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

struct Height {
    uint64 revisionNumber;
    uint64 revisionHeight;
}

interface IIBC {
  event ClientCreated(string clientId);
  event ConnectionCreated(string clientId, string connectionId);
  event ChannelCreated(string clientId, string connectionId);
  event PacketSent(
    bytes data,
    string timeoutHeight,
    uint64 timeoutTimestamp,
    uint64 sequence,
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    int32 channelOrdering
  );
  event PacketReceived(
    bytes data,
    string timeoutHeight,
    uint64 timeoutTimestamp,
    uint64 sequence,
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    int32 channelOrdering
  );
  event AcknowledgementWritten(
    bytes data,
    string timeoutHeight,
    uint64 timeoutTimestamp,
    uint64 sequence,
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    bytes ack,
    string connectionID
  );
  event AcknowledgePacket(
    string timeoutHeight,
    uint64 timeoutTimestamp,
    uint64 sequence,
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    int32 channelOrdering,
    string connectionID
  );
  event TimeoutPacket(
    string timeoutHeight,
    uint64 timeoutTimestamp,
    uint64 sequence,
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    int32 channelOrdering,
    string connectionID
  );


  struct MsgRecvPacket {
    Packet packet;
    bytes proofCommitment;
    Height proofHeight;
    string signer;
  }

  struct MsgAcknowledgement {
    Packet packet;
    bytes acknowledgement;
    bytes proofAcked;
    Height proofHeight;
    string signer;

  }

  struct MsgTimeoutOnClose {
    Packet packet;
    bytes proofUnreceived;
    bytes proofClose;
    Height proofHeight;
    uint64 nextSequenceRecv;
    string signer;
  }

  struct MsgTimeout {
    Packet packet;
    bytes proofUnreceived;
    Height proofHeight;
    uint64 nextSequenceRecv;
    string signer;
  }

  function RecvPacket(MsgRecvPacket memory message) external;

  function SendPacket(
      uint64 channelCapability,
      string memory sourcePort,
      string memory sourceChannel,
      Height memory timeoutHeight,
      uint64 timeoutTimestamp,
      bytes memory data
  ) external;

  function Acknowledgement(MsgAcknowledgement memory message) external;
  function TimeoutOnClose(MsgTimeoutOnClose memory message) external;
  function Timeout(MsgTimeout memory message) external;

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

  function bindPort(string memory portID) external;
}