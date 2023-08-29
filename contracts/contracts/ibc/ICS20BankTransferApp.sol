//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import {IICS20Bank} from "./ICS20Bank.sol";
import {BaseFungibleTokenApp} from "./BaseFungibleTokenApp.sol";

contract ICS20BankTransferApp is Ownable, BaseFungibleTokenApp {
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
    try bank.transferFrom(sender, receiver, denom, amount) {
      return true;
    } catch (bytes memory) {
      return false;
    }
  }

  function _mint(address account, string memory denom, uint256 amount) internal override returns (bool) {
    try bank.mint(account, denom, amount) {
      return true;
    } catch (bytes memory) {
      return false;
    }
  }
}
