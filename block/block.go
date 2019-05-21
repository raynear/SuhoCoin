package block

import (
    "SuhoCoin/blockheader"
)

type Block struct {
    Header       blockheader.BlockHeader
    TxCnt        int64
    Transactions [][]byte
    Data         string
}
