package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	an "pzcs-lab1/an"
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
