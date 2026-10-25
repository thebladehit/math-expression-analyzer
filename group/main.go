package group

import "pzcs-lab1/tree"

type Item struct {
	Node    *tree.Node
	Inverse bool
}

type Group struct {
	Direct  string
	Inverse string
}

var (
	Additive       = Group{"+", "-"}
	Multiplicative = Group{"*", "/"}
)

func Of(n *tree.Node) Group {
	if n.Token.Text == "*" || n.Token.Text == "/" {
		return Multiplicative
	}
	return Additive
}

func Collect(n *tree.Node, g Group, inverse bool) []Item {
	if n.Type == tree.NodeOP && (n.Token.Text == g.Direct || n.Token.Text == g.Inverse) {
		rightInv := inverse
		if n.Token.Text == g.Inverse {
			rightInv = !rightInv
		}
		return append(Collect(n.Left, g, inverse), Collect(n.Right, g, rightInv)...)
	}
	if n.Type == tree.NodeUM && g == Additive {
		return Collect(n.Left, g, !inverse)
	}
	return []Item{{Node: n, Inverse: inverse}}
}
