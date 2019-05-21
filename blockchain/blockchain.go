package blockchain

import (
    "SuhoCoin/Consensus/POW"
    "SuhoCoin/block"
)

type Blockchain struct {
    Blocks []*block.Block
}

func GenesisBlock() *block.Block {
    return POW.FindAnswer("GenesisBlock", []byte{}, int64(0), int64(0), []byte{})
}

func (bc *Blockchain) AddBlock(data string) {
    prevBlock := bc.Blocks[len(bc.Blocks)-1]
    newBlock := POW.FindAnswer(data, prevBlock.Header.Hash, int64(len(bc.Blocks)-1), 0, []byte{})

    bc.Blocks = append(bc.Blocks, newBlock)
}

func NewBlockchain() *Blockchain {
    return &Blockchain{[]*block.Block{GenesisBlock()}}
}
