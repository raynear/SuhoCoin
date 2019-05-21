package blockchain

import (
    "SuhoCoin/Consensus/POW"
    "SuhoCoin/block"
    "SuhoCoin/blockheader"
    "SuhoCoin/config"
    "time"
)

type Blockchain struct {
    Blocks []*block.Block
}

func NewBlock(data string, prevBlockHash []byte, height int64, difficulty int64, merkleRoot []byte) *block.Block {
    block := &block.Block{blockheader.BlockHeader{config.BlockchainVersion, []byte{}, prevBlockHash, height, time.Now().Unix(), difficulty, 0, merkleRoot}, 0, [][]byte{}, data}

    pow := POW.NewPOW(block)
    nonce, hash := pow.Run()

    block.Header.Hash = hash[:]
    block.Header.Nonce = nonce

    return block
}

func GenesisBlock() *block.Block {
    return NewBlock("GenesisBlock", []byte{}, int64(0), int64(0), []byte{})
}

func (bc *Blockchain) AddBlock(data string) {
    prevBlock := bc.Blocks[len(bc.Blocks)-1]
    newBlock := NewBlock(data, prevBlock.Header.Hash, int64(len(bc.Blocks)-1), 0, []byte{})

    bc.Blocks = append(bc.Blocks, newBlock)
}

func NewBlockchain() *Blockchain {
    return &Blockchain{[]*block.Block{GenesisBlock()}}
}
