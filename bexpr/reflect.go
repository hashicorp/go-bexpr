package bexpr

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

var primitiveCoercionFns = map[reflect.Kind]FieldValueCoercionFn{
	reflect.Bool:    CoerceBool,
	reflect.Int:     CoerceInt,
	reflect.Int8:    CoerceInt8,
	reflect.Int16:   CoerceInt16,
	reflect.Int32:   CoerceInt32,
	reflect.Int64:   CoerceInt64,
	reflect.Uint:    CoerceUint,
	reflect.Uint8:   CoerceUint8,
	reflect.Uint16:  CoerceUint16,
	reflect.Uint32:  CoerceUint32,
	reflect.Uint64:  CoerceUint64,
	reflect.Float32: CoerceFloat32,
	reflect.Float64: CoerceFloat64,
	reflect.String:  nil,
}

var primitiveEqualityFns = map[reflect.Kind]func(first interface{}, second interface{}) bool{
	reflect.Bool:    reflectEqualBool,
	reflect.Int:     reflectEqualInt,
	reflect.Int8:    reflectEqualInt8,
	reflect.Int16:   reflectEqualInt16,
	reflect.Int32:   reflectEqualInt32,
	reflect.Int64:   reflectEqualInt64,
	reflect.Uint:    reflectEqualUint,
	reflect.Uint8:   reflectEqualUint8,
	reflect.Uint16:  reflectEqualUint16,
	reflect.Uint32:  reflectEqualUint32,
	reflect.Uint64:  reflectEqualUint64,
	reflect.Float32: reflectEqualFloat32,
	reflect.Float64: reflectEqualFloat64,
	reflect.String:  reflectEqualString,
}

func reflectEqualBool(first interface{}, second interface{}) bool {
	return first.(bool) == second.(bool)
}

func reflectEqualInt(first interface{}, second interface{}) bool {
	return first.(int) == second.(int)
}

func reflectEqualInt8(first interface{}, second interface{}) bool {
	return first.(int8) == second.(int8)
}

func reflectEqualInt16(first interface{}, second interface{}) bool {
	return first.(int16) == second.(int16)
}

func reflectEqualInt32(first interface{}, second interface{}) bool {
	return first.(int32) == second.(int32)
}

func reflectEqualInt64(first interface{}, second interface{}) bool {
	return first.(int64) == second.(int64)
}

func reflectEqualUint(first interface{}, second interface{}) bool {
	return first.(uint) == second.(uint)
}

func reflectEqualUint8(first interface{}, second interface{}) bool {
	return first.(uint8) == second.(uint8)
}

func reflectEqualUint16(first interface{}, second interface{}) bool {
	return first.(uint16) == second.(uint16)
}

func reflectEqualUint32(first interface{}, second interface{}) bool {
	return first.(uint32) == second.(uint32)
}

func reflectEqualUint64(first interface{}, second interface{}) bool {
	return first.(uint64) == second.(uint64)
}

func reflectEqualFloat32(first interface{}, second interface{}) bool {
	return first.(float32) == second.(float32)
}

func reflectEqualFloat64(first interface{}, second interface{}) bool {
	return first.(float64) == second.(float64)
}

func reflectEqualString(first interface{}, second interface{}) bool {
	return first.(string) == second.(string)
}

// Get rid of 0 to many levels of pointers to get at the real type
func derefType(rtype reflect.Type) reflect.Type {
	for rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}
	return rtype
}

