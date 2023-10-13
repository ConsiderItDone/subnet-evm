//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IBCApp.sol";

contract IBCCounterApp is IBCApp {
    int64 public counter = 0;
    function OnRecvPacket(Packet memory packet, bytes memory) override external {
        int64 value = abi.decode(packet.data, (int64));
        counter += value;
    }
}