package vm

import (
	"fmt"

	"github.com/ajtroup1/clear-compiler/src/code"
	"github.com/ajtroup1/clear-compiler/src/compiler"
	"github.com/ajtroup1/clear-compiler/src/object"
)

// Set the maximum stack size
const StackSize = 2048

// Define global boolean objects
// This is because we only need one instance of each
// and the instances are immutable
var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

// Overall state and data for the VM
type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // Stack Pointer
}

// Create a new VM with the given bytecode representation
func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

// Run the VM
// Sequentially execute the instructions in the bytecode
func (vm *VM) Run() error {
	// Loop through the instructions and keep and instruction pointer
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip]) // Get the opcode at the current instruction pointer
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		// Primitive infix operations
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		// Primitive boolean values
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		// Primitive comparison operations
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		// Pop the top of the stack
		// Used after every expression statement
		case code.OpPop:
			vm.pop()
		}
	}

	return nil
}

// Found an infix expression, execute it according to the operator and data types
// If the data types or operator are not supported, return an error
// Otherwise, execute the operation and push the result to the stack
func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	// Pop the top two elements from the stack
	// These should be the operands for the infix operation
	right := vm.pop()
	left := vm.pop()
	// Check the types of the operands
	leftType := left.Type()
	rightType := right.Type()
	// Integer operations
	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s",
		leftType, rightType)
}

// Evaluate the infix operation for two integers based on the operator
// Push the result to the stack
func (vm *VM) executeBinaryIntegerOperation(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	var result int64
	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}
	return vm.push(&object.Integer{Value: result})
}

// Execute a comparison operation based on the operator
func (vm *VM) executeComparison(op code.Opcode) error {
	// Pop the top two elements from the stack
	// These should be the operands for the comparison operation
	right := vm.pop()
	left := vm.pop()
	// Integer comparison
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}
	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)",
			op, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}
