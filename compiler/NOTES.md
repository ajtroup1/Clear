# Compiler Notes

## Contents
1. [What is bytecode?](#what-is-bytecode)
2. [Compilation process / flow](#compilation-process--flow)
3. [Random Notes](#random-notes)

## What is bytecode?
Bytecode is sort of an intermediate representation, lying between source and machine code (or binaries). Bytecode acts as a middle-man between interpreting source code and generating asm based on source code. Since machine code depends so heavily on system architecture and must be designed dynamically to run on all devices, bytecode gives us an "easy" way to make a compiler without worrying about: Linux vs Windows assembly, differences in CPUs or GPUs, or what 'x86' even means.

The good thing about building a bytecode compiler is that 2/3 of the process is similar to making an interpreter. Before you actually generate bytecode, you can sill build a parse tree the same way as you do in an interpreter (like `goclear`).

Bytecode can also be thought of as assembly for virtual machines. It abstracts low-level functionality via bytes (`0x00`, `0x01`, ...) to actual functionality that would be executed by asm. Take this bytecode example and see how similar it is to asm, just without worrying about system architecture or what assembly language you should be using

BYTECODE:

```asm
0x01 0x05        // LOAD_CONST 5
0x02 0x01        // STORE_VAR x
0x03 0x01        // LOAD_VAR x
0x01 0x03        // LOAD_CONST 3
0x04             // ADD
0x02 0x02        // STORE_VAR y
0x03 0x02        // LOAD_VAR y
0x05             // PRINT
```

ASSEMBLY:

```asm
MOV R0, 5       ; Move 5 into register R0
STORE R0, x     ; Store the value in R0 into variable x
LOAD R1, x      ; Load the value of x into register R1
MOV R2, 3       ; Move 3 into register R2
ADD R1, R2      ; Add the values in R1 and R2
STORE R1, y     ; Store the result in variable y
LOAD R0, y      ; Load the value of y into register R0
CALL PRINT      ; Call the print function
```

# Compilation process / flow

1. Generate the AST
    - 1a. Tokenize / lex the source code into a stream of tokens
    - 1b. Parse individual tokens into a structured parse tree
2. ...

## Random Notes
- There are `single-pass` and `double-pass` compilers, which either go through the IR once or twice depending on how advanced you want the compiler to be.
    - Example:
      - When compiling conditionals (if stmts), you use jump instructions to "skip" the generated consequence or alternative. Figuring out where to jump to is the difficult part. You can `back-patch` the location to jump to after compiling the consquence and alternative. Or you can leave them blank and pass through the IR a second time to fill those in.