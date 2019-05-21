package blockchain

import (
    "test/block/block"
)

type Blockchain struct {
    Blocks []*block.Block
}

func (bc *Blockchain) AddBlock(data string) {
    prevBlock := bc.Blocks[len(bc.Blocks)-1]
    newBlock := block.NewBlock(data, prevBlock.Header.Hash, int64(len(bc.Blocks)-1), 0, []byte{}, []byte{})

    bc.Blocks = append(bc.Blocks, newBlock)
}

func NewBlockchain() *Blockchain {
    return &Blockchain{[]*block.Block{block.GenesisBlock()}}
}
