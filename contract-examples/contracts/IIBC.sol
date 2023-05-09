//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IIBC {
  // Create IBC Client
  function createClient(string clientType, bytes clientStateBytes, bytes consensusStateBytes) external;

  // Update IBC Client
  function updateClient(int64 clientId, bytes message) external;

  // Upgrade IBC Client
  function upgradeClient(int64 clientId, bytes upgradedClient, bytes upgradedConsState, bytes proofUpgradeClient, bytes proofUpgradeConsState) external;
}
