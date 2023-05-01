// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"

	"github.com/hashicorp/go-bexpr"
)

type Example struct {
	X int

	// Can rename a field with the struct tag
	Y string `bexpr:"y"`
	Z bool   `bexpr:"baz"`

	// Tag with "-" to prevent allowing this field from being used
	Hidden string `bexpr:"-"`

	// Unexported fields are not available for evaluation
	unexported string
}

func main() {
	value := map[string]Example{
		"foo": {X: 5, Y: "foo", Z: true, Hidden: "yes", unexported: "no"},
		"bar": {X: 42, Y: "bar", Z: false, Hidden: "no", unexported: "yes"},
	}

	expressions := []string{
		"foo.X == 5",
		"bar.y == bar",
		"foo.baz == true",

		// will error in evaluator creation
		"bar.Hidden != yes",

		// will error in evaluator creation
		"foo.unexported == no",
	}

	for _, expression := range expressions {
		eval, err := bexpr.CreateEvaluator(expression)

		if err != nil {
			fmt.Printf("Failed to create evaluator for expression %q: %v\n", expression, err)
			continue
		}

		result, err := eval.Evaluate(value)
		if err != nil {
			fmt.Printf("Failed to run evaluation of expression %q: %v\n", expression, err)
			continue
		}

		fmt.Printf("Result of expression %q evaluation: %t\n", expression, result)
	}
}
