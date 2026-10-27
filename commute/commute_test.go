package commute

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

func TestForms(t *testing.T) {
	cases := map[string]int{
		"a":           1,
		"a+b":         2,
		"a-b":         2,
		"a+b+c":       6,
		"a*b+c":       4,
		"a-b*c":       4,
		"(a+b)*(c+d)": 8,
		"sin(a+b)":    2,
		"x+x":         1,
	}
	for expr, want := range cases {
		forms := NewCommuter().Forms(build(expr), 1000)
		if len(forms) != want {
			var got []string
			for _, f := range forms {
				got = append(got, tree.ToInfix(f))
			}
			t.Errorf("%s: форм %d, очікувалось %d: %v", expr, len(forms), want, got)
		}
	}
}
