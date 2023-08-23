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
    await expect(ics20Bank.setOperator(ics20app.address)).not.reverted;
    await expect(ics20Bank.setOperator(owner.address)).not.reverted;
  });

  it("chekc address", async function () {
    expect(ics20Bank.address).not.eq(ethers.constants.AddressZero);
    expect(ics20app.address).not.eq(ethers.constants.AddressZero);
  });

  it("escrow addr can alloc only owner", async function () {
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

  it("success mint", async function () {
    const codec = new ethers.utils.AbiCoder();
    const data = codec.encode(
      ["tuple(string, uint256, string, address, string)"],
      [["ETH", ethers.BigNumber.from("1000"), owner.address, user.address, "some memo"]],
    )

    const onRecvPacketCall = ics20app.OnRecvPacket(
      {
        sequence: "0",
        sourcePort: "port",
        sourceChannel: "channel",
        destinationPort: "",
        destinationChannel: "",
        data: data,
        timeoutHeight: {
          revisionNumber: "0",
          revisionHeight: "0",
        },
        timeoutTimestamp: "0",
      },
      "0x00"
    );

    await expect(onRecvPacketCall)
      .to.emit(ics20Bank, "Transfer")
      .withArgs(ethers.constants.AddressZero, user.address, ethers.BigNumber.from("1000"));
  });

  it("success transfer", async function () {
    const codec = new ethers.utils.AbiCoder();
    const data = codec.encode(
      ["tuple(string, uint256, string, address, string)"],
      [["port/channel/ETH", ethers.BigNumber.from("1000"), user.address, owner.address, "some memo"]],
    )
    await expect(ics20app.setChannelEscrowAddresses("channel", user.address)).not.reverted;
    await expect(ics20Bank.mint(user.address,  "port/channel/ETH", 1000)).not.reverted;
    const onRecvPacketCall = ics20app.OnRecvPacket(
      {
        sequence: "0",
        sourcePort: "port",
        sourceChannel: "channel",
        destinationPort: "",
        destinationChannel: "",
        data: data,
        timeoutHeight: {
          revisionNumber: "0",
          revisionHeight: "0",
        },
        timeoutTimestamp: "0",
      },
      "0x00"
    );
    await expect(onRecvPacketCall)
      .to.emit(ics20Bank, "Transfer")
      .withArgs(user.address, owner.address, ethers.BigNumber.from("1000"));
  });
});