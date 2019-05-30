package cli

import (
    "SuhoCoin/Consensus/POW"
    "SuhoCoin/block"
    "SuhoCoin/blockchain"
    "SuhoCoin/config"
    "SuhoCoin/transaction"
    "SuhoCoin/util"
    "SuhoCoin/wallet"
    "bytes"
    "fmt"
    "os"
    "strconv"

    "github.com/urfave/cli"
)

func Run(bc *blockchain.Blockchain) {
    app := cli.NewApp()

    app.Commands = []cli.Command{
        {
            Name:  "createwallet",
            Usage: "createwallet for use",
            Action: func(c *cli.Context) error {
                myWallet := wallet.NewWallet()
                fmt.Println("Your new address:", string(myWallet.GetAddress()))

                return nil
            },
        },
        {
            Name:  "getbalance",
            Usage: "getbalance 'address'",
            Action: func(c *cli.Context) error {
                address := c.Args()[0]
                balance := 0
                UTXOs := bc.FindUTXO(address)

                for _, out := range UTXOs {
                    balance += out.Value
                }

                fmt.Printf("Balance of '%s': %d\n", address, balance)

                return nil
            },
        },
        {
            Name:    "addblock",
            Aliases: []string{"add"},
            Usage:   "addblock 'a send 1 to b'",
            Action: func(c *cli.Context) error {
                fmt.Println("addblock:", c.Args())
                bc.AddBlock(c.Args().First())

                return nil
            },
        },
        {
            Name:    "transaction",
            Aliases: []string{"tx"},
            Usage:   "transaction 'TXHASH'",
            Action: func(c *cli.Context) error {
                fmt.Println("tx detail:", c.Args())
                aTxPayload, e := bc.DB.Get([]byte(c.Args()[0]), nil)
                err.ERR("Get Tx Error", e)
                aTx := transaction.DeserializeTx(aTxPayload)
                fmt.Println("tx:", aTx)

                return nil
            },
        },
        {
            Name:    "printchain",
            Aliases: []string{"pc"},
            Usage:   "print all chain",
            Action: func(c *cli.Context) error {
                bci := bc.Iterator()

                for {
                    block := bci.Next()
                    fmt.Printf("Prev Hash: %x\n", block.Header.PrevBlockHash)
                    fmt.Printf("Data: %s\n", block.Data)
                    fmt.Printf("Hash: %x\n", block.Header.Hash)

                    pow := POW.NewPOW(block)
                    fmt.Printf("POW: %s\n\n", strconv.FormatBool(pow.Validate()))

                    if len(block.Header.PrevBlockHash) == 0 {
                        break
                    }
                }

                return nil
            },
        },
        {
            Name:    "printdb",
            Aliases: []string{"pd"},
            Usage:   "print all chain on db",
            Action: func(c *cli.Context) error {
                iter := bc.DB.NewIterator(nil, nil)
                for iter.Next() {
                    key := iter.Key()
                    value := iter.Value()
                    if string(key) == "l" {
                        fmt.Printf("Key: %s | Value: %x\n", string(key), value)
                    } else {
                        if bytes.Compare(key[:1], []byte("b")) == 0 {
                            fmt.Printf("Key: %x | Value: ", key)
                            aBlock := block.DeserializeBlock(value)
                            aBlock.Print()
                        }
                        if bytes.Compare(key[:1], []byte("t")) == 0 {
                            fmt.Printf("Key: %x | Value: ", key)
                            aTx := transaction.DeserializeTx(value)
                            aTx.Print()
                        }
                    }
                }

                return nil
            },
        },
        {
            Name:  "clearDB",
            Usage: "delete all DB",
            Action: func(c *cli.Context) error {
                bc.DB.Close()
                e := os.RemoveAll("./" + config.V.GetString("Default_db"))
                err.ERR("del error", e)
                return nil
            },
        },
    }

    e := app.Run(os.Args)
    err.ERR("Cli Error", e)
}
