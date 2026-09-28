package an

type TokenType string

const (
	TokenVAR   TokenType = "VAR"
	TokenOP    TokenType = "OP"
	TokenCONST TokenType = "CONST"
	TokenOB    TokenType = "OB"
	TokenCB    TokenType = "CB"
	TokenFN    TokenType = "FN"
	TokenUM    TokenType = "UNARY_M"
	TokenEND   TokenType = "END"
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
