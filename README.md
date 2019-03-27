Generic boolean expression evaluation on Go Data Structures [![GoDoc](https://godoc.org/github.com/hashicorp/go-bexpr?status.svg)](https://godoc.org/github.com/hashicorp/go-bexpr) [![CircleCI](https://circleci.com/gh/hashicorp/go-bexpr.svg?style=svg)](https://circleci.com/gh/hashicorp/go-bexpr)

`go-bexpr` is a Go (golang) library to provide generic boolean expression evaluation and filtering for Go data structures.

## Usage (Reflection)

This example program is available in [examples/simple](examples/simple)

```go
package main

import (
   "fmt"
   "github.com/hashicorp/go-bexpr"
)

type Example struct {
   X int

   // Can renamed a field with the struct tag
   Y string `bexpr:"y"`

   // Fields can use multiple names for accessing
   Z bool `bexpr:"Z,z,foo"`

   // Tag with "-" to prevent allowing this field from being used
   Hidden string `bexpr:"-"`

   // Unexported fields are not available for evaluation
   unexported string
}

func main() {
   value := map[string]Example{
      "foo": Example{X: 5, Y: "foo", Z: true, Hidden: "yes", unexported: "no"},
      "bar": Example{X: 42, Y: "bar", Z: false, Hidden: "no", unexported: "yes"},
   }

   expressions := []string{
      "foo.X == 5",
      "bar.y == bar",
      "foo.foo != false",
      "foo.z == true",
      "foo.Z == true",

      // will error in evaluator creation
      "bar.Hidden != yes",

      // will error in evaluator creation
      "foo.unexported == no",
   }

   for _, expression := range expressions {
      eval, err := bexpr.CreateEvaluatorForType(expression, nil, (*map[string]Example)(nil))

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
```

This will output:

```
Result of expression "foo.X == 5" evaluation: true
Result of expression "bar.y == bar" evaluation: true
Result of expression "foo.foo != false" evaluation: true
Result of expression "foo.z == true" evaluation: true
Result of expression "foo.Z == true" evaluation: true
Failed to create evaluator for expression "bar.Hidden != yes": Selector "bar.Hidden" is not valid
Failed to create evaluator for expression "foo.unexported == no": Selector "foo.unexported" is not valid
```
