package bexpr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterLangParsing(t *testing.T) {
	type testCase struct {
		input    string
		expected Expr
		err      error
	}

	tests := map[string]testCase{
		"Match Equality": {
			input:    "foo == 3",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "3"}},
			err:      nil,
		},
		"Match Inequality": {
			input:    "foo != xyz",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchNotEqual, Value: &Value{Raw: "xyz"}},
			err:      nil,
		},
		"Match Is Empty": {
			input:    "list is empty",
			expected: &MatchExpr{Selector: Selector{"list"}, Operator: MatchIsEmpty, Value: nil},
			err:      nil,
		},
		"Match Is Not Empty": {
			input:    "list is not empty",
			expected: &MatchExpr{Selector: Selector{"list"}, Operator: MatchIsNotEmpty, Value: nil},
			err:      nil,
		},
		"Match In": {
			input:    "foo in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchIn, Value: &Value{Raw: "foo"}},
			err:      nil,
		},
		"Match Not In": {
			input:    "foo not in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchNotIn, Value: &Value{Raw: "foo"}},
			err:      nil,
		},
		"Logical Not": {
			input: "not prod in tags",
			expected: &UnaryExpr{
				Operator: UnaryOpNot,
				Operand:  &MatchExpr{Selector: Selector{"tags"}, Operator: MatchIn, Value: &Value{Raw: "prod"}},
			},
			err: nil,
		},
		"Logical And": {
			input: "port != 80 and port != 8080",
			expected: &BinaryExpr{
				Operator: BinaryOpAnd,
				Left:     &MatchExpr{Selector: Selector{"port"}, Operator: MatchNotEqual, Value: &Value{Raw: "80"}},
				Right:    &MatchExpr{Selector: Selector{"port"}, Operator: MatchNotEqual, Value: &Value{Raw: "8080"}},
			},
			err: nil,
		},
		"Logical Or": {
			input: "port == 80 or port == 443",
			expected: &BinaryExpr{
				Operator: BinaryOpOr,
				Left:     &MatchExpr{Selector: Selector{"port"}, Operator: MatchEqual, Value: &Value{Raw: "80"}},
				Right:    &MatchExpr{Selector: Selector{"port"}, Operator: MatchEqual, Value: &Value{Raw: "443"}},
			},
			err: nil,
		},
	}

	for name, tcase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			raw, err := Parse(name, []byte(tcase.input))
			if tcase.err != nil {
				require.Error(t, err)
				require.Equal(t, tcase.err, err)
				require.Nil(t, raw)
			} else {
				require.NoError(t, err)
				require.NotNil(t, raw)
				require.Equal(t, tcase.expected, raw)
			}
		})
	}
}
