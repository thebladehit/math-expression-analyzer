package tree

import (
	"fmt"
	"strings"
)

func (n *Node) Draw() string {
	if n == nil {
		return ""
	}

	const gap = 3
	pos := map[*Node]int{}
	next := 0
	place(n, pos, &next, gap)

	top := n.Height()
	rows := 2*top + 1
	width := next + 1
	grid := make([][]rune, rows)
	for i := range grid {
		grid[i] = []rune(strings.Repeat(" ", width))
	}

	var drawNode func(*Node)
	drawNode = func(m *Node) {
		if m == nil {
			return
		}
		row := 2 * (top - m.Height())
		writeLabel(grid[row], pos[m], label(m))
		for _, child := range []*Node{m.Left, m.Right} {
			if child != nil {
				drawEdge(grid, pos[m], row, label(m), pos[child], 2*(top-child.Height()))
				drawNode(child)
			}
		}
	}
	drawNode(n)

	var b strings.Builder
	for i, line := range grid {
		prefix := "    "
		if i%2 == 0 {
			prefix = fmt.Sprintf("%-4d", top-i/2)
		}
		b.WriteString(strings.TrimRight(prefix+string(line), " "))
		b.WriteByte('\n')
	}
	return b.String()
}

func place(n *Node, pos map[*Node]int, next *int, gap int) {
	switch {
	case n.Left == nil && n.Right == nil:
		w := len([]rune(label(n)))
		pos[n] = *next + w/2
		*next += w + gap
	case n.Right == nil:
		place(n.Left, pos, next, gap)
		pos[n] = pos[n.Left]
	default:
		place(n.Left, pos, next, gap)
		place(n.Right, pos, next, gap)
		pos[n] = (pos[n.Left] + pos[n.Right]) / 2
	}
}

func label(n *Node) string {
	switch {
	case n.Type == NodeUM:
		return "-u"
	case n.Type == NodeOP && n.Token.Text == "/":
		return "÷"
	}
	return n.Token.Text
}

func labelStart(x int, text string) int {
	return max(0, x-len([]rune(text))/2)
}

func writeLabel(line []rune, x int, text string) {
	start := labelStart(x, text)
	for i, r := range []rune(text) {
		if start+i < len(line) {
			line[start+i] = r
		}
	}
}

func drawEdge(grid [][]rune, px, prow int, text string, cx, crow int) {
	set := func(r, x int, ch rune) {
		if r >= 0 && r < len(grid) && x >= 0 && x < len(grid[r]) && grid[r][x] == ' ' {
			grid[r][x] = ch
		}
	}

	rows := crow - prow - 1
	dx := cx - px
	if dx == 0 {
		for r := prow + 1; r < crow; r++ {
			set(r, cx, '|')
		}
		return
	}

	dir, slash := 1, '\\'
	if dx < 0 {
		dir, slash = -1, '/'
	}
	dist := dx * dir

	if dist > rows {
		top := cx - dir*rows
		edge := labelStart(px, text) - 1
		if dir > 0 {
			edge = labelStart(px, text) + len([]rune(text))
		}
		for x := edge; x*dir < top*dir; x += dir {
			set(prow, x, '_')
		}
		for i := 1; i <= rows; i++ {
			set(prow+i, cx-dir*(rows-i+1), slash)
		}
		return
	}

	for i := 1; i < dist; i++ {
		set(prow+i, px+dir*i, slash)
	}
	for r := prow + dist; r < crow; r++ {
		set(r, cx, '|')
	}
}
