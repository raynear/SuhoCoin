package server

import (
	"SuhoCoin/blockchain"
	"SuhoCoin/config"
	"SuhoCoin/transaction"
	"SuhoCoin/util"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/urfave/negroni"
)

type Message struct {
	JsonRPC string
	ID      int
	Method  string
	Params  []interface{}
}

type Result struct {
	JsonRPC string
	ID      int
	Result  []byte
}

func JsonRPC(w http.ResponseWriter, r *http.Request, bc *blockchain.Blockchain) {
	b, e := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	util.ERR("RPC ReadBody Error", e)

	var msg Message
	e = json.Unmarshal(b, &msg)
	util.ERR("json unmarshal Error", e)

	result := Result{"2.0", msg.ID, []byte("")}

	switch msg.Method {
	case "GetBestHeight":
		height := bc.GetBestHeight()
		result.Result = []byte(strconv.FormatInt(height, 10))
	case "GetBestBlockHash":
		BestBlockHash := bc.GetBestBlockHash()
		fmt.Println("BestBlockHash:", BestBlockHash)
		result.Result = BestBlockHash
	case "GetBlock":
		hash, e := GetBytes(msg.Params[0])
		util.ERR("convert byte error", e)
		fmt.Println("hash", hash)
		aBlock := bc.GetBlockByHash(hash)
		fmt.Println("Block Height:", aBlock.Header.Height)
		result.Result = []byte(string(hash))
	default:
	}

	output, e := json.Marshal(result)
	if e != nil {
		http.Error(w, e.Error(), 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func StartServer(bc *blockchain.Blockchain, Port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		JsonRPC(w, r, bc)
	})

	n := negroni.New()
	n.UseHandler(mux)

	http.ListenAndServe(":"+Port, n)
}

type version struct {
	Version    int
	BestHeight int64
	AddrFrom   string
}

type getblocks struct {
	AddrFrom string
}

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var knownNodes = []string{"localhost:3000"}
var blocksInTransit = [][]byte{}
var mempool = make(map[string]transaction.Tx)

func _StartServer(bc *blockchain.Blockchain) {
	nodeAddress = fmt.Sprintf("localhost:%s", config.V.GetString("Port"))
	miningAddress = config.V.GetString("Coinbase")
	ln, e := net.Listen(protocol, nodeAddress)
	util.ERR("net listen error", e)
	defer ln.Close()

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, e := ln.Accept()
		util.ERR("listen accept Error", e)
		go handleConnection(conn, bc)
	}
}

func handleConnection(conn net.Conn, bc *blockchain.Blockchain) {
	request, e := ioutil.ReadAll(conn)
	util.ERR("connection read Error", e)
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func sendVersion(addr string, bc *blockchain.Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(version{nodeVersion, bestHeight, nodeAddress})

	request := append(commandToBytes("version"), payload...)
	sendData(addr, request)
}

func sendData(addr string, data []byte) {
	conn, e := net.Dial(protocol, addr)
	if e != nil {
		fmt.Println(addr, "is not available")
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}
		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, e = io.Copy(conn, bytes.NewReader(data))
	util.ERR("SendData Error", e)
}

func nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func handleVersion(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	e := dec.Decode(&payload)
	util.ERR("decode Error", e)

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return string(command)
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	e := enc.Encode(data)
	util.ERR("encode error", e)

	return buff.Bytes()
}
