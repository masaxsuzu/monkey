package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	FreeScope    SymbolScope = "FREE"
	BuiltinScope SymbolScope = "BUILTIN"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer          *SymbolTable
	FreeSymbols    []Symbol
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	free := []Symbol{}
	return &SymbolTable{store: s, FreeSymbols: free}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

func (st *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: st.numDefinitions}
	if st.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	st.store[name] = symbol
	st.numDefinitions++
	return symbol
}

func (st *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	st.store[name] = symbol
	return symbol

}

func (st *SymbolTable) defineFree(original Symbol) Symbol {
	st.FreeSymbols = append(st.FreeSymbols, original)
	sym := Symbol{Name: original.Name, Index: len(st.FreeSymbols) - 1}
	sym.Scope = FreeScope
	st.store[original.Name] = sym
	return sym
}
func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	s, ok := st.store[name]
	if !ok && st.Outer != nil {
		s, ok := st.Outer.Resolve(name)
		if !ok {
			return s, ok
		}
		if s.Scope == GlobalScope || s.Scope == BuiltinScope {
			return s, ok
		}
		free := st.defineFree(s)
		return free, true
	}
	return s, ok
}
