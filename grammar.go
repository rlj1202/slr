package slr

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	EPSILON = '^'
	EOS     = '$' // End Of String
)

type symbol struct {
	name string
	// true means terminal symbol
	// and false means non-terminal symbols
	terminal bool
}

type production struct {
	// Symbol id
	lhs int
	// Symbol ids
	rhs []int

	len int
}

type grammar struct {
	symbols           []*symbol
	symbolIndex       map[*symbol]int
	symbolIndexByName map[string]int

	productions          []*production
	productionIndex      map[*production]int
	productionIndexByLHS map[int][]int
}

func newGrammar(source string) *grammar {
	g := new(grammar)
	g.symbols = make([]*symbol, 0)
	g.symbolIndex = make(map[*symbol]int)
	g.symbolIndexByName = make(map[string]int)
	g.productions = make([]*production, 0)
	g.productionIndex = make(map[*production]int)
	g.productionIndexByLHS = make(map[int][]int)

	stringReader := strings.NewReader(source)
	scanner := bufio.NewScanner(stringReader)

	g.getSymbolIDByName(string([]byte{EOS}))
	g.getSymbolIDByName(string([]byte{EPSILON}))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		strs := strings.Split(string(line), "->")

		if len(strs) != 2 {
			continue
		}

		lhsName := strings.TrimSpace(strs[0])
		rhsNames := strings.Split(strings.TrimSpace(strs[1]), " ")

		p := &production{}
		p.lhs = g.getSymbolIDByName(lhsName)
		p.rhs = make([]int, len(rhsNames))
		for i, rhsName := range rhsNames {
			p.rhs[i] = g.getSymbolIDByName(rhsName)
		}
		p.len = len(p.rhs)

		g.symbols[p.lhs].terminal = false

		g.addProduction(p)
	}

	return g
}

func (g *grammar) getSymbolIDByName(name string) int {
	if id, okay := g.symbolIndexByName[name]; okay {
		return id
	} else {
		id := len(g.symbols)
		symbol := &symbol{name, true}
		g.symbols = append(g.symbols, symbol)
		g.symbolIndex[symbol] = id
		g.symbolIndexByName[name] = id

		return id
	}
}

func (g *grammar) addProduction(p *production) {
	id := len(g.productions)
	g.productions = append(g.productions, p)
	g.productionIndex[p] = id

	ids := g.productionIndexByLHS[p.lhs]
	if ids == nil {
		ids = make([]int, 0)
	}
	ids = append(ids, id)
	g.productionIndexByLHS[p.lhs] = ids
}

func (g *grammar) productionString(pId int) string {
	p := g.productions[pId]

	result := ""

	result += g.symbols[p.lhs].name
	result += " -> "
	for i, rhsId := range p.rhs {
		result += g.symbols[rhsId].name

		if i != len(p.rhs)-1 {
			result += " "
		}
	}

	return result
}

func (g *grammar) String() string {
	result := "Grammar{\n\tSymbols: ["

	for i, symbol := range g.symbols {
		result += fmt.Sprintf("%v", *symbol)

		if i != len(g.symbols)-1 {
			result += ", "
		}
	}

	result += "],\n\tProductions: [\n"

	for i := range g.productions {
		result += "\t\t\""
		result += g.productionString(i)
		result += "\""

		if i != len(g.productions)-1 {
			result += ",\n"
		}
	}

	result += "\n\t]\n}"

	return result
}
