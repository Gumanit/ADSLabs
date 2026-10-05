package main

/*
Убрать логирование коллизий, добавить возможность добавлять выбираемое количество сгенерированных элементов,
при замене всегда возвращать статус 2 ключей и их индексы.
*/

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func (h *hashTable) printTable() { // функция для читаемого вида хеш-функции
	fmt.Println()
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

	fmt.Println("=== Параметры хеш-таблицы ===")
	fmt.Printf("Коэффициент заполнения:    %.2f\n", float64(numElements-failed)/float64(t))
	fmt.Printf("Среднее число шагов:        %.2f\n", float64(totalSteps)/float64(numElements))
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

var (
	totalSteps  int = 0                           // суммарное число шагов
	failed      int = 0                           // количество ключей, которые не влезли
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

	h.printTable()
}

func (h *hashTable) search(key int) (int, bool) {
	if key > 9999 || key < 1000 {
		panic("Значение ключа не 4-х значное!")
	}

	a0 := hashFunc(key) // первичный адрес
	switch h.Table[a0] {
	case 0:
		return -1, false
	case key:
		return a0, true
	default:
		sqr := false
		for i := 1; i < 20; i++ {
			ai := (a0 + i*i) % t
			if h.Table[ai] == key {
				sqr = true
				return ai, true
			}
		}

		if !sqr {
			for j := 1; j <= t; j++ {
				ai := (a0 + j) % t
				if h.Table[ai] == key {
					return ai, true
				}
			}
		}
	}
	return -1, false

}

func (h *hashTable) delete(key int) (int, bool) {
	idx, found := h.search(key)
	if found {
		delete(h.Table, idx)
		h.Table[idx] = 0
		numElements--
		return idx, true
	} else {
		return -1, false
	}
}

func (h *hashTable) add(key int) (int, bool) {
	if numElements == t {
		fmt.Println("Хеш-таблица полностью заполнена, вставка новых элементов невозможна!")
		return -1, false
	}
	idx, found := h.search(key)
	if found {
		return idx, false
	} else {
		a0 := hashFunc(key) // первичный адрес
		if h.Table[a0] == 0 {
			h.Table[a0] = key
			numElements++
			totalSteps++
			return a0, true
		} else {
			// fmt.Printf("Коллизия - ключ %v, адрес %v занят ключом %v\n", key, a0, h.Table[a0])
			inserted := false
			for i := 1; i < 20; i++ {
				ai := (a0 + i*i) % t
				if h.Table[ai] == 0 {
					h.Table[ai] = key
					inserted = true
					totalSteps += i + 1
					numElements++
					return ai, true
				} else {
					// fmt.Printf("Коллизия - ключ %v, адрес %v занят ключом %v\n", key, ai, h.Table[ai])
				}
			}

			if !inserted {
				// fmt.Printf("Квадратичные пробы не помогли, линейная проба для ключа %v\n", key)
				for j := 1; j <= t; j++ {
					ai := (a0 + j) % t
					if h.Table[ai] == 0 {
						h.Table[ai] = key
						inserted = true
						totalSteps += 20 + j
						numElements++
						return ai, true
					} else {
						// fmt.Printf("Коллизия - ключ %v, адрес %v занят ключом %v\n", key, ai, h.Table[ai])
					}
				}
			}

			if !inserted {
				failed++
			}
		}
	}
	return -1, false
}

func (h *hashTable) replace(delKey, addKey int) (int, bool) {
	if idx, found := h.search(addKey); found {
		return idx, false
	}

	if _, delFound := h.delete(delKey); !delFound {
		return -1, false
	}

	addIndx, addFound := h.add(addKey)
	if !addFound {
		return -1, false
	}

	return addIndx, true
}

// func (h *hashTable) addGen(num int) (int, int) {
// 	for num < 0 {
// 		key := rand.Intn(10000-1000) + 1000
// 		if idx, found := h.add(key); found {
// 			continue
// 		} else {
// 			num--
// 			return idx, key
// 		}
// 	}

// }

func main() {
	myTable := hashTable{}
	myTable.newHashTable(t)
	myTable.printKeys()
	myTable.fill()

	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			return
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Введите число от 1 до 5")
			continue
		}

		if !handleChoice(choice, &myTable) {
			fmt.Println("До свидания!")
			return
		}
	}
}

