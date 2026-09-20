package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func (h *hashTable) printTable() { // функция для читаемого вида хеш-функции
	const perRow = 10
	rows := (t + perRow - 1) / perRow

	fmt.Println("╔" + strings.Repeat("═══════╤", perRow-1) + "═══════╗")

	for r := 0; r < rows; r++ {
		fmt.Print("║")
		for c := 0; c < perRow; c++ {
			i := r*perRow + c
			if i >= t {
				fmt.Print("       │")
			} else {
				fmt.Printf(" %5d │", i)
			}
		}
		fmt.Println()

		fmt.Print("║")
		for c := 0; c < perRow; c++ {
			i := r*perRow + c
			if i >= t {
				fmt.Print("       │")
			} else if h.Table[i] == 0 {
				fmt.Print("   —   │")
			} else {
				fmt.Printf(" %5d │", h.Table[i])
			}
		}
		fmt.Println()

		if r < rows-1 {
			fmt.Println("╟" + strings.Repeat("───────┼", perRow-1) + "───────╢")
		}
	}

	fmt.Println("╚" + strings.Repeat("═══════╧", perRow-1) + "═══════╝")
}

func (h *hashTable) printKeys() { // функция для читаемого вида сгенерированных ключей
	list := make([]int, 0, len(h.Keys))
	for k := range h.Keys {
		list = append(list, k)
	}

	const perLine = 10
	for i, k := range list {
		fmt.Printf("%5d ", k)
		if (i+1)%perLine == 0 {
			fmt.Println()
		}
	}
	if len(list)%perLine != 0 {
		fmt.Println()
	}
}

const (
	numElements int = 45                          // число неповторяющихся элементов
	t           int = numElements + numElements/2 // размер хеш-таблицы
)

type hashTable struct {
	Table map[int]int
	Keys  map[int]struct{}
}

func (h *hashTable) newHashTable(t int) {
	h.Table = make(map[int]int, t)
	h.Keys = make(map[int]struct{})

	for i := 0; i < t; i++ {
		h.Table[i] = 0
	}

	for len(h.Keys) < numElements {
		newKey := rand.Intn(10000-1000) + 1000 // очередной сгенерированный ключ
		if _, ok := h.Keys[newKey]; ok {
			continue
		}
		h.Keys[newKey] = struct{}{}
	}

}

func hashFunc(key int) int {
	if key > 9999 || key < 1000 {
		panic("Значение ключа не 4-х значное!")
	}

	return ((key/100)%10 + (key/10)%10 + key%10) % t
}

func (h *hashTable) fill() {
	totalSteps := 0 // суммарное число шагов
	failed := 0     // количество ключей, которые не влезли

	for newKey := range h.Keys {
		a0 := hashFunc(newKey) // первичный адрес
		if h.Table[a0] == 0 {
			h.Table[a0] = newKey
			totalSteps++
		} else {
			// fmt.Printf("Коллизия - ключ %v, адрес %v занят ключом %v\n", newKey, a0, h.Table[a0])
			inserted := false
			for i := 1; i < 20; i++ {
				ai := (a0 + i*i) % t
				if h.Table[ai] == 0 {
					h.Table[ai] = newKey
					inserted = true
					totalSteps += i + 1
					break
				} else {
					// fmt.Printf("Коллизия - ключ %v, адрес %v занят ключом %v\n", newKey, ai, h.Table[ai])
				}
			}

			if !inserted {
				// fmt.Printf("Квадратичные пробы не помогли, линейная проба для ключа %v\n", newKey)
				for j := 1; j <= t; j++ {
					ai := (a0 + j) % t
					if h.Table[ai] == 0 {
						h.Table[ai] = newKey
						inserted = true
						totalSteps += 20 + j
						break
					} else {
						// fmt.Printf("Коллизия - ключ %v, адрес %v занят ключом %v\n", newKey, ai, h.Table[ai])
					}
				}
			}

			if !inserted {
				// fmt.Printf("Ключ %v не удалось вставить в хеш-таблицу\n", newKey)
				failed++
			}
		}
	}

	fmt.Println()
	h.printTable()
	fmt.Println("=== Параметры хеш-таблицы ===")
	fmt.Printf("Коэффициент заполнения:    %.2f\n", float64(numElements-failed)/float64(t))
	fmt.Printf("Среднее число шагов:        %.2f\n", float64(totalSteps)/float64(numElements))
}

func main() {
	myTable := hashTable{}
	myTable.newHashTable(t)
	myTable.printKeys()
	myTable.fill()
}
