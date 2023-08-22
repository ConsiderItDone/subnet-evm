//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IICS20Bank} from "./ICS20Bank.sol";
import {BaseIBCApp} from "./IBCApp.sol";

contract ICS20TransferBank is BaseIBCApp {
  IICS20Bank bank;

  constructor(IICS20Bank _bank) {
    bank = _bank;
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
}
