package an

import (
	"slices"
	"testing"
)

func newTestAnalyzer() Analyzer {
	return NewAnalyzer(NewLexer(), NewParser(StState))
}

func errorPositions(errs []Error) []int {
	pos := make([]int, 0, len(errs))
	for _, e := range errs {
		pos = append(pos, e.Pos)
	}
	slices.Sort(pos)
	return pos
}

func TestValidExpressions(t *testing.T) {
	cases := []string{
		"a+b",
		"3.14",
		"-a",
		"x*-y",
		"((a))",
		"-(a+-b)",
		"sin(x)*2",
		"sin(x)*(2.5-y)/abs(-z)",
		"x1+my_var*2",
		"total/count",
		"sin(angle)*radius",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			_, errs := newTestAnalyzer().Start(input)
			if len(errs) != 0 {
				t.Errorf("очікувався коректний вираз, отримано помилки: %v", errs)
			}
		})
	}
}

// Позиції рахуються з 0, як у Error.Pos.
func TestErrorPositions(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []int
	}{
		// початок виразу
		{"операція на початку", "*a", []int{0}},
		{"порожній вираз", "", []int{0}},
		{"закриваюча дужка на початку", ")a", []int{0, 0, 1}},

		// кінець виразу
		{"операція в кінці", "a+", []int{2}},
		{"функція без аргументу", "sin x", []int{4}},

		// середина виразу
		{"подвійна операція", "a++b", []int{2}},
		{"немає операції між дужками", "(a)(b)", []int{3}},
		{"немає операції перед дужкою", "2(a)", []int{1}},
		{"операція після відкриваючої дужки", "(*a)", []int{1}},
		{"операція перед закриваючою дужкою", "(a+)", []int{3}},

		// дужки
		{"порожні дужки", "()", []int{1}},
		{"незакрита дужка", "((a)", []int{0}},
		{"лише відкриваюча дужка", "(", []int{0, 1}},
		{"зайва закриваюча дужка", "a+b)", []int{3}},

		// написання
		{"зайва крапка", "1.2.3", []int{3}},
		{"дужка після змінної (невідома функція)", "foo(x)", []int{3}},
		{"ім'я починається з цифри", "2abc", []int{1}},
		{"недопустимий символ", "a@b", []int{1, 2}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, errs := newTestAnalyzer().Start(tc.input)
			if got := errorPositions(errs); !slices.Equal(got, tc.want) {
				t.Errorf("%q: позиції помилок %v, очікувалось %v\n помилки: %v", tc.input, got, tc.want, errs)
			}
		})
	}
}

func TestTokenTypes(t *testing.T) {
	cases := []struct {
		input string
		want  []TokenType
	}{
		{"sin(x)", []TokenType{TokenFN, TokenOB, TokenVAR, TokenCB, TokenEND}},
		{"-(a)-2", []TokenType{TokenUM, TokenOB, TokenVAR, TokenCB, TokenOP, TokenCONST, TokenEND}},
		{"a*-b", []TokenType{TokenVAR, TokenOP, TokenUM, TokenVAR, TokenEND}},
		{"3.14 + y", []TokenType{TokenCONST, TokenOP, TokenVAR, TokenEND}},
		{"my_var2*sinx", []TokenType{TokenVAR, TokenOP, TokenVAR, TokenEND}},
		{"cos(foo)", []TokenType{TokenFN, TokenOB, TokenVAR, TokenCB, TokenEND}},
		{"", []TokenType{TokenEND}},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			tokens, _ := NewLexer().Tokenize(tc.input)
			got := make([]TokenType, 0, len(tokens))
			for _, tok := range tokens {
				got = append(got, tok.Type)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("типи токенів %v, очікувалось %v", got, tc.want)
			}
		})
	}
}

// Той самий аналізатор не має переносити помилки чи стан з попереднього виклику.
func TestAnalyzerReuse(t *testing.T) {
	a := newTestAnalyzer()
	if _, errs := a.Start("(a++"); len(errs) == 0 {
		t.Fatal("очікувались помилки для \"(a++\"")
	}
	if _, errs := a.Start("a+b"); len(errs) != 0 {
		t.Errorf("після попереднього виклику коректний вираз дав помилки: %v", errs)
	}
}
