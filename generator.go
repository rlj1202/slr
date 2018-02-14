package slr

import "fmt"

type itemSet struct {
	items []item
	marks map[item]struct{}

	// map[symId]itemSetId
	gotos map[int]int
}

type item struct {
	// production id
	id int
	// position of dot
	mark int
}

type Generator struct {
	*grammar

	firsts  []map[int]struct{}
	follows []map[int]struct{}

	itemSets     []*itemSet
	itemSetIndex map[item]int
}

func newItemSet() *itemSet {
	set := new(itemSet)
	set.items = make([]item, 0)
	set.marks = make(map[item]struct{})
	set.gotos = make(map[int]int)

	return set
}

func (set *itemSet) add(i item) bool {
	if _, exists := set.marks[i]; !exists {
		set.items = append(set.items, i)
		set.marks[i] = struct{}{}
		return true
	}
	return false
}

func (set *itemSet) equals(a *itemSet) bool {
	if len(set.items) != len(set.items) {
		return false
	}
	for _, itm := range set.items {
		if _, okay := a.marks[itm]; !okay {
			return false
		}
	}
	return true
}

// LHS symbol of first production rule will be the start symbol.
func NewGenerator(source string) *Generator {
	generator := new(Generator)
	generator.grammar = newGrammar(source)
	generator.firsts = make([]map[int]struct{}, len(generator.symbols))
	generator.follows = make([]map[int]struct{}, len(generator.symbols))
	generator.itemSets = make([]*itemSet, 0)
	generator.itemSetIndex = make(map[item]int)

	fmt.Println("input source = " + source)
	fmt.Println()

	fmt.Println("symbols = ")
	for i, symbol := range generator.grammar.symbols {
		fmt.Printf("\t%02d: %s, %t\n", i, symbol.name, symbol.terminal)
	}
	fmt.Println()

	i0 := newItemSet()
	i0.add(item{0, 0}) // first production and mark at zero
	generator.addItemSet(i0)

	i0 = generator.closure(i0)

	fmt.Println("closure(I0) = ")
	for _, itm := range i0.items {
		fmt.Printf("\t%s\n", generator.itemString(itm))
	}
	fmt.Println()

	for i := 0; i < len(generator.itemSets); i++ {
		curItemSet := generator.itemSets[i]

		for symId, symbol := range generator.symbols {
			newSet := generator.go2(curItemSet, symId)
			if len(newSet.items) == 0 {
				continue
			}
			gotoSetId := generator.addItemSet(newSet)
			curItemSet.gotos[symId] = gotoSetId

			fmt.Printf("goto(I%d, %s) = I%d = \n", i, symbol.name, gotoSetId)
			for _, itm := range newSet.items {
				fmt.Printf("\t%s\n", generator.itemString(itm))
			}
			fmt.Println()
		}
	}

	for i, set := range generator.itemSets {
		fmt.Printf("I%d = \n", i)
		for _, itm := range set.items {
			fmt.Printf("\t%s\n", generator.itemString(itm))
		}
		fmt.Println()
	}

	generator.first()
	for symId, symbols := range generator.firsts {
		fmt.Printf("first(\"%s\") = {", generator.symbols[symId].name)
		for firstId := range symbols {
			fmt.Printf(" \"%s\"", generator.symbols[firstId].name)
		}
		fmt.Println(" }")
	}
	fmt.Println()

	generator.follow()
	for symId, symbols := range generator.follows {
		if generator.symbols[symId].terminal {
			continue
		}
		fmt.Printf("follow(\"%s\") = {", generator.symbols[symId].name)
		for firstId := range symbols {
			fmt.Printf(" \"%s\"", generator.symbols[firstId].name)
		}
		fmt.Println(" }")
	}
	fmt.Println()

	table := make([][]string, len(generator.symbols))
	for i := range table {
		table[i] = make([]string, len(generator.itemSets))
	}
	for itemSetId, itemSet := range generator.itemSets {
		for symId, gotoSetId := range itemSet.gotos {
			if generator.symbols[symId].terminal {
				table[symId][itemSetId] = fmt.Sprintf("s%d", gotoSetId)
			} else {
				table[symId][itemSetId] = fmt.Sprintf("%d", gotoSetId)
			}
		}
		for _, singleItem := range itemSet.items {
			pro := generator.productions[singleItem.id]

			if singleItem.mark == pro.len {
				followIds := generator.follows[pro.lhs]
				for followId := range followIds {
					table[followId][itemSetId] = fmt.Sprintf("r%d", singleItem.id)
				}
			}
		}
	}

	fmt.Print("PS ")
	for _, symbol := range generator.symbols {
		if symbol.terminal {
			fmt.Printf(" | %2s", symbol.name)
		}
	}
	for _, symbol := range generator.symbols {
		if !symbol.terminal {
			fmt.Printf(" | %2s", symbol.name)
		}
	}
	fmt.Println()
	for itemSetId := range generator.itemSets {
		fmt.Printf("%02d ", itemSetId)
		for symId := range generator.symbols {
			if generator.symbols[symId].terminal {
				fmt.Printf("   %2s", table[symId][itemSetId])
			}
		}
		for symId := range generator.symbols {
			if !generator.symbols[symId].terminal {
				fmt.Printf("   %2s", table[symId][itemSetId])
			}
		}
		fmt.Println()
	}

	return generator
}

