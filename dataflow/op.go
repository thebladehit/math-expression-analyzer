package dataflow

import "pzcs-lab1/tree"

type Op string

const (
	OpAdd Op = "+"
	OpSub Op = "-"
	OpMul Op = "*"
	OpDiv Op = "/"
	OpNeg Op = "neg"
	OpFn  Op = "fn"
	OpMem Op = "mem"
)

func opName(n *tree.Node) Op {
	if n.Type == tree.NodeUM {
		return OpNeg
	}
	return Op(n.Token.Text)
}
