package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Тернарный оператор
func Ternary[T any](flag bool, a, b T) T {
	if flag {
		return a
	}
	return b
}

// каждую строчку утверждение, можно представить в вибе объекта
type Statement struct {
	Action   string // действие о котором говорится в утверждении
	Actor    string // лицо, сделавшее утверждение
	Subject  string // лицо, о котором говрится в утверждении
	Negation bool   // имеет ли утверждение отрицание
	Self     bool   // тоесть говорит ли актор о самом себе, поможет проще считать очки
}

// строковые константы, чтобы меньше алоцировать память
const (
	ToBeI           = "am"
	ToBeHe          = "is" // she or it
	personalPronoun = "I"
	negativeAdverb  = "not"
	exclamationMark = "!"
	colon           = ":"
	delim           = '\n' //кроме этой, это руна
)

// метод конструктор, как бы пытаемся превратить строчку в Statement, а заодно провалидировать
func createStatement(sentence string) (*Statement, error) {

	words := strings.Fields(sentence)
	if !(len(words) == 4 || len(words) == 5) {
		return nil, errors.New("invalid format of sentence: unexpected statement")
	}

	actor := words[0]
	if !strings.HasSuffix(actor, colon) {
		return nil, errors.New("invalid format of sentence: try add \":\" for actor")
	}
	actor = strings.TrimSuffix(actor, colon)

	var self bool
	subject := words[1]
	beVerb := words[2]
	if beVerb == ToBeI && subject == personalPronoun {
		self = true
		subject = actor
	}

	if self && beVerb == ToBeHe {
		return nil, errors.New("invalid format of sentence: \"are's\" ARE not allowed :-) (если вы понимаете о чём я) ")
	}

	var negation bool
	action := words[3]
	if action == negativeAdverb {
		negation = true
		action = strings.TrimSuffix(words[4], exclamationMark)
	} else {
		action = strings.TrimSuffix(action, exclamationMark)
	}

	return &Statement{action, actor, subject, negation, self}, nil

}

// Прочитать, распарсить
func getStatements(reader *bufio.Reader) ([][]Statement, error) {

	parseQuantity := func(numStr string) (int, error) {
		return strconv.Atoi(strings.Trim(numStr, "\n"))
	}

	unitsQs, err := reader.ReadString(delim)
	if err != nil {
		return nil, err
	}

	unitQuantity, err := parseQuantity(unitsQs)
	if err != nil {
		return nil, err
	}

	dataSet := make([][]Statement, unitQuantity)

	for i := 0; i < unitQuantity; i++ {
		unitSizeS, err := reader.ReadString(delim)
		if err != nil {
			return nil, err
		}

		unitSize, err := parseQuantity(unitSizeS)
		if err != nil {
			return nil, err
		}

		unit := make([]Statement, unitSize)
		for j := 0; j < unitSize; j++ {
			word, err := reader.ReadString(delim)
			if err != nil {
				return nil, err
			}
			statement, err := createStatement(strings.TrimSpace(word))
			if err != nil {
				return nil, err
			}
			unit[j] = *statement
		}
		dataSet[i] = unit
	}
	return dataSet, nil
}

// решаем сет из данных
func ResolveStatements(set []Statement) ([]string, error) {
	setAction := set[0].Action
	scoreCap := Ternary(len(set) > 10, len(set)/4, len(set)/2)
	scoreTable := make(map[string]int, scoreCap)

	for _, statement := range set {
		if statement.Action != setAction {
			return nil, errors.New("all statements in set must have the same action")
		}

		subjectPointCounter := scoreTable[statement.Subject]

		if statement.Negation {
			scoreTable[statement.Subject] = subjectPointCounter - 1
		} else {
			scoreTable[statement.Subject] = subjectPointCounter + Ternary(statement.Self, 2, 1)
		}

		//если человек не выступил субъектом хотя бы один раз, он может набрать ноль очков. Инициализируем его здесь
		if actorPointCounter, ok := scoreTable[statement.Actor]; !ok {
			scoreTable[statement.Actor] = actorPointCounter
		}

	}

	createConclusion := func(name string) string {
		return fmt.Sprintf("%s %s %s.", name, ToBeHe, setAction)
	}

	var maxScore int = -100
	var conclusions []string
	for contester, pointCounter := range scoreTable {
		if pointCounter > maxScore {
			maxScore = pointCounter
			conclusions = []string{createConclusion(contester)}
		} else if pointCounter == maxScore {
			conclusions = append(conclusions, createConclusion(contester))
		}
	}

	return conclusions, nil

}

func main() {
	dataSet, err := getStatements(bufio.NewReader(os.Stdin))
	if err != nil {
		panic(err)
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for _, set := range dataSet {
		conclusions, err := ResolveStatements(set)
		if err != nil {
			panic(err)
		}
		sort.Strings(conclusions)
		for _, conclusion := range conclusions {
			fmt.Fprintln(out, conclusion)
		}
	}
}
