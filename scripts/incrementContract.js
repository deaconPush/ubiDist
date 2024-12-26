const hre = require("hardhat");

async function main() {
const [owner, otherAccount] = await hre.ethers.getSigners();

console.log("address of owner: ", owner.address);
console.log("address of otherAccount: ", otherAccount.address);

const incrementContractFactory = await hre.ethers.getContractFactory("Increment");
const incrementContractAddress = "0x5fbdb2315678afecb367f032d93f642f64180aa3";
const incrementContract = incrementContractFactory.attach(incrementContractAddress);

// connect to the contract with different accounts
// await incrementContract.connect(owner).increment();
// await incrementContract.connect(owner).increment();

// this will fail because the other account is not the owner of the contract
// await incrementContract.connect(otherAccount).decrement();

// The counter value should be 1 instead of 0
let counterValue = await incrementContract.get();

console.log("Counter value: ", counterValue);
}

main()
 .then(() => process.exit(0))
 .catch((error) => {
   console.error(error);
   process.exit(1);
 });
