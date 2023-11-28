//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Packet} from "./interfaces/IIBC.sol";

interface IBCApp {
  function OnRecvPacket(Packet memory packet, bytes memory relayer) external returns (bool);
  function OnAcknowledgementPacket(Packet memory packet, bytes memory, bytes memory) external returns (bool);
}