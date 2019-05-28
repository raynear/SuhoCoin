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

    data := bytes.Join(
        [][]byte{
            pow.block.Header.PrevBlockHash,
            height,
            timestamp,
            difficulty,
            nonce,
            pow.block.Header.MerkleRoot,
        },
        []byte{},
    )

    return data
}

func CoinbaseTx(to string, data string) *transaction.Tx {
    if data == "" {
        data = fmt.Sprintf("Reward to '%s'", to)
    }
    txin := transaction.TXInput{[]byte{}, -1, data}
    txout := transaction.TXOutput{config.V.GetInt("Reward"), to}
    tx := transaction.Tx{nil, []transaction.TXInput{txin}, []transaction.TXOutput{txout}}
    return &tx
}

func (pow *POW) Run() (int64, []byte) {
    var hashInt big.Int
    var hash [32]byte
    var nonce int64

    nonce = 0

    fmt.Printf("Mining Block %s\n", pow.block.Data)

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
    block := &block.Block{blockheader.BlockHeader{config.V.GetInt64("BlockchainVersion"), []byte{}, prevBlockHash, height, time.Now().Unix(), difficulty, 0, merkleRoot}, 0, TXs, data}

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
