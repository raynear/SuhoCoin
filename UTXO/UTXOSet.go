package utxo

import (
	"SuhoCoin/block"
	"SuhoCoin/blockchain"
	"SuhoCoin/config"
	"SuhoCoin/transaction"
	"SuhoCoin/util"
	"SuhoCoin/wallet"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/btcsuite/btcutil/base58"
	"github.com/syndtr/goleveldb/leveldb"
)

type UTXO struct {
	Blockchain *blockchain.Blockchain
}

func NewUTXOTransaction(Wallet *wallet.Wallet, to string, amount int, UTXO *UTXO) *transaction.Tx {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput

	pubKeyHash := wallet.HashPubKey(Wallet.PublicKey)
	acc, validOutputs := UTXO.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, e := hex.DecodeString(txid)
		err.ERR("Decode Error", e)

		for _, out := range outs {
			input := transaction.TXInput{TxID: txID, Vout: out, Signature: nil, PubKey: Wallet.PublicKey, ScriptSig: "", Data: ""}
			inputs = append(inputs, input)
		}
	}

	from := Wallet.GetAddress()
	outputs = append(outputs, *transaction.NewTXO(amount, to))
	if acc > amount {
		outputs = append(outputs, *transaction.NewTXO(acc-amount, from)) // a change
	}

	tx := transaction.Tx{ID: nil, Vin: inputs, Vout: outputs}
	tx.SetID()
	UTXO.Blockchain.SignTransaction(&tx, Wallet.PrivateKey)

	return &tx
}

func (u UTXO) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.UTXODB

	iter := db.NewIterator(nil, nil)

	for ok := iter.First(); ok; ok = iter.Next() {
		k := iter.Key()
		v := iter.Value()
		txID := hex.EncodeToString(k)
		outs := transaction.DeserializeOutputs(v)

		for outIdx, out := range outs.Outputs {
			fmt.Println("    pubKeyHash", base58.Encode(pubKeyHash))
			fmt.Println("out.PubKeyHash", base58.Encode(out.PubKeyHash))
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
			}
		}
	}

	return accumulated, unspentOutputs
}

func (u UTXO) FindUTXO(pubKeyHash []byte) []transaction.TXOutput {
	var UTXOs []transaction.TXOutput
	db := u.Blockchain.UTXODB

	iter := db.NewIterator(nil, nil)

	for ok := iter.First(); ok; ok = iter.Next() {
		//		k := iter.Key()
		v := iter.Value()
		outs := transaction.DeserializeOutputs(v)

		for _, out := range outs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (u UTXO) CountTransactions() int {
	db := u.Blockchain.UTXODB
	counter := 0

	iter := db.NewIterator(nil, nil)

	for ok := iter.First(); ok; ok = iter.Next() {
		counter++
	}

	return counter
}

func (u UTXO) Reindex() {
	db := u.Blockchain.UTXODB

	db.Close()
	e := os.RemoveAll("./" + config.V.GetString("Default_db") + "_UTXO")
	err.ERR("del error", e)

	db, e = leveldb.OpenFile(config.V.GetString("Default_db")+"_UTXO", nil)
	u.Blockchain.UTXODB = db

	UTXO := u.Blockchain.FindUTXO()

	for txID, outs := range UTXO {
		key, e := hex.DecodeString(txID)
		err.ERR("Decode Error", e)

		e = db.Put(key, outs.Serialize(), nil)
		err.ERR("DB Put in Error", e)
	}
}

func (u UTXO) Update(block *block.Block) {
	db := u.Blockchain.UTXODB

	for _, tx := range block.Transactions {
		if tx.IsCoinbase() == false {
			for _, vin := range tx.Vin {
				updatedOuts := transaction.TXOutputs{}
				outsBytes, e := db.Get(vin.TxID, nil)
				err.ERR("TXID Get Error", e)
				outs := transaction.DeserializeOutputs(outsBytes)

				for outIdx, out := range outs.Outputs {
					if outIdx != vin.Vout {
						updatedOuts.Outputs = append(updatedOuts.Outputs, out)
					}
				}

				if len(updatedOuts.Outputs) == 0 {
					err := db.Delete(vin.TxID, nil)
					if err != nil {
						log.Panic(err)
					}
				} else {
					err := db.Put(vin.TxID, updatedOuts.Serialize(), nil)
					if err != nil {
						log.Panic(err)
					}
				}

			}
		}

		newOutputs := transaction.TXOutputs{}
		for _, out := range tx.Vout {
			newOutputs.Outputs = append(newOutputs.Outputs, out)
		}

		err := db.Put(tx.ID, newOutputs.Serialize(), nil)
		if err != nil {
			log.Panic(err)
		}
	}
}
