#ifndef clear_debug_h
#define clear_debug_h

/*
  Contains functionality for debugging chunks
  Visualizes chunks for a structured understanding of what the computer is working with
*/

#include "chunk.h"

void disassembleChunk(Chunk *chunk, const char *name);
int disassembleInstruction(Chunk *chunk, int offset);

#endif