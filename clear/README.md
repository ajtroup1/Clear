# GoClear

Follows and extends *Thorston Ball*'s [Writing an Interpreter in Go](https://interpreterbook.com/)

### Overview


### Usage
Clear is ran via a Makefile in the `Clear\clear` directory

**Makefile**:
- Automates building the project, making it efficient to run various commands on the src code
- **Make commands**:
  - `make build`
    - Simply builds the `clear` executable into the `bin` directory to be ran at a later time
  - `make repl`
    - Runs Clear in the `repl` mode, which is a real-time, interactive code executor as opposed to a predefined file. Excellent for quick testing and development
      - More about `repl` in Python context [here](https://codewith.mu/en/tutorials/1.2/repl)
  - `make run` --> `make run ARGS="../examples/00/00.clr -d"`
    - *This make run command actually works!*
    - The run command executes a predefined `.clr` file 
    - The command expects:
      - *Optionally* a debug flag `-d` OR `-debug`
        - The debug flag will enable:
          - Printing a JSON file containing the source code's Abstract Syntax Tree
          - A "talking" log file detailing every step the interpreter took to interpret the source code
      - A path pointing to the `.clr` script to be executed (**Required!**)
  - `make test`
    - Runs all Go test files in the src
    - All this does is call `go test ./...` with the verbose flag
  - `make fmt`
    - Formats every file in the src according to the `go fmt` standard
    - All this does is call `go fmt ./...`


### A Talking Interpreter
The Clear interpreter is designed to "talk" to you as it

The talking is actually just detailed explainations of every step in each process in the interpeter (lexing, parsing, evaluating)

Make sure to execute a file in [script / "make run"](#usage) mode with the `-d` (or `-debug`) flag to have the interpreter generate a log file with this content