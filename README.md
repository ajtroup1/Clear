# Clear

---

## GoClear

Follows Tyler Laceby's parsing series with the Pikachu on the thumbnails

## JClear

Follows the classic Crafting Interpreters

The Java implementation and also my first Java project ever

### Grammar

$$

% \begin{align}

% expression     → literal

%                | unary

%                | binary

%                | grouping ;

% literal        → NUMBER | STRING | "true" | "false" | "nil" ;

% grouping       → "(" expression ")" ;

% unary          → ( "-" | "!" ) expression ;

% binary         → expression operator expression ;

% operator       → "==" | "!=" | "<" | "<=" | ">" | ">="

%                | "+"  | "-"  | "*" | "/" ;

% [\text{Expression}] &\to&

% \begin{cases}

% [\text{Literal}] \\

% [\text{Unary}] \\

% [\text{Binary}] \\

% [\text{Grouping}] \\

% \end{cases}

% \\

% [\text{Literal}] &\to&

% \begin{cases}

% NUMBER \\

% STRING \\

% \text{"true"} \\

% \text{"false"} \\

% \text{"nil"} \\

% \end{cases}

% \\

% [\text{Grouping}] &\to& \text{"(" [Expression] ")"}

% \\

% [\text{Unary}] &\to& \text{("-" | "!") [Expression]}

% \\

% [\text{Binary}] &\to& \text{[Expression] \, [Operator] \, [Expression]}

% \\

% [\text{Operator}] &\to&

% \begin{cases}

% \text{"+"} \\

% \text{"-"} \\

% \text{"*"} \\

% \text{"/"} \\

% \text{"=="} \\

% \text{"!="} \\

% \text{"<"} \\

% \text{"<="} \\

% \text{">"} \\

% \text{">="} \\

% \end{cases}

% \end{align}

$$

**Expression** →  
- Literal  
- Unary  
- Binary  
- Grouping  

**Literal** →  
- NUMBER  
- STRING  
- "true"  
- "false"  
- "nil"  

**Grouping** →  
- "(" Expression ")"  

**Unary** →  
- ("-" | "!") Expression  

**Binary** →  
- Expression Operator Expression  

**Operator** →  
- "+"  
- "-"  
- "*"  
- "/"  
- "=="  
- "!="  
- "<"  
- "<="  
- ">"  
- ">=" 