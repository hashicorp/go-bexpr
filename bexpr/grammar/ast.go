package grammar

import (
	"fmt"
	"io"
	"strings"
)

//go:generate pigeon -o grammar.go grammar.peg
//go:generate goimports -w grammar.go

type Expr interface {
	ExprNode() // no-op function to indicate that the node implements the Expr interface

	Dump(w io.Writer, indent string, level int)
}

type UnaryOperator int

const (
	UnaryOpNot UnaryOperator = iota
)

func (op UnaryOperator) String() string {
	switch op {
	case UnaryOpNot:
		return "Not"
	default:
		return "UNKNOWN"
	}
}

type BinaryOperator int

const (
	BinaryOpAnd BinaryOperator = iota
	BinaryOpOr
)

func (op BinaryOperator) String() string {
	switch op {
	case BinaryOpAnd:
		return "And"
	case BinaryOpOr:
		return "Or"
	default:
		return "UNKNOWN"
	}
}

type MatchOperator int

const (
	MatchEqual MatchOperator = iota
	MatchNotEqual
	MatchIn
	MatchNotIn
	MatchIsEmpty
	MatchIsNotEmpty
)

func (op MatchOperator) String() string {
	switch op {
	case MatchEqual:
		return "Equal"
	case MatchNotEqual:
		return "Not Equal"
	case MatchIn:
		return "In"
	case MatchNotIn:
		return "Not In"
	case MatchIsEmpty:
		return "Is Empty"
	case MatchIsNotEmpty:
		return "Is Not Empty"
	default:
		return "UNKNOWN"
	}
}

type Value struct {
	Raw       string
	Converted interface{}
}

type UnaryExpr struct {
	Operator UnaryOperator
	Operand  Expr
}

type BinaryExpr struct {
	Left     Expr
	Operator BinaryOperator
	Right    Expr
}

type Selector []string

type MatchExpr struct {
	Selector Selector
	Operator MatchOperator
	Value    *Value
}

func (expr *UnaryExpr) ExprNode()  {}
func (expr *BinaryExpr) ExprNode() {}
func (expr *MatchExpr) ExprNode()  {}

func (expr *UnaryExpr) Dump(w io.Writer, indent string, level int) {
	localIndent := strings.Repeat(indent, level)
	fmt.Fprintf(w, "%s%s {\n", localIndent, expr.Operator.String())
	expr.Operand.Dump(w, indent, level+1)
	fmt.Fprintf(w, "%s}\n", localIndent)
}

func (expr *BinaryExpr) Dump(w io.Writer, indent string, level int) {
	localIndent := strings.Repeat(indent, level)
	fmt.Fprintf(w, "%s%s {\n", localIndent, expr.Operator.String())
	expr.Left.Dump(w, indent, level+1)
	expr.Right.Dump(w, indent, level+1)
	fmt.Fprintf(w, "%s}\n", localIndent)
}

func (expr *MatchExpr) Dump(w io.Writer, indent string, level int) {
	switch expr.Operator {
	case MatchEqual, MatchNotEqual, MatchIn, MatchNotIn:
		fmt.Fprintf(w, "%[1]s%[3]s {\n%[2]sSelector: %[4]v\n%[2]sValue: %[5]q\n%[1]s}\n", strings.Repeat(indent, level), strings.Repeat(indent, level+1), expr.Operator.String(), expr.Selector, expr.Value.Raw)
	default:
		fmt.Fprintf(w, "%[1]s%[3]s {\n%[2]sSelector: %[4]v\n%[1]s}\n", strings.Repeat(indent, level), strings.Repeat(indent, level+1), expr.Operator.String(), expr.Selector)
	}

}
