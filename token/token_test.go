package token

import "testing"

func TestLookupIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{input: "fn", expected: FUNCTION},
		{input: "let", expected: LET},
		{input: "true", expected: TRUE},
		{input: "false", expected: FALSE},
		{input: "if", expected: IF},
		{input: "else", expected: ELSE},
		{input: "return", expected: RETURN},
	}

	for _, tt := range tests {
		val := LookupIdentifier(tt.input)
		if val != tt.expected {
			t.Fatalf("LookupIdentifier does not return %s. got=%s", tt.expected, val)
		}
	}
}
