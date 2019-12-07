package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const targetBits = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func (pow *ProofOfWork) Run() (int, []byte){
	var(
		intHash big.Int
		hash [32]byte
	)
	nonce := 0
	fmt.Printf("Mining  block containing data -->\"%s\"\n", pow.Block.Transactions)

	for nonce <= math.MaxInt64{
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		fmt.Println()
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1{
			break
		}else {
			nonce++
		}
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) ValidatePow() bool{
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(pow.Target) == -1
}

func NewProofWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte{
	data := bytes.Join([][]byte{
		pow.Block.PreviousHash,
		pow.Block.HashTransactions(),
		IntToHex(pow.Block.TimeStamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),},
		[]byte{},)
	return  data
}

func IntToHex(num int64) []byte{
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil{
		log.Panic(err)
	}
	return buff.Bytes()
}
