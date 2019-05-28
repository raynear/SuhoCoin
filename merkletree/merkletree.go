package merkletree

import (
    "crypto/sha256"
    "fmt"
)

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

func NewMerkleTree(data [][]byte) *MerkleNode {
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

    fmt.Println("bi", bi)
    fmt.Println("len(nodes)", len(nodes))
    //    fmt.Println("nodes", nodes)

    var newLevel []MerkleNode
    k := 0
    for k = 0; bi != len(nodes)-k/2; k += 2 {
        node := NewMerkleNode(&nodes[k], &nodes[k+1], nil)
        newLevel = append(newLevel, *node)
    }
    nodes = nodes[k:]
    nodes = append(nodes, newLevel...)

    fmt.Println("nodes", nodes, "len", len(nodes))
    fmt.Println("newLevel", newLevel, "len", len(newLevel))

    for i := 0; i < len(nodes)/2; i++ {
        var newLevel []MerkleNode

        for j := 0; j < len(nodes)-1; j += 2 {
            node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
            newLevel = append(newLevel, *node)
            fmt.Println(j, "newLevel", newLevel)
        }
        nodes = newLevel

        fmt.Println("newLevel", newLevel)
    }

    fmt.Println("Root", nodes[0])

    MerkleRoot := nodes[0]

    return &MerkleRoot
}
