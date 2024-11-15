#include <stdio.h>

#include "debug.h"

/*
  Disassemble a chunk of bytecode
*/
void disassembleChunk(Chunk *chunk, const char *name)
{
  printf("== %s ==\n", name);

  for (int offset = 0; offset < chunk->count;)
  {
    offset = disassembleInstruction(chunk, offset);
  }
}


/*
  Disassemble the individual instructions of a chunk
*/
int disassembleInstruction(Chunk *chunk, int offset)
{
  printf("%04d ", offset);

  uint8_t instruction = chunk->code[offset];
  switch (instruction)
  {
  case OP_RETURN:
    return simpleInstruction("OP_RETURN", offset);
  default:
    printf("Unknown opcode %d\n", instruction);
    return offset + 1;
  }
}

// Prints the name of a simple instruction, meaning an instruction that requires no additional inputs or data
int simpleInstruction(const char *name, int offset)
{
  printf("%s\n", name);
  return offset + 1;
}