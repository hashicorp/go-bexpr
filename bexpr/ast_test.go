package bexpr

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAST_Dump(t *testing.T) {
	type testCase struct {
		expr     Expr
		expected string
	}

	tests := map[string]testCase{
		"MatchEqual": {
			expr:     &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchEqual, Value: &Value{Raw: "baz"}},
			expected: "Equal {\n   Selector: foo.bar\n   Value: \"baz\"\n}\n",
		},
		"MatchNotEqual": {
			expr:     &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchNotEqual, Value: &Value{Raw: "baz"}},
			expected: "Not Equal {\n   Selector: foo.bar\n   Value: \"baz\"\n}\n",
		},
		"MatchIn": {
			expr:     &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchIn, Value: &Value{Raw: "baz"}},
			expected: "In {\n   Selector: foo.bar\n   Value: \"baz\"\n}\n",
		},
		"MatchNotIn": {
			expr:     &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchNotIn, Value: &Value{Raw: "baz"}},
			expected: "Not In {\n   Selector: foo.bar\n   Value: \"baz\"\n}\n",
		},
		"MatchIsEmpty": {
			expr:     &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchIsEmpty, Value: nil},
			expected: "Is Empty {\n   Selector: foo.bar\n}\n",
		},
		"MatchIsNotEmpty": {
			expr:     &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchIsNotEmpty, Value: nil},
			expected: "Is Not Empty {\n   Selector: foo.bar\n}\n",
		},
		"MatchUnknown": {
			expr:     &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchOperator(42), Value: nil},
			expected: "UNKNOWN {\n   Selector: foo.bar\n}\n",
		},
		"UnaryOpNot": {
			expr:     &UnaryExpr{Operator: UnaryOpNot, Operand: &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchIsEmpty, Value: nil}},
			expected: "Not {\n   Is Empty {\n      Selector: foo.bar\n   }\n}\n",
		},
		"UnaryOpUnknown": {
			expr:     &UnaryExpr{Operator: UnaryOperator(42), Operand: &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchIsEmpty, Value: nil}},
			expected: "UNKNOWN {\n   Is Empty {\n      Selector: foo.bar\n   }\n}\n",
		},
		"BinaryOpAnd": {
			expr: &BinaryExpr{
				Operator: BinaryOpAnd,
				Left:     &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchIsEmpty, Value: nil},
				Right:    &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchIsEmpty, Value: nil},
			},
			expected: "And {\n   Is Empty {\n      Selector: foo.bar\n   }\n   Is Empty {\n      Selector: foo.bar\n   }\n}\n",
		},
		"BinaryOpOr": {
			expr: &BinaryExpr{
				Operator: BinaryOpOr,
				Left:     &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchIsEmpty, Value: nil},
				Right:    &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchIsEmpty, Value: nil},
			},
			expected: "Or {\n   Is Empty {\n      Selector: foo.bar\n   }\n   Is Empty {\n      Selector: foo.bar\n   }\n}\n",
		},
		"BinaryOpUnknown": {
			expr: &BinaryExpr{
				Operator: BinaryOperator(42),
				Left:     &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchIsEmpty, Value: nil},
				Right:    &MatchExpr{Selector: Selector{"foo", "bar"}, Operator: MatchIsEmpty, Value: nil},
			},
			expected: "UNKNOWN {\n   Is Empty {\n      Selector: foo.bar\n   }\n   Is Empty {\n      Selector: foo.bar\n   }\n}\n",
		},
	}

	for name, tcase := range tests {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)
			tcase.expr.Dump(buf, "   ", 0)
			actual := buf.String()

			require.Equal(t, tcase.expected, actual)
		})
	}
}
