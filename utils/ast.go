package utils

import (
	"errors"
	"fmt"

	"github.com/kr/pretty"
)

type BinaryExpr struct {
	Left  interface{}
	Op    Token
	Right interface{}
}

type BetweenExpr struct {
	Expr  interface{}
	Lower interface{}
	Upper interface{}
}

type BinaryExprCombine struct {
	Left  interface{}
	Op    Token
	Right interface{}
}
type AST struct {
	Columns      []Token
	From         Token
	Where        Expr
	WhereCombine BinaryExprCombine
}

func Expect(token Token, expectedType TokenType) (bool, error) {
	if token.Type != expectedType {
		return false, errors.New(fmt.Sprintf("Syntax error: No SELECT statement found"))
	}
	return true, nil
}

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

func ParseFrom(tokens []Token, pointer int) (Token, int) {
	return tokens[pointer], pointer
}

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

	(*p).RegisterPrefix(TokenNumber, 0)
	(*p).RegisterPrefix(TokenIdent, 0)
	(*p).RegisterInfix(TokenOr, 100)
	(*p).RegisterBetweenInfix(TokenBetween, 400)
	(*p).RegisterInfix(TokenAnd, 200)
	(*p).RegisterInfix(TokenGreater, 500)
	(*p).RegisterInfix(TokenGreaterEqual, 500)
	(*p).RegisterInfix(TokenLess, 500)
	(*p).RegisterInfix(TokenLessEqual, 500)
	(*p).RegisterInfix(TokenEqual, 500)
	(*p).RegisterInfix(TokenNotEqual, 500)
	fmt.Println("whereTokens", whereTokens)
	ast := p.ParseExpression(0)
	// fmt.Printf("ast: %#v", pretty(ast))
	fmt.Printf("ast:%# v", pretty.Formatter(ast))

	return ast, pointer
	// binaryExpr := BinaryExpr{}
	// fmt.Println("binaryExpr", binaryExpr.Left)
	// // binaryExprCombine := BinaryExprCombine{}
	// isCombineExpr := false
	// for pointer < len(tokens) {
	// 	token := tokens[pointer]
	// 	if IsOperator(fmt.Sprintf("%v", token.Value)) {
	// 		binaryExpr.Op = token
	// 	} else {
	// 		if binaryExpr.Left == nil {
	// 			binaryExpr.Left = token
	// 		} else if binaryExpr.Right == nil {
	// 			binaryExpr.Right = token
	// 		}
	// 	}

	// 	if token.Type == TokenAnd {
	// 		if _, ok := binaryExpr.Left.(BinaryExpr); !ok {
	// 			binaryExpr.Left = binaryExpr
	// 		}
	// 		isCombineExpr = true
	// 	}

	// 	// _, ok := binaryExprCombine.Right.(BinaryExpr)
	// 	// if isCombineExpr && binaryExpr.Right == nil {
	// 	// 	binaryExpr.Right = binaryExpr
	// 	// }

	// 	pointer++
	// }
	// if isCombineExpr {
	// 	// _, ok := binaryExprCombine.Right.(BinaryExpr)
	// 	// if ok {
	// 	// 	binaryExprCombine.Right = binaryExpr
	// 	// }
	// 	if binaryExpr.Right == nil {
	// 		binaryExpr.Right = binaryExpr
	// 	}
	// }
	// return binaryExpr, pointer
}

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
	pointer = endIdx + 1
	ast.Where = where
	// ast.WhereCombine = whereCombine

	return ast, nil

}
