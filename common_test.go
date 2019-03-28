package bexpr

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var benchFull *bool = flag.Bool("bench-full", false, "Run all benchmarks rather than a subset")

func FullBenchmarks() bool {
	return benchFull != nil && *benchFull
}

type testFlatStruct struct {
	Int        int
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Float32    float32
	Float64    float64
	Bool       bool
	String     string
	unexported string
	Hidden     bool `bexpr:"-"`
}

var testFlatStructKindMap map[string]reflect.Kind = map[string]reflect.Kind{
	"Int":     reflect.Int,
	"Int8":    reflect.Int8,
	"Int16":   reflect.Int16,
	"Int32":   reflect.Int32,
	"Int64":   reflect.Int64,
	"Uint":    reflect.Uint,
	"Uint8":   reflect.Uint8,
	"Uint16":  reflect.Uint16,
	"Uint32":  reflect.Uint32,
	"Uint64":  reflect.Uint64,
	"Float32": reflect.Float32,
	"Float64": reflect.Float64,
	"Bool":    reflect.Bool,
	"String":  reflect.String,
}

type testNestedLevel2_1 struct {
	Foo int
	bar string
	Baz string
}

type testNestedLevel2_2 struct {
	X int
	Y int
	z int
}

type testNestedLevel1 struct {
	Map              map[string]string
	MapOfStructs     map[string]testNestedLevel2_1
	MapInfInf        map[interface{}]interface{}
	SliceOfInts      []int
	SliceOfStructs   []testNestedLevel2_2
	SliceOfMapInfInf []map[interface{}]interface{}
}

type testNestedTypes struct {
	Nested testNestedLevel1
	TopInt int
}

type testStructInterfaceImpl struct {
	storage map[string]*testFlatStruct
}

func (t *testStructInterfaceImpl) FieldConfigurations() FieldConfigurations {
	// only going to allow foo, bar and baz for selectors

	subfields, _ := GenerateFieldConfigurations((*testFlatStruct)(nil))

	fields := make(FieldConfigurations)

	subfield := &FieldConfiguration{
		SubFields: subfields,
	}
	fields[FieldName("foo")] = subfield
	fields[FieldName("bar")] = subfield
	fields[FieldName("baz")] = subfield

	return fields
}

func (t *testStructInterfaceImpl) EvaluateMatch(sel Selector, op MatchOperator, value interface{}) (bool, error) {
	switch sel[0] {
	case "foo", "bar", "baz":
		storageVal, ok := t.storage[sel[0]]
		if !ok {
			// default to no match if this struct isn't stored
			return false, nil
		}

		if len(sel) < 2 {
			return false, fmt.Errorf("Need more selector")
		}

		dataType, ok := testFlatStructKindMap[sel[1]]
		if !ok {
			return false, fmt.Errorf("Invalid selector")
		}

		eqFn, ok := primitiveEqualityFns[dataType]
		if !ok {
			return false, fmt.Errorf("Invalid data type")
		}

		result := false
		switch sel[1] {
		case "Int":
			result = eqFn(storageVal.Int, value)
		case "Int8":
			result = eqFn(storageVal.Int8, value)
		case "Int16":
			result = eqFn(storageVal.Int16, value)
		case "Int32":
			result = eqFn(storageVal.Int32, value)
		case "Int64":
			result = eqFn(storageVal.Int64, value)
		case "Uint":
			result = eqFn(storageVal.Uint, value)
		case "Uint8":
			result = eqFn(storageVal.Uint8, value)
		case "Uint16":
			result = eqFn(storageVal.Uint16, value)
		case "Uint32":
			result = eqFn(storageVal.Uint32, value)
		case "Uint64":
			result = eqFn(storageVal.Uint64, value)
		case "Float32":
			result = eqFn(storageVal.Float32, value)
		case "Float64":
			result = eqFn(storageVal.Float64, value)
		case "Bool":
			result = eqFn(storageVal.Bool, value)
		case "String":
			result = eqFn(storageVal.String, value)
		default:
			return false, fmt.Errorf("Invalid data type")
		}

		if op == MatchNotEqual {
			return !result, nil
		}
		return result, nil
	default:
		return false, fmt.Errorf("Invalid selector")
	}
}

func validateFieldConfigurationsRecurse(t *testing.T, expected, actual FieldConfigurations, path string) bool {
	t.Helper()

	ok := assert.Len(t, actual, len(expected), "Actual FieldConfigurations length of %d != expected length of %d for path %q", len(actual), len(expected), path)

	for fieldName, expectedConfig := range expected {
		actualConfig, ok := actual[fieldName]
		ok = ok && assert.True(t, ok, "Actual configuration is missing field %q", fieldName)
		ok = ok && assert.Equal(t, expectedConfig.StructFieldName, actualConfig.StructFieldName, "Field %q on path %q have different StructFieldNames - Expected: %q, Actual: %q", fieldName, path, expectedConfig.StructFieldName, actualConfig.StructFieldName)
		ok = ok && assert.ElementsMatch(t, expectedConfig.SupportedOperations, actualConfig.SupportedOperations, "Fields %q on path %q have different SupportedOperations - Expected: %v, Actual: %v", fieldName, path, expectedConfig.SupportedOperations, actualConfig.SupportedOperations)

		newPath := string(fieldName)
		if newPath == "" {
			newPath = "*"
		}
		if path != "" {
			newPath = fmt.Sprintf("%s.%s", path, newPath)
		}
		ok = ok && validateFieldConfigurationsRecurse(t, expectedConfig.SubFields, actualConfig.SubFields, newPath)

		if !ok {
			break
		}
	}

	return ok
}

func validateFieldConfigurations(t *testing.T, expected, actual FieldConfigurations) {
	t.Helper()
	require.True(t, validateFieldConfigurationsRecurse(t, expected, actual, ""))
}

func dumpFieldConfigurationsRecurse(fields FieldConfigurations, level int, path string) {
	for fieldName, cfg := range fields {
		fmt.Printf("%sPath: %s Field: %s, StructFieldName: %s, CoerceFn: %p, SupportedOperations: %v\n", strings.Repeat("   ", level), path, fieldName, cfg.StructFieldName, cfg.CoerceFn, cfg.SupportedOperations)
		newPath := string(fieldName)
		if path != "" {
			newPath = fmt.Sprintf("%s.%s", path, fieldName)
		}
		dumpFieldConfigurationsRecurse(cfg.SubFields, level+1, newPath)
	}
}

func dumpFieldConfigurations(name string, fields FieldConfigurations) {
	fmt.Printf("===== %s =====\n", name)
	dumpFieldConfigurationsRecurse(fields, 1, "")
}
