package main

import (
    "fmt"
    "test/block/blockchain"
)

func main() {
    fmt.Println("testing blockchain")

    bc := blockchain.NewBlockchain()

    bc.AddBlock("test1")
    bc.AddBlock("test2")
    bc.AddBlock("running?")

    for _, block := range bc.Blocks {
        fmt.Printf("prev hash: %x\n", block.Header.PrevBlockHash)
        fmt.Println("Data", block.Data)
        fmt.Printf("Hash: %x\n", block.Header.Hash)
        fmt.Println()
    }
}
