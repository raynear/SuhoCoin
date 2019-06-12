package transaction

import (
	"SuhoCoin/config"
	"SuhoCoin/util"
	"SuhoCoin/wallet"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"github.com/syndtr/goleveldb/leveldb"
)

type TXInput struct {
	TxID      []byte
	Vout      int
	Signature []byte
	PubKey    []byte

	ScriptSig string
	Data      string
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

type TXOutput struct {
	Value      int
	PubKeyHash []byte

	ScriptPubKey string
}

func (out *TXOutput) Print() {
	fmt.Println(" aVout")
	fmt.Printf("  Value(%d) ", out.Value)
	fmt.Printf("PubKeyHash(%s) ", base58.Encode(out.PubKeyHash))
	fmt.Printf("ScriptPubKey(%s) ", out.ScriptPubKey)
	fmt.Println()
}

func (out *TXOutput) Lock(address []byte) {
	fmt.Println("address : ", string(address[:]))
	pubKeyHash := base58.Decode(string(address[:]))
	pubKeyHash = pubKeyHash[2 : len(pubKeyHash)-4]
	fmt.Println("pubKeyHash : ", base58.Encode(pubKeyHash[:]))
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func NewTXO(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil, ""}
	txo.Lock([]byte(address))

	return txo
}

type Tx struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Tx) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxID) == 0 && tx.Vin[0].Vout == -1
}

func CoinbaseTx(to string, data string) *Tx {
	if data == "" {
		data = fmt.Sprintf("Coinbase Reward to '%s'", to)
	}

	txin := TXInput{TxID: []byte{}, Vout: -1, Signature: []byte{}, PubKey: []byte{}, Data: data}
	txout := TXOutput{Value: config.V.GetInt("Reward"), PubKeyHash: nil, ScriptPubKey: string(to)}
	txout.Lock([]byte(to))
	tx := Tx{ID: nil, Vin: []TXInput{txin}, Vout: []TXOutput{txout}}
	tx.SetID()
	return &tx
}

func (tx *Tx) TrimmedCopy() Tx {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.TxID, vin.Vout, nil, nil, "", ""})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash, ""})
	}

	return Tx{tx.ID, inputs, outputs}
}

func (tx *Tx) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Tx) {
	if tx.IsCoinbase() {
		return
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.TxID)].ID == nil {
			log.Panic("ERROR: Previous tx is not correct")
		}
	}

	txCopy := tx

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.TxID)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, e := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
		util.ERR("signing error", e)

		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
		txCopy.Vin[inID].PubKey = nil
	}
}

func (tx *Tx) Verify(prevTXs map[string]Tx) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.TxID)].ID == nil {
			log.Panic("ERROR: Previous tx is not correct")
		}
	}

	txCopy := tx

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.TxID)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Vin[inID].PubKey = nil
	}

	return true
}

func (tx *Tx) Print() {
	fmt.Printf("ID(%x)\n", tx.ID)
	for _, aVin := range tx.Vin {
		fmt.Println(" aVin")
		fmt.Printf("  TxID(%x) ", aVin.TxID)
		fmt.Printf("Vout(%d) ", aVin.Vout)
		fmt.Printf("ScriptSig(%s) ", aVin.ScriptSig)
		fmt.Printf("Signature(%s) ", base58.Encode(aVin.Signature))
		fmt.Printf("PubKey(%x) ", aVin.PubKey)
		fmt.Printf("Data(%s) ", aVin.Data)
		fmt.Println()
	}
	for _, aVout := range tx.Vout {
		fmt.Println(" aVout")
		fmt.Printf("  Value(%d) ", aVout.Value)
		fmt.Printf("PubKeyHash(%s) ", base58.Encode(aVout.PubKeyHash))
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

func GetTxFromDB(TxPoolDB *leveldb.DB) []*Tx {
	var Txs []*Tx
	iter := TxPoolDB.NewIterator(nil, nil)

	for iter.Next() {
		value := iter.Value()

		NewTx := DeserializeTx(value)
		Txs = append(Txs, NewTx)
	}

	return Txs
}

func (tx *Tx) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	e := encoder.Encode(tx)

	util.ERR("Tx Encode Error", e)

	return result.Bytes()
}

func DeserializeTx(txb []byte) *Tx {
	var tx Tx

	decoder := gob.NewDecoder(bytes.NewReader(txb))
	e := decoder.Decode(&tx)

	util.ERR("Decode Error", e)

	return &tx
}

type TXOutputs struct {
	Outputs []TXOutput
}

func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
