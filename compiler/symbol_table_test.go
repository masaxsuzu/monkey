package compiler

import "testing"

func TestDefine(t *testing.T) {
	want := map[string]Symbol{
		"a": Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": Symbol{Name: "b", Scope: GlobalScope, Index: 1},
		"c": Symbol{Name: "c", Scope: LocalScope, Index: 0},
		"d": Symbol{Name: "d", Scope: LocalScope, Index: 1},
		"e": Symbol{Name: "e", Scope: LocalScope, Index: 0},
		"f": Symbol{Name: "f", Scope: LocalScope, Index: 1},
	}

	global := NewSymbolTable()
	l1 := NewEnclosedSymbolTable(global)
	l2 := NewEnclosedSymbolTable(l1)

	a := global.Define("a")
	if a != want["a"] {
		t.Errorf("want a=%+v, got=%+v", want["a"], a)
	}

	b := global.Define("b")
	if b != want["b"] {
		t.Errorf("want a=%+v, got=%+v", want["a"], b)
	}

	c := l1.Define("c")
	if c != want["c"] {
		t.Errorf("want c=%+v, got=%+v", want["c"], c)
	}

	d := l1.Define("d")
	if d != want["d"] {
		t.Errorf("want d=%+v, got=%+v", want["d"], d)
	}

	e := l2.Define("e")
	if e != want["e"] {
		t.Errorf("want e=%+v, got=%+v", want["e"], e)
	}

	f := l2.Define("f")
	if f != want["f"] {
		t.Errorf("want f=%+v, got=%+v", want["f"], f)
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

func TestResolveLocal(t *testing.T) {
	g := NewSymbolTable()
	g.Define("a")
	g.Define("b")
	l := NewEnclosedSymbolTable(g)
	l.Define("c")
	l.Define("d")

	want := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: "d", Scope: LocalScope, Index: 1},
	}

	for _, sym := range want {
		ret, ok := l.Resolve(sym.Name)

		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if ret != sym {
			t.Errorf("want %s to resolve to %+v, got=%+v", sym.Name, sym, ret)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	g := NewSymbolTable()
	g.Define("a")
	g.Define("b")

	l1 := NewEnclosedSymbolTable(g)
	l1.Define("c")
	l1.Define("d")

	l2 := NewEnclosedSymbolTable(l1)
	l2.Define("e")
	l2.Define("f")

	tests := []struct {
		table *SymbolTable
		want  []Symbol
	}{
		{
			table: l1,
			want: []Symbol{

				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			table: l2,
			want: []Symbol{

				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.want {
			ret, ok := tt.table.Resolve(sym.Name)

			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}

			if ret != sym {
				t.Errorf("want %s to resolve to %+v, got=%+v", sym.Name, sym, ret)
			}
		}
	}
}

func TestResolveFree(t *testing.T) {
	g := NewSymbolTable()
	g.Define("a")
	g.Define("b")

	l1 := NewEnclosedSymbolTable(g)
	l1.Define("c")
	l1.Define("d")

	l2 := NewEnclosedSymbolTable(l1)
	l2.Define("e")
	l2.Define("f")

	tests := []struct {
		table    *SymbolTable
		want     []Symbol
		wantFree []Symbol
	}{
		{
			table: l1,
			want: []Symbol{

				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
			wantFree: []Symbol{},
		},
		{
			table: l2,
			want: []Symbol{

				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: FreeScope, Index: 0},
				{Name: "d", Scope: FreeScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
			wantFree: []Symbol{
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.want {
			ret, ok := tt.table.Resolve(sym.Name)

			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}

			if ret != sym {
				t.Errorf("want %s to resolve to %+v, got=%+v", sym.Name, sym, ret)
			}
		}
		if len(tt.table.FreeSymbols) != len(tt.wantFree) {
			t.Errorf("wrong number of free symbols. got=%d, want=%d", len(tt.table.FreeSymbols), len(tt.wantFree))
		}

		for i, sym := range tt.wantFree {
			ret := tt.table.FreeSymbols[i]
			if ret != sym {
				t.Errorf("wrong free symbol. got=%+v, want=%+v", ret, sym)
			}
		}
	}
}

func TestResolveUnresolvableFree(t *testing.T) {
	g := NewSymbolTable()
	g.Define("a")

	l1 := NewEnclosedSymbolTable(g)
	l1.Define("c")

	l2 := NewEnclosedSymbolTable(l1)
	l2.Define("e")
	l2.Define("f")

	want := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "c", Scope: FreeScope, Index: 0},
		{Name: "e", Scope: LocalScope, Index: 0},
		{Name: "f", Scope: LocalScope, Index: 1},
	}

	for _, sym := range want {
		ret, ok := l2.Resolve(sym.Name)

		if !ok {
			t.Errorf("name %s not resolvavle", sym.Name)
			continue
		}
		if ret != sym {
			t.Errorf("expecte %s to resolve to %+v, got=%+v", sym.Name, sym, ret)
		}
	}
	expectedUnresolvable := []string{
		"b",
		"d",
	}
	for _, name := range expectedUnresolvable {
		_, ok := l2.Resolve(name)

		if ok {
			t.Errorf("name %s resolved, but was expected not to", name)
		}
	}
}

func TestResolveBuiltins(t *testing.T) {
	g := NewSymbolTable()
	s1 := NewEnclosedSymbolTable(g)
	s2 := NewEnclosedSymbolTable(s1)

	want := []Symbol{
		{Name: "a", Scope: BuiltinScope, Index: 0},
		{Name: "c", Scope: BuiltinScope, Index: 1},
		{Name: "e", Scope: BuiltinScope, Index: 2},
		{Name: "f", Scope: BuiltinScope, Index: 3},
	}

	for i, v := range want {
		g.DefineBuiltin(i, v.Name)
	}
	for _, table := range []*SymbolTable{g, s1, s2} {
		for _, sym := range want {
			ret, ok := table.Resolve(sym.Name)

			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}

			if ret != sym {
				t.Errorf("want %s to resolve to %+v, got=%+v", sym.Name, sym, ret)
			}
		}
	}

}
