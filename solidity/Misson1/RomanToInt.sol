// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomanToInt{
    
    function roman2Int(string memory romanString)  public pure returns (uint){
            bytes memory sourceBytes = bytes(romanString);
            uint len = sourceBytes.length;

            uint[] memory num = new uint[](len);
            for (uint i = 0; i < len; i++) {
                if (sourceBytes[i] == 'M') num[i] = 1000;
                else if (sourceBytes[i] == 'D') num[i] = 500;
                else if (sourceBytes[i] == 'C') num[i] = 100;
                else if (sourceBytes[i] == 'L') num[i] = 50;
                else if (sourceBytes[i] == 'X') num[i] = 10;
                else if (sourceBytes[i] == 'V') num[i] = 5;
                    else if (sourceBytes[i] == 'I') num[i] = 1;
            }

            uint  res;
            for(uint i=0;i<len-1;i++){
                if(num[i] < num[i+1]){
                    res = res -num[i];
                }else {
                    res = res + num[i];
                }
            }
            return res + num[len-1];
    }
}