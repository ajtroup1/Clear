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
		{"let arr = [1, 2, 3]; arrays.push(arr, 4);", []int{1, 2, 3, 4}},
		{"let arr = [\"Dog\", \"Cat\", \"Fish\"]; arrays.pop(arr);", "Fish"},
		{`let mathFunc = fn(x) { 
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
		{"let arr = [1, 2, 3]; arrays.last(arr);", 3},
		{"let arr = [1, 2, 3]; arrays.rest(arr);", []int{2, 3}},
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