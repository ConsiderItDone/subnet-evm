//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/Context.sol";
import "solidity-stringutils/src/strings.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import {Packet} from "./interfaces/IIBC.sol";
import {IBCApp} from "./IBCApp.sol";

struct FungibleTokenPacketData {
  string denom;
  uint256 amount;
  bytes sender;
  bytes receiver;
  bytes memo;
}

abstract contract IBCBaseFungibleTokenApp is IBCApp, Context {
  using strings for *;
  using BytesLib for bytes;

  mapping(string => address) channelEscrowAddresses;

  /**
   * @dev Throws if called by any account other than the IBC contract.
   */
  modifier onlyIBC() {
    require(ibcAddress() == _msgSender(), "BaseFungibleTokenApp: caller is not the IBC contract");
    _;
  }

  /**
   * @dev Returns the address of the IBC contract.
   */
  function ibcAddress() public view virtual returns (address);

  function OnRecvPacket(Packet memory packet, bytes memory) external override onlyIBC returns (bool) {
    FungibleTokenPacketData memory data = abi.decode(packet.data, (FungibleTokenPacketData));
    strings.slice memory denom = data.denom.toSlice();
    strings.slice memory trimedDenom = data.denom.toSlice().beyond(
      _makeDenomPrefix(packet.sourcePort, packet.sourceChannel)
    );
    if (!denom.equals(trimedDenom)) {
      return
        _transferFrom(_getEscrowAddress(packet.sourceChannel), data.receiver.toAddress(0), data.denom, data.amount);
    }
    string memory prefixedDenom = _makeDenomPrefix(packet.destinationPort, packet.destinationChannel).concat(denom);
    return _mint(data.receiver.toAddress(0), prefixedDenom, data.amount);
  }

  function OnAcknowledgementPacket(
    Packet memory packet,
    bytes memory ack,
    bytes memory
  ) external virtual override onlyIBC returns (bool) {
    if (!_isSuccessAcknowledgement(ack)) {
      FungibleTokenPacketData memory data = abi.decode(packet.data, (FungibleTokenPacketData));
      _refundTokens(data, packet.sourcePort, packet.sourceChannel);
    }
    return true;
  }

  function _isSuccessAcknowledgement(bytes memory ack) internal pure virtual returns (bool) {
    require(ack.length == 1);
    return uint8(ack[0]) == 1;
  }

  function _refundTokens(
    FungibleTokenPacketData memory data,
    string memory sourcePort,
    string memory sourceChannel
  ) internal virtual {
    if (!data.denom.toSlice().startsWith(_makeDenomPrefix(sourcePort, sourceChannel))) {
      require(_transferFrom(_getEscrowAddress(sourceChannel), data.sender.toAddress(0), data.denom, data.amount));
    } else {
      require(_mint(data.sender.toAddress(0), data.denom, data.amount));
    }
  }

  function _mint(address account, string memory denom, uint256 amount) internal virtual returns (bool);

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

  function _makeDenomPrefix(
    string memory port,
    string memory channel
  ) internal pure virtual returns (strings.slice memory) {
    return
      port
        .toSlice()
        .concat("/".toSlice())
        .toSlice()
        .concat(channel.toSlice())
        .toSlice()
        .concat("/".toSlice())
        .toSlice();
  }
}