package utils

import (
	"log"
	erc165 "try-ethereum/contracts/erc165"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func IsERC721Contract(address common.Address, client *ethclient.Client) bool {
	instance, err := erc165.NewErc165(address, client)
	if err != nil {
		log.Fatal(err)
	}

	var interfaceId [4]byte
	copy(interfaceId[:], common.Hex2Bytes("80ac58cd"))

	var success bool = false
	success, _ = instance.SupportsInterface(&bind.CallOpts{}, interfaceId)

	return success
}

func IsERC1155Contract(address common.Address, client *ethclient.Client) bool {
	instance, err := erc165.NewErc165(address, client)
	if err != nil {
		log.Fatal(err)
	}

	var interfaceId [4]byte
	copy(interfaceId[:], common.Hex2Bytes("d9b67a26"))

	var success bool = false
	success, _ = instance.SupportsInterface(&bind.CallOpts{}, interfaceId)

	return success
}
