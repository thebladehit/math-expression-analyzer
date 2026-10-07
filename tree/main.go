package tree

import "pzcs-lab1/an"

type TreeBuilder interface {
	Build(tokens []an.Token) *Node
}

type treeBuilder struct{}

func NewTreeBuilder() TreeBuilder {
	return &treeBuilder{}
}

func (b *treeBuilder) Build(tokens []an.Token) *Node {
	stack := make([]*Node, 0, len(tokens))

	for _, t := range tokens {
		switch t.Type {
		case an.TokenVAR, an.TokenCONST:
			stack = append(stack, NewNode(t))
		case an.TokenOP:
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			node := NewNode(t)
			node.Left = left
			node.Right = right

			stack = append(stack, node)
		case an.TokenFN, an.TokenUM:
			left := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			node := NewNode(t)
			node.Left = left

			stack = append(stack, node)
		}
	}

	return stack[0]
}
