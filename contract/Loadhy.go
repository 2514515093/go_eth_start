package contract

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	contractAddr = "0x692B8d7a67D75996924EF1FE6c3F011A1FEc97fe"
)

func Load() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/3ac6493099004557a0879796935a9ef1")
	if err != nil {
		log.Fatal(err)
	}
	countContract, err := NewContract(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("==")
	if err != nil {
		log.Fatal(err)
	}
	chainID, _ := client.NetworkID(context.Background())
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	tx, _ := countContract.Add(opt)
	fmt.Println(tx.Hash().Hex())

	callOpt := &bind.CallOpts{Context: context.Background()}
	valueInContract, err := countContract.GetCount(callOpt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("value:", valueInContract)

	receipt, _ := bind.WaitMined(context.Background(), client, tx)

	for _, vLog := range receipt.Logs {
		evt, err := countContract.ParseCount(*vLog)
		if err == nil {
			fmt.Println("counter:", evt.Send)
		}
	}
}
