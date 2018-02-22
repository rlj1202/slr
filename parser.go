package slr

import "fmt"

type Parser struct {
	*grammar

	// [symbol][state]data
	//
	// if data == 0 {
	//     null
	// } else if data > 0 {
	//     shift for terminal symbols
	//     goto for non-terminal symbols
	// } else if data < 0 {
	//     reduce
	// }
	table [][]int

	stateStack  []int
	symbolStack []int
}

type TreeNode struct {
	Name   string
	Leaves []*TreeNode
}

func newParser(g *grammar, table [][]int) *Parser {
	parser := new(Parser)
	parser.grammar = g
	parser.table = table
	parser.stateStack = make([]int, 0)
	parser.symbolStack = make([]int, 0)

	return parser
}

func newTreeNode(name string) *TreeNode {
	node := new(TreeNode)
	node.Name = name
	node.Leaves = make([]*TreeNode, 0)

	return node
}

func (parser *Parser) Parse(symbols []string) *TreeNode {
	nodes := make([]*TreeNode, 0)
	symbols = append(symbols, "$")
	parser.pushState(0)

	for {
		symName := symbols[0]
		symId, exists := parser.grammar.symbolIndexByName[symName]
		if !exists {
			return nil
		}
		symbol := parser.grammar.symbols[symId]
		state := parser.peekState()

		inst := parser.table[symId][state]

		if inst == 0 {
			// error
			return nil
		} else if inst > 0 {
			// shift
			parser.pushState(inst - 1)
			parser.pushSymbol(symId)
			symbols = symbols[1:]

			nodes = append(nodes, newTreeNode(symName))

			fmt.Printf("ACTION(%d, %s) = shift %d\n", state, symbol.name, parser.peekState())
		} else if inst < 0 {
			// reduce

			proId := -inst - 1
			pro := parser.grammar.productions[proId]

			node := newTreeNode(parser.grammar.symbols[pro.lhs].name)
			for i := 0; i < pro.len; i++ {
				parser.popSymbol()
				parser.popState()
				node.Leaves = append(node.Leaves, nodes[len(nodes)-1])
				nodes = nodes[:len(nodes)-1]
			}
			nodes = append(nodes, node)

			parser.pushSymbol(pro.lhs)

			if proId == 0 {
				fmt.Println("ACCEPTED")
				break
			}
			fmt.Printf("ACTION(%d, %s) = reduce %d\n", state, symbol.name, proId)

			state = parser.peekState()
			gotoState := parser.table[pro.lhs][state]
			if gotoState == 0 {
				return nil
			}
			gotoState--
			parser.pushState(gotoState)

			fmt.Printf("ACTION(%d, %s) = goto %d\n", state, parser.grammar.symbols[pro.lhs].name, parser.peekState())

		}

		if len(symbols) == 0 {
			break
		}
	}

	return nodes[0]
}

func (parser *Parser) peekState() int {
	return parser.stateStack[len(parser.stateStack)-1]
}

func (parser *Parser) pushState(state int) {
	parser.stateStack = append(parser.stateStack, state)
}

func (parser *Parser) popState() int {
	len := len(parser.stateStack)
	state := parser.stateStack[len-1]
	parser.stateStack = parser.stateStack[:len-1]
	return state
}

func (parser *Parser) pushSymbol(symId int) {
	parser.symbolStack = append(parser.symbolStack, symId)
}

func (parser *Parser) popSymbol() int {
	len := len(parser.symbolStack)
	symId := parser.symbolStack[len-1]
	parser.symbolStack = parser.symbolStack[:len-1]
	return symId
}

func (treeNode *TreeNode) addLeaf(leaf *TreeNode) {
	treeNode.Leaves = append(treeNode.Leaves, leaf)
}
