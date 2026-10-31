package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
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

// Build строит идеально сбалансированное дерево.
// Первый элемент становится корнем, следующие n/2 элементов идут
// в левое поддерево, остальные — в правое.
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

// Print выводит дерево сверху вниз: корень сверху, потомки ниже.
func (t *Tree) Print() {
	if t.root == nil {
		return
	}

	const keyLen = 3 // все ключи трёхзначные
	const gap = 2    // отступ между соседними поддеревьями

	// Ширина блока, который занимает поддерево вершины.
	width := make(map[*Node]int)
	var calcWidth func(n *Node) int
	calcWidth = func(n *Node) int {
		if n == nil {
			return 0
		}
		if n.Left == nil && n.Right == nil {
			width[n] = keyLen
			return keyLen
		}
		wl, wr := calcWidth(n.Left), calcWidth(n.Right)
		w := keyLen
		switch {
		case n.Left != nil && n.Right != nil:
			w = wl + gap + wr
		case n.Left != nil:
			w = wl + gap
		default:
			w = wr + gap
		}
		width[n] = w
		return w
	}
	calcWidth(t.root)

	type point struct {
		col   int // столбец начала ключа
		depth int
	}
	pos := make(map[*Node]point)
	height := 0

	// Размещаем вершины: левое поддерево — слева в блоке,
	// правое — справа, сама вершина — по центру блока.
	var place func(n *Node, start, depth int)
	place = func(n *Node, start, depth int) {
		if n == nil {
			return
		}
		w := width[n]
		pos[n] = point{start + (w-keyLen)/2, depth}
		if depth+1 > height {
			height = depth + 1
		}
		if n.Left != nil {
			place(n.Left, start, depth+1)
		}
		if n.Right != nil {
			place(n.Right, start+w-width[n.Right], depth+1)
		}
	}
	place(t.root, 0, 0)

	canvas := make([][]rune, 2*height-1)
	for i := range canvas {
		canvas[i] = []rune(strings.Repeat(" ", width[t.root]))
	}

	label := func(n *Node) string { return strconv.Itoa(n.Key) }
	center := func(n *Node) int { return pos[n].col + keyLen/2 }

	// Ключи вершин.
	for n, p := range pos {
		for i, r := range label(n) {
			canvas[2*p.depth][p.col+i] = r
		}
	}

	// Линии, соединяющие родителя с потомками.
	var connect func(n *Node)
	connect = func(n *Node) {
		if n == nil {
			return
		}
		row := 2*pos[n].depth + 1
		if row < len(canvas) {
			pc := center(n)
			line := func(from, to int) {
				if from > to {
					from, to = to, from
				}
				for c := from; c <= to; c++ {
					canvas[row][c] = '-'
				}
				canvas[row][from] = '+'
				canvas[row][to] = '+'
				canvas[row][pc] = '+'
			}
			if n.Left != nil {
				line(center(n.Left), pc)
			}
			if n.Right != nil {
				line(pc, center(n.Right))
			}
		}
		connect(n.Left)
		connect(n.Right)
	}
	connect(t.root)

	for _, line := range canvas {
		fmt.Println(strings.TrimRight(string(line), " "))
	}
}

// NLR — прямой обход: узел, левое поддерево, правое поддерево.
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

// generateKeys возвращает n уникальных трёхзначных чисел.
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
