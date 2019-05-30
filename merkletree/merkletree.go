package merkletree

import (
    "crypto/sha256"
)

type MerkleTree struct {
    Root *MerkleNode
}
type MerkleNode struct {
    Left  *MerkleNode
    Right *MerkleNode
    Data  []byte
}

func NewMerkleNode(left *MerkleNode, right *MerkleNode, data []byte) *MerkleNode {
    mNode := MerkleNode{}
    if left == nil && right == nil {
        hash := sha256.Sum256(data)
        mNode.Data = hash[:]
    } else {
        prevHashes := append(left.Data, right.Data...)
        hash := sha256.Sum256(prevHashes)
        mNode.Data = hash[:]
    }
    mNode.Left = left
    mNode.Right = right
    return &mNode
}
func NewMerkleTree(data [][]byte) *MerkleTree {
    var nodes []MerkleNode
    if len(data)%2 != 0 {
        data = append(data, data[len(data)-1])
    }
    for _, aData := range data {
        node := NewMerkleNode(nil, nil, aData)
        nodes = append(nodes, *node)
    }
    bi := 1
    for len(nodes) > bi<<1 {
        bi = bi << 1
    }
    var newLevel []MerkleNode
    k := 0
    for k = 0; bi != len(nodes)-k/2; k += 2 {
        node := NewMerkleNode(&nodes[k], &nodes[k+1], nil)
        newLevel = append(newLevel, *node)
    }
    nodes = nodes[k:]
    nodes = append(nodes, newLevel...)
    for i := 0; i < len(nodes)/2; i++ {
        var newLevel []MerkleNode
        for j := 0; j < len(nodes); j += 2 {
            node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
            newLevel = append(newLevel, *node)
        }
        nodes = newLevel
    }
    mMerkleTree := MerkleTree{&nodes[0]}
    return &mMerkleTree
}
