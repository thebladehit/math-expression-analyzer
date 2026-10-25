package parallel

import (
	"cmp"
	"pzcs-lab1/group"
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

func (p *parallelizer) Parallelize(n *tree.Node) *tree.Node {
	switch n.Type {
	case tree.NodeOP:
		g := group.Of(n)
		items := group.Collect(n, g, false)
		for i := range items {
			items[i].Node = p.Parallelize(items[i].Node)
		}
		return combine(items, g)
	case tree.NodeFN, tree.NodeUM:
		n.Left = p.Parallelize(n.Left)
	}
	return n
}

func combine(items []group.Item, g group.Group) *tree.Node {
	for len(items) > 1 {
		slices.SortStableFunc(items, func(a, b group.Item) int {
			return cmp.Compare(a.Node.Height(), b.Node.Height())
		})
		items = append(items[2:], join(items[0], items[1], g))
	}
	if items[0].Inverse {
		return negate(items[0].Node, g)
	}
	return items[0].Node
}

func join(a, b group.Item, g group.Group) group.Item {
	switch {
	case a.Inverse == b.Inverse:
		return group.Item{Node: newOp(g.Direct, a.Node, b.Node), Inverse: a.Inverse}
	case b.Inverse:
		return group.Item{Node: newOp(g.Inverse, a.Node, b.Node)}
	default:
		return group.Item{Node: newOp(g.Inverse, b.Node, a.Node)}
	}
}

func negate(n *tree.Node, g group.Group) *tree.Node {
	if g == group.Additive {
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
