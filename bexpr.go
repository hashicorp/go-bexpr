// bexpr is an implementation of a generic boolean expression evaluator.
// The general goal is to be able to evaluate some expression against some
// arbitrary data and get back a boolean of whether or not the data
// was matched by the expression
package bexpr

import (
	"fmt"
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
	return CreateEvaluatorForType(expression, config, nil)
}

func CreateEvaluatorForType(expression string, config *EvaluatorConfig, dataType interface{}) (*Evaluator, error) {
	ast, err := Parse("", []byte(expression))

	if err != nil {
		return nil, err
	}

	eval := &Evaluator{ast: ast.(Expression)}

	if config == nil {
		config = &eval.config
	}
	err = eval.validate(config, dataType, true)
	if err != nil {
		return nil, err
	}

	return eval, nil
}

func (eval *Evaluator) Evaluate(datum interface{}) (bool, error) {
	if eval.fields == nil {
		err := eval.validate(&eval.config, datum, true)
		if err != nil {
			return false, err
		}
	} else if reflect.TypeOf(datum) != eval.boundType {
		return false, fmt.Errorf("This evaluator can only be used to evaluate matches against %s", eval.boundType)
	}

	return evaluate(eval.ast, datum, eval.fields)
}

func (eval *Evaluator) validate(config *EvaluatorConfig, dataType interface{}, updateEvaluator bool) error {
	if config == nil {
		return fmt.Errorf("Invalid config")
	}

	var fields FieldConfigurations
	var err error
	var rtype reflect.Type
	if dataType != nil {
		registry := DefaultRegistry

		switch t := dataType.(type) {
		case reflect.Type:
			rtype = t
		case *reflect.Type:
			rtype = *t
		case reflect.Value:
			rtype = t.Type()
		case *reflect.Value:
			rtype = t.Type()
		default:
			rtype = reflect.TypeOf(dataType)
		}

		fields, err = registry.GetFieldConfigurations(rtype)
		if err != nil {
			return err
		}

		if len(fields) < 1 {
			return fmt.Errorf("Data type %s has no evaluatable fields", rtype.String())
		}
	}

	maxMatches := config.MaxMatches
	if maxMatches == 0 {
		maxMatches = defaultMaxMatches
	}

	maxRawValueLength := config.MaxRawValueLength
	if maxRawValueLength == 0 {
		maxRawValueLength = defaultMaxRawValueLength
	}

	err = validate(eval.ast, fields, maxMatches, maxRawValueLength)
	if err != nil {
		return err
	}

	if updateEvaluator {
		eval.config = *config
		eval.fields = fields
		eval.boundType = rtype
	}

	return nil
}

// Validates an existing expression against a possibly different configuration
func (eval *Evaluator) Validate(config *EvaluatorConfig, dataType interface{}) error {
	return eval.validate(config, dataType, false)
}
