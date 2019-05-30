package blockchain

import (
    "SuhoCoin/Consensus/POW"
    "SuhoCoin/block"
    "SuhoCoin/config"
    "SuhoCoin/transaction"
    "SuhoCoin/util"
    "encoding/hex"
    "fmt"

    "github.com/syndtr/goleveldb/leveldb"
)

type Blockchain struct {
    LastBlockHash []byte
    DB            *leveldb.DB
}

func GenesisBlock(coinbase *transaction.Tx) *block.Block {
    return POW.FindAnswer("GenesisBlock", []byte{}, int64(0), int64(0), []byte{}, []*transaction.Tx{coinbase})
}

func NewBlockchain() *Blockchain {
    var LastBlockHash []byte
    db, e := leveldb.OpenFile(config.V.GetString("Default_db"), nil)
    err.ERR("NewBlockchain open DB Error", e)

    LastBlockHash, e = db.Get([]byte("l"), nil)
    if e != nil {
        cbtx := transaction.CoinbaseTx([]byte(config.V.GetString("Coinbase")), "GENESIS of SuhoCoin")
        db.Put(append([]byte("t"), cbtx.ID...), cbtx.Serialize(), nil)
        genesis := GenesisBlock(cbtx)
        genesis.Print()
        e = db.Put(append([]byte("b"), genesis.Header.Hash...), genesis.Serialize(), nil)
        err.ERR("Genesis Block put in DB Error", e)
        e = db.Put([]byte("l"), genesis.Header.Hash, nil)
        err.ERR("lastBlockHash put in DB Error", e)
        LastBlockHash = genesis.Header.Hash
    }

    bc := Blockchain{LastBlockHash, db}

    return &bc
}

func (bc *Blockchain) AddBlock(data string) {
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

    cbtx := transaction.CoinbaseTx([]byte(config.V.GetString("Coinbase")), "")
    bc.DB.Put(append([]byte("t"), cbtx.ID...), cbtx.Serialize(), nil)
    newBlock := POW.FindAnswer(data, lastHash, lastBlock.Header.Height+1, 0, []byte{}, []*transaction.Tx{cbtx})
    newBlock.Print()

    e = bc.DB.Put(append([]byte("b"), newBlock.Header.Hash...), newBlock.Serialize(), nil)
    err.ERR("new block put in db Error", e)

    e = bc.DB.Put([]byte("l"), newBlock.Header.Hash, nil)
    err.ERR("new block hash put in db(l) Error", e)
    bc.LastBlockHash = newBlock.Header.Hash
}

func (bc *Blockchain) FindUnspentTransactions(address string) []transaction.Tx {
    var UTXO []transaction.Tx
    STXO := make(map[string][]int)
    bci := bc.Iterator()

    //    fmt.Println(bci.currentHash)

    //    ttt, e := bci.DB.Get(bci.currentHash, nil)
    //    err.ERR("er?", e)

    //    newBlock := block.DeserializeBlock(ttt)

    //    fmt.Println(newBlock.Header.PrevBlockHash)

    for {
        block := bci.Next()

        if block.Header.PrevBlockHash == nil {
            return UTXO
        }

        for _, tx := range block.Transactions {
            txid := hex.EncodeToString(tx.ID)

        Outputs:
            for outIdx, out := range tx.Vout {
                if STXO[txid] != nil {
                    for _, spentOut := range STXO[txid] {
                        if spentOut == outIdx {
                            continue Outputs
                        }
                    }
                }
                if out.CanBeUnlockedWith(address) {
                    UTXO = append(UTXO, *tx)
                }
            }

            if tx.IsCoinbase() == false {
                for _, in := range tx.Vin {
                    if in.CanUnlockOutputWith(address) {
                        inTxID := hex.EncodeToString(in.Txid)
                        STXO[inTxID] = append(STXO[inTxID], in.Vout)
                    }
                }
            }
        }
    }
}

func (bc *Blockchain) FindUTXO(address string) []transaction.TXOutput {
    var UTXO []transaction.TXOutput

    unspentTransactions := bc.FindUnspentTransactions(address)
    fmt.Println("utx cnt :", len(unspentTransactions))
    for _, tx := range unspentTransactions {
        for _, out := range tx.Vout {
            if out.CanBeUnlockedWith(address) {
                UTXO = append(UTXO, out)
            }
        }
    }

    return UTXO
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
    err.ERR("read DB Error", e)
    Block = block.DeserializeBlock(encodedBlock)
    i.currentHash = Block.Header.PrevBlockHash
    return Block
}
