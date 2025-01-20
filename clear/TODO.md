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
  4. Object oriented to some extent
  5. Built-in testing suite:
      - Automated testing framework / prebuilt
  6. Implement a lexer, parser, and evaluator in Clear
      - Self host the language
      - More as an example of Clear's features, not for efficiency 

### General Additions
  1. For / while Statements
  2. Add a DateTime type to aid the time module
  3. Implement compound operators like `+=` `-=` `*=` ...

### Quick Fixes
  1. if statement condition booleans
    - if (myBoolean)
    - if (!myBoolean)
  2. Empty var declaration