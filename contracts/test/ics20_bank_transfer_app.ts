import { expect } from "chai"
import { ethers } from "hardhat"
import { Contract } from "ethers";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";

function encodeFungibleTokenPacketData(
  denom: string,
  amount: string,
  sender: string,
  receiver: string,
  memo: string
): string {
  const codec = new ethers.utils.AbiCoder();
  return codec.encode(
    ["tuple(string, uint256, string, address, string)"],
    [[denom, ethers.BigNumber.from(amount), sender, receiver, memo]]
  );
}

function makeOnRecvPacketData(data?: object): object {
  const def = {
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
  };
  return { ...def, ...(data ?? {}) }
}

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
    await expect(ics20app.connect(user).OnRecvPacket(makeOnRecvPacketData(), "0x00"))
      .revertedWith("BaseFungibleTokenApp: caller is not the IBC contract");
  });

  it("mint success", async function () {
    const data = makeOnRecvPacketData({
      sourcePort: "port",
      sourceChannel: "channel",
      data: encodeFungibleTokenPacketData(
        "ETH",
        "1000",
        owner.address,
        user.address,
        "some memo"
      ),
    });
    await expect(ics20app.OnRecvPacket(data, "0x00"))
      .to.emit(ics20Bank, "Transfer")
      .withArgs(ethers.constants.AddressZero, user.address, ethers.BigNumber.from("1000"));
  });

  it("transfer channel doesn't have address", async function () {
    const data = makeOnRecvPacketData({
      sourcePort: "port",
      sourceChannel: "channel",
      data: encodeFungibleTokenPacketData(
        "port/channel/ETH",
        "1000",
        user.address,
        owner.address,
        "some memo"
      ),
    });
    await expect(ics20app.OnRecvPacket(data, "0x00")).to.be.reverted;
  });

  it("transfer escrow doesn't have balance", async function () {
    const data = makeOnRecvPacketData({
      sourcePort: "port",
      sourceChannel: "channel",
      data: encodeFungibleTokenPacketData(
        "port/channel/BTC",
        "1000",
        user.address,
        owner.address,
        "some memo"
      ),
    });
    await expect(ics20app.setChannelEscrowAddresses("channel", user.address)).not.reverted;
    const tx = await ics20app.OnRecvPacket(data, "0x00");
    expect(await tx.wait()).to.deep.include({ logs: [] });
  });

  it("transfer success", async function () {
    const data = makeOnRecvPacketData({
      sourcePort: "port",
      sourceChannel: "channel",
      data: encodeFungibleTokenPacketData(
        "port/channel/ETH",
        "1000",
        user.address,
        owner.address,
        "some memo"
      ),
    });
    await expect(ics20Bank.mint(user.address, "port/channel/ETH", 1000)).not.reverted;
    await expect(ics20app.OnRecvPacket(data, "0x00"))
      .to.emit(ics20Bank, "Transfer")
      .withArgs(user.address, owner.address, ethers.BigNumber.from("1000"));
  });
});