// Supported Primitive Types:
//    string
//    integers (all width types and signedness)
//    floats
//    bools
//
//    Each one of these supports the "==" and "!=" checks"
//
// Supported Compound Types
//   map[*]*
//       - is empty, is not empty
//   map[string]*
//       - in, not in, contains, not contains to check keys
//   map[string]<supported type>
//       - Will have a single subfield with name "" (wildcard) and the rest of
//         the field configuration will come from the <supported type>
//   array or slice of any type
//   array or slice of <supported primitive type>
//       -  in, not in, contains, not contains to check against the values
//   array or slice of <supported compound type>
//       - Will have subfields with configurations of whatever the supported
//         compound type is
//
//   structs
//       - no supported ops on it (just used for selecting the final value)
//       - sub fields for each exported field of the struct that is a supported
//         type (primitive and complex)
//
func reflectFieldConfigurationInternal(name string, rtype reflect.Type) *FieldConfiguration {

	// handle the primitive types
	if coerceFn, ok := primitiveCoercionFns[rtype.Kind()]; ok {
		return &FieldConfiguration{
			Name:                name,
			CoerceFn:            coerceFn,
			SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual},
		}
	}

	switch rtype.Kind() {
	case reflect.Map:
		switch derefType(rtype.Key()).Kind() {
		case reflect.String:
			cfg := &FieldConfiguration{
				Name:                name,
				SupportedOperations: []MatchOperator{MatchIsEmpty, MatchIsNotEmpty, MatchIn, MatchNotIn},
			}

			subfield := reflectFieldConfigurationInternal("", derefType(rtype.Elem()))

			if subfield != nil {
				cfg.SubFields = []*FieldConfiguration{subfield}
			}

			return cfg
		default:
			// For maps with non string keys we can really only do emptiness checks
			// and cannot index into them at all.
			return &FieldConfiguration{
				Name:                name,
				SupportedOperations: []MatchOperator{MatchIsEmpty, MatchIsNotEmpty},
			}
		}

	case reflect.Array, reflect.Slice:
		elemType := derefType(rtype.Elem())

		if coerceFn, ok := primitiveCoercionFns[elemType.Kind()]; ok {
			return &FieldConfiguration{
				Name:                name,
				CoerceFn:            coerceFn,
				SupportedOperations: []MatchOperator{MatchIn, MatchNotIn, MatchIsEmpty, MatchIsNotEmpty},
			}
		}

		// TODO (mkeeler) - could maybe make this a little cleaner
		subfields := reflectFieldConfigurationInternal("fake", elemType)

		// TODO (mkeeler) - this wont support a slice/array of other non-struct containers
		// This shows a bit of a deficiency with the original syntax in not being able to select
		// each element of the slice by an index and perform an operation against it.
		// We really would need an any/all syntax or we could introduce injected selectors
		//
		// foo.bar is empty would check if the foo.bar slice is empty and then
		// foo.bar.any is empty would check if any elements are empty.
		// I need to think about this syntax a bit more.
		return &FieldConfiguration{
			Name:                name,
			SubFields:           subfields.SubFields,
			SupportedOperations: []MatchOperator{MatchIsEmpty, MatchIsNotEmpty},
		}

	case reflect.Struct:
		cfg := &FieldConfiguration{
			Name: name,
		}
		for i, maxFields := 0, rtype.NumField(); i < maxFields; i++ {
			field := rtype.Field(i)

			if field.PkgPath != "" {
				continue
			}

			subfield := reflectFieldConfigurationInternal(field.Name, derefType(field.Type))
			cfg.SubFields = append(cfg.SubFields, subfield)
		}

		return cfg

	default:
		return nil
	}

}

func ReflectFieldConfigurations(topLevelType interface{}) ([]*FieldConfiguration, error) {
	rtype := derefType(reflect.TypeOf(topLevelType))

	switch rtype.Kind() {
	case reflect.Struct:
		var fields []*FieldConfiguration
		for i, maxFields := 0, rtype.NumField(); i < maxFields; i++ {
			field := rtype.Field(i)

			if field.PkgPath != "" {
				continue
			}

			cfg := reflectFieldConfigurationInternal(field.Name, derefType(field.Type))
			if cfg != nil {
				fields = append(fields, cfg)
			}
		}

		return fields, nil
	case reflect.Map:
		if rtype.Key().Kind() != reflect.String {
			return nil, fmt.Errorf("Cannot generate FieldConfigurations for maps with keys that are not strings")
		}

		elemType := derefType(rtype.Elem())

		field := reflectFieldConfigurationInternal("", elemType)
		if field != nil {
			return []*FieldConfiguration{field}, nil
		}

		return nil, nil
	}

	return nil, fmt.Errorf("Invalid top level type - can only use structs or map[string]*")
}

func getMatcherValue(matcher *MatchExpr, kind reflect.Kind) (interface{}, error) {
	if matcher.Value == nil {
		return nil, fmt.Errorf("Matching expression for selector %q has no value", matcher.Selector)
	}

	if coerceFn, ok := primitiveCoercionFns[kind]; ok {
		if coerceFn != nil && matcher.Value.Converted == nil {
			conv, err := coerceFn(matcher.Value.Raw)

			if err != nil {
				return nil, fmt.Errorf("Failed to convert value to type %s for selector %q - %v", kind, matcher.Selector, err)
			}

			matcher.Value.Converted = conv
			return conv, nil
		} else if matcher.Value.Converted == nil {
			return matcher.Value.Raw, nil
		} else if reflect.TypeOf(matcher.Value.Converted).Kind() != kind {
			return nil, fmt.Errorf("Invalid converted value stored for expression with selector %q", matcher.Selector)
		}

		return matcher.Value.Converted, nil
	}

	return nil, fmt.Errorf("Invalid non-primitive value type %s exists for selector %q", kind, matcher.Selector)
}

func reflectMatchEqual(matcher *MatchExpr, value reflect.Value) (bool, error) {
	if eqFn, ok := primitiveEqualityFns[value.Kind()]; ok {
		matchValue, err := getMatcherValue(matcher, value.Kind())
		if err != nil {
			return false, fmt.Errorf("Failed to execute equality match: %v", err)
		}

		return eqFn(matchValue, value.Interface()), nil
	}
	return false, fmt.Errorf("Cannot use non-primitive value types for equality operations for selector: %q", matcher.Selector)
}

