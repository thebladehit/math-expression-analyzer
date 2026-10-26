package parallel

import (
	"cmp"
	"slices"

	"pzcs-lab1/group"
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
		return group.Negate(items[0].Node, g)
	}
	return items[0].Node
}

func join(a, b group.Item, g group.Group) group.Item {
	switch {
	case a.Inverse == b.Inverse:
		return group.Item{Node: tree.NewOp(g.Direct, a.Node, b.Node), Inverse: a.Inverse}
	case b.Inverse:
		return group.Item{Node: tree.NewOp(g.Inverse, a.Node, b.Node)}
	default:
		return group.Item{Node: tree.NewOp(g.Inverse, b.Node, a.Node)}
	}
}
