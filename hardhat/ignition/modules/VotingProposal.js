const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

const votingProposalModule = buildModule("VotingProposal", (m) => {
    const votingProposal = m.contract("VotingProposal", []);
    return { votingProposal };  
});

module.exports = votingProposalModule;