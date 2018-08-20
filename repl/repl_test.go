package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestStart(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"\"Hello monkey!\"",
			"Hello monkey!\n",
		},
		{
			"let x = true;x;",
			"true\n",
		},

		{
			"let y = 2;y;",
			"2\n",
		},
		{
			"let add = fn(x,y) { return x + y ;}; add(2,3);",
			"5\n",
		},
		{
			"puts(\"monkey\");",
			"null\n",
		},
	}

	for _, tt := range tests {
		r := strings.NewReader(tt.input)
		w := &fakeWriter{Buffer: bytes.NewBuffer(nil)}

		Start(r, w, "")
		out := w.String()
		if out != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, out)
		}
	}
}

type fakeWriter struct {
	Buffer *bytes.Buffer
}

func (w *fakeWriter) Write(p []byte) (n int, err error) {
	w.Buffer.Write(p)
	return 0, nil
}

func (w *fakeWriter) String() string {
	return w.Buffer.String()
}
