const hre = require("hardhat");

async function main() {
    const [owner, voter1, voter2, voter3] = await hre.ethers.getSigners();
    const votingProposalContractFactory = await hre.ethers.getContractFactory("VotingProposal");
    // retrieve the contract address from running the ignition script for the contract
    const votingProposalContractAddress = "";
    const votingProposalContract = votingProposalContractFactory.attach(votingProposalContractAddress);
    // create a proposal and vote on it
    await votingProposalContract.connect(owner).create("4 day workweek proposal ", 3);
    // check if voter1 has voted before calling vote function
    let hasVoted = await votingProposalContract.hasVoted(voter1.address);
    console.log("voter1 has voted? ", hasVoted? "yes": "no");
    await votingProposalContract.connect(voter1).vote(1);
    // check if voter1 has voted after calling vote function
    hasVoted = await votingProposalContract.hasVoted(voter1.address);
    console.log("voter1 has voted? ", hasVoted? "yes": "no");
    await votingProposalContract.connect(voter2).vote(1);
    await votingProposalContract.connect(voter3).vote(0);
    // check proposal object    
    const proposal = await votingProposalContract.getCurrentProposal();
    console.log("current state of proposal: ", proposal.current_state ? "passed" : "rejected");
    console.log("is the proposal closed? ", proposal.is_active ? "no" : "yes");
}

main()
 .then(() => process.exit(0))
 .catch((error) => {
   console.error(error);
   process.exit(1);
 });