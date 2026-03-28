package scim_filter

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type TokenType string

const (
	typeOpenParen        TokenType = "typeOpenParen"
	typeCloseParen                 = "typeCloseParen"
	typeOpenSquareBrace            = "typeOpenSquareBrace"
	typeCloseSquareBrace           = "typeCloseSquareBrace"
	typeBareword                   = "typeBareword"
	typeQuotedStr                  = "typeQuotedStr"
	typeBool                       = "typeBool"
	typeNumber                     = "typeNumber"
	typeNull                       = "typeNull"
)

type Loc struct {
	Start int
	End   int
}

type Token struct {
	Loc   Loc
	Value string
	Type  TokenType
}

func (t Token) String() string {
	return fmt.Sprintf("<%s:%s:loc(%d:%d)>", t.Type, t.Value, t.Loc.Start, t.Loc.End)
}

type tokenizer struct {
	i int
	s string
}

func Tokenize(s string) ([]Token, error) {
	res := make([]Token, 0)
	t := tokenizer{i: 0, s: s}

	for t.i < len(t.s) {
		w, err := t.nextToken()
		if err != nil {
			return nil, err
		}
		if w.Type == "" {
			break
		}
		res = append(res, w)
	}
	return res, nil
}

func (t *tokenizer) nextToken() (Token, error) {
	switch t.s[t.i] {
	case '[':
		t.i += 1
		return Token{Type: typeOpenSquareBrace, Value: "[", Loc: Loc{Start: t.i - 1, End: t.i}}, nil
	case ']':
		t.i += 1
		return Token{Type: typeCloseSquareBrace, Value: "]", Loc: Loc{Start: t.i - 1, End: t.i}}, nil
	case '(':
		t.i += 1
		return Token{Type: typeOpenParen, Value: "(", Loc: Loc{Start: t.i - 1, End: t.i}}, nil
	case ')':
		t.i += 1
		return Token{Type: typeCloseParen, Value: ")", Loc: Loc{Start: t.i - 1, End: t.i}}, nil
	case '"':
		t.i += 1
		for j := t.i; t.i < len(t.s); t.i++ {
			if t.s[t.i] == '"' {
				token := Token{Type: typeQuotedStr, Value: t.s[j:t.i], Loc: Loc{Start: j, End: t.i}}
				t.i += 1
				return token, nil
			}
		}
		return Token{}, fmt.Errorf("unterminated quoted string")
	default:
		if unicode.IsSpace(rune(t.s[t.i])) {
			for t.i < len(t.s) && unicode.IsSpace(rune(t.s[t.i])) {
				t.i += 1
			}
			if t.i >= len(t.s) {
				return Token{}, nil
			}
			return t.nextToken()
		}
		j := t.i
		for ; t.i < len(t.s); t.i++ {
			if unicode.IsSpace(rune(t.s[t.i])) || strings.ContainsRune("()[]\"", rune(t.s[t.i])) {
				break
			}
		}
		w := t.s[j:t.i]
		if w == "null" {
			return Token{Type: typeNull, Value: "null", Loc: Loc{Start: j, End: t.i}}, nil
		} else if w == "true" {
			return Token{Type: typeBool, Value: "true", Loc: Loc{Start: j, End: t.i}}, nil
		} else if w == "false" {
			return Token{Type: typeBool, Value: "false", Loc: Loc{Start: j, End: t.i}}, nil
		} else if t.isNumber(w) {
			return Token{Type: typeNumber, Value: w, Loc: Loc{Start: j, End: t.i}}, nil
		}

		return Token{Type: typeBareword, Value: w, Loc: Loc{Start: j, End: t.i}}, nil
	}
}

func (*tokenizer) loc(i, j int) Loc {
	return Loc{Start: i, End: j}
}

func (*tokenizer) isNumber(s string) bool {
	// number parsing is hard
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
