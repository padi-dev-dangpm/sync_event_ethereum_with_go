package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	txUtils "try-ethereum/transactions"
	helper "try-ethereum/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TransferEvent struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func queryAndDecodeTransactionByBlock(client *ethclient.Client, blockNum int) {
	blockNumber := big.NewInt(int64(blockNum))

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	transactions := block.Transactions()

	// for _, transaction := range transactions {
	// 	func(transaction *types.Transaction, client *ethclient.Client) {
	// 		data := txUtils.ParseTransactionBaseInfo(client, transaction)
	// 		contractAddress := data.To
	// 		switch {
	// 		case helper.IsERC721Contract(common.HexToAddress(contractAddress), client) == true:
	// 			fmt.Println(txUtils.DecodeTransferLog(data.Logs))
	// 		case helper.IsERC1155Contract(common.HexToAddress(contractAddress), client) == true:
	// 			fmt.Println(txUtils.DecodeTransferSingleLog(data.Logs))
	// 			fmt.Println(txUtils.DecodeTransferBatchLog(data.Logs))
	// 		}
	// 	}(transaction, client)
	// }

	var wg sync.WaitGroup
	for _, transaction := range transactions {
		wg.Add(1)
		go func(transaction *types.Transaction, client *ethclient.Client) {
			defer wg.Done()
			data := txUtils.ParseTransactionBaseInfo(client, transaction)
			contractAddress := data.To
			switch {
			case helper.IsERC721Contract(common.HexToAddress(contractAddress), client) == true:
				fmt.Println(txUtils.DecodeTransferLog(data.Logs))
			case helper.IsERC1155Contract(common.HexToAddress(contractAddress), client) == true:
				fmt.Println(txUtils.DecodeTransferSingleLog(data.Logs))
				fmt.Println(txUtils.DecodeTransferBatchLog(data.Logs))
			}
		}(transaction, client)
	}
	wg.Wait()
}

func main() {
	client, err := ethclient.Dial("https://data-seed-prebsc-2-s3.binance.org:8545")
	if err != nil {
		log.Fatal(err)
	}

	// blockNumber := big.NewInt(int64(26750437))
	// blockNumber := big.NewInt(int64(26767279))

	defer timeTrack(time.Now(), "Transaction Run")

	var wg sync.WaitGroup
	for i := 26767269; i <= 26767279; i++ {
		wg.Add(1)
		go func(client *ethclient.Client, i int) {
			defer wg.Done()
			queryAndDecodeTransactionByBlock(client, i)
		}(client, i)
	}
	wg.Wait()

	fmt.Println("Done")
}
