const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

const demoTokenModule = buildModule("DemoToken", (m) => {
    const DemoToken = m.contract("DemoToken", []);
    
    return { DemoToken };
  });
  
  module.exports = demoTokenModule;