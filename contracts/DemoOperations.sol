 // SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import {console} from "hardhat/console.sol";

contract DemoOperations {
    uint256 variable;
    uint256 otherVariable;
    
    struct Account {
        address addr;
        string name;
        uint256 amount;
    }

    struct Community {
        Account[] accounts;
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

    function createAccount(address _addr, string memory _name, uint256 _amount) public pure {
    // In Solidity, the data location must be explicitly specified for complex types like arrays and structs
    //  when they are used as function parameters. However, for simple types like address, 
    // string, and uint256, the compiler can infer the data location.
        Account memory account = Account(_addr, _name, _amount);
        console.log("Account created: %s, %s, %s", account.addr, account.name, account.amount);
    }

    function createCommunity() public pure {
        Account memory account1 = Account(0x5B38Da6a701c568545dCfcB03FcB875f56beddC4, "Alice", 100);
        Account memory account2 = Account(0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2, "Bob", 200);
        Account memory account3 = Account(0xCA35b7d915458EF540aDe6068dFe2F44E8fa733c, "Charlie", 300);
        Account[] memory accounts = new Account[](3);
        accounts[0] = account1;
        accounts[1] = account2;
        accounts[2] = account3;
        Community memory community = Community(accounts);
        console.log("Community created: %s, %s, %s", community.accounts[0].name, community.accounts[1].name, community.accounts[2].name);
    }

}



