package formatter_test

import (
	"github.com/masa-suzu/monkey/ast"
	"github.com/masa-suzu/monkey/evaluator"
	"github.com/masa-suzu/monkey/formatter"
	"github.com/masa-suzu/monkey/lexer"
	"github.com/masa-suzu/monkey/object"
	"github.com/masa-suzu/monkey/parser"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"10", "10;"},
		{`"a";`, `"a";`},
		{"10;20;", `10;
20;`},
		{"let x=1", "let x = 1;"},
		{"return 10", "return 10;"},
		{"if(true){1;}", `if(true) {
    1;
};`},
		{"if(true){3} else{false;}", `if(true) {
    3;
} else {
    false;
};`},
		{"fn(x,y,z){return x*y +z}", `fn(x, y, z) {
    return ((x * y) + z);
};`},
		{"[1,2,3]", "[1, 2, 3];"},
		{"[1,2,3] (2)", "[1, 2, 3](2);"},
		{
			`let f = fn(x,y,z){return x*y +z};
f(1,4,5);`,
			`let f = fn(x, y, z) {
    return ((x * y) + z);
};
f(1, 4, 5);`},
		{`{"a":1};`, `{"a":1};`},
		{`if(true){ if(false){1}else{1}}`,
			`if(true) {
    if(false) {
        1;
    } else {
        1;
    };
};`},
		{`fn(x,y,z){fn(x,y,z){return x*y +z}}`,
			`fn(x, y, z) {
    fn(x, y, z) {
        return ((x * y) + z);
    };
};`},
		{`let f = fn(x,y,z){fn(x,y,z){return x*y +z}}`,
			`let f = fn(x, y, z) {
    fn(x, y, z) {
        return ((x * y) + z);
    };
};`},
		{`let f = fn(x,y,z){if(true){return x*y +z}}`,
			`let f = fn(x, y, z) {
    if(true) {
        return ((x * y) + z);
    };
};`},
		{`macro(x){x}`, `macro(x) {
    x;
};`,
		},
		{`macro(x){fn(){x}}`, `macro(x) {
    fn() {
        x;
    };
};`,
		},
		{
			`
let p = macro(x) {
    quote(if (unquote(x)) {
    unquote(x);
} else {
    "not truthy";
});
};`,
			`let p = macro(x) {
    quote(if(unquote(x)) {
        unquote(x);
    } else {
        "not truthy";
    });
};`,
		},
	}

	for _, tt := range tests {
		p := parse(tt.input)
		got := formatter.Format(p, 0)
		if got != tt.want {
			t.Errorf("Format(%v)\ngot:\n%v\nwant:\n%v", tt.input, got, tt.want)
		}

		formatProgram := parse(got)
		e1 := eval(formatProgram)
		e2 := eval(p)
		if e1 != nil && e2 != nil {
			if e1.Inspect() != e2.Inspect() {
				t.Fatalf("eval(parse(Format(%v))) got %v, want %v", tt.input, e1, e2)
			}
		} else if !(e1 == nil && e2 == nil) {
			t.Fatalf("eval(parse(Format(%v))) got %v, want %v", tt.input, e1, e2)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func eval(p *ast.Program) object.Object {
	env := object.NewEnvironment()
	return evaluator.Eval(p, env)
}
