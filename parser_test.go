package scim_filter

import (
	"strings"
	"testing"
)

func TestParserSimple(t *testing.T) {
	inputs := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "SimpleUserName",
			input:    `userName eq "bjensen"`,
			expected: `userName eq "bjensen"`,
		},
		{
			name:     "UserNameOrActive",
			input:    `userName eq "bjensen" or active eq true`,
			expected: `userName eq "bjensen" or active eq true`,
		},
		{
			name:     "UserNameOrActiveAndEmail",
			input:    `userName eq "bjensen" or active eq true and emails[type eq "work"]`,
			expected: `userName eq "bjensen" or active eq true and emails[type eq "work"]`,
		},
	}

	for _, test := range inputs {
		expr, err := Parse(test.input)
		if err != nil {
			t.Fatalf("%s: Parse returned error: %v", test.name, err)
		}
		if expr.String() != test.expected {
			t.Errorf("%s: Parse returned unexpected result: `%s`", test.name, expr.String())
		}
	}
}

func TestParserOperatorsAndLiterals(t *testing.T) {
	inputs := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "PresentExpression",
			input:    `emails pr`,
			expected: `emails pr`,
		},
		{
			name:     "SubAttributeAndStringOperators",
			input:    `name.familyName sw "Jen"`,
			expected: `name.familyName sw "Jen"`,
		},
		{
			name:     "ContainsOperator",
			input:    `userName co "js"`,
			expected: `userName co "js"`,
		},
		{
			name:     "EndsWithOperator",
			input:    `email.value ew "@example.com"`,
			expected: `email.value ew "@example.com"`,
		},
		{
			name:     "NotEqualsNull",
			input:    `manager ne null`,
			expected: `manager ne null`,
		},
		{
			name:     "GreaterThanNumber",
			input:    `meta.version gt 42`,
			expected: `meta.version gt 42`,
		},
		{
			name:     "GreaterThanOrEqualsNumber",
			input:    `meta.version ge 42.5`,
			expected: `meta.version ge 42.5`,
		},
		{
			name:     "LessThanNumber",
			input:    `meta.version lt 100`,
			expected: `meta.version lt 100`,
		},
		{
			name:     "LessThanOrEqualsBoolean",
			input:    `active le false`,
			expected: `active le false`,
		},
	}

	for _, test := range inputs {
		expr, err := Parse(test.input)
		if err != nil {
			t.Fatalf("%s: Parse returned error: %v", test.name, err)
		}
		if expr.String() != test.expected {
			t.Errorf("%s: Parse returned unexpected result: `%s`", test.name, expr.String())
		}
	}
}

func TestParserNestedExpressions(t *testing.T) {
	inputs := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ParenthesizedOrUnderAnd",
			input:    `active eq true and (userName eq "bjensen" or userName eq "jsmith")`,
			expected: `active eq true and (userName eq "bjensen" or userName eq "jsmith")`,
		},
		{
			name:     "NestedValuePathFilter",
			input:    `emails[type eq "work" and value co "@example.com"]`,
			expected: `emails[type eq "work" and value co "@example.com"]`,
		},
		{
			name:     "LogicalExpressionInsideValuePath",
			input:    `members[value eq "2819c223-7f76-453a-919d-413861904646" or display eq "Babs Jensen"]`,
			expected: `members[value eq "2819c223-7f76-453a-919d-413861904646" or display eq "Babs Jensen"]`,
		},
	}

	for _, test := range inputs {
		expr, err := Parse(test.input)
		if err != nil {
			t.Fatalf("%s: Parse returned error: %v", test.name, err)
		}
		if expr.String() != test.expected {
			t.Errorf("%s: Parse returned unexpected result: `%s`", test.name, expr.String())
		}
	}
}

func TestParserErrors(t *testing.T) {
	inputs := []struct {
		name        string
		input       string
		errContains string
	}{
		{
			name:        "UnknownOperator",
			input:       `userName xx "bjensen"`,
			errContains: `unknown operator`,
		},
		{
			name:        "MissingClosingParen",
			input:       `(userName eq "bjensen"`,
			errContains: `unexpected end of input`,
		},
		{
			name:        "InvalidAttributePath",
			input:       `name.formatted.extra eq "x"`,
			errContains: `invalid attribute path`,
		},
	}

	for _, test := range inputs {
		_, err := Parse(test.input)
		if err == nil {
			t.Fatalf("%s: Parse expected error", test.name)
		}
		if !strings.Contains(err.Error(), test.errContains) {
			t.Errorf("%s: error %q does not contain %q", test.name, err.Error(), test.errContains)
		}
	}
}
