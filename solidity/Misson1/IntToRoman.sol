// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract IntToRoman {

    function int2Roman(uint sourceInt) public pure returns(string memory){
        uint256[13] memory values = [
            uint256(1000), uint256(900), uint256(500), uint256(400),uint256(100), uint256(90), uint256(50), uint256(40),uint256(10), uint256(9), uint256(5), uint256(4), uint256(1)
        ];

        string[13] memory symbols = [
            "M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"
        ];


        string memory res;
        for(uint8 i=0; i < values.length;i++){
            uint256 value = values[i];
            string memory symbol = symbols[i];
            while(sourceInt >= value){
                sourceInt -= value;
                res= string.concat(res,symbol);
            }
            if(sourceInt == 0){
                break ;
            }

        }
        return res;
    }

}