package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-filter/bexpr/grammar"
)

func main() {
	args := os.Args[1:]

	for i, exp := range args {
		ast, err := grammar.Parse(fmt.Sprintf("Expression %d", i), []byte(exp))

		if err != nil {
			fmt.Println(err)
		} else {
			ast.(grammar.Expr).Dump(os.Stdout, "   ", 1)
		}
	}
}
