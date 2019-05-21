package main

import (
    "SuhoCoin/blockchain"
    "fmt"
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
        fmt.Println("Nonce", block.Header.Nonce)
        fmt.Printf("Hash: %x\n", block.Header.Hash)
        fmt.Println()
    }
}
