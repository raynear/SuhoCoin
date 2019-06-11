package blockchain

import (
	"SuhoCoin/Consensus/POW"
	"SuhoCoin/block"
	"SuhoCoin/config"
	"SuhoCoin/transaction"
	"SuhoCoin/util"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

type Blockchain struct {
	LastBlockHash []byte
	DB            *leveldb.DB
	UTXODB        *leveldb.DB
	TxPoolDB      *leveldb.DB
}

func GenesisBlock(coinbase *transaction.Tx) *block.Block {
	return POW.FindAnswer("GenesisBlock", []byte{}, int64(0), config.V.GetInt64("TargetBits"), []byte{}, []*transaction.Tx{coinbase})
}

func NewBlockchain() *Blockchain {
	var LastBlockHash []byte
	db, e := leveldb.OpenFile(config.V.GetString("Default_db"), nil)
	util.ERR("NewBlockchain open DB Error", e)

	LastBlockHash, e = db.Get([]byte("l"), nil)
	if e != nil {
		cbtx := transaction.CoinbaseTx(config.V.GetString("Coinbase"), "GENESIS of SuhoCoin")
		db.Put(append([]byte("t"), cbtx.ID...), cbtx.Serialize(), nil)
		genesis := GenesisBlock(cbtx)
		genesis.Print()
		e = db.Put(append([]byte("b"), genesis.Header.Hash...), genesis.Serialize(), nil)
		util.ERR("Genesis Block put in DB Error", e)
		e = db.Put([]byte("l"), genesis.Header.Hash, nil)
		util.ERR("lastBlockHash put in DB Error", e)
		LastBlockHash = genesis.Header.Hash
	}

	utxo_db, e := leveldb.OpenFile(config.V.GetString("Default_db")+"_UTXO", nil)
	util.ERR("NewBlockchain open DB Error", e)

	txpool_db, e := leveldb.OpenFile(config.V.GetString("Default_db")+"_TxPool", nil)
	util.ERR("NewBlockchain open DB Error", e)

	bc := Blockchain{LastBlockHash, db, utxo_db, txpool_db}

	return &bc
}

func (bc *Blockchain) AddBlock(data string) *block.Block {
	lastHash, e := bc.DB.Get([]byte("l"), nil)
	if e != nil {
		fmt.Println("Blockchain not in DB")
		fmt.Println(lastHash)
	}
	lastBlockByte, e := bc.DB.Get(append([]byte("b"), lastHash...), nil)
	if e != nil {
		fmt.Println("lastHash exist, lastHash's block is not in DB")
		fmt.Println(lastBlockByte)
	}
	lastBlock := block.DeserializeBlock(lastBlockByte)

	cbtx := transaction.CoinbaseTx(config.V.GetString("Coinbase"), config.V.GetString("Coinbase"))
	bc.AddTx(cbtx)
	TXs := transaction.GetTxFromDB(bc.TxPoolDB)
	newBlock := POW.FindAnswer(data, lastHash, lastBlock.Header.Height+1, config.V.GetInt64("TargetBits"), []byte{}, TXs)
	newBlock.Print()

	e = bc.DB.Put(append([]byte("b"), newBlock.Header.Hash...), newBlock.Serialize(), nil)
	util.ERR("new block put in db Error", e)

	e = bc.DB.Put([]byte("l"), newBlock.Header.Hash, nil)
	util.ERR("new block hash put in db(l) Error", e)
	bc.LastBlockHash = newBlock.Header.Hash

	util.ClearDB(bc.TxPoolDB)

	return newBlock
}

func (bc *Blockchain) AddTx(Tx *transaction.Tx) {
	fmt.Println("Add Tx")
	Tx.Print()
	bc.TxPoolDB.Put(Tx.ID, Tx.Serialize(), nil)
}

func (bc *Blockchain) FindTransaction(ID []byte) (transaction.Tx, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.Header.PrevBlockHash) == 0 {
			break
		}
	}

	return transaction.Tx{}, errors.New("Tx Not Found")
}

func (bc *Blockchain) FindUTXO() map[string]transaction.TXOutputs {
	UTXO := make(map[string]transaction.TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.TxID)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.Header.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

func (bc *Blockchain) SignTransaction(tx *transaction.Tx, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]transaction.Tx)

	for _, vin := range tx.Vin {
		prevTX, e := bc.FindTransaction(vin.TxID)
		util.ERR("FindTransaction Error", e)

		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) VerifyTransaction(tx *transaction.Tx) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]transaction.Tx)

	for _, vin := range tx.Vin {
		prevTX, e := bc.FindTransaction(vin.TxID)
		util.ERR("FindTransaction Error", e)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

type BlockchainIterator struct {
	currentHash []byte
	DB          *leveldb.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.LastBlockHash, bc.DB}
	return bci
}

func (i *BlockchainIterator) Next() *block.Block {
	var Block *block.Block

	encodedBlock, e := i.DB.Get(append([]byte("b"), i.currentHash...), nil)
	util.ERR("read DB Error", e)
	Block = block.DeserializeBlock(encodedBlock)
	i.currentHash = Block.Header.PrevBlockHash
	return Block
}
