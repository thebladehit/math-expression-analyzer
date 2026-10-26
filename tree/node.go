package tree

import (
	"fmt"
	"pzcs-lab1/an"
)

type NodeType string

const (
	NodeVAR   NodeType = "VAR"
	NodeOP    NodeType = "OP"
	NodeCONST NodeType = "CONST"
	NodeUM    NodeType = "UM"
	NodeFN    NodeType = "FN"
)

type Node struct {
	Type  NodeType
	Token an.Token
	Left  *Node
	Right *Node
}

func NewNode(token an.Token) *Node {
	return &Node{
		Type:  nodeType(token),
		Token: token,
	}
}

func NewOp(op string, left, right *Node) *Node {
	node := NewNode(an.NewToken(an.TokenOP, op, left.Token.Pos))
	node.Left = left
	node.Right = right
	return node
}

func (n *Node) Height() int {
	if n == nil {
		return -1
	}

	lH := n.Left.Height() + 1
	rH := n.Right.Height() + 1

	return max(lH, rH)
}

func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}
	return &Node{
		Type:  n.Type,
		Token: n.Token,
		Left:  n.Left.Clone(),
		Right: n.Right.Clone(),
	}
}

func nodeType(token an.Token) NodeType {
	switch token.Type {
	case an.TokenVAR:
		return NodeVAR
	case an.TokenCONST:
		return NodeCONST
	case an.TokenOP:
		return NodeOP
	case an.TokenUM:
		return NodeUM
	case an.TokenFN:
		return NodeFN
	default:
		panic(fmt.Sprintf("Undefined token type in nodeType: %v", token.Type))
	}
}
