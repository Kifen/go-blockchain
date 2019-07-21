package blockchain

import "errors"

var (
	b = &Blockchain{[]*Block{genesisBlock()}}
)

func BlockChain()*Blockchain{
	return b
}

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) (*Block, error){
	prevBlock := bc.Blocks[(len(b.Blocks))-1]
	index := prevBlock.Index + 1
	newBlock := createBlock(data, index)

	if newBlock.isBlockValid(prevBlock){
		bc.Blocks = append(bc.Blocks, newBlock)
		return newBlock, nil
	}

	return nil, errors.New("Invalid Block...")
}


func genesisBlock() *Block {
	return newBlock("Genesis Block", []byte{}, 0)
}


