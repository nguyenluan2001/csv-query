package utils

import (
	"fmt"
	"slices"
	"strconv"
)

type TokenType int
type TokenValue any
type Token struct {
	Type   TokenType
	Value  any
	Pos    int
	EndPos int
}

const (
	TokenEOF TokenType = iota
	TokenIdent
	TokenSelect
	TokenAs
	TokenNumber
	TokenString
	TokenComma
	TokenStar
	TokenEqual
	TokenNotEqual
	TokenLess
	TokenGreater
	TokenLessEqual
	TokenGreaterEqual
	TokenAnd
	TokenOr
	TokenFrom
	TokenWhere
	TokenLimit
	TokenBetween
	TokenIn
	// parentheses
	TokenLParen
	TokenRParen

	TokenOrderBy
	TokenAsc
	TokenDesc

	TokenGroupBy
	TokenSum
	TokenAverage
	TokenCount
	TokenMax
	TokenMin

	TokenJoin
	TokenLeftJoin
	TokenRightJoin
	TokenOn
	TokenDot
)

func IsNumber(code int) bool {
	return code >= 48 && code <= 57
}

func IsSingleQuote(char string) bool {
	return char == `'`
}

func IsOperator(char string) bool {
	operators := []string{">", "<", ">=", "<=", "=", "<>", "!"}
	return slices.Contains(operators, char)
}

func IsComma(char string) bool {
	return char == ","
}

func IsDot(char string) bool {
	return char == "."
}

func IsWhiteSpace(char string) bool {
	return char == " "
}

func IsNewLine(char string) bool {
	return char == `\n`
}

func IsTab(char string) bool {
	return char == `\t`
}

func IsSkipParsing(char string) bool {
	return IsWhiteSpace(char) || IsNewLine(char) || IsTab(char)
}

func IsHyphen(char string) bool {
	return char == "-"
}

func IsIdentifier(code int) bool {
	isUppercase := code >= 65 && code <= 90
	isLowercase := code >= 97 && code <= 122
	return isUppercase || isLowercase
}

func IsStar(char string) bool {
	return char == "*"
}

func IsLeftParen(char string) bool {
	return char == "("
}

func IsRightParen(char string) bool {
	return char == ")"
}

func IsUnderscore(char string) bool {
	return char == "_"
}

func IsOrderByStart(str string) bool {
	return str == "ORDER"
}

func IsOrderBy(str string) bool {
	return str == "ORDER BY"
}

func IsGroupByStart(str string) bool {
	return str == "GROUP"
}

func IsGroupBy(str string) bool {
	return str == "GROUP BY"
}

func IsLeftJoinStart(str string) bool {
	return str == "LEFT"
}

func IsLeftJoin(str string) bool {
	return str == "LEFT JOIN"
}

func IsRightJoinStart(str string) bool {
	return str == "RIGHT"
}

func IsRightJoin(str string) bool {
	return str == "RIGHT JOIN"
}

func IsSum(str string) bool {
	return str == "SUM"
}

func IsAggregateFn(token Column) bool {
	aggregateFnType := []TokenType{
		TokenSum, TokenCount, TokenMax, TokenMin, TokenAverage,
	}
	return slices.Contains(aggregateFnType, token.Type)
}

func IsEOF(token Token) bool {
	return token.Type == TokenEOF
}

