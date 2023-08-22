import { expect } from "chai"
import { ethers } from "hardhat"
import { Contract } from "ethers";

describe("ICS20TransferBank", function () {
  let ics20Bank: Contract;
  let ics20tb: Contract;

  beforeEach("deploy", async function () {
    const [ICS20Bank, ICS20TransferBank] = await Promise.all([
      ethers.getContractFactory("ICS20Bank"),
      ethers.getContractFactory("ICS20TransferBank"),
    ]);
    ics20Bank = await ICS20Bank.deploy();
    ics20tb = await ICS20TransferBank.deploy(ics20Bank.address);
  });

  it("chekc address", async function () {
    expect(ics20Bank.address).not.eq(ethers.constants.AddressZero);
    expect(ics20tb.address).not.eq(ethers.constants.AddressZero);
  })
});