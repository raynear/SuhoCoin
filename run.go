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

    config.V = config.ReadConfig("suho.conf")
    bc := blockchain.NewBlockchain()
    defer bc.DB.Close()

    cli.Run(bc)
}
