package pkg

import (
	"slices"
)

type ParserFrom struct {
	tokens  []Token
	pointer int
	current Token
	opTable map[TokenType]OpInfoFrom
}

type NudFrom func() Token
type LedFrom func(left interface{}) any

type OpInfoFrom struct {
	lbp int
	nud NudFrom
	led LedFrom
}

type JoinExpr struct {
	Type      Token
	Left      interface{}
	Right     interface{}
	Condition interface{}
}

type TableIdentifier struct {
	Table Token
	Field Token
}

func NewParserFrom(tokens []Token) *ParserFrom {
	p := &ParserFrom{
		tokens:  tokens,
		current: tokens[0],
		pointer: 0,
		opTable: map[TokenType]OpInfoFrom{},
	}
	return p
}

// Register binding power for common operators
func (p *ParserFrom) RegisterInfix(tokenType TokenType, lbp int) {
	(*p).opTable[tokenType] = OpInfoFrom{
		lbp: lbp,
		led: func(left interface{}) interface{} {
			op := p.tokens[p.pointer-1]
			right := p.ParseOnCondition(lbp)
			return BinaryExpr{
				Left:  left,
				Op:    op,
				Right: right,
			}
		},
	}
}

// Register precedence for prefix token
func (p *ParserFrom) RegisterPrefix(tokenType TokenType, lbp int) {
	(*p).opTable[tokenType] = OpInfoFrom{
		lbp: lbp,
		nud: func() Token {
			return p.tokens[p.pointer-1]
		},
	}
}

// Get binding power of token
func (p *ParserFrom) GetLBP(tokenType TokenType) int {
	return p.opTable[tokenType].lbp
}

func (p *ParserFrom) Advance() {
	if p.pointer < len(p.tokens)-1 {
		p.pointer++
		p.current = p.tokens[p.pointer]
	}
}

func (p *ParserFrom) ParseJoin(joinToken Token, left interface{}) JoinExpr {
	right, condition := p.ParserFromExpression(0)
	return JoinExpr{
		Type:      joinToken,
		Left:      left,
		Right:     right,
		Condition: condition,
	}
}

func (p *ParserFrom) ParseTableIdentifier() TableIdentifier {
	table := p.current
	p.Advance()

	if p.current.Type != TokenDot {
		panic("Syntax error")
	}

	p.Advance()

	field := p.current
	return TableIdentifier{
		Table: table,
		Field: field,
	}
}

func (p *ParserFrom) ParseOnCondition(minBP int) interface{} {

	var left interface{}
	t := p.current
	if t.Type == TokenIdent {
		left = p.ParseTableIdentifier()
		p.Advance()
	}

	if t.Type == TokenNumber || t.Type == TokenString {
		p.Advance()
		left = p.opTable[t.Type].nud()
	}

	for p.current.Type != TokenEOF && p.GetLBP(p.current.Type) > minBP {
		t = p.current
		p.Advance()
		left = p.opTable[t.Type].led(left)
	}

	return left
}

func CheckJoinToken(token Token) bool {
	joinTokenType := []TokenType{TokenJoin, TokenLeftJoin, TokenRightJoin}
	return slices.Contains(joinTokenType, token.Type)
}

func (p *ParserFrom) ParserFromExpression(minBP int) (interface{}, interface{}) {
	t := p.current
	p.Advance()
	var left interface{}
	var condition interface{}

	if t.Type == TokenIdent {
		left = t
	}

	if p.current.Type == TokenOn {
		p.Advance()
		condition = p.ParseOnCondition(0)
		// condition = nil
	}

	for CheckJoinToken(p.current) && condition == nil {
		joinToken := p.current
		p.Advance()
		left = p.ParseJoin(joinToken, left)
	}
	return left, condition
}
