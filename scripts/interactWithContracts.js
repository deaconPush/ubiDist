const hre = require("hardhat");

async function main() { 
    const rocketContractFactory  = await hre.ethers.getContractFactory("Rocket");
    const helloWorldContractFactory = await hre.ethers.getContractFactory("HelloWorld");  
    // retrieve the contract addresses from running the ignition scripts for the contracts
    // the node should be running before executing this script
    const rocketContractAddress = "";
    const helloWorldContractAddress = "";
    const rocketContract = rocketContractFactory.attach(rocketContractAddress);
    const helloWorldContract = helloWorldContractFactory.attach(helloWorldContractAddress);
    // interact with the contracts
    await rocketContract.launch();
    const message = await helloWorldContract.getMessage();
    console.log('message from helloWorld contract: ', message);
}

main()
 .then(() => process.exit(0))
 .catch((error) => {
   console.error(error);
   process.exit(1);
 });
