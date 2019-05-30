package POW

import (
    "SuhoCoin/block"
    "SuhoCoin/blockheader"
    "SuhoCoin/config"
    "SuhoCoin/transaction"
    "bytes"
    "crypto/sha256"
    "fmt"
    "math"
    "math/big"
    "strconv"
    "time"
)

type POW struct {
    block  *block.Block
    target *big.Int
}

func NewPOW(b *block.Block) *POW {
    target := big.NewInt(1)
    target.Lsh(target, uint(256-config.V.GetInt("TargetBits")))

    pow := &POW{b, target}

    return pow
}

func (pow *POW) prepareData(_nonce int64) []byte {
    height := []byte(strconv.FormatInt(pow.block.Header.Height, 10))
    timestamp := []byte(strconv.FormatInt(pow.block.Header.TimeStamp, 10))
    difficulty := []byte(config.V.GetString("TargetBits"))
    nonce := []byte(strconv.FormatInt(_nonce, 10))
    merkleroot := pow.block.Header.MerkleRoot

    data := bytes.Join(
        [][]byte{
            pow.block.Header.PrevBlockHash,
            height,
            timestamp,
            difficulty,
            nonce,
            merkleroot,
        },
        []byte{},
    )

    return data
}

func (pow *POW) Run() (int64, []byte) {
    var hashInt big.Int
    var hash [32]byte
    var nonce int64

    nonce = 0

    fmt.Printf("Mining Block : %s\n", pow.block.Data)
    pow.block.Header.MerkleRoot = pow.block.NewTxMerkleTree()

    for nonce < math.MaxInt64 {
        data := pow.prepareData(nonce)
        hash = sha256.Sum256(data)
        //        fmt.Printf("%x\n", hash)

        hashInt.SetBytes(hash[:])

        if hashInt.Cmp(pow.target) == -1 {
            break
        } else {
            nonce++
        }
    }
    fmt.Println()

    return nonce, hash[:]
}

func FindAnswer(data string, prevBlockHash []byte, height int64, difficulty int64, merkleRoot []byte, TXs []*transaction.Tx) *block.Block {
    block := &block.Block{Header: blockheader.BlockHeader{Version: config.V.GetInt64("BlockchainVersion"), Hash: []byte{}, PrevBlockHash: prevBlockHash, Height: height, TimeStamp: time.Now().Unix(), Difficulty: difficulty, Nonce: 0, MerkleRoot: merkleRoot}, TxCnt: 0, Transactions: TXs, Data: data}

    pow := NewPOW(block)
    nonce, hash := pow.Run()

    block.Header.Hash = hash[:]
    block.Header.Nonce = nonce

    return block
}

func (pow *POW) Validate() bool {
    var hashInt big.Int

    data := pow.prepareData(pow.block.Header.Nonce)
    hash := sha256.Sum256(data)
    hashInt.SetBytes(hash[:])

    isValid := hashInt.Cmp(pow.target) == -1

    return isValid
}
