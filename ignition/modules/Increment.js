const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

const incrementModule = buildModule("Increment", (m) => {
    const increment = m.contract("Increment", [5, "counter"]);
    return { increment };
});

module.exports = incrementModule;