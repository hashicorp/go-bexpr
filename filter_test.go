package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	X int
	Y string
}

func TestFilterSlice(t *testing.T) {
	t.Parallel()
	data := []testStruct{
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

	t.Run("X==1", func(t *testing.T) {
		t.Parallel()
		expected := []testStruct{
			testStruct{
				X: 1,
				Y: "a",
			},
			testStruct{
				X: 1,
				Y: "b",
			},
		}
		flt, err := Create("X==1", nil, (*testStruct)(nil))
		require.NoError(t, err)
		require.NotNil(t, flt)

		results, err := flt.Execute(data)
		require.NoError(t, err)
		require.Equal(t, expected, results)
	})

	t.Run("Y==`c`", func(t *testing.T) {
		t.Parallel()

		expected := []testStruct{
			testStruct{
				X: 3,
				Y: "c",
			},
		}
		flt, err := Create("Y==`c`", nil, (*testStruct)(nil))
		require.NoError(t, err)
		require.NotNil(t, flt)

		results, err := flt.Execute(data)
		require.NoError(t, err)
		require.Equal(t, expected, results)
	})
}

func TestFilterArray(t *testing.T) {
	t.Parallel()
	data := [5]testStruct{
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

	t.Run("X==1", func(t *testing.T) {
		t.Parallel()
		expected := []testStruct{
			testStruct{
				X: 1,
				Y: "a",
			},
			testStruct{
				X: 1,
				Y: "b",
			},
		}
		flt, err := Create("X==1", nil, (*testStruct)(nil))
		require.NoError(t, err)
		require.NotNil(t, flt)

		results, err := flt.Execute(data)
		require.NoError(t, err)
		require.Equal(t, expected, results)
	})

	t.Run("Y==`c`", func(t *testing.T) {
		t.Parallel()

		expected := []testStruct{
			testStruct{
				X: 3,
				Y: "c",
			},
		}
		flt, err := Create("Y==`c`", nil, (*testStruct)(nil))
		require.NoError(t, err)
		require.NotNil(t, flt)

		results, err := flt.Execute(data)
		require.NoError(t, err)
		require.Equal(t, expected, results)
	})
}

func TestFilterMap(t *testing.T) {
	t.Parallel()
	data := map[string]testStruct{
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

	t.Run("X==1", func(t *testing.T) {
		t.Parallel()
		expected := map[string]testStruct{
			"one": testStruct{
				X: 1,
				Y: "a",
			},
			"two": testStruct{
				X: 1,
				Y: "b",
			},
		}
		flt, err := Create("X==1", nil, (*testStruct)(nil))
		require.NoError(t, err)
		require.NotNil(t, flt)

		results, err := flt.Execute(data)
		require.NoError(t, err)
		require.Equal(t, expected, results)
	})

	t.Run("Y==`c`", func(t *testing.T) {
		t.Parallel()

		expected := map[string]testStruct{
			"five": testStruct{
				X: 3,
				Y: "c",
			},
		}
		flt, err := Create("Y==`c`", nil, (*testStruct)(nil))
		require.NoError(t, err)
		require.NotNil(t, flt)

		results, err := flt.Execute(data)
		require.NoError(t, err)
		require.Equal(t, expected, results)
	})
}
