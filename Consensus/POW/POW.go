package POW

import (
    "SuhoCoin/block"
    "SuhoCoin/config"
    "bytes"
    "crypto/sha256"
    "fmt"
    "math"
    "math/big"
    "strconv"
)

type POW struct {
    block  *block.Block
    target *big.Int
}

func NewPOW(b *block.Block) *POW {
    target := big.NewInt(1)
    target.Lsh(target, uint(256-config.TargetBits))

    pow := &POW{b, target}

    return pow
}

func (pow *POW) prepareData(_nonce int64) []byte {
    height := []byte(strconv.FormatInt(pow.block.Header.Height, 10))
    timestamp := []byte(strconv.FormatInt(pow.block.Header.TimeStamp, 10))
    difficulty := []byte(strconv.FormatInt(config.TargetBits, 10))
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

func (pow *POW) Validate() bool {
    var hashInt big.Int

    data := pow.prepareData(pow.block.Header.Nonce)
    hash := sha256.Sum256(data)
    hashInt.SetBytes(hash[:])

    isValid := hashInt.Cmp(pow.target) == -1

    return isValid
}
