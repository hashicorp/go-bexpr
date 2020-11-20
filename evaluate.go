package bexpr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/go-bexpr/grammar"
	"github.com/mitchellh/pointerstructure"
)

var byteSliceTyp reflect.Type = reflect.TypeOf([]byte{})

func primitiveEqualityFn(kind reflect.Kind) func(first interface{}, second reflect.Value) bool {
	switch kind {
	case reflect.Bool:
		return doEqualBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return doEqualInt64
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return doEqualUint64
	case reflect.Float32:
		return doEqualFloat32
	case reflect.Float64:
		return doEqualFloat64
	case reflect.String:
		return doEqualString
	default:
		return nil
	}
}

func doEqualBool(first interface{}, second reflect.Value) bool {
	return first.(bool) == second.Bool()
}

func doEqualInt64(first interface{}, second reflect.Value) bool {
	return first.(int64) == second.Int()
}

func doEqualUint64(first interface{}, second reflect.Value) bool {
	return first.(uint64) == second.Uint()
}

func doEqualFloat32(first interface{}, second reflect.Value) bool {
	return first.(float32) == float32(second.Float())
}

func doEqualFloat64(first interface{}, second reflect.Value) bool {
	return first.(float64) == second.Float()
}

func doEqualString(first interface{}, second reflect.Value) bool {
	return first.(string) == second.String()
}

// Get rid of 0 to many levels of pointers to get at the real type
func derefType(rtype reflect.Type) reflect.Type {
	for rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}
	return rtype
}

func doMatchMatches(expression *grammar.MatchExpression, value reflect.Value) (bool, error) {
	if !value.Type().ConvertibleTo(byteSliceTyp) {
		return false, fmt.Errorf("Value of type %s is not convertible to []byte", value.Type())
	}

	var re *regexp.Regexp
	var ok bool
	if expression.Value.Converted != nil {
		re, ok = expression.Value.Converted.(*regexp.Regexp)
	}
	if !ok || re == nil {
		var err error
		re, err = regexp.Compile(expression.Value.Raw)
		if err != nil {
			return false, fmt.Errorf("Failed to compile regular expression %q: %v", expression.Value.Raw, err)
		}
		expression.Value.Converted = re
	}

	return re.Match(value.Convert(byteSliceTyp).Interface().([]byte)), nil
}

func doMatchEqual(expression *grammar.MatchExpression, value reflect.Value) (bool, error) {
	// NOTE: see preconditions in evaluategrammar.MatchExpressionRecurse
	eqFn := primitiveEqualityFn(value.Kind())
	matchValue, err := getMatchExprValue(expression, value.Kind())
	if err != nil {
		return false, fmt.Errorf("error getting match value in expression: %w", err)
	}
	return eqFn(matchValue, value), nil
}

func doMatchIn(expression *grammar.MatchExpression, value reflect.Value) (bool, error) {
	matchValue, err := getMatchExprValue(expression, value.Kind())
	if err != nil {
		return false, fmt.Errorf("error getting match value in expression: %w", err)
	}

	switch kind := value.Kind(); kind {
	case reflect.Map:
		found := value.MapIndex(reflect.ValueOf(matchValue))
		return found.IsValid(), nil

	case reflect.Slice, reflect.Array:
		itemType := derefType(value.Type().Elem())
		// Once we know the item type, we need to re-derive the match value for
		// equality assertion
		matchValue, err = getMatchExprValue(expression, itemType.Kind())
		if err != nil {
			return false, fmt.Errorf("error getting match value in expression: %w", err)
		}
		eqFn := primitiveEqualityFn(itemType.Kind())

		for i := 0; i < value.Len(); i++ {
			item := value.Index(i)

			// the value will be the correct type as we verified the itemType
			if eqFn(matchValue, reflect.Indirect(item)) {
				return true, nil
			}
		}

		return false, nil

	case reflect.String:
		return strings.Contains(value.String(), matchValue.(string)), nil

	default:
		return false, fmt.Errorf("Cannot perform in/contains operations on type %s for selector: %q", kind, expression.Selector)
	}
}

func doMatchIsEmpty(matcher *grammar.MatchExpression, value reflect.Value) (bool, error) {
	// NOTE: see preconditions in evaluategrammar.MatchExpressionRecurse
	return value.Len() == 0, nil
}

func getMatchExprValue(expression *grammar.MatchExpression, rvalue reflect.Kind) (interface{}, error) {
	if expression.Value == nil {
		return nil, nil
	}

	switch rvalue {
	case reflect.Bool:
		return CoerceBool(expression.Value.Raw)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return CoerceInt64(expression.Value.Raw)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return CoerceUint64(expression.Value.Raw)

	case reflect.Float32:
		return CoerceFloat32(expression.Value.Raw)

	case reflect.Float64:
		return CoerceFloat64(expression.Value.Raw)

	default:
		return expression.Value.Raw, nil
	}
}

func evaluateMatchExpression(expression *grammar.MatchExpression, datum interface{}) (bool, error) {
	ptr := pointerstructure.Pointer{
		Parts: expression.Selector.Path,
		Config: pointerstructure.Config{
			TagName: "bexpr",
		},
	}
	val, err := ptr.Get(datum)
	if err != nil {
		return false, fmt.Errorf("error finding value in datum: %w", err)
	}

	if jn, ok := val.(json.Number); ok {
		if jni, err := jn.Int64(); err == nil {
			val = jni
		} else if jnf, err := jn.Float64(); err == nil {
			val = jnf
		} else {
			return false, fmt.Errorf("unable to convert json number %s to int or float", jn)
		}
	}

	rvalue := reflect.Indirect(reflect.ValueOf(val))
	switch expression.Operator {
	case grammar.MatchEqual:
		return doMatchEqual(expression, rvalue)
	case grammar.MatchNotEqual:
		result, err := doMatchEqual(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	case grammar.MatchIn:
		return doMatchIn(expression, rvalue)
	case grammar.MatchNotIn:
		result, err := doMatchIn(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	case grammar.MatchIsEmpty:
		return doMatchIsEmpty(expression, rvalue)
	case grammar.MatchIsNotEmpty:
		result, err := doMatchIsEmpty(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	case grammar.MatchMatches:
		return doMatchMatches(expression, rvalue)
	case grammar.MatchNotMatches:
		result, err := doMatchMatches(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	default:
		return false, fmt.Errorf("Invalid match operation: %d", expression.Operator)
	}
}

func evaluate(ast grammar.Expression, datum interface{}) (bool, error) {
	switch node := ast.(type) {
	case *grammar.UnaryExpression:
		switch node.Operator {
		case grammar.UnaryOpNot:
			result, err := evaluate(node.Operand, datum)
			return !result, err
		}
	case *grammar.BinaryExpression:
		switch node.Operator {
		case grammar.BinaryOpAnd:
			result, err := evaluate(node.Left, datum)
			if err != nil || !result {
				return result, err
			}

			return evaluate(node.Right, datum)

		case grammar.BinaryOpOr:
			result, err := evaluate(node.Left, datum)
			if err != nil || result {
				return result, err
			}

			return evaluate(node.Right, datum)
		}
	case *grammar.MatchExpression:
		return evaluateMatchExpression(node, datum)
	}
	return false, fmt.Errorf("Invalid AST node")
}
