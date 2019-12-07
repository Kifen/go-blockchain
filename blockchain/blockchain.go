package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"os"
	"runtime"
)

const (
	dbPath  = "./tmp/blocks"
	dbFile = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

type Blockchain struct {
	PreviousHash []byte
	Database     *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DBexists() bool{
	if _, err := os.Stat(dbFile); os.IsNotExist(err){
		return false
	}
	return true
}

func (bc *Blockchain) AddBlock(transaction []*Transaction){
	var lastHash []byte

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleErr(err)
		err = item.Value(func(val []byte) error {
			HandleErr(err)
			lastHash = val

			return err
		})
		return err
	})
	HandleErr(err)

	newBlock := createBlock(transaction, lastHash)
	err = bc.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.serialize())
		HandleErr(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		bc.PreviousHash = newBlock.Hash

		return err
	})
	HandleErr(err)
}

func (bc *Blockchain) FindUnspentTransactions(addr string) []Transaction{
	var unspentTxs []Transaction
	spentTxos := make(map[string][]int)
	iter := bc.Iterator()

	for {
		block := iter.Next()
		for _, tx := range block.Transactions{
			txID := hex.EncodeToString(tx.ID)

			Outputs:
				for outIdx, out := range tx.Outputs{
					if spentTxos[txID] != nil{
						for _, spentout := range spentTxos[txID]{
							if spentout == outIdx{
								continue Outputs
							}
						}
					}
					if out.CanBeUnlocked(addr) {
						unspentTxs = append(unspentTxs, *tx)
					}
				}
				if tx.IsCoinbase() == false {
					for _, in := range tx.Inputs {
						if in.CanUnlock(addr) {
							inTxID := hex.EncodeToString(in.ID)
							spentTxos[inTxID] = append(spentTxos[inTxID], in.Out)
						}
					}
				}

		}
		if len(block.PreviousHash) == 0{
			break
		}

	}

	return unspentTxs
}

func (bc *Blockchain) FndUTXO(address string) []TxOutput {
	var UTXOS []TxOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions{
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address){
				UTXOS = append(UTXOS, out)
			}
		}
	}
	return UTXOS
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string] []int){
	unspentOuts := make(map[string][]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

	Work:
		for _, tx := range unspentTxs{
			txID := hex.EncodeToString(tx.ID)

			for outIdx, out := range tx.Outputs{
				if out.CanBeUnlocked(address) && accumulated < amount{
					accumulated += out.Value
					unspentOuts[txID] = append(unspentOuts[txID], outIdx)

					if accumulated >= amount {
						break Work
					}
				}
			}
		}
	return accumulated, unspentOuts
}

func (bc *Blockchain) Iterator() *BlockchainIterator{
	iter := &BlockchainIterator{bc.PreviousHash, bc.Database}
	return iter
}

func (iter *BlockchainIterator) Next() *Block{
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		HandleErr(err)
		err = item.Value(func(val []byte) error {
			HandleErr(err)
			encodedBlock := val
			block = deserialize(encodedBlock)

			return err
		})
		return err
	})
	HandleErr(err)
	iter.CurrentHash = block.PreviousHash

	return block
}

func InitBlockchain(addr string) *Blockchain{
	var lastHash []byte

	if DBexists(){
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	HandleErr(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinBaseTx(addr, genesisData)
		genesis := genesisBlock(cbtx)
		fmt.Println("Genesis Created...")
		err = txn.Set(genesis.Hash, genesis.serialize())
		HandleErr(err)
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err
	})

	HandleErr(err)
	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

func ContinueBlockChain(addr string) *Blockchain{
	if DBexists() == false {
		fmt.Println("No existing blockchain found, create one...")
		runtime.Goexit()
	}
	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	HandleErr(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleErr(err)
		err = item.Value(func(val []byte) error {
			HandleErr(err)
			lastHash = val
			return err
		})
		return err
	})
	HandleErr(err)
	chain := Blockchain{lastHash, db}
	return &chain
}


func HandleErr(e error) {
	if e != nil{
		log.Panic(e)
	}
}


func genesisBlock(coinbase *Transaction) *Block {
	return createBlock([]*Transaction{coinbase}, []byte{})
}


