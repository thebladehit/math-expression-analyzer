package tree

import "strings"

func ToInfix(n *Node) string {
	switch n.Type {
	case NodeVAR, NodeCONST:
		return n.Token.Text
	case NodeFN:
		return n.Token.Text + "(" + ToInfix(n.Left) + ")"
	case NodeUM:
		return "-" + wrapIf(n.Left, priority(n.Left) <= priority(n))
	}

	p := priority(n)
	left := wrapIf(n.Left, priority(n.Left) < p)
	right := wrapIf(n.Right, priority(n.Right) <= p || isNegative(n.Right))
	return left + n.Token.Text + right
}

func priority(n *Node) int {
	switch {
	case isNegative(n):
		return 3
	case n.Type == NodeOP && (n.Token.Text == "+" || n.Token.Text == "-"):
		return 1
	case n.Type == NodeOP:
		return 2
	}
	return 4
}

func isNegative(n *Node) bool {
	return n.Type == NodeUM || (n.Type == NodeCONST && strings.HasPrefix(n.Token.Text, "-"))
}

func wrapIf(n *Node, wrap bool) string {
	if wrap {
		return "(" + ToInfix(n) + ")"
	}
	return ToInfix(n)
}
