cmake_minimum_required(VERSION 3.10)

# Set the project name
project(Clear)

# Set the C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED True)

# Add the executable
add_executable(clear
    src/main.cpp
    src/lexer.cpp
)

# Include directories
target_include_directories(clear PRIVATE src/include)