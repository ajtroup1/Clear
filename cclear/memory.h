#ifndef clear_memory_h
#define clear_memory_h

#include "common.h"

// Increases the capacity property of a chunk based on a given old capacity
// If old capacity is 0, make new capacity 8
// If the old capacity is >= 8, double it
#define GROW_CAPACITY(capacity) ((capacity) < 8 ? 8 : (capacity) * 2)

#define GROW_ARRAY(type, pointer, oldCount, newCount)   \
  (type*)reallocate(pointer, sizeof(type) * (oldCount), \
                    sizeof(type) * (newCount))

void* reallocate(void* pointer, size_t oldSize, size_t newSize);

#endif