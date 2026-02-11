package bexpr

import (
	"testing"

	"github.com/hashicorp/go-bexpr/grammar"
)

func TestContainsVsIn(t *testing.T) {
	claims := map[string]any{
		"userinfo": map[string]any{
			"groups": "totallynotanadmin",
			"email":  "admin@company.com",
		},
	}

	tests := []struct {
		name   string
		filter string
		expect bool
	}{
		{
			name:   "in does not find admin in totallynotanadmin",
			filter: `"admin" in "/userinfo/groups"`,
			expect: false,
		},
		{
			name:   "contains finds admin in totallynotanadmin",
			filter: `"/userinfo/groups" contains "admin"`,
			expect: true,
		},
		{
			name:   "contains on email does substring match",
			filter: `"/userinfo/email" contains "@company.com"`,
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, _ := grammar.Parse("", []byte(tt.filter))
			expr := ast.(*grammar.MatchExpression)

			eval, _ := CreateEvaluator(tt.filter)
			result, _ := eval.Evaluate(claims)

			t.Logf("\n  Filter:   %s\n  Operator: %s\n  Result:   %v [expect: %v]",
				tt.filter, expr.Operator, result, tt.expect)

			if result != tt.expect {
				t.Errorf("expected %v, got %v", tt.expect, result)
			}
		})
	}
}
