package main

import (
	"fmt"

	bexpr "github.com/hashicorp/go-bexpr"
)

type Filterable struct {
	X int
	Y int
}

func main() {
	before := [5]Filterable{
		Filterable{X: 1, Y: 2},
		Filterable{X: 2, Y: 3},
		Filterable{X: 3, Y: 4},
		Filterable{X: 4, Y: 5},
		Filterable{X: 5, Y: 6},
	}

	beforeMap := map[string]Filterable{
		"one":   Filterable{X: 1, Y: 2},
		"two":   Filterable{X: 2, Y: 3},
		"three": Filterable{X: 3, Y: 4},
		"four":  Filterable{X: 4, Y: 5},
		"five":  Filterable{X: 5, Y: 6},
	}

	filter, err := bexpr.CreateFilter("X == 2 or Y == 2", nil, (*Filterable)(nil))
	if err != nil {
		fmt.Printf("Failed to create filter: %v\n", err)
		return

	}
	after, err := filter.Execute(before)
	if err != nil {
		fmt.Printf("Failed to execute the filter: %v\n", err)
		return
	}

	fmt.Printf("Before: %v\n", before)
	fmt.Printf("After: %v\n", after)

	afterMap, err := filter.Execute(beforeMap)
	if err != nil {
		fmt.Printf("Failed to execute the filter: %v\n", err)
		return
	}

	fmt.Printf("Before: %v\n", beforeMap)
	fmt.Printf("After: %v\n", afterMap)
}
