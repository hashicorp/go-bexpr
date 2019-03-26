package main

import (
	"fmt"
	"os"

	bexpr "github.com/hashicorp/go-bexpr"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Please input an expression to parse.")
		return
	}

	for i, exp := range args {
		ast, err := bexpr.Parse(fmt.Sprintf("Expression %d", i), []byte(exp))

		if err != nil {
			fmt.Println(err)
		} else {
			ast.(bexpr.Expression).ExpressionDump(os.Stdout, "   ", 1)
		}
	}
}
