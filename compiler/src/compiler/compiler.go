/*
	THE COMPILER PACKAGE defines the compiler that will take in an AST and generate bytecode.
	The compiler will:
		1. Traverse the AST
		2. Generate bytecode instructions
		3. Add the bytecode instructions to the compiler's instruction list
		4. Return the bytecode instructions and constants generated by the compiler
*/

package compiler

import (
	"fmt"

	"github.com/ajtroup1/clear-compiler/src/ast"
	"github.com/ajtroup1/clear-compiler/src/code"
	"github.com/ajtroup1/clear-compiler/src/object"
)

// Overall structure for the state and data of the compiler
type Compiler struct {
	instructions        code.Instructions  // The bytecode instructions generated by the compiler
	constants           []object.Object    // The constants generated by the compiler
	lastInstruction     EmittedInstruction // The last instruction generated by the compiler
	previousInstruction EmittedInstruction // The instruction before the last instruction generated by the compiler
	symbolTable         *SymbolTable       // The symbol table used by the compiler
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

// Create a new compiler
func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         NewSymbolTable(),
	}
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	// The first case and initial hit in this switch is the root node of the AST
	case *ast.Program:
		// Recursively compile all of the statements in the program
		for _, statement := range node.Statements {
			err := c.Compile(statement)
			if err != nil {
				return err
			}
		}
	// Compiling block statements is the same as compiling programs
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		// Compile the expression in the statement
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		// Pop the result of the expression off the stack
		// This is because we don't need the result of the expression for anything
		c.emit(code.OpPop)
	case *ast.LetStatement:
		// Firstly, compile the value that the variable is being assigned to
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		// Now, define the variable in the symbol table
		symbol := c.symbolTable.Define(node.Name.Value)
		// Emit the `OpSetGlobal` instruction with the index of the symbol in the symbol table
		c.emit(code.OpSetGlobal, symbol.Index)
	case *ast.InfixExpression:
		if node.Operator == "<" { // Special case for the `<` operator
			// Compile the right side of the infix expression first
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			// Then compile the left side of the infix expression
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			// Emit the `OpGreaterThan` instruction
			// This is because we want to check if the left side is less than the right side
			// So this bascially reverses the order of the operands and uses the same opcode
			c.emit(code.OpGreaterThan)
			return nil
		}

		// Compile the left side of the infix expression first
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		// Then compile the right side of the infix expression
		err = c.Compile(node.Right)
		if err != nil {
			return nil
		}

		// Simply emit the appropriate opcode based on the operator
		// And that's it
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.PrefixExpression:
		// Compile the vale that the prefix operator is being applied to
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		// Emit the appropriate opcode based on the operator
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IfExpression:
		// Must compile the condition first
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Jump instructions are required to skip either the consequence or the alternative
		// based on the result of the condition

		// Emit an `OpJumpNotTruthy` with a bogus value
		// This will be used to know where to skip to if the condition is false
		// This is because we don't know the position of after the consequence yet
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		// Compile the consequence as a block statement
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		// Remove the last `OpPop` instruction if it exists
		// This is because the last instruction in the consequence block is always an `OpPop`
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		// Emit an `OpJump` with a bogus value
		jumpPos := c.emit(code.OpJump, 9999)

		// Calculate the position of the consequence and set the operand of the `OpJumpNotTruthy` instruction
		afterConsequencePos := len(c.instructions)
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		// Handle the alternative if present
		if node.Alternative == nil {
			// If there is no alternative, simply emit an `OpNull` instruction
			c.emit(code.OpNull)
		} else {
			// Compile the alternative as a block statement
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			// Again, remove the last `OpPop` instruction if it exists
			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		}

		// Calculate the position of after the alternative and set the operand of the `OpJump` instruction
		// So, if the condition is true, it runs the consequence and skips the alternative
		afterAlternativePos := len(c.instructions)
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.Identifier:
		// An identifier was found, so we must be invoking a variable's value
		// Find that variable in the symbol table
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		// Emit the `OpGetGlobal` instruction with the index of the symbol in the symbol table
		c.emit(code.OpGetGlobal, symbol.Index)
	case *ast.IntegerLiteral:
		// Just emit an `OpConstant` instruction with the integer in object form
		value := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(value))
	case *ast.Boolean:
		// Just emit an `OpTrue` or `OpFalse` instruction based on the boolean's value
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	}

	// No errors occurred, so return nil
	return nil
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
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

	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}
	c.previousInstruction = previous
	c.lastInstruction = last
}

// Helper to add an instruction to the compiler's instruction list
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

// Return the information generated by the compiler as bytecode
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
