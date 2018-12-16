package main

import (
	"testing"
)

func TestStartRep(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`let hello = fn(){ "Hello, Monkey!"; }; hello();`,
			"Hello, Monkey!\n",
		},
		{
			`let hello = fn(){
				"Hello, Monkey!";
			};
			hello();`,
			"Hello, Monkey!\n",
		},
	}

	for _, tt := range tests {
		out := startRep(tt.input)
		if out != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, out)
		}
	}
}
