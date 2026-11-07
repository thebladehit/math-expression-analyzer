package dataflow

type Metrics struct {
	T1       int
	Tp       int
	Speedup  float64
	Used     int
	Total    int
	EffUsed  float64
	EffTotal float64
}

func Measure(res Result, cfg Config) Metrics {
	m := Metrics{Tp: res.Time, Total: cfg.Processors}
	if res.Time == 0 {
		return m
	}

	used := map[Unit]bool{}
	for _, s := range res.Slots {
		m.T1 += s.End - s.Start
		if s.Unit.Kind == Proc {
			used[s.Unit] = true
		}
	}
	m.Used = len(used)

	m.Speedup = float64(m.T1) / float64(m.Tp)
	m.EffUsed = m.Speedup / float64(m.Used)
	m.EffTotal = m.Speedup / float64(m.Total)
	return m
}
