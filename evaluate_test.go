// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package bexpr

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/mitchellh/pointerstructure"
	"github.com/stretchr/testify/require"
)

type expressionCheck struct {
	expression string
	result     bool
	err        string
	benchQuick bool
	hook       ValueTransformationHookFn
}

type expressionTest struct {
	value       interface{}
	expressions []expressionCheck
}

var evaluateTests map[string]expressionTest = map[string]expressionTest{
	"Flat Struct": {
		testFlatStruct{
			Int:         -1,
			Int8:        -2,
			Int16:       -3,
			Int32:       -4,
			Int64:       -5,
			Uint:        6,
			Uint8:       7,
			Uint16:      8,
			Uint32:      9,
			Uint64:      10,
			Float32:     1.1,
			Float64:     1.2,
			Bool:        true,
			String:      "exported",
			ColonString: "expo:rted",
			Slash:       "hello",
			unexported:  "unexported",
			Hidden:      true,
		},
		[]expressionCheck{
			{expression: "Int == -1", result: true, benchQuick: true},
			{expression: "Int == `foo`", result: true, hook: func(reflect.Value) reflect.Value { return reflect.ValueOf("foo") }},
			{expression: "Int == -99", result: false, benchQuick: true},
			{expression: "Int != -1", result: false},
			{expression: "Int != -99", result: true},
			{expression: "Int8 == -2", result: true},
			{expression: "Int8 == -99", result: false},
			{expression: "Int8 != -2", result: false},
			{expression: "Int8 != -99", result: true},
			{expression: "Int16 == -3", result: true},
			{expression: "Int16 == -99", result: false},
			{expression: "Int16 != -3", result: false},
			{expression: "Int16 != -99", result: true},
			{expression: "Int32 == -4", result: true},
			{expression: "Int32 == -99", result: false},
			{expression: "Int32 != -4", result: false},
			{expression: "Int32 != -99", result: true},
			{expression: "Int64 == -5", result: true},
			{expression: "Int64 == -99", result: false},
			{expression: "Int64 != -5", result: false},
			{expression: "Int64 != -99", result: true},
			{expression: "Uint == 6", result: true},
			{expression: "Uint == 99", result: false},
			{expression: "Uint != 6", result: false},
			{expression: "Uint != 99", result: true},
			{expression: "Uint8 == 7", result: true},
			{expression: "Uint8 == 99", result: false},
			{expression: "Uint8 != 7", result: false},
			{expression: "Uint8 != 99", result: true},
			{expression: "Uint16 == 8", result: true},
			{expression: "Uint16 == 99", result: false},
			{expression: "Uint16 != 8", result: false},
			{expression: "Uint16 != 99", result: true},
			{expression: "Uint32 == 9", result: true},
			{expression: "Uint32 == 99", result: false},
			{expression: "Uint32 != 9", result: false},
			{expression: "Uint32 != 99", result: true},
			{expression: "Uint64 == 10", result: true},
			{expression: "Uint64 == 99", result: false},
			{expression: "Uint64 != 10", result: false, benchQuick: true},
			{expression: "Uint64 != 99", result: true, benchQuick: true},
			{expression: "Float32 == 1.1", result: true},
			{expression: "Float32 == 9.9", result: false},
			{expression: "Float32 != 1.1", result: false},
			{expression: "Float32 != 9.9", result: true},
			{expression: "Float64 == 1.2", result: true},
			{expression: "Float64 == 9.9", result: false},
			{expression: "Float64 != 1.2", result: false, benchQuick: true},
			{expression: "Float64 != 9.9", result: true, benchQuick: true},
			{expression: "Bool == true", result: true},
			{expression: "Bool == false", result: false},
			{expression: "Bool != true", result: false},
			{expression: "Bool != false", result: true},
			{expression: "String == `exported`", result: true, benchQuick: true},
			{expression: "String == `not-it`", result: false, benchQuick: true},
			{expression: "String != `exported`", result: false},
			{expression: "String != `not-it`", result: true},
			{expression: "String != `not-it`", result: false, hook: func(value reflect.Value) reflect.Value { return reflect.ValueOf("not-it") }},
			{expression: "port in String", result: true, benchQuick: true},
			{expression: "part in String", result: false},
			{expression: "port not in String", result: false},
			{expression: "part not in String", result: true},
			{expression: "ColonString == `expo:rted`", result: true},
			{expression: "ColonString != `expor:ted`", result: true},
			{expression: "slash/value == `hello`", result: true},
			{expression: "unexported == `unexported`", result: false, err: `error finding value in datum: /unexported at part 0: couldn't find key: struct field with name "unexported"`},
			{expression: "Hidden == false", result: false, err: "error finding value in datum: /Hidden at part 0: struct field \"Hidden\" is ignored and cannot be used"},
			{expression: "String matches 	`^ex.*`", result: true, benchQuick: true},
			{expression: "String not matches `^anchored.*`", result: true, benchQuick: true},
			{expression: "String matches 	`^anchored.*`", result: false},
			{expression: "String not matches `^ex.*`", result: false},
		},
	},
	"Flat Struct Alt Types": {
		testFlatStructAlt{
			Int:        -1,
			Int8:       -2,
			Int16:      -3,
			Int32:      -4,
			Int64:      -5,
			Uint:       6,
			Uint8:      7,
			Uint16:     8,
			Uint32:     9,
			Uint64:     10,
			Float32:    1.1,
			Float64:    1.2,
			Bool:       true,
			String:     "exported",
			unexported: "unexported",
			Hidden:     true,
		},
		[]expressionCheck{
			{expression: "Int == -1", result: true, benchQuick: true},
			{expression: "Int == -99", result: false, benchQuick: true},
			{expression: "Int != -1", result: false},
			{expression: "Int != -99", result: true},
			{expression: "Int8 == -2", result: true},
			{expression: "Int8 == -99", result: false},
			{expression: "Int8 != -2", result: false},
			{expression: "Int8 != -99", result: true},
			{expression: "Int16 == -3", result: true},
			{expression: "Int16 == -99", result: false},
			{expression: "Int16 != -3", result: false},
			{expression: "Int16 != -99", result: true},
			{expression: "Int32 == -4", result: true},
			{expression: "Int32 == -99", result: false},
			{expression: "Int32 != -4", result: false},
			{expression: "Int32 != -99", result: true},
			{expression: "Int64 == -5", result: true},
			{expression: "Int64 == -99", result: false},
			{expression: "Int64 != -5", result: false},
			{expression: "Int64 != -99", result: true},
			{expression: "Uint == 6", result: true},
			{expression: "Uint == 99", result: false},
			{expression: "Uint != 6", result: false},
			{expression: "Uint != 99", result: true},
			{expression: "Uint8 == 7", result: true},
			{expression: "Uint8 == 99", result: false},
			{expression: "Uint8 != 7", result: false},
			{expression: "Uint8 != 99", result: true},
			{expression: "Uint16 == 8", result: true},
			{expression: "Uint16 == 99", result: false},
			{expression: "Uint16 != 8", result: false},
			{expression: "Uint16 != 99", result: true},
			{expression: "Uint32 == 9", result: true},
			{expression: "Uint32 == 99", result: false},
			{expression: "Uint32 != 9", result: false},
			{expression: "Uint32 != 99", result: true},
			{expression: "Uint64 == 10", result: true},
			{expression: "Uint64 == 99", result: false},
			{expression: "Uint64 != 10", result: false, benchQuick: true},
			{expression: "Uint64 != 99", result: true, benchQuick: true},
			{expression: "Float32 == 1.1", result: true},
			{expression: "Float32 == 9.9", result: false},
			{expression: "Float32 != 1.1", result: false},
			{expression: "Float32 != 9.9", result: true},
			{expression: "Float64 == 1.2", result: true},
			{expression: "Float64 == 9.9", result: false},
			{expression: "Float64 != 1.2", result: false, benchQuick: true},
			{expression: "Float64 != 9.9", result: true, benchQuick: true},
			{expression: "Bool == true", result: true},
			{expression: "Bool == false", result: false},
			{expression: "Bool != true", result: false},
			{expression: "Bool != false", result: true},
			{expression: "String == `exported`", result: true, benchQuick: true},
			{expression: "String == `not-it`", result: false, benchQuick: true},
			{expression: "String != `exported`", result: false},
			{expression: "String != `not-it`", result: true},
			{expression: "unexported == `unexported`", result: false, err: `error finding value in datum: /unexported at part 0: couldn't find key: struct field with name "unexported"`},
			{expression: "Hidden == false", result: false, err: "error finding value in datum: /Hidden at part 0: struct field \"Hidden\" is ignored and cannot be used"},
		},
	},
	"map[string]map[string]bool": {
		map[string]map[string]bool{
			"foo": {
				"bar": true,
				"baz": false,
			},
			"abc": nil,
			"co:lon": {
				"bar": true,
			},
		},
		[]expressionCheck{
			{expression: "foo == true", result: true, hook: func(v reflect.Value) reflect.Value {
				if r := v.MapIndex(reflect.ValueOf("bar")); !r.IsZero() {
					return r
				}
				return v
			}},
			{expression: "bar in foo", result: true},
			{expression: `bar in "/co:lon"`, result: true},
			{expression: "arg in foo", result: false},
			{expression: "arg not in foo", result: true},
			{expression: "baz not in foo", result: false},
			{expression: "foo is empty", result: false},
			{expression: "foo is not empty", result: true},
			{expression: "abc is empty", result: true},
			{expression: "abc is not empty", result: false},
			{expression: "foo in abc", result: false, benchQuick: true},
			{expression: "foo not in abc", result: true},
			{expression: "foo.bar == true", result: true},
			{expression: "foo.bar == false", result: false},
			{expression: "foo.baz == false", result: true},
			{expression: "foo.baz == true", result: false, benchQuick: true},
			{expression: "foo.bar != true", result: false},
			{expression: "foo.bar != false", result: true},
			{expression: "foo.baz != false", result: false},
			{expression: "foo.baz != true", result: true},
			{expression: "foo.bar.baz == 3", result: false, err: `error finding value in datum: /foo/bar/baz: at part 2, invalid value kind: bool`},
		},
	},
	"Nested Structs and Maps": {
		testNestedTypes{
			Nested: testNestedLevel1{
				Map: map[string]string{
					"foo":    "bar",
					"bar":    "baz",
					"abc":    "123",
					"colon":  "co:lon",
					"co:lon": "co:lon",
					"email":  "foo@example.com",
				},
				MapOfStructs: map[string]testNestedLevel2_1{
					"one": {
						Foo: 42,
						bar: "unexported",
						Baz: "exported",
					},
					"two": {
						Foo: 77,
						bar: "unexported",
						Baz: "consul",
					},
				},
				SliceOfInts: []int{1, 3, 5, 7, 9},
				SliceOfStructs: []testNestedLevel2_2{
					{
						X: 1,
						Y: 2,
						z: 10,
					},
					{
						X: 3,
						Y: 5,
						z: 10,
					},
				},
				SliceOfInfs: []interface{}{"foobar", 1, true},
			},
			TopInt: 5,
		},
		[]expressionCheck{
			{expression: "Nested.Map == bar", result: true, benchQuick: true, hook: func(v reflect.Value) reflect.Value {
				if r, ok := v.Interface().(map[string]string); ok {
					return reflect.ValueOf(r["foo"])
				}
				return v
			}},
			{expression: "Nested.Map.foo == bar", result: true, benchQuick: true},
			{expression: "Nested.Map.foo contains ba", result: true, benchQuick: true},
			{expression: "Nested.Map.foo == baz", result: false},
			{expression: "Nested.Map is not empty", result: true},
			{expression: "Nested.Map is not empty", result: true},
			{expression: "Nested.Map contains foo and Nested.Map contains bar", result: true, benchQuick: true},
			{expression: `Nested.Map.colon == "co:lon"`, result: true},
			{expression: `"/Nested/Map/co:lon" == "co:lon"`, result: true},
			{expression: "Nested.Map contains nope", result: false},
			{expression: "Nested.Map contains bar", result: true},
			{expression: "Nested.Map.bar == `bazel`", result: false, benchQuick: true},
			{expression: "TopInt != 0", result: true},
			{expression: "Nested.Map contains nope or (Nested.Map contains bar and Nested.Map.bar == `bazel`) or TopInt != 0", result: true, benchQuick: true},
			{expression: "Nested.MapOfStructs.one.Foo == 42", result: true},
			{expression: "Nested.MapOfStructs.one.bar == `unexported`", result: false, err: `error finding value in datum: /Nested/MapOfStructs/one/bar at part 3: couldn't find key: struct field with name "bar"`},
			{expression: "Nested.MapOfStructs.one.Baz == `exported`", result: true},
			{expression: "Nested.MapOfStructs.two.Foo == 77", result: true},
			{expression: "Nested.MapOfStructs.two.bar == `unexported`", result: false, err: `error finding value in datum: /Nested/MapOfStructs/two/bar at part 3: couldn't find key: struct field with name "bar"`},
			{expression: "Nested.MapOfStructs.two.Baz == `consul`", result: true},
			{expression: "7 in Nested.SliceOfInts", result: true},
			{expression: `"/Nested/SliceOfInts" == "7"`, result: false, err: `unable to find suitable primitive comparison function for matching`},
			{expression: "Nested.MapOfStructs is empty or (Nested.SliceOfInts contains 7 and 9 in Nested.SliceOfInts)", result: true, benchQuick: true},
			{expression: "Nested.SliceOfStructs.0.X == 1", result: true},
			{expression: "Nested.SliceOfStructs.0.Y == 4", result: false},
			{expression: "Map in Nested", result: false, err: "Cannot perform in/contains operations on type struct for selector: \"Nested\""},
			{expression: `"foobar" in "/Nested/SliceOfInfs"`, result: true},
			{expression: `"1" in "/Nested/SliceOfInfs"`, result: true},
			{expression: `"2" in "/Nested/SliceOfInfs"`, result: false},
			{expression: `"true" in "/Nested/SliceOfInfs"`, result: true},
			{expression: `"/Nested/Map/email" matches "(foz|foo)@example.com"`, result: true},
			// Missing key in map tests
			{expression: "Nested.Map.notfound == 4", result: false},
			{expression: "Nested.Map.notfound != 4", result: true},
			{expression: "4 in Nested.Map.notfound", result: false},
			{expression: "4 not in Nested.Map.notfound", result: true},
			{expression: "Nested.Map.notfound is empty", result: true},
			{expression: "Nested.Map.notfound is not empty", result: false},
			{expression: `Nested.Map.notfound matches ".*"`, result: false},
			{expression: `Nested.Map.notfound not matches ".*"`, result: true},
			// Missing field in struct tests
			{expression: "Nested.Notfound == 4", result: false, err: `error finding value in datum: /Nested/Notfound at part 1: couldn't find key: struct field with name "Notfound"`},
			{expression: "Nested.Notfound != 4", result: false, err: `error finding value in datum: /Nested/Notfound at part 1: couldn't find key: struct field with name "Notfound"`},
			{expression: "4 in Nested.Notfound", result: false, err: `error finding value in datum: /Nested/Notfound at part 1: couldn't find key: struct field with name "Notfound"`},
			{expression: "4 not in Nested.Notfound", result: false, err: `error finding value in datum: /Nested/Notfound at part 1: couldn't find key: struct field with name "Notfound"`},
			{expression: "Nested.Notfound is empty", result: false, err: `error finding value in datum: /Nested/Notfound at part 1: couldn't find key: struct field with name "Notfound"`},
			{expression: "Nested.Notfound is not empty", result: false, err: `error finding value in datum: /Nested/Notfound at part 1: couldn't find key: struct field with name "Notfound"`},
			{expression: `Nested.Notfound matches ".*"`, result: false, err: `error finding value in datum: /Nested/Notfound at part 1: couldn't find key: struct field with name "Notfound"`},
			{expression: `Nested.Notfound not matches ".*"`, result: false, err: `error finding value in datum: /Nested/Notfound at part 1: couldn't find key: struct field with name "Notfound"`},
			// all
			{expression: `all Nested.SliceOfInts as i { i != 42 }`, result: true},
			{expression: `all Nested.SliceOfInts as i { i == 1 }`, result: false},
			{expression: `all Nested.Map as v { v == "bar" }`, result: false},
			{expression: `all Nested.Map as v { v != "hello" }`, result: true},
			{expression: `all Nested.Map as k, k { TopInt == 5 }`, err: `"k" cannot be used as a placeholder for both the index and the value`},
			{expression: `all Nested.Map as k, _ { k != "foo" }`, result: false},
			{expression: `all Nested.Map as k, _ { k != "hello" }`, result: true},
			{expression: `all Nested.Map as k, v { k != "foo" or v != "baz" }`, result: true},
			{expression: `all TopInt as k, v { k != "foo" or v != "baz" }`, err: "TopInt is not a list or a map"},
			// any
			{expression: `any Nested.SliceOfInts as i { i == 1 }`, result: true},
			{expression: `any Nested.SliceOfInts as i { i == 42 }`, result: false},
			{expression: `any Nested.SliceOfStructs as i { "/i/X" == 1 }`, result: true},
			{expression: `any Nested.Map as k { k != "bar" }`, result: true},
			{expression: `any Nested.Map as k { k == "bar" }`, result: true},
			{expression: `any Nested.Map as k { k == "hello" }`, result: false},
			{expression: `any Nested.Map as k, v { k == "foo" and v == "bar" }`, result: true},
			{expression: `any Nested.Map as k { k.Color == "red" }`, err: "/k references a string so /k/Color is invalid"},
			{expression: `any Nested.SliceOfInts as i, _ { i.Color == "red" }`, err: "/i references a int so /i/Color is invalid"},
		},
	},
}

