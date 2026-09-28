package an

type Analyzer interface {
	Start(raw string) ([]Token, []Error)
}

type analyzer struct {
	lex    Lexer
	parser Parser
	errors []Error
}

func NewAnalyzer(lexer Lexer, parser Parser) Analyzer {
	return &analyzer{
		lex:    lexer,
		parser: parser,
	}
}

func (a *analyzer) Start(raw string) ([]Token, []Error) {
	tokens, errors := a.lex.Tokenize(raw)
	a.errors = append(a.errors, errors...)

	errors = a.parser.Parse(tokens)
	a.errors = append(a.errors, errors...)

	errors = a.errors
	a.clear()

	return tokens, errors
}

func (a *analyzer) clear() {
	a.lex.Clear()
	a.parser.Clear()
	a.errors = make([]Error, 0)
}
