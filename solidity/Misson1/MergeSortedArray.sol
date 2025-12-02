// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MergeSortedArray{

    function mergeSortedArray(uint[] memory arr1,uint[] memory arr2) public pure returns(uint[] memory){
        uint len1 = arr1.length;
        uint len2 = arr2.length;

        if(len1 == 0){
            return arr2;
        }else if(len2 == 0){
            return arr1;
        }

        uint i = 0;
        uint j = 0;
        uint k = 0;
        uint[] memory res = new uint[](len1+len2);
        while(i < len1 && j < len2){
            if(arr1[i]<=arr2[j]){
                res[k] = arr1[i];
                i++;
            }else{
                res[k] = arr2[j];
                j++;
            }
            k++;
        }
        if(i < len1){
            for(uint c=i;c < len1;c++){
                res[k] = arr1[c];
                k++;
            }
        }else{
            for(uint c=j;c < len2;c++){
                res[k] = arr2[c];
                k++;
            }

        }



        return res;


    }
}