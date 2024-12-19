 // SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import {console} from "hardhat/console.sol";

contract demoOperations {
    uint256 variable;
    uint256 otherVariable;
    struct Account {
        address addr;
        uint256 amount;
    }

    function set(uint256 value) public {
        variable = value;
    }

    function get() public view returns (uint256) {
        return variable;
    }

    function callAnotherFunction() public pure {
        internalFunction();
        privateFunction();
    }

    function privateFunction() private pure  {
        console.log("privateFunction was called");
        // pure functions do not read from or write to the state.
        //  writing to the state incurs gas costs, as it involves making changes to the blockchain
        // This function can only be called from this contract+
        // Private functions are not visible to external contracts or accounts.
    }

    function internalFunction() internal pure {
        console.log("internalFunction was called");
        // This function can be called from this contract or contracts that derive from it.
    }

    function externalFunction() external pure {
        console.log("externalFunction was called");
        // This function can be called from other contracts and accounts.
    }

    function createAccount(address _addr, uint256 _amount) public pure {
        Account memory account = Account(_addr, _amount);
        console.log("Account created with address: ", account.addr, " and amount: ", account.amount);
    }

}



