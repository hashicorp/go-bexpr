package bexpr

import (
	"flag"
	"fmt"
	"reflect"
)

var benchFull *bool = flag.Bool("bench-full", false, "Run all benchmarks rather than a subset")

func FullBenchmarks() bool {
	return benchFull != nil && *benchFull
}

type testFlatStruct struct {
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64
	Bool    bool
	String  string
	Hidden  bool `bexpr:"-"`
}

type CustomInt int
type CustomInt8 int8
type CustomInt16 int16
type CustomInt32 int32
type CustomInt64 int64
type CustomUint uint
type CustomUint8 uint8
type CustomUint16 uint16
type CustomUint32 uint32
type CustomUint64 uint64
type CustomFloat32 float32
type CustomFloat64 float64
type CustomBool bool
type CustomString string

type testFlatStructAlt struct {
	Int     CustomInt
	Int8    CustomInt8
	Int16   CustomInt16
	Int32   CustomInt32
	Int64   CustomInt64
	Uint    CustomUint
	Uint8   CustomUint8
	Uint16  CustomUint16
	Uint32  CustomUint32
	Uint64  CustomUint64
	Float32 CustomFloat32
	Float64 CustomFloat64
	Bool    CustomBool
	String  CustomString
	Hidden  CustomBool `bexpr:"-"`
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

func (t *testStructInterfaceImpl) EvaluateMatch(selector Selector, op MatchOperator, value interface{}) (bool, error) {
	sel := selector.Path
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
			result = eqFn(value, reflect.ValueOf(storageVal.Int))
		case "Int8":
			result = eqFn(value, reflect.ValueOf(storageVal.Int8))
		case "Int16":
			result = eqFn(value, reflect.ValueOf(storageVal.Int16))
		case "Int32":
			result = eqFn(value, reflect.ValueOf(storageVal.Int32))
		case "Int64":
			result = eqFn(value, reflect.ValueOf(storageVal.Int64))
		case "Uint":
			result = eqFn(value, reflect.ValueOf(storageVal.Uint))
		case "Uint8":
			result = eqFn(value, reflect.ValueOf(storageVal.Uint8))
		case "Uint16":
			result = eqFn(value, reflect.ValueOf(storageVal.Uint16))
		case "Uint32":
			result = eqFn(value, reflect.ValueOf(storageVal.Uint32))
		case "Uint64":
			result = eqFn(value, reflect.ValueOf(storageVal.Uint64))
		case "Float32":
			result = eqFn(value, reflect.ValueOf(storageVal.Float32))
		case "Float64":
			result = eqFn(value, reflect.ValueOf(storageVal.Float64))
		case "Bool":
			result = eqFn(value, reflect.ValueOf(storageVal.Bool))
		case "String":
			result = eqFn(value, reflect.ValueOf(storageVal.String))
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
