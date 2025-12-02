// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ReverseString {

    function reverseString(string memory sourceString) public pure returns(string memory){
        bytes memory bs = bytes(sourceString);
        uint256 length = bs.length;

        if(length <= 1){
            return sourceString;
        }

        for(uint i =0 ; i < length/2;i++){
            bytes1 temp = bs[i];
            bs[i] = bs[length-1-i];
            bs[length-1-i]=temp;
        }

        return string(bs);

    }
    
}