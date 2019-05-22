package block

import (
    "SuhoCoin/blockheader"
    "bytes"
    "encoding/gob"
    "fmt"
)

type Block struct {
    Header       blockheader.BlockHeader
    TxCnt        int64
    Transactions [][]byte
    Data         string
}

func (b *Block) Serialize() []byte {
    var result bytes.Buffer

    encoder := gob.NewEncoder(&result)
    err := encoder.Encode(b)

    if err != nil {
        fmt.Println("{p:Block, f:Serialize} Error", err)
    }

    return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
    var block Block

    decoder := gob.NewDecoder(bytes.NewReader(d))
    err := decoder.Decode(&block)

    if err != nil {
        fmt.Println("{p:Block, f:Deserialize} Error", err)
    }

    return &block
}
