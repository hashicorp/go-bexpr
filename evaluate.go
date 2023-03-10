package bexpr

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/solo-finance/go-bexpr/grammar"

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

func primitiveGreaterThanFn(kind reflect.Kind) func(first interface{}, second reflect.Value) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return doGreaterThanInt64
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return doGreaterThanUint64
	case reflect.Float32:
		return doGreaterThanFloat32
	case reflect.Float64:
		return doGreaterThanFloat64
	case reflect.String:
		return doGreaterThanString
	default:
		return nil
	}
}

func primitiveLesserThanFn(kind reflect.Kind) func(first interface{}, second reflect.Value) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return doLesserThanInt64
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return doLesserThanUint64
	case reflect.Float32:
		return doLesserThanFloat32
	case reflect.Float64:
		return doLesserThanFloat64
	case reflect.String:
		return doLesserThanString
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

func doGreaterThanInt64(first interface{}, second reflect.Value) bool {
	return first.(int64) > second.Int()
}

func doLesserThanInt64(first interface{}, second reflect.Value) bool {
	return first.(int64) < second.Int()
}

func doEqualUint64(first interface{}, second reflect.Value) bool {
	return first.(uint64) == second.Uint()
}

func doGreaterThanUint64(first interface{}, second reflect.Value) bool {
	return first.(uint64) > second.Uint()
}

func doLesserThanUint64(first interface{}, second reflect.Value) bool {
	return first.(uint64) < second.Uint()
}

func doEqualFloat32(first interface{}, second reflect.Value) bool {
	return first.(float32) == float32(second.Float())
}

func doGreaterThanFloat32(first interface{}, second reflect.Value) bool {
	return first.(float32) > float32(second.Float())
}

func doLesserThanFloat32(first interface{}, second reflect.Value) bool {
	return first.(float32) < float32(second.Float())
}

func doEqualFloat64(first interface{}, second reflect.Value) bool {
	return first.(float64) == second.Float()
}

func doGreaterThanFloat64(first interface{}, second reflect.Value) bool {
	return first.(float64) > second.Float()
}

func doLesserThanFloat64(first interface{}, second reflect.Value) bool {
	return first.(float64) < second.Float()
}

func doEqualString(first interface{}, second reflect.Value) bool {
	dateTimeFirst, errFirst := time.Parse("2020-10-30", first.(string))
	dateTimeSecond, errSecond := time.Parse("2020-10-30", second.String())

	if errFirst != nil || errSecond != nil {
		return first.(string) == second.String()
	}
	return dateTimeFirst.Unix() == dateTimeSecond.Unix()
}

func doGreaterThanString(first interface{}, second reflect.Value) bool {
	dateTimeFirst, errFirst := time.Parse("2020-10-30", first.(string))
	dateTimeSecond, errSecond := time.Parse("2020-10-30", second.String())

	if errFirst != nil || errSecond != nil {
		return first.(string) > second.String()
	}
	return dateTimeFirst.Unix() > dateTimeSecond.Unix()
}

func doLesserThanString(first interface{}, second reflect.Value) bool {
	dateTimeFirst, errFirst := time.Parse("2020-10-30", first.(string))
	dateTimeSecond, errSecond := time.Parse("2020-10-30", second.String())

	if errFirst != nil || errSecond != nil {
		return first.(string) < second.String()
	}
	return dateTimeFirst.Unix() < dateTimeSecond.Unix()
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
	if eqFn == nil {
		return false, errors.New("unable to find suitable primitive comparison function for matching")
	}
	matchValue, err := getMatchExprValue(expression, value.Kind())
	if err != nil {
		return false, fmt.Errorf("error getting match value in expression: %w", err)
	}
	return eqFn(matchValue, value), nil
}

func doMatchGreaterThan(expression *grammar.MatchExpression, value reflect.Value) (bool, error) {
	// NOTE: see preconditions in evaluategrammar.MatchExpressionRecurse
	gthanFn := primitiveGreaterThanFn(value.Kind())
	if gthanFn == nil {
		return false, errors.New("unable to find suitable primitive comparison function for matching")
	}
	matchValue, err := getMatchExprValue(expression, value.Kind())
	if err != nil {
		return false, fmt.Errorf("error getting match value in expression: %w", err)
	}
	return gthanFn(matchValue, value), nil
}

