package main

import (
	"SuhoCoin/cli"
	"SuhoCoin/config"
	"fmt"
)

func main() {
	fmt.Println("SuhoCoin Running")
	fmt.Println()

	config.V = config.ReadConfig("suho.conf")
	//	bc := blockchain.NewBlockchain("3001")
	//	defer bc.DB.Close()
	//	defer bc.UTXODB.Close()
	//	defer bc.TxPoolDB.Close()

	cli.Run()
}
