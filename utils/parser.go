package main

import (
	"fmt"
	"slices"
	"strconv"
)

type TokenType int

type Token struct {
	Type  TokenType
	Value interface{}
}

const (
	TokenNumber = iota
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
)

func IsNumber(code int) bool {
	return code >= 48 && code <= 57
}

func IsOperator(char string) bool {
	operators := []string{"+", "-", "*", "/"}
	return slices.Contains(operators, char)
}

func ParseOperator(char string) Token {
	switch char {
	case "+":
		{
			return Token{
				Type:  TokenPlus,
				Value: char,
			}
		}
	case "-":
		{
			return Token{
				Type:  TokenMinus,
				Value: char,
			}
		}
	case "*":
		{
			return Token{
				Type:  TokenMultiply,
				Value: char,
			}
		}
	case "/":
		{
			return Token{
				Type:  TokenDivide,
				Value: char,
			}
		}
	default:
		{
			return Token{}
		}
	}
}

type Parser struct {
	tokens  []Token
	current Token
	pointer int
	opTable map[TokenType]OpInfo
}

type Nud func(p *Parser, t Token) Token
type Led func(p *Parser, left interface{}, t Token) BinaryExpr

type OpInfo struct {
	lbp int
	nud Nud
	led Led
}

type BinaryExpr struct {
	Left  interface{}
	Op    Token
	Right interface{}
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

func (p *Parser) RegisterInfix(tokenType TokenType, lbp int) {
	(*p).opTable[tokenType] = OpInfo{
		lbp: lbp,
		led: func(p *Parser, left interface{}, t Token) BinaryExpr {
			return BinaryExpr{
				Left:  left,
				Op:    t,
				Right: nil,
			}
		},
	}
}

func (p *Parser) ParseExpression(minBp int) interface{} {
	t := p.current
	return nil
}

func main() {
	exp := "1+2*3"
	tokens := []Token{}

	for _, char := range exp {
		if IsOperator(string(char)) {
			tokens = append(tokens, ParseOperator(string(char)))
		} else {
			number, _ := strconv.Atoi(
				string(char),
			)
			tokens = append(tokens, Token{
				Type:  TokenNumber,
				Value: number,
			})
		}

	}

	p := NewParser(tokens)
	(*p).RegisterInfix(TokenPlus, 1)
	(*p).RegisterInfix(TokenMinus, 1)
	(*p).RegisterInfix(TokenMultiply, 2)
	(*p).RegisterInfix(TokenDivide, 2)

	fmt.Println(tokens)
	fmt.Println(p)
}
