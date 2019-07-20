package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type Block struct {
	Index     int    // the position of the data record in the blockchain
	Timestamp string //the time the data is written
	Data      string
	Hash      string // a SHA256 identifier representing this data record
	PrevHash  string // hash of the previous block
}

func (b *Block) DeriveHash() string {
	record := b.PrevHash + b.Timestamp + string(b.Index) + b.Data
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (b *Block) IsBlockValid(oldBlock *Block) bool {
	if oldBlock.Index+1 != b.Index {
		return false
	}

	if oldBlock.Hash != b.PrevHash {
		return false
	}

	if b.DeriveHash() != b.Hash {
		return false
	}

	return true
}

func CreateBlock(data string, oldBLock *Block) (*Block, error) {
	newBLock := &Block{}

	t := time.Now()
	newBLock.Index = oldBLock.Index + 1
	newBLock.Timestamp = t.String()
	newBLock.Data = data
	newBLock.PrevHash = oldBLock.Hash
	newBLock.Hash = newBLock.DeriveHash()

	return newBLock, nil
}

type blockchain struct {
	Blocks []*Block
	mu    sync.RWMutex
}

var (
	b *blockchain = &blockchain{}
)

func (b *blockchain) NewBlock(data string) (*Block, *Block, error){
	prevBlock := b.Blocks[len(b.Blocks)-1]
	newBlock, err := CreateBlock(data, prevBlock)
	//b.Blocks = append(b.Blocks, newBlock)

	return newBlock, prevBlock, err
}

func genesis() *Block {
	t := time.Now().String()
	index := -1
	genesis := &Block{
		Index:     index,
		Timestamp: t,
	}
	block, _ := CreateBlock("Genesis", genesis)
	return block
}

func InitBlockchain() *blockchain{
	newBLock := genesis()
	b.Blocks = append(b.Blocks, newBLock)
	return b
}

func ReplaceChain(newBlocks[] *Block){
	if len(newBlocks) > len(Blockchain().Blocks){
		b.Blocks = newBlocks
	}

}

func Blockchain() *blockchain{
	return b
}