func ParseIdentifier(startIdx int, sql string) (TokenType, string, int) {
	byteArr := []byte{}
	endIdx := startIdx
	isHaveUnderscore := false
	for i := startIdx; i < len(sql); i++ {
		char := sql[i]
		code, _ := strconv.Atoi(
			fmt.Sprintf("%d", char),
		)

		if IsHyphen(string(char)) {
			panic(fmt.Sprintf("Syntax error at position %v", i))
		}

		if IsUnderscore(string(char)) {
			isHaveUnderscore = true
		}

		if IsOrderByStart(string(byteArr)) {
			if IsOrderBy(string(byteArr)) {
				break
			}
		} else if IsGroupByStart(string(byteArr)) {
			if IsGroupBy(string(byteArr)) {
				break
			}
		} else if IsLeftJoinStart(string(byteArr)) {
			if IsLeftJoin(string(byteArr)) {
				break
			}
		} else if IsRightJoinStart(string(byteArr)) {
			if IsRightJoin(string(byteArr)) {
				break
			}
		} else if !IsIdentifier(code) {
			// If have underscore => Only stop parse when meet "," or whitespace
			if isHaveUnderscore {
				if IsComma(string(char)) || IsSkipParsing(string(char)) || IsRightParen(string(char)) || IsDot(string(char)) {
					break
				}
			} else {
				if IsComma(string(char)) || IsSkipParsing(string(char)) || IsDot(string(char)) {
					break
				}

				if !IsNumber(code) {
					break
				}
			}
		}

		byteArr = append(byteArr, char)
		endIdx++
	}
	identifier := string(byteArr)
	fmt.Println("identifier", identifier, startIdx, endIdx)
	endIdx--

	switch identifier {
	case "SELECT":
		{
			return TokenSelect, string(byteArr), endIdx
		}
	case "AS":
		{
			return TokenAs, string(byteArr), endIdx
		}
	case "FROM":
		{
			return TokenFrom, string(byteArr), endIdx
		}
	case "WHERE":
		{
			return TokenWhere, string(byteArr), endIdx
		}
	case "AND":
		{
			return TokenAnd, string(byteArr), endIdx
		}
	case "BETWEEN":
		{
			return TokenBetween, string(byteArr), endIdx
		}
	case "IN":
		{
			return TokenIn, string(byteArr), endIdx
		}
	case "ORDER BY":
		{
			return TokenOrderBy, string(byteArr), endIdx
		}
	case "GROUP BY":
		{
			return TokenGroupBy, string(byteArr), endIdx
		}
	case "ASC":
		{
			return TokenAsc, string(byteArr), endIdx
		}
	case "DESC":
		{
			return TokenDesc, string(byteArr), endIdx
		}
	case "OR":
		{
			return TokenOr, string(byteArr), endIdx
		}
	case "LIMIT":
		{
			return TokenLimit, string(byteArr), endIdx
		}
	case "SUM":
		{
			return TokenSum, string(byteArr), endIdx
		}
	case "COUNT":
		{
			return TokenCount, string(byteArr), endIdx
		}
	case "AVG":
		{
			return TokenAverage, string(byteArr), endIdx
		}
	case "MAX":
		{
			return TokenMax, string(byteArr), endIdx
		}
	case "MIN":
		{
			return TokenMin, string(byteArr), endIdx
		}
	case "JOIN":
		{
			return TokenJoin, string(byteArr), endIdx
		}
	case "LEFT JOIN":
		{
			return TokenLeftJoin, string(byteArr), endIdx
		}
	case "RIGHT JOIN":
		{
			return TokenRightJoin, string(byteArr), endIdx
		}
	case "ON":
		{
			return TokenOn, string(byteArr), endIdx
		}
	default:
		{

			return TokenIdent, string(byteArr), endIdx
		}
	}
}

func ParseNumber(startIdx int, sql string) (int, int) {
	numberArr := []byte{}
	endIdx := startIdx
	for i := startIdx; i < len(sql); i++ {
		char := sql[i]
		code, _ := strconv.Atoi(
			fmt.Sprintf("%d", char),
		)
		if !IsNumber(code) || IsSkipParsing(string(char)) {
			break
		}
		numberArr = append(numberArr, char)
		endIdx += 1

	}
	endIdx--
	number, _ := strconv.Atoi(
		string(numberArr),
	)
	return number, endIdx
}

