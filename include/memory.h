#ifndef clear_memory_h
#define clear_memory_h

/*
  Defines basic memory management systems
*/

#include "common.h"

// Macros

/*
  Expands the capacity attribute of a dynamic array
    If the capacity is below 8, the capacity is set to 8
      By default, the capacity is 0, so it will always set to 8 when used
    If the capacity is already 8 or higher, double it

*/
#define GROW_CAPACITY(capacity) \
  ((capacity) < 8 ? 8 : (capacity) * 2)

// Pretty way to call the reallocate() function
#define GROW_ARRAY(type, pointer, oldCount, newCount)    \
  (type *)reallocate(pointer, sizeof(type) * (oldCount), \
                     sizeof(type) * (newCount))

// Calls reallocate() with a new size of 0, which frees it in its implementation
#define FREE_ARRAY(type, pointer, oldCount) \
  reallocate(pointer, sizeof(type) * (oldCount), 0)


void *reallocate(void *pointer, size_t oldSize, size_t newSize);

#endif