// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Voting{

    mapping(string => uint) private votingCounter;

   // 存储所有候选人的数组，用于重置功能
    string[] private candidates;

    function vote(string memory userId) public {
        if(votingCounter[userId] == 0){
            candidates.push(userId);
        }
        votingCounter[userId]++;
    }

    function getVotes(string memory userId) public view returns (uint) {
        return votingCounter[userId];
    } 

    function resetVotes() public {
        for(uint i=0;i < candidates.length;i++){
            votingCounter[candidates[i]] = 0;
        }
        delete candidates;
    }


}