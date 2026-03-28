package scim_filter

import (
	"reflect"
	"testing"
)

func TestTokenizeSCIMExpression(t *testing.T) {
	input := `userName eq "bjensen" and (emails[type eq "work"] pr or active eq true)`

	got, err := Tokenize(input)
	if err != nil {
		t.Fatalf("Tokenize returned error: %v", err)
	}

	want := []Token{
		{Type: typeBareword, Value: "userName", Loc: Loc{Start: 0, End: 8}},
		{Type: typeBareword, Value: "eq", Loc: Loc{Start: 9, End: 11}},
		{Type: typeQuotedStr, Value: "bjensen", Loc: Loc{Start: 13, End: 20}},
		{Type: typeBareword, Value: "and", Loc: Loc{Start: 22, End: 25}},
		{Type: typeOpenParen, Value: "(", Loc: Loc{Start: 26, End: 27}},
		{Type: typeBareword, Value: "emails", Loc: Loc{Start: 27, End: 33}},
		{Type: typeOpenSquareBrace, Value: "[", Loc: Loc{Start: 33, End: 34}},
		{Type: typeBareword, Value: "type", Loc: Loc{Start: 34, End: 38}},
		{Type: typeBareword, Value: "eq", Loc: Loc{Start: 39, End: 41}},
		{Type: typeQuotedStr, Value: "work", Loc: Loc{Start: 43, End: 47}},
		{Type: typeCloseSquareBrace, Value: "]", Loc: Loc{Start: 48, End: 49}},
		{Type: typeBareword, Value: "pr", Loc: Loc{Start: 50, End: 52}},
		{Type: typeBareword, Value: "or", Loc: Loc{Start: 53, End: 55}},
		{Type: typeBareword, Value: "active", Loc: Loc{Start: 56, End: 62}},
		{Type: typeBareword, Value: "eq", Loc: Loc{Start: 63, End: 65}},
		{Type: typeBool, Value: "true", Loc: Loc{Start: 66, End: 70}},
		{Type: typeCloseParen, Value: ")", Loc: Loc{Start: 70, End: 71}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokenize mismatch\n got: %#v\nwant: %#v", got, want)
	}
}

func TestTokenizeLiteralTypesAndWhitespace(t *testing.T) {
	input := "  userType   eq null or loginCount eq 42.5 and active eq false  "

	got, err := Tokenize(input)
	if err != nil {
		t.Fatalf("Tokenize returned error: %v", err)
	}

	want := []Token{
		{Type: typeBareword, Value: "userType", Loc: Loc{Start: 2, End: 10}},
		{Type: typeBareword, Value: "eq", Loc: Loc{Start: 13, End: 15}},
		{Type: typeNull, Value: "null", Loc: Loc{Start: 16, End: 20}},
		{Type: typeBareword, Value: "or", Loc: Loc{Start: 21, End: 23}},
		{Type: typeBareword, Value: "loginCount", Loc: Loc{Start: 24, End: 34}},
		{Type: typeBareword, Value: "eq", Loc: Loc{Start: 35, End: 37}},
		{Type: typeNumber, Value: "42.5", Loc: Loc{Start: 38, End: 42}},
		{Type: typeBareword, Value: "and", Loc: Loc{Start: 43, End: 46}},
		{Type: typeBareword, Value: "active", Loc: Loc{Start: 47, End: 53}},
		{Type: typeBareword, Value: "eq", Loc: Loc{Start: 54, End: 56}},
		{Type: typeBool, Value: "false", Loc: Loc{Start: 57, End: 62}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokenize mismatch\n got: %#v\nwant: %#v", got, want)
	}
}

func TestTokenizeQuotedStringError(t *testing.T) {
	_, err := Tokenize(`displayName eq "unterminated`)
	if err == nil {
		t.Fatal("Tokenize expected error for unterminated quoted string")
	}

	if err.Error() != "unterminated quoted string" {
		t.Fatalf("unexpected error: %v", err)
	}
}
