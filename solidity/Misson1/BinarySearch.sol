// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BinarySearch {
    function binarySearch(int256[] memory nums,int256 target) public pure returns(int256){

        uint256 left = 0;
        uint256 right = nums.length;
        
        if (right == 0) {
            return -1;
        }
        
      
        while (left < right) {
            uint256 mid = left + (right - left) / 2;
            
            if (nums[mid] == target) {
                return int256(mid);
            } else if (nums[mid] < target) {
                left = mid + 1;
            } else {
                right = mid;
            }
        }
        
        // 未找到
        return -1;

    }
}