package an

import (
	"fmt"
	"slices"
	"strings"
)

type Parser interface {
	Parse([]Token) []Error
	Clear()
}

type SharedData struct {
	brackets []int
}

type parser struct {
	state  StateType
	errors []Error
	sData  SharedData
}

func NewParser(initialState StateType) Parser {
	return &parser{
		state: initialState,
		sData: SharedData{brackets: make([]int, 0)},
	}
}

func (p *parser) Clear() {
	p.sData.brackets = make([]int, 0)
	p.errors = make([]Error, 0)
	p.state = StState
}

func (p *parser) Parse(tokens []Token) []Error {
	for _, token := range tokens {
		p.checkBrackets(token)
		p.checkFun(token)

		tr := StateTable[p.state][token.Type]
		if tr.Err != "" {
			msg := strings.ReplaceAll(tr.Err, "{t}", token.Text)
			p.errors = append(p.errors, Error{Pos: token.Pos, Msg: msg})
		}
		p.state = tr.Next
	}

	if len(p.sData.brackets) != 0 {
		for _, pos := range p.sData.brackets {
			p.errors = append(p.errors, Error{Pos: pos, Msg: "Відкрита дужка '('"})
		}
	}

	return p.errors
}

func (p *parser) checkBrackets(token Token) {
	if token.Type == TokenOB {
		p.sData.brackets = append(p.sData.brackets, token.Pos)
	} else if token.Type == TokenCB {
		if len(p.sData.brackets) == 0 {
			p.errors = append(p.errors, Error{Pos: token.Pos, Msg: "Зайва закриваюча дужка ')'"})
			return
		}
		p.sData.brackets = p.sData.brackets[:len(p.sData.brackets)-1]
	}
}

var ALLOWED_FN = []string{"abs", "cos", "sin"}

func (p *parser) checkFun(token Token) {
	if token.Type == TokenFN && !slices.Contains(ALLOWED_FN, token.Text) {
		p.errors = append(p.errors, Error{Pos: token.Pos, Msg: fmt.Sprintf("Невідома функція: '%s'", token.Text)})
	}
}
