#include <stdlib.h>

#include "memory.h"

/*
  Dynamically allocates memory based on both a new and old size
  Receives:
    A pointer to the array that's being managed
    The current size of that array
    The new desired size of the array

  Memory allocation will occur differently based on the new and old size:
    ___________________________________________________________________
    | oldSize    | newSize              | Operation                   |
    | 0	         | Non‑zero             | Allocate new block.         |
    | Non‑zero	 | 0                    | Free allocation.            |
    | Non‑zero	 | Smaller than oldSize	| Shrink existing allocation. |
    | Non‑zero	 | Larger than oldSize	| Grow existing allocation.   |
    |________________________________________________________________ |
*/
void *reallocate(void *pointer, size_t oldSize, size_t newSize)
{
  if (newSize == 0)
  {
    free(pointer);
    return NULL;
  }

  void *result = realloc(pointer, newSize);
  if (result == NULL)
    exit(1);
  return result;
}