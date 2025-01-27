/*
	THE COMPILER PACKAGE defines the compiler that will take in an AST and generate bytecode.
*/

package compiler

import (
	"github.com/ajtroup1/clear-compiler/src/ast"
	"github.com/ajtroup1/clear-compiler/src/code"
	"github.com/ajtroup1/clear-compiler/src/object"
)

// Overall structure for the state and data of the compiler
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

// Create a new compiler
func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, statement := range node.Statements {
			err := c.Compile(statement)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return nil
		}
	case *ast.IntegerLiteral:
		value := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(value))
	}

	return nil
}

// Helper to add a constant to the compiler's constant pool and return its index on the stack
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// "emit" simply means generate an instruction and add it to results
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	// Generate the instruction given an opcode and operands
	ins := code.Make(op, operands...)

	// Add the instruction to the compiler's instruction list and return its index in the slice
	pos := c.addInstruction(ins)
	return pos
}

// Helper to add an instruction to the compiler's instruction list
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
