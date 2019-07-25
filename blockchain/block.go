package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Index        int
	Data         []byte
	TimeStamp    int64
	PreviousHash []byte
	Hash         []byte
	Nonce int
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

func (b *Block) serialize() []byte{
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)

	if err != nil{
		log.Panic(err)
	}

	return res.Bytes()
}

func createBlock(data string, prevHash []byte) *Block{
	block := &Block{
		TimeStamp: time.Now().Unix(),
		Data:      []byte(data),
		//PreviousHash: prevHash,
		//Nonce: 0,
	}
	pow := NewProofWork(block)
	block.PreviousHash = prevHash
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func deserialize(data []byte) *Block{
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	if err != nil{
		log.Panicln(err)
	}

	return &block
}



