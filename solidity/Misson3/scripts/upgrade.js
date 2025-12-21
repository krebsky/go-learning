const { ethers, upgrades } = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    
    console.log("升级账户:", deployer.address);
    
    // 从环境变量或配置中获取代理地址
    const PROXY_ADDRESS = process.env.PROXY_ADDRESS || "0x..."; // 替换为实际的代理地址
    
    if (PROXY_ADDRESS === "0x...") {
        console.error("请设置 PROXY_ADDRESS 环境变量");
        process.exit(1);
    }
    
    console.log("代理合约地址:", PROXY_ADDRESS);
    
    // 部署新版本的实现合约
    console.log("\n部署 V2 实现合约...");
    const AuctionMarketplaceV2 = await ethers.getContractFactory("AuctionMarketplaceV2");
    
    const upgraded = await upgrades.upgradeProxy(PROXY_ADDRESS, AuctionMarketplaceV2);
    await upgraded.waitForDeployment();
    
    console.log("升级完成！");
    console.log("代理地址:", await upgraded.getAddress());
    
    // 获取新的实现地址
    const implementationAddress = await upgrades.erc1967.getImplementationAddress(await upgraded.getAddress());
    console.log("新的实现地址:", implementationAddress);
    
    // 初始化 V2 版本
    console.log("\n初始化 V2 版本...");
    const tx = await upgraded.initializeV2();
    await tx.wait();
    console.log("V2 版本初始化完成");
    
    // 验证版本
    const version = await upgraded.version();
    console.log("合约版本:", version);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

