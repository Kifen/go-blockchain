package blockchain

type TxOutput struct {
	Value int
	PubKey string
}

type TxInput struct {
	ID []byte
	Out int
	Sig string
}

type Transaction struct {
	ID []byte
	Inputs []TxInput
	Outputs []TxOutput
}


