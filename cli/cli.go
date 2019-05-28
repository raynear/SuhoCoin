package cli

import (
    "SuhoCoin/Consensus/POW"
    "SuhoCoin/block"
    "SuhoCoin/blockchain"
    "SuhoCoin/config"
    "fmt"
    "log"
    "os"
    "strconv"

    "github.com/spf13/viper"
    "github.com/urfave/cli"
)

func err(msg string, e error) {
    if e != nil {
        fmt.Println(msg, e)
        log.Fatal(msg, e)
    }
}

func Run(bc *blockchain.Blockchain) {
    file, e := os.Open("test.config")
    err("File Error", e)

    viper.SetConfigType("prop")
    viper.ReadConfig(file)

    app := cli.NewApp()

    app.Commands = []cli.Command{
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
                        fmt.Printf("Key: %s | Value: %x", string(key), value)
                    } else {
                        fmt.Printf("Key: %x | Value: ", key)
                        fmt.Println(block.DeserializeBlock(value))
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
                err := os.RemoveAll("./" + config.V.GetString("Default_db"))
                if err != nil {
                    fmt.Println("del error", err)
                }
                return nil
            },
        },
    }

    e = app.Run(os.Args)
    err("Cli Error", e)
}
