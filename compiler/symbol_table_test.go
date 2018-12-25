package compiler

import "testing"

func TestDefine(t *testing.T) {
	want := map[string]Symbol{
		"a": Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": Symbol{Name: "b", Scope: GlobalScope, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	if a != want["a"] {
		t.Errorf("want a=%+v, got=%+v", want["a"], a)
	}

	b := global.Define("b")
	if b != want["b"] {
		t.Errorf("want a=%+v, got=%+v", want["a"], b)
	}
}

func TestResolve(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	want := map[string]Symbol{
		"a": Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": Symbol{Name: "b", Scope: GlobalScope, Index: 1},
	}

	for _, sym := range want {
		got, ok := global.Resolve(sym.Name)

		if !ok {
			t.Errorf("name %s not resolved", sym.Name)
		}

		if got != sym {
			t.Errorf("want %s to resolve to %+v, got=%+v", sym.Name, sym, got)
		}
	}

}
