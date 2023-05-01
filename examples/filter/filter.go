// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
		{X: 1, Y: 2},
		{X: 2, Y: 3},
		{X: 3, Y: 4},
		{X: 4, Y: 5},
		{X: 5, Y: 6},
	}

	beforeMap := map[string]Filterable{
		"one":   {X: 1, Y: 2},
		"two":   {X: 2, Y: 3},
		"three": {X: 3, Y: 4},
		"four":  {X: 4, Y: 5},
		"five":  {X: 5, Y: 6},
	}

	filter, err := bexpr.CreateFilter("X == 2 or Y == 2")
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
