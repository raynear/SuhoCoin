package block

import (
    "SuhoCoin/blockheader"
    "SuhoCoin/merkletree"
    "SuhoCoin/transaction"
    "SuhoCoin/util"
    "bytes"
    "encoding/gob"
)

type Block struct {
    Header       blockheader.BlockHeader
    TxCnt        int64
    Transactions []*transaction.Tx
    Data         string
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
