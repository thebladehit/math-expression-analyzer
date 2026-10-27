package associate

import (
	"iter"

	"pzcs-lab1/group"
	"pzcs-lab1/tree"
)

type Associator interface {
	Forms(n *tree.Node, limit int) []*tree.Node
}

type associator struct{}

func NewAssociator() Associator {
	return &associator{}
}

func (a *associator) Forms(n *tree.Node, limit int) []*tree.Node {
	return tree.Unique(associate(n, limit))
}

func associate(n *tree.Node, limit int) []*tree.Node {
	var result []*tree.Node

	switch n.Type {
	case tree.NodeFN, tree.NodeUM:
		for _, arg := range associate(n.Left, limit) {
			result = append(result, &tree.Node{Type: n.Type, Token: n.Token, Left: arg})
		}

	case tree.NodeOP:
		g := group.Of(n)
		items := group.Collect(n, g, false)

		variants := make([][]group.Item, len(items))
		for i, it := range items {
			for _, f := range associate(it.Node, limit) {
				if i == 0 && it.Inverse {
					variants[i] = append(variants[i], group.Item{Node: group.Negate(f, g)})
				} else {
					variants[i] = append(variants[i], group.Item{Node: f, Inverse: it.Inverse})
				}
			}
		}

		for t := range bracketing(variants, g) {
			result = append(result, t.Node.Clone())
			if len(result) == limit {
				return result
			}
		}

	default:
		result = append(result, n.Clone())
	}
	return result
}

func bracketing(operands [][]group.Item, g group.Group) iter.Seq[group.Item] {
	var build func(lo, hi int, yield func(group.Item) bool) bool

	build = func(lo, hi int, yield func(group.Item) bool) bool {
		if hi-lo == 1 {
			for _, f := range operands[lo] {
				if !yield(f) {
					return false
				}
			}
			return true
		}
		for cut := lo + 1; cut < hi; cut++ {
			ok := build(lo, cut, func(left group.Item) bool {
				return build(cut, hi, func(right group.Item) bool {
					return yield(join(left, right, g))
				})
			})
			if !ok {
				return false
			}
		}
		return true
	}

	return func(yield func(group.Item) bool) {
		build(0, len(operands), yield)
	}
}

func join(left, right group.Item, g group.Group) group.Item {
	op := g.Direct
	if left.Inverse != right.Inverse {
		op = g.Inverse
	}
	return group.Item{Node: tree.NewOp(op, left.Node, right.Node), Inverse: left.Inverse}
}
