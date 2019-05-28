package main

import (
    "SuhoCoin/blockchain"
    "SuhoCoin/cli"
    "SuhoCoin/config"
    "fmt"
)

func main() {
    fmt.Println("SuhoCoin Running")
    fmt.Println()

    config.V = config.ReadConfig("test.config")
    bc := blockchain.NewBlockchain("raynear")
    defer bc.DB.Close()

    cli.Run(bc)
}
