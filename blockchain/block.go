package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"sync"
	"time"
)

type Block struct {
	Index        int
	Data         []byte
	TimeStamp    int64
	PreviousHash []byte
	Hash         []byte
	Nonce int
	Mu           sync.RWMutex
}

func (block *Block) calculteHash() [32]byte{
	timeStamp := []byte(strconv.FormatInt(block.TimeStamp, 10))
	headers := bytes.Join([][]byte{block.PreviousHash, block.Data,
		timeStamp}, []byte{})
	hash := sha256.Sum256(headers)
	return hash
}

func (b *Block) isBlockValid(prevBlock *Block) bool{
	if  !bytes.Equal(b.PreviousHash, prevBlock.Hash){
		return false
	}
	if b.Index != prevBlock.Index+1{
		return  false
	}
	return true
}


func createBlock(data string, index int) *Block{
	bc := BlockChain()
	prevBlock := bc.Blocks[(len(bc.Blocks)) - 1]
	prevHash := prevBlock.Hash
	newBlock := newBlock(data,prevHash, index)
	return newBlock
}

func newBlock(data string, prevHash []byte, index int) *Block{
	block := &Block{
		TimeStamp: time.Now().Unix(),
		Data:      []byte(data),
		PreviousHash: prevHash,
		//Nonce: 0,
	}

	block.Index = index
	pow := NewProofWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}



