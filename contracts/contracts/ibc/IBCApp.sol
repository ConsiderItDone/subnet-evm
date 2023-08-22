//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Packet} from "../interfaces/IIBC.sol";

struct FungibleTokenPacketData {
  string denom;
  uint256 amount;
  string sender;
  address receiver;
  string memo;
}

abstract contract BaseIBCApp {
  mapping(string => address) channelEscrowAddresses;

  function OnRecvPacket(Packet memory packet, bytes memory) external returns (bytes memory) {
    FungibleTokenPacketData memory data = abi.decode(packet.data, (FungibleTokenPacketData));
    return
      _newAcknowledgement(
        _transferFrom(_getEscrowAddress(packet.sourceChannel), data.receiver, data.denom, data.amount)
      );
  }

  function _transferFrom(
    address sender,
    address receiver,
    string memory denom,
    uint256 amount
  ) internal virtual returns (bool);

  function _getEscrowAddress(string memory sourceChannel) internal view virtual returns (address) {
    address escrow = channelEscrowAddresses[sourceChannel];
    require(escrow != address(0));
    return escrow;
  }

  function _newAcknowledgement(bool success) internal pure virtual returns (bytes memory) {
    bytes memory acknowledgement = new bytes(1);
    if (success) {
      acknowledgement[0] = 0x01;
    } else {
      acknowledgement[0] = 0x00;
    }
    return acknowledgement;
  }
}


