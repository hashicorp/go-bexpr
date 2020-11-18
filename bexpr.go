// bexpr is an implementation of a generic boolean expression evaluator.
// The general goal is to be able to evaluate some expression against some
// arbitrary data and get back a boolean of whether or not the data
// was matched by the expression
package bexpr

type Evaluator struct {
	// The syntax tree
	ast Expression

	// A few configurations for extra validation of the AST
	config EvaluatorConfig
}

// Extra configuration used to perform further validation on a parsed
// expression. Currently this does not hold any fields, but it avoids changing
// the function signature.
//
// TODO: Remove this? Perhaps in favor of an Options approach for the calls that
// need it?
type EvaluatorConfig struct {
}

func CreateEvaluator(expression string, config *EvaluatorConfig) (*Evaluator, error) {
	ast, err := Parse("", []byte(expression))
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = &EvaluatorConfig{}
	}

	eval := &Evaluator{
		ast:    ast.(Expression),
		config: *config,
	}

	return eval, nil
}

func (eval *Evaluator) Evaluate(datum interface{}) (bool, error) {
	return evaluate(eval.ast, datum)
}
