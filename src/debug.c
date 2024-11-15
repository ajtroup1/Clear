#include <stdio.h>

#include "debug.h"
#include "value.h"

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
  // Switch for the type of instruction since each are handled relatively uniquely
  switch (instruction)
  {
  case OP_CONSTANT:
    return constantInstruction("OP_CONSTANT", chunk, offset);
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

int constantInstruction(const char *name, Chunk *chunk,
                               int offset)
{
  uint8_t constant = chunk->code[offset + 1];
  printf("%-16s %4d '", name, constant);
  printValue(chunk->constants.values[constant]);
  printf("'\n");
  return offset + 2;
}