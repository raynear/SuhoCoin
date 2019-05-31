package transaction

import (
	"SuhoCoin/config"
	"SuhoCoin/util"
	"SuhoCoin/wallet"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
	Signature []byte
	PubKey    []byte
	Data      string
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

type TXOutput struct {
	Value        int
	PubKeyHash   []byte
	ScriptPubKey string
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := base58.Decode(string(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

type Tx struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Tx) IsCoinbase() bool {
	return false
}

func CoinbaseTx(to []byte, data string) *Tx {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{Txid: []byte{}, Vout: -1, Signature: []byte{}, PubKey: []byte{}, Data: data}
	txout := TXOutput{Value: config.V.GetInt("Reward"), PubKeyHash: to}
	tx := Tx{ID: nil, Vin: []TXInput{txin}, Vout: []TXOutput{txout}}
	tx.SetID()
	return &tx
}

func (tx *Tx) Print() {
	fmt.Printf("ID(%x)\n", tx.ID)
	for _, aVin := range tx.Vin {
		fmt.Println("aVin")
		fmt.Printf("  Txid(%x) ", aVin.Txid)
		fmt.Printf("Vout(%d) ", aVin.Vout)
		fmt.Printf("ScriptSig(%s) ", aVin.ScriptSig)
		fmt.Printf("Signature(%x) ", aVin.Signature)
		fmt.Printf("PubKey(%x) ", aVin.PubKey)
		fmt.Printf("Data(%s) ", aVin.Data)
		fmt.Println()
	}
	for _, aVout := range tx.Vout {
		fmt.Println("aVout")
		fmt.Printf("  Value(%d) ", aVout.Value)
		fmt.Printf("PubKeyHash(%x) ", aVout.PubKeyHash)
		fmt.Printf("ScriptPubKey(%s) ", aVout.ScriptPubKey)
		fmt.Println()
	}

	fmt.Println()
}

func (tx *Tx) SetID() {
	txPayload := tx.Serialize()

	hash := sha256.Sum256(txPayload)

	tx.ID = hash[:]
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
