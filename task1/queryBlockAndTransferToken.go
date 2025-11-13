package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	//获取环境变量
	initEnv()
	rpcUrl := os.Getenv("RPC_URL")
	pk := os.Getenv("PRIVATE_KEY")

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}
	///////////////////////获取指定区块的区块信息//////////////////////////////////////
	blockNumber := big.NewInt(9598788)
	blockInfo, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("区块哈希:", blockInfo.Hash().Hex())
	fmt.Println("区块高度:", blockInfo.Number())
	fmt.Println("区块时间:", blockInfo.Time())
	fmt.Println("区块交易数量:", len(blockInfo.Transactions()))
	fmt.Println("区块难度:", blockInfo.Difficulty())
	fmt.Println("交易序号:", blockInfo.Nonce())
	fmt.Println("区块GasLimit:", blockInfo.GasLimit())
	fmt.Println("区块GasUsed:", blockInfo.GasUsed())

	header := blockInfo.Header()
	headerJson, err := json.Marshal(header)
	if err != nil {
		log.Panic("区块头序列化失败", err)
	}
	fmt.Println("区块头json:", string(headerJson))
	///////////////////////////////ERC20转账/////////////////////////////////////////
	/**
	  1. 获取转出方私钥
	  2. 由私钥获取公钥，公钥获取转出方地址
	  3. 获取下一个交易序号
	  4. 获取接收方地址和代币合约地址
	  5. 构造转账数据
	  6. 获取网络id
	  7. 使用私钥签名交易
	  8. 发送交易
	  9. 获取交易hash
	*/

	privateKey, err := crypto.HexToECDSA(pk) //测试账户私钥
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	//打印公钥地址
	fmt.Println("fromAddress-公钥地址：", crypto.PubkeyToAddress(*publicKeyECDSA).Hex())
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	fmt.Println("交易序号nonce:", nonce)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0) //当转账代币时转账金额设置为0
	gasPrice, err := client.SuggestGasPrice(context.Background())
	fmt.Println("获取建议的gas价格gasPrice:", gasPrice)
	if err != nil {
		log.Fatal(err)
	}

	// toAddress := common.HexToAddress("0xb1b29850e895add42661f51cf8cda44280404a3b") //接收方账户地址
	toAddress := common.HexToAddress("0x31dc153670ce03ad2628d23555cc161ac3f62b4e")    //接收方账户地址
	tokenAddress := common.HexToAddress("0x654CBcFAA14608b8b428B2C108C8b281A7acEf24") //代币合约地址

	//构造转账数据
	transferFnSignature := []byte("transfer(address,uint256)") //转账函数的方法签名
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID))
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32) //将地址和金额填充为32字节(ABI编码规范)
	fmt.Println(hexutil.Encode(paddedAddress))
	amount := new(big.Int)
	amount.SetString("5000000000000000000", 10) //转账的代币数量(这里表示5个代币，18位的精度)，第二个参数10代表采用十进制
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount))
	var data []byte
	data = append(data, methodID...)      //函数选择器  注意！！！！这里的三个append都不能交换顺序
	data = append(data, paddedAddress...) //接收方地址
	data = append(data, paddedAmount...)  //代币数量

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	//打印预估的gasLimit
	fmt.Println("预估的gasLimit:", gasLimit)
	if err != nil {
		log.Fatal(err)
	}
	gasLimit = 300000 //手动重置，预估的不准，可能造成交易失败
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	//打印链ID和交易签名
	fmt.Println("链ID:", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("交易签名:", signedTx)
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("获取到交易hash:", signedTx.Hash().Hex())
}

func initEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
