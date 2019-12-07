package cli

import (
	"flag"
	"fmt"
	"github.com/Kifen/go-blockchain/blockchain"
	"github.com/Kifen/go-blockchain/wallet"
	"log"
	"os"
	"runtime"
	"strconv"
)

type CLI struct {}

func (cli *CLI) printUsage(){
	fmt.Println("Usage:")
	fmt.Println("getbalance -address ADDRESS ->> get the balance for an address")
	fmt.Println("createblockchain -address ADDRESS ->> creates a blockchain")
	fmt.Println("printchain ->> Prints the blocks in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT ->> Send amount of coins")
	fmt.Println("createwallet ->> Creates a new wallet")
	fmt.Println("listaddresses ->> List the addresses in wallet file")
}

func (cli *CLI) validateArgs(){
	if len(os.Args) <2{
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CLI) createBlockChain(address string){
	chain := blockchain.InitBlockchain(address)
	chain.Database.Close()
	fmt.Println("Finished!")
}

func (cli *CLI) listAddresses() {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddesses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CLI) createWallet(){
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)
}

func (cli *CLI) getBalance(address string){
	fmt.Println("IN getBalance()...")
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOS := chain.FndUTXO(address)
	for _, out := range UTXOS {
		balance += out.Value
	}
	fmt.Printf("Balance of %v: %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int){
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([] *blockchain.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CLI) printChain(){
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Previous-hash: %x\n", block.PreviousHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.ValidatePow()))
		fmt.Println()

		if len(block.PreviousHash) == 0{
			break
		}
	}
}

func (cli *CLI) Run(){
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)



	getBalanceAddress := getBalanceCmd.String("address", "", "address balance")
	createBlockChainAddress := createBlockChainCmd.String("address", "", "")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}

	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}

	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}

	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("In case: createwallet 3...")
			log.Panic(err)
		}

	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed(){
		if *getBalanceAddress == ""{
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockChainCmd.Parsed(){
		if *createBlockChainAddress == ""{
			createBlockChainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockChainAddress)
	}

	if sendCmd.Parsed(){
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed(){
		cli.printChain()
	}

	if listAddressesCmd.Parsed(){
		cli.listAddresses()
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
}
