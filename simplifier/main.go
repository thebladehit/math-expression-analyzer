package simplifier

import (
	"fmt"
	"strconv"

	"pzcs-lab1/an"
	"pzcs-lab1/tree"
)

type Simplifier interface {
	Simplify(node *tree.Node) (*tree.Node, []string)
}

type simplifier struct {
	log []string
}

func NewSimplifier() Simplifier {
	return &simplifier{}
}

func (s *simplifier) Simplify(node *tree.Node) (*tree.Node, []string) {
	s.log = nil
	return s.simplify(node), s.log
}

func (s *simplifier) simplify(n *tree.Node) *tree.Node {
	if n == nil {
		return nil
	}

	n.Left = s.simplify(n.Left)
	n.Right = s.simplify(n.Right)

	switch n.Type {
	case tree.NodeUM:
		if v, ok := constValue(n.Left); ok {
			return s.replace(n, newConst(-v, n), "обчислення константи")
		}
	case tree.NodeOP:
		return s.simplifyOp(n)
	}
	return n
}

func (s *simplifier) simplifyOp(n *tree.Node) *tree.Node {
	l, r, op := n.Left, n.Right, n.Token.Text
	lv, lConst := constValue(l)
	rv, rConst := constValue(r)

	switch {
	case lConst && rConst:
		if v, ok := calculate(op, lv, rv); ok {
			return s.replace(n, newConst(v, n), "обчислення констант")
		}
	case op == "*" && (isConst(l, 0) || isConst(r, 0)):
		return s.replace(n, newConst(0, n), "множення на 0")
	case op == "*" && isConst(r, 1):
		return s.replace(n, l, "множення на 1")
	case op == "*" && isConst(l, 1):
		return s.replace(n, r, "множення на 1")
	case op == "/" && isConst(r, 1):
		return s.replace(n, l, "ділення на 1")
	case op == "/" && isConst(l, 0):
		return s.replace(n, newConst(0, n), "ділення 0 на вираз (якщо дільник ≠ 0)")
	case (op == "+" || op == "-") && isConst(r, 0):
		return s.replace(n, l, "додавання/віднімання 0")
	case op == "+" && isConst(l, 0):
		return s.replace(n, r, "додавання 0")
	}
	return n
}

func (s *simplifier) replace(old, new *tree.Node, reason string) *tree.Node {
	s.log = append(s.log, fmt.Sprintf("%s: %s → %s", reason, old, new))
	return new
}

func calculate(op string, a, b float64) (float64, bool) {
	switch op {
	case "+":
		return a + b, true
	case "-":
		return a - b, true
	case "*":
		return a * b, true
	case "/":
		if b == 0 {
			return 0, false
		}
		return a / b, true
	}
	return 0, false
}

func constValue(n *tree.Node) (float64, bool) {
	if n == nil || n.Type != tree.NodeCONST {
		return 0, false
	}
	v, err := strconv.ParseFloat(n.Token.Text, 64)
	return v, err == nil
}

func isConst(n *tree.Node, want float64) bool {
	v, ok := constValue(n)
	return ok && v == want
}

func newConst(v float64, at *tree.Node) *tree.Node {
	text := strconv.FormatFloat(v, 'g', -1, 64)
	return tree.NewNode(an.NewToken(an.TokenCONST, text, at.Token.Pos))
}
