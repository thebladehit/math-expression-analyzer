package commute

import (
	"iter"

	"pzcs-lab1/group"
	"pzcs-lab1/tree"
)

type Commuter interface {
	Forms(n *tree.Node, limit int) []*tree.Node
}

type commuter struct{}

func NewCommuter() Commuter {
	return &commuter{}
}

func (c *commuter) Forms(n *tree.Node, limit int) []*tree.Node {
	return unique(commute(n, limit))
}

func unique(forms []*tree.Node) []*tree.Node {
	seen := map[string]bool{}
	var result []*tree.Node
	for _, f := range forms {
		if key := tree.ToInfix(f); !seen[key] {
			seen[key] = true
			result = append(result, f)
		}
	}
	return result
}

func commute(n *tree.Node, limit int) []*tree.Node {
	var result []*tree.Node

	switch n.Type {
	case tree.NodeFN, tree.NodeUM:
		for _, arg := range commute(n.Left, limit) {
			result = append(result, &tree.Node{Type: n.Type, Token: n.Token, Left: arg})
		}

	case tree.NodeOP:
		g := group.Of(n)
		items := group.Collect(n, g, false)

		variants := make([][]group.Item, len(items))
		for i, it := range items {
			for _, f := range commute(it.Node, limit) {
				variants[i] = append(variants[i], group.Item{Node: f, Inverse: it.Inverse})
			}
		}

		for chosen := range arrange(variants) {
			result = append(result, combine(chosen, g))
			if len(result) == limit {
				return result
			}
		}

	default:
		result = append(result, n.Clone())
	}
	return result
}

func combine(items []group.Item, g group.Group) *tree.Node {
	acc := items[0].Node.Clone()
	if items[0].Inverse {
		acc = group.Negate(acc, g)
	}
	for _, it := range items[1:] {
		op := g.Direct
		if it.Inverse {
			op = g.Inverse
		}
		acc = tree.NewOp(op, acc, it.Node.Clone())
	}
	return acc
}

func arrange(operands [][]group.Item) iter.Seq[[]group.Item] {
	return func(yield func([]group.Item) bool) {
		var walk func(chosen []group.Item, rest [][]group.Item) bool
		walk = func(chosen []group.Item, rest [][]group.Item) bool {
			if len(rest) == 0 {
				return yield(chosen)
			}
			for i, forms := range rest {
				others := append(append([][]group.Item{}, rest[:i]...), rest[i+1:]...)
				for _, f := range forms {
					next := append(append([]group.Item{}, chosen...), f)
					if !walk(next, others) {
						return false
					}
				}
			}
			return true
		}
		walk(nil, operands)
	}
}
