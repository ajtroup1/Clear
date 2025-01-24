package parser

import (
	"strings"
	"testing"

	"github.com/ajtroup1/compiled_clear/src/ast"
	"github.com/ajtroup1/compiled_clear/src/errorlogger"
	"github.com/ajtroup1/compiled_clear/src/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		// Basic variable declarations
		{"let x = 5;", "x", "5"},
		{"let y = true;", "y", "true"},
		{"let foobar = y;", "foobar", "y"},

		// Strings and booleans
		{"let name = \"John\";", "name", "John"},
		{"let isActive = false;", "isActive", "false"},

		// Arithmetic operations
		{"let sum = 5 + 10;", "sum", "(5 + 10);"},
		{"let product = 4 * 8;", "product", "(4 * 8);"},
		{"let difference = 15 - 3;", "difference", "(15 - 3);"},
		{"let quotient = 20 / 4;", "quotient", "(20 / 4);"},

		// Complex expressions
		{"let complex = 5 + 2 * 10;", "complex", "(5 + (2 * 10););"},
		{"let complex1 = 5 * 2 + 10;", "complex1", "((5 * 2); + 10);"},
		{"let complex2 = 5 + 2 * 10 + 3;", "complex2", "((5 + (2 * 10);); + 3);"},

		// Comparisons
		{"let isEqual = 5 == 5;", "isEqual", "(5 == 5);"},
		{"let isNotEqual = 10 != 5;", "isNotEqual", "(10 != 5);"},
		{"let isGreater = 15 > 10;", "isGreater", "(15 > 10);"},
		{"let isLesser = 5 < 10;", "isLesser", "(5 < 10);"},

		// Logical operators
		{"let andResult = true && false;", "andResult", "(true && false);"},
		{"let orResult = true || false;", "orResult", "(true || false);"},

		// // Arrays and indexing
		// {"let array = [1, 2, 3];", "array", "[1, 2, 3]"},
		// {"let firstElement = array[0];", "firstElement", "array[0]"},

		// // Functions
		// {
		// 	"let add = fn(x, y) { x + y; };",
		// 	"add",
		// 	"fn(x, y) { x + y; }",
		// },
		// {
		// 	"let result = add(5, 10);",
		// 	"result",
		// 	"add(5, 10)",
		// },

		// // Multiline input
		// {
		// 	`let x = 5;
		//     let y = 10;
		//     let sum = x + y;`,
		// 	"sum",
		// 	"x + y",
		// },

		// // Nested blocks
		// {
		// 	`let nested = fn(a) {
		//         if (a > 0) {
		//             return true;
		//         } else {
		//             return false;
		//         }
		//     };`,
		// 	"nested",
		// 	`fn(a) { if (a > 0) { return true; } else { return false; } }`,
		// },

		// // Edge cases
		// {"let emptyString = \"\";", "emptyString", ""},
		// {"let zero = 0;", "zero", "0"},
		// {"let largeNum = 1234567890;", "largeNum", "1234567890"},
		// {"let negNum = -42;", "negNum", "-42"},
		// {"let escaped = \"hello \\\"world\\\"\";", "escaped", "hello \"world\""},
	}

	for _, tt := range tests {
		testLetStatement(t, tt.input, tt.expectedIdentifier, tt.expectedValue)
	}
}

func testLetStatement(t *testing.T, input string, expectedIdentifier string, expectedValue interface{}) {
	el := errorlogger.New(strings.Split(input, "\n"), false)
	p := New(lexer.New(input, el, false), el)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *LetStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != expectedIdentifier {
		t.Fatalf("stmt.Name.Value not '%s'. got=%s", expectedIdentifier, stmt.Name.Value)
	}

	if stmt.Name.TokenLiteral() != expectedIdentifier {
		t.Fatalf("stmt.Name.TokenLiteral() not '%s'. got=%s", expectedIdentifier, stmt.Name.TokenLiteral())
	}

	if stmt.Value.ToString() != expectedValue {
		t.Fatalf("stmt.Value.TokenLiteral() not '%v'. got=%v", expectedValue, stmt.Value.ToString())
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", "5"},
		{"return true;", "true"},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		testReturnStatement(t, tt.input, tt.expectedValue)
	}
}

func testReturnStatement(t *testing.T, input string, expectedValue interface{}) {
	el := errorlogger.New(strings.Split(input, "\n"), false)
	p := New(lexer.New(input, el, false), el)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ReturnStatement. got=%T", program.Statements[0])
	}

	if stmt.TokenLiteral() != "return" {
		t.Fatalf("stmt.TokenLiteral not 'return'. got=%q", stmt.TokenLiteral())
	}

	if stmt.ReturnValue.ToString() != expectedValue {
		t.Fatalf("stmt.Value not '%v'. got=%v", expectedValue, stmt.ReturnValue.ToString())
	}
}