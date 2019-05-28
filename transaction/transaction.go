package transaction

import (
    "SuhoCoin/util"
    "bytes"
    "encoding/gob"
)

type TXInput struct {
    Txid      []byte
    Vout      int
    ScriptSig string
}

type TXOutput struct {
    Value        int
    ScriptPubKey string
}

type Tx struct {
    ID   []byte
    Vin  []TXInput
    Vout []TXOutput
}

func (tx *Tx) Serialize() []byte {
    var result bytes.Buffer

    encoder := gob.NewEncoder(&result)
    e := encoder.Encode(tx)

    err.ERR("Tx Encode Error", e)

    return result.Bytes()
}

func DeserializeTx(txb []byte) *Tx {
    var tx Tx

    decoder := gob.NewDecoder(bytes.NewReader(txb))
    e := decoder.Decode(&tx)

    err.ERR("Decode Error", e)

    return &tx
}
