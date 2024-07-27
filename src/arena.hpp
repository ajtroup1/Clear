#pragma once

#include <cstddef> 
#include <cstdlib>

class ArenaAllocator {
public:
    // Constructor to initialize the allocator with a given size
    inline explicit ArenaAllocator(size_t bytes)
        : m_size(bytes), m_offset(nullptr)
    {
        m_buffer = static_cast<std::byte*>(malloc(m_size));
        m_offset = m_buffer;
    }

    // Template function to allocate memory for an object of type T
    template<typename T>
    inline T* alloc() {
        void* offset = m_offset;
        m_offset += sizeof(T);
        return offset;
    }

    // Delete copy constructor
    inline ArenaAllocator(const ArenaAllocator& other) = delete;

    // Delete copy assignment operator
    inline ArenaAllocator& operator=(const ArenaAllocator& other) = delete;

    // Destructor to free the allocated memory
    inline ~ArenaAllocator() {
        free(m_buffer);
    }

private:
    size_t m_size;          // Total size of the memory block
    std::byte* m_buffer;    // Pointer to the allocated memory block
    std::byte* m_offset;    // Current offset within the memory block
};
