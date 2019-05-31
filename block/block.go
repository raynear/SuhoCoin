package block

import (
	"SuhoCoin/blockheader"
	"SuhoCoin/merkletree"
	"SuhoCoin/transaction"
	"SuhoCoin/util"
	"bytes"
	"encoding/gob"
	"fmt"
)

type Block struct {
	Header       blockheader.BlockHeader
	TxCnt        int64
	Transactions []*transaction.Tx
	Data         string
}

func (b *Block) Print() {
	fmt.Printf("Version(%d) ", b.Header.Version)
	fmt.Printf("Hash(%x) ", b.Header.Hash)
	fmt.Printf("PrevBlockHash(%x) ", b.Header.PrevBlockHash)
	fmt.Printf("Height(%d) ", b.Header.Height)
	fmt.Printf("TimeStamp(%d) ", b.Header.TimeStamp)
	fmt.Printf("Difficulty(%d) ", b.Header.Difficulty)
	fmt.Printf("Nonce(%d) ", b.Header.Nonce)
	fmt.Printf("MerkleRoot(%x) ", b.Header.MerkleRoot)
	fmt.Printf("TxCnt(%d) ", b.TxCnt)
	fmt.Printf("Data(%s)\n", b.Data)
}

func (b *Block) NewTxMerkleTree() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := merkletree.NewMerkleTree(transactions)

	return mTree.Root.Data
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	e := encoder.Encode(b)

	err.ERR("Encode Error", e)

	return result.Bytes()
}

func DeserializeBlock(b []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(b))
	e := decoder.Decode(&block)

	err.ERR("Decode Error", e)

	return &block
}
