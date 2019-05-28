package block

import (
    "SuhoCoin/blockheader"
    "SuhoCoin/transaction"
    "bytes"
    "encoding/gob"
    "fmt"
    "runtime"
)

type Block struct {
    Header       blockheader.BlockHeader
    TxCnt        int64
    Transactions []*transaction.Tx
    Data         string
}

func (b *Block) Serialize() []byte {
    var result bytes.Buffer

    encoder := gob.NewEncoder(&result)
    err := encoder.Encode(b)

    if err != nil {
        fmt.Println("{p:Block, f:Serialize} Error", err)
    }

    pc := make([]uintptr, 10)
    runtime.Callers(1, pc)
    f := runtime.FuncForPC(pc[0])
    fmt.Println("currentFunction:", f.Name())

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
