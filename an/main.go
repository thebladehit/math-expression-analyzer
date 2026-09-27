package an

import (
	"fmt"
	"unicode"
)

type TokenType string

const (
	TokenVAR   TokenType = "VAR"
	TokenOP    TokenType = "OP"
	TokenCONST TokenType = "CONST"
	TokenOB    TokenType = "OB"
	TokenCB    TokenType = "CB"
	TokenFN    TokenType = "FN"
	TokenUM    TokenType = "UNARY_M"
)

type Token struct {
	Type TokenType
	Text string
	Pos  int
}

func NewToken(tokenType TokenType, text string, pos int) Token {
	return Token{
		Type: tokenType,
		Text: text,
		Pos:  pos,
	}
}

type Analyzer interface {
	Start(raw string)
}

type analyzer struct {
	errors []Error
}

func NewAnalyzer() Analyzer {
	return &analyzer{}
}

func (a *analyzer) Start(raw string) {
	tokens := a.parseTokens(raw)
	fmt.Println(tokens)
	fmt.Println(a.errors)
}

func (a *analyzer) parseTokens(raw string) []Token {
	var tokens []Token
	runes := []rune(raw)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		switch {
		case unicode.IsDigit(r):
			token, idx := a.parseDigits(i, runes)
			tokens = append(tokens, token)
			i = idx
		case unicode.IsLetter(r):
			token, idx := a.parseLetters(i, runes)
			tokens = append(tokens, token)
			i = idx
		case unicode.IsSpace(r):
		case r == '-':
			token, idx := a.parseMinus(i, runes)
			tokens = append(tokens, token)
			i = idx
		case r == '+' || r == '*' || r == '/':
			tokens = append(tokens, NewToken(TokenOP, string(r), i))
		case r == '(':
			tokens = append(tokens, NewToken(TokenOB, string(r), i))
		case r == ')':
			tokens = append(tokens, NewToken(TokenCB, string(r), i))
		default:
			a.errors = append(a.errors, Error{Pos: i, Msg: fmt.Sprintf("Неочікуваний символ '%c'", r)})
		}

	}

	return tokens
}

func (a *analyzer) parseDigits(idx int, runes []rune) (Token, int) {
	digits := make([]rune, 0)
	isDot := false
	i := idx
	for ; i < len(runes); i++ {
		r := runes[i]
		if r == '.' {
			if isDot {
				a.errors = append(a.errors, (Error{Msg: "Зайва крапка у числі", Pos: i}))
			}
			isDot = true
		} else if !unicode.IsDigit(r) {
			break
		}
		digits = append(digits, r)
	}
	return NewToken(TokenCONST, string(digits), idx), i - 1
}

func (a *analyzer) parseLetters(idx int, runes []rune) (Token, int) {
	letters := make([]rune, 0)
	i := idx
	for ; i < len(runes); i++ {
		r := runes[i]
		if !unicode.IsLetter(r) {
			break
		}
		letters = append(letters, r)
	}
	str := string(letters)
	if len(letters) == 1 {
		return NewToken(TokenVAR, str, idx), i - 1
	}
	return NewToken(TokenFN, str, idx), i - 1
}

func (a *analyzer) parseMinus(idx int, runes []rune) (Token, int) {
	if idx == 0 {
		return NewToken(TokenUM, "-", idx), idx
	}

	prevR := runes[idx-1]
	for i := idx - 1; i >= 0; i-- {
		if prevR != ' ' {
			break
		}
		prevR = runes[i]
	}
	if prevR == '(' || prevR == '+' || prevR == '-' || prevR == '*' || prevR == '/' || prevR == ' ' {
		return NewToken(TokenUM, "-", idx), idx
	}

	return NewToken(TokenOP, "-", idx), idx
}
