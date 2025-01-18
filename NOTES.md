# Notes

## Interpreter
1. Read in a file as a single string.
    - Example File:
      ```
      let x = 7;
      return x;
      ```
    - Example output for step 1:
      - `“let x = 7;\nreturn x;”`
2. Establish the tokens for your language’s grammar.
    - aka “what words, symbols, keywords, etc will people be able to type in your language’s source code files?”
    - A `token` is a data structure containing the type of token and a literal value:
      - The type of the token is what the token is:
      - PLUS
      - ASTERISK
      - IDENTIFIER
        - Variable or function identifier, for example
        - x, myFunction, myObject
        - But, IDENTIFIER can be a string that identifies anything, even a class
      - RETURN
      - LET
      - CONST
      - RIGHT_PARENTHESIS
      - UNDERSCORE
      - NUMBER
      - STRING
    - The literal value of the token is the raw text that was read from the source code:
      - “+”
      - “*”
      - “x”
      - “return”
      - “let”
      - “const”
      - “(“
      - “_”
      - “7”
      - “myString”
    - The goal of a `Token` is to represent raw source code as a simple data structure of information.
    - Example input:
      - `“let x = 7;”`
    - Example output:
      1. The `let` would be ‘tokenized’ or ‘lexed’ and returned as an object:
        - ```
          Token { TokenType: LET, Literal: “let” }
          ```
      2. The “x” would then be tokenized and returned as an object:
        - ```
          Token { TokenType: IDENTIFIER, Literal: “x” }
          ```
      3. The entire example would return a stream of tokens:
          - ```
            Tokens token.Token[] = 
            [ 
              Token{ TokenType: LET, Literal: “let” }, 
              Token{ TokenType: IDENTIFIER, Literal: “x” },
              Token{ TokenType: EQUALS, Literal: “+” },
              Token{ TokenType: INT, Literal: “7” },
              Token{ TokenType: SEMICOLON, Literal: “;” }
            ]
            ```
    - Tokens can be divided into groups based on what raw text they represent:
      - Variable identifiers and data type indicators
        - x
        - foobar
        - myVariable
        - int
        - string
        - bool
      - Operators
        -  `+`
        -  `-`
        -  `/`
        -  `*`
        -  `!`
        -  `++`
        -  `--`
      - Delimiters
        -  ,
        -  ;
        -  :
        -  (
        -  {
      - Keywords:
        -  let
        -  return
        -  function / func / fn
        -  if / else
        -  true / false
3. Run the source code string through a Lexer / Tokenizer to convert the source code into a stream of tokens.
    - Example input:
      - ```
        let x = 7;
        ```
    - Example output:
      - ```
        Tokens token.Token[] = [
          Token{ TokenType: LET, Literal: "let" },
          Token{ TokenType: IDENTIFIER, Literal: "x" },
          Token{ TokenType: EQUALS, Literal: "=" },
          Token{ TokenType: INT, Literal: "7" }
          [let], [x], [=], [7]
        ]
        ```
4. Now, you should have a stream of tokens (an array, list, slice, vector, whatever). Now, the `parser` will inspect these tokens one after the other, making decisions on the fly to structure the tokens in a hierarchical structure.
  - The structure it creates is actually just a tree. An Abstract Syntax Tree (or AST) to be specific. Since this is a basic tree structure, it can be extremely helpful to think of the program's parse tree as a [JSON](https://www.w3schools.com/js/js_json_intro.asp) structure as well.
    - <img src="https://ruslanspivak.com/lsbasi-part7/lsbasi_part7_astprecedence_01.png"/>
    - Example code:
      - ```
        let x = 7;
        return x;
        ```
    - Example output:
      - ```json
        Program: { <-- the root Program node
          "NoStatements": false, <-- flag for whether there are no statements found in the program
          "Statements": [
            LetStatement: {
              "Token": {
                "Type": "LET",
                "Literal": "let"
              },
              "Name": {
                "Token": {
                  "Type": "IDENT",
                  "Literal": "x"
                },
                "Value": "x"
              },
              "Value": {
                "Token": {
                  "Type": "INT",
                  "Literal": "7"
                },
                "Value": 7
              }
            }, <-- END LET STATEMENT
            ReturnStatement: {
              "Token": {
                "Type": "RETURN",
                "Literal": "return"
              },
              "ReturnValue": {
                "Token": {
                  "Type": "IDENT",
                  "Literal": "x"
                },
                "Value": "x"
              }
            }
          ],
          "Modules": null <-- Imported module information is stored at the highest level
        }
        ```
  - ASTs also have to be mindful of operator precedence, which can be a challenge when first learning how to parse source code.
    - Operator precedence is handled in parsing techniques such as [Pratt Parsing](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/) and [Recursive Descent](https://www.cs.rochester.edu/users/faculty/nelson/courses/csc_173/grammars/parsing.html).
  - Parsing is maybe the toughest to grasp and most extensive/complex step in interpreting source code. Parsing is what truly defines the syntax, structure, and immediate flow of your language from one token to the next.
    - Speaking from experience, you will spend around 50% -> 75% of your time coding *(and debugging AND testing)* the parser.
  - View the [README](./README.md) to see resources for learning parsing, as parsing takes up most of the total content there.
5. Now that the source code has been parsed into an AST, it's time to "walk" that AST and [evalute](https://mariusbancila.ro/blog/2009/02/06/evaluate-expressions-%E2%80%93-part-4-evaluate-the-abstract-syntax-tree/) its nodes on the fly.
    - A [Tree Walking Interpreter](https://lutzhamel.github.io/CSC402/notes/csc402-ln006a.pdf) is a common beginner approach to interpreting source code. It is technically the most inefficient way to evalute source code, but it is easy to understand. However, the difference in complexity and difficulty between tree walking interpreters and [JIT compilers](https://en.wikipedia.org/wiki/Just-in-time_compilation) is pretty large in my opinion.