func doMatchGreaterOrEqualThan(expression *grammar.MatchExpression, value reflect.Value) (bool, error) {
	eqFn := primitiveEqualityFn(value.Kind())
	if eqFn == nil {
		return false, errors.New("unable to find suitable primitive comparison function for matching")
	}
	matchValue, err := getMatchExprValue(expression, value.Kind())
	if err != nil {
		return false, fmt.Errorf("error getting match value in expression: %w", err)
	}
	// NOTE: see preconditions in evaluategrammar.MatchExpressionRecurse
	gthanFn := primitiveGreaterThanFn(value.Kind())
	if gthanFn == nil {
		return false, errors.New("unable to find suitable primitive comparison function for matching")
	}
	return eqFn(matchValue, value) || gthanFn(matchValue, value), nil
}

func doMatchLesserThan(expression *grammar.MatchExpression, value reflect.Value) (bool, error) {
	// NOTE: see preconditions in evaluategrammar.MatchExpressionRecurse
	lthanFn := primitiveLesserThanFn(value.Kind())
	if lthanFn == nil {
		return false, errors.New("unable to find suitable primitive comparison function for matching")
	}
	matchValue, err := getMatchExprValue(expression, value.Kind())
	if err != nil {
		return false, fmt.Errorf("error getting match value in expression: %w", err)
	}
	return lthanFn(matchValue, value), nil
}

