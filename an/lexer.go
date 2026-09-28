package an

import (
	"fmt"
	"unicode"
)

type Lexer interface {
	Tokenize(raw string) ([]Token, []Error)
	Clear()
}

type lexer struct {
	errors []Error
}

func NewLexer() Lexer {
	return &lexer{}
}

func (l *lexer) Clear() {
	l.errors = []Error{}
}

func (l *lexer) Tokenize(raw string) ([]Token, []Error) {
	var tokens []Token
	runes := []rune(raw)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		switch {
		case unicode.IsDigit(r):
			token, idx := l.parseDigits(i, runes)
			tokens = append(tokens, token)
			i = idx
		case unicode.IsLetter(r):
			token, idx := l.parseLetters(i, runes)
			tokens = append(tokens, token)
			i = idx
		case unicode.IsSpace(r):
		case r == '-':
			token, idx := l.parseMinus(i, runes)
			tokens = append(tokens, token)
			i = idx
		case r == '+' || r == '*' || r == '/':
			tokens = append(tokens, NewToken(TokenOP, string(r), i))
		case r == '(':
			tokens = append(tokens, NewToken(TokenOB, string(r), i))
		case r == ')':
			tokens = append(tokens, NewToken(TokenCB, string(r), i))
		default:
			l.errors = append(l.errors, Error{Pos: i, Msg: fmt.Sprintf("Неочікуваний символ '%c'", r)})
		}
	}
	tokens = append(tokens, NewToken(TokenEND, "", len(runes)))

	return tokens, l.errors
}

func (l *lexer) parseDigits(idx int, runes []rune) (Token, int) {
	digits := make([]rune, 0)
	isDot := false
	i := idx
	for ; i < len(runes); i++ {
		r := runes[i]
		if r == '.' {
			if isDot {
				l.errors = append(l.errors, (Error{Msg: "Зайва крапка у числі", Pos: i}))
			}
			isDot = true
		} else if !unicode.IsDigit(r) {
			break
		}
		digits = append(digits, r)
	}
	return NewToken(TokenCONST, string(digits), idx), i - 1
}

func (l *lexer) parseLetters(idx int, runes []rune) (Token, int) {
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

func (l *lexer) parseMinus(idx int, runes []rune) (Token, int) {
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
