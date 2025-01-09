// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";


contract DemoToken is ERC20 {

        uint256 constant initialSupply = 1000000 * (10**18);

    // DemoToken is the name of the contract and it inherits from ERC20
    constructor() ERC20("DemoToken", "DT") {
        // The mint function mins an initial supply of tokens
        _mint(msg.sender, initialSupply);
    }
}