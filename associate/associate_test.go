package associate

import (
	"slices"
	"testing"

	"pzcs-lab1/an"
	"pzcs-lab1/postfix"
	"pzcs-lab1/tree"
)

func build(expr string) *tree.Node {
	tokens, _ := an.NewLexer().Tokenize(expr)
	return tree.NewTreeBuilder().Build(postfix.NewPostfix().Convert(tokens))
}

func forms(expr string) []string {
	var result []string
	for _, f := range NewAssociator().Forms(build(expr), 1000) {
		result = append(result, tree.ToInfix(f))
	}
	return result
}

func TestFormsCount(t *testing.T) {
	cases := map[string]int{
		"a":             1,
		"a+b":           1,
		"a+b+c":         2,
		"a+b+c+d":       5,
		"a*b*c*d":       5,
		"a+b+c+d+e":     14,
		"a-b-c":         2,
		"a-b-c-d":       5,
		"a+b*c*d":       2,
		"(a+b+c)*(d+e)": 2,
		"sin(a+b+c)":    2,
		"-a+b":          1,
	}
	for expr, want := range cases {
		if got := forms(expr); len(got) != want {
			t.Errorf("%s: форм %d, очікувалось %d: %v", expr, len(got), want, got)
		}
	}
}

func TestFormsSigns(t *testing.T) {
	cases := map[string][]string{
		"a+b+c": {"a+(b+c)", "a+b+c"},
		"a-b-c": {"a-(b+c)", "a-b-c"},
		"a-b+c": {"a-(b-c)", "a-b+c"},
		"a/b/c": {"a/(b*c)", "a/b/c"},
		"-a-b":  {"-a-b"},
	}
	for expr, want := range cases {
		if got := forms(expr); !slices.Equal(got, want) {
			t.Errorf("%s: отримано %v, очікувалось %v", expr, got, want)
		}
	}
}
