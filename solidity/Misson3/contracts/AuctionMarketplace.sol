// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "./AuctionNFT.sol";

/**
 * @title AuctionMarketplace
 * @dev NFT 拍卖市场合约，支持 ERC20 和以太坊出价，集成 Chainlink 预言机
 */
contract AuctionMarketplace is 
    Initializable, 
    UUPSUpgradeable, 
    OwnableUpgradeable,
    ReentrancyGuardUpgradeable
{
    using SafeERC20 for IERC20;
    
    // 拍卖结构体
    struct Auction {
        address seller;              // 卖家地址
        address nftContract;         // NFT 合约地址
        uint256 tokenId;             // NFT Token ID
        uint256 startTime;           // 拍卖开始时间
        uint256 endTime;             // 拍卖结束时间
        uint256 startingPrice;       // 起拍价（以 wei 为单位）
        address paymentToken;        // 支付代币地址（address(0) 表示以太坊）
        address highestBidder;       // 当前最高出价者
        uint256 highestBid;          // 当前最高出价
        bool ended;                  // 拍卖是否已结束
    }
    
    // 拍卖 ID 到拍卖信息的映射
    mapping(uint256 => Auction) public auctions;
    
    // 用户出价记录（拍卖 ID => 出价者 => 出价金额）
    mapping(uint256 => mapping(address => uint256)) public bids;
    
    // 用户可提取的资金（出价者 => 可提取金额）
    mapping(address => uint256) public pendingReturns;
    
    // Chainlink 价格预言机地址映射（代币地址 => 价格预言机地址）
    mapping(address => address) public priceFeeds;
    
    // 平台手续费率,基点为单位，100 = 1%
    uint256 public platformFeeRate;
    
    // 平台手续费接收地址
    address public platformFeeRecipient;
    
    // 拍卖计数器
    uint256 private _auctionIdCounter;
    
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }
    
    /**
     * @dev 初始化函数（UUPS 代理模式）
     * @param _platformFeeRate 平台手续费率（基点）
     * @param _platformFeeRecipient 平台手续费接收地址
     */
    function initialize(
        uint256 _platformFeeRate,
        address _platformFeeRecipient
    ) public initializer {
        __Ownable_init(msg.sender);
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
        
        require(_platformFeeRate < 10000, "Fee rate too high"); //验证手续费率合理性
        require(_platformFeeRecipient != address(0), "Invalid fee recipient");
        
        platformFeeRate = _platformFeeRate;
        platformFeeRecipient = _platformFeeRecipient;
        _auctionIdCounter = 1;
    }
    
    /**
     * @dev UUPS 升级授权函数
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
    
    /**
     * @dev 设置价格预言机地址
     * @param token 代币地址（address(0) 表示以太坊）
     * @param priceFeed Chainlink 价格预言机地址
     */
    function setPriceFeed(address token, address priceFeed) external onlyOwner {
        require(priceFeed != address(0), "Invalid price feed");
        priceFeeds[token] = priceFeed;
    }
    
    /**
     * @dev 设置平台手续费率
     * @param _platformFeeRate 新的手续费率（基点）
     */
    function setPlatformFeeRate(uint256 _platformFeeRate) external onlyOwner {
        require(_platformFeeRate < 10000, "Fee rate too high");
        platformFeeRate = _platformFeeRate;
    }
    
    /**
     * @dev 设置平台手续费接收地址
     * @param _platformFeeRecipient 新的手续费接收地址
     */
    function setPlatformFeeRecipient(address _platformFeeRecipient) external onlyOwner {
        require(_platformFeeRecipient != address(0), "Invalid fee recipient");
        platformFeeRecipient = _platformFeeRecipient;
    }
    
    /**
     * @dev 创建拍卖
     * @param nftContract NFT 合约地址
     * @param tokenId NFT Token ID
     * @param duration 拍卖持续时间（秒）
     * @param startingPrice 起拍价
     * @param paymentToken 支付代币地址（address(0) 表示以太坊）
     * @return auctionId 创建的拍卖 ID
     */
    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 duration,
        uint256 startingPrice,
        address paymentToken
    ) external returns (uint256) {
        require(duration > 0, "Duration must be greater than 0");
        require(startingPrice > 0, "Starting price must be greater than 0");
        
        // 检查 NFT 所有权
        IERC721 nft = IERC721(nftContract);
        require(nft.ownerOf(tokenId) == msg.sender, "Not the owner");
        
        // 检查是否已授权给本合约
        require(
            nft.getApproved(tokenId) == address(this) || 
            nft.isApprovedForAll(msg.sender, address(this)),
            "NFT not approved"
        );
        
        // 转移 NFT 到合约
        nft.transferFrom(msg.sender, address(this), tokenId);
        
        uint256 auctionId = _auctionIdCounter++;
        uint256 startTime = block.timestamp;
        uint256 endTime = startTime + duration;
        
        auctions[auctionId] = Auction({
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            startTime: startTime,
            endTime: endTime,
            startingPrice: startingPrice,
            paymentToken: paymentToken,
            highestBidder: address(0),
            highestBid: 0,
            ended: false
        });
        
        return auctionId;
    }
    
    /**
     * @dev 出价（以太坊）
     * @param auctionId 拍卖 ID
     */
    function bid(uint256 auctionId) external payable nonReentrant {
        Auction storage auction = auctions[auctionId];
        
        require(!auction.ended, "Auction already ended");
        require(block.timestamp >= auction.startTime, "Auction not started");
        require(block.timestamp < auction.endTime, "Auction ended");
        require(auction.paymentToken == address(0), "Use bidWithToken for ERC20");
        require(msg.value >= auction.startingPrice, "Bid too low");
        require(msg.value > auction.highestBid, "Bid must be higher");
        
        // 退还前一个最高出价者的资金
        if (auction.highestBidder != address(0)) {
            pendingReturns[auction.highestBidder] += auction.highestBid;
        }
        
        // 更新最高出价
        auction.highestBidder = msg.sender;
        auction.highestBid = msg.value;
        bids[auctionId][msg.sender] = msg.value;
    }
    
    /**
     * @dev 出价（ERC20）
     * @param auctionId 拍卖 ID
     * @param bidAmount 出价金额
     */
    function bidWithToken(uint256 auctionId, uint256 bidAmount) external nonReentrant {
        Auction storage auction = auctions[auctionId];
        
        require(!auction.ended, "Auction already ended");
        require(block.timestamp >= auction.startTime, "Auction not started");
        require(block.timestamp < auction.endTime, "Auction ended");
        require(auction.paymentToken != address(0), "Use bid for ETH");
        require(bidAmount >= auction.startingPrice, "Bid too low");
        require(bidAmount > auction.highestBid, "Bid must be higher");
        
        IERC20 paymentToken = IERC20(auction.paymentToken);
        
        // 检查余额和授权
        require(
            paymentToken.balanceOf(msg.sender) >= bidAmount,
            "not enough token balance"
        );
        require(
            paymentToken.allowance(msg.sender, address(this)) >= bidAmount,
            "not enough token allowance"
        );
        
        // 退还前一个最高出价者的资金
        if (auction.highestBidder != address(0)) {
            pendingReturns[auction.highestBidder] += auction.highestBid;
        }
        
        // 转移代币到合约
        paymentToken.safeTransferFrom(msg.sender, address(this), bidAmount);
        
        // 更新最高出价
        auction.highestBidder = msg.sender;
        auction.highestBid = bidAmount;
        bids[auctionId][msg.sender] = bidAmount;
    }
    
    /**
     * @dev 结束拍卖
     * @param auctionId 拍卖 ID
     */
    function endAuction(uint256 auctionId) external virtual nonReentrant {
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
            
            // 计算手续费
            uint256 feeAmount = (auction.highestBid * platformFeeRate) / 10000;
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
     * @dev 提取退款
     */
    function withdraw() external nonReentrant {
        uint256 amount = pendingReturns[msg.sender];
        require(amount > 0, "No funds to withdraw");
        
        pendingReturns[msg.sender] = 0;
        payable(msg.sender).transfer(amount);
    }
    
    /**
     * @dev 获取出价金额的美元价值（使用 Chainlink 预言机）
     * @param token 代币地址（address(0) 表示以太坊）
     * @param amount 代币数量
     * @return 美元价值（以 8 位小数表示）
     */
    function getBidAmountUSD(address token, uint256 amount) public view returns (uint256) {
        address priceFeedAddress = priceFeeds[token];
        require(priceFeedAddress != address(0), "Price feed not set");
        
        AggregatorV3Interface priceFeed = AggregatorV3Interface(priceFeedAddress);
        (, int256 price, , , ) = priceFeed.latestRoundData();
        
        require(price > 0, "Invalid price");
        
        uint256 priceValue = uint256(price);
        
        // 计算美元价值（保持 8 位小数精度）
        // amount 是代币的最小单位（通常是 18 位小数，如 ETH）
        // priceValue 是价格（8 位小数，来自 Chainlink 预言机）
        // 结果 = (amount * priceValue) / 10^18，得到 8 位小数的美元价值
        return (amount * priceValue) / (10 ** 18);
    }
    
    /**
     * @dev 获取拍卖信息
     * @param auctionId 拍卖 ID
     * @return 拍卖信息结构体
     */
    function getAuction(uint256 auctionId) external view returns (Auction memory) {
        return auctions[auctionId];
    }
    
    /**
     * @dev 获取当前最高出价的美元价值
     * @param auctionId 拍卖 ID
     * @return 美元价值
     */
    function getHighestBidUSD(uint256 auctionId) external view returns (uint256) {
        Auction memory auction = auctions[auctionId];
        if (auction.highestBid == 0) {
            return 0;
        }
        return getBidAmountUSD(auction.paymentToken, auction.highestBid);
    }
}

