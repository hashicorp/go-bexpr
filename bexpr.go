// bexpr is an implementation of a generic boolean expression evaluator.
// The general goal is to be able to evaluate some expression against some
// arbitrary data and get back a boolean of whether or not the data
// was matched by the expression
package bexpr

import (
	"reflect"
)

const (
	defaultMaxMatches        = 32
	defaultMaxRawValueLength = 512
)

type Evaluator struct {
	// The syntax tree
	ast Expression

	// A few configurations for extra validation of the AST
	config EvaluatorConfig

	// Once an expression has been run against a particular data type it cannot be executed
	// against a different data type. Some coerced value memoization occurs which would
	// be invalid against other data types.
	boundType reflect.Type

	// The field configuration of the boundType
	fields FieldConfigurations
}

// Extra configuration used to perform further validation on a parsed
// expression and to aid in the evaluation process
type EvaluatorConfig struct {
	// Maximum number of matching expressions allowed. 0 means unlimited
	// This does not include and, or and not expressions within the AST
	MaxMatches int
	// Maximum length of raw values. 0 means unlimited
	MaxRawValueLength int
}

func CreateEvaluator(expression string, config *EvaluatorConfig) (*Evaluator, error) {
	ast, err := Parse("", []byte(expression))
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = &EvaluatorConfig{}
	}

	eval := &Evaluator{
		ast:    ast.(Expression),
		config: *config,
	}

	return eval, nil
}

func (eval *Evaluator) Evaluate(datum interface{}) (bool, error) {
	return evaluate(eval.ast, datum)
}