func userInput() (int, bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите 4-х значный ключ: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return 0, false
	}
	input = strings.TrimSpace(input)
	key, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Введите корректный 4-х значный ключ")
		return 0, false
	}
	return key, true
}

func userInputGen() (int, bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите количество генерируемых ключей: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return 0, false
	}
	input = strings.TrimSpace(input)
	num, err := strconv.Atoi(input)
	if err != nil || num == 0 {
		fmt.Println("Введите корректное значение")
		return 0, false
	}
	avaiableNum := t - numElements
	if num > avaiableNum {
		fmt.Printf("Введенное количество элементов %d превышает свободный объем хеш-таблицы, будет добавлен возможный максимум: %d элементов\n", num, avaiableNum)
		num = avaiableNum
	}
	return num, true
}

func handleChoice(choice int, h *hashTable) bool {
	fmt.Println()

	switch choice {
	case 1: // Выполнить поиск элемента в хеш-таблице

		key, valid := userInput()
		if !valid {
			return true
		}

		index, found := h.search(key)
		if found {
			fmt.Printf("Ключ %d найден в хеш-таблице по индексу %d\n", key, index)
		} else {
			fmt.Printf("Ключ %d не найден в хеш-таблице\n", key)
		}

	case 2: // Удалить элемент из хеш-таблицы
		key, valid := userInput()
		if !valid {
			return true
		}

		idx, isDeleted := h.delete(key)
		if isDeleted {
			fmt.Printf("Элемент %d удален из индекса %d\n", key, idx)
		} else {
			fmt.Printf("Элемент не найден\n")
		}

	case 3: // Добавить элемент в хеш-таблицу

		key, valid := userInput()
		if !valid {
			return true
		}

		idx, isAdded := h.add(key)
		if idx != -1 && isAdded {
			fmt.Printf("Элемент %d добавлен в индекс %d\n", key, idx)
		} else if idx != -1 && !isAdded {
			fmt.Printf("Элемент %d уже есть в хеш-таблице по индексу %d\n", key, idx)
		} else if idx == -1 {
			fmt.Printf("Элемент %d не удалось вставить в хеш-таблицу\n", key)
		}

	case 4: // Заменить элемент в хеш-таблице

		fmt.Print("Введите заменяемый и добавляемый элементы\n")
		delKey, valid := userInput()
		if !valid {
			return true
		}

		addKey, valid := userInput()
		if !valid {
			return true
		}

		idx, isReplased := h.replace(delKey, addKey)
		if idx != -1 && isReplased {
			fmt.Printf("Элемент %d удален, элемент %d добавлен в хеш-таблицу под индексом %d", delKey, addKey, idx)
		} else if idx != -1 && !isReplased {
			fmt.Printf("Элемент %d уже существует в хеш-таблице", addKey)
		} else if idx == -1 && !isReplased {
			fmt.Printf("Заменяемый элемент %d не найден", delKey)
		}

	case 5: // сгенерировать элементы в хеш-таблице
		num, valid := userInputGen()
		if !valid {
			return true
		}

		for num > 0 {
			key := rand.Intn(10000-1000) + 1000
			if idx, isAdded := h.add(key); idx != -1 && isAdded {
				fmt.Printf("Элемент %d добавлен по индексу %d\n", key, idx)
				num--
			} else if idx != -1 && !isAdded {
				continue
			} else if idx == -1 && !isAdded {
				fmt.Println("Вставка невозможна")
			}
		}

	case 6: // Выход
		return false
	default:
		fmt.Println("Введите число от 1 до 6")
	}
	h.printTable()
	return true
}

func printMenu() {
	fmt.Println()
	fmt.Println("========== МЕНЮ ==========")
	fmt.Println("1. Найти элемент в хеш-таблице")
	fmt.Println("2. Удалить элемент в хеш-таблице")
	fmt.Println("3. Добавить элемент в хеш-таблицу")
	fmt.Println("4. Заменить элемент в хеш-таблице")
	fmt.Println("5. Добавить сгенерированные элементы в хеш-таблице")
	fmt.Println("6. Выход")
	fmt.Println("==========================")
	fmt.Print("Ваш выбор: ")
}
