package compiler

// Define the different types of scopes
// Using strings instead of integers for better readability
type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

// Define the Symbol struct
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int // The numerical position of the symbol in the symbol table
}

// The symbol table tracks all the symbols in the program
// It stores the symbols in a map where the key is the symbol's name
// It also stores the number of definitions in the symbol table, which is used to assign the index to each symbol
type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

// Define a new symbol in the symbol table
// TODO - Actually develop these methods once other scopes are implemented
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// Resolve (return) a symbol by name key
// If the symbol is not found, return a zeroed Symbol and false
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	return obj, ok
}
