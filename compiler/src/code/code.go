/*
	THE CODE PACKAGE defines the bytecode instructions that the compiler will generate.
	These instructions are what the VM will execute to run the compiled code.
	Each instruction is a single byte, and each instruction can have zero or more operands.
	An operand is a value that the instruction needs to do its job.
	For example, the OpConstant instruction needs an operand to know which constant to load onto the stack.
	The operands needed for the above example are: the constant opcode and a pointer to the constant within the stack.
*/

package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Represent a stream of byte instructions, each corresponding to an opcode
type Instructions []byte

// Disassemble the instructions into a human-readable format
// Useful for testing and debugging
func (ins Instructions) String() string {
	var out bytes.Buffer // Initialize a buffer to store the disassembled instructions
	i := 0

	// Loop through all the instructions
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}
		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}
	return out.String()
}

// Format the instruction based on the definition and operands
func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}
	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// Represent a single byte opcode
type Opcode byte

// Define all the different opcode types
const (
	// Load a constant value onto the stack
	// A constant value is a value that will never change at runtime
	// This means the value is solely determined at compile time
	OpConstant Opcode = iota

	// Stack instructions
	OpPop

	// Arithmetic instructions
	OpAdd
	OpSub
	OpMul
	OpDiv

	// Boolean instructions
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan

	// Prefix instructions
	OpMinus
	OpBang

	// Jump instructions
	OpJumpNotTruthy
	OpJump

	// Global bindings
	OpSetGlobal
	OpGetGlobal

	// Special instructions
	OpNull
)

// Define a human-readable representation of the opcode
type Definition struct {
	Name          string
	OperandWidths []int // Slice of ints representing how many bytes each opcode takes up
}

// Actually define each opcode definition with
// Its human-readable name
// The number of bytes each opcode takes up
var definitions = map[Opcode]*Definition{
	// Here, OpConstant is defined in the map with:
	// Its readable name
	// And the number of bytes it takes up
	// - So here it has one operand that is 2 bytes long
	// - This is signified by only have one list item in the OperandWidths slice and that value being 2
	OpConstant:      {"OpConstant", []int{2}},

	// Stack instructions
	OpPop:           {"OpPop", []int{}},

	// Arithmetic instructions
	OpAdd:           {"OpAdd", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},

	// Boolean instructions
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},

	// Prefix instructions
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},

	// Jump instructions
	// For conditional branching
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},
	
	// Global bindings
	// Takes one 2 byte operand which is the index of the global binding
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},

	// Special instructions
	OpNull:          {"OpNull", []int{}},
}

// Retrieve a Definition object based on a raw opcode
// Recieve nil if the opcode is not defined
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	// Retrieve the definition of the opcode
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	// Sum up the total length of the instruction based on all the operand widths in the instruction
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	// Create a byte slice with the length of the instruction
	instruction := make([]byte, instructionLen)
	// Set the first byte to the opcode
	instruction[0] = byte(op)
	offset := 1 // Tracks where the next operand should be written in the byte slice

	// Loop through all the operands and write them to the instruction byte slice
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}

// ReadOperands is the counterpart to Make
// It reads the operands from the instruction byte slice
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	// Initialize a slice to store the operands
	operands := make([]int, len(def.OperandWidths))
	offset := 0 // Tracks where the next operand should be read from in the byte slice

	// Loop through all the operand widths in the definition
	for i, width := range def.OperandWidths {
		// Read the operand based on the width
		switch width {
		case 2:
			// Assume the operand is a 16-bit unsigned integer
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width // Make sure to increment the offset by the width of the operand
	}
	return operands, offset
}
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
