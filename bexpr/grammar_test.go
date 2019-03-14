package bexpr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpressionParsing(t *testing.T) {
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
		"Double Quoted Value (Equal)": {
			input:    "foo == \"bar\"",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "bar"}},
			err:      nil,
		},
		"Double Quoted Value (Not Equal)": {
			input:    "foo != \"bar\"",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchNotEqual, Value: &Value{Raw: "bar"}},
			err:      nil,
		},
		"Double Quoted Value (In)": {
			input:    "\"foo\" in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchIn, Value: &Value{Raw: "foo"}},
			err:      nil,
		},
		"Double Quoted Value (Not In)": {
			input:    "\"foo\" not in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchNotIn, Value: &Value{Raw: "foo"}},
			err:      nil,
		},
		"Backtick Quoted Value (Equal)": {
			input:    "foo == `bar`",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "bar"}},
			err:      nil,
		},
		"Backtick Quoted Value (Not Equal)": {
			input:    "foo != `bar`",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchNotEqual, Value: &Value{Raw: "bar"}},
			err:      nil,
		},
		"Backtick Quoted Value (In)": {
			input:    "`foo` in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchIn, Value: &Value{Raw: "foo"}},
			err:      nil,
		},
		"Backtick Quoted Value (Not In)": {
			input:    "`foo` not in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchNotIn, Value: &Value{Raw: "foo"}},
			err:      nil,
		},
		// This is standard boolean expression precedence
		//
		// 1 - not
		// 2 - and
		// 3 - or
		"Logical Operator Precedence": {
			input: "x in foo and not str == something or list is empty",
			expected: &BinaryExpr{
				Operator: BinaryOpOr,
				Left: &BinaryExpr{
					Operator: BinaryOpAnd,
					Left:     &MatchExpr{Selector: Selector{"foo"}, Operator: MatchIn, Value: &Value{Raw: "x"}},
					Right: &UnaryExpr{
						Operator: UnaryOpNot,
						Operand:  &MatchExpr{Selector: Selector{"str"}, Operator: MatchEqual, Value: &Value{Raw: "something"}},
					},
				},
				Right: &MatchExpr{Selector: Selector{"list"}, Operator: MatchIsEmpty, Value: nil},
			},
			err: nil,
		},
		// not in the absence of parentheses would normally get applied
		// to a MatchExpr
		//
		// or operators normal are the last to be applied but here they
		// happen earlier
		"Logical Operator Precedence (Parenthesized)": {
			input: "x in foo and not (str == something or list is empty)",
			expected: &BinaryExpr{
				Operator: BinaryOpAnd,
				Left:     &MatchExpr{Selector: Selector{"foo"}, Operator: MatchIn, Value: &Value{Raw: "x"}},
				Right: &UnaryExpr{
					Operator: UnaryOpNot,
					Operand: &BinaryExpr{
						Operator: BinaryOpOr,
						Left:     &MatchExpr{Selector: Selector{"str"}, Operator: MatchEqual, Value: &Value{Raw: "something"}},
						Right:    &MatchExpr{Selector: Selector{"list"}, Operator: MatchIsEmpty, Value: nil},
					},
				},
			},
			err: nil,
		},
		"Extra Whitespace (Equal)": {
			input:    "\t\r\n  foo \t\r\n == \t\r\n x \t\r\n",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "x"}},
			err:      nil,
		},
		"Extra Whitespace (Not Equal)": {
			input:    "\t\r\n  foo \t\r\n != \t\r\n x \t\r\n",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchNotEqual, Value: &Value{Raw: "x"}},
			err:      nil,
		},
		"Extra Whitespace (In)": {
			input:    "\t\r\n  foo \t\r\n in \t\r\n x \t\r\n",
			expected: &MatchExpr{Selector: Selector{"x"}, Operator: MatchIn, Value: &Value{Raw: "foo"}},
			err:      nil,
		},
		"Extra Whitespace (Not In)": {
			input:    "\t\r\n  foo \t\r\n not \t\r\n in \t\r\n x \t\r\n",
			expected: &MatchExpr{Selector: Selector{"x"}, Operator: MatchNotIn, Value: &Value{Raw: "foo"}},
			err:      nil,
		},
		"Extra Whitespace (Is Empty)": {
			input:    "\t\r\n  foo \t\r\n is \t\r\n empty \t\r\n",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchIsEmpty, Value: nil},
			err:      nil,
		},
		"Extra Whitespace (Is Not Empty)": {
			input:    "\t\r\n  foo \t\r\n is \t\r\n not \t\r\n empty \t\r\n",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchIsNotEmpty, Value: nil},
			err:      nil,
		},
		"Extra Whitespace (Not)": {
			input: "\t\r\n not \t\r\n  foo \t\r\n in \t\r\n x \t\r\n",
			expected: &UnaryExpr{
				Operator: UnaryOpNot,
				Operand:  &MatchExpr{Selector: Selector{"x"}, Operator: MatchIn, Value: &Value{Raw: "foo"}},
			},
			err: nil,
		},
		"Extra Whitespace (And)": {
			input: "\t\r\n foo \t\r\n == \t\r\n x \t\r\n and \t\r\n y \t\r\n is \t\r\n empty \t\r\n",
			expected: &BinaryExpr{
				Operator: BinaryOpAnd,
				Left:     &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "x"}},
				Right:    &MatchExpr{Selector: Selector{"y"}, Operator: MatchIsEmpty, Value: nil},
			},
			err: nil,
		},
		"Extra Whitespace (Or)": {
			input: "\t\r\n foo \t\r\n == \t\r\n x \t\r\n or \t\r\n y \t\r\n is \t\r\n empty \t\r\n",
			expected: &BinaryExpr{
				Operator: BinaryOpOr,
				Left:     &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "x"}},
				Right:    &MatchExpr{Selector: Selector{"y"}, Operator: MatchIsEmpty, Value: nil},
			},
			err: nil,
		},
		"Extra Whitespace (Parentheses)": {
			input:    "\t\r\n ( \t\r\n foo \t\r\n == \t\r\n x \t\r\n ) \t\r\n",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "x"}},
			err:      nil,
		},
		"Selector Path": {
			input:    "`environment` in foo.bar[\"meta\"].tags[\t`ENV` ]",
			expected: &MatchExpr{Selector: Selector{"foo", "bar", "meta", "tags", "ENV"}, Operator: MatchIn, Value: &Value{Raw: "environment"}},
			err:      nil,
		},
		"Selector All Indexes": {
			input:    `environment in foo["bar"]["meta"]["tags"]["ENV"]`,
			expected: &MatchExpr{Selector: Selector{"foo", "bar", "meta", "tags", "ENV"}, Operator: MatchIn, Value: &Value{Raw: "environment"}},
			err:      nil,
		},
		"Selector All Dotted": {
			input:    "environment in foo.bar.meta.tags.ENV",
			expected: &MatchExpr{Selector: Selector{"foo", "bar", "meta", "tags", "ENV"}, Operator: MatchIn, Value: &Value{Raw: "environment"}},
			err:      nil,
		},
		// selectors can contain almost any character set when index expressions are used
		// This includes whitespace, hyphens, unicode, etc.
		"Selector Index Chars": {
			input:    "environment in foo[\"abc-def ghi åß∂ƒ\"]",
			expected: &MatchExpr{Selector: Selector{"foo", "abc-def ghi åß∂ƒ"}, Operator: MatchIn, Value: &Value{Raw: "environment"}},
			err:      nil,
		},
	}

	for name, tcase := range tests {
		tcase := tcase
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
