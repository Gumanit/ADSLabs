package main

import (
	"fmt"
	"math/rand/v2"
)

// Узел бинарного дерева.
type Node struct {
	Key         int
	Left, Right *Node
}

// Дерево хранит указатель на корень.
type Tree struct {
	root *Node
}

func (t *Tree) Build(keys []int) {
	t.root = build(keys)
}

func build(keys []int) *Node {
	if len(keys) == 0 {
		return nil
	}
	nl := len(keys) / 2
	return &Node{
		Key:   keys[0],
		Left:  build(keys[1 : 1+nl]),
		Right: build(keys[1+nl:]),
	}
}

func (t *Tree) Print() {
	printNode(t.root, 0)
}

func printNode(n *Node, level int) {
	if n == nil {
		return
	}
	printNode(n.Right, level+1)
	for i := 0; i < level; i++ {
		fmt.Print("   ")
	}
	fmt.Println(n.Key)
	printNode(n.Left, level+1)
}

func (t *Tree) NLR() {
	nlr(t.root)
}

func nlr(n *Node) {
	if n == nil {
		return
	}
	fmt.Print(n.Key, " ")
	nlr(n.Left)
	nlr(n.Right)
}

func generateKeys(n int) []int {
	used := make(map[int]bool, n)
	keys := make([]int, 0, n)
	for len(keys) < n {
		k := rand.IntN(900) + 100 // от 100 до 999
		if !used[k] {
			used[k] = true
			keys = append(keys, k)
		}
	}
	return keys
}

func main() {
	const n = 26

	keys := generateKeys(n)
	fmt.Println("Сгенерированные ключи:")
	for _, k := range keys {
		fmt.Printf("%d ", k)
	}
	fmt.Println()

	var tree Tree
	tree.Build(keys)

	fmt.Println("\nПостроенное идеально сбалансированное дерево:")
	tree.Print()

	fmt.Println("\nПрямой обход (NLR):")
	tree.NLR()
	fmt.Println()
}
