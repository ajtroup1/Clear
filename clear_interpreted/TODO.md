# TODO

### High-Level
  1. Implement the "Clear" aspect of the language
      - The interpreter / compiler will be completely transparent about what it's doing
      - It will be a "talking interpreter"
      - Append strings to a log file that 'trace' what the interpreter is doing at any relevant time and dump that log in the example folder next to the AST, src code, etc.
  2. Builtins
  3. Error handling:
  4. Object oriented to some extent
  5. Built-in testing suite:
      - Automated testing framework / prebuilt
  6. Implement a lexer, parser, and evaluator in Clear
      - Self host the language
      - More as an example of Clear's features, not for efficiency 

### General Additions
  1. Add a DateTime type to aid the time module
  2. Warnings:
      - Variable `x` is unused
      - ...
  3. Add logical operators `&&` `||`

### Quick Fixes
  1. if statement condition booleans
    - if (myBoolean)
    - if (!myBoolean)
  2. Returns nil dereference: "x * = 2;" because of the space