require("@nomicfoundation/hardhat-toolbox");
require("@openzeppelin/hardhat-upgrades");
require("hardhat-deploy");

// 可选加载 .env 文件
try {
    require("dotenv").config();
} catch (e) {
    // dotenv 未安装或 .env 文件不存在，继续执行
}

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: {
    version: "0.8.22",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  networks: {
    hardhat: {
      chainId: 1337,
      // Fork Sepolia 测试网以使用真实的 Chainlink 价格预言机
      // 注意：启用 fork 需要网络连接和稳定的 RPC 端点
      // 如果网络不稳定，可以设置 DISABLE_FORK=true 来禁用 fork
      // 禁用 fork 后，需要在测试网环境下运行测试才能使用真实的 Chainlink 价格预言机
      forking: {
        // 强烈建议使用 Alchemy 或 Infura 等专业 RPC 服务
        // 公共端点 rpc.sepolia.org 可能不稳定或响应慢
        url: process.env.SEPOLIA_RPC_URL || "https://rpc.sepolia.org",
        enabled: process.env.DISABLE_FORK !== "true", // 默认启用 fork，可通过 DISABLE_FORK=true 禁用
        // 固定区块号以提高性能（强烈推荐）
        // 如果不设置，每次都会 fork 最新区块，速度很慢
        blockNumber: process.env.FORK_BLOCK_NUMBER 
          ? parseInt(process.env.FORK_BLOCK_NUMBER) 
          : 6000000, // 默认使用一个较新的区块号（Sepolia 测试网）
        httpHeaders: {}, // 可以添加自定义 HTTP 头
      },
      // 增加超时时间
      timeout: 180000, // 180 秒（3 分钟）
      // 增加 gas 限制
      gas: 12000000,
      blockGasLimit: 12000000,
    },
    sepolia: {
      url: process.env.SEPOLIA_RPC_URL || "https://rpc.sepolia.org",
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
      chainId: 11155111,
    },
  },
  etherscan: {
    apiKey: process.env.ETHERSCAN_API_KEY || "",
  },
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts",
  },
};

