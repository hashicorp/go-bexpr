package grammar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpressionParsing(t *testing.T) {
	type testCase struct {
		input    string
		expected Expr
		err      string
	}

	tests := map[string]testCase{
		"Match Equality": {
			input:    "foo == 3",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "3"}},
			err:      "",
		},
		"Match Inequality": {
			input:    "foo != xyz",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchNotEqual, Value: &Value{Raw: "xyz"}},
			err:      "",
		},
		"Match Is Empty": {
			input:    "list is empty",
			expected: &MatchExpr{Selector: Selector{"list"}, Operator: MatchIsEmpty, Value: nil},
			err:      "",
		},
		"Match Is Not Empty": {
			input:    "list is not empty",
			expected: &MatchExpr{Selector: Selector{"list"}, Operator: MatchIsNotEmpty, Value: nil},
			err:      "",
		},
		"Match In": {
			input:    "foo in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchIn, Value: &Value{Raw: "foo"}},
			err:      "",
		},
		"Match Not In": {
			input:    "foo not in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchNotIn, Value: &Value{Raw: "foo"}},
			err:      "",
		},
		"Logical Not": {
			input: "not prod in tags",
			expected: &UnaryExpr{
				Operator: UnaryOpNot,
				Operand:  &MatchExpr{Selector: Selector{"tags"}, Operator: MatchIn, Value: &Value{Raw: "prod"}},
			},
			err: "",
		},
		"Logical And": {
			input: "port != 80 and port != 8080",
			expected: &BinaryExpr{
				Operator: BinaryOpAnd,
				Left:     &MatchExpr{Selector: Selector{"port"}, Operator: MatchNotEqual, Value: &Value{Raw: "80"}},
				Right:    &MatchExpr{Selector: Selector{"port"}, Operator: MatchNotEqual, Value: &Value{Raw: "8080"}},
			},
			err: "",
		},
		"Logical Or": {
			input: "port == 80 or port == 443",
			expected: &BinaryExpr{
				Operator: BinaryOpOr,
				Left:     &MatchExpr{Selector: Selector{"port"}, Operator: MatchEqual, Value: &Value{Raw: "80"}},
				Right:    &MatchExpr{Selector: Selector{"port"}, Operator: MatchEqual, Value: &Value{Raw: "443"}},
			},
			err: "",
		},
		"Double Quoted Value (Equal)": {
			input:    "foo == \"bar\"",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "bar"}},
			err:      "",
		},
		"Double Quoted Value (Not Equal)": {
			input:    "foo != \"bar\"",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchNotEqual, Value: &Value{Raw: "bar"}},
			err:      "",
		},
		"Double Quoted Value (In)": {
			input:    "\"foo\" in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchIn, Value: &Value{Raw: "foo"}},
			err:      "",
		},
		"Double Quoted Value (Not In)": {
			input:    "\"foo\" not in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchNotIn, Value: &Value{Raw: "foo"}},
			err:      "",
		},
		"Backtick Quoted Value (Equal)": {
			input:    "foo == `bar`",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "bar"}},
			err:      "",
		},
		"Backtick Quoted Value (Not Equal)": {
			input:    "foo != `bar`",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchNotEqual, Value: &Value{Raw: "bar"}},
			err:      "",
		},
		"Backtick Quoted Value (In)": {
			input:    "`foo` in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchIn, Value: &Value{Raw: "foo"}},
			err:      "",
		},
		"Backtick Quoted Value (Not In)": {
			input:    "`foo` not in bar",
			expected: &MatchExpr{Selector: Selector{"bar"}, Operator: MatchNotIn, Value: &Value{Raw: "foo"}},
			err:      "",
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
			err: "",
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
			err: "",
		},
		"Extra Whitespace (Equal)": {
			input:    "\t\r\n  foo \t\r\n == \t\r\n x \t\r\n",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "x"}},
			err:      "",
		},
		"Extra Whitespace (Not Equal)": {
			input:    "\t\r\n  foo \t\r\n != \t\r\n x \t\r\n",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchNotEqual, Value: &Value{Raw: "x"}},
			err:      "",
		},
		"Extra Whitespace (In)": {
			input:    "\t\r\n  foo \t\r\n in \t\r\n x \t\r\n",
			expected: &MatchExpr{Selector: Selector{"x"}, Operator: MatchIn, Value: &Value{Raw: "foo"}},
			err:      "",
		},
		"Extra Whitespace (Not In)": {
			input:    "\t\r\n  foo \t\r\n not \t\r\n in \t\r\n x \t\r\n",
			expected: &MatchExpr{Selector: Selector{"x"}, Operator: MatchNotIn, Value: &Value{Raw: "foo"}},
			err:      "",
		},
		"Extra Whitespace (Is Empty)": {
			input:    "\t\r\n  foo \t\r\n is \t\r\n empty \t\r\n",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchIsEmpty, Value: nil},
			err:      "",
		},
		"Extra Whitespace (Is Not Empty)": {
			input:    "\t\r\n  foo \t\r\n is \t\r\n not \t\r\n empty \t\r\n",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchIsNotEmpty, Value: nil},
			err:      "",
		},
		"Extra Whitespace (Not)": {
			input: "\t\r\n not \t\r\n  foo \t\r\n in \t\r\n x \t\r\n",
			expected: &UnaryExpr{
				Operator: UnaryOpNot,
				Operand:  &MatchExpr{Selector: Selector{"x"}, Operator: MatchIn, Value: &Value{Raw: "foo"}},
			},
			err: "",
		},
		"Extra Whitespace (And)": {
			input: "\t\r\n foo \t\r\n == \t\r\n x \t\r\n and \t\r\n y \t\r\n is \t\r\n empty \t\r\n",
			expected: &BinaryExpr{
				Operator: BinaryOpAnd,
				Left:     &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "x"}},
				Right:    &MatchExpr{Selector: Selector{"y"}, Operator: MatchIsEmpty, Value: nil},
			},
			err: "",
		},
		"Extra Whitespace (Or)": {
			input: "\t\r\n foo \t\r\n == \t\r\n x \t\r\n or \t\r\n y \t\r\n is \t\r\n empty \t\r\n",
			expected: &BinaryExpr{
				Operator: BinaryOpOr,
				Left:     &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "x"}},
				Right:    &MatchExpr{Selector: Selector{"y"}, Operator: MatchIsEmpty, Value: nil},
			},
			err: "",
		},
		"Extra Whitespace (Parentheses)": {
			input:    "\t\r\n ( \t\r\n foo \t\r\n == \t\r\n x \t\r\n ) \t\r\n",
			expected: &MatchExpr{Selector: Selector{"foo"}, Operator: MatchEqual, Value: &Value{Raw: "x"}},
			err:      "",
		},
		"Selector Path": {
			input:    "`environment` in foo.bar[\"meta\"].tags[\t`ENV` ]",
			expected: &MatchExpr{Selector: Selector{"foo", "bar", "meta", "tags", "ENV"}, Operator: MatchIn, Value: &Value{Raw: "environment"}},
			err:      "",
		},
		"Selector All Indexes": {
			input:    `environment in foo["bar"]["meta"]["tags"]["ENV"]`,
			expected: &MatchExpr{Selector: Selector{"foo", "bar", "meta", "tags", "ENV"}, Operator: MatchIn, Value: &Value{Raw: "environment"}},
			err:      "",
		},
		"Selector All Dotted": {
			input:    "environment in foo.bar.meta.tags.ENV",
			expected: &MatchExpr{Selector: Selector{"foo", "bar", "meta", "tags", "ENV"}, Operator: MatchIn, Value: &Value{Raw: "environment"}},
			err:      "",
		},
		// selectors can contain almost any character set when index expressions are used
		// This includes whitespace, hyphens, unicode, etc.
		"Selector Index Chars": {
			input:    "environment in foo[\"abc-def ghi åß∂ƒ\"]",
			expected: &MatchExpr{Selector: Selector{"foo", "abc-def ghi åß∂ƒ"}, Operator: MatchIn, Value: &Value{Raw: "environment"}},
			err:      "",
		},
		"Unterminated String Literal 1": {
			input:    "foo == \"12x",
			expected: nil,
			err:      "1:12 (11): rule \"string\": Unterminated string literal",
		},
		"Unterminated String Literal 2": {
			input:    "foo == `12x",
			expected: nil,
			err:      "1:12 (11): rule \"string\": Unterminated string literal",
		},
		"Invalid Integer": {
			input:    "foo == 3x",
			expected: nil,
			err:      "1:9 (8): rule \"integer\": Invalid integer literal",
		},
		"Invalid Index Key": {
			input:    "foo[3] == abc",
			expected: nil,
			err:      "1:5 (4): rule \"index\": Invalid index",
		},
		"Unclosed Index Expression 1": {
			input:    "x in foo[\"abc\"",
			expected: nil,
			err:      "1:15 (14): rule \"index\": Unclosed index expression",
		},
		"Unclosed Index Expression 2": {
			input:    "foo[\"abc\" == 3",
			expected: nil,
			err:      "1:11 (10): rule \"index\": Unclosed index expression",
		},
		"Invalid Selector 1": {
			input:    "x in 32",
			expected: nil,
			err:      "1:6 (5): rule \"match\": Invalid selector",
		},
		"Invalid Selector 2": {
			input:    "32 == 32",
			expected: nil,
			err:      `1:4 (3): no match found, expected: "in", "not" or [ \t\r\n]`,
		},
		"Invalid Selector 3": {
			input:    "32 is empty",
			expected: nil,
			err:      `1:4 (3): no match found, expected: "in", "not" or [ \t\r\n]`,
		},
		"Junk at the end 1": {
			input:    "x in foo abc",
			expected: nil,
			err:      `1:10 (9): no match found, expected: "and", "or", [ \t\r\n] or EOF`,
		},
		"Junk at the end 2": {
			input:    "x in foo and ",
			expected: nil,
			err:      "1:14 (13): no match found, expected: \"(\", \"-\", \"\\\"\", \"`\", \"not\", [ \\t\\r\\n], [1-9] or [a-zA-Z]",
		},
		"Junk at the end 3": {
			input:    "x in foo or ",
			expected: nil,
			err:      "1:13 (12): no match found, expected: \"(\", \"-\", \"\\\"\", \"`\", \"not\", [ \\t\\r\\n], [1-9] or [a-zA-Z]",
		},
		"Junk at the end 4": {
			input:    "x in foo or not ",
			expected: nil,
			err:      "1:17 (16): no match found, expected: \"!=\", \"(\", \"-\", \"==\", \"\\\"\", \"`\", \"in\", \"is\", \"not\", [ \\t\\r\\n], [1-9] or [a-zA-Z]",
		},
	}

	for name, tcase := range tests {
		tcase := tcase
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			raw, err := Parse("", []byte(tcase.input))
			if tcase.err != "" {
				require.Error(t, err)
				require.EqualError(t, err, tcase.err)
				require.Nil(t, raw)
			} else {
				require.NoError(t, err)
				require.NotNil(t, raw)
				require.Equal(t, tcase.expected, raw)
			}
		})
	}
}