func ParseString(startIdx int, sql string) (string, int) {
	strBytes := []byte{}
	endIdx := startIdx
	for i := startIdx + 1; i < len(sql); i++ {
		char := sql[i]
		if IsSingleQuote(string(char)) {
			endIdx += 1
			break
		}
		strBytes = append(strBytes, char)
		endIdx += 1

	}
	return string(strBytes), endIdx
}

func ParseOperator(startIdx int, sql string) (TokenType, string, int) {
	byteArr := []byte{}
	endIdx := startIdx
	for i := startIdx; i < len(sql); i++ {
		char := sql[i]
		if !IsOperator(string(char)) || IsSkipParsing(string(char)) {
			break
		}
		byteArr = append(byteArr, char)
		endIdx++

	}
	identifier := string(byteArr)
	fmt.Println("identifier", identifier)
	endIdx--

	switch identifier {
	case ">":
		{
			return TokenGreater, string(byteArr), endIdx
		}
	case ">=":
		{
			return TokenGreaterEqual, string(byteArr), endIdx
		}
	case "<":
		{
			return TokenLess, string(byteArr), endIdx
		}
	case "<=":
		{
			return TokenLessEqual, string(byteArr), endIdx
		}
	case "<>":
		{
			return TokenNotEqual, string(byteArr), endIdx
		}
	default:
		{

			return TokenEqual, string(byteArr), endIdx
		}
	}
}

func Tokenizer(sql string) []Token {
	tokens := []Token{}
	pointer := 0
	for pointer < len(sql) {
		char := sql[pointer]
		code, _ := strconv.Atoi(
			fmt.Sprintf("%d", char),
		)

		// If character is a number => parsing to get the whole number
		if IsNumber(code) {
			value, endIdx := ParseNumber(pointer, sql)
			tokens = append(tokens, Token{
				Type:   TokenNumber,
				Value:  value,
				Pos:    pointer,
				EndPos: endIdx,
			})
			pointer = endIdx
		} else if IsSingleQuote(string(char)) { // If character is single quote => Parsing to get whole string within 2 single quotes
			value, endIdx := ParseString(pointer, sql)
			tokens = append(tokens, Token{
				Type:   TokenString,
				Value:  value,
				Pos:    pointer,
				EndPos: endIdx,
			})
			pointer = endIdx
		} else if IsIdentifier(code) { // If character is a string => Parsing to get the whole identifier
			tokenType, value, endIdx := ParseIdentifier(pointer, sql)
			tokens = append(tokens, Token{
				Type:   tokenType,
				Value:  value,
				Pos:    pointer,
				EndPos: endIdx,
			})
			pointer = endIdx
		} else if IsComma(string(char)) { // If character is a comma => Parsing to get comma
			tokens = append(tokens, Token{
				Type:  TokenComma,
				Value: string(char),
				Pos:   pointer,
			})
		} else if IsDot(string(char)) {
			tokens = append(tokens, Token{
				Type:  TokenDot,
				Value: string(char),
				Pos:   pointer,
			})

		} else if IsOperator(string(char)) { // If charater is a operator => Parsing to get operator
			tokenType, value, endIdx := ParseOperator(pointer, sql)
			tokens = append(tokens, Token{
				Type:   tokenType,
				Value:  value,
				Pos:    pointer,
				EndPos: endIdx,
			})
			pointer = endIdx

		} else if IsStar(string(char)) { // If character is "*" => Parsing to get "*"
			tokens = append(tokens, Token{
				Type:  TokenStar,
				Value: string(char),
				Pos:   pointer,
			})
		} else if IsLeftParen(string(char)) {
			tokens = append(tokens, Token{
				Type:  TokenLParen,
				Value: string(char),
				Pos:   pointer,
			})
		} else if IsRightParen(string(char)) {
			tokens = append(tokens, Token{
				Type:  TokenRParen,
				Value: string(char),
				Pos:   pointer,
			})
		}
		pointer++
	}
	tokens = append(tokens, Token{
		Type:  TokenEOF,
		Value: nil,
	})
	return tokens
}
