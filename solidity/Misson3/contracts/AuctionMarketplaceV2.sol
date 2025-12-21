// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import "./AuctionMarketplace.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/**
 * @title AuctionMarketplaceV2
 * @dev 升级版本的拍卖市场合约，添加动态手续费功能，假设只有两个阈值
 */
contract AuctionMarketplaceV2 is AuctionMarketplace {
    using SafeERC20 for IERC20;
    // 动态手续费配置
    struct DynamicFeeConfig {
        uint256 threshold1;      // 第一个阈值
        uint256 feeRate1;        // 第一个阈值的手续费率
        uint256 threshold2;      // 第二个阈值
        uint256 feeRate2;        // 第二个阈值的手续费率
        uint256 feeRate3;        // 超过第二个阈值的手续费率
    }
    
    DynamicFeeConfig public dynamicFeeConfig;
    bool public useDynamicFee;
    
    /**
     * @dev 初始化 V2 版本（仅用于新部署，升级时不需要）
     */
    function initializeV2() public reinitializer(2) {
        // V2 的初始化逻辑
    }
    
    /**
     * @dev 设置动态手续费配置
     * @param threshold1 第一个阈值
     * @param feeRate1 第一个阈值的手续费率
     * @param threshold2 第二个阈值
     * @param feeRate2 第二个阈值的手续费率
     * @param feeRate3 超过第二个阈值的手续费率
     */
    function setDynamicFeeConfig(
        uint256 threshold1,
        uint256 feeRate1,
        uint256 threshold2,
        uint256 feeRate2,
        uint256 feeRate3
    ) external onlyOwner {
        require(feeRate1 <= 10000 && feeRate2 <= 10000 && feeRate3 <= 10000, "Fee rate too high");
        require(threshold1 < threshold2, "Invalid thresholds");
        
        dynamicFeeConfig = DynamicFeeConfig({
            threshold1: threshold1,
            feeRate1: feeRate1,
            threshold2: threshold2,
            feeRate2: feeRate2,
            feeRate3: feeRate3
        });
    }
    
    /**
     * @dev 启用/禁用动态手续费
     * @param enabled 是否启用动态手续费
     */
    function toggleDynamicFee(bool enabled) external onlyOwner {
        useDynamicFee = enabled;
    }
    
    /**
     * @dev 根据金额计算动态手续费率
     * @param amount 拍卖金额
     * @return 手续费率（基点）
     */
    function calculateDynamicFeeRate(uint256 amount) public view returns (uint256) {
        if (!useDynamicFee) {
            return platformFeeRate;
        }
        
        if (amount < dynamicFeeConfig.threshold1) {
            return dynamicFeeConfig.feeRate1;
        } else if (amount < dynamicFeeConfig.threshold2) {
            return dynamicFeeConfig.feeRate2;
        } else {
            return dynamicFeeConfig.feeRate3;
        }
    }
    
    /**
     * @dev 重写结束拍卖函数，使用动态手续费
     * @param auctionId 拍卖 ID
     */
    function endAuction(uint256 auctionId) external override nonReentrant {
        Auction storage auction = auctions[auctionId];
        
        require(!auction.ended, "Auction already ended");
        require(
            block.timestamp >= auction.endTime || msg.sender == auction.seller,
            "Auction not ended"
        );
        
        auction.ended = true;
        
        IERC721 nft = IERC721(auction.nftContract);
        
        if (auction.highestBidder != address(0)) {
            // 有出价者，转移 NFT 给最高出价者
            nft.transferFrom(address(this), auction.highestBidder, auction.tokenId);
            
            // 计算手续费（使用动态手续费）
            uint256 feeRate = calculateDynamicFeeRate(auction.highestBid);
            uint256 feeAmount = (auction.highestBid * feeRate) / 10000;
            uint256 sellerAmount = auction.highestBid - feeAmount;
            
            // 转移资金
            if (auction.paymentToken == address(0)) {
                // 以太坊支付
                payable(auction.seller).transfer(sellerAmount);
                payable(platformFeeRecipient).transfer(feeAmount);
            } else {
                // ERC20 支付
                IERC20 paymentToken = IERC20(auction.paymentToken);
                paymentToken.safeTransfer(auction.seller, sellerAmount);
                paymentToken.safeTransfer(platformFeeRecipient, feeAmount);
            }
        } else {
            // 没有出价者，退还 NFT 给卖家
            nft.transferFrom(address(this), auction.seller, auction.tokenId);
        }
    }
    
    /**
     * @dev V2 新增功能：获取版本号
     */
    function version() external pure returns (string memory) {
        return "2.0.0";
    }
}

