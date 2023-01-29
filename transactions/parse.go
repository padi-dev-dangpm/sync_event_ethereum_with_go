package transactions

import (
	"context"
	"encoding/hex"
	"math/big"

	// "fmt"
	"log"
	"strings"
	transferAbi "try-ethereum/contracts/transferNFT"
	utils "try-ethereum/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type (
	TransactionData struct {
		Hash             string
		ChainId          *big.Int
		Value            *big.Int
		From             string
		To               string
		Gas              uint64
		GasPrice         *big.Int
		Nonce            uint64
		TransactionData  string
		Logs             []*types.Log
		BlockNumber      uint64
		TransactionIndex uint64
	}

	LogTransfer struct {
		From    common.Address
		To      common.Address
		TokenId *big.Int
	}

	LogTransferSingle struct {
		Operator common.Address
		From     common.Address
		To       common.Address
		Id       *big.Int
		Value    *big.Int
	}

	LogTransferBatch struct {
		Operator common.Address
		From     common.Address
		To       common.Address
		Ids      []*big.Int
		Values   []*big.Int
	}
)

func GetTransactionMessage(tx *types.Transaction) types.Message {
	msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

func ParseTransactionBaseInfo(client *ethclient.Client, tx *types.Transaction) TransactionData {
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	var toAddress string

	utils.TryCatch{
		Try: func() {
			toAddress = GetTransactionMessage(tx).To().Hex()
		},
		Catch: func(e utils.Exception) {
			toAddress = receipt.ContractAddress.Hex()
		},
		Finally: func() {},
	}.Do()

	return TransactionData{
		Hash:             tx.Hash().Hex(),
		ChainId:          tx.ChainId(),
		Value:            tx.Value(),
		From:             GetTransactionMessage(tx).From().Hex(),
		To:               toAddress,
		Gas:              tx.Gas(),
		GasPrice:         tx.GasPrice(),
		Nonce:            tx.Nonce(),
		TransactionData:  hex.EncodeToString(tx.Data()),
		Logs:             receipt.Logs,
		BlockNumber:      receipt.BlockNumber.Uint64(),
		TransactionIndex: uint64(receipt.TransactionIndex),
	}
}

func DecodeTransferLog(logs []*types.Log) []LogTransfer {
	var transferEvents []LogTransfer
	var transferEvent LogTransfer

	transferEventHash := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

	for _, vLog := range logs {
		if strings.Compare(vLog.Topics[0].Hex(), transferEventHash.Hex()) == 0 && len(vLog.Topics) >= 4 {
			func() {
				transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
				transferEvent.TokenId = vLog.Topics[3].Big()

				transferEvents = append(transferEvents, transferEvent)
			}()
		}
	}

	return transferEvents
}

func DecodeTransferSingleLog(logs []*types.Log) []LogTransferSingle {
	var transferEvents []LogTransferSingle
	var transferEvent LogTransferSingle

	transferEventHash := crypto.Keccak256Hash([]byte("TransferSingle(address,address,address,uint256,uint256)"))
	contractAbi, err := abi.JSON(strings.NewReader(string(transferAbi.ContractsABI)))

	for _, vLog := range logs {
		if strings.Compare(vLog.Topics[0].Hex(), transferEventHash.Hex()) == 0 && len(vLog.Topics) >= 4 {
			func() {
				err = contractAbi.UnpackIntoInterface(&transferEvent, "TransferSingle", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}
				transferEvent.Operator = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.From = common.HexToAddress(vLog.Topics[2].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[3].Hex())

				transferEvents = append(transferEvents, transferEvent)
			}()
		}
	}

	return transferEvents
}

func DecodeTransferBatchLog(logs []*types.Log) []LogTransferBatch {
	var transferEvents []LogTransferBatch
	var transferEvent LogTransferBatch

	transferEventHash := crypto.Keccak256Hash([]byte("TransferBatch(address,address,address,uint256[],uint256[])"))
	contractAbi, err := abi.JSON(strings.NewReader(string(transferAbi.ContractsABI)))

	for _, vLog := range logs {
		if strings.Compare(vLog.Topics[0].Hex(), transferEventHash.Hex()) == 0 && len(vLog.Topics) >= 4 {
			func() {
				err = contractAbi.UnpackIntoInterface(&transferEvent, "TransferBatch", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}
				transferEvent.Operator = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.From = common.HexToAddress(vLog.Topics[2].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[3].Hex())

				transferEvents = append(transferEvents, transferEvent)
			}()
		}
	}

	return transferEvents
}
