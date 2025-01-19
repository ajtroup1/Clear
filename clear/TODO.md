# TODO

### High-Level
  1. Implement the "Clear" aspect of the language
      - The interpreter / compiler will be completely transparent about what it's doing
      - It will be a "talking interpreter"
      - Append strings to a log file that 'trace' what the interpreter is doing at any relevant time and dump that log in the example folder next to the AST, src code, etc.
  2. Builtins
      - Array
      - Math
      - Strings
  3. Error handling:
      - Complete rework of error handling system (look at how Go src code handles errors)
      - Outsource error handling to another package
      - Line and col numbers
      - Errors AND warnings
      - Colored and well formatted messaged
      - Store line content to be able to point to the correct area
  4. Built-in testing suite:
      - Automated testing framework / prebuilt

### General Additions
  2. For / while Statements

### Quick Fixes
  1. Nil pointer dereference upon no return value specified
    - Default return 0 or something?