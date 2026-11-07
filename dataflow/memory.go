package dataflow

import "slices"

func Distribute(tasks []Task, banks int) map[string]int {
	bankOf := map[string]int{}
	for _, t := range tasks {
		for _, v := range t.Reads {
			if _, ok := bankOf[v]; !ok {
				bankOf[v] = len(bankOf) % banks
			}
		}
	}
	return bankOf
}

type banks struct {
	bankOf map[string]int
	free   []int
	time   int
}

func newBanks(bankOf map[string]int, cfg Config) *banks {
	return &banks{bankOf: bankOf, free: make([]int, cfg.Banks), time: cfg.Duration(OpMem)}
}

func (m *banks) read(v string) Slot {
	b := m.bankOf[v]
	return m.use(b, m.free[b], v)
}

func (m *banks) write(at int) Slot {
	b := slices.Index(m.free, slices.Min(m.free))
	return m.use(b, max(at, m.free[b]), "w")
}

func (m *banks) use(b, start int, label string) Slot {
	m.free[b] = start + m.time
	return Slot{Unit{Bank, b}, start, m.free[b], label}
}

func unique(vs []string) []string {
	var out []string
	for _, v := range vs {
		if !slices.Contains(out, v) {
			out = append(out, v)
		}
	}
	return out
}
