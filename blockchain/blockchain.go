package blockchain

import (
    "SuhoCoin/Consensus/POW"
    "SuhoCoin/block"
    "SuhoCoin/config"
    "SuhoCoin/transaction"
    "SuhoCoin/util"
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

func NewBlockchain(address string) *Blockchain {
    var LastBlockHash []byte
    db, e := leveldb.OpenFile(config.V.GetString("Default_db"), nil)
    err.ERR("NewBlockchain open DB Error", e)

    LastBlockHash, e = db.Get([]byte("l"), nil)
    if e != nil {
        cbtx := POW.CoinbaseTx(address, "GENESIS of SuhoCoin")
        genesis := GenesisBlock(cbtx)
        fmt.Println(genesis)
        e = db.Put(genesis.Header.Hash, genesis.Serialize(), nil)
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
    lastBlockByte, e := bc.DB.Get(lastHash, nil)
    if e != nil {
        fmt.Println("lastHash exist, lastHash's block is not in DB")
        fmt.Println(lastBlockByte)
    }
    lastBlock := block.DeserializeBlock(lastBlockByte)

    cbtx := POW.CoinbaseTx(config.V.GetString("Coinbase"), "")
    newBlock := POW.FindAnswer(data, lastHash, lastBlock.Header.Height, 0, []byte{}, []*transaction.Tx{cbtx})

    e = bc.DB.Put(newBlock.Header.Hash, newBlock.Serialize(), nil)
    err.ERR("new block put in db Error", e)

    e = bc.DB.Put([]byte("l"), newBlock.Header.Hash, nil)
    err.ERR("new block hash put in db(l) Error", e)
    bc.LastBlockHash = newBlock.Header.Hash

}

type BlockchainIterator struct {
    currentHash []byte
    db          *leveldb.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
    bci := &BlockchainIterator{bc.LastBlockHash, bc.DB}
    return bci
}

func (i *BlockchainIterator) Next() *block.Block {
    var Block *block.Block

    encodedBlock, e := i.db.Get(i.currentHash, nil)
    err.ERR("read DB Error", e)
    Block = block.DeserializeBlock(encodedBlock)
    i.currentHash = Block.Header.PrevBlockHash
    return Block
}
