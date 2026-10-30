package dataflow

import (
	"fmt"
	"testing"

	"pzcs-lab1/an"
	"pzcs-lab1/postfix"
	"pzcs-lab1/tree"
)

func build(t *testing.T, expr string) *tree.Node {
	t.Helper()
	tokens, errs := an.NewAnalyzer(an.NewLexer(), an.NewParser(an.StState)).Start(expr)
	if errs != nil {
		t.Fatalf("invalid expr %q: %v", expr, errs)
	}
	return tree.NewTreeBuilder().Build(postfix.NewPostfix().Convert(tokens))
}

func TestBuildGraph(t *testing.T) {
	tests := []struct {
		expr string
		want []string
	}{
		{"a", nil},
		{"a*2", []string{"op1: *    (a, 2)"}},
		{"a*b+c/d", []string{
			"op1: *    (a, b)",
			"op2: /    (c, d)",
			"op3: +    (op1, op2)",
		}},
		{"b/(b-a)", []string{
			"op1: -    (b, a)",
			"op2: /    (b, op1)",
		}},
		{"-sin(x)+1", []string{
			"op1: sin  (x)",
			"op2: neg  (op1)",
			"op3: +    (op2, 1)",
		}},
	}
	for _, tt := range tests {
		got := BuildGraph(build(t, tt.expr))
		if len(got) != len(tt.want) {
			t.Fatalf("%s: got %d tasks %v, want %d", tt.expr, len(got), got, len(tt.want))
		}
		for i := range got {
			if s := fmt.Sprint(got[i]); s != tt.want[i] {
				t.Errorf("%s: task %d = %q, want %q", tt.expr, i, s, tt.want[i])
			}
		}
	}
}
