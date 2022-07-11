package btree

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type BTree struct {
	Root          *BTreeNode
	MinimumDegree int
}

// BTree initializer
func NewBTree(minimumDegree int) *BTree {
	return &BTree{Root: nil, MinimumDegree: minimumDegree}
}

func (b *BTree) Traverse() string {
	if b.Root == nil {
		return ""
	}

	return b.Root.traverse()
}

func (b *BTree) Search(key int) (int, int) {
	if b.Root == nil {
		return -1, -1
	}

	return b.Root.search(key)
}

func (b *BTree) Insert(key, record int) {
	if b.Root == nil {
		b.Root = NewBTreeNode(b.MinimumDegree, true)
		b.Root.KeyCount = 1
		b.Root.Keys[0] = key
		b.Root.Records[0] = record

		return
	}

	// Check if root is full
	if b.Root.KeyCount == 2*b.MinimumDegree-1 {
		newRoot := NewBTreeNode(b.MinimumDegree, false)
		newRoot.Children[0] = b.Root

		// Split the old root and move 1 key to the new root
		newRoot.splitChild(0, b.Root)

		// Decide which of the two children is going to have new key
		i := 0
		if newRoot.Keys[0] < key {
			i++
		}
		newRoot.Children[i].insertNonFull(key, record)

		// Change root
		b.Root = newRoot

		return
	}

	b.Root.insertNonFull(key, record)
}

type BTreeNode struct {
	Offset        int
	KeyCount      int
	MinimumDegree int
	Leaf          bool

	Keys     []int
	Records  []int
	Children []*BTreeNode
}

// BTreeNode initializer
func NewBTreeNode(MinimumDegree int, Leaf bool) *BTreeNode {
	newNode := &BTreeNode{KeyCount: 0, Leaf: Leaf}
	newNode.MinimumDegree = MinimumDegree

	newNode.Keys = make([]int, 2*MinimumDegree-1)
	newNode.Records = make([]int, 2*MinimumDegree-1)
	newNode.Children = make([]*BTreeNode, 2*MinimumDegree)

	return newNode
}

func (node *BTreeNode) ToBuffer() bytes.Buffer {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	err := enc.Encode(node)
	if err != nil {
		fmt.Println(node)
		panic(err)
	}

	return buff
}

func (node *BTreeNode) Save() {

}

func (node *BTreeNode) traverse() string {
	output := ""

	for i := 0; i < node.KeyCount; i++ {
		if !node.Leaf {
			// for each Children of node, when node isn't Leaf, should be called "traverse"
			output += node.Children[i].traverse()
		}

		// after Children's print, should print the key of parent of right
		output += fmt.Sprintf(" %d", node.Keys[i])
	}

	if !node.Leaf {
		// the last Children should be called too
		output += node.Children[node.KeyCount].traverse()
	}

	return output
}

func (node *BTreeNode) search(key int) (int, int) {
	// find the first key greater than or equal to "key"
	i := 0
	for i < node.KeyCount && key > node.Keys[i] {
		i++
	}

	// if match, return this node
	if node.Keys[i] == key {
		return node.Keys[i], node.Records[i]
	}

	// else then check if this node is Leaf
	if node.Leaf {
		// if true then record was not found
		return -1, -1
	}

	// else then call search to Children "i"
	return node.Children[i].search(key)
}

// The node must be non-full
func (node *BTreeNode) insertNonFull(key, record int) {
	// Initialize index as index of rightmost element
	idx := node.KeyCount - 1

	// If this node is leaf
	if node.Leaf {
		for idx >= 0 && node.Keys[idx] > key {
			node.Keys[idx+1] = node.Keys[idx]
			node.Records[idx+1] = node.Records[idx]

			idx--
		}

		node.Keys[idx+1] = key
		node.Records[idx+1] = record
		node.KeyCount += 1
	} else {
		for idx >= 0 && node.Keys[idx] > key {
			idx--
		}

		// Check if child node is full
		if node.Children[idx+1].KeyCount == 2*node.MinimumDegree-1 {
			node.splitChild(idx+1, node.Children[idx+1])

			// Decide which of the two children is going to have new key
			if node.Keys[idx+1] < key {
				idx++
			}
		}

		node.Children[idx+1].insertNonFull(key, record)
	}
}

func (node *BTreeNode) splitChild(idx int, child *BTreeNode) {
	newChild := NewBTreeNode(child.MinimumDegree, child.Leaf)
	newChild.KeyCount = node.MinimumDegree - 1

	for j := 0; j < node.MinimumDegree-1; j++ {
		newChild.Keys[j] = child.Keys[j+node.MinimumDegree]
		newChild.Records[j] = child.Records[j+node.MinimumDegree]
	}

	if !child.Leaf {
		for j := 0; j < node.MinimumDegree; j++ {
			newChild.Children[j] = child.Children[j+node.MinimumDegree]
		}
	}

	child.KeyCount = node.MinimumDegree - 1

	for j := node.MinimumDegree - 1; j >= idx; j-- {
		node.Children[j+1] = node.Children[j]
	}
	node.Children[idx+1] = newChild

	for j := node.KeyCount - 1; j >= idx; j-- {
		node.Keys[j+1] = node.Keys[j]
		node.Records[j+1] = node.Records[j]
	}

	node.Keys[idx] = child.Keys[node.MinimumDegree-1]
	node.Records[idx] = child.Records[node.MinimumDegree-1]
	node.KeyCount += 1
}
