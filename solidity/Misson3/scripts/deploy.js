const { ethers, upgrades } = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    
    console.log("部署账户:", deployer.address);
    console.log("账户余额:", ethers.formatEther(await ethers.provider.getBalance(deployer.address)), "ETH");

    // 部署 NFT 合约
    console.log("\n部署 NFT 合约...");
    const AuctionNFT = await ethers.getContractFactory("AuctionNFT");
    const nftContract = await AuctionNFT.deploy("Auction NFT", "ANFT");
    await nftContract.waitForDeployment();
    const nftAddress = await nftContract.getAddress();
    console.log("NFT 合约地址:", nftAddress);

    // 部署拍卖市场合约（UUPS 代理）
    console.log("\n部署拍卖市场合约...");
    const AuctionMarketplace = await ethers.getContractFactory("AuctionMarketplace");
    
    const PLATFORM_FEE_RATE = 250; // 2.5%
    const PLATFORM_FEE_RECIPIENT = deployer.address; // 
    
    const marketplace = await upgrades.deployProxy(
        AuctionMarketplace,
        [PLATFORM_FEE_RATE, PLATFORM_FEE_RECIPIENT],
        { initializer: "initialize" }
    );
    await marketplace.waitForDeployment();
    const marketplaceAddress = await marketplace.getAddress();
    console.log("拍卖市场合约地址:", marketplaceAddress);
    
    // 获取代理地址和实现地址
    const implementationAddress = await upgrades.erc1967.getImplementationAddress(marketplaceAddress);
    console.log("实现合约地址:", implementationAddress);

    // 设置 Chainlink 价格预言机
    // 使用真实的 Chainlink 价格预言机地址
    console.log("\n设置 Chainlink 价格预言机...");
    
    // Chainlink 价格预言机地址（Sepolia 测试网）
    const CHAINLINK_ETH_USD_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306";
    const CHAINLINK_LINK_USD_FEED = "0xc59E3633BAAC79493d908e63626716e8A18Ad701";
    const CHAINLINK_USDC_USD_FEED = "0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E"; // USDC/USD on Sepolia
    
    // 设置 ETH/USD 价格预言机（使用真实的 Chainlink 价格预言机）
    await marketplace.setPriceFeed(ethers.ZeroAddress, CHAINLINK_ETH_USD_FEED);
    console.log("已设置 ETH/USD 价格预言机:", CHAINLINK_ETH_USD_FEED);

    console.log("\n部署完成！");
    console.log("\n合约地址汇总:");
    console.log("NFT 合约:", nftAddress);
    console.log("拍卖市场合约（代理）:", marketplaceAddress);
    console.log("拍卖市场合约（实现）:", implementationAddress);
    
    // 保存部署信息
    const deploymentInfo = {
        network: network.name,
        nftContract: nftAddress,
        marketplaceProxy: marketplaceAddress,
        marketplaceImplementation: implementationAddress,
        platformFeeRate: PLATFORM_FEE_RATE,
        platformFeeRecipient: PLATFORM_FEE_RECIPIENT,
        deployer: deployer.address,
        timestamp: new Date().toISOString()
    };
    
    console.log("\n部署信息:", JSON.stringify(deploymentInfo, null, 2));
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

