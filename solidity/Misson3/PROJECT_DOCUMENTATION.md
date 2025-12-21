# NFT 拍卖市场项目文档

## 目录

1. [项目概述](#项目概述)
2. [功能说明](#功能说明)
3. [技术架构](#技术架构)
4. [环境要求](#环境要求)
5. [安装与配置](#安装与配置)
6. [部署步骤](#部署步骤)
7. [使用说明](#使用说明)
8. [测试说明](#测试说明)
9. [合约升级](#合约升级)
10. [安全考虑](#安全考虑)
11. [常见问题](#常见问题)
12. [附录](#附录)

---

## 项目概述

本项目是一个基于 Hardhat 框架开发的 NFT 拍卖市场系统，实现了完整的 NFT 拍卖功能，支持 ERC20 代币和以太坊出价，集成了 Chainlink 价格预言机进行美元价格计算，并使用 UUPS 代理模式实现合约升级。

### 核心特性

- ✅ **NFT 合约**：基于 ERC721 标准的 NFT 合约，支持铸造和转移
- ✅ **拍卖功能**：支持创建拍卖、出价（ERC20/以太坊）、结束拍卖
- ✅ **Chainlink 集成**：使用 Chainlink 价格预言机计算代币的美元价值
- ✅ **UUPS 升级**：使用 UUPS 代理模式实现合约升级
- ✅ **动态手续费**：V2 版本支持根据拍卖金额动态调整手续费（可选功能）
- ✅ **安全特性**：重入攻击保护、权限控制、输入验证

### 项目结构

```
Misson3/
├── contracts/                    # 智能合约
│   ├── AuctionNFT.sol           # NFT 合约（ERC721）
│   ├── AuctionMarketplace.sol   # 拍卖市场合约（V1）
│   ├── AuctionMarketplaceV2.sol # 拍卖市场合约（V2，支持动态手续费）
│   ├── MockERC20.sol            # 测试用 ERC20 代币
│   └── MockPriceFeed.sol        # 测试用价格预言机
├── scripts/                      # 部署脚本
│   ├── deploy.js                # 部署脚本
│   ├── upgrade.js               # 升级脚本
│   └── check-rpc.js             # RPC 连接检查脚本
├── test/                         # 测试文件
│   └── AuctionMarketplace.test.js
├── hardhat.config.js            # Hardhat 配置
├── package.json                 # 项目依赖
└── README.md                    # 项目文档
```

---

## 功能说明

### 1. NFT 合约（AuctionNFT）

NFT 合约基于 OpenZeppelin 的 ERC721URIStorage 标准实现，支持 NFT 的铸造、转移和元数据管理。

#### 主要功能

- **`mintNFT(address recipient, string tokenUri)`**
  - 功能：铸造新的 NFT
  - 参数：
    - `recipient`: NFT 接收者地址
    - `tokenUri`: NFT 元数据 URI（JSON 格式）
  - 返回：新铸造的 token ID

- **`batchMintNFT(address recipient, string[] tokenUris)`**
  - 功能：批量铸造 NFT
  - 参数：
    - `recipient`: NFT 接收者地址
    - `tokenUris`: NFT 元数据 URI 数组
  - 返回：新铸造的 token ID 数组

- **`totalSupply()`**
  - 功能：获取当前总供应量
  - 返回：已铸造的 NFT 总数

- **`nextTokenId()`**
  - 功能：获取下一个 token ID
  - 返回：下一个 token ID

- **`tokenURI(uint256 tokenId)`**
  - 功能：获取指定 token ID 的元数据 URI
  - 返回：元数据 URI 字符串

#### 使用示例

```javascript
// 铸造单个 NFT
const tokenId = await nftContract.mintNFT(
    seller.address,
    "https://example.com/metadata/1.json"
);

// 批量铸造 NFT
const tokenUris = [
    "https://example.com/metadata/1.json",
    "https://example.com/metadata/2.json",
    "https://example.com/metadata/3.json"
];
const tokenIds = await nftContract.batchMintNFT(seller.address, tokenUris);
```

### 2. 拍卖市场合约（AuctionMarketplace）

拍卖市场合约是系统的核心，实现了完整的拍卖功能，支持 ERC20 代币和以太坊出价。

#### 核心功能

##### 2.1 创建拍卖

**`createAuction(...)`**
- 功能：创建新的 NFT 拍卖
- 参数：
  - `nftContract`: NFT 合约地址
  - `tokenId`: NFT Token ID
  - `duration`: 拍卖持续时间（秒）
  - `startingPrice`: 起拍价（以 wei 为单位）
  - `paymentToken`: 支付代币地址（`address(0)` 表示以太坊）
- 要求：
  - NFT 所有者必须授权合约转移 NFT
  - 起拍价必须大于 0
  - 持续时间必须大于 0
- 返回：创建的拍卖 ID

**使用流程：**
1. 铸造或拥有 NFT
2. 授权拍卖市场合约转移 NFT
3. 调用 `createAuction` 创建拍卖

##### 2.2 出价功能

**`bid(uint256 auctionId)`**
- 功能：使用以太坊出价
- 参数：
  - `auctionId`: 拍卖 ID
- 要求：
  - 拍卖未结束
  - 拍卖已开始
  - 出价金额必须大于等于起拍价
  - 出价金额必须大于当前最高出价
- 自动处理：
  - 退还前一个最高出价者的资金到 `pendingReturns` 映射

**`bidWithToken(uint256 auctionId, uint256 bidAmount)`**
- 功能：使用 ERC20 代币出价
- 参数：
  - `auctionId`: 拍卖 ID
  - `bidAmount`: 出价金额
- 要求：
  - 拍卖未结束
  - 拍卖已开始
  - 出价金额必须大于等于起拍价
  - 出价金额必须大于当前最高出价
  - 用户必须有足够的代币余额和授权
- 自动处理：
  - 转移代币到合约
  - 退还前一个最高出价者的资金

##### 2.3 结束拍卖

**`endAuction(uint256 auctionId)`**
- 功能：结束拍卖并分配资产
- 参数：
  - `auctionId`: 拍卖 ID
- 要求：
  - 拍卖未结束
  - 拍卖已结束或卖家提前结束
- 处理逻辑：
  - 如果有出价者：
    - NFT 转移给最高出价者
    - 计算平台手续费
    - 资金分配给卖家和平台
  - 如果没有出价者：
    - NFT 退还给卖家

##### 2.4 提取退款

**`withdraw()`**
- 功能：提取被超过的出价退款
- 要求：
  - 用户有可提取的资金
- 处理：
  - 将 `pendingReturns` 中的资金转移给用户

#### Chainlink 价格预言机集成

##### 2.5 价格计算功能

**`getBidAmountUSD(address token, uint256 amount)`**
- 功能：获取出价金额的美元价值
- 参数：
  - `token`: 代币地址（`address(0)` 表示以太坊）
  - `amount`: 代币数量
- 返回：美元价值（以 8 位小数表示）
- 说明：
  - 使用 Chainlink 价格预言机获取实时价格
  - 自动处理代币精度转换

**`getHighestBidUSD(uint256 auctionId)`**
- 功能：获取当前最高出价的美元价值
- 参数：
  - `auctionId`: 拍卖 ID
- 返回：美元价值

##### 2.6 价格预言机配置

**`setPriceFeed(address token, address priceFeed)`**
- 功能：设置价格预言机地址（仅所有者）
- 参数：
  - `token`: 代币地址（`address(0)` 表示以太坊）
  - `priceFeed`: Chainlink 价格预言机地址
- 要求：
  - 调用者必须是合约所有者
  - 价格预言机地址不能为零地址

#### 管理员功能

**`setPlatformFeeRate(uint256 feeRate)`**
- 功能：设置平台手续费率（仅所有者）
- 参数：
  - `feeRate`: 手续费率（基点，100 = 1%，10000 = 100%）
- 要求：
  - 调用者必须是合约所有者
  - 手续费率必须小于 10000

**`setPlatformFeeRecipient(address recipient)`**
- 功能：设置平台手续费接收地址（仅所有者）
- 参数：
  - `recipient`: 手续费接收地址
- 要求：
  - 调用者必须是合约所有者
  - 接收地址不能为零地址

**`getAuction(uint256 auctionId)`**
- 功能：获取拍卖信息
- 参数：
  - `auctionId`: 拍卖 ID
- 返回：拍卖信息结构体

### 3. 合约升级（UUPS）

#### 3.1 UUPS 代理模式

本项目使用 UUPS（Universal Upgradeable Proxy Standard）代理模式实现合约升级：

- **代理合约**：存储所有状态数据
- **实现合约**：包含业务逻辑
- **升级机制**：通过更新实现合约地址实现升级

#### 3.2 V2 版本新增功能

**AuctionMarketplaceV2** 在 V1 基础上新增了动态手续费功能：

**`setDynamicFeeConfig(...)`**
- 功能：设置动态手续费配置（仅所有者）
- 参数：
  - `threshold1`: 第一个阈值
  - `feeRate1`: 第一个阈值的手续费率
  - `threshold2`: 第二个阈值
  - `feeRate2`: 第二个阈值的手续费率
  - `feeRate3`: 超过第二个阈值的手续费率

**`toggleDynamicFee(bool enabled)`**
- 功能：启用/禁用动态手续费（仅所有者）
- 参数：
  - `enabled`: 是否启用动态手续费

**`calculateDynamicFeeRate(uint256 amount)`**
- 功能：根据金额计算动态手续费率
- 参数：
  - `amount`: 拍卖金额
- 返回：手续费率（基点）
- 逻辑：
  - 金额 < threshold1: 使用 feeRate1
  - threshold1 ≤ 金额 < threshold2: 使用 feeRate2
  - 金额 ≥ threshold2: 使用 feeRate3

**`version()`**
- 功能：获取合约版本号
- 返回：版本号字符串（"2.0.0"）

---

## 技术架构

### 技术栈

- **开发框架**: Hardhat
- **Solidity 版本**: ^0.8.22
- **OpenZeppelin**: 最新版本
  - `ERC721URIStorage`: NFT 标准实现
  - `OwnableUpgradeable`: 权限管理
  - `UUPSUpgradeable`: 代理升级
  - `ReentrancyGuardUpgradeable`: 重入攻击保护
  - `SafeERC20`: 安全的 ERC20 转账
- **Chainlink**: 价格预言机
  - `AggregatorV3Interface`: 价格数据接口
- **测试框架**: Mocha + Chai
- **部署工具**: Hardhat Deploy

### 安全机制

1. **重入攻击保护**
   - 使用 `ReentrancyGuardUpgradeable` 修饰符
   - 所有涉及资金转移的函数都使用 `nonReentrant` 修饰符

2. **权限控制**
   - 使用 `OwnableUpgradeable` 进行管理员权限控制
   - 关键操作仅允许所有者执行

3. **输入验证**
   - 所有用户输入都经过验证
   - 使用 `require` 语句进行条件检查

4. **资金安全**
   - 使用 `SafeERC20` 进行安全的 ERC20 转账
   - 处理非标准 ERC20 代币

5. **时间锁定**
   - 拍卖有明确的时间限制
   - 防止时间相关的攻击

---

## 环境要求

### 系统要求

- Node.js >= 16.0.0
- npm >= 7.0.0
- Git

### 网络要求

- 稳定的网络连接（用于 Fork 模式和测试网部署）
- RPC 端点（推荐使用 Alchemy 或 Infura）

---

## 安装与配置

### 1. 安装依赖

```bash
cd /Volumes/raid1/downloads/go-learning/solidity/Misson3
npm install
```

### 2. 配置环境变量

创建 `.env` 文件：

```env
# RPC 端点配置
# 使用 Alchemy（推荐）
SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY

# 或使用 Infura
# SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_PROJECT_ID

# Fork 配置
# 固定区块号以提高性能（强烈推荐）
FORK_BLOCK_NUMBER=6000000

# 可选：如果网络不稳定，可以禁用 fork
# DISABLE_FORK=true

# 部署到测试网时需要的配置
PRIVATE_KEY=your_private_key_here
ETHERSCAN_API_KEY=your_etherscan_api_key_here
```

### 3. 验证 RPC 连接（可选）

运行 RPC 连接检查脚本：

```bash
npx hardhat run scripts/check-rpc.js
```

---

## 部署步骤

### 1. 编译合约

```bash
npm run compile
```

### 2. 运行测试（推荐）

在部署前运行测试确保合约正常工作：

```bash
npm test
```

### 3. 部署到本地网络

#### 3.1 启动本地节点

```bash
npx hardhat node
```

#### 3.2 部署合约

在另一个终端运行：

```bash
npm run deploy
```

或者：

```bash
npx hardhat run scripts/deploy.js
```

#### 3.3 部署输出示例

```
部署账户: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
账户余额: 10000.0 ETH

部署 NFT 合约...
NFT 合约地址: 0x5FbDB2315678afecb367f032d93F642f64180aa3

部署拍卖市场合约...
拍卖市场合约地址: 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
实现合约地址: 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0

设置 Chainlink 价格预言机...
已设置 ETH/USD 价格预言机: 0x694AA1769357215DE4FAC081bf1f309aDC325306

部署完成！

合约地址汇总:
NFT 合约: 0x5FbDB2315678afecb367f032d93F642f64180aa3
拍卖市场合约（代理）: 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
拍卖市场合约（实现）: 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0
```

### 4. 部署到 Sepolia 测试网

#### 4.1 准备测试网 ETH

确保部署账户有足够的 Sepolia ETH（用于支付 gas 费用）。

#### 4.2 配置环境变量

确保 `.env` 文件中包含：

```env
SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
PRIVATE_KEY=your_private_key_here
ETHERSCAN_API_KEY=your_etherscan_api_key_here
```

#### 4.3 部署合约

```bash
npm run deploy:sepolia
```

或者：

```bash
npx hardhat run scripts/deploy.js --network sepolia
```

#### 4.4 验证合约（可选）

如果配置了 `ETHERSCAN_API_KEY`，可以使用 Hardhat 的验证插件验证合约。

### 5. 配置价格预言机

部署后，需要为每个支持的代币设置价格预言机：

```javascript
// 使用 ethers.js 或 Hardhat Console
const marketplace = await ethers.getContractAt("AuctionMarketplace", MARKETPLACE_ADDRESS);

// 设置 ETH/USD 价格预言机
await marketplace.setPriceFeed(
    ethers.ZeroAddress, // 以太坊使用零地址
    "0x694AA1769357215DE4FAC081bf1f309aDC325306" // Sepolia ETH/USD
);

// 设置 ERC20/USD 价格预言机（例如 USDC）
await marketplace.setPriceFeed(
    ERC20_TOKEN_ADDRESS,
    "0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E" // Sepolia USDC/USD
);
```

### 6. 保存部署信息

部署完成后，请保存以下信息：

- NFT 合约地址
- 拍卖市场合约地址（代理）
- 拍卖市场合约地址（实现）
- 部署网络
- 部署时间
- 部署账户

---

## 使用说明

### 1. 创建拍卖

#### 步骤 1: 铸造 NFT

```javascript
const nftContract = await ethers.getContractAt("AuctionNFT", NFT_ADDRESS);
const tokenId = await nftContract.mintNFT(
    seller.address,
    "https://example.com/metadata/1.json"
);
```

#### 步骤 2: 授权合约

```javascript
await nftContract.connect(seller).approve(MARKETPLACE_ADDRESS, tokenId);
// 或者使用 setApprovalForAll 授权所有 NFT
await nftContract.connect(seller).setApprovalForAll(MARKETPLACE_ADDRESS, true);
```

#### 步骤 3: 创建拍卖

```javascript
const marketplace = await ethers.getContractAt("AuctionMarketplace", MARKETPLACE_ADDRESS);

// 使用以太坊支付
const tx1 = await marketplace.connect(seller).createAuction(
    NFT_ADDRESS,
    tokenId,
    3600, // 1 小时（秒）
    ethers.parseEther("1"), // 起拍价 1 ETH
    ethers.ZeroAddress // 以太坊支付
);
await tx1.wait();

// 使用 ERC20 代币支付
const tx2 = await marketplace.connect(seller).createAuction(
    NFT_ADDRESS,
    tokenId,
    3600, // 1 小时
    ethers.parseEther("100"), // 起拍价 100 代币
    ERC20_TOKEN_ADDRESS // ERC20 代币地址
);
await tx2.wait();
```

### 2. 出价

#### 使用以太坊出价

```javascript
const marketplace = await ethers.getContractAt("AuctionMarketplace", MARKETPLACE_ADDRESS);

const tx = await marketplace.connect(bidder1).bid(auctionId, {
    value: ethers.parseEther("2") // 出价 2 ETH
});
await tx.wait();
```

#### 使用 ERC20 代币出价

```javascript
const erc20Token = await ethers.getContractAt("MockERC20", ERC20_TOKEN_ADDRESS);

// 1. 授权合约使用代币
await erc20Token.connect(bidder1).approve(
    MARKETPLACE_ADDRESS,
    ethers.parseEther("200")
);

// 2. 出价
const marketplace = await ethers.getContractAt("AuctionMarketplace", MARKETPLACE_ADDRESS);
const tx = await marketplace.connect(bidder1).bidWithToken(
    auctionId,
    ethers.parseEther("200") // 出价 200 代币
);
await tx.wait();
```

### 3. 查看拍卖信息

```javascript
const auction = await marketplace.getAuction(auctionId);
console.log("卖家:", auction.seller);
console.log("NFT 合约:", auction.nftContract);
console.log("Token ID:", auction.tokenId.toString());
console.log("起拍价:", ethers.formatEther(auction.startingPrice));
console.log("当前最高出价:", ethers.formatEther(auction.highestBid));
console.log("最高出价者:", auction.highestBidder);
console.log("结束时间:", new Date(Number(auction.endTime) * 1000));
```

### 4. 查看美元价格

```javascript
// 获取出价的美元价值
const usdValue = await marketplace.getBidAmountUSD(
    ethers.ZeroAddress, // 以太坊
    ethers.parseEther("1")
);
console.log("1 ETH =", ethers.formatUnits(usdValue, 8), "USD");

// 获取当前最高出价的美元价值
const highestBidUSD = await marketplace.getHighestBidUSD(auctionId);
console.log("最高出价（USD）:", ethers.formatUnits(highestBidUSD, 8));
```

### 5. 结束拍卖

```javascript
// 等待拍卖结束或卖家提前结束
const tx = await marketplace.connect(seller).endAuction(auctionId);
await tx.wait();
```

### 6. 提取退款

如果出价被超过，可以提取退款：

```javascript
const tx = await marketplace.connect(bidder).withdraw();
await tx.wait();
```

---

## 测试说明

### 运行测试

```bash
npm test
```

### 运行特定测试

```bash
# 运行特定测试套件
npx hardhat test --grep "部署和初始化"

# 运行特定测试用例
npx hardhat test --grep "应该正确初始化合约"
```

### 测试覆盖率

```bash
npm run coverage
```

### 测试说明

测试文件 `test/AuctionMarketplace.test.js` 包含以下测试套件：

1. **部署和初始化**
   - 合约正确初始化
   - 平台手续费率设置
   - 平台手续费接收地址设置

2. **NFT 拍卖功能**
   - 创建拍卖
   - 以太坊出价
   - ERC20 代币出价
   - 结束拍卖
   - 提取退款

3. **Chainlink 价格预言机**
   - 价格计算
   - 美元价值转换

4. **管理员功能**
   - 设置价格预言机
   - 更新手续费率
   - 权限控制

5. **合约升级**
   - UUPS 升级流程
   - V2 版本功能

### 测试注意事项

1. **Fork 模式**：测试使用 Hardhat Fork 模式连接到 Sepolia 测试网，使用真实的 Chainlink 价格预言机。
2. **超时设置**：测试超时时间设置为 180 秒（3 分钟），以适应 Fork 模式的网络延迟。
3. **RPC 端点**：建议使用稳定的 RPC 端点（如 Alchemy 或 Infura）以确保测试稳定性。

---

## 合约升级

### 升级到 V2 版本

#### 步骤 1: 准备升级脚本

确保 `scripts/upgrade.js` 中的代理地址正确。

#### 步骤 2: 设置环境变量

```bash
export PROXY_ADDRESS=0x... # 替换为实际的代理地址
```

#### 步骤 3: 运行升级脚本

```bash
npx hardhat run scripts/upgrade.js --network sepolia
```

#### 步骤 4: 验证升级

```javascript
const marketplaceV2 = await ethers.getContractAt("AuctionMarketplaceV2", PROXY_ADDRESS);
const version = await marketplaceV2.version();
console.log("合约版本:", version); // 应该输出 "2.0.0"
```

### 升级后配置

升级到 V2 后，需要配置动态手续费（可选）：

```javascript
// 设置动态手续费配置
await marketplaceV2.setDynamicFeeConfig(
    ethers.parseEther("1"),    // threshold1: 1 ETH
    100,                        // feeRate1: 1%
    ethers.parseEther("10"),   // threshold2: 10 ETH
    200,                        // feeRate2: 2%
    300                         // feeRate3: 3%
);

// 启用动态手续费
await marketplaceV2.toggleDynamicFee(true);
```

---

## 安全考虑

### 1. 重入攻击保护

- 所有涉及资金转移的函数都使用 `nonReentrant` 修饰符
- 使用 OpenZeppelin 的 `ReentrancyGuardUpgradeable`

### 2. 权限控制

- 关键操作仅允许所有者执行
- 使用 OpenZeppelin 的 `OwnableUpgradeable`

### 3. 输入验证

- 所有用户输入都经过验证
- 使用 `require` 语句进行条件检查

### 4. 资金安全

- 使用 `SafeERC20` 进行安全的 ERC20 转账
- 处理非标准 ERC20 代币

### 5. 时间锁定

- 拍卖有明确的时间限制
- 防止时间相关的攻击

### 6. 升级安全

- 使用 UUPS 代理模式
- 升级需要所有者授权
- 升级后保留所有状态数据

---

## 常见问题

### 1. 测试超时

**问题**：测试运行时出现超时错误。

**解决方案**：
- 检查 RPC 端点是否稳定
- 使用 Alchemy 或 Infura 等专业 RPC 服务
- 设置 `FORK_BLOCK_NUMBER` 固定区块号
- 增加测试超时时间

### 2. 价格预言机连接失败

**问题**：无法连接到 Chainlink 价格预言机。

**解决方案**：
- 检查 RPC 端点配置
- 运行 `scripts/check-rpc.js` 检查连接
- 确保使用正确的价格预言机地址

### 3. 部署失败

**问题**：部署到测试网失败。

**解决方案**：
- 检查账户余额是否足够支付 gas 费用
- 检查网络配置是否正确
- 检查私钥是否正确配置

### 4. 升级失败

**问题**：合约升级失败。

**解决方案**：
- 检查代理地址是否正确
- 确保调用者是合约所有者
- 检查新版本合约是否兼容旧版本

### 5. NFT 授权失败

**问题**：创建拍卖时提示 NFT 未授权。

**解决方案**：
- 确保已调用 `approve` 或 `setApprovalForAll`
- 检查授权地址是否为拍卖市场合约地址
- 检查 NFT 所有者是否正确

---

## 附录

### A. Chainlink 价格预言机地址

#### Sepolia 测试网

- **ETH/USD**: `0x694AA1769357215DE4FAC081bf1f309aDC325306`
- **BTC/USD**: `0x1b44F3514812e835E1D5Ac89b8B52F4C6d5F5c5A`
- **LINK/USD**: `0xc59E3633BAAC79493d908e63626716e8A18Ad701`
- **USDC/USD**: `0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E`

更多价格预言机地址请参考 [Chainlink 文档](https://docs.chain.link/data-feeds/price-feeds/addresses)。

### B. 合约接口

#### AuctionNFT

```solidity
function mintNFT(address recipient, string memory tokenUri) external returns (uint256);
function batchMintNFT(address recipient, string[] memory tokenUris) external returns (uint256[] memory);
function totalSupply() external view returns (uint256);
function nextTokenId() external view returns (uint256);
function tokenURI(uint256 tokenId) external view returns (string memory);
```

#### AuctionMarketplace

```solidity
function createAuction(
    address nftContract,
    uint256 tokenId,
    uint256 duration,
    uint256 startingPrice,
    address paymentToken
) external returns (uint256);

function bid(uint256 auctionId) external payable;
function bidWithToken(uint256 auctionId, uint256 bidAmount) external;
function endAuction(uint256 auctionId) external;
function withdraw() external;

function getBidAmountUSD(address token, uint256 amount) external view returns (uint256);
function getHighestBidUSD(uint256 auctionId) external view returns (uint256);
function getAuction(uint256 auctionId) external view returns (Auction memory);

function setPriceFeed(address token, address priceFeed) external;
function setPlatformFeeRate(uint256 feeRate) external;
function setPlatformFeeRecipient(address recipient) external;
```

#### AuctionMarketplaceV2

```solidity
function setDynamicFeeConfig(
    uint256 threshold1,
    uint256 feeRate1,
    uint256 threshold2,
    uint256 feeRate2,
    uint256 feeRate3
) external;

function toggleDynamicFee(bool enabled) external;
function calculateDynamicFeeRate(uint256 amount) external view returns (uint256);
function version() external pure returns (string memory);
```

### C. 相关文档

- [README.md](./README.md) - 项目快速开始指南
- [CHAINLINK_SETUP.md](./CHAINLINK_SETUP.md) - Chainlink 配置指南
- [RPC_SETUP_GUIDE.md](./RPC_SETUP_GUIDE.md) - RPC 端点配置指南
- [DEPLOYMENT.md](./DEPLOYMENT.md) - 部署详细说明
- [PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md) - 项目总结

### D. 许可证

MIT License

### E. 联系方式

如有问题，请提交 Issue 或联系项目维护者。

---

**文档版本**: 1.0.0  
**最后更新**: 2024-12-21

