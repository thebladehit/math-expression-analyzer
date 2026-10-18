package simplifier

import (
	"testing"

	"pzcs-lab1/an"
	"pzcs-lab1/postfix"
	"pzcs-lab1/tree"
)

func build(expr string) *tree.Node {
	tokens, _ := an.NewLexer().Tokenize(expr)
	return tree.NewTreeBuilder().Build(postfix.NewPostfix().Convert(tokens))
}

func TestSimplify(t *testing.T) {
	cases := map[string]int{
		"a+b":     0,
		"a*1":     1,
		"a/1":     1,
		"a+0":     1,
		"a-0":     1,
		"a*0":     1,
		"0/a":     1,
		"2+3":     1,
		"1/a":     0,
		"0-a":     0,
		"1/0":     0,
		"(a*1)*1": 2,
		"x*0+1":   2,
	}
	for expr, want := range cases {
		_, opts := NewSimplifier().Simplify(build(expr))
		if len(opts) != want {
			t.Errorf("%s: оптимізацій %d, очікувалось %d", expr, len(opts), want)
		}
	}
}
