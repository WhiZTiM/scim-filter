package scim_filter

import "testing"

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
