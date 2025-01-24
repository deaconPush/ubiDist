const hre = require("hardhat");
const { ethers } = require("ethers");

async function main() {
  const [owner, account1] = await hre.ethers.getSigners();
  const ownerAddress = owner.address;
  const demoTokenFactory = await hre.ethers.getContractFactory("DemoToken");
  const demoTokenContractAddress = "";
  const demoToken = demoTokenFactory.attach(demoTokenContractAddress);
  // Get Token name
  const tokenName = await demoToken.connect(owner).name();
  // Get Token symbol  
  const tokenSymbol = await demoToken.connect(owner).symbol();
  console.log("Token name: ", tokenName);
  // Get token decimals
  const decimals = await demoToken.connect(owner).decimals();
  // Get total supply and convert to whole number
  const totalSupply = await demoToken.connect(owner).totalSupply();
  console.log("Total supply: ", ethers.formatUnits(totalSupply, decimals));
  // Get the balcance of the owner
  let ownerBalance = await demoToken.balanceOf(ownerAddress);
  console.log(`Contract owner at ${ownerAddress} has a ${tokenSymbol} balance of ${ethers.formatUnits(ownerBalance, decimals)}`)
  // transfer tokens to another address
  // transfer(address to, uint256 amount)
  const transferAmount = 100000;
  console.log(`Transfering ${transferAmount} ${tokenSymbol} to ${account1.address}`);
  await demoToken.transfer(account1.address, ethers.parseUnits(transferAmount.toString(), decimals));
  // check balance of owner and recipient after transfer
  ownerBalance = await demoToken.balanceOf(ownerAddress);
  console.log(`Balance of owner (${ownerAddress}): ${ethers.formatUnits(ownerBalance, decimals)} ${tokenSymbol}`);
  let recipientBalance = await demoToken.balanceOf(account1.address);
  console.log(`Balance of recipient (${account1.address}): ${ethers.formatUnits(recipientBalance, decimals)} ${tokenSymbol}`)
  // Try to spend tokens of another account without approval
  let spendAmount = 1000;
  try {
    await demoToken.transferFrom(account1.address, ownerAddress, ethers.parseUnits(spendAmount.toString(), decimals));
  } catch(error){
    console.error("Failed to spend tokens without approval");
  }
  // Approve spender to spend tokens on behalf of owner
  // approve(address spender, uint256 amount)
  const approveAmmount = 10000;
  const account1Contract = demoToken.connect(account1);
  await account1Contract.approve(ownerAddress, ethers.parseUnits(approveAmmount.toString(), decimals));
  // Get allowance of spender
  // allowance(address owner, address spender)
  let allowance = await demoToken.allowance(account1.address, ownerAddress);
  console.log(`allowance of owner ${ethers.formatUnits(allowance, decimals)} ${tokenSymbol}`)
  await demoToken.transferFrom(account1.address, ownerAddress, ethers.parseUnits(spendAmount.toString(), decimals));
  ownerBalance = await demoToken.balanceOf(ownerAddress);
  recipientBalance = await demoToken.balanceOf(account1.address);
  console.log(`new balance of owner (${ownerAddress}): ${ethers.formatUnits(ownerBalance, decimals)} ${tokenSymbol}`);
  console.log(`new balance of recipient (${account1.address}): ${ethers.formatUnits(recipientBalance, decimals)} ${tokenSymbol}`)
}
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
