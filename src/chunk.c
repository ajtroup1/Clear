#include <stdlib.h>

#include "chunk.h"
#include "memory.h"

void initChunk(Chunk *chunk)
{
  chunk->count = 0;
  chunk->capacity = 0;
  chunk->code = NULL;
  initValueArray(&chunk->constants);
}

void freeChunk(Chunk *chunk)
{
  FREE_ARRAY(uint8_t, chunk->code, chunk->capacity);
  freeValueArray(&chunk->constants);
  initChunk(chunk);
}

void writeChunk(Chunk *chunk, uint8_t byte)
{
  // Is the capacity reached by count yet?
    // Or defaulted at 0?
  if (chunk->capacity < chunk->count + 1)
  {
    // Yes, grow the capacity using a macro
    int oldCapacity = chunk->capacity;
    chunk->capacity = GROW_CAPACITY(oldCapacity);
    chunk->code = GROW_ARRAY(uint8_t, chunk->code,
                             oldCapacity, chunk->capacity);
  }

  // Asign the byte to the array
  chunk->code[chunk->count] = byte;
  chunk->count++;
}

int addConstant(Chunk *chunk, Value value)
{
  writeValueArray(&chunk->constants, value);
  return chunk->constants.count - 1;
}