func doMatchLesserOrEqualThan(expression *grammar.MatchExpression, value reflect.Value) (bool, error) {
	eqFn := primitiveEqualityFn(value.Kind())
	if eqFn == nil {
		return false, errors.New("unable to find suitable primitive comparison function for matching")
	}
	matchValue, err := getMatchExprValue(expression, value.Kind())
	if err != nil {
		return false, fmt.Errorf("error getting match value in expression: %w", err)
	}
	// NOTE: see preconditions in evaluategrammar.MatchExpressionRecurse
	lthanFn := primitiveLesserThanFn(value.Kind())
	if lthanFn == nil {
		return false, errors.New("unable to find suitable primitive comparison function for matching")
	}
	return eqFn(matchValue, value) || lthanFn(matchValue, value), nil
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
		kind := itemType.Kind()
		switch kind {
		case reflect.Interface:
			// If it's an interface, that is, the type was []interface{}, we
			// have to treat each element individually, checking each element's
			// type/kind and rederiving the match value.
			for i := 0; i < value.Len(); i++ {
				item := value.Index(i).Elem()
				itemType := derefType(item.Type())
				kind := itemType.Kind()
				// We need to special case errors here. The reason is that in an
				// interface slice there can be a mix/match of types, but the
				// coerce functions expect a certain type. So the expression
				// passed in might be `"true" in "/my/slice"` but the value it's
				// checking against might be an integer, thus it will try to
				// coerce "true" to an integer and fail. However, all of the
				// functions use strconv which has a specific error type for
				// syntax errors, so as a special case in this situation, don't
				// error on a strconv.ErrSyntax, just continue on to the next
				// element.
				matchValue, err = getMatchExprValue(expression, kind)
				if err != nil {
					if errors.Is(err, strconv.ErrSyntax) {
						continue
					}
					return false, errors.New(`error getting interface slice match value in expression`)
				}
				eqFn := primitiveEqualityFn(kind)
				if eqFn == nil {
					return false, fmt.Errorf(`unable to find suitable primitive comparison function for "in" comparison in interface slice: %s`, kind)
				}
				// the value will be the correct type as we verified the itemType
				if eqFn(matchValue, reflect.Indirect(item)) {
					return true, nil
				}
			}
			return false, nil

		default:
			// Otherwise it's a concrete type and we can essentially cache the
			// answers. First we need to re-derive the match value for equality
			// assertion.
			matchValue, err = getMatchExprValue(expression, kind)
			if err != nil {
				return false, fmt.Errorf("error getting match value in expression: %w", err)
			}
			eqFn := primitiveEqualityFn(kind)
			if eqFn == nil {
				return false, errors.New(`unable to find suitable primitive comparison function for "in" comparison`)
			}
			for i := 0; i < value.Len(); i++ {
				item := value.Index(i)
				// the value will be the correct type as we verified the itemType
				if eqFn(matchValue, reflect.Indirect(item)) {
					return true, nil
				}
			}
			return false, nil
		}

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

func evaluateCollectionExpression(expression *grammar.CollectionExpression, datum interface{}, opt ...Option) (bool, error) {
	opts := getOpts(opt...)
	ptr := pointerstructure.Pointer{
		Parts: expression.Selector.Path,
		Config: pointerstructure.Config{
			TagName:                 opts.withTagName,
			ValueTransformationHook: opts.withHookFn,
		},
	}
	val, err := ptr.Get(datum)
	if err != nil {
		if errors.Is(err, pointerstructure.ErrNotFound) && opts.withUnknown != nil {
			err = nil
			val = *opts.withUnknown
		}

		if err != nil {
			return false, fmt.Errorf("error finding value in datum: %w", err)
		}
	}
	fmt.Println(opts)
	buf := new(bytes.Buffer)
	expression.Expression.ExpressionDump(buf, " ", 0)
	fmt.Println(buf.String())
	fmt.Println(val)
	rvalue := reflect.Indirect(reflect.ValueOf(val))
	switch expression.Operator {
	case grammar.CollectionOpAll:
		for i := 0; i < rvalue.Len(); i++ {
			sliceValue := rvalue.Index(i)
			result, _ := evaluate(expression.Expression, sliceValue, opt...)
			if !result && expression.Operator == grammar.CollectionOpAny {
				return false, nil
			}
		}
		return true, nil
	case grammar.CollectionOpAny:
		for i := 0; i < rvalue.Len(); i++ {
			sliceValue := rvalue.Index(i)
			result, _ := evaluate(expression.Expression, sliceValue, opt...)
			if result && expression.Operator == grammar.CollectionOpAny {
				return true, nil
			}
		}
		return false, nil

	}
	return false, nil

}

func evaluateMatchExpression(expression *grammar.MatchExpression, datum interface{}, opt ...Option) (bool, error) {
	opts := getOpts(opt...)
	ptr := pointerstructure.Pointer{
		Parts: expression.Selector.Path,
		Config: pointerstructure.Config{
			TagName:                 opts.withTagName,
			ValueTransformationHook: opts.withHookFn,
		},
	}
	val, err := ptr.Get(datum)
	if err != nil {
		if errors.Is(err, pointerstructure.ErrNotFound) && opts.withUnknown != nil {
			err = nil
			val = *opts.withUnknown
		}

		if err != nil {
			return false, fmt.Errorf("error finding value in datum: %w", err)
		}
	}

	rvalue := reflect.Indirect(reflect.ValueOf(val))
	switch expression.Operator {
	case grammar.MatchGreaterThan:
		return doMatchGreaterThan(expression, rvalue)
	case grammar.MatchGreaterOrEqualThan:
		return doMatchGreaterOrEqualThan(expression, rvalue)
	case grammar.MatchLesserThan:
		return doMatchLesserThan(expression, rvalue)
	case grammar.MatchLesserOrEqualThan:
		return doMatchLesserOrEqualThan(expression, rvalue)
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

func evaluate(ast grammar.Expression, datum interface{}, opt ...Option) (bool, error) {

	switch node := ast.(type) {
	case *grammar.UnaryExpression:
		switch node.Operator {
		case grammar.UnaryOpNot:
			result, err := evaluate(node.Operand, datum, opt...)
			return !result, err
		}
	case *grammar.BinaryExpression:
		switch node.Operator {
		case grammar.BinaryOpAnd:
			result, err := evaluate(node.Left, datum, opt...)
			if err != nil || !result {
				return result, err
			}

			return evaluate(node.Right, datum, opt...)

		case grammar.BinaryOpOr:
			result, err := evaluate(node.Left, datum, opt...)
			if err != nil || result {
				return result, err
			}

			return evaluate(node.Right, datum, opt...)
		}
	case *grammar.MatchExpression:
		return evaluateMatchExpression(node, datum, opt...)
	case *grammar.CollectionExpression:
		return evaluateCollectionExpression(node, datum, opt...)
	}

	return false, fmt.Errorf("Invalid AST node")
}
