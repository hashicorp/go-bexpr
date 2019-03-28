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
	testStruct{
		X: 1,
		Y: "a",
	},
	testStruct{
		X: 1,
		Y: "b",
	},
	testStruct{
		X: 2,
		Y: "a",
	},
	testStruct{
		X: 2,
		Y: "b",
	},
	testStruct{
		X: 3,
		Y: "c",
	},
}

var testArray [5]testStruct = [5]testStruct{
	testStruct{
		X: 1,
		Y: "a",
	},
	testStruct{
		X: 1,
		Y: "b",
	},
	testStruct{
		X: 2,
		Y: "a",
	},
	testStruct{
		X: 2,
		Y: "b",
	},
	testStruct{
		X: 3,
		Y: "c",
	},
}

var testMap map[string]testStruct = map[string]testStruct{
	"one": testStruct{
		X: 1,
		Y: "a",
	},
	"two": testStruct{
		X: 1,
		Y: "b",
	},
	"three": testStruct{
		X: 2,
		Y: "a",
	},
	"four": testStruct{
		X: 2,
		Y: "b",
	},
	"five": testStruct{
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
		"Slice X==1": testCase{
			"X==1",
			testSlice,
			[]testStruct{
				testStruct{X: 1, Y: "a"},
				testStruct{X: 1, Y: "b"},
			},
		},
		"Slice Y==`c`": testCase{
			"Y==`c`",
			testSlice,
			[]testStruct{
				testStruct{X: 3, Y: "c"},
			},
		},
		"Array X==1": testCase{
			"X==1",
			testArray,
			[]testStruct{
				testStruct{X: 1, Y: "a"},
				testStruct{X: 1, Y: "b"},
			},
		},
		"Array Y==`c`": testCase{
			"Y==`c`",
			testArray,
			[]testStruct{
				testStruct{X: 3, Y: "c"},
			},
		},
		"Map X==1": testCase{
			"X==1",
			testMap,
			map[string]testStruct{
				"one": testStruct{X: 1, Y: "a"},
				"two": testStruct{X: 1, Y: "b"},
			},
		},
		"Map Y==`c`": testCase{
			"Y==`c`",
			testMap,
			map[string]testStruct{
				"five": testStruct{X: 3, Y: "c"},
			},
		},
	}

	t.Parallel()

	for name, tcase := range cases {
		tcase := tcase
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			flt, err := CreateFilter(tcase.expression, nil, tcase.input)
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
		"Slice": benchCase{
			"X==1",
			testSlice,
		},
		"Array": benchCase{
			"X==1",
			testArray,
		},
		"Map": benchCase{
			"X==1",
			testMap,
		},
	}

	for name, bcase := range cases {
		bcase := bcase
		name := name
		b.Run(name, func(b *testing.B) {
			flt, err := CreateFilter(bcase.expression, nil, bcase.input)
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
