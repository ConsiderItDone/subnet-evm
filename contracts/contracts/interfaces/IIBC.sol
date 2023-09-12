//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

struct Packet {
  uint sequence;
  string sourcePort;
  string sourceChannel;
  string destinationPort;
  string destinationChannel;
  bytes data;
  Height timeoutHeight;
  uint timeoutTimestamp;
}

struct Height {
  uint revisionNumber;
  uint revisionHeight;
}

interface IIBC {
  event ClientCreated(string clientId);
  event ConnectionCreated(string clientId, string connectionId);
  event ChannelCreated(string clientId, string connectionId);
  event PacketSent(
    bytes data,
    string timeoutHeight,
    uint timeoutTimestamp,
    uint sequence,
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    int32 channelOrdering
  );
  event PacketReceived(
    bytes data,
    string timeoutHeight,
    uint timeoutTimestamp,
    uint sequence,
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    int32 channelOrdering
  );
  event AcknowledgementWritten(
    bytes data,
    string timeoutHeight,
    uint timeoutTimestamp,
    uint sequence,
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    bytes ack,
    string connectionID
  );
  event AcknowledgePacket(
    string timeoutHeight,
    uint timeoutTimestamp,
    uint sequence,
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    int32 channelOrdering,
    string connectionID
  );
  event TimeoutPacket(
    string timeoutHeight,
    uint timeoutTimestamp,
    uint sequence,
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    int32 channelOrdering,
    string connectionID
  );
  event AcknowledgementError(
    bytes data,
    string timeoutHeight,
    uint timeoutTimestamp,
    uint sequence,
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    string error
  );
  event TypeSubmitMisbehaviour(
    string clientID,
    string clientType
  );
  event TypeChannelClosed(
    string sourcePort,
    string sourceChannel,
    string destPort,
    string destChannel,
    string ConnectionID,
    string ChannelOrdering
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

  function recvPacket(MsgRecvPacket memory message) external;

  function sendPacket(
    uint channelCapability,
    string memory sourcePort,
    string memory sourceChannel,
    Height memory timeoutHeight,
    uint timeoutTimestamp,
    bytes memory data
  ) external;

  function acknowledgement(
    Packet memory packet,
    bytes memory acknowledgement,
    bytes memory proofAcked,
    Height memory proofHeight,
    string memory signer
  ) external;

  function timeoutOnClose(
    Packet memory packet,
    bytes memory proofUnreceived,
    bytes memory proofClose,
    Height memory proofHeight,
    uint nextSequenceRecv,
    string memory signer
  ) external;

  function timeout(
    Packet memory packet,
    bytes memory proofUnreceived,
    Height memory proofHeight,
    uint nextSequenceRecv,
    string memory signer
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

  function OnRecvPacket(Packet memory packet, bytes memory Relayer) external;

  function OnTimeout(
	  Packet  memory packet,
	  bytes memory Relayer
  ) external;

  function OnTimeoutOnClose(
	  Packet  memory packet,
	  bytes memory Relayer
  ) external;

  function OnAcknowledgementPacket(
    Packet memory packet, 
    bytes memory ack, 
    bytes memory
  ) external;

    // query methods
  function queryClientState(string memory clientId) external returns (bytes memory);

  function queryConsensusState(string memory clientId) external returns (bytes memory);

  function queryConnection(string memory connectionID) external returns (bytes memory);

  function queryChannel(string memory portID, string memory channelID) external returns (bytes memory);

  function queryPacketCommitment(string memory portID, string memory channelID, uint sequence) external returns (bytes memory);

  function queryPacketAcknowledgement(string memory portID, string memory channelID, uint sequence) external returns (bytes memory);
}