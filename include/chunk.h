#ifndef clear_chunk_h
#define clear_chunk_h

/*
  Create / manage chunks that act as bits of individual instructions
  Bytecode is the middle man between a tree walking interpreter and machine code generating compiler
  A program is simply a string of byte instructions executed sequentially
*/

#include "common.h"

typedef enum
{
  OP_RETURN,
} OpCode;

typedef struct
{
  int count;
  int capacity;
  uint8_t *code; // Pointer to a dynamically-allocated array of bytecode instructions
} Chunk;

void initChunk(Chunk *chunk);
void freeChunk(Chunk *chunk);
void writeChunk(Chunk *chunk, uint8_t byte);

#endif