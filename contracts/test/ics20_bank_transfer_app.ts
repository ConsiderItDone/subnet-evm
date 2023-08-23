import { expect } from "chai"
import { ethers } from "hardhat"
import { Contract } from "ethers";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";

describe("ICS20BankTransferApp", function () {
  let ics20Bank: Contract;
  let ics20app: Contract;
  let owner: SignerWithAddress;
  let user: SignerWithAddress;

  before("deploy", async function () {
    const [[_owner, _user], ICS20Bank, ICS20BankTransferApp] = await Promise.all([
      ethers.getSigners(),
      ethers.getContractFactory("ICS20Bank"),
      ethers.getContractFactory("ICS20BankTransferApp"),
    ]);
    owner = _owner;
    user = _user;
    ics20Bank = await ICS20Bank.deploy();
    ics20app = await ICS20BankTransferApp.deploy(owner.address, ics20Bank.address);
  });

  it("chekc address", async function () {
    expect(ics20Bank.address).not.eq(ethers.constants.AddressZero);
    expect(ics20app.address).not.eq(ethers.constants.AddressZero);
  });

  it("escrow addr can alloca only owner", async function () {
    await expect(ics20app.connect(user).setChannelEscrowAddresses("channel", ethers.constants.AddressZero))
      .revertedWith("Ownable: caller is not the owner");
    await expect(ics20app.setChannelEscrowAddresses("", ethers.constants.AddressZero)).not.reverted;
  });

  it("can call only ibc", async function () {
    const tx = ics20app.connect(user).OnRecvPacket(
      {
        sequence: "0",
        sourcePort: "",
        sourceChannel: "",
        destinationPort: "",
        destinationChannel: "",
        data: "0x00",
        timeoutHeight: {
          revisionNumber: "0",
          revisionHeight: "0",
        },
        timeoutTimestamp: "0",
      },
      "0x00"
    );
    await expect(tx).revertedWith("BaseFungibleTokenApp: caller is not the IBC contract")
  });
});