// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bexpr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	X int
	Y string
}

var testSlice []testStruct = []testStruct{
	{
		X: 1,
		Y: "a",
	},
	{
		X: 1,
		Y: "b",
	},
	{
		X: 2,
		Y: "a",
	},
	{
		X: 2,
		Y: "b",
	},
	{
		X: 3,
		Y: "c",
	},
}

var testArray [5]testStruct = [5]testStruct{
	{
		X: 1,
		Y: "a",
	},
	{
		X: 1,
		Y: "b",
	},
	{
		X: 2,
		Y: "a",
	},
	{
		X: 2,
		Y: "b",
	},
	{
		X: 3,
		Y: "c",
	},
}

var testMap map[string]testStruct = map[string]testStruct{
	"one": {
		X: 1,
		Y: "a",
	},
	"two": {
		X: 1,
		Y: "b",
	},
	"three": {
		X: 2,
		Y: "a",
	},
	"four": {
		X: 2,
		Y: "b",
	},
	"five": {
		X: 3,
		Y: "c",
	},
}

func TestFilter(t *testing.T) {
	type testCase struct {
		expression string
		input      interface{}
		expected   interface{}
	}

	cases := map[string]testCase{
		"Slice X==1": {
			"X==1",
			testSlice,
			[]testStruct{
				{X: 1, Y: "a"},
				{X: 1, Y: "b"},
			},
		},
		"Slice Y==`c`": {
			"Y==`c`",
			testSlice,
			[]testStruct{
				{X: 3, Y: "c"},
			},
		},
		"Array X==1": {
			"X==1",
			testArray,
			[]testStruct{
				{X: 1, Y: "a"},
				{X: 1, Y: "b"},
			},
		},
		"Array Y==`c`": {
			"Y==`c`",
			testArray,
			[]testStruct{
				{X: 3, Y: "c"},
			},
		},
		"Map X==1": {
			"X==1",
			testMap,
			map[string]testStruct{
				"one": {X: 1, Y: "a"},
				"two": {X: 1, Y: "b"},
			},
		},
		"Map Y==`c`": {
			"Y==`c`",
			testMap,
			map[string]testStruct{
				"five": {X: 3, Y: "c"},
			},
		},
	}

	t.Parallel()

	for name, tcase := range cases {
		tcase := tcase
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			flt, err := CreateFilter(tcase.expression)
			require.NoError(t, err)
			require.NotNil(t, flt)

			results, err := flt.Execute(tcase.input)
			require.NoError(t, err)
			require.Equal(t, tcase.expected, results)
		})
	}
}

func BenchmarkFilter(b *testing.B) {
	type benchCase struct {
		expression string
		input      interface{}
	}

	// The expressions used here are purposefully simple\
	// This is meant to bencharmk the filter execution timing
	// more than the boolean expression evaluation. That includes
	// handling the top level container and generating a new one.
	// The BenchmarkEvaluate function handles benchmarking most
	// of the underlying evaluation.
	cases := map[string]benchCase{
		"Slice": {
			"X==1",
			testSlice,
		},
		"Array": {
			"X==1",
			testArray,
		},
		"Map": {
			"X==1",
			testMap,
		},
	}

	for name, bcase := range cases {
		bcase := bcase
		name := name
		b.Run(name, func(b *testing.B) {
			flt, err := CreateFilter(bcase.expression)
			require.NoError(b, err)
			require.NotNil(b, flt)

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				_, err := flt.Execute(bcase.input)
				require.NoError(b, err)
			}
		})
	}
}
