# Clear
> "There are two kinds of programmers: those who have written compilers and those who haven't" - Terry A. Davis

## Grammar
$$
\begin{align}
[\text{prog}] &\to [\text{stmt}]^* \\
[\text{stmt}] &\to 
\begin{cases}
    \text{exit}(\text{[expr]}); \\
    \text{let}\space\text{ident} = \text{[expr]};
\end{cases} \\
[\text{expr}] &\to 
\begin{cases}
\text{int\_lit} \\
\text{ident}
\end{cases}
\end{align}
$$
