package main

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Kifen/go-blockchain/blockchain"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		// Init the blockchain and create genesis block
		genesisBlock := blockchain.InitBlockchain()
		spew.Dump(genesisBlock)
	}()
	log.Fatal(run())

}



func run() error{
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", httpAddr)

	srv := &http.Server{
		Addr: ":" + httpAddr,
		Handler: mux,
		ReadHeaderTimeout: 10*time.Second,
		WriteTimeout: 10*time.Second,
		MaxHeaderBytes: 1<<20,
	}

	if err := srv.ListenAndServe(); err != nil{
		return err
	}

	return nil
}

type Message struct {
	Data string
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleWriteBlock(writer http.ResponseWriter, request *http.Request) {
	var m Message

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&m); err != nil{
		respondWithJSON(writer, request, http.StatusBadRequest, request.Body)
		return
	}

	defer request.Body.Close()

	blockChain := blockchain.Blockchain()
	newBlock, prevBlock, err := blockChain.NewBlock(m.Data)

	if err != nil{
		respondWithJSON(writer, request, http.StatusInternalServerError, m)
		return
	}

	//prevBlock := blockChain.Blocks[len(blockChain.Blocks)-1]
	if newBlock.IsBlockValid(prevBlock){
		newBlocks := append(blockChain.Blocks, newBlock)
		blockChain.ReplaceChain(newBlocks)
		spew.Dump(blockChain.Blocks)
	}
}

func respondWithJSON(writer http.ResponseWriter, request *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	writer.WriteHeader(code)
	writer.Write(response)
}

func handleGetBlockchain(writer http.ResponseWriter, request *http.Request){
	blockchain := blockchain.Blockchain()
	bytes, err := json.MarshalIndent(blockchain, "", " ")
	if err != nil{
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(writer, string(bytes))
}
