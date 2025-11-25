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
	// parentheses
	TokenLParen
	TokenRParen
)

func IsNumber(code int) bool {
	return code >= 48 && code <= 57
}

func IsOperator(char string) bool {
	operators := []string{">", "<", ">=", "<=", "=", "!=", "!"}
	return slices.Contains(operators, char)
}

func IsComma(char string) bool {
	return char == ","
}

func IsWhiteSpace(char string) bool {
	return char == " "
}

func IsIdentifier(code int) bool {
	isUppercase := code >= 65 && code <= 90
	isLowercase := code >= 97 && code <= 122
	return isUppercase || isLowercase
}

func IsStar(char string) bool {
	return char == "*"
}

func ParseIdentifier(startIdx int, sql string) (TokenType, string, int) {
	byteArr := []byte{}
	endIdx := startIdx
	for i := startIdx; i < len(sql); i++ {
		char := sql[i]
		code, _ := strconv.Atoi(
			fmt.Sprintf("%d", char),
		)
		if !IsIdentifier(code) || IsWhiteSpace(string(char)) {
			break
		}
		byteArr = append(byteArr, char)
		endIdx++

	}
	identifier := string(byteArr)
	endIdx--

	switch identifier {
	case "SELECT":
		{
			return TokenSelect, string(byteArr), endIdx
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
	case "OR":
		{
			return TokenOr, string(byteArr), endIdx
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
		if !IsNumber(code) || IsWhiteSpace(string(char)) {
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

func ParseOperator(startIdx int, sql string) (TokenType, string, int) {
	byteArr := []byte{}
	endIdx := startIdx
	for i := startIdx; i < len(sql); i++ {
		char := sql[i]
		if !IsOperator(string(char)) || IsWhiteSpace(string(char)) {
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
	case "!=":
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

		if IsNumber(code) {
			// Number
			value, endIdx := ParseNumber(pointer, sql)
			tokens = append(tokens, Token{
				Type:   TokenNumber,
				Value:  value,
				Pos:    pointer,
				EndPos: endIdx,
			})
			pointer = endIdx
		} else if IsIdentifier(code) {
			tokenType, value, endIdx := ParseIdentifier(pointer, sql)
			tokens = append(tokens, Token{
				Type:   tokenType,
				Value:  value,
				Pos:    pointer,
				EndPos: endIdx,
			})
			pointer = endIdx
		} else if IsComma(string(char)) {
			tokens = append(tokens, Token{
				Type:  TokenComma,
				Value: string(char),
				Pos:   pointer,
			})
		} else if IsOperator(string(char)) {
			tokenType, value, endIdx := ParseOperator(pointer, sql)
			tokens = append(tokens, Token{
				Type:   tokenType,
				Value:  value,
				Pos:    pointer,
				EndPos: endIdx,
			})
			pointer = endIdx

		} else if IsStar(string(char)) {
			tokens = append(tokens, Token{
				Type:  TokenStar,
				Value: string(char),
				Pos:   pointer,
			})
		}
		pointer++
	}
	return tokens
}
