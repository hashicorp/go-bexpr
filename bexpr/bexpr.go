// bexpr is an implementation of a generic boolean expression evaluator.
// The general goal is to be able to evaluate some expression against some
//
package bexpr

import (
	"fmt"
	"reflect"
)

// ExpressionEvaluator is the interface to implement to provide custom evaluation
// logic for a selector. This could be used to enable synthetic fields or other
// more complex logic that the default behavior does not support
//
// Currently there is no way to go back to using the default behavior so if this is
// implemented all subfields
type ExpressionEvaluator interface {
	// FieldConfigurations returns the configuration for this field and any subfields
	// it may have. It must be valid to call this method on nil
	FieldConfigurations() FieldConfigurations

	// EvaluateExpression returns whether there was a match or not. We are not also
	// expecting any errors because all the validation bits are handled
	// during parsing and cross checking against the output of FieldConfigurations.
	EvaluateExpression(sel Selector, op MatchOperator, value interface{}) (bool, error)
}

type Expression struct {
	// The syntax tree
	ast Expr

	// A few configurations for extra validation of the AST
	config ExpressionConfig

	// Once an expression has been run against a particular data type it cannot be executed
	// against a different data type. Some coerced value memoization occurs which would
	// be invalid against other data types.
	boundType reflect.Type

	// The field configuration of the boundType
	fields FieldConfigurations
}

func Create(expression string, config *ExpressionConfig) (*Expression, error) {
	return CreateForType(expression, config, nil)
}

func CreateForType(expression string, config *ExpressionConfig, dataType interface{}) (*Expression, error) {
	ast, err := Parse("", []byte(expression))

	if err != nil {
		return nil, err
	}

	exp := &Expression{ast: ast.(Expr)}

	if config == nil {
		config = &exp.config
	}
	err = exp.validate(config, dataType, true)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (exp *Expression) Evaluate(datum interface{}) (bool, error) {
	if exp.fields == nil {
		err := exp.validate(&exp.config, datum, true)
		if err != nil {
			return false, err
		}
	} else if derefType(reflect.TypeOf(datum)) != exp.boundType {
		return false, fmt.Errorf("This expression can only be used to evaluate matches against %s", exp.boundType)
	}

	return evaluate(exp.ast, datum, exp.fields)
}

func (exp *Expression) validate(config *ExpressionConfig, dataType interface{}, updateExpression bool) error {
	if config == nil {
		return fmt.Errorf("Invalid config")
	}

	var fields FieldConfigurations
	var err error
	var rtype reflect.Type
	if dataType != nil {
		fields, rtype, err = generateFieldConfigurationsAndType(dataType)
		if err != nil {
			return err
		}
	}

	err = validate(exp.ast, fields, config.MaxMatches, config.MaxRawValueLength)
	if err != nil {
		return err
	}

	if updateExpression {
		exp.config = *config
		exp.fields = fields
		exp.boundType = rtype
	}

	return nil
}

// Validates an existing expression against a possibly different configuration
func (exp *Expression) Validate(config *ExpressionConfig, dataType interface{}) error {
	return exp.validate(config, dataType, false)
}
