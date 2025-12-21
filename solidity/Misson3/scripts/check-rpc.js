/**
 * 检查 RPC 端点连接和响应速度
 * 运行: npx hardhat run scripts/check-rpc.js
 */

// 加载 .env 文件
try {
    require("dotenv").config();
} catch (e) {
    console.warn("dotenv 未安装，将使用环境变量或默认值");
}

const { ethers } = require("hardhat");

async function main() {
    console.log("检查 RPC 端点连接...\n");
    
    const rpcUrl = process.env.SEPOLIA_RPC_URL || "https://rpc.sepolia.org";
    console.log("RPC URL:", rpcUrl);
    
    try {
        // 创建 provider
        const provider = new ethers.JsonRpcProvider(rpcUrl);
        
        // 测试连接
        console.log("\n1. 测试基本连接...");
        const startTime = Date.now();
        const blockNumber = await provider.getBlockNumber();
        const endTime = Date.now();
        const responseTime = endTime - startTime;
        
        console.log("✅ 连接成功！");
        console.log("当前区块号:", blockNumber);
        console.log("响应时间:", responseTime, "ms");
        
        if (responseTime > 5000) {
            console.log("⚠️  警告：响应时间较慢，建议使用更快的 RPC 端点（如 Alchemy 或 Infura）");
        }
        
        // 测试 Chainlink 价格预言机地址
        console.log("\n2. 测试 Chainlink 价格预言机...");
        const CHAINLINK_ETH_USD_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306";
        
        const priceFeedABI = [
            "function latestRoundData() external view returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)",
            "function decimals() external view returns (uint8)"
        ];
        
        const priceFeed = new ethers.Contract(CHAINLINK_ETH_USD_FEED, priceFeedABI, provider);
        
        const priceStartTime = Date.now();
        const [roundId, price, startedAt, updatedAt, answeredInRound] = await priceFeed.latestRoundData();
        const priceEndTime = Date.now();
        const priceResponseTime = priceEndTime - priceStartTime;
        
        console.log("✅ Chainlink 价格预言机可访问！");
        console.log("ETH/USD 价格:", ethers.formatUnits(price, 8), "USD");
        console.log("响应时间:", priceResponseTime, "ms");
        
        if (priceResponseTime > 10000) {
            console.log("⚠️  警告：价格查询响应时间较慢");
        }
        
        // 建议
        console.log("\n3. 建议：");
        if (rpcUrl.includes("rpc.sepolia.org")) {
            console.log("⚠️  您正在使用公共 RPC 端点，可能较慢且不稳定");
            console.log("💡 建议使用专业的 RPC 服务：");
            console.log("   - Alchemy: https://www.alchemy.com/");
            console.log("   - Infura: https://www.infura.io/");
            console.log("   - QuickNode: https://www.quicknode.com/");
        }
        
        console.log("\n✅ RPC 端点检查完成！");
        
    } catch (error) {
        console.error("❌ 连接失败:", error.message);
        console.log("\n故障排除建议：");
        console.log("1. 检查网络连接");
        console.log("2. 验证 RPC URL 是否正确");
        console.log("3. 如果使用 Alchemy/Infura，检查 API 密钥是否有效");
        console.log("4. 尝试使用其他 RPC 端点");
        process.exit(1);
    }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

