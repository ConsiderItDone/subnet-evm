//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Packet} from "./interfaces/IIBC.sol";

abstract contract IBCApp {
  function OnRecvPacket(Packet memory packet, bytes memory relayer) virtual external returns (bool);
  function OnAcknowledgementInput(Packet calldata packet, bytes calldata acknowledgement, bytes calldata) external;
}