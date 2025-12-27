// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
    uint private counter;

    function increase() public returns (uint) {
        counter++;
        return counter;
    }

    function get() public view returns (uint) {
        return counter;
    }
}

