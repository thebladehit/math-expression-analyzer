package dataflow

import (
	"fmt"
	"strings"

	"pzcs-lab1/tree"
)

type Task struct {
	ID    int
	Op    Op
	Deps  []int
	Reads []string
	Imm   []string
	Args  []string
}

func (t Task) String() string {
	return fmt.Sprintf("op%d: %-4s (%s)", t.ID, t.Op, strings.Join(t.Args, ", "))
}

func BuildGraph(root *tree.Node) []Task {
	var tasks []Task

	var visit func(n *tree.Node) int
	visit = func(n *tree.Node) int {
		if n.Type == tree.NodeVAR || n.Type == tree.NodeCONST {
			return -1
		}

		t := Task{Op: opName(n)}
		for _, child := range []*tree.Node{n.Left, n.Right} {
			if child == nil {
				continue
			}
			switch child.Type {
			case tree.NodeVAR:
				t.Reads = append(t.Reads, child.Token.Text)
				t.Args = append(t.Args, child.Token.Text)
			case tree.NodeCONST:
				t.Imm = append(t.Imm, child.Token.Text)
				//t.Reads = append(t.Reads, child.Token.Text) // Load to Read
				t.Args = append(t.Args, child.Token.Text)
			default:
				dep := visit(child)
				t.Deps = append(t.Deps, dep)
				t.Args = append(t.Args, fmt.Sprintf("op%d", dep))
			}
		}

		t.ID = len(tasks) + 1
		tasks = append(tasks, t)
		return t.ID
	}

	visit(root)
	return tasks
}
