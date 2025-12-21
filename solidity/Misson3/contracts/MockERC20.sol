// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/**
 * @title MockERC20
 * @dev 用于测试的 ERC20 代币合约
 */
contract MockERC20 is ERC20 {
    constructor(string memory name, string memory symbol) ERC20(name, symbol) {
        // 给部署者铸造 1000000 个代币
        _mint(msg.sender, 1000000 * 10**decimals());
    }
    
    /**
     * @dev 铸造代币（用于测试）
     * @param to 接收地址
     * @param amount 数量
     */
    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

