package bexpr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		expression string
		config     EvaluatorConfig
		err        string
	}

	tests := map[string]testCase{
		"basic": {
			expression: "foo == 3",
		},
	}

	for name, tcase := range tests {
		name := name
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			expr, err := Create(tcase.expression, &tcase.config)
			if tcase.err == "" {
				require.NoError(t, err)
				require.NotNil(t, expr)
			}
		})
	}
}
