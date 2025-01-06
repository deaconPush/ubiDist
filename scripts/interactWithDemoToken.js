const hre = require("hardhat");
const { ethers } = require("ethers");

async function main() {
  const [owner] = await hre.ethers.getSigners();
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
  // balanceOf owner
  const ownerBalance = await demoToken.connect(owner).balanceOf(ownerAddress);  
  console.log("Owner at " + ownerAddress + " has a " + tokenSymbol +  " balance of:" + ownerBalance);    
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
