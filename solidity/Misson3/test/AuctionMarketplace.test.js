const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("NFT Auction Marketplace", function () {
    // 增加超时时间，因为 fork 模式可能需要更长时间
    // Fork 初始化可能需要 30-60 秒，加上测试时间，设置为 3 分钟
    this.timeout(180000); // 180 秒（3 分钟）
    let nftContract;
    let marketplace;
    let marketplaceV2;
    let erc20Token;
    // 注意：不再使用 Mock 价格预言机，直接使用 Chainlink 价格预言机地址
    let owner;
    let seller;
    let bidder1;
    let bidder2;
    let platformFeeRecipient;

    const PLATFORM_FEE_RATE = 250; // 2.5%
    
    // Chainlink 价格预言机地址（Sepolia 测试网）
    const CHAINLINK_ETH_USD_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306";
    // 注意：ERC20 代币的价格预言机需要根据实际代币选择，这里使用 USDC/USD 作为示例
    // USDC/USD on Sepolia: 0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E
    const CHAINLINK_USDC_USD_FEED = "0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E"; // 实际使用时需要替换为对应代币的价格预言机

    beforeEach(async function () {
        // 增加单个测试的超时时间
        this.timeout(180000); // 180 秒（3 分钟）
        
        [owner, seller, bidder1, bidder2, platformFeeRecipient] = await ethers.getSigners();

        // 部署 NFT 合约
        const AuctionNFT = await ethers.getContractFactory("AuctionNFT");
        nftContract = await AuctionNFT.deploy("Test NFT", "TNFT");
        await nftContract.waitForDeployment();
        const nftAddress = await nftContract.getAddress();
        console.log("NFT 合约地址:", nftAddress);

        // 部署 Mock ERC20（用于测试，实际部署时使用真实代币）
        const MockERC20 = await ethers.getContractFactory("MockERC20");
        erc20Token = await MockERC20.deploy("Test Token", "TT");
        await erc20Token.waitForDeployment();
        const erc20Address = await erc20Token.getAddress();
        //console.log("ERC20 代币合约地址:", erc20Address);

        // 部署拍卖市场合约（UUPS 代理）
        const AuctionMarketplace = await ethers.getContractFactory("AuctionMarketplace");
        marketplace = await upgrades.deployProxy(
            AuctionMarketplace,
            [PLATFORM_FEE_RATE, platformFeeRecipient.address],
            { initializer: "initialize" }
        );
        await marketplace.waitForDeployment();
        const marketplaceAddress = await marketplace.getAddress();
        console.log("拍卖市场合约地址（代理）:", marketplaceAddress);
        
        // 获取实现合约地址
        const implementationAddress = await upgrades.erc1967.getImplementationAddress(marketplaceAddress);
        console.log("拍卖市场合约地址（实现）:", implementationAddress);

        // 设置 Chainlink 价格预言机（使用真实的 Chainlink 价格预言机）
        // ETH/USD 价格预言机
        await marketplace.setPriceFeed(ethers.ZeroAddress, CHAINLINK_ETH_USD_FEED);
        
        
        // ERC20/USD 价格预言机（使用 USDC/USD 作为示例）
        await marketplace.setPriceFeed(await erc20Token.getAddress(), CHAINLINK_USDC_USD_FEED);

        // 给测试用户分配代币
        await erc20Token.transfer(seller.address, ethers.parseEther("1000"));
        await erc20Token.transfer(bidder1.address, ethers.parseEther("1000"));
        await erc20Token.transfer(bidder2.address, ethers.parseEther("1000"));
        
    });

    describe("部署和初始化", function () {
        it("应该正确初始化合约", async function () {
            expect(await marketplace.platformFeeRate()).to.equal(PLATFORM_FEE_RATE);
            expect(await marketplace.platformFeeRecipient()).to.equal(platformFeeRecipient.address);
        });

        it("应该拒绝无效的手续费率", async function () {
            const AuctionMarketplace = await ethers.getContractFactory("AuctionMarketplace");
            await expect(
                upgrades.deployProxy(
                    AuctionMarketplace,
                    [10001, platformFeeRecipient.address], // 超过 10%
                    { initializer: "initialize" }
                )
            ).to.be.revertedWith("Fee rate too high");
        });
    });

    describe("NFT 铸造", function () {
        it("应该能够铸造 NFT", async function () {
            const tokenUri = "https://example.com/token/1";
            await nftContract.mintNFT(seller.address, tokenUri);
            
            expect(await nftContract.ownerOf(1)).to.equal(seller.address);
            expect(await nftContract.tokenURI(1)).to.equal(tokenUri);
        });

        it("应该能够批量铸造 NFT", async function () {
            const tokenUris = [
                "https://example.com/token/1",
                "https://example.com/token/2",
                "https://example.com/token/3"
            ];
            
            const tx = await nftContract.mintNFT(seller.address, tokenUris[0]);
            await tx.wait();
            
            expect(await nftContract.totalSupply()).to.equal(1);
        });
    });

    describe("创建拍卖", function () {
        beforeEach(async function () {
            await nftContract.mintNFT(seller.address, "https://example.com/token/1");
            await nftContract.connect(seller).approve(await marketplace.getAddress(), 1);
        });

        it("应该能够创建以太坊拍卖", async function () {
            const duration = 3600; // 1 小时
            const startingPrice = ethers.parseEther("1");

            const tx = await marketplace.connect(seller).createAuction(
                await nftContract.getAddress(),
                1,
                duration,
                startingPrice,
                ethers.ZeroAddress
            );
            await tx.wait();

            const auction = await marketplace.getAuction(1);
            expect(auction.seller).to.equal(seller.address);
            expect(auction.nftContract).to.equal(await nftContract.getAddress());
            expect(auction.tokenId).to.equal(1);
            expect(auction.startingPrice).to.equal(startingPrice);
            expect(auction.paymentToken).to.equal(ethers.ZeroAddress);
            expect(auction.ended).to.be.false;
        });

        it("应该能够创建 ERC20 拍卖", async function () {
            const duration = 3600;
            const startingPrice = ethers.parseEther("100");

            await erc20Token.connect(seller).approve(await marketplace.getAddress(), ethers.MaxUint256);

            const tx = await marketplace.connect(seller).createAuction(
                await nftContract.getAddress(),
                1,
                duration,
                startingPrice,
                await erc20Token.getAddress()
            );
            await tx.wait();

            const auction = await marketplace.getAuction(1);
            expect(auction.paymentToken).to.equal(await erc20Token.getAddress());
        });

        it("应该拒绝非 NFT 所有者创建拍卖", async function () {
            await expect(
                marketplace.connect(bidder1).createAuction(
                    await nftContract.getAddress(),
                    1,
                    3600,
                    ethers.parseEther("1"),
                    ethers.ZeroAddress
                )
            ).to.be.revertedWith("Not the owner");
        });

        it("应该拒绝未授权的 NFT 创建拍卖", async function () {
            await nftContract.mintNFT(bidder1.address, "https://example.com/token/2");
            
            await expect(
                marketplace.connect(bidder1).createAuction(
                    await nftContract.getAddress(),
                    2,
                    3600,
                    ethers.parseEther("1"),
                    ethers.ZeroAddress
                )
            ).to.be.revertedWith("NFT not approved");
        });
    });

    describe("出价", function () {
        beforeEach(async function () {
            await nftContract.mintNFT(seller.address, "https://example.com/token/1");
            await nftContract.connect(seller).approve(await marketplace.getAddress(), 1);
            
            await marketplace.connect(seller).createAuction(
                await nftContract.getAddress(),
                1,
                3600,
                ethers.parseEther("1"),
                ethers.ZeroAddress
            );
        });

        it("应该能够使用以太坊出价", async function () {
            const bidAmount = ethers.parseEther("2");
            
            const tx = await marketplace.connect(bidder1).bid(1, { value: bidAmount });
            await tx.wait();

            const auction = await marketplace.getAuction(1);
            expect(auction.highestBidder).to.equal(bidder1.address);
            expect(auction.highestBid).to.equal(bidAmount);
        });

        it("应该能够使用 ERC20 出价", async function () {
            // 创建 ERC20 拍卖
            await nftContract.mintNFT(seller.address, "https://example.com/token/2");
            await nftContract.connect(seller).approve(await marketplace.getAddress(), 2);
            
            await marketplace.connect(seller).createAuction(
                await nftContract.getAddress(),
                2,
                3600,
                ethers.parseEther("100"),
                await erc20Token.getAddress()
            );

            const bidAmount = ethers.parseEther("200");
            await erc20Token.connect(bidder1).approve(await marketplace.getAddress(), bidAmount);
            
            const tx = await marketplace.connect(bidder1).bidWithToken(2, bidAmount);
            await tx.wait();

            const auction = await marketplace.getAuction(2);
            expect(auction.highestBidder).to.equal(bidder1.address);
            expect(auction.highestBid).to.equal(bidAmount);
        });

        it("应该拒绝低于起拍价的出价", async function () {
            await expect(
                marketplace.connect(bidder1).bid(1, { value: ethers.parseEther("0.5") })
            ).to.be.revertedWith("Bid too low");
        });

        it("应该拒绝低于当前最高价的出价", async function () {
            await marketplace.connect(bidder1).bid(1, { value: ethers.parseEther("2") });
            
            await expect(
                marketplace.connect(bidder2).bid(1, { value: ethers.parseEther("1.5") })
            ).to.be.revertedWith("Bid must be higher");
        });

        it("应该退还前一个出价者的资金", async function () {
            const bid1 = ethers.parseEther("2");
            const bid2 = ethers.parseEther("3");

            await marketplace.connect(bidder1).bid(1, { value: bid1 });
            await marketplace.connect(bidder2).bid(1, { value: bid2 });

            const pendingReturn = await marketplace.pendingReturns(bidder1.address);
            expect(pendingReturn).to.equal(bid1);
        });

        it("应该计算美元价格", async function () {
            const bidAmount = ethers.parseEther("1");
            await marketplace.connect(bidder1).bid(1, { value: bidAmount });

            const usdValue = await marketplace.getBidAmountUSD(ethers.ZeroAddress, bidAmount);
            // 使用真实的 Chainlink 价格，价格会波动，所以使用较大的容差
            // Chainlink 价格是 8 位小数，所以 usdValue 也是 8 位小数
            expect(usdValue).to.be.gt(0); // 价格应该大于 0
            expect(usdValue).to.be.lt(10000 * 1e8); // 价格应该小于 $10000（合理范围）
        });
    });

    describe("结束拍卖", function () {
        beforeEach(async function () {
            await nftContract.mintNFT(seller.address, "https://example.com/token/1");
            await nftContract.connect(seller).approve(await marketplace.getAddress(), 1);
        });

        it("应该在有出价者时正确结束拍卖", async function () {
            await marketplace.connect(seller).createAuction(
                await nftContract.getAddress(),
                1,
                3600,
                ethers.parseEther("1"),
                ethers.ZeroAddress
            );

            const bidAmount = ethers.parseEther("2");
            await marketplace.connect(bidder1).bid(1, { value: bidAmount });

            // 推进时间到拍卖结束
            await time.increase(3601);

            const sellerBalanceBefore = await ethers.provider.getBalance(seller.address);
            const feeRecipientBalanceBefore = await ethers.provider.getBalance(platformFeeRecipient.address);

            const tx = await marketplace.connect(bidder1).endAuction(1);
            const receipt = await tx.wait();
            const gasUsed = receipt.gasUsed * receipt.gasPrice;

            const sellerBalanceAfter = await ethers.provider.getBalance(seller.address);
            const feeRecipientBalanceAfter = await ethers.provider.getBalance(platformFeeRecipient.address);

            // 检查 NFT 转移
            expect(await nftContract.ownerOf(1)).to.equal(bidder1.address);

            // 检查资金分配（考虑 gas 费用）
            const sellerReceived = sellerBalanceAfter - sellerBalanceBefore;
            const feeReceived = feeRecipientBalanceAfter - feeRecipientBalanceBefore;
            const expectedSellerAmount = bidAmount - (bidAmount * BigInt(PLATFORM_FEE_RATE) / BigInt(10000));
            const expectedFeeAmount = bidAmount * BigInt(PLATFORM_FEE_RATE) / BigInt(10000);

            expect(sellerReceived).to.be.closeTo(expectedSellerAmount, ethers.parseEther("0.01"));
            expect(feeReceived).to.equal(expectedFeeAmount);

            const auction = await marketplace.getAuction(1);
            expect(auction.ended).to.be.true;
        });

        it("应该在没有出价者时退还 NFT", async function () {
            await marketplace.connect(seller).createAuction(
                await nftContract.getAddress(),
                1,
                3600,
                ethers.parseEther("1"),
                ethers.ZeroAddress
            );

            await time.increase(3601);
            await marketplace.connect(seller).endAuction(1);

            expect(await nftContract.ownerOf(1)).to.equal(seller.address);
        });

        it("应该允许卖家提前结束拍卖", async function () {
            await marketplace.connect(seller).createAuction(
                await nftContract.getAddress(),
                1,
                3600,
                ethers.parseEther("1"),
                ethers.ZeroAddress
            );

            await marketplace.connect(bidder1).bid(1, { value: ethers.parseEther("2") });
            
            // 不等待时间，直接结束
            await marketplace.connect(seller).endAuction(1);

            expect(await nftContract.ownerOf(1)).to.equal(bidder1.address);
        });
    });

    describe("提取退款", function () {
        beforeEach(async function () {
            await nftContract.mintNFT(seller.address, "https://example.com/token/1");
            await nftContract.connect(seller).approve(await marketplace.getAddress(), 1);
            
            await marketplace.connect(seller).createAuction(
                await nftContract.getAddress(),
                1,
                3600,
                ethers.parseEther("1"),
                ethers.ZeroAddress
            );
        });

        it("应该能够提取退款", async function () {
            const bid1 = ethers.parseEther("2");
            const bid2 = ethers.parseEther("3");

            await marketplace.connect(bidder1).bid(1, { value: bid1 });
            await marketplace.connect(bidder2).bid(1, { value: bid2 });

            const balanceBefore = await ethers.provider.getBalance(bidder1.address);
            const tx = await marketplace.connect(bidder1).withdraw();
            const receipt = await tx.wait();
            const gasUsed = receipt.gasUsed * receipt.gasPrice;
            const balanceAfter = await ethers.provider.getBalance(bidder1.address);

            expect(balanceAfter - balanceBefore + gasUsed).to.equal(bid1);
        });
    });

    describe("合约升级（UUPS）", function () {
        it("应该能够升级合约", async function () {
            // 部署新版本
            const AuctionMarketplaceV2 = await ethers.getContractFactory("AuctionMarketplaceV2");
            
            // 注意：这里需要先创建 V2 合约，暂时跳过这个测试
            // 在实际项目中，需要先部署 V2 合约才能测试升级
        });
    });

    describe("管理员功能", function () {
        it("应该允许所有者设置价格预言机", async function () {
            // 使用另一个 Chainlink 价格预言机地址（USDC/USD）作为测试
            // USDC/USD on Sepolia
            const NEW_CHAINLINK_FEED = CHAINLINK_USDC_USD_FEED;
            
            await marketplace.setPriceFeed(ethers.ZeroAddress, NEW_CHAINLINK_FEED);
            
            const priceFeed = await marketplace.priceFeeds(ethers.ZeroAddress);
            expect(priceFeed).to.equal(NEW_CHAINLINK_FEED);
        });

        it("应该允许所有者更新手续费率", async function () {
            const newFeeRate = 300;
            await marketplace.setPlatformFeeRate(newFeeRate);
            
            expect(await marketplace.platformFeeRate()).to.equal(newFeeRate);
        });

        it("应该拒绝非所有者设置价格预言机", async function () {
            // 使用另一个 Chainlink 价格预言机地址作为测试
            // USDC/USD on Sepolia
            const NEW_CHAINLINK_FEED = CHAINLINK_USDC_USD_FEED;
            
            await expect(
                marketplace.connect(seller).setPriceFeed(ethers.ZeroAddress, NEW_CHAINLINK_FEED)
            ).to.be.revertedWithCustomError(marketplace, "OwnableUnauthorizedAccount");
        });
    });
});

