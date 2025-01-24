const { ethers } = require("hardhat");
 
async function main() {
  const HelloWorldContract = await ethers.getContractFactory("HelloWorld");
  const helloWorld = await HelloWorldContract.deploy();
  
  const resp = await helloWorld.getMessage();
  console.log('Get message response', resp)
}

main()
 .then(() => process.exit(0))
 .catch((error) => {
   console.error(error);
   process.exit(1);
 });


