package btree

import (
	"testing"
)

/*
 BTree for tests
                              [.100.]
                      /                      \
             [.35.65.]                           [.130.180.]
             /    |    \                        /    |      \
    [.10.20.] [.40.50.] [.70.80.90]  [.110.120.] [.140.160.] [.190.240.260.]
*/
func buildBTree() *BTree {
	minimumDegree := 3

	tree := NewBTree(minimumDegree)

	tree.Root = NewBTreeNode(minimumDegree, false)
	// [.100.]
	tree.Root.KeyCount = 1
	tree.Root.Keys[0] = 100
	tree.Root.Children[0] = NewBTreeNode(minimumDegree, false)
	tree.Root.Children[1] = NewBTreeNode(minimumDegree, false)

	// [.35.65.]
	tree.Root.Children[0].KeyCount = 2
	tree.Root.Children[0].Keys[0] = 35
	tree.Root.Children[0].Keys[1] = 65
	tree.Root.Children[0].Children[0] = NewBTreeNode(minimumDegree, true)
	tree.Root.Children[0].Children[1] = NewBTreeNode(minimumDegree, true)
	tree.Root.Children[0].Children[2] = NewBTreeNode(minimumDegree, true)

	// [.130.180.]
	tree.Root.Children[1].KeyCount = 2
	tree.Root.Children[1].Keys[0] = 130
	tree.Root.Children[1].Keys[1] = 180
	tree.Root.Children[1].Children[0] = NewBTreeNode(minimumDegree, true)
	tree.Root.Children[1].Children[1] = NewBTreeNode(minimumDegree, true)
	tree.Root.Children[1].Children[2] = NewBTreeNode(minimumDegree, true)

	// [.10.20.]
	tree.Root.Children[0].Children[0].KeyCount = 2
	tree.Root.Children[0].Children[0].Keys[0] = 10
	tree.Root.Children[0].Children[0].Keys[1] = 20

	// [.40.50.]
	tree.Root.Children[0].Children[1].KeyCount = 2
	tree.Root.Children[0].Children[1].Keys[0] = 40
	tree.Root.Children[0].Children[1].Keys[1] = 50

	// [.70.80.90]
	tree.Root.Children[0].Children[2].KeyCount = 3
	tree.Root.Children[0].Children[2].Keys[0] = 70
	tree.Root.Children[0].Children[2].Keys[1] = 80
	tree.Root.Children[0].Children[2].Keys[2] = 90

	// [.110.120.]
	tree.Root.Children[1].Children[0].KeyCount = 2
	tree.Root.Children[1].Children[0].Keys[0] = 110
	tree.Root.Children[1].Children[0].Keys[1] = 120

	// [.140.160.]
	tree.Root.Children[1].Children[1].KeyCount = 2
	tree.Root.Children[1].Children[1].Keys[0] = 140
	tree.Root.Children[1].Children[1].Keys[1] = 160

	// [.190.240.260.]
	tree.Root.Children[1].Children[2].KeyCount = 3
	tree.Root.Children[1].Children[2].Keys[0] = 190
	tree.Root.Children[1].Children[2].Keys[1] = 240
	tree.Root.Children[1].Children[2].Keys[2] = 260

	return tree
}

func TestTraverse(t *testing.T) {
	tree := buildBTree()

	expected := " 10 20 35 40 50 65 70 80 90 100 110 120 130 140 160 180 190 240 260"
	result := tree.Traverse()

	if result != expected {
		t.Errorf("result: %s\nexpected: %s", result, expected)
	}
}

func TestSearch(t *testing.T) {
	tree := NewBTree(3)

	tree.Insert(10, 1)
	tree.Insert(20, 2)
	tree.Insert(30, 3)
	tree.Insert(40, 4)
	tree.Insert(50, 5)
	tree.Insert(60, 6)
	tree.Insert(70, 7)
	tree.Insert(80, 8)
	tree.Insert(90, 9)

	key, offset := tree.Search(30)
	if key != 30 || offset != 3 {
		t.Errorf("record not found")
	}

	key, offset = tree.Search(70)
	if key != 70 || offset != 7 {
		t.Errorf("record not found")
	}
}

func TestInsert(t *testing.T) {
	tree := NewBTree(3)

	tree.Insert(10, 1)
	tree.Insert(20, 2)
	tree.Insert(30, 3)
	tree.Insert(40, 4)
	tree.Insert(50, 5)
	tree.Insert(60, 6)
	tree.Insert(70, 7)
	tree.Insert(80, 8)
	tree.Insert(90, 9)

	expected := " 10 20 30 40 50 60 70 80 90"
	result := tree.Traverse()

	if result != expected {
		t.Errorf("result: %s\nexpected: %s", result, expected)
	}
}
