package blockheader

type BlockHeader struct {
	Version       int64
	Hash          []byte
	PrevBlockHash []byte
	Height        int64
	TimeStamp     int64
	Difficulty    int64
	Nonce         int64
	MerkleRoot    []byte
}
