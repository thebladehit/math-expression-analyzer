package parallel

import (
	"cmp"
	"slices"

	"pzcs-lab1/an"
	"pzcs-lab1/tree"
)

type Parallelizer interface {
	Parallelize(n *tree.Node) *tree.Node
}

type parallelizer struct{}

func NewParallelizer() Parallelizer {
	return &parallelizer{}
}

type item struct {
	node    *tree.Node
	inverse bool
}

type group struct {
	direct  string
	inverse string
}

var (
	additive       = group{"+", "-"}
	multiplicative = group{"*", "/"}
)

func (p *parallelizer) Parallelize(n *tree.Node) *tree.Node {
	switch n.Type {
	case tree.NodeOP:
		g := groupOf(n)
		items := collect(n, g, false)
		for i := range items {
			items[i].node = p.Parallelize(items[i].node)
		}
		return combine(items, g)
	case tree.NodeFN, tree.NodeUM:
		n.Left = p.Parallelize(n.Left)
	}
	return n
}

func groupOf(n *tree.Node) group {
	if n.Token.Text == "*" || n.Token.Text == "/" {
		return multiplicative
	}
	return additive
}

func collect(n *tree.Node, g group, inverse bool) []item {
	if n.Type == tree.NodeOP && (n.Token.Text == g.direct || n.Token.Text == g.inverse) {
		rightInv := inverse
		if n.Token.Text == g.inverse {
			rightInv = !rightInv
		}
		return append(collect(n.Left, g, inverse), collect(n.Right, g, rightInv)...)
	}
	if n.Type == tree.NodeUM && g == additive {
		return collect(n.Left, g, !inverse)
	}
	return []item{{node: n, inverse: inverse}}
}

func combine(items []item, g group) *tree.Node {
	for len(items) > 1 {
		slices.SortStableFunc(items, func(a, b item) int {
			return cmp.Compare(a.node.Height(), b.node.Height())
		})
		items = append(items[2:], join(items[0], items[1], g))
	}
	if items[0].inverse {
		return negate(items[0].node, g)
	}
	return items[0].node
}

func join(a, b item, g group) item {
	switch {
	case a.inverse == b.inverse:
		return item{newOp(g.direct, a.node, b.node), a.inverse}
	case b.inverse:
		return item{newOp(g.inverse, a.node, b.node), false}
	default:
		return item{newOp(g.inverse, b.node, a.node), false}
	}
}

func negate(n *tree.Node, g group) *tree.Node {
	if g == additive {
		um := tree.NewNode(an.NewToken(an.TokenUM, "-", n.Token.Pos))
		um.Left = n
		return um
	}
	one := tree.NewNode(an.NewToken(an.TokenCONST, "1", n.Token.Pos))
	return newOp("/", one, n)
}

func newOp(op string, left, right *tree.Node) *tree.Node {
	node := tree.NewNode(an.NewToken(an.TokenOP, op, left.Token.Pos))
	node.Left = left
	node.Right = right
	return node
}
