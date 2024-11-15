#ifndef clear_chunk_h
#define clear_chunk_h

/*
  Create / manage chunks that act as bits of individual instructions
  Bytecode is the middle man between a tree walking interpreter and machine code generating compiler
  A program is simply a string of byte instructions executed sequentially
*/

#include "common.h"
#include "value.h"

typedef enum
{
  OP_CONSTANT,
  OP_RETURN,
} OpCode;

typedef struct
{
  int count;
  int capacity;
  uint8_t *code; // Pointer to a dynamically-allocated array of bytecode instructions
  ValueArray constants; // Stuct containing a dynamic array acting as a pool of constants for the program
  int *lines;
} Chunk;

// Basic dynamic array functionality
void initChunk(Chunk *chunk);
void freeChunk(Chunk *chunk);
void writeChunk(Chunk *chunk, uint8_t byte, int line);

// Writes a value to a given chunk's constants array
int addConstant(Chunk *chunk, Value value);

#endif