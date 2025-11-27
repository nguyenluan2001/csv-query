package utils

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current Token
	pointer int
	opTable map[TokenType]OpInfo
}

type Nud func(p *Parser) Token
type Led func(left interface{}) any

type OpInfo struct {
	lbp int
	nud Nud
	led Led
}

type Expr interface {
}

func NewParser(tokens []Token) *Parser {
	p := &Parser{
		tokens:  tokens,
		current: tokens[0],
		pointer: 0,
		opTable: map[TokenType]OpInfo{},
	}
	return p
}

// Register binding power for BETWEEN operator
func (p *Parser) RegisterBetweenInfix(tokenType TokenType, lbp int) {
	(*p).opTable[tokenType] = OpInfo{
		lbp: lbp,
		led: func(left interface{}) interface{} {
			ranges := p.ParseBetweenExpression()
			return BetweenExpr{
				Expr:  left,
				Lower: ranges[0],
				Upper: ranges[1],
			}
		},
	}
}

// Register binding power for IN operator
func (p *Parser) RegisterInInfix(tokenType TokenType, lbp int) {
	(*p).opTable[tokenType] = OpInfo{
		lbp: lbp,
		led: func(left interface{}) interface{} {
			collector := p.ParseInExpression()
			return InExpr{
				Expr:      left,
				Collector: collector,
			}
		},
	}
}

// Register binding power for common operators
func (p *Parser) RegisterInfix(tokenType TokenType, lbp int) {
	(*p).opTable[tokenType] = OpInfo{
		lbp: lbp,
		led: func(left interface{}) interface{} {
			op := p.tokens[p.pointer-1]
			right := p.ParseExpression(lbp)
			return BinaryExpr{
				Left:  left,
				Op:    op,
				Right: right,
			}
		},
	}
}

// Register precedence for prefix token
func (p *Parser) RegisterPrefix(tokenType TokenType, lbp int) {
	(*p).opTable[tokenType] = OpInfo{
		lbp: lbp,
		nud: func(p *Parser) Token {
			return p.tokens[p.pointer-1]
		},
	}
}

// Moving to next token
func (p *Parser) Advance() {
	if p.pointer < len(p.tokens)-1 {
		p.pointer++
		p.current = p.tokens[p.pointer]
	}
}

// Get binding power of token
func (p *Parser) GetLBP(tokenType TokenType) int {
	return p.opTable[tokenType].lbp
}

// Pratt parser function
func (p *Parser) ParseExpression(minBp int) Expr {
	t := p.current
	p.Advance()
	var left Expr
	if t.Type == TokenNumber || t.Type == TokenIdent || t.Type == TokenString {
		left = p.opTable[t.Type].nud(p)
	}

	for p.current.Type != TokenEOF && p.GetLBP(p.current.Type) > minBp {
		t = p.current
		p.Advance()
		left = p.opTable[t.Type].led(left)
	}

	return left
}

// Handle parse tokens into BETWEEN expression format
func (p *Parser) ParseBetweenExpression() []interface{} {
	ranges := []interface{}{} // Store lower and upper tokens
	for len(ranges) < 2 {
		t := p.current
		if t.Type != TokenAnd {
			ranges = append(ranges, t)
		}
		p.Advance()
	}
	return ranges
}

// Handle parse tokens into IN expression format
func (p *Parser) ParseInExpression() []interface{} {
	if p.current.Type != TokenLParen {
		panic("Missing '(' symbol")
	}
	p.Advance()
	collector := []interface{}{} // Store tokens within "(" and ")"
	for p.current.Type != TokenRParen {
		if p.current.Type != TokenComma {
			collector = append(collector, p.current)
		}
		p.Advance()
	}
	p.Advance()
	return collector
}

