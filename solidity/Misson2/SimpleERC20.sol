// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleERC20{
    string public name="Great";
    mapping(address=>uint) balance;
    mapping(address=>mapping(address=>uint)) aprovels;

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);


    function  balanceOf(address account) external view returns (uint256){
        return balance[account];
    }
    function transfer(address to, uint256 amount) external  returns (bool){
        address from =msg.sender;
        require(to != address(0), "to zero address");

        //uint fromBalance = balance[from];
        //require(fromBalance >= amount, "transfer ammount exceeds banlance");
        //balance[from] -= amount;
        balance[to] += amount;

        emit Transfer(from, to, amount);

        return true;
    }

    function approve(address spender, uint256 amount) external returns (bool){
        address tokenOwner = msg.sender;
        require(spender != address(0), "to zero address");

        aprovels[tokenOwner][spender] = amount;
        emit Approval(tokenOwner, spender, amount);
        return true;
    }

    function transferFrom(address from, address to, uint256 amount) external returns (bool){
        require(from != address(0), "from zero address");
        require(to != address(0), "to zero address");

        uint currBalance = balance[from];
        uint currApprovalLeft = aprovels[msg.sender][from];
        require(currBalance >= amount, "transfer exceeds balance");
        require(currApprovalLeft >= amount, "approval exceeds limit");
        aprovels[msg.sender][from] = currBalance-amount;
        emit Approval(msg.sender, from, amount);

        balance[from] -= amount;
        balance[to] += amount;
        emit Transfer(from, to, amount);

        return true;
    }



}