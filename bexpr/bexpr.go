// bexpr is an implementation of a generic boolean expression evaluator.
// The general goal is to be able to evaluate some expression against some
//
package bexpr

import (
	"fmt"
)

// Function type for usage with a SelectorConfiguration
type FieldValueCoercionFn func(value string) (interface{}, error)

// The FieldConfiguration struct represents how boolean expression
// validation and preparation should work for the given field. A field
// in this case is a single element of a selector.
//
// Example: foo.bar.baz has 3 fields separate by '.' characters.
type FieldConfiguration struct {
	// The field's name. If this is an empty string then it will be treated
	// as a wildcard and allow anything.
	Name string

	// SubFields
	SubFields []*FieldConfiguration

	// Function to run on the raw string value present in the expression
	// syntax to coerce into whatever form the ExpressionMatcher wants
	// The coercion happens only once and will then be passed as the `value`
	// parameter to all ExecuteMatcher invocations on the ExpressionMatcher.
	CoerceFn FieldValueCoercionFn

	// List of MatchOperators supported for this field. This configuration
	// is used to pre-validate an expressions fields before execution.
	SupportedOperations []MatchOperator
}

// ExpressionMatcher is the interface to implement to enable evaluating the boolean expressions
// against them.
type ExpressionMatcher interface {
	// ExecuteMatch returns whether there was a match or not. We are not also
	// expecting any errors because all the validation bits are handled
	// during parsing and cross checking against the output of FieldConfigurations.
	ExecuteMatcher(sel Selector, op MatchOperator, value interface{}) (bool, error)
}

type Expression struct {
	ast Expr
}

// TODO (mkeeler) The interface approach here isn't working out so fantastically
// 1 issue so far is that we use a type during Creation which presents the
// field configurations to be used for validation but then there is nothing to
// prevent passing a different value type to the Evaluate function

func Create(expression string, fields []*FieldConfiguration) (*Expression, error) {
	ast, err := Parse("", []byte(expression))

	if err != nil {
		return nil, err
	}

	expr := ast.(Expr)

	// Perform extra field validations if we were given a list of FieldConfigurations
	if len(fields) > 0 {
		err = Validate(expr, fields)

		if err != nil {
			return nil, err
		}
	}

	return &Expression{expr}, nil
}

func evaluateInternal(ast Expr, datum interface{}) (bool, error) {
	switch node := ast.(type) {
	case *UnaryExpr:
		switch node.Operator {
		case UnaryOpNot:
			result, err := evaluateInternal(node.Operand, datum)
			return !result, err
		}
	case *BinaryExpr:
		switch node.Operator {
		case BinaryOpAnd:
			result, err := evaluateInternal(node.Left, datum)
			if err != nil || result == false {
				return result, err
			}

			return evaluateInternal(node.Right, datum)

		case BinaryOpOr:
			result, err := evaluateInternal(node.Left, datum)
			if err != nil || result == true {
				return result, err
			}

			return evaluateInternal(node.Right, datum)
		}
	case *MatchExpr:
		var value interface{}
		if node.Value != nil {
			if node.Value.Converted != nil {
				value = node.Value.Converted
			} else {
				value = node.Value.Raw
			}
		}

		if matcher, ok := datum.(ExpressionMatcher); ok {
			return matcher.ExecuteMatcher(node.Selector, node.Operator, value)
		}

		return reflectEvaluateMatcher(node, datum)

		return false, fmt.Errorf("Reflection based evaluation not implemented")
	}
	return false, fmt.Errorf("Invalid AST node")
}

func (exp *Expression) Evaluate(datum interface{}) (bool, error) {
	return evaluateInternal(exp.ast, datum)
}

func Validate(ast Expr, fields []*FieldConfiguration) error {
	switch node := ast.(type) {
	case *UnaryExpr:
		return Validate(node.Operand, fields)
	case *BinaryExpr:
		if err := Validate(node.Left, fields); err != nil {
			return err
		}

		return Validate(node.Right, fields)
	case *MatchExpr:
		configs := fields
		var lastConfig *FieldConfiguration
		// validate the selector
		for idx, field := range node.Selector {
			found := false
			for _, fcfg := range configs {
				if fcfg.Name == field || fcfg.Name == "" {
					found = true
					lastConfig = fcfg
					configs = fcfg.SubFields
					break
				}
			}

			if !found {
				return fmt.Errorf("Selector %q is not valid", node.Selector[:idx])
			}
		}

		if lastConfig == nil {
			// TODO (mkeeler) - maybe a better error message here. I don't think this is possible with the parser
			// but I may want to take a closer look to be sure.
			return fmt.Errorf("Invalid expression without a selector")
		}

		// check the operator
		found := false
		for _, op := range lastConfig.SupportedOperations {
			if op == node.Operator {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Invalid match operator %s for selector %q", node.Operator, node.Selector)
		}

		// coerce/validate the value
		if node.Value != nil {
			if lastConfig.CoerceFn != nil {
				coerced, err := lastConfig.CoerceFn(node.Value.Raw)
				if err != nil {
					return fmt.Errorf("Failed to coerce value %q for selector %q: %v", node.Value.Raw, node.Selector, err)
				}

				node.Value.Converted = coerced
			}
		}
	default:
		return fmt.Errorf("Cannot validate: Invalid AST")
	}

	return nil
}
