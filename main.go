package main

import (
	"flag"
	"fmt"
	"github.com/Kifen/go-blockchain/blockchain"
	"os"
	"runtime"
	"strconv"
)

func main(){
	blockChain := blockchain.InitBlockchain()

	defer blockChain.Database.Close()

	cli := CLI{blockChain}
	cli.run()
}

type CLI struct {
	blockchain *blockchain.Blockchain
}

func (cli *CLI) printUsage(){
	fmt.Println("Usage:")
	fmt.Println("add -block BLOCK_DATA - add a block to the chain")
	fmt.Println("print - Prints the blocks in the chain")
}

func (cli *CLI) validateArgs(){
	if len(os.Args) <2{
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CLI) addBlock(data string){
	cli.blockchain.AddBlock(data)
	fmt.Println("Block added!")
}

func (cli *CLI) printChain(){
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Previous-hash: %x\n", block.PreviousHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.ValidatePow()))
		fmt.Println()

		if len(block.PreviousHash) == 0{
			break
		}
	}
}

func (cli *CLI) run(){
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "BLock data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.HandleErr(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.HandleErr(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed(){
		if *addBlockData == ""{
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed(){
		cli.printChain()
	}
}