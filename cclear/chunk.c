#include <stdlib.h>

#include "chunk.h"
#include "memory.h"

// Instantiate the array with 0 count and capacity
void initChunk(Chunk* chunk) {
  chunk->count = 0;
  chunk->capacity = 0;
  chunk->code = NULL;
}

void writeChunk(Chunk* chunk, uint8_t byte) {
  // Does the array have the capacity for a new byte?
  if (chunk->capacity < chunk->count + 1) {
    int oldCapacity = chunk->capacity;
    // Grow the capacity of the array based on its current capacity
    chunk->capacity = GROW_CAPACITY(oldCapacity);
    chunk->code =
        GROW_ARRAY(uint8_t, chunk->code, oldCapacity, chunk->capacity);
  }

  // Append the byte to the end of the code and incriment count
  chunk->code[chunk->count] = byte;
  chunk->count++;
}