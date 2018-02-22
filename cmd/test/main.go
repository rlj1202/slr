// Production rules
//
// S -> E
// E -> E + T
// E -> T
// T -> id
// T -> ( E )
//
package main

import (
	"fmt"
	"github.com/rlj1202/slr"
)

const (
	rules = `
		S -> E
		E -> E + T
		E -> T
		T -> id
		T -> ( E )
	`
)

func main() {
	generator := slr.NewGenerator(rules)
	parser := generator.BuildParser()
	tree := parser.Parse([]string{
		"id", "+", "(", "id", "+", "id", ")",
	})
	if tree == nil {
		fmt.Println("parsing error")
	} else {
		printNode(tree, 0)
	}
}

func printBlank(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("\t")
	}
}

func printNode(node *slr.TreeNode, depth int) {
	printBlank(depth)
	fmt.Printf("%s {\n", node.Name)
	for _, leaf := range node.Leaves {
		printNode(leaf, depth+1)
	}
	printBlank(depth)
	fmt.Print("}\n")
}
