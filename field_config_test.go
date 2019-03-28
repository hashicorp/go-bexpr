package bexpr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type fieldConfigTest struct {
	dataType   interface{}
	expected   FieldConfigurations
	err        string
	benchQuick bool
}

var fieldConfigTests map[string]fieldConfigTest = map[string]fieldConfigTest{
	"Flat Struct": {
		dataType: (*testFlatStruct)(nil),
		expected: FieldConfigurations{
			"Int":     &FieldConfiguration{StructFieldName: "Int", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Int8":    &FieldConfiguration{StructFieldName: "Int8", CoerceFn: CoerceInt8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Int16":   &FieldConfiguration{StructFieldName: "Int16", CoerceFn: CoerceInt16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Int32":   &FieldConfiguration{StructFieldName: "Int32", CoerceFn: CoerceInt32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Int64":   &FieldConfiguration{StructFieldName: "Int64", CoerceFn: CoerceInt64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Uint":    &FieldConfiguration{StructFieldName: "Uint", CoerceFn: CoerceUint, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Uint8":   &FieldConfiguration{StructFieldName: "Uint8", CoerceFn: CoerceUint8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Uint16":  &FieldConfiguration{StructFieldName: "Uint16", CoerceFn: CoerceUint16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Uint32":  &FieldConfiguration{StructFieldName: "Uint32", CoerceFn: CoerceUint32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Uint64":  &FieldConfiguration{StructFieldName: "Uint64", CoerceFn: CoerceUint64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Float32": &FieldConfiguration{StructFieldName: "Float32", CoerceFn: CoerceFloat32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Float64": &FieldConfiguration{StructFieldName: "Float64", CoerceFn: CoerceFloat64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"Bool":    &FieldConfiguration{StructFieldName: "Bool", CoerceFn: CoerceBool, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			"String":  &FieldConfiguration{StructFieldName: "String", CoerceFn: CoerceString, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
		},
		benchQuick: true,
	},
	"map[string]bool": {
		dataType: (*map[string]bool)(nil),
		expected: FieldConfigurations{
			FieldNameAny: &FieldConfiguration{CoerceFn: CoerceBool, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
		},
		benchQuick: true,
	},
	"map[string]interface{}": {
		dataType: (*map[string]interface{})(nil),
	},
	"map[interface{}]interface{}": {
		dataType: (*map[interface{}]interface{})(nil),
		err:      "Cannot generate FieldConfigurations for maps with keys that are not strings",
	},
	"[]map[string]string": {
		dataType: (*[]map[string]string)(nil),
		err:      "Invalid top level type - can only use structs, map[string]* or an MatchExpressionEvaluator",
	},
	"Nested Structs and Maps": {
		dataType: (*testNestedTypes)(nil),
		expected: FieldConfigurations{
			"Nested": &FieldConfiguration{StructFieldName: "Nested", SubFields: FieldConfigurations{
				"Map": &FieldConfiguration{StructFieldName: "Map", SupportedOperations: []MatchOperator{MatchIn, MatchNotIn, MatchIsEmpty, MatchIsNotEmpty}, SubFields: FieldConfigurations{
					FieldNameAny: &FieldConfiguration{StructFieldName: "", SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				}},
				"MapOfStructs": &FieldConfiguration{StructFieldName: "MapOfStructs", SupportedOperations: []MatchOperator{MatchIsEmpty, MatchIsNotEmpty, MatchIn, MatchNotIn}, SubFields: FieldConfigurations{
					FieldNameAny: &FieldConfiguration{StructFieldName: "", SubFields: FieldConfigurations{
						"Foo": &FieldConfiguration{StructFieldName: "Foo", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
						"Baz": &FieldConfiguration{StructFieldName: "Baz", SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					}},
				}},
				"MapInfInf":   &FieldConfiguration{StructFieldName: "MapInfInf", SupportedOperations: []MatchOperator{MatchIsEmpty, MatchIsNotEmpty}},
				"SliceOfInts": &FieldConfiguration{StructFieldName: "SliceOfInts", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchIn, MatchNotIn, MatchIsEmpty, MatchIsNotEmpty}},
				"SliceOfStructs": &FieldConfiguration{StructFieldName: "SliceOfStructs", SupportedOperations: []MatchOperator{MatchIsEmpty, MatchIsNotEmpty}, SubFields: FieldConfigurations{
					"X": &FieldConfiguration{StructFieldName: "X", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Y": &FieldConfiguration{StructFieldName: "Y", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				}},
				"SliceOfMapInfInf": &FieldConfiguration{StructFieldName: "SliceOfMapInfInf", SupportedOperations: []MatchOperator{MatchIsEmpty, MatchIsNotEmpty}},
			}},
			"TopInt": &FieldConfiguration{StructFieldName: "TopInt", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
		},
		benchQuick: true,
	},
	"Interface Implementor": {
		dataType: (*testStructInterfaceImpl)(nil),
		expected: FieldConfigurations{
			"foo": &FieldConfiguration{SubFields: FieldConfigurations{
				"Int":     &FieldConfiguration{StructFieldName: "Int", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int8":    &FieldConfiguration{StructFieldName: "Int8", CoerceFn: CoerceInt8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int16":   &FieldConfiguration{StructFieldName: "Int16", CoerceFn: CoerceInt16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int32":   &FieldConfiguration{StructFieldName: "Int32", CoerceFn: CoerceInt32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int64":   &FieldConfiguration{StructFieldName: "Int64", CoerceFn: CoerceInt64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint":    &FieldConfiguration{StructFieldName: "Uint", CoerceFn: CoerceUint, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint8":   &FieldConfiguration{StructFieldName: "Uint8", CoerceFn: CoerceUint8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint16":  &FieldConfiguration{StructFieldName: "Uint16", CoerceFn: CoerceUint16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint32":  &FieldConfiguration{StructFieldName: "Uint32", CoerceFn: CoerceUint32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint64":  &FieldConfiguration{StructFieldName: "Uint64", CoerceFn: CoerceUint64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Float32": &FieldConfiguration{StructFieldName: "Float32", CoerceFn: CoerceFloat32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Float64": &FieldConfiguration{StructFieldName: "Float64", CoerceFn: CoerceFloat64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Bool":    &FieldConfiguration{StructFieldName: "Bool", CoerceFn: CoerceBool, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"String":  &FieldConfiguration{StructFieldName: "String", CoerceFn: CoerceString, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			}},
			"bar": &FieldConfiguration{SubFields: FieldConfigurations{
				"Int":     &FieldConfiguration{StructFieldName: "Int", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int8":    &FieldConfiguration{StructFieldName: "Int8", CoerceFn: CoerceInt8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int16":   &FieldConfiguration{StructFieldName: "Int16", CoerceFn: CoerceInt16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int32":   &FieldConfiguration{StructFieldName: "Int32", CoerceFn: CoerceInt32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int64":   &FieldConfiguration{StructFieldName: "Int64", CoerceFn: CoerceInt64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint":    &FieldConfiguration{StructFieldName: "Uint", CoerceFn: CoerceUint, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint8":   &FieldConfiguration{StructFieldName: "Uint8", CoerceFn: CoerceUint8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint16":  &FieldConfiguration{StructFieldName: "Uint16", CoerceFn: CoerceUint16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint32":  &FieldConfiguration{StructFieldName: "Uint32", CoerceFn: CoerceUint32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint64":  &FieldConfiguration{StructFieldName: "Uint64", CoerceFn: CoerceUint64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Float32": &FieldConfiguration{StructFieldName: "Float32", CoerceFn: CoerceFloat32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Float64": &FieldConfiguration{StructFieldName: "Float64", CoerceFn: CoerceFloat64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Bool":    &FieldConfiguration{StructFieldName: "Bool", CoerceFn: CoerceBool, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"String":  &FieldConfiguration{StructFieldName: "String", CoerceFn: CoerceString, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			}},
			"baz": &FieldConfiguration{SubFields: FieldConfigurations{
				"Int":     &FieldConfiguration{StructFieldName: "Int", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int8":    &FieldConfiguration{StructFieldName: "Int8", CoerceFn: CoerceInt8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int16":   &FieldConfiguration{StructFieldName: "Int16", CoerceFn: CoerceInt16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int32":   &FieldConfiguration{StructFieldName: "Int32", CoerceFn: CoerceInt32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Int64":   &FieldConfiguration{StructFieldName: "Int64", CoerceFn: CoerceInt64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint":    &FieldConfiguration{StructFieldName: "Uint", CoerceFn: CoerceUint, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint8":   &FieldConfiguration{StructFieldName: "Uint8", CoerceFn: CoerceUint8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint16":  &FieldConfiguration{StructFieldName: "Uint16", CoerceFn: CoerceUint16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint32":  &FieldConfiguration{StructFieldName: "Uint32", CoerceFn: CoerceUint32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Uint64":  &FieldConfiguration{StructFieldName: "Uint64", CoerceFn: CoerceUint64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Float32": &FieldConfiguration{StructFieldName: "Float32", CoerceFn: CoerceFloat32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Float64": &FieldConfiguration{StructFieldName: "Float64", CoerceFn: CoerceFloat64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"Bool":    &FieldConfiguration{StructFieldName: "Bool", CoerceFn: CoerceBool, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				"String":  &FieldConfiguration{StructFieldName: "String", CoerceFn: CoerceString, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
			}},
		},
		benchQuick: true,
	},
	"Map of Interface Implementor": {
		dataType: (*map[string]*testStructInterfaceImpl)(nil),
		expected: FieldConfigurations{
			FieldNameAny: &FieldConfiguration{SubFields: FieldConfigurations{
				"foo": &FieldConfiguration{SubFields: FieldConfigurations{
					"Int":     &FieldConfiguration{StructFieldName: "Int", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int8":    &FieldConfiguration{StructFieldName: "Int8", CoerceFn: CoerceInt8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int16":   &FieldConfiguration{StructFieldName: "Int16", CoerceFn: CoerceInt16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int32":   &FieldConfiguration{StructFieldName: "Int32", CoerceFn: CoerceInt32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int64":   &FieldConfiguration{StructFieldName: "Int64", CoerceFn: CoerceInt64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint":    &FieldConfiguration{StructFieldName: "Uint", CoerceFn: CoerceUint, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint8":   &FieldConfiguration{StructFieldName: "Uint8", CoerceFn: CoerceUint8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint16":  &FieldConfiguration{StructFieldName: "Uint16", CoerceFn: CoerceUint16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint32":  &FieldConfiguration{StructFieldName: "Uint32", CoerceFn: CoerceUint32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint64":  &FieldConfiguration{StructFieldName: "Uint64", CoerceFn: CoerceUint64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Float32": &FieldConfiguration{StructFieldName: "Float32", CoerceFn: CoerceFloat32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Float64": &FieldConfiguration{StructFieldName: "Float64", CoerceFn: CoerceFloat64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Bool":    &FieldConfiguration{StructFieldName: "Bool", CoerceFn: CoerceBool, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"String":  &FieldConfiguration{StructFieldName: "String", CoerceFn: CoerceString, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				}},
				"bar": &FieldConfiguration{SubFields: FieldConfigurations{
					"Int":     &FieldConfiguration{StructFieldName: "Int", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int8":    &FieldConfiguration{StructFieldName: "Int8", CoerceFn: CoerceInt8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int16":   &FieldConfiguration{StructFieldName: "Int16", CoerceFn: CoerceInt16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int32":   &FieldConfiguration{StructFieldName: "Int32", CoerceFn: CoerceInt32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int64":   &FieldConfiguration{StructFieldName: "Int64", CoerceFn: CoerceInt64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint":    &FieldConfiguration{StructFieldName: "Uint", CoerceFn: CoerceUint, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint8":   &FieldConfiguration{StructFieldName: "Uint8", CoerceFn: CoerceUint8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint16":  &FieldConfiguration{StructFieldName: "Uint16", CoerceFn: CoerceUint16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint32":  &FieldConfiguration{StructFieldName: "Uint32", CoerceFn: CoerceUint32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint64":  &FieldConfiguration{StructFieldName: "Uint64", CoerceFn: CoerceUint64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Float32": &FieldConfiguration{StructFieldName: "Float32", CoerceFn: CoerceFloat32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Float64": &FieldConfiguration{StructFieldName: "Float64", CoerceFn: CoerceFloat64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Bool":    &FieldConfiguration{StructFieldName: "Bool", CoerceFn: CoerceBool, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"String":  &FieldConfiguration{StructFieldName: "String", CoerceFn: CoerceString, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				}},
				"baz": &FieldConfiguration{SubFields: FieldConfigurations{
					"Int":     &FieldConfiguration{StructFieldName: "Int", CoerceFn: CoerceInt, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int8":    &FieldConfiguration{StructFieldName: "Int8", CoerceFn: CoerceInt8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int16":   &FieldConfiguration{StructFieldName: "Int16", CoerceFn: CoerceInt16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int32":   &FieldConfiguration{StructFieldName: "Int32", CoerceFn: CoerceInt32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Int64":   &FieldConfiguration{StructFieldName: "Int64", CoerceFn: CoerceInt64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint":    &FieldConfiguration{StructFieldName: "Uint", CoerceFn: CoerceUint, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint8":   &FieldConfiguration{StructFieldName: "Uint8", CoerceFn: CoerceUint8, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint16":  &FieldConfiguration{StructFieldName: "Uint16", CoerceFn: CoerceUint16, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint32":  &FieldConfiguration{StructFieldName: "Uint32", CoerceFn: CoerceUint32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Uint64":  &FieldConfiguration{StructFieldName: "Uint64", CoerceFn: CoerceUint64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Float32": &FieldConfiguration{StructFieldName: "Float32", CoerceFn: CoerceFloat32, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Float64": &FieldConfiguration{StructFieldName: "Float64", CoerceFn: CoerceFloat64, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"Bool":    &FieldConfiguration{StructFieldName: "Bool", CoerceFn: CoerceBool, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
					"String":  &FieldConfiguration{StructFieldName: "String", CoerceFn: CoerceString, SupportedOperations: []MatchOperator{MatchEqual, MatchNotEqual}},
				}},
			}},
		},
		benchQuick: true,
	},
}

func TestGenerateFieldConfigurations(t *testing.T) {
	t.Parallel()
	for name, tcase := range fieldConfigTests {
		// capture these values in the closure
		name := name
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fields, err := GenerateFieldConfigurations(tcase.dataType)
			if tcase.err == "" {
				require.NoError(t, err)
				validateFieldConfigurations(t, tcase.expected, fields)
			} else {
				require.Len(t, fields, 0)
				require.Error(t, err)
				require.EqualError(t, err, tcase.err)
			}
		})
	}
}

func BenchmarkGenerateFieldConfigurations(b *testing.B) {
	for name, tcase := range fieldConfigTests {
		if tcase.benchQuick || FullBenchmarks() {
			b.Run(name, func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					_, err := GenerateFieldConfigurations(tcase.dataType)
					if tcase.err == "" {
						require.NoError(b, err)
					} else {
						require.Error(b, err)
					}
				}
			})
		}
	}
}