func TestEvaluate(t *testing.T) {
	t.Parallel()
	for name, tcase := range evaluateTests {
		// capture these values in the closure
		name := name
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for i, expTest := range tcase.expressions {
				// capture these values in the closure
				expTest := expTest
				t.Run(fmt.Sprintf("#%d - %s", i, expTest.expression), func(t *testing.T) {
					t.Parallel()

					expr, err := CreateEvaluator(expTest.expression, WithHookFn(expTest.hook))
					require.NoError(t, err)

					match, err := expr.Evaluate(tcase.value)
					if expTest.err != "" {
						require.Error(t, err)
						require.EqualError(t, err, expTest.err)
					} else {
						require.NoError(t, err)
					}
					require.Equal(t, expTest.result, match)
				})
			}
		})
	}
}

func TestWithHookFn(t *testing.T) {
	t.Parallel()
	type testStruct struct {
		I interface{}
		S *testStruct
	}
	cases := []struct {
		name string
		hook ValueTransformationHookFn
		in   *testStruct
		eval []expressionCheck
	}{
		{
			name: "simple",
			hook: func(v reflect.Value) reflect.Value { return v },
			in:   &testStruct{I: "foo"},
			eval: []expressionCheck{
				{expression: `"/I"=="foo"`, result: true},
			},
		},
		{
			name: "dive to pointer",
			hook: func(v reflect.Value) reflect.Value {
				if r, ok := v.Interface().(*testStruct); ok {
					return reflect.ValueOf(r.I)
				}
				return v
			},
			in: &testStruct{S: &testStruct{I: "foo"}, I: &testStruct{I: &testStruct{I: "bar"}}},
			eval: []expressionCheck{
				{expression: `"/S"=="foo"`, result: true},
				{expression: `"/I/I"=="bar"`, result: true},
				{
					expression: `"/S/I"=="foo"`, result: false,
					err: "error finding value in datum: /S/I: at part 1, invalid value kind: string",
				},
			},
		},
		{
			name: "valueTransformationHook returns nil interface{}",
			hook: func(v reflect.Value) reflect.Value { return reflect.ValueOf(nil) },
			in:   &testStruct{I: "foo"},
			eval: []expressionCheck{
				{
					expression: `"/I"=="foo"`, result: false,
					err: "error finding value in datum: /I at part 0: ValueTransformationHook returned the value of a nil interface",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, eval := range tc.eval {
				expr, err := CreateEvaluator(eval.expression, WithHookFn(tc.hook))
				require.NoError(t, err)

				match, err := expr.Evaluate(tc.in)
				if eval.err != "" {
					require.Error(t, err)
					require.Equal(t, eval.err, err.Error())
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, eval.result, match)
			}
		})
	}
}

func TestUnknownVal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		expression string
		unknownVal interface{}
		result     bool
		err        string
	}{
		{
			name:       "key exists",
			expression: `key == "foo"`,
			unknownVal: "bar",
			result:     true,
		},
		{
			name:       "key does not exist",
			expression: `unknown == "bar"`,
			unknownVal: "bar",
			result:     true,
		},
		{
			name:       "key does not exist (fail)",
			expression: `unknown == "qux"`,
			unknownVal: "bar",
			result:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := CreateEvaluator(tc.expression, WithUnknownValue(tc.unknownVal))
			require.NoError(t, err)

			match, err := expr.Evaluate(map[string]string{
				"key": "foo",
			})
			if tc.err != "" {
				require.Error(t, err)
				require.EqualError(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.result, match)
		})
	}
}

