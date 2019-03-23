package bexpr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
