# GoClear

This language implementation follows Thorston Ball's classic book "Writing an Interpreter in Go", but will extend it (hopefully) tremendously

So, this interpreter implementation is written in Go

## src

#### Structure
The `src/` folder is structured according to the packages that correspond to different stages of interpreting source code.
- `lexer` & `token`
  - Handle the tokenizing stage, where source code is converted to tokens.
- `parser`
  - Handles parsing tokens into an AST (Abstract Syntax Tree) to be evaluated

## examples

In the `/examples` folder there are examples of Clear code chunks in folders

In each folder, there is:
- A source file containing Clear code (ex: `script.clr`)
- A JSON file representing the AST of that source code (ex. `script.clr.ast.json`)
- *Optionally* a `.md` file containing information about the source file or the features it implements
  - This information may be contained within comments in the source (`.clr`) file if it is brief

## Testing

If tests are necessary for any package, it will contain a *example*_test.go file within the corresponding folder
  ```
  parser/
    parser.go
    parser_test.go
  ```