func TestUnknownVal_struct(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		expression string
		unknownVal interface{}
		result     bool
		err        string
	}{
		{
			name:       "key exists",
			expression: `key == "foo"`,
			unknownVal: "bar",
			result:     true,
		},
		{
			name:       "key does not exist",
			expression: `unknown == "bar"`,
			unknownVal: "bar",
			result:     true,
		},
		{
			name:       "key does not exist (fail)",
			expression: `unknown == "qux"`,
			unknownVal: "bar",
			result:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := CreateEvaluator(tc.expression, WithUnknownValue(tc.unknownVal))
			require.NoError(t, err)

			match, err := expr.Evaluate(struct {
				Key string `bexpr:"key"`
			}{
				Key: "foo",
			})
			if tc.err != "" {
				require.Error(t, err)
				require.EqualError(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.result, match)
		})
	}
}

func TestCustomTag(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		BexprName string `bexpr:"bname"`
		JsonName  string `json:"jname"`
	}
	ts := testStruct{BexprName: "foo", JsonName: "bar"}

	cases := []struct {
		name       string
		expression string
		jsonTag    bool
		bnameFound bool
		jnameFound bool
	}{
		{
			name:       "bexpr tag, bname",
			expression: `"/bname" == "foo"`,
			bnameFound: true,
		},
		{
			name:       "bexpr tag, jname",
			expression: `"/jname" == "bar"`,
		},
		{
			name:       "json tag, bname",
			expression: `"/bname" == "foo"`,
			jsonTag:    true,
		},
		{
			name:       "json tag, jname",
			expression: `"/jname" == "bar"`,
			jsonTag:    true,
			jnameFound: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var opts []Option
			if tc.jsonTag {
				opts = append(opts, WithTagName("json"))
			}
			expr, err := CreateEvaluator(tc.expression, opts...)
			require.NoError(t, err)

			match, err := expr.Evaluate(ts)
			if tc.jsonTag {
				if tc.jnameFound {
					require.NoError(t, err)
					require.True(t, match)
				} else {
					require.True(t, errors.Is(err, pointerstructure.ErrNotFound))
				}
			} else {
				if tc.bnameFound {
					require.NoError(t, err)
					require.True(t, match)
				} else {
					require.True(t, errors.Is(err, pointerstructure.ErrNotFound))
				}
			}
		})
	}
}

