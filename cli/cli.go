package cli

import (
	"SuhoCoin/Consensus/POW"
	utxo "SuhoCoin/UTXO"
	"SuhoCoin/block"
	"SuhoCoin/blockchain"
	"SuhoCoin/config"
	"SuhoCoin/transaction"
	"SuhoCoin/util"
	"SuhoCoin/wallet"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

func Run(bc *blockchain.Blockchain) {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:  "createwallet",
			Usage: "createwallet for use",
			Action: func(c *cli.Context) error {
				myAddress := wallet.NewWallet()
				fmt.Println("Your new address:", myAddress)

				return nil
			},
		},
		{
			Name:  "listaddress",
			Usage: "listaddress",
			Action: func(c *cli.Context) error {
				files, err := ioutil.ReadDir("./")
				if err != nil {
					log.Fatal(err)
				}

				for _, f := range files {
					if strings.Contains(f.Name(), ".wallet") {
						address := f.Name()
						a := []rune(address)

						myAddress := string(a[:len(a)-7])

						fmt.Println("address :", myAddress)
					}
				}

				return nil
			},
		},
		{
			Name: "address2hash",
			Action: func(c *cli.Context) error {
				address := c.Args()[0]
				PubKeyHash := wallet.GetPubKeyHashFromAddress(address)
				fmt.Println("PubKeyHash", base58.Encode(PubKeyHash))
				return nil
			},
		},
		{
			Name:  "send",
			Usage: "send 'sender address' 'receiver address' amount",
			Action: func(c *cli.Context) error {
				sender := c.Args()[0]
				receiver := c.Args()[1]
				amount, e := strconv.Atoi(c.Args()[2])

				UTXO := utxo.UTXO{Blockchain: bc}

				EncryptedWallet := wallet.LoadFromFile(sender)

				address := []byte(sender)
				fmt.Println("address : ", string(address[:]))
				pubKeyHash := base58.Decode(string(address[:]))
				pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
				fmt.Println("pubKeyHash : ", base58.Encode(pubKeyHash[:]))

				fmt.Printf("Input Password: ")
				silentPassword, e := gopass.GetPasswdMasked()
				util.ERR("Password Input Error", e)

				myWallet := wallet.DecryptWallet(EncryptedWallet, string(silentPassword))

				hashedpubkey := wallet.HashPubKey(myWallet.PublicKey)

				fmt.Println("HashedPubKey", base58.Encode(hashedpubkey))

				tx := utxo.NewUTXOTransaction(&myWallet, receiver, amount, &UTXO)
				bc.AddTx(tx)

				return nil
			},
		},
		{
			Name:  "setcoinbase",
			Usage: "setcoinbase 'address'",
			Action: func(c *cli.Context) error {
				/////////////
				// NOT WORK
				/////////////
				address := c.Args()[0]
				fmt.Println("address:", address)

				config.V.Set("Coinbase", address)
				config.V.WriteConfigAs("suho.conf")

				fmt.Printf("NewCoinbase is %s\n", address)

				return nil
			},
		},
		{
			Name:  "getbalance",
			Usage: "getbalance 'address'",
			Action: func(c *cli.Context) error {
				address := c.Args()[0]
				pubKeyHash := base58.Decode(address)
				pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
				balance := 0
				UTXO := utxo.UTXO{Blockchain: bc}
				UTXOs := UTXO.FindUTXO(pubKeyHash)

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

				UTXO := utxo.UTXO{Blockchain: bc}
				UTXO.Reindex()

				return nil
			},
		},
		{
			Name:  "pendingtx",
			Usage: "Show pending tx",
			Action: func(c *cli.Context) error {
				iter := bc.TxPoolDB.NewIterator(nil, nil)
				for iter.Next() {
					value := iter.Value()
					aTx := transaction.DeserializeTx(value)
					aTx.Print()
				}

				return nil
			},
		},
		{
			Name:    "reindex",
			Aliases: []string{"re"},
			Usage:   "reindex utxo db",
			Action: func(c *cli.Context) error {
				UTXO := utxo.UTXO{Blockchain: bc}
				// Reindex => update로 바꿔야 함.
				UTXO.Reindex()

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
				util.ERR("Get Tx Error", e)
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
			Name:  "printutxo",
			Usage: "print all utxo on db",
			Action: func(c *cli.Context) error {
				iter := bc.UTXODB.NewIterator(nil, nil)
				for iter.Next() {
					key := iter.Key()
					value := iter.Value()
					fmt.Printf("Key: %x | Value: ", key)
					TXOs := transaction.DeserializeOutputs(value)
					for _, aOut := range TXOs.Outputs {
						aOut.Print()
					}
				}

				return nil
			},
		},
		{
			Name:  "clearDB",
			Usage: "delete all DB",
			Action: func(c *cli.Context) error {
				DBIter := bc.DB.NewIterator(nil, nil)
				for DBIter.Next() {
					key := DBIter.Key()
					e := bc.DB.Delete(key, nil)
					util.ERR("DB Delete Error", e)
				}
				DBIter.Release()
				DBIter = bc.UTXODB.NewIterator(nil, nil)
				for DBIter.Next() {
					key := DBIter.Key()
					e := bc.UTXODB.Delete(key, nil)
					util.ERR("UTXODB Delete Error", e)
				}
				DBIter.Release()
				DBIter = bc.TxPoolDB.NewIterator(nil, nil)
				for DBIter.Next() {
					key := DBIter.Key()
					e := bc.TxPoolDB.Delete(key, nil)
					util.ERR("TxPoolDB Delete Error", e)
				}
				DBIter.Release()

				return nil
			},
		},
	}

	e := app.Run(os.Args)
	util.ERR("Cli Error", e)
}
