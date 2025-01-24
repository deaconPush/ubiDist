const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

const helloWorldModule = buildModule("HelloWorld", (m) => {
  const helloWorld = m.contract("HelloWorld", []);
  
  return { helloWorld };
});

module.exports = helloWorldModule;