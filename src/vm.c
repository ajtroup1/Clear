#include <stdio.h>

#include "common.h"
#include "vm.h"
#include "debug.h"

VM vm;

void initVM()
{
  resetStack();
}

static void resetStack()
{
  vm.stackTop = vm.stack;
}

void freeVM()
{
}

/*
  The most important function in Clear
*/
static InterpretResult run()
{
  // Macro that reads the byte currently pointed at by ip and then advances the instruction pointer
#define READ_BYTE() (*vm.ip++)
// Reads the next byte from the bytecode, treats the resulting number as an index, and looks up the corresponding Value in the chunk’s constant table.
#define READ_CONSTANT() (vm.chunk->constants.values[READ_BYTE()])

  for (;;)
  {
    uint8_t instruction;
    switch (instruction = READ_BYTE())
    {
    case OP_CONSTANT:
    {
      Value constant = READ_CONSTANT();
      push(constant);
      break;
    }
    case OP_RETURN:
    {
      printValue(pop());
      printf("\n");
      return INTERPRET_OK;
    }
    }
  }

#ifdef DEBUG_TRACE_EXECUTION
  printf("          ");
  for (Value *slot = vm.stack; slot < vm.stackTop; slot++)
  {
    printf("[ ");
    printValue(*slot);
    printf(" ]");
  }
  printf("\n");
  disassembleInstruction(vm.chunk,
                         (int)(vm.ip - vm.chunk->code));
#endif

#undef READ_BYTE
#undef READ_CONSTANT
}

InterpretResult interpret(Chunk *chunk)
{
  vm.chunk = chunk;
  vm.ip = vm.chunk->code;
  return run();
}

void push(Value value)
{
  *vm.stackTop = value;
  vm.stackTop++;
}

// "We don’t need to explicitly “remove” it from the array—moving stackTop down is enough to mark that slot as no longer in use."
Value pop()
{
  vm.stackTop--;
  return *vm.stackTop;
}