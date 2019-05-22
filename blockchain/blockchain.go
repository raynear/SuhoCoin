package blockchain

import (
    "SuhoCoin/Consensus/POW"
    "SuhoCoin/block"
    "SuhoCoin/config"
    "fmt"

    "github.com/syndtr/goleveldb/leveldb"
)

type Blockchain struct {
    LastBlockHash []byte
    DB            *leveldb.DB
}

func GenesisBlock() *block.Block {
    return POW.FindAnswer("GenesisBlock", []byte{}, int64(0), int64(0), []byte{})
}

func NewBlockchain() *Blockchain {
    var LastBlockHash []byte
    db, err := leveldb.OpenFile(config.Default_db, nil)
    if err != nil {
        fmt.Println("NewBlockchain open DB Error", err)
    }

    LastBlockHash, err = db.Get([]byte("l"), nil)
    if err != nil {
        genesis := GenesisBlock()
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

        fmt.Println(db.Get([]byte("l"), nil))
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

    newBlock := POW.FindAnswer(data, lastHash, lastBlock.Header.Height, 0, []byte{})

    err = bc.DB.Put(newBlock.Header.Hash, newBlock.Serialize(), nil)
    if err != nil {
        fmt.Println("new block put in db Error")
    }
}
