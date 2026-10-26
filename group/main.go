package group

import (
	"pzcs-lab1/an"
	"pzcs-lab1/tree"
)

type Item struct {
	Node    *tree.Node
	Inverse bool
}

func (it Item) String() string {
	sign := "+"
	if it.Inverse {
		sign = "~"
	}
	return sign + tree.ToInfix(it.Node)
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

func Negate(n *tree.Node, g Group) *tree.Node {
	if g == Additive {
		um := tree.NewNode(an.NewToken(an.TokenUM, "-", n.Token.Pos))
		um.Left = n
		return um
	}
	one := tree.NewNode(an.NewToken(an.TokenCONST, "1", n.Token.Pos))
	return tree.NewOp("/", one, n)
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
