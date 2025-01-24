const hre = require("hardhat");

async function main() {
const [owner, otherAccount] = await hre.ethers.getSigners();
const incrementContractFactory = await hre.ethers.getContractFactory("Increment");
// retrieve the contract address from running the ignition script for the contract
const incrementContractAddress = "";
const incrementContract = incrementContractFactory.attach(incrementContractAddress);
// connect to the contract with different accounts
await incrementContract.connect(owner).increment();
// this will fail because the other account is not the owner of the contract
try {
  await incrementContract.connect(otherAccount).increment();
} catch (error) {
  console.error("Failed to increment the counter from account different from the owner");
}
// The counter value should be 6 instead of 7 because the other account is not the owner of the contract 
// account different from the owner can read the counter value
const counterValueFromOtherAccount = await incrementContract.connect(otherAccount).getCounter();
console.log("Counter value from other account: ", counterValueFromOtherAccount.toString());
}

main()
 .then(() => process.exit(0))
 .catch((error) => {
   console.error(error);
   process.exit(1);
 });
