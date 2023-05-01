// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package grammar

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAST_Dump(t *testing.T) {
	t.Parallel()
	type testCase struct {
		expr     Expression
		expected string
	}

	tests := map[string]testCase{
		"MatchEqual": {
			expr:     &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchEqual, Value: &MatchValue{Raw: "baz"}},
			expected: "Equal {\n   Selector: foo.bar\n   Value: \"baz\"\n}\n",
		},
		"MatchNotEqual": {
			expr:     &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchNotEqual, Value: &MatchValue{Raw: "baz"}},
			expected: "Not Equal {\n   Selector: foo.bar\n   Value: \"baz\"\n}\n",
		},
		"MatchIn": {
			expr:     &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchIn, Value: &MatchValue{Raw: "baz"}},
			expected: "In {\n   Selector: foo.bar\n   Value: \"baz\"\n}\n",
		},
		"MatchNotIn": {
			expr:     &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchNotIn, Value: &MatchValue{Raw: "baz"}},
			expected: "Not In {\n   Selector: foo.bar\n   Value: \"baz\"\n}\n",
		},
		"MatchIsEmpty": {
			expr:     &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchIsEmpty, Value: nil},
			expected: "Is Empty {\n   Selector: foo.bar\n}\n",
		},
		"MatchIsNotEmpty": {
			expr:     &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchIsNotEmpty, Value: nil},
			expected: "Is Not Empty {\n   Selector: foo.bar\n}\n",
		},
		"MatchUnknown": {
			expr:     &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchOperator(42), Value: nil},
			expected: "UNKNOWN {\n   Selector: foo.bar\n}\n",
		},
		"UnaryOpNot": {
			expr:     &UnaryExpression{Operator: UnaryOpNot, Operand: &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchIsEmpty, Value: nil}},
			expected: "Not {\n   Is Empty {\n      Selector: foo.bar\n   }\n}\n",
		},
		"UnaryOpUnknown": {
			expr:     &UnaryExpression{Operator: UnaryOperator(42), Operand: &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchIsEmpty, Value: nil}},
			expected: "UNKNOWN {\n   Is Empty {\n      Selector: foo.bar\n   }\n}\n",
		},
		"BinaryOpAnd": {
			expr: &BinaryExpression{
				Operator: BinaryOpAnd,
				Left:     &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchIsEmpty, Value: nil},
				Right:    &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchIsEmpty, Value: nil},
			},
			expected: "And {\n   Is Empty {\n      Selector: foo.bar\n   }\n   Is Empty {\n      Selector: foo.bar\n   }\n}\n",
		},
		"BinaryOpOr": {
			expr: &BinaryExpression{
				Operator: BinaryOpOr,
				Left:     &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchIsEmpty, Value: nil},
				Right:    &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchIsEmpty, Value: nil},
			},
			expected: "Or {\n   Is Empty {\n      Selector: foo.bar\n   }\n   Is Empty {\n      Selector: foo.bar\n   }\n}\n",
		},
		"BinaryOpUnknown": {
			expr: &BinaryExpression{
				Operator: BinaryOperator(42),
				Left:     &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchIsEmpty, Value: nil},
				Right:    &MatchExpression{Selector: Selector{Type: SelectorTypeBexpr, Path: []string{"foo", "bar"}}, Operator: MatchIsEmpty, Value: nil},
			},
			expected: "UNKNOWN {\n   Is Empty {\n      Selector: foo.bar\n   }\n   Is Empty {\n      Selector: foo.bar\n   }\n}\n",
		},
	}

	for name, tcase := range tests {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)
			tcase.expr.ExpressionDump(buf, "   ", 0)
			actual := buf.String()

			require.Equal(t, tcase.expected, actual)
		})
	}
}