func reflectMatchIn(matcher *MatchExpr, value reflect.Value) (bool, error) {
	switch kind := value.Kind(); kind {
	case reflect.Map:
		keyType := derefType(value.Type().Key())

		matchValue, err := getMatcherValue(matcher, keyType.Kind())
		if err != nil {
			return false, fmt.Errorf("Failed to execute in/contains match: %v", err)
		}

		found := value.MapIndex(reflect.ValueOf(matchValue))

		return found.IsValid(), nil
	case reflect.Slice, reflect.Array:
		itemType := derefType(value.Type().Elem())

		if eqFn, ok := primitiveEqualityFns[itemType.Kind()]; ok {
			matchValue, err := getMatcherValue(matcher, itemType.Kind())
			if err != nil {
				return false, fmt.Errorf("Failed to execute in/contains match: %v", err)
			}

			for i := 0; i < value.Len(); i++ {
				item := value.Index(i)

				// the value will be the correct type as we verified the itemType
				if eqFn(matchValue, reflect.Indirect(item).Interface()) {
					return true, nil
				}
			}

			return false, nil
		}

		return false, fmt.Errorf("Cannot use non-primitive value types for in/contains operations for selector: %q", matcher.Selector)
	case reflect.String:
		// getMatcherValue will ensure we get a string back or error
		matchValue, err := getMatcherValue(matcher, value.Kind())
		if err != nil {
			return false, fmt.Errorf("Failed to execute in/contains match: %v", err)
		}
		return strings.Contains(value.String(), matchValue.(string)), nil
	default:
		return false, fmt.Errorf("Cannot perform in/contains operations on type %s for selector: %q", kind, matcher.Selector)
	}
}

func reflectMatchIsEmpty(matcher *MatchExpr, value reflect.Value) (bool, error) {
	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		return value.Len() == 0, nil
	default:
		return false, fmt.Errorf("Cannot use value of type %s for emptiness checks", value.Kind())
	}
}

func reflectEvaluateMatcherRecurse(matcher *MatchExpr, depth int, value reflect.Value) (bool, error) {
	kind := value.Kind()
	if depth >= len(matcher.Selector) {
		// we have reached the end of the selector - execute the match operations
		switch matcher.Operator {
		case MatchEqual:
			return reflectMatchEqual(matcher, value)
		case MatchNotEqual:
			result, err := reflectMatchEqual(matcher, value)
			if err == nil {
				return !result, nil
			}
			return false, err
		case MatchIn:
			return reflectMatchIn(matcher, value)
		case MatchNotIn:
			result, err := reflectMatchIn(matcher, value)
			if err == nil {
				return !result, nil
			}
			return false, err
		case MatchIsEmpty:
			return reflectMatchIsEmpty(matcher, value)
		case MatchIsNotEmpty:
			result, err := reflectMatchIsEmpty(matcher, value)
			if err == nil {
				return !result, nil
			}
			return false, err
		default:
			return false, fmt.Errorf("Invalid match operation: %d", matcher.Operator)
		}
	} else {
		// still more selectors to go through

		switch kind {
		case reflect.Struct:
			sel := []rune(matcher.Selector[depth])

			if len(sel) == 0 {
				return false, fmt.Errorf("Invalid selector with empty field: %q", matcher.Selector[:depth+1])
			}

			// disallow selecting through unexported fields
			if unicode.IsLower(sel[0]) {
				return false, fmt.Errorf("Invalid selector: %q", matcher.Selector[:depth+1])
			}

			value = reflect.Indirect(value.FieldByName(matcher.Selector[depth]))

			if !value.IsValid() {
				return false, fmt.Errorf("Invalid selector: %q", matcher.Selector[:depth+1])
			}

			return reflectEvaluateMatcherRecurse(matcher, depth+1, value)
		case reflect.Slice, reflect.Array:
			// TODO (mkeeler) - maybe check the elem type here and report an error if unable to be used
			for i := 0; i < value.Len(); i++ {
				item := reflect.Indirect(value.Index(i))
				// we use the same depth because right now we are not allowing
				// selection of individual slice/array elements
				result, err := reflectEvaluateMatcherRecurse(matcher, depth, item)
				if err != nil {
					return false, err
				}

				if result {
					return true, nil
				}
			}
			return false, nil
		case reflect.Map:
			keyType := derefType(value.Type().Key())

			if keyType.Kind() != reflect.String {
				return false, fmt.Errorf("Invalid map key type for selector: %q - %s", matcher.Selector[:depth+1], keyType.Kind())
			}

			value = reflect.Indirect(value.MapIndex(reflect.ValueOf(matcher.Selector[depth])))

			if !value.IsValid() {
				// TODO (mkeeler) should this error? Probably not but it would require more logic
				// to determine whether it should be true or false depending on the operator
				return false, fmt.Errorf("Invalid selector - %s: key not found in map", matcher.Selector[:depth+1])
			}

			return reflectEvaluateMatcherRecurse(matcher, depth+1, value)

		default:
			return false, fmt.Errorf("Value at selector %q with type %s does not support nested field selection", matcher.Selector[:depth], kind)
		}
	}
}

func reflectEvaluateMatcher(matcher *MatchExpr, datum interface{}) (bool, error) {
	value := reflect.Indirect(reflect.ValueOf(datum))

	return reflectEvaluateMatcherRecurse(matcher, 0, value)

}
