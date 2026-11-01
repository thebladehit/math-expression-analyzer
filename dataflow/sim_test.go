package dataflow

import "testing"

var testCfg = Config{
	Processors: 4,
	Banks:      2,
	Times:      map[Op]int{OpAdd: 1, OpSub: 1, OpMul: 2, OpDiv: 4, OpNeg: 1, OpFn: 4, OpMem: 1},
}

func simulate(t *testing.T, expr string, cfg Config) Result {
	t.Helper()
	tasks := BuildGraph(build(t, expr))
	return Simulate(tasks, Distribute(tasks, cfg.Banks), cfg)
}

func TestSimulateTime(t *testing.T) {
	one := testCfg
	one.Processors, one.Banks = 1, 1

	tests := []struct {
		expr string
		cfg  Config
		want int
	}{
		{"a", testCfg, 0},
		{"a*b+c*d", testCfg, 6}, // читання a,b | c,d; * на P1 і P2; +; запис
		{"a*a", testCfg, 4},     // a читається один раз
		{"a*2", testCfg, 4},     // константа не читається
		{"a*b+c*d", one, 8},     // читання c,d йде паралельно з a*b
	}
	for _, tt := range tests {
		if got := simulate(t, tt.expr, tt.cfg).Time; got != tt.want {
			t.Errorf("%s (P=%d, B=%d): Tp = %d, want %d", tt.expr, tt.cfg.Processors, tt.cfg.Banks, got, tt.want)
		}
	}
}

func TestSimulateNoOverlap(t *testing.T) {
	exprs := []string{
		"a*2 + b/(b-a) - 1/x4/x3",
		"a+b+c+d+e+f+g+h",
		"a/b/c/d/e/f/g/h",
		"a+b+c+d+e+f+g+h+i+j+k+l+m*n+o+p+q+a+a+t+u+v+w+x+y+z",
	}
	for _, expr := range exprs {
		res := simulate(t, expr, testCfg)
		for i, a := range res.Slots {
			for _, b := range res.Slots[i+1:] {
				if a.Unit == b.Unit && a.Start < b.End && b.Start < a.End {
					t.Errorf("%s: %s зайнятий двічі: %v і %v", expr, a.Unit, a, b)
				}
			}
		}
	}
}