// Handle validate input data based on AST expression
func Eval(ast interface{}, row []string, headerIndex map[string]int) interface{} {

	//If ast is instance of BetweenExpr
	betweenExpr, isBetweenExpr := ast.(BetweenExpr)
	if isBetweenExpr {
		return ComputeBetween(betweenExpr, row, headerIndex)
	}

	//If ast is instance of InExpr
	inExpr, isInExpr := ast.(InExpr)
	if isInExpr {
		return ComputeIn(inExpr, row, headerIndex)
	}

	token, isToken := ast.(Token)
	//If ast is token means the smallest unit to compute
	if isToken {
		if token.Type == TokenNumber {
			number, _ := StringToInt(token.Value)
			return number
		}
		if token.Type == TokenIdent {
			fieldIdx := headerIndex[fmt.Sprintf("%v", token.Value)]
			number, numberErr := StringToInt(row[fieldIdx])

			//If data can perform number
			if numberErr == nil {
				return number
			}
			return Stringify(row[fieldIdx])
		}
		if token.Type == TokenString {
			str := Stringify(token.Value)
			return str
		}
	}

	// If ast is expression => recurrive Eval()
	binary, _ := ast.(BinaryExpr)
	left := Eval(binary.Left, row, headerIndex)
	right := Eval(binary.Right, row, headerIndex)
	op := binary.Op.Type
	return Compute(left, op, right)
}

// Compute proxy to specific data type
func Compute(left interface{}, op TokenType, right interface{}) int {
	leftStr, isString := left.(string)
	rightStr, _ := right.(string)
	if isString {
		return ComputeString(leftStr, op, rightStr)
	}

	leftNumber, _ := left.(int)
	rightNumber, _ := right.(int)
	return ComputeNumber(leftNumber, op, rightNumber)
}

// Handle compute for string
func ComputeString(left string, op TokenType, right string) int {
	switch op {
	case TokenEqual:
		{
			return BooleanToInt(left == right)
		}
	default:
		{
			return 0
		}
	}
}

// Handle compute for number (int or float)
func ComputeNumber(left int, op TokenType, right int) int {
	fmt.Println("ComputeNumber", left, right)
	switch op {
	case TokenGreater:
		{
			return BooleanToInt(left > right)
		}
	case TokenGreaterEqual:
		{
			return BooleanToInt(left >= right)
		}
	case TokenLess:
		{
			return BooleanToInt(left < right)
		}
	case TokenLessEqual:
		{
			return BooleanToInt(left <= right)
		}
	case TokenEqual:
		{
			return BooleanToInt(left == right)
		}
	case TokenNotEqual:
		{
			return BooleanToInt(left != right)
		}
	case TokenAnd:
		{
			return BooleanToInt(left == 1 && right == 1)
		}
	case TokenOr:
		{
			return BooleanToInt(left == 1 || right == 1)
		}
	default:
		{
			return 0
		}
	}
}

// Handle compute for BETWEEN expression
func ComputeBetween(ast BetweenExpr, row []string, headerIndex map[string]int) int {
	tokenExpr, isTokenExpr := ast.Expr.(Token)
	tokenLower, isTokenLower := ast.Lower.(Token)
	tokenUpper, isTokenUpper := ast.Upper.(Token)
	if !isTokenExpr || !isTokenLower || !isTokenUpper {
		return 0
	}
	fieldIdx := headerIndex[fmt.Sprintf("%v", tokenExpr.Value)]
	number, _ := StringToInt(row[fieldIdx])
	lowerNumber, _ := StringToInt(Stringify(tokenLower.Value))
	upperNumber, _ := StringToInt(Stringify(tokenUpper.Value))

	return BooleanToInt(number >= lowerNumber && number <= upperNumber)
}

// Handle compute for IN expression
func ComputeIn(ast InExpr, row []string, headerIndex map[string]int) int {
	fmt.Println("ComputeIn")
	tokenExpr, isTokenExpr := ast.Expr.(Token)
	collector, isCollector := ast.Collector.([]interface{})
	if !isTokenExpr || !isCollector {
		return 0
	}
	fieldIdx := headerIndex[fmt.Sprintf("%v", tokenExpr.Value)]
	exprVal := row[fieldIdx]
	fmt.Println("exprVal", exprVal)
	isValid := 0
	for _, val := range collector {
		token, _ := val.(Token)
		if Stringify(exprVal) == Stringify(token.Value) {
			isValid = 1
		}
	}
	return isValid
}
