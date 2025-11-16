# Geth 私有链搭建完整教程

## 目录
1. [环境准备](#1-环境准备)
2. [Geth 安装](#2-geth-安装)
3. [创世区块配置](#3-创世区块配置)
4. [初始化私有链](#4-初始化私有链)
5. [启动私有链](#5-启动私有链)
6. [账户管理](#6-账户管理)
7. [挖矿配置](#7-挖矿配置)
8. [交易操作](#8-交易操作)
9. [网络配置](#9-网络配置)
10. [常见问题](#10-常见问题)

---

## 1. 环境准备

### 1.1 系统要求
- **操作系统**: Linux、macOS 或 Windows
- **内存**: 至少 4GB RAM
- **存储**: 至少 10GB 可用空间
- **网络**: 稳定的网络连接

### 1.2 前置依赖
```bash
# macOS
brew install git

# Ubuntu/Debian
sudo apt update
sudo apt install git build-essential

# CentOS/RHEL
sudo yum install git gcc gcc-c++ make
```

---

## 2. Geth 安装

### 2.1 方法一：预编译二进制文件（推荐）

#### macOS
```bash
# 使用 Homebrew 安装
brew tap ethereum/ethereum
brew install ethereum
```

#### Linux
```bash
# 下载最新版本
wget https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.13.5-916d6a44.tar.gz

# 解压并安装
tar -xzf geth-linux-amd64-1.13.5-916d6a44.tar.gz
sudo cp geth-linux-amd64-1.13.5-916d6a44/geth /usr/local/bin/
```

#### Windows
1. 访问 [Geth 官方下载页面](https://geth.ethereum.org/downloads/)
2. 下载 Windows 版本的 zip 文件
3. 解压到指定目录
4. 将 geth.exe 路径添加到系统环境变量

### 2.2 方法二：源码编译
```bash
# 克隆源码
git clone https://github.com/ethereum/go-ethereum.git
cd go-ethereum

# 编译安装
make geth

# 将编译好的 geth 复制到系统路径
sudo cp build/bin/geth /usr/local/bin/
```

### 2.3 验证安装
```bash
geth version
```

---

## 3. 创世区块配置

### 3.1 创建工作目录
```bash
mkdir ~/private-ethereum
cd ~/private-ethereum
```

### 3.2 创建创世区块配置文件
创建 `genesis.json` 文件：

```json
{
  "config": {
    "chainId": 12345,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "clique": {
      "period": 15,
      "epoch": 30000
    }
  },
  "difficulty": "0x400",
  "gasLimit": "0x8000000",
  "alloc": {
    "0x7df9a875a174b3bc565e6424a0050ebc1b2d1d82": {
      "balance": "300000000000000000000"
    }
  }
}
```

### 3.3 配置参数说明

#### 基础配置
- **chainId**: 链ID，用于区分不同的以太坊网络
- **difficulty**: 初始挖矿难度
- **gasLimit**: 区块gas限制
- **alloc**: 预分配账户和余额

#### 共识算法配置
- **clique**: PoA (Proof of Authority) 共识算法
  - **period**: 出块间隔（秒）
  - **epoch**: 投票周期

#### 硬分叉配置
各种硬分叉的激活区块号，设置为0表示从创世区块开始激活。

---

## 4. 初始化私有链

### 4.1 创建数据目录
```bash
mkdir -p ~/private-ethereum/data
```

### 4.2 初始化创世区块
```bash
geth --datadir ~/private-ethereum/data init ~/private-ethereum/genesis.json
```

### 4.3 验证初始化
```bash
ls ~/private-ethereum/data
# 应该看到 geth/ 和 keystore/ 目录
```

---

## 5. 启动私有链

### 5.1 基础启动命令
```bash
geth --datadir ~/private-ethereum/data \
     --networkid 12345 \
     --http \
     --http.addr "0.0.0.0" \
     --http.port 8545 \
     --http.api "eth,net,web3,personal,miner" \
     --http.corsdomain "*" \
     --ws \
     --ws.addr "0.0.0.0" \
     --ws.port 8546 \
     --ws.api "eth,net,web3,personal,miner" \
     --ws.origins "*" \
     --allow-insecure-unlock \
     --console
```

### 5.2 参数说明
- **--datadir**: 数据目录路径
- **--networkid**: 网络ID，应与创世区块中的chainId一致
- **--http**: 启用HTTP-RPC服务器
- **--http.addr**: HTTP服务器监听地址
- **--http.port**: HTTP服务器端口
- **--http.api**: 启用的API模块
- **--http.corsdomain**: CORS域名设置
- **--ws**: 启用WebSocket-RPC服务器
- **--allow-insecure-unlock**: 允许不安全的账户解锁
- **--console**: 启动交互式控制台

### 5.3 后台运行
```bash
nohup geth --datadir ~/private-ethereum/data \
           --networkid 12345 \
           --http \
           --http.addr "0.0.0.0" \
           --http.port 8545 \
           --http.api "eth,net,web3,personal,miner" \
           --http.corsdomain "*" \
           --allow-insecure-unlock \
           > ~/private-ethereum/geth.log 2>&1 &
```

---

## 6. 账户管理

### 6.1 创建新账户
```javascript
// 在 geth console 中执行
personal.newAccount("your_password")
```

### 6.2 查看账户列表
```javascript
eth.accounts
```

### 6.3 查看账户余额
```javascript
eth.getBalance(eth.accounts[0])
// 或者指定具体地址
eth.getBalance("0x7df9a875a174b3bc565e6424a0050ebc1b2d1d82")
```

### 6.4 解锁账户
```javascript
personal.unlockAccount(eth.accounts[0], "your_password", 0)
```

### 6.5 设置挖矿账户
```javascript
miner.setEtherbase(eth.accounts[0])
```

---

## 7. 挖矿配置

### 7.1 开始挖矿
```javascript
// 开始挖矿，参数为挖矿线程数
miner.start(1)
```

### 7.2 停止挖矿
```javascript
miner.stop()
```

### 7.3 查看挖矿状态
```javascript
eth.mining
```

### 7.4 查看算力
```javascript
eth.hashrate
```

### 7.5 查看区块信息
```javascript
// 查看最新区块号
eth.blockNumber

// 查看指定区块信息
eth.getBlock(1)

// 查看最新区块
eth.getBlock("latest")
```

---

## 8. 交易操作

### 8.1 发送交易
```javascript
// 基础转账交易
eth.sendTransaction({
    from: eth.accounts[0],
    to: "0x目标地址",
    value: web3.toWei(1, "ether"),
    gas: 21000,
    gasPrice: web3.toWei(20, "gwei")
})
```

### 8.2 查看交易
```javascript
// 根据交易哈希查看交易
eth.getTransaction("0x交易哈希")

// 查看交易收据
eth.getTransactionReceipt("0x交易哈希")
```

### 8.3 查看交易池
```javascript
txpool.status
txpool.inspect
```

### 8.4 估算Gas费用
```javascript
eth.estimateGas({
    from: eth.accounts[0],
    to: "0x目标地址",
    value: web3.toWei(1, "ether")
})
```

---

## 9. 网络配置

### 9.1 多节点网络搭建

#### 节点1配置
```bash
geth --datadir ~/private-ethereum/node1 \
     --networkid 12345 \
     --port 30303 \
     --http \
     --http.port 8545 \
     --http.api "eth,net,web3,personal,miner" \
     --allow-insecure-unlock \
     --console
```

#### 节点2配置
```bash
geth --datadir ~/private-ethereum/node2 \
     --networkid 12345 \
     --port 30304 \
     --http \
     --http.port 8546 \
     --http.api "eth,net,web3,personal,miner" \
     --allow-insecure-unlock \
     --console
```

### 9.2 节点连接
```javascript
// 在节点1中获取enode信息
admin.nodeInfo.enode

// 在节点2中连接节点1
admin.addPeer("enode://节点1的enode信息")

// 查看连接的节点
admin.peers
```

### 9.3 静态节点配置
创建 `static-nodes.json` 文件：
```json
[
  "enode://节点1的enode@IP:端口",
  "enode://节点2的enode@IP:端口"
]
```

---

## 10. 常见问题

### 10.1 端口被占用
```bash
# 查看端口占用
lsof -i :8545

# 杀死占用进程
kill -9 PID
```

### 10.2 账户解锁失败
确保使用 `--allow-insecure-unlock` 参数启动geth。

### 10.3 挖矿不出块
1. 检查是否设置了挖矿账户
2. 确认账户已解锁
3. 检查网络连接

### 10.4 交易pending状态
1. 检查gas价格是否足够
2. 确认账户余额充足
3. 检查nonce值是否正确

### 10.5 数据目录权限问题
```bash
# 修改数据目录权限
chmod -R 755 ~/private-ethereum/data
```

---

## 附录

### A. 常用命令速查

#### 账户操作
```javascript
personal.newAccount("password")    // 创建账户
eth.accounts                       // 查看账户
eth.getBalance(address)           // 查看余额
personal.unlockAccount(address, "password", 0)  // 解锁账户
```

#### 挖矿操作
```javascript
miner.setEtherbase(address)       // 设置挖矿账户
miner.start(1)                    // 开始挖矿
miner.stop()                      // 停止挖矿
eth.mining                        // 挖矿状态
```

#### 区块链操作
```javascript
eth.blockNumber                   // 最新区块号
eth.getBlock(number)             // 获取区块信息
eth.syncing                      // 同步状态
```

### B. 配置文件模板

#### 开发环境配置
```json
{
  "config": {
    "chainId": 12345,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "clique": {
      "period": 5,
      "epoch": 30000
    }
  },
  "difficulty": "0x1",
  "gasLimit": "0x8000000",
  "alloc": {}
}
```

#### 生产环境配置
```json
{
  "config": {
    "chainId": 54321,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "clique": {
      "period": 15,
      "epoch": 30000
    }
  },
  "difficulty": "0x400",
  "gasLimit": "0x8000000",
  "alloc": {}
}
```

---

**注意**: 本教程适用于学习和开发环境，生产环境部署需要额外的安全配置和优化。