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
		{Type: typeBareword, Value: "userName"},
		{Type: typeBareword, Value: "eq"},
		{Type: typeQuotedStr, Value: "bjensen"},
		{Type: typeBareword, Value: "and"},
		{Type: typeOpenParen, Value: "("},
		{Type: typeBareword, Value: "emails"},
		{Type: typeOpenSquareBrace, Value: "["},
		{Type: typeBareword, Value: "type"},
		{Type: typeBareword, Value: "eq"},
		{Type: typeQuotedStr, Value: "work"},
		{Type: typeCloseSquareBrace, Value: "]"},
		{Type: typeBareword, Value: "pr"},
		{Type: typeBareword, Value: "or"},
		{Type: typeBareword, Value: "active"},
		{Type: typeBareword, Value: "eq"},
		{Type: typeBool, Value: "true"},
		{Type: typeCloseParen, Value: ")"},
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
		{Type: typeBareword, Value: "userType"},
		{Type: typeBareword, Value: "eq"},
		{Type: typeNull, Value: "null"},
		{Type: typeBareword, Value: "or"},
		{Type: typeBareword, Value: "loginCount"},
		{Type: typeBareword, Value: "eq"},
		{Type: typeNumber, Value: "42.5"},
		{Type: typeBareword, Value: "and"},
		{Type: typeBareword, Value: "active"},
		{Type: typeBareword, Value: "eq"},
		{Type: typeBool, Value: "false"},
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
