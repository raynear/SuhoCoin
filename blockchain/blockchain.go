package blockchain

import (
    "SuhoCoin/Consensus/POW"
    "SuhoCoin/block"
    "SuhoCoin/config"
    "SuhoCoin/transaction"
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
    db, err := leveldb.OpenFile(config.V.GetString("Default_db"), nil)
    if err != nil {
        fmt.Println("NewBlockchain open DB Error", err)
    }

    LastBlockHash, err = db.Get([]byte("l"), nil)
    if err != nil {
        cbtx := POW.CoinbaseTx(address, "GENESIS of SuhoCoin")
        genesis := GenesisBlock(cbtx)
        fmt.Println(genesis)
        err = db.Put(genesis.Header.Hash, genesis.Serialize(), nil)
        if err != nil {
            fmt.Println("Genesis Block put in DB Error")
        }
        err = db.Put([]byte("l"), genesis.Header.Hash, nil)
        if err != nil {
            fmt.Println("lastBlockHash put in DB Error")
        }
        LastBlockHash = genesis.Header.Hash

        lvalue, _ := db.Get([]byte("l"), nil)
        fmt.Println("DB(l)", lvalue)
    }

    bc := Blockchain{LastBlockHash, db}

    return &bc
}

func (bc *Blockchain) AddBlock(data string) {
    lastHash, err := bc.DB.Get([]byte("l"), nil)
    if err != nil {
        fmt.Println("Blockchain not in DB")
        fmt.Println(lastHash)
    }
    lastBlockByte, err := bc.DB.Get(lastHash, nil)
    if err != nil {
        fmt.Println("lastHash exist, lastHash's block is not in DB")
        fmt.Println(lastBlockByte)
    }
    lastBlock := block.DeserializeBlock(lastBlockByte)

    cbtx := POW.CoinbaseTx(config.V.GetString("Coinbase"), "")
    newBlock := POW.FindAnswer(data, lastHash, lastBlock.Header.Height, 0, []byte{}, []*transaction.Tx{cbtx})

    err = bc.DB.Put(newBlock.Header.Hash, newBlock.Serialize(), nil)
    if err != nil {
        fmt.Println("new block put in db Error")
    }

    err = bc.DB.Put([]byte("l"), newBlock.Header.Hash, nil)
    if err != nil {
        fmt.Println("new block hash put in db(l) Error")
    }
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

    encodedBlock, err := i.db.Get(i.currentHash, nil)
    if err != nil {
        fmt.Println("read DB Error:", err)
    }
    Block = block.DeserializeBlock(encodedBlock)
    i.currentHash = Block.Header.PrevBlockHash
    return Block
}
