const { ethers } = require("hardhat");
const { generatePrivateKey, generateWallet } = require('./generateAddress');

async function main() {
    const demoOperationsContract = await ethers.getContractFactory("demoOperations");
    const demoOperations = await demoOperationsContract.deploy();

    await demoOperations.set(5);

    const resp = await demoOperations.get();
    console.log('Get message response', resp.toString())

    await demoOperations.externalFunction();

    const wallet = generateWallet(generatePrivateKey());

    await demoOperations.createAccount(wallet.address, 1000);
}

main()
 .then(() => process.exit(0))
 .catch((error) => {
   console.error(error);
   process.exit(1);
 });
