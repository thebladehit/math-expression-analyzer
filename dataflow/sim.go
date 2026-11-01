package dataflow

import (
	"fmt"
	"slices"
)

type Config struct {
	Processors int
	Banks      int
	Times      map[Op]int
}

func (c Config) Duration(op Op) int {
	if d, ok := c.Times[op]; ok {
		return d
	}
	return c.Times[OpFn]
}

type Slot struct {
	Unit  string
	Start int
	End   int
	Label string
}

type Result struct {
	Slots []Slot
	Time  int
}

func Simulate(tasks []Task, bankOf map[string]int, cfg Config) Result {
	n := len(tasks)
	if n == 0 {
		return Result{}
	}

	prio := priorities(tasks, cfg)
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	slices.SortStableFunc(order, func(a, b int) int { return prio[b] - prio[a] })

	end := make([]int, n)
	procFree := make([]int, cfg.Processors)
	mem := newBanks(bankOf, cfg)
	var res Result

	for _, i := range order {
		ready := 0
		for _, d := range tasks[i].Deps {
			ready = max(ready, end[d-1])
		}

		for _, v := range unique(tasks[i].Reads) {
			s := mem.read(v)
			res.Slots = append(res.Slots, s)
			ready = max(ready, s.End)
		}

		p := slices.Index(procFree, slices.Min(procFree))
		start := max(ready, procFree[p])
		end[i] = start + cfg.Duration(tasks[i].Op)
		procFree[p] = end[i]
		res.Slots = append(res.Slots, Slot{fmt.Sprintf("P%d", p+1), start, end[i], fmt.Sprintf("op%d", tasks[i].ID)})
	}

	s := mem.write(end[n-1])
	res.Slots = append(res.Slots, s)
	res.Time = s.End
	return res
}

func priorities(tasks []Task, cfg Config) []int {
	prio := make([]int, len(tasks))
	for i := len(tasks) - 1; i >= 0; i-- {
		prio[i] += cfg.Duration(tasks[i].Op)
		for _, d := range tasks[i].Deps {
			prio[d-1] = prio[i]
		}
	}
	return prio
}
