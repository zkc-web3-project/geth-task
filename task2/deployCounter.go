package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	counter "github.com/zkc-web3-project/geth-task/task2/bindings" // 引入生成的 Go 绑定包
)

func main() {
	initEnv()
	rpcUrl := os.Getenv("RPC_URL")
	pk := os.Getenv("PRIVATE_KEY")

	///////////////////////////////////////部署合约////////////////////////////////////

	// 1. 连接到以太坊节点（Sepolia）
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 加载私钥
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 构建交易签名器
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法转换公钥类型")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 4. 设置交易参数
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // 无 ETH 转账
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	// 5. 部署合约（传递构造函数参数）  后面调用合约时可注释此行
	deployContract(auth, client)

	/////////////////////////////////////////调用合约/////////////////////////////////////
	address2 := common.HexToAddress("0x883f2fbbbc7b03fcb40965fb4f0d7224b56ad7a9") //步骤5中部署的合约地址
	counterInstance, err := counter.NewCounter(address2, client)
	if err != nil {
		log.Fatal(err)
	}
	_, err = counterInstance.Increment(auth)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("调用Increment成功！")
	counter, _ := counterInstance.Count(nil)
	fmt.Println("当前计数为:", counter)
}

/**
* 部署合约
 */
func deployContract(auth *bind.TransactOpts, client *ethclient.Client) {
	address, tx, _, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	// 6. 输出结果
	fmt.Printf("合约地址: 0x%x\n", address)       //0x883f2fbbbc7b03fcb40965fb4f0d7224b56ad7a9
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex()) //0x3d2e3d7e9515c28664125b9e2446f401f02989d04b22613190494875b82350b5

	//等待一下
	time.Sleep(30 * time.Second)

	// 7. 验证合约是否部署成功
	code, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	if len(code) == 0 {
		log.Fatal("合约部署失败")
	}
	fmt.Println("合约部署成功！")
}

/*
初始化环境变量
*/
func initEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
