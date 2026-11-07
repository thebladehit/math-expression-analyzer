package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	an "pzcs-lab1/an"
	df "pzcs-lab1/dataflow"
	"pzcs-lab1/tree"
)

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
	colorBold  = "\033[1m"
	colorDim   = "\033[2m"
)

var useColor = func() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	fi, err := os.Stdout.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}()

func paint(s, color string) string {
	if !useColor || s == "" {
		return s
	}
	return color + s + colorReset
}

func PrintOptimizations(optimizations []string) {
	if len(optimizations) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(paint("Оптимізації:", colorBold))
	for i, o := range optimizations {
		fmt.Printf("  %d. %s\n", i+1, o)
	}
}

func PrintForms(title string, forms []*tree.Node, limit int) {
	fmt.Println()
	fmt.Println(paint(fmt.Sprintf("%s (%d):", title, len(forms)), colorBold))
	for i, f := range forms {
		fmt.Printf("  %2d. %s\n", i+1, tree.ToInfix(f))
	}
	if len(forms) == limit {
		fmt.Println(paint(fmt.Sprintf("  … показано перші %d", limit), colorDim))
	}
}

func PrintTree(title string, root *tree.Node) {
	fmt.Println()
	fmt.Println(paint(title+":", colorBold))
	fmt.Print(root.Draw())
	fmt.Println(paint(fmt.Sprintf("ярусів: %d", root.Height()), colorDim))
}

func PrintTokens(tokens []an.Token) {
	fmt.Println()
	fmt.Println(paint("Токени:", colorBold))
	fmt.Println(paint("   №  поз.  тип      текст", colorDim))
	for i, t := range tokens {
		text := t.Text
		if t.Type == an.TokenEND {
			text = paint("(кінець виразу)", colorDim)
		}
		fmt.Printf("  %2d  %4d  %s %s\n", i+1, t.Pos+1, paint(fmt.Sprintf("%-8s", t.Type), colorCyan), text)
	}
}

func PrintResult(input string, tokens []an.Token, errs []an.Error) {
	PrintTokens(tokens)
	fmt.Println()
	if len(errs) == 0 {
		fmt.Println("  " + input)
		fmt.Println(paint("✓ Вираз коректний", colorGreen+colorBold))
		return
	}

	errs = slices.Clone(errs)
	slices.SortStableFunc(errs, func(a, b an.Error) int { return a.Pos - b.Pos })

	runes := []rune(input)
	marks := make([]rune, len(runes)+1)
	for _, e := range errs {
		n := tokenLen(tokens, e.Pos)
		for i := e.Pos; i < e.Pos+n && i < len(marks); i++ {
			if i < 0 {
				continue
			}
			if i == e.Pos {
				marks[i] = '^'
			} else if marks[i] == 0 {
				marks[i] = '~'
			}
		}
	}

	var expr, under strings.Builder
	for i, m := range marks {
		if i == len(runes) && m == 0 {
			break
		}
		ch := " "
		if i < len(runes) {
			ch = string(runes[i])
		}
		if m == 0 {
			expr.WriteString(ch)
			under.WriteRune(' ')
			continue
		}
		expr.WriteString(paint(ch, colorRed+colorBold))
		under.WriteRune(m)
	}

	fmt.Println("  " + expr.String())
	fmt.Println("  " + paint(strings.TrimRight(under.String(), " "), colorRed))
	fmt.Println(paint(fmt.Sprintf("✗ Знайдено помилок: %d", len(errs)), colorRed+colorBold))
	for i, e := range errs {
		pos := paint(fmt.Sprintf("[поз. %d]", e.Pos+1), colorDim)
		fmt.Printf("  %d. %s %s\n", i+1, pos, e.Msg)
	}
}

func tokenLen(tokens []an.Token, pos int) int {
	for _, t := range tokens {
		if t.Pos == pos {
			return max(1, len([]rune(t.Text)))
		}
	}
	return 1
}

func TokensToString(tokens []an.Token) string {
	parts := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if t.Type == an.TokenUM {
			parts = append(parts, "-u")
			continue
		}
		parts = append(parts, t.Text)
	}
	return strings.Join(parts, " ")
}

func PrintGraph(tasks []df.Task) {
	fmt.Println("\nГраф задачі:")
	if len(tasks) == 0 {
		fmt.Println("  вираз не містить операцій")
	}
	for _, t := range tasks {
		fmt.Println(" ", t)
	}
}

func PrintBanks(bankOf map[string]int, banks int) {
	vars := make([][]string, banks)
	for v, b := range bankOf {
		vars[b] = append(vars[b], v)
	}
	fmt.Println("\nРозкладка змінних по банках:")
	for b, vs := range vars {
		slices.Sort(vs)
		fmt.Printf("  B%d: %s\n", b+1, strings.Join(vs, ", "))
	}
}

func PrintGantt(res df.Result, processors, banks int) {
	fmt.Println("\nДіаграма Ганта:")
	if res.Time == 0 {
		fmt.Println("  вираз не містить операцій")
		return
	}

	cell := 3
	for _, s := range res.Slots {
		cell = max(cell, len(s.Label)+2)
	}

	var units []df.Unit
	for p := range processors {
		units = append(units, df.Unit{Kind: df.Proc, N: p})
	}
	for b := range banks {
		units = append(units, df.Unit{Kind: df.Bank, N: b})
	}

	fmt.Print("    ")
	for t := range res.Time {
		fmt.Printf("%-*d", cell, t)
	}
	fmt.Println()

	for _, u := range units {
		var slots []df.Slot
		for _, s := range res.Slots {
			if s.Unit == u {
				slots = append(slots, s)
			}
		}
		slices.SortFunc(slots, func(a, b df.Slot) int { return a.Start - b.Start })

		color := colorCyan
		if u.Kind == df.Bank {
			color = colorGreen
		}

		var line strings.Builder
		t := 0
		for _, s := range slots {
			for ; t < s.Start; t++ {
				line.WriteString(paint("·"+strings.Repeat(" ", cell-1), colorDim))
			}
			width := (s.End-s.Start)*cell - 2
			line.WriteString(paint("["+s.Label+strings.Repeat(" ", width-len(s.Label))+"]", color))
			t = s.End
		}
		for ; t < res.Time; t++ {
			line.WriteString(paint("·"+strings.Repeat(" ", cell-1), colorDim))
		}
		fmt.Printf("%-4s%s\n", u, line.String())
	}
	fmt.Println("Tp =", res.Time)
}

func PrintMetrics(m df.Metrics) {
	if m.Tp == 0 {
		return
	}
	fmt.Println("\nМетрики:")
	fmt.Printf("  Час послідовного виконання T1:          %d\n", m.T1)
	fmt.Printf("  Час паралельного виконання Tp:          %d\n", m.Tp)
	fmt.Printf("  Коефіцієнт прискорення K = T1/Tp:       %.2f\n", m.Speedup)
	fmt.Printf("  Задіяно процесорів:                     %d з %d\n", m.Used, m.Total)
	fmt.Printf("  Ефективність щодо задіяних K/%d:         %.2f\n", m.Used, m.EffUsed)
	fmt.Printf("  Ефективність щодо всіх K/%d:             %.2f\n", m.Total, m.EffTotal)
}
