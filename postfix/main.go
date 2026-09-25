package postfix

import . "pzcs-lab1/an"

type Postfix interface {
	Convert([]Token) []Token
}

type postfix struct{}

func NewPostfix() Postfix {
	return &postfix{}
}

func (p *postfix) Convert(tokens []Token) []Token {
	output := make([]Token, 0, len(tokens))
	operations := make([]Token, 0, len(tokens))

	for _, t := range tokens {
		switch t.Type {
		case TokenVAR, TokenCONST:
			output = append(output, t)
		case TokenFN, TokenUM, TokenOB:
			operations = append(operations, t)
		case TokenOP:
			for i := len(operations) - 1; i >= 0; i-- {
				top := operations[i]
				if (top.Type == TokenOP || top.Type == TokenUM) && getOpPriority(top) >= getOpPriority(t) {
					output = append(output, top)
					operations = operations[:len(operations)-1]
					continue
				}
				break
			}
			operations = append(operations, t)
		case TokenCB:
			for top := operations[len(operations)-1]; top.Type != TokenOB; top = operations[len(operations)-1] {
				output = append(output, top)
				operations = operations[:len(operations)-1]
			}
			operations = operations[:len(operations)-1]

			if len(operations) > 0 {
				if top := operations[len(operations)-1]; top.Type == TokenFN {
					output = append(output, top)
					operations = operations[:len(operations)-1]
				}
			}
		}
	}

	for i := len(operations) - 1; i >= 0; i-- {
		output = append(output, operations[i])
	}

	return output
}

var opPriority = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

func getOpPriority(t Token) int {
	if t.Type == TokenUM {
		return 3
	}
	return opPriority[t.Text]
}
