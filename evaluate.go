package bexpr

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/mitchellh/pointerstructure"
)

var byteSliceTyp reflect.Type = reflect.TypeOf([]byte{})

var primitiveEqualityFns = map[reflect.Kind]func(first interface{}, second reflect.Value) bool{
	reflect.Bool:    doEqualBool,
	reflect.Int:     doEqualInt,
	reflect.Int8:    doEqualInt8,
	reflect.Int16:   doEqualInt16,
	reflect.Int32:   doEqualInt32,
	reflect.Int64:   doEqualInt64,
	reflect.Uint:    doEqualUint,
	reflect.Uint8:   doEqualUint8,
	reflect.Uint16:  doEqualUint16,
	reflect.Uint32:  doEqualUint32,
	reflect.Uint64:  doEqualUint64,
	reflect.Float32: doEqualFloat32,
	reflect.Float64: doEqualFloat64,
	reflect.String:  doEqualString,
}

func doEqualBool(first interface{}, second reflect.Value) bool {
	return first.(bool) == second.Bool()
}

func doEqualInt(first interface{}, second reflect.Value) bool {
	return first.(int) == int(second.Int())
}

func doEqualInt8(first interface{}, second reflect.Value) bool {
	return first.(int8) == int8(second.Int())
}

func doEqualInt16(first interface{}, second reflect.Value) bool {
	return first.(int16) == int16(second.Int())
}

func doEqualInt32(first interface{}, second reflect.Value) bool {
	return first.(int32) == int32(second.Int())
}

func doEqualInt64(first interface{}, second reflect.Value) bool {
	return first.(int64) == second.Int()
}

func doEqualUint(first interface{}, second reflect.Value) bool {
	return first.(uint) == uint(second.Uint())
}

func doEqualUint8(first interface{}, second reflect.Value) bool {
	return first.(uint8) == uint8(second.Uint())
}

func doEqualUint16(first interface{}, second reflect.Value) bool {
	return first.(uint16) == uint16(second.Uint())
}

func doEqualUint32(first interface{}, second reflect.Value) bool {
	return first.(uint32) == uint32(second.Uint())
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

func doMatchMatches(expression *MatchExpression, value reflect.Value) (bool, error) {
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

func doMatchEqual(expression *MatchExpression, value reflect.Value) (bool, error) {
	// NOTE: see preconditions in evaluateMatchExpressionRecurse
	eqFn := primitiveEqualityFns[value.Kind()]
	matchValue, err := getMatchExprValue(expression, value)
	if err != nil {
		return false, fmt.Errorf("error getting match value in expression: %w", err)
	}
	log.Println(fmt.Sprintf("matchValue type %T, value %v", matchValue, matchValue))
	log.Println(fmt.Sprintf("value kind %s", value.Kind()))
	return eqFn(matchValue, value), nil
}

func doMatchIn(expression *MatchExpression, value reflect.Value) (bool, error) {
	// NOTE: see preconditions in evaluateMatchExpressionRecurse
	matchValue, err := getMatchExprValue(expression, value)
	if err != nil {
		return false, fmt.Errorf("error getting match value in expression: %w", err)
	}

	switch kind := value.Kind(); kind {
	case reflect.Map:
		found := value.MapIndex(reflect.ValueOf(matchValue))
		return found.IsValid(), nil
	case reflect.Slice, reflect.Array:
		itemType := derefType(value.Type().Elem())
		eqFn := primitiveEqualityFns[itemType.Kind()]

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
		// this shouldn't be possible but we have to have something to return to keep the compiler happy
		return false, fmt.Errorf("Cannot perform in/contains operations on type %s for selector: %q", kind, expression.Selector)
	}
}

func doMatchIsEmpty(matcher *MatchExpression, value reflect.Value) (bool, error) {
	// NOTE: see preconditions in evaluateMatchExpressionRecurse
	return value.Len() == 0, nil
}

func getMatchExprValue(expression *MatchExpression, rvalue reflect.Value) (interface{}, error) {
	// NOTE: see preconditions in evaluateMatchExpressionRecurse
	if expression.Value == nil {
		return nil, nil
	}

	switch rvalue.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(expression.Value.Raw)
		if err != nil {
			return nil, fmt.Errorf("error parsing value as bool: %w", err)
		}
		return b, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f, err := strconv.ParseInt(expression.Value.Raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing value as int: %w", err)
		}
		return f, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f, err := strconv.ParseUint(expression.Value.Raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing value as uint: %w", err)
		}
		return f, nil
	case reflect.Float32:
		f, err := strconv.ParseFloat(expression.Value.Raw, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing value as float32: %w", err)
		}
		return float32(f), nil
	case reflect.Float64:
		f, err := strconv.ParseFloat(expression.Value.Raw, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing value as float64: %w", err)
		}
		return f, nil
	case reflect.String:
		return expression.Value.Raw, nil
	}
	return expression.Value.Raw, nil
}

func evaluateMatchExpression(expression *MatchExpression, datum interface{}) (bool, error) {
	path := fmt.Sprintf("/%s", strings.Join(expression.Selector.Path, "/"))
	ptr, err := pointerstructure.Parse(path)
	if err != nil {
		return false, fmt.Errorf("error parsing path: %w", err)
	}
	ptr.Config.TagName = "bexpr"
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

	pvalue := reflect.Indirect(reflect.ValueOf(val))
	switch pvalue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = pvalue.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = pvalue.Uint()
	case reflect.Float32:
		val = float32(pvalue.Float())
	case reflect.Float64:
		val = pvalue.Float()
	}

	log.Println(fmt.Sprintf("before rvalue %T, %v", val, val))

	rvalue := reflect.Indirect(reflect.ValueOf(val))
	log.Println(fmt.Sprintf("rvalue %s, %v", rvalue.Kind(), rvalue))
	switch expression.Operator {
	case MatchEqual:
		return doMatchEqual(expression, rvalue)
	case MatchNotEqual:
		result, err := doMatchEqual(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	case MatchIn:
		return doMatchIn(expression, rvalue)
	case MatchNotIn:
		result, err := doMatchIn(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	case MatchIsEmpty:
		return doMatchIsEmpty(expression, rvalue)
	case MatchIsNotEmpty:
		result, err := doMatchIsEmpty(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	case MatchMatches:
		return doMatchMatches(expression, rvalue)
	case MatchNotMatches:
		result, err := doMatchMatches(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	default:
		return false, fmt.Errorf("Invalid match operation: %d", expression.Operator)
	}
}

func evaluate(ast Expression, datum interface{}) (bool, error) {
	switch node := ast.(type) {
	case *UnaryExpression:
		switch node.Operator {
		case UnaryOpNot:
			result, err := evaluate(node.Operand, datum)
			return !result, err
		}
	case *BinaryExpression:
		switch node.Operator {
		case BinaryOpAnd:
			result, err := evaluate(node.Left, datum)
			if err != nil || !result {
				return result, err
			}

			return evaluate(node.Right, datum)

		case BinaryOpOr:
			result, err := evaluate(node.Left, datum)
			if err != nil || result {
				return result, err
			}

			return evaluate(node.Right, datum)
		}
	case *MatchExpression:
		return evaluateMatchExpression(node, datum)
	}
	return false, fmt.Errorf("Invalid AST node")
}
