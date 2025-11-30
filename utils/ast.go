package utils

import (
	"errors"
	"fmt"
	"slices"

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

type OrderBySingle struct {
	Field     string
	Direction TokenType
}

type AST struct {
	Columns []Token
	From    Token
	Where   Expr
	OrderBy []OrderBySingle
	Limit   int
	GroupBy []Token
}

// Check whether token existed
func Expect(token Token, expectedType TokenType) (bool, error) {
	if token.Type != expectedType {
		return false, errors.New(fmt.Sprintf("Syntax error: No SELECT statement found"))
	}
	return true, nil
}

// Parse aggregate function
func ParseAggregateFn(tokens []Token, pointer int) (Token, int) {
	if tokens[pointer+1].Type != TokenLParen || tokens[pointer+3].Type != TokenRParen {
		panic("Syntax error")
	}
	return tokens[pointer+2], pointer + 3
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
		isAggregateFn := IsAggregateFn(token)
		// isSum, _ := Expect(token, TokenSum)
		// isMax, _ := Expect(token, TokenMax)

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

		if isAggregateFn {
			fieldToken, endIdx := ParseAggregateFn(tokens, pointer)
			columns = append(columns, Token{
				Type:  token.Type,
				Value: fieldToken.Value,
			})
			pointer = endIdx
		}

		// if isSum {
		// 	fieldToken, endIdx := ParseAggregateFn(tokens, pointer)
		// 	columns = append(columns, Token{
		// 		Type:  TokenSum,
		// 		Value: fieldToken.Value,
		// 	})
		// 	pointer = endIdx
		// }

		// if isMax {
		// 	fieldToken, endIdx := ParseSum(tokens, pointer)
		// 	columns = append(columns, Token{
		// 		Type:  TokenMax,
		// 		Value: fieldToken.Value,
		// 	})
		// 	pointer = endIdx
		// }
		pointer++
	}
	pointer--
	return columns, pointer
}

// Get table of FROM statement
func ParseFrom(tokens []Token, pointer int) (Token, int) {
	return tokens[pointer], pointer
}

func CheckStopParseWhere(token Token) bool {
	stopTokens := []TokenType{TokenOrderBy, TokenEOF}
	return slices.Contains(stopTokens, token.Type)
}

// Get AST expression of WHERE statement
func ParseWhere(tokens []Token, pointer int) (Expr, int) {
	whereTokens := []Token{}
	for pointer < len(tokens) {
		token := tokens[pointer]
		if CheckStopParseWhere(token) {
			whereTokens = append(whereTokens, tokens[len(tokens)-1])
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

func ParseGroupBy(tokens []Token, pointer int) ([]Token, int) {
	groupBy := []Token{}
	for pointer < len(tokens) {
		token := tokens[pointer]
		if IsEOF(token) {
			break
		}
		if token.Type != TokenIdent && !IsComma(Stringify(token.Value)) {
			break
		}
		if token.Type == TokenComma {
			pointer++
			continue
		}
		if token.Type == TokenIdent {
			groupBy = append(groupBy, token)
		}
		pointer++
	}
	return groupBy, pointer
}

func ParseOrderBy(tokens []Token, pointer int) ([]OrderBySingle, int) {
	orderBy := []OrderBySingle{}
	for pointer < len(tokens) {
		token := tokens[pointer]
		if token.Type == TokenComma {
			pointer++
			continue
		}
		if token.Type == TokenIdent {
			orderBy = append(orderBy, OrderBySingle{
				Field:     Stringify(token.Value),
				Direction: TokenAsc,
			})
		}
		if token.Type == TokenAsc || token.Type == TokenDesc {
			orderBy[len(orderBy)-1].Direction = token.Type
		}
		pointer++
	}
	return orderBy, pointer
}

func ParseLimit(tokens []Token, pointer int) (int, int) {
	limit := -1
	for pointer < len(tokens) {
		token := tokens[pointer]

		if token.Type == TokenNumber {
			limit, _ = StringToInt(Stringify(token.Value))
			break
		}

		pointer++
	}
	return limit, pointer
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

	if isNext {
		// === Parse where ===
		where, endIdx := ParseWhere(tokens, pointer+1)
		ast.Where = where
		pointer = endIdx
	}

	// === Expect GROUP BY
	isNext, err = Expect(tokens[pointer], TokenGroupBy)

	if isNext {
		// === Parse GROUP BY ===
		groupBy, endIdx := ParseGroupBy(tokens, pointer+1)
		ast.GroupBy = groupBy
		pointer = endIdx
	}

	// === Expect ORDER BY ===
	isNext, err = Expect(tokens[pointer], TokenOrderBy)
	if isNext {
		orderBy, endIdx := ParseOrderBy(tokens, pointer)
		ast.OrderBy = orderBy
		pointer = endIdx + 1
	}

	// === Expect LIMIT ===
	isNext, err = Expect(tokens[pointer], TokenLimit)
	if isNext {
		limit, endIdx := ParseLimit(tokens, pointer)
		ast.Limit = limit
		pointer = endIdx + 1
	}

	return ast, nil

}
