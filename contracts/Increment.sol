// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Increment {

    address owner;
    struct Counter {
        uint256 number;
        string description;
    }
    Counter counter;

    modifier  onlyOwner {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    constructor(uint256 initial_value, string memory description){
        owner = msg.sender;
        counter = Counter(initial_value, description);
    }

        // Since these funcions modify data on chain, they require gas
    function increment() external onlyOwner {
        counter.number += 1;
    }

    function decrement() external onlyOwner {
        counter.number -= 1;
    }
    
    function getCounter() public view returns (uint256) {
        return counter.number;
    }
}