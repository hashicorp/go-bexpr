// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bexpr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEvaluator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		expression string
		err        string
	}

	tests := map[string]testCase{
		"basic": {
			expression: "foo == 3",
		},
		"default max expressions": {
			expression: "((((((((foo == 1))))))))",
			// typo in pigeon code-gen
			err: "max number of expresssions parsed",
		},
	}

	for name, tcase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			expr, err := CreateEvaluator(tcase.expression)
			if tcase.err == "" {
				require.NoError(t, err)
				require.NotNil(t, expr)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tcase.err)
			}
		})
	}
}
