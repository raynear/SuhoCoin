package main

import (
    "SuhoCoin/block"
    "SuhoCoin/blockchain"
    "fmt"
)

func main() {
    fmt.Println("testing blockchain")

    bc := blockchain.NewBlockchain()

    bc.AddBlock("test1")
    bc.AddBlock("test2")
    bc.AddBlock("running?")

    iter := bc.DB.NewIterator(nil, nil)
    for iter.Next() {
        key := iter.Key()
        value := iter.Value()
        if string(key) == "l" {
            fmt.Printf("Key: %s | Value: %x", string(key), value)
        } else {
            fmt.Printf("Key: %x | Value: ", key)
            fmt.Println(block.DeserializeBlock(value))
        }
    }
}
