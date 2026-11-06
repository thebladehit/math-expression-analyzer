package dataflow

import (
	"math"
	"testing"
)

func TestMeasure(t *testing.T) {
	m := Measure(simulate(t, "a*b+c*d", testCfg), testCfg)

	if m.T1 != 10 || m.Tp != 6 || m.Used != 3 || m.Total != 4 {
		t.Fatalf("got %+v", m)
	}
	for name, got := range map[string]float64{
		"Speedup":  m.Speedup - 10.0/6,
		"EffUsed":  m.EffUsed - 10.0/6/3,
		"EffTotal": m.EffTotal - 10.0/6/4,
	} {
		if math.Abs(got) > 1e-9 {
			t.Errorf("%s: off by %v", name, got)
		}
	}
}

func TestMeasureEmpty(t *testing.T) {
	m := Measure(simulate(t, "a", testCfg), testCfg)
	if m.T1 != 0 || m.Tp != 0 || m.Speedup != 0 {
		t.Errorf("got %+v", m)
	}
}
