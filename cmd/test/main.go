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
	slr.NewGenerator(rules)
}
