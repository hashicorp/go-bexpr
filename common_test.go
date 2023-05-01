// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bexpr

import (
	"flag"
	"reflect"
)

var benchFull *bool = flag.Bool("bench-full", false, "Run all benchmarks rather than a subset")

func FullBenchmarks() bool {
	return benchFull != nil && *benchFull
}

type testFlatStruct struct {
	Int         int
	Int8        int8
	Int16       int16
	Int32       int32
	Int64       int64
	Uint        uint
	Uint8       uint8
	Uint16      uint16
	Uint32      uint32
	Uint64      uint64
	Float32     float32
	Float64     float64
	Bool        bool
	String      string
	ColonString string
	Slash       string `bexpr:"slash/value"`
	unexported  string
	Hidden      bool `bexpr:"-"`
}

type (
	CustomInt     int
	CustomInt8    int8
	CustomInt16   int16
	CustomInt32   int32
	CustomInt64   int64
	CustomUint    uint
	CustomUint8   uint8
	CustomUint16  uint16
	CustomUint32  uint32
	CustomUint64  uint64
	CustomFloat32 float32
	CustomFloat64 float64
	CustomBool    bool
	CustomString  string
)

type testFlatStructAlt struct {
	Int        CustomInt
	Int8       CustomInt8
	Int16      CustomInt16
	Int32      CustomInt32
	Int64      CustomInt64
	Uint       CustomUint
	Uint8      CustomUint8
	Uint16     CustomUint16
	Uint32     CustomUint32
	Uint64     CustomUint64
	Float32    CustomFloat32
	Float64    CustomFloat64
	Bool       CustomBool
	String     CustomString
	unexported CustomString
	Hidden     CustomBool `bexpr:"-"`
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
	SliceOfInfs      []interface{}
}

type testNestedTypes struct {
	Nested testNestedLevel1
	TopInt int
}
