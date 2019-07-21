package main

import (
	"fmt"
	"github.com/Kifen/go-blockchain/blockchain"
	"strconv"
)

func main(){
	blockChain := blockchain.BlockChain()

	blockChain.AddBlock("Send 1 btc to Ivan")
	blockChain.AddBlock("Send 2 btc ti KIfen")
	blockChain.AddBlock("Send 3 btc to KIta")

	for _, block := range blockChain.Blocks{
		fmt.Println()
		fmt.Printf("Index: %v\n", block.Index)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Prev-hash: %x\n", block.PreviousHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProofWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.ValidatePow()))
		fmt.Println()
	}
}