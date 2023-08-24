//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma experimental ABIEncoderV2;

struct FungibleTokenPacketData {
  string denom;
  uint256 amount;
  string sender;
  address receiver;
  string memo;
}

contract ICS20PacketDataCodec {
  function encode(FungibleTokenPacketData memory data) public pure returns (bytes memory) {
    return abi.encode(data);
  }

  function decode(bytes memory rawdata) public pure returns (FungibleTokenPacketData memory) {
    return abi.decode(rawdata, (FungibleTokenPacketData));
  }
}
