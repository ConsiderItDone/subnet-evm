//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IIBC {
  // Create IBC Client
  function createClient(string memory clientType, bytes memory clientState, bytes memory consensusState) external returns (string memory clientID);
}