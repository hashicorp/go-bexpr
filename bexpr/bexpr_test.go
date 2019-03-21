package bexpr

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockExpressionMatcher struct {
	mock.Mock
}

func (m *mockExpressionMatcher) FieldConfigurations() []FieldConfiguration {
	args := m.Called()
	return args.Get(0).([]FieldConfiguration)
}

func (m *mockExpressionMatcher) ExecuteMatcher(field Selector, op MatchOperator, value interface{}) (bool, error) {
	args := m.Called(field, op, value)
	return args.Bool(0), args.Error(1)
}

func TestEvaluation(t *testing.T) {
	t.Parallel()
	type execMatchCall struct {
		selector    Selector
		op          MatchOperator
		value       interface{}
		returnValue bool
	}

	type testCase struct {
		config     []*FieldConfiguration
		expression string
		matches    []execMatchCall
		valid      bool
		err        string
	}

	tests := map[string]testCase{
		"basic": {
			config: []*FieldConfiguration{
				&FieldConfiguration{
					Name:     "foo",
					CoerceFn: CoerceInt,
					SupportedOperations: []MatchOperator{
						MatchEqual,
					},
				},
			},
			expression: "foo == 3",
			matches: []execMatchCall{
				execMatchCall{
					selector:    Selector{"foo"},
					op:          MatchEqual,
					value:       3,
					returnValue: true,
				},
			},
			valid: true,
			err:   "",
		},
	}

	for name, tcase := range tests {
		name := name
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := new(mockExpressionMatcher)
			for _, call := range tcase.matches {
				m.On("ExecuteMatcher", call.selector, call.op, call.value).Return(call.returnValue, nil)
			}

			expr, err := Create(tcase.expression, tcase.config)
			if tcase.err != "" {
				require.Error(t, err)
				require.EqualError(t, err, tcase.err)
				require.Nil(t, expr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, expr)

				match, err := expr.Evaluate(m)
				require.NoError(t, err)
				require.Equal(t, tcase.valid, match)
			}

			m.AssertExpectations(t)
		})
	}
}

func TestFieldConfiguration(t *testing.T) {
	t.Parallel()
	type testCase struct {
		config     []*FieldConfiguration
		expression string
		err        string
	}

	tests := map[string]testCase{
		"basic": {
			config: []*FieldConfiguration{
				&FieldConfiguration{
					Name:     "foo",
					CoerceFn: CoerceInt,
					SupportedOperations: []MatchOperator{
						MatchEqual,
					},
				},
			},
			expression: "foo == 3",
			err:        "",
		},
	}

	for name, tcase := range tests {
		name := name
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			expr, err := Create(tcase.expression, tcase.config)
			if tcase.err != "" {
				require.Error(t, err)
				require.EqualError(t, err, tcase.err)
				require.Nil(t, expr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, expr)
			}
		})
	}
}
