package main

import (
	"bufio"
	"encoding/json"
	"github.com/Kifen/go-blockchain/blockchain"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	// blockChainServer handles incoming concurrent Blocks
	blockChainServer chan [] *blockchain.Block
	mutex =  &sync.Mutex{}
)

func main(){

	err := godotenv.Load()
	blockChainServer = make(chan []*blockchain.Block)

	if err != nil {
		log.Fatal(err)
	}

	// Init the blockchain and create the genesis block
	genesisBlock := blockchain.InitBlockchain()
	spew.Dump(genesisBlock)


	// start TCP and serve TCP server
	tcpPort := os.Getenv("TCP-PORT")
	server, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil{
		log.Fatal(err)
	}

	log.Println("TCP Server Listening on port :", tcpPort)

	defer server.Close()

	for {
		conn, err := server.Accept()
		log.Println("NEW CONNECTION ACCEPTED...")
		if err != nil{
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(connection net.Conn) {
	defer connection.Close()
	io.WriteString(connection, "Enter new data:")
	scanner := bufio.NewScanner(connection)
	blockChain := blockchain.Blockchain()

	//take in the data from stdin and add it to blockchain
	// after conducting necessary validation
	go func() {
		for scanner.Scan(){
			data := scanner.Text()
			newBlock, prevBlock, err := blockChain.NewBlock(data)
			if err != nil{
				log.Panicln(err)
				continue
			}

			if newBlock.IsBlockValid(prevBlock){
				newBlockChain := append(blockChain.Blocks, newBlock)
				blockChain.ReplaceChain(newBlockChain)
			}

			blockChainServer <- blockChain.Blocks
			io.WriteString(connection, "\nEnter new data:")
		}
	}()

	// simulate receiving the broadcast
	go func() {
		for {
			time.Sleep(30*time.Second)
			blockChain.Mu.Lock()
			output, err := json.Marshal(blockchain.Blockchain().Blocks)
			if err != nil{
				log.Fatal(err)
			}
			blockChain.Mu.Unlock()
			io.WriteString(connection, string((output)))
		}
	}()

	for _= range blockChainServer{
		spew.Dump(blockchain.Blockchain().Blocks)
	}

}
