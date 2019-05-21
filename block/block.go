package block

import (
    "test/block/blockheader"
    "test/block/config"
    "time"
)

type Block struct {
    Header       blockheader.BlockHeader
    TxCnt        int64
    Transactions [][]byte
    Data         string
}

func NewBlock(data string, prevBlockHash []byte, height int64, difficulty int64, nonce []byte, merkleRoot []byte) *Block {
    block := &Block{blockheader.BlockHeader{config.BlockchainVersion, []byte{}, prevBlockHash, height, time.Now().Unix(), difficulty, nonce, merkleRoot}, 0, [][]byte{}, data}

    block.Header.SetHash()

    return block
}

func GenesisBlock() *Block {
    return NewBlock("GenesisBlock", []byte{}, int64(0), int64(0), []byte{}, []byte{})
}
