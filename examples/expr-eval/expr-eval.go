package main

import (
	"fmt"

	"github.com/hashicorp/go-filter/bexpr"
)

type Internal struct {
	Name   string
	Values []int
}

type Matchable struct {
	Map           map[string]string
	X             int
	Internal      Internal
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
		Internal{
			Name:   "odd",
			Values: []int{1, 3, 5, 7, 9},
		},
		Internal{
			Name:   "even",
			Values: []int{2, 4, 6, 8, 10},
		},
		Internal{
			Name:   "fib",
			Values: []int{0, 1, 1, 2, 3, 5},
		},
	},
	MapInternal: map[string]Internal{
		"odd": Internal{
			Name:   "odd",
			Values: []int{1, 3, 5, 7, 9},
		},
		"even": Internal{
			Name:   "even",
			Values: []int{2, 4, 6, 8, 10},
		},
		"fib": Internal{
			Name:   "fib",
			Values: []int{0, 1, 1, 2, 3, 5},
		},
	},
}

var expressions []string = []string{
	// should error out in creating the evaluator as Foo is not a valid selector
	"Foo == 3",
	// should evaluate to true
	"Map[`abc`] == `def`",
	// should evaluate to false
	"X == 3",
	// should evaluate to true
	"Internal.Values is not empty",
	// should evaluate to true
	"0 in SliceInternal.Values",
	// should evaluate to false
	"4 not in SliceInternal.Values",
	// should evaluate to false
	"MapInternal.fib.Name != fib",
	// should evaluate to true
	"odd in MapInternal",
}

func main() {
	for _, expression := range expressions {
		eval, err := bexpr.CreateForType(expression, nil, (*Matchable)(nil))

		if err != nil {
			fmt.Printf("Failed to create evaluator for expression %q: %v\n", expression, err)
			continue
		}

		result, err := eval.Evaluate(data)
		if err != nil {
			fmt.Printf("Failed to run evaluation of expression %q: %v\n", expression, err)
			continue
		}

		fmt.Printf("Result of expression %q evaluation: %t\n", expression, result)
	}
}
