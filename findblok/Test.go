package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"go_eth_start/contract"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

//编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
//实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
//输出查询结果到控制台。

func main() {
	//findblock()
	//findtx()
	//findtxtwo()
	//threetest()
	//contract.Cs()
	contract.Load()
}

func findblock() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/3ac6493099004557a0879796935a9ef1")
	if err != nil {
		log.Fatal(err)
	}
	blockNumber := big.NewInt(9848921)
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	fmt.Println(header.Number.Uint64())
	fmt.Println(header.Time)
	fmt.Println(header.Difficulty.Uint64())
	fmt.Println(header.Hash().Hex())

	if err != nil {
		log.Fatal(err)
	}
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(block.Number().Uint64())
	fmt.Println(block.Time())
	fmt.Println(block.Difficulty().Uint64())
	fmt.Println(block.Hash().Hex())
	fmt.Println(len(block.Transactions()))
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(count) // 70
}

func findtx() {
	client, _ := ethclient.Dial("https://sepolia.infura.io/v3/3ac6493099004557a0879796935a9ef1")
	chainID, err := client.ChainID(context.Background())
	blockNumber := big.NewInt(9848921)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	fmt.Println(block.Transactions().Len())
	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())
		if err != nil {
			log.Fatal(err)
		}
		if sender, err := types.Sender(types.NewEIP155Signer(chainID), tx); err == nil {
			fmt.Println("sender", sender.Hex())
		}
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(receipt.Status)
	}
}

func findtxtwo() {
	client, _ := ethclient.Dial("https://sepolia.infura.io/v3/3ac6493099004557a0879796935a9ef1")
	chainID, _ := client.ChainID(context.Background())
	blockNumber := big.NewInt(9848921)
	block, _ := client.BlockByNumber(context.Background(), blockNumber)
	fmt.Println(block.Transactions().Len())
	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())
		if tx.Hash().Hex() != "0x12cd0b5a97f835b561b9d71d3e25fbebfb1ed06b6c6ec0e033847fb2f02c7c69" {
			continue
		}
		// ===== from（重点）=====
		from, err := types.Sender(
			types.LatestSignerForChainID(chainID),
			tx,
		)
		if err != nil {
			panic(err)
		}

		// ===== to =====
		var to string
		if tx.To() == nil {
			to = "contract creation"
		} else {
			to = tx.To().Hex()
		}

		// ===== gas price =====
		var gasPrice *big.Int
		if tx.Type() == types.DynamicFeeTxType {
			// EIP-1559 交易
			tip, err := tx.EffectiveGasTip(block.BaseFee())
			if err != nil {
				panic(err)
			}
			gasPrice = new(big.Int).Add(block.BaseFee(), tip)
		} else {
			// legacy 交易
			gasPrice = tx.GasPrice()
		}

		fmt.Println("From:", from.Hex())
		fmt.Println("To:", to)
		fmt.Println("GasPrice (wei):", gasPrice.String())
		fmt.Println("GasPrice (gwei):", new(big.Int).Div(gasPrice, big.NewInt(1e9)))
	}
}

func createQb() {
	privateKey, _ := crypto.HexToECDSA("0e0e0e7978b4997d1dc196c2cd1225e3bdeb3d73e2de72d1facada7583eac461")
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 去掉'0x'
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("from pubKey:", hexutil.Encode(publicKeyBytes)[4:]) // 去掉'0x04'
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println("full:", hexutil.Encode(hash.Sum(nil)[:]))
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 原长32位，截去12位，保留后20位
}

func threetest() {
	client, _ := ethclient.Dial("https://sepolia.infura.io/v3/3ac6493099004557a0879796935a9ef1")

	privateKey, _ := crypto.HexToECDSA("0e0e0e7978b4997d1dc196c2cd1225e3bdeb3d73e2de72d1facada7583eac461")
	publicKey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(publicKey)
	fmt.Println(address.Hex())

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}
	value := big.NewInt(100000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)               // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x36B017A4961A83960EdA1F116038cCa59e4D5276")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
