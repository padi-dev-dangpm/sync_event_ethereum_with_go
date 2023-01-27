package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
    "time"
    "sync"
    "strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/accounts/abi"
	nfts "try-ethereum/contracts/transferNFT"
    txUtils "try-ethereum/transactions"
)


type TransferEvent struct {
    From common.Address
    To common.Address
    TokenId *big.Int
}

func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s", name, elapsed)
}


func main() {
    client, err := ethclient.Dial("https://data-seed-prebsc-1-s3.binance.org:8545")
    if err != nil {
        log.Fatal(err)
    }

    blockNumber := big.NewInt(int64(26416514))
    block, err := client.BlockByNumber(context.Background(), blockNumber)
    if err != nil {
        log.Fatal(err)
    }

    
    contractAbi, err := abi.JSON(strings.NewReader(string(nfts.ContractsABI)))
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(contractAbi)

    transactions := block.Transactions()

    // for _, transaction := range transactions {
    //     func(transaction *types.Transaction){
    //         data := txUtils.ParseTransactionBaseInfo(client,transaction)
    //         fmt.Println(txUtils.DecodeTransferLog(data.Logs))
    //     }(transaction)
    // }

    defer timeTrack(time.Now(), "Transaction Run")
    var wg sync.WaitGroup
    for _, transaction := range transactions {
        wg.Add(1)
        go func(transaction *types.Transaction){
            defer wg.Done()
            data := txUtils.ParseTransactionBaseInfo(client,transaction)
            fmt.Println(txUtils.DecodeTransferLog(data.Logs))
        }(transaction)
    }
    wg.Wait()

    fmt.Println("Done")
}