func TestInOnOperator(t *testing.T) {
	type testStruct struct {
		Role  any
		Match bool
	}

	exprs := []string{
		"Role contains admin",
		"admin in Role",
	}
	cases := []testStruct{
		{
			Role:  []string{"admin", "foo"},
			Match: true,
		},
		{
			Role:  "admin",
			Match: true,
		},
		{
			Role:  "dbadmin",
			Match: true, // dangerous!
		},
		{
			Role:  []string{"dbadmin", "foo"},
			Match: false,
		},
	}

	for _, q := range exprs {
		expr, err := CreateEvaluator(q)
		require.NoError(t, err)

		for _, tc := range cases {
			match, err := expr.Evaluate(tc)
			require.NoError(t, err)
			require.Equal(t, tc.Match, match)
		}
	}
}

func BenchmarkEvaluate(b *testing.B) {
	for name, tcase := range evaluateTests {
		// capture these values in the closure
		name := name
		tcase := tcase
		b.Run(name, func(b *testing.B) {
			for i, expTest := range tcase.expressions {
				// capture these values in the closure
				expTest := expTest
				b.Run(fmt.Sprintf("#%d", i), func(b *testing.B) {
					if !expTest.benchQuick && !FullBenchmarks() {
						b.Skip("Skipping benchmark - rerun with -bench-full to enable")
					}

					expr, err := CreateEvaluator(expTest.expression, WithHookFn(expTest.hook))
					require.NoError(b, err)

					b.ResetTimer()
					for n := 0; n < b.N; n++ {
						_, err = expr.Evaluate(tcase.value)
						if expTest.err != "" {
							require.Error(b, err)
						} else {
							require.NoError(b, err)
						}
					}
				})
			}
		})
	}
}
