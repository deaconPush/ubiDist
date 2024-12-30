const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

const rocketModule = buildModule("Rocket", (m) => {
  const rocket = m.contract("Rocket", ["Saturn V"]);
  m.call(rocket , "launch", []);
  
  return { rocket };
});

module.exports = rocketModule;