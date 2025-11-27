package utils

import (
	"fmt"
)

// func ParseOperator(char string) Token {
// 	switch char {
// 	case "+":
// 		{
// 			return Token{
// 				Type:  TokenPlus,
// 				Value: char,
// 			}
// 		}
// 	case "-":
// 		{
// 			return Token{
// 				Type:  TokenMinus,
// 				Value: char,
// 			}
// 		}
// 	case "*":
// 		{
// 			return Token{
// 				Type:  TokenMultiply,
// 				Value: char,
// 			}
// 		}
// 	case "/":
// 		{
// 			return Token{
// 				Type:  TokenDivide,
// 				Value: char,
// 			}
// 		}
// 	default:
// 		{
// 			return Token{}
// 		}
// 	}
// }

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
	// Eval() int
}

// type BinaryExpr struct {
// 	Left  interface{}
// 	Op    Token
// 	Right interface{}
// }

func NewParser(tokens []Token) *Parser {
	p := &Parser{
		tokens:  tokens,
		current: tokens[0],
		pointer: 0,
		opTable: map[TokenType]OpInfo{},
	}
	return p
}

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

func (p *Parser) RegisterBetweenInfix(tokenType TokenType, lbp int) {
	(*p).opTable[tokenType] = OpInfo{
		lbp: lbp,
		led: func(left interface{}) interface{} {
			// op := p.tokens[p.pointer-1]
			// right := p.ParseExpression(lbp)
			ranges := p.ParseBetweenExpression()
			return BetweenExpr{
				Expr:  left,
				Lower: ranges[0],
				Upper: ranges[1],
			}
		},
	}
}

func (p *Parser) RegisterPrefix(tokenType TokenType, lbp int) {
	(*p).opTable[tokenType] = OpInfo{
		lbp: lbp,
		nud: func(p *Parser) Token {
			return p.tokens[p.pointer-1]
		},
	}
}

func (p *Parser) Advance() {
	if p.pointer < len(p.tokens) {
		p.pointer++
		p.current = p.tokens[p.pointer]
	}
}

func (p *Parser) GetLBP(tokenType TokenType) int {
	return p.opTable[tokenType].lbp
}

func (p *Parser) ParseExpression(minBp int) Expr {
	t := p.current
	p.Advance()
	var left Expr
	if t.Type == TokenNumber || t.Type == TokenIdent {
		left = p.opTable[t.Type].nud(p)
	}

	for p.current.Type != TokenEOF && p.GetLBP(p.current.Type) > minBp {
		t = p.current
		p.Advance()
		left = p.opTable[t.Type].led(left)
	}

	return left
}

func (p *Parser) ParseBetweenExpression() []interface{} {
	ranges := []interface{}{}
	for len(ranges) < 2 {
		t := p.current
		if t.Type != TokenAnd {
			ranges = append(ranges, t)
		}
		p.Advance()
	}
	return ranges
}

// func (ast *BinaryExpr) Eval() int {
// 	var left int
// 	left := ast.Left.Eval()
// 	right := ast.Left.Eval()
// 	op := ast.Op
// 	// fmt.Printf(left)
// 	return Compute(left, op.Type, right)
// }

// func Eval(ast interface{}) int {
// 	token, isToken := ast.(Token)
// 	if isToken {
// 		number, _ := strconv.Atoi(fmt.Sprintf("%v", token.Value))
// 		return number
// 	}

// 	binary, _ := ast.(BinaryExpr)
// 	left := Eval(binary.Left)
// 	right := Eval(binary.Right)
// 	op := binary.Op.Type
// 	return Compute(left, op, right)
// }

func Eval(ast interface{}, row []string, headerIndex map[string]int) int {

	//If ast is instance of BetweenExpr
	betweenExpr, isBetweenExpr := ast.(BetweenExpr)
	if isBetweenExpr {
		return ComputeBetween(betweenExpr, row, headerIndex)
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
			number, _ := StringToInt(row[fieldIdx])
			return number
		}
	}

	// If ast is expression
	binary, _ := ast.(BinaryExpr)
	left := Eval(binary.Left, row, headerIndex)
	right := Eval(binary.Right, row, headerIndex)
	op := binary.Op.Type
	return Compute(left, op, right)
}

func Compute(left int, op TokenType, right int) int {
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

// func Compute(left int, op TokenType, right int) int {
// 	switch op {
// 	case TokenPlus:
// 		{
// 			return left + right
// 		}
// 	case TokenMinus:
// 		{
// 			return left - right
// 		}
// 	case TokenMultiply:
// 		{
// 			return left * right
// 		}
// 	case TokenDivide:
// 		{
// 			return left / right
// 		}
// 	default:
// 		{
// 			return 0
// 		}
// 	}
// }

// func main() {
// 	exp := "1+2+3*3+4/2"
// 	tokens := []Token{}

// 	for _, char := range exp {
// 		if IsOperator(string(char)) {
// 			tokens = append(tokens, ParseOperator(string(char)))
// 		} else {
// 			number, _ := strconv.Atoi(
// 				string(char),
// 			)
// 			tokens = append(tokens, Token{
// 				Type:  TokenNumber,
// 				Value: number,
// 			})
// 		}

// 	}

// 	tokens = append(tokens, Token{
// 		Type:  TokenEOF,
// 		Value: nil,
// 	})

// 	p := NewParser(tokens)
// 	(*p).RegisterPrefix(TokenNumber, 0)
// 	(*p).RegisterInfix(TokenPlus, 1)
// 	(*p).RegisterInfix(TokenMinus, 1)
// 	(*p).RegisterInfix(TokenMultiply, 2)
// 	(*p).RegisterInfix(TokenDivide, 2)
// 	ast := p.ParseExpression(0)

// 	fmt.Println(tokens)
// 	fmt.Println(p)
// 	fmt.Printf("ast: %+v", ast)
// 	fmt.Println("Eval: ", Eval(ast))
// }
