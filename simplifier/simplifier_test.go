package simplifier

import (
	"testing"

	"pzcs-lab1/an"
	"pzcs-lab1/postfix"
	"pzcs-lab1/tree"
)

func build(t *testing.T, expr string) *tree.Node {
	t.Helper()
	tokens, errs := an.NewLexer().Tokenize(expr)
	if len(errs) > 0 {
		t.Fatalf("%q: лексичні помилки %v", expr, errs)
	}
	return tree.NewTreeBuilder().Build(postfix.NewPostfix().Convert(tokens))
}

func TestSimplify(t *testing.T) {
	cases := []struct {
		in, want string
		opts     int
	}{
		{"a*1", "a", 1},
		{"1*a", "a", 1},
		{"a/1", "a", 1},
		{"1/a", "(1/a)", 0},
		{"a+0", "a", 1},
		{"0+a", "a", 1},
		{"a-0", "a", 1},
		{"0-a", "(0-a)", 0},
		{"a*0", "0", 1},
		{"0*a", "0", 1},
		{"0/a", "0", 1},
		{"(a*1)*1", "a", 2},
		{"(a+0)*1", "a", 2},
		{"x*0+1", "1", 2},
		{"2+3", "5", 1},
		{"2+3*a", "(2+(3*a))", 0},
		{"(2+3)*a", "(5*a)", 1},
		{"2*0", "0", 1},
		{"2.5+1", "3.5", 1},
		{"7/2", "3.5", 1},
		{"1/0", "(1/0)", 0},
		{"-2+5", "3", 2},
		{"-2*a", "((-2)*a)", 1},
		{"sin(a*1)", "sin(a)", 1},
		{"a+b", "(a+b)", 0},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, opts := NewSimplifier().Simplify(build(t, tc.in))
			if got.String() != tc.want {
				t.Errorf("отримано %s, очікувалось %s", got, tc.want)
			}
			if len(opts) != tc.opts {
				t.Errorf("оптимізацій %d, очікувалось %d: %v", len(opts), tc.opts, opts)
			}
		})
	}
}

func TestSimplifyLogReset(t *testing.T) {
	s := NewSimplifier()
	s.Simplify(build(t, "a*1"))
	if _, opts := s.Simplify(build(t, "a+b")); len(opts) != 0 {
		t.Errorf("журнал не очищено: %v", opts)
	}
}
