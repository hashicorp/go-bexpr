package main

import (
	"fmt"

	bexpr "github.com/hashicorp/go-bexpr"
)

type Internal struct {
	Name string

	// Use an alternative name for referencing this field in expressions
	Values []int `bexpr:"fields"`

	// Hides this field so it cannot be used in expression evaluation
	Hidden int `bexpr:"-"`

	// Unexported fields are not available for use by the evaluator
	unexported int
}

type Matchable struct {
	Map      map[string]string
	X        int
	Internal Internal

	// Slices are handled specially. If any value within the slice
	// evaluates to true the overall evaluation for the slice will
	// be true.
	SliceInternal []Internal
	MapInternal   map[string]Internal
}

var data Matchable = Matchable{
	Map: map[string]string{
		"abc": "def",
		"ghi": "jkl",
		"mno": "pqr",
		"stu": "vwx",
		"y":   "z",
	},
	X: 5,
	Internal: Internal{
		Name:   "main",
		Values: []int{1, 2, 3, 4, 5},
	},
	SliceInternal: []Internal{
		{
			Name:   "odd",
			Values: []int{1, 3, 5, 7, 9},
		},
		{
			Name:   "even",
			Values: []int{2, 4, 6, 8, 10},
		},
		{
			Name:   "fib",
			Values: []int{0, 1, 1, 2, 3, 5},
		},
	},
	MapInternal: map[string]Internal{
		"odd": {
			Name:   "odd",
			Values: []int{1, 3, 5, 7, 9},
		},
		"even": {
			Name:   "even",
			Values: []int{2, 4, 6, 8, 10},
		},
		"fib": {
			Name:   "fib",
			Values: []int{0, 1, 1, 2, 3, 5},
		},
	},
}

type example struct {
	expression string
	variables  map[string]string
}

var examples []example = []example{
	// should error out in creating the evaluator as Foo is not a valid selector
	{expression: "Foo == 3"},
	// should error out because the field is hidden
	{expression: "Internal.Hidden == 5"},
	// should error out because the field is not exported
	{expression: "Internal.unexported == 3"},
	// should evaluate to true
	{expression: "Map[`abc`] == `def`"},
	// should evaluate to false
	{expression: "X == 3"},
	// should evaluate to true
	{expression: "Internal.fields is not empty"},
	// should evaluate to false
	{expression: "MapInternal.fib.Name != fib"},
	// should evaluate to true
	{expression: "odd in MapInternal"},
	// variable interpolation - should evaluate to true
	{expression: "X == ${value}", variables: map[string]string{"value": "5"}},
	// variable interpolation - should evaluate to false
	{expression: "X == ${value}", variables: map[string]string{"value": "4"}},
	// variable interpolation default value - should evaluate to false
	{expression: "X == ${value}"},
}

func main() {
	for _, ex := range examples {
		eval, err := bexpr.CreateEvaluator(ex.expression)

		if err != nil {
			fmt.Printf("Failed to create evaluator for expression %q: %v\n", ex.expression, err)
			continue
		}

		result, err := eval.Evaluate(data, ex.variables)
		if err != nil {
			fmt.Printf("Failed to run evaluation of expression %q (variables %#v): %v\n", ex.expression, ex.variables, err)
			continue
		}

		fmt.Printf("Result of expression %q evaluation (variables: %#v): %t\n", ex.expression, ex.variables, result)
	}
}
