package utils

import (
	"errors"
	"fmt"

	"github.com/kr/pretty"
)

// Struct of common expression
type BinaryExpr struct {
	Left  interface{}
	Op    Token
	Right interface{}
}

// Struct of BETWEEN expression
type BetweenExpr struct {
	Expr  interface{}
	Lower interface{}
	Upper interface{}
}

// Struct of IN expression
type InExpr struct {
	Expr      interface{}
	Collector interface{}
}

type AST struct {
	Columns []Token
	From    Token
	Where   Expr
}

// Check whether token existed
func Expect(token Token, expectedType TokenType) (bool, error) {
	if token.Type != expectedType {
		return false, errors.New(fmt.Sprintf("Syntax error: No SELECT statement found"))
	}
	return true, nil
}

// Get columns of SELECT statement
func ParseColumns(tokens []Token, pointer int) ([]Token, int) {
	columns := []Token{}

	for pointer < len(tokens) {
		token := tokens[pointer]
		fmt.Println("pointer", pointer, token)
		isNext, _ := Expect(token, TokenIdent)
		isFromStmt, _ := Expect(token, TokenFrom)
		isComma := IsComma(fmt.Sprintf("%v", token.Value))
		isStar := IsStar(fmt.Sprintf("%v", token.Value))

		if isStar {
			columns = append(columns, token)
			pointer++
			break
		}

		if isFromStmt {
			break
		}

		if isNext && !isComma {
			columns = append(columns, token)
		}
		pointer++
	}
	pointer--
	return columns, pointer
}

// Get table of FROM statement
func ParseFrom(tokens []Token, pointer int) (Token, int) {
	return tokens[pointer], pointer
}

// Get AST expression of WHERE statement
func ParseWhere(tokens []Token, pointer int) (Expr, int) {
	whereTokens := []Token{}
	for pointer < len(tokens) {
		token := tokens[pointer]
		if IsEOF(token) {
			whereTokens = append(whereTokens, token)
			break
		}
		whereTokens = append(whereTokens, token)
		pointer++
	}

	p := NewParser(whereTokens)

	//Register operator precedence
	(*p).RegisterPrefix(TokenNumber, 0)
	(*p).RegisterPrefix(TokenString, 0)
	(*p).RegisterPrefix(TokenIdent, 0)
	(*p).RegisterInfix(TokenOr, 100)
	(*p).RegisterInfix(TokenAnd, 200)
	(*p).RegisterBetweenInfix(TokenBetween, 400)
	(*p).RegisterInInfix(TokenIn, 400)
	(*p).RegisterInfix(TokenGreater, 500)
	(*p).RegisterInfix(TokenGreaterEqual, 500)
	(*p).RegisterInfix(TokenLess, 500)
	(*p).RegisterInfix(TokenLessEqual, 500)
	(*p).RegisterInfix(TokenEqual, 500)
	(*p).RegisterInfix(TokenNotEqual, 500)

	// Start parse expression from min_bp=0
	ast := p.ParseExpression(0)

	fmt.Println("whereTokens", whereTokens)
	fmt.Printf("ast:%# v", pretty.Formatter(ast))

	return ast, pointer
}

// Build AST of whole query
func BuildAST(tokens []Token) (AST, error) {
	ast := AST{}
	pointer := 0
	//=== Expect SELECT ===
	isNext, err := Expect(tokens[pointer], TokenSelect)
	if !isNext {
		return ast, err
	}
	pointer++

	// === Parse columns ===
	columns, endIdx := ParseColumns(tokens, pointer)
	pointer = endIdx + 1
	ast.Columns = columns

	// === Expect FROM ===
	isNext, err = Expect(tokens[pointer], TokenFrom)
	if !isNext {
		return ast, err
	}
	pointer++

	// === Parse from ===
	from, endIdx := ParseFrom(tokens, pointer)
	pointer = endIdx + 1
	ast.From = from

	// === Expect WHERE
	isNext, err = Expect(tokens[pointer], TokenWhere)
	if !isNext {
		return ast, err
	}
	pointer++

	// === Parse where ===
	where, endIdx := ParseWhere(tokens, pointer)
	ast.Where = where
	pointer = endIdx + 1

	return ast, nil

}
