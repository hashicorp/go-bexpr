// bexpr is an implementation of a generic boolean expression evaluator.
// The general goal is to be able to evaluate some expression against some
package bexpr

import (
	"fmt"
	"strconv"
)

// Function type for usage with a SelectorConfiguration
type FieldValueCoercionFn func(value string) (interface{}, error)

// The FieldConfiguration struct represents how boolean expression
// validation and preparation should work for the given field. A field
// in this case is a single element of a selector.
//
// Example: foo.bar.baz has 3 fields separate by '.' characters.
type FieldConfiguration struct {
	// The field's name
	Name string

	// SubFields
	SubFields []FieldConfiguration

	// Another type that implements the ExpressionMatcher interface to
	// which operations on sub-fields should be delegated.
	Delegate ExpressionMatcher

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
	// FieldConfigurations is used during parsing of the expression to validate
	// that operations are valid for the particular fields and to further
	// validate/coerce the string input into something the ExpressionMatcher
	// wants to receive for its checks.
	FieldConfigurations() []FieldConfiguration

	// ExecuteMatch returns whether there was a match or not. We are not also
	// expecting any errors because all the validation bits are handled
	// during parsing and cross checking against the output of FieldConfigurations.
	ExecuteMatcher(sel Selector, op MatchOperator, value interface{}) bool
}

func CoerceInt(value string) (interface{}, error) {
	i, err := strconv.ParseInt(value, 10, 0)
	return int(i), err
}

func CoerceInt8(value string) (interface{}, error) {
	i, err := strconv.ParseInt(value, 10, 8)
	return int8(i), err
}

func CoerceInt16(value string) (interface{}, error) {
	i, err := strconv.ParseInt(value, 10, 16)
	return int16(i), err
}

func CoerceInt32(value string) (interface{}, error) {
	i, err := strconv.ParseInt(value, 10, 32)
	return int32(i), err
}

func CoerceInt64(value string) (interface{}, error) {
	i, err := strconv.ParseInt(value, 10, 64)
	return int64(i), err
}

type Expression struct {
	ast Expr
}

// TODO (mkeeler) The interface approach here isn't working out so fantastically
// 1 issue so far is that we use a type during Creation which presents the
// field configurations to be used for validation but then there is nothing to
// prevent passing a different value type to the Evaluate function

func Create(expression string, matcher ExpressionMatcher) (*Expression, error) {
	ast, err := Parse("", []byte(expression))
	expr := ast.(Expr)

	if err != nil {
		return nil, err
	}

	fields := matcher.FieldConfigurations()

	err = Validate(expr, fields)

	if err != nil {
		return nil, err
	}

	return &Expression{expr}, nil
}

func evaluateInternal(ast Expr, matcher ExpressionMatcher) bool {
	switch node := ast.(type) {
	case *UnaryExpr:
		switch node.Operator {
		case UnaryOpNot:
			return !evaluateInternal(node.Operand, matcher)
		}
	case *BinaryExpr:
		switch node.Operator {
		case BinaryOpAnd:
			return evaluateInternal(node.Left, matcher) && evaluateInternal(node.Right, matcher)
		case BinaryOpOr:
			return evaluateInternal(node.Left, matcher) || evaluateInternal(node.Right, matcher)
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
		return matcher.ExecuteMatcher(node.Selector, node.Operator, value)
	}
	return false
}

func (exp *Expression) Evaluate(matcher ExpressionMatcher) bool {
	return evaluateInternal(exp.ast, matcher)
}

func Validate(ast Expr, fields []FieldConfiguration) error {
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
				if fcfg.Name == field {
					found = true
					lastConfig = &fcfg
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
