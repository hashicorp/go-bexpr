// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-bexpr/grammar"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Please input an expression to parse.")
		return
	}

	for i, exp := range args {
		ast, err := grammar.Parse(fmt.Sprintf("Expression %d", i), []byte(exp))

		if err != nil {
			fmt.Println(err)
		} else {
			ast.(grammar.Expression).ExpressionDump(os.Stdout, "   ", 1)
		}
	}
}
