// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import {console} from "hardhat/console.sol";

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
        console.log("increment constructor called with initial value %s and description %s", initial_value, description);
        owner = msg.sender;
        counter = Counter(initial_value, description);
        console.log("Counter initialized with value %s and description %s", counter.number, counter.description);
    }

        // Since these funcions modify data on chain, they require gas

    function increment() external onlyOwner {
        counter.number += 1;
    }

    function decrement() external onlyOwner {
        counter.number -= 1;
    }

    function get() public returns (uint256) {
        console.log("get called with value %s", counter.number);
        return counter.number;
    }
}