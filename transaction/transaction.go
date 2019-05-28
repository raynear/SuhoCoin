package transaction

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
