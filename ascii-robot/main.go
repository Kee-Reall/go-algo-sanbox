package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Up = iota + 1
	Down
	lineDelimiterRune = '\n'
	lineDelimiterStr  = "\n"
)

func parseSetQuantity(reader *bufio.Reader) (int, error) {
	numStr, err := reader.ReadString(lineDelimiterRune)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.Trim(numStr, lineDelimiterStr))
}

func parseUnitQuantity(reader *bufio.Reader) (a int, b int, err error) {
	numsStr, err := reader.ReadString(lineDelimiterRune)
	if err != nil {
		return
	}
	nums := strings.Fields(numsStr)
	if a, err = strconv.Atoi(nums[0]); err != nil {
		return
	}
	b, err = strconv.Atoi(nums[1])
	return
}

func fillRowForStore(reader *bufio.Reader, row []string) { //вообще можно определить координаты уже здесь, прокинув ординату
	for i, _ := range row { // но я это понял пока писал их поиск, и сюда один фиг почти никто не посмотрит
		run, _ := reader.ReadByte() // типо если нашли A или B как то вернуть координаты заодно
		row[i] = string(run)
	}
	reader.ReadByte() //skip 'n'
}

func ReadSetAndResolve(reader *bufio.Reader) [][][]string {
	setQuantity, _ := parseSetQuantity(reader)

	r := make([][][]string, setQuantity)

	for i := 0; i < setQuantity; i++ {
		n, m, _ := parseUnitQuantity(reader)
		storageMap := make(StorageMap, n)
		for j := 0; j < m; j++ {
			storageMap[j] = make([]string, m)
			fillRowForStore(reader, storageMap[j])

		}

		DrawRoute(storageMap)

		r[i] = storageMap

	}

	return r
}

type StorageMap [][]string

func (sm StorageMap) GetMaxCoords() (int, int) {
	return len(sm[0]) - 1, len(sm) - 1
}

type Robot struct {
	Abscissa int
	Ordinate int
	Symbol   string
}

func (bot *Robot) IsInPlace(sm StorageMap) (bool, int) {
	x, y := sm.GetMaxCoords()
	if bot.Abscissa == 0 && bot.Ordinate == 0 {
		return true, Up
	}

	if bot.Abscissa == x && bot.Ordinate == y {
		return true, Down
	}

	return false, 0
}

func (bot *Robot) markWithStore(sm StorageMap) {
	sm[bot.Abscissa][bot.Ordinate] = bot.GetRouteSymbol()
}

func (bot *Robot) StepUp(sm StorageMap) {
	bot.Ordinate -= 1
	bot.markWithStore(sm)
}

func (bot *Robot) StepDown(sm StorageMap) {
	bot.Ordinate += 1
	bot.markWithStore(sm)
}

func (bot *Robot) StepLeft(sm StorageMap) {
	bot.Abscissa -= 1
	bot.markWithStore(sm)
}

func (bot *Robot) StepRight(sm StorageMap) {
	bot.Abscissa += 1
	bot.markWithStore(sm)
}

func (bot *Robot) GetRouteSymbol() string {
	return strings.ToLower(bot.Symbol)
}

func FindRobots(store StorageMap) (*Robot, *Robot) {
	var first *Robot = nil
	var second *Robot = nil
loop:
	for i := 0; i < len(store); i++ {
		for j := 0; j < len(store[i]); j++ {
			if store[i][j] == "A" || store[i][j] == "B" {
				if first == nil {
					first = &Robot{j, i, store[i][j]}
					continue
				}
				second = &Robot{j, i, store[i][j]}
				break loop
			}
		}
	}
	return first, second
}

func DrawRoute(storeMap StorageMap) { // решение начинается здесь
	firstBot, secondBot := FindRobots(storeMap)
	if firstBot == nil || secondBot == nil {
		panic("ТУТ ДОЛЖНГЫ БЫТЬ БОТЫ!")
	}
	firstInPlace, firstDir := firstBot.IsInPlace(storeMap)
	secondInPlace, secondDir := secondBot.IsInPlace(storeMap)

	if firstInPlace && secondInPlace {
		return // ничего не надо делать, боты в точках
	}

	var upper *Robot = nil
	var downer *Robot = nil

	isAtLeastOneInPlace := firstInPlace || secondInPlace

	if isAtLeastOneInPlace {
		if firstInPlace {
			if firstDir == Up {
				upper = firstBot
				downer = secondBot
			} else {
				upper = secondBot
				downer = firstBot
			}
		} else {
			if secondDir == Up {
				upper = secondBot
				downer = firstBot
			} else {
				downer = secondBot
				upper = firstBot
			}
		}
	} else { // сюда проваливаться будем чаще всего
		LinkRobots(upper, downer, firstBot, secondBot, storeMap)
	}

	if inPlace, _ := upper.IsInPlace(storeMap); !inPlace {
		MoveUpper(upper, storeMap)
	}

	if inPlace, _ := downer.IsInPlace(storeMap); !inPlace {
		MoveDowner(downer, storeMap)
	}
}

func MoveUpper(upper *Robot, storeMap [][]string) {
	if upper.Ordinate&1 == 0 { // идем влево до конца
		for upper.Abscissa != 0 {
			upper.StepLeft(storeMap)
		}
	}
	inPlace, _ := upper.IsInPlace(storeMap)
	if inPlace {
		return // пришли
	}
	for upper.Ordinate == 0 { //
		upper.StepUp(storeMap)
	}
}

func MoveDowner(downer *Robot, storeMap [][]string) {

}

func LinkRobots(upper, downer, first, second *Robot, sm StorageMap) {
	maxX, maxY := sm.GetMaxCoords()
	firstToStart := first.Abscissa + first.Ordinate
	firstToEnd := (maxX - first.Abscissa) + (maxY - first.Ordinate)
	secondToStart := second.Abscissa + second.Ordinate
	secondToEnd := (maxX - second.Abscissa) + (maxY - second.Ordinate)
	if (firstToStart + secondToEnd) < (firstToEnd + secondToStart) {
		upper = first
		downer = second
	} else {
		upper = second
		downer = first
	}
}

func main() {
	set := ReadSetAndResolve(bufio.NewReader(os.Stdin))
	PrintSet(set)
}

func PrintSet(set [][][]string) {
	for _, storage := range set {
		for _, row := range storage {
			fmt.Println(strings.Join(row, ""))
		}
	}
}
