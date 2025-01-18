package modules

import (
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
		// Testing arrays with a large number of elements
		{
			"mod arrays: [len]; let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; arrays.len(arr);",
			10,
		},
		// Pushing elements to an array repeatedly to test large array functionality
		{
			"mod arrays: [push, len]; let arr = [1]; arrays.push(arr, 2); arrays.push(arr, 3); arrays.push(arr, 4); arrays.push(arr, 5); arrays.push(arr, 6); arrays.push(arr, 7); arrays.push(arr, 8); arrays.push(arr, 9); arrays.push(arr, 10); arrays.len(arr);",
			10,
		},
		// Testing pop on a large array
		{
			"mod arrays: [pop, len]; let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; arrays.pop(arr); arrays.len(arr);",
			9,
		},
		// Test pop until empty
		{
			"mod arrays: [pop, len]; let arr = [1, 2, 3]; arrays.pop(arr); arrays.pop(arr); arrays.pop(arr); arrays.len(arr);",
			0,
		},
		// Test first and last with large array
		{
			"mod arrays: [first, last]; let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; arrays.first(arr);",
			1,
		},
		{
			"mod arrays: [first, last]; let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; arrays.last(arr);",
			10,
		},
		// Testing rest with a large array
		{
			"mod arrays: [rest]; let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; arrays.rest(arr);",
			[]int{2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		// Testing push with a nested array
		{
			"mod arrays: [push, len]; let arr = [[1, 2], [3, 4]]; arrays.push(arr, [5, 6]); arrays.len(arr);",
			3,
		},
		// Test if `first` works correctly for a nested array
		{
			"mod arrays: [first]; let arr = [[1, 2], [3, 4], [5, 6]]; arrays.first(arr);",
			[]int{1, 2},
		},
		// Test if `last` works correctly for a nested array
		{
			"mod arrays: [last]; let arr = [[1, 2], [3, 4], [5, 6]]; arrays.last(arr);",
			[]int{5, 6},
		},
		// Test rest with a nested array
		{
			"mod arrays: [rest]; let arr = [[1, 2], [3, 4], [5, 6]]; arrays.rest(arr);",
			[][]int{{3, 4}, {5, 6}},
		},
		// Test if len works with both arrays and strings
		{
			"mod arrays: [len]: mod strings: [len]; let arr = [1, 2, 3]; let str = \"Hello\"; arrays.len(arr) + strings.len(str);",
			8,
		},
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
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testExpectedObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	Register(env)

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
