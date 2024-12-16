#ifndef clear_chunk_h
#define clear_chunk_h

#include "common.h"

// Enum for the type of bytecode instruction stored in a chunk
typedef enum {
  OP_RETURN,
} OpCode;

// Data struct for a chunk
// Uses dynamic array structure to hold a sequence, or "chunk", of bytecode
typedef struct {
  int count;
  int capacity;
  uint8_t* code;
} Chunk;

// Initialize a new chunk object to default values
void initChunk(Chunk* chunk);

// Append a byte to the end of a given chunk
void writeChunk(Chunk* chunk, uint8_t byte);

#endif