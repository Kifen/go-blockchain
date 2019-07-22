package blockchain

import "errors"


type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) (*Block, error){
	prevBlock := bc.Blocks[(len(bc.Blocks))-1]
	index := prevBlock.Index + 1
	newBlock := bc.createBlock(data, index)

	if newBlock.isBlockValid(prevBlock){
		bc.Blocks = append(bc.Blocks, newBlock)
		return newBlock, nil
	}

	return nil, errors.New("Invalid Block...")
}

func (bc *Blockchain)createBlock(data string, index int) *Block{
	prevBlock := bc.Blocks[(len(bc.Blocks)) - 1]
	prevHash := prevBlock.Hash
	newBlock := newBlock(data,prevHash, index)
	return newBlock
}

func InitBlockchain() *Blockchain{
	return &Blockchain{[]*Block{genesisBlock()}}
}


func genesisBlock() *Block {
	return newBlock("Genesis Block", []byte{}, 0)
}


