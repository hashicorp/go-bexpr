package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-filter/filter"
)

func main() {
	args := os.Args[1:]

	for i, flt := range args {
		ast, err := filter.Parse(fmt.Sprintf("Filter %d", i), []byte(flt))

		if err != nil {
			fmt.Println(err)
		} else {
			ast.(filter.Expr).Dump(os.Stdout, "   ", 1)
		}
	}
}