func (generator *Generator) first() {
	for symId, symbol := range generator.symbols {
		if symbol.terminal {
			generator.addFirstSymbol(symId, symId)
		}
	}
}

func (generator *Generator) follow() {
	generator.addFollowSymbol(generator.productions[0].lhs, 0)

	for _, pro := range generator.productions {
		if len(pro.rhs) < 2 {
			continue
		}
		for i, symId := range pro.rhs[:len(pro.rhs)-1] {
			followId := pro.rhs[i+1]

			if !generator.symbols[symId].terminal {
				for followSymId := range generator.firsts[followId] {
					generator.addFollowSymbol(symId, followSymId)
				}
			}
		}
	}
}

func (generator *Generator) addFirstSymbol(nter, ter int) {
	symbols := generator.firsts[nter]
	if symbols == nil {
		symbols = make(map[int]struct{})
	}
	symbols[ter] = struct{}{}
	generator.firsts[nter] = symbols

	for _, pro := range generator.productions {
		if nter == pro.rhs[0] && pro.lhs != pro.rhs[0] {
			generator.addFirstSymbol(pro.lhs, ter)
		}
	}
}

func (generator *Generator) addFollowSymbol(nter, ter int) {
	symbols := generator.follows[nter]
	if symbols == nil {
		symbols = make(map[int]struct{})
	}
	symbols[ter] = struct{}{}
	generator.follows[nter] = symbols

	for _, proId := range generator.productionIndexByLHS[nter] {
		pro := generator.productions[proId]
		symId := pro.rhs[len(pro.rhs)-1]
		sym := generator.symbols[symId]

		if !sym.terminal && pro.lhs != symId {
			generator.addFollowSymbol(symId, ter)
		}
	}
}

func (generator *Generator) addItemSet(set *itemSet) int {
	if id, exists := generator.getItemSetId(set); exists {
		return id
	}

	id := len(generator.itemSets)
	generator.itemSets = append(generator.itemSets, set)
	generator.itemSetIndex[set.items[0]] = id

	return id
}

// itemSetId, exists
func (generator *Generator) getItemSetId(set *itemSet) (int, bool) {
	for i, a := range generator.itemSets {
		if a.equals(set) {
			return i, true
		}
	}

	return -1, false
}

func (generator *Generator) go2(set *itemSet, symId int) *itemSet {
	newSet := newItemSet()

	for _, itm := range set.items {
		if generator.peek(itm) == symId {
			newSet.add(item{itm.id, itm.mark + 1})
		}
	}

	return generator.closure(newSet)
}

func (generator *Generator) closure(set *itemSet) *itemSet {
	for i := 0; i < len(set.items); i++ {
		curItem := set.items[i]
		symId := generator.peek(curItem)
		if symId == -1 {
			continue
		}
		if !generator.grammar.symbols[symId].terminal {
			for _, pId := range generator.grammar.productionIndexByLHS[symId] {
				set.add(item{pId, 0})
			}
		}
	}

	return set
}

func (generator *Generator) peek(i item) int {
	production := generator.grammar.productions[i.id]
	if i.mark >= production.len {
		return -1
	}
	return production.rhs[i.mark]
}

func (generator *Generator) itemString(itm item) string {
	p := generator.productions[itm.id]

	result := ""

	result += generator.symbols[p.lhs].name
	result += " -> "
	for i, rhsId := range p.rhs {
		if i == itm.mark {
			result += ". "
		}

		result += generator.symbols[rhsId].name

		if i != len(p.rhs)-1 {
			result += " "
		}
	}

	if itm.mark == len(p.rhs) {
		result += " ."
	}

	return result
}

func (generator *Generator) String() string {
	return generator.grammar.String()
}
