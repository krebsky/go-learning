package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func ReadWrite() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/YoIi2VFhjh1Dri5yhB2ic")
	if err != nil {
		log.Fatal(err)
	}

	BlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(BlockNumber)

	Block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(BlockNumber)))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Block.Hash().Hex())
	fmt.Println(Block.Number().Uint64())
	fmt.Println(Block.Time())
	fmt.Println(Block.Nonce())

}
