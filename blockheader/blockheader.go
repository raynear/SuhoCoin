package blockheader

import (
    "bytes"
    "crypto/sha256"
    "strconv"
)

type BlockHeader struct {
    Version       int64
    Hash          []byte
    PrevBlockHash []byte
    Height        int64
    TimeStamp     int64
    Difficulty    int64
    Nonce         []byte
    MerkleRoot    []byte
}

func (bh *BlockHeader) SetHash() {
    version := []byte(strconv.FormatInt(bh.Version, 10))
    height := []byte(strconv.FormatInt(bh.Height, 10))
    difficulty := []byte(strconv.FormatInt(bh.Difficulty, 10))
    timestamp := []byte(strconv.FormatInt(bh.TimeStamp, 10))
    Serialize := bytes.Join([][]byte{version,
        bh.PrevBlockHash,
        height,
        timestamp,
        difficulty,
        bh.Nonce,
        bh.MerkleRoot},
        []byte{})
    hash := sha256.Sum256(Serialize)

    bh.Hash = hash[:]
}
