package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
)

const dbPath  = "./tmp/blocks"

type Blockchain struct {
	PreviousHash []byte
	Database     *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (bc *Blockchain) AddBlock(data string){
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

	newBlock := createBlock(data, lastHash)
	err = bc.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.serialize())
		HandleErr(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		bc.PreviousHash = newBlock.Hash

		return err
	})
	HandleErr(err)
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

func InitBlockchain() *Blockchain{
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	HandleErr(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound{
			fmt.Println("No existing blockchain found...")
			genesis := genesisBlock()
			fmt.Println("Genesis block proved...")
			err = txn.Set(genesis.Hash, genesis.serialize())
			HandleErr(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		}else{
			item, err := txn.Get([]byte("lh"))
			HandleErr(err)
			err = item.Value(func(val []byte) error {
				HandleErr(err)
				lastHash = val

				return err
			})
			return err
		}
	})

	HandleErr(err)
	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

func HandleErr(e error) {
	if e != nil{
		log.Panic(e)
	}
}


func genesisBlock() *Block {
	return createBlock("Genesis Block", []byte{})
}


