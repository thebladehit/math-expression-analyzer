package dataflow

import "testing"

func TestDistribute(t *testing.T) {
	tests := []struct {
		expr  string
		banks int
	}{
		{"a*b+c*d", 2},
		{"a*2 + b/(b-a) - 1/x4/x3", 2},
		{"a+b+c+d+e+f+g+h", 2},
		{"a*b*c*d", 4},
	}
	for _, tt := range tests {
		tasks := BuildGraph(build(t, tt.expr))
		bankOf := Distribute(tasks, tt.banks)

		for _, task := range tasks {
			for _, v := range task.Reads {
				if _, ok := bankOf[v]; !ok {
					t.Errorf("%s: змінна %s не має банку", tt.expr, v)
				}
			}
			if len(task.Reads) == 2 && task.Reads[0] != task.Reads[1] &&
				bankOf[task.Reads[0]] == bankOf[task.Reads[1]] {
				t.Errorf("%s: op%d читає %v з одного банку", tt.expr, task.ID, task.Reads)
			}
		}
	}
}
