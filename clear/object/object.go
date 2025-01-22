package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/ajtroup1/clear/ast"
)

type ObjectType string

const (
	NULL_OBJ  = "NULL"
	ERROR_OBJ = "ERROR"

	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ   = "FLOAT"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"
	ARRAY_OBJ   = "ARRAY"
	HASH_OBJ    = "HASH"

	BREAK_OBJ    = "BREAK"
	CONTINUE_OBJ = "CONTINUE"

	RETURN_VALUE_OBJ = "RETURN_VALUE"

	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ  = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Line() int
	Col() int
}

type Position struct {
	Line int
	Col  int
}

type Integer struct {
	Position
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Line() int        { return i.Position.Line }
func (i *Integer) Col() int         { return i.Position.Col }

type Float struct {
	Position
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Line() int        { return f.Position.Line }
func (f *Float) Col() int         { return f.Position.Col }

type Boolean struct {
	Position
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Line() int        { return b.Position.Line }
func (b *Boolean) Col() int         { return b.Position.Col }

type String struct {
	Position
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) Line() int        { return s.Position.Line }
func (s *String) Col() int         { return s.Position.Col }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) Line() int        { return 0 }
func (n *Null) Col() int         { return 0 }

type ReturnValue struct {
	Position
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Line() int        { return rv.Position.Line }
func (rv *ReturnValue) Col() int         { return rv.Position.Col }

type Error struct {
	Position
	Message string
	Context string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Line() int        { return e.Position.Line }
func (e *Error) Col() int         { return e.Position.Col }

type Function struct {
	Position
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (f *Function) Line() int { return f.Position.Line }
func (f *Function) Col() int  { return f.Position.Col }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Position
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Line() int        { return b.Position.Line }
func (b *Builtin) Col() int         { return b.Position.Col }

type Array struct {
	Position
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range ao.Elements {
		elements = append(elements, el.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (ao *Array) Line() int { return ao.Position.Line }
func (ao *Array) Col() int  { return ao.Position.Col }

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	Position
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
func (h *Hash) Line() int { return h.Position.Line }
func (h *Hash) Col() int  { return h.Position.Col }

type Hashable interface {
	HashKey() HashKey
}

type Continue struct {
	Position
}

func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }
func (c *Continue) Inspect() string  { return "continue" }
func (c *Continue) Line() int        { return c.Position.Line }
func (c *Continue) Col() int         { return c.Position.Col }

type Break struct {
	Position
}

func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) Inspect() string  { return "break" }
func (b *Break) Line() int        { return b.Position.Line }
func (b *Break) Col() int         { return b.Position.Col }
