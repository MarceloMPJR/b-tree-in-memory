package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"

	btree "github.com/MarceloMPJR/b-tree-in-memory"
)

func main() {
	fd, err := syscall.Open("sample.tsv", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	tree := buildTreeByKey(fd, "tconst")

	searchKey(fd, 4, tree)
	searchKey(fd, 1, tree)
	searchKey(fd, 9, tree)
	searchKey(fd, 18, tree)
	searchKey(fd, 3, tree)

	if err = syscall.Close(fd); err != nil {
		panic(err)
	}
}

func searchKey(fd int, key int, tree *btree.BTree) {
	_, offset := tree.Search(key)
	if offset == -1 {
		fmt.Printf("[%d] => RECORD NOT FOUND!\n", key)
		return
	}

	syscall.Seek(fd, int64(offset), io.SeekStart)
	line, _, err := getLine(fd)

	if err != nil {
		panic(err)
	}

	fmt.Printf("[%d] => %v\n", key, line)
}

func buildTreeByKey(fd int, key string) *btree.BTree {
	keyIndex := -1
	tree := btree.NewBTree(5)

	head, eof, err := getLine(fd)
	if err != nil {
		panic(err)
	}

	if eof {
		return nil
	}

	for i, k := range head {
		if k == key {
			keyIndex = i
			break
		}
	}

	if keyIndex == -1 {
		panic("key not found")
	}

	var offset int64
	var line []string

	for {
		offset, err = syscall.Seek(fd, 0, os.SEEK_CUR)
		if err != nil {
			panic(err)
		}

		line, eof, err = getLine(fd)
		if err != nil {
			panic(err)
		}

		if eof {
			break
		}

		tree.Insert(keyToInt(line[keyIndex]), int(offset))
	}

	return tree
}

func getLine(fd int) ([]string, bool, error) {
	var line []string

	for {
		str, eof, newline, err := nextString(fd)
		if err != nil {
			return line, false, err
		}

		line = append(line, str)

		if eof {
			return line, true, nil
		}

		if newline {
			break
		}
	}

	return line, false, nil
}

func nextString(fd int) (string, bool, bool, error) {
	buff := make([]byte, 1)

	var str string
	var lastChar string

	for {
		n, err := syscall.Read(fd, buff)
		if err != nil {
			return "", false, false, err
		}

		if n == 0 {
			return str, true, false, nil
		}

		lastChar = string(buff)
		if lastChar == "\t" {
			return str, false, false, nil
		}

		if lastChar == "\n" {
			return str, false, true, nil
		}
		str += lastChar
	}
}

func keyToInt(key string) int {
	key = key[2:]
	keyInt, err := strconv.Atoi(key)

	if err != nil {
		panic(err)
	}

	return keyInt
}
