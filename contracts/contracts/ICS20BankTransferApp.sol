//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import {IICS20Bank} from "./ICS20Bank.sol";
import {IBCBaseFungibleTokenApp} from "./IBCBaseFungibleTokenApp.sol";

interface PortBinder{
  function bindPort(string memory portID) external;
}

contract ICS20BankTransferApp is Ownable, IBCBaseFungibleTokenApp {
  address ibcAddr;
  IICS20Bank bank;

  constructor(address _ibcAddr, IICS20Bank _bank) {
    ibcAddr = _ibcAddr;
    bank = _bank;
  }

  function ibcAddress() public view override returns (address) {
    return ibcAddr;
  }

  function setChannelEscrowAddresses(string memory chan, address chanAddr) external onlyOwner {
    channelEscrowAddresses[chan] = chanAddr;
  }

  function _transferFrom(
    address sender,
    address receiver,
    string memory denom,
    uint256 amount
  ) internal override returns (bool) {
    bank.transferFrom(sender, receiver, denom, amount);
    return true;
  }

  function _mint(address account, string memory denom, uint256 amount) internal override returns (bool) {
    bank.mint(account, denom, amount);
    return true;
  }

  function bindPort(address someAddr, string memory portId) external {
    PortBinder(someAddr).bindPort(portId);
  }
}
