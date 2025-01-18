package modules

import (
	"fmt"
	"testing"

	"github.com/ajtroup1/clear/evaluator"
	"github.com/ajtroup1/clear/lexer"
	"github.com/ajtroup1/clear/object"
	"github.com/ajtroup1/clear/parser"
)

func TestArraysBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"mod arrays: [len]; let arr = [1, 2, 3]; arrays.len(arr);", 3},
		{"mod arrays: [push]; let arr = [1, 2, 3]; arrays.push(arr, 4);", []int{1, 2, 3, 4}},
		{"mod arrays: [pop]; let arr = [\"Dog\", \"Cat\", \"Fish\"]; arrays.pop(arr);", "Fish"},
		{`mod arrays: [first];
			let mathFunc = fn(x) { 
				x = x * 2; 
				x = (x - 7) * 2; 
				x = x * x + x; 
				x; 
			}; 
			let arr = [mathFunc(7), 2, 3]; 
			arrays.first(arr);
			`,
			210,
		},
		{"mod arrays: [last]; let arr = [1, 2, 3]; arrays.last(arr);", 3},
		{"mod arrays: [rest]; let arr = [1, 2, 3]; arrays.rest(arr);", []int{2, 3}},
		{"mod arrays: [len]; let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; arrays.len(arr);", 10},
		{"mod arrays: [push, len]; let arr = [1]; arrays.push(arr, 2); arrays.push(arr, 3); arrays.push(arr, 4); arrays.push(arr, 5); arrays.push(arr, 6); arrays.push(arr, 7); arrays.push(arr, 8); arrays.push(arr, 9); arrays.push(arr, 10); arrays.len(arr);", 10},
		{"mod arrays: [pop, len]; let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; arrays.pop(arr); arrays.len(arr);", 9},
		{"mod arrays: [pop, len]; let arr = [1, 2, 3]; arrays.pop(arr); arrays.pop(arr); arrays.pop(arr); arrays.len(arr);", 0},
		{"mod arrays: [first, last]; let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; arrays.first(arr);", 1},
		{"mod arrays: [first, last]; let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; arrays.last(arr);", 10},
		{"mod arrays: [rest]; let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; arrays.rest(arr);", []int{2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{"mod arrays: [push, len]; let arr = [[1, 2], [3, 4]]; arrays.push(arr, [5, 6]); arrays.len(arr);", 3},
		{"mod arrays: [first]; let arr = [[1, 2], [3, 4], [5, 6]]; arrays.first(arr);", []int{1, 2}},
		{"mod arrays: [last]; let arr = [[1, 2], [3, 4], [5, 6]]; arrays.last(arr);", []int{5, 6}},
		{"mod arrays: [rest]; let arr = [[1, 2], [3, 4], [5, 6]]; arrays.rest(arr);", [][]int{{3, 4}, {5, 6}}},
		{"mod arrays: [len]: mod strings: [len]; let arr = [1, 2, 3]; let str = \"Hello\"; arrays.len(arr) + strings.len(str);", 8},
		{"mod arrays: [reverse]; let arr = [1, 2, 3]; arrays.reverse(arr);", []int{3, 2, 1}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testExpectedObject(t, evaluated, tt.expected)
	}
}

func TestStringsBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"mod strings: [len]; strings.len(\"Hello, World!\");", 13},
		{"mod strings: [concat]; strings.concat(\"Hello\", \"World\");", "HelloWorld"},
		{"mod strings: [concatDelim]; strings.concatDelim(\", \", \"Hello\", \"World\");", "Hello, World"},
		{"mod strings: [concatDelim]; strings.concatDelim(\", \", \"Hello\", \"World\", \"How\", \"Are\", \"You\");", "Hello, World, How, Are, You"},
		{"mod strings: [split]; strings.split(\"Hello, World\", \", \");", []string{"Hello", "World"}},
		{"mod strings: [split]; strings.split(\"Hello, World, How, Are, You\", \", \");", []string{"Hello", "World", "How", "Are", "You"}},
		{"mod strings: [upper]; let str = \"hello, world\"; strings.upper(str);", "HELLO, WORLD"},
		{"mod strings: [lower] let str = \"HELLO, WORLD\"; strings.lower(str);", "hello, world"},
		{"mod strings: [replace]; strings.replace(\"Hello, World\", \"World\", \"Everyone\");", "Hello, Everyone"},
		{"mod strings: [replace]; strings.replace(\"Hello, World, World, World\", \"World\", \"Everyone\");", "Hello, Everyone, Everyone, Everyone"},
		{"mod strings: [trimSpace]; let str = \"  Hello, World  \"; strings.trimSpace(str);", "Hello, World"},
		{"mod strings: [trimSpace, len]; let str = \"  Hello, World  \"; strings.len(strings.trimSpace(str));", 12},
		{"mod strings: [trimPrefix]; let str = \"Hello, World\"; strings.trimPrefix(str, \"Hello\");", ", World"},
		{"mod strings: [trimSuffix]; let str = \"Hello, World\"; strings.trimSuffix(str, \"World\");", "Hello, "},
		{"mod strings: [trimPrefix, trimSuffix]; let str = \"Hello, World\"; strings.trimSuffix(strings.trimPrefix(str, \"Hello\"), \"World\");", ", "},
		{"mod strings: [hasPrefix]; let str = \"Hello, World\"; strings.hasPrefix(str, \"Hello\");", true},
		{"mod strings: [hasSuffix]; let str = \"Hello, World\"; strings.hasSuffix(str, \"World\");", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testExpectedObject(t, evaluated, tt.expected)
	}
}

func TestMathBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"mod math: [abs]; math.abs(5);", 5},
		{"mod math: [abs]; math.abs(-5);", 5},
		{"mod math: [abs]; math.abs(5.5);", 5.5},
		{"mod math: [abs]; math.abs(-5.5);", 5.5},
		{"mod math: [round]; math.round(5);", 5},
		{"mod math: [round]; math.round(5.5);", 6},
		{"mod math: [round]; math.round(5.4);", 5},
		{"mod math: [pow]; math.pow(2, 3);", 8},
		{"mod math: [pow]; math.pow(2.5, 3);", 15.625},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testExpectedObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	fmt.Print()
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	Register(env)
	// fmt.Printf("env: %+v\n", env)

	return evaluator.Eval(program, env)
}

func testExpectedObject(t *testing.T, obj object.Object, expected interface{}) {
	switch expected := expected.(type) {
	case int:
		testIntegerObject(t, obj, int64(expected))
	case string:
		str, ok := obj.(*object.String)
		if !ok {
			t.Errorf("object is not String. got=%T (%+v)", obj, obj)
			return
		}
		testStringObject(t, str, expected)
	case []int:
		arr, ok := obj.(*object.Array)
		if !ok {
			t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
			return
		}

		if len(arr.Elements) != len(expected) {
			t.Errorf("array has wrong num of elements. got=%d", len(arr.Elements))
			return
		}

		for i, expectedElem := range expected {
			testIntegerObject(t, arr.Elements[i], int64(expectedElem))
		}
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
}

func testStringObject(t *testing.T, obj *object.String, expected string) {
	if obj.Value != expected {
		t.Errorf("object has wrong value. got=%q, want=%q", obj.Value, expected)
	}
}
