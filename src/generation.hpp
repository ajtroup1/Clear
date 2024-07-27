#pragma once

#include <sstream>
#include <unordered_map>

class Generator {
public:
    // Constructor initializes the generator with the given NodeProg.
    // NodeProg contains the abstract syntax tree (AST) of the program.
    inline explicit Generator(NodeProg prog)
        : m_prog(std::move(prog))
    { 
    }

    // Generates assembly code for expression.
    void gen_expr(const NodeExpr& expr) {
        struct ExprVisitor {
            Generator* gen;
            
            // Handles integer literal expressions.
            void operator()(const NodeExprIntLit& expr_int_lit) {
                // Generate assembly to move integer literal value into the rax register.
                gen->m_output << "    mov rax, " << expr_int_lit.int_lit.value.value() << "\n";
                gen->push("rax"); // Push the rax register value onto the stack.
            }
            
            // Handles identifier expressions (variables).
            void operator()(const NodeExprIdent& expr_ident) {
                // Check if the identifier is declared. If not, report an error.
                if (gen->m_vars.find(expr_ident.ident.value.value()) == gen->m_vars.end()) {
                    std::cerr << "Identifier `" << expr_ident.ident.value.value() << "` is undeclared" << std::endl;
                    exit(EXIT_FAILURE);
                }
                const auto& var = gen->m_vars.at(expr_ident.ident.value.value());
                // Calculate the stack offset for the variable and generate assembly code to access it.
                std::stringstream offset;
                offset << "QWORD [rsp + " << (gen->m_stack_size - var.stack_loc - 1) * 8 << "]\n";
                gen->push(offset.str()); // Push the value at the calculated offset onto the stack.
            }
        };

        ExprVisitor visitor { .gen = this };
        // Visit the expression based on its type (int_lit or ident) and generate corresponding assembly code.
        std::visit(visitor, expr.var);
    }

    // Generates assembly code for the given statement.
    void gen_stmt(const NodeStmt& stmt) {
        struct StmtVisitor {
            Generator* gen;
            
            // Handles exit statements. 
            void operator()(const NodeStmtExit& stmt_exit) const {
                // Generate code to evaluate the expression.
                gen->gen_expr(stmt_exit.expr);
                // Prepare to exit the program with the return value.
                gen->m_output << "    mov rax, 60\n"; // System call number for exit.
                gen->pop("rdi"); // Pop the return value into the rdi register.
                gen->m_output << "    syscall\n"; // Make the system call to exit the program.
            }
            
            // Handles let statements (variable declarations).
            void operator()(const NodeStmtLet& stmt_let) {
                // Check if the identifier is already used. If so, report an error.
                if (gen->m_vars.find(stmt_let.ident.value.value()) != gen->m_vars.end()) {
                    std::cerr << "Identifier `" << stmt_let.ident.value.value() << "` already used" << std::endl;
                    exit(EXIT_FAILURE);
                }
                // Insert the new variable into the map with its stack location.
                gen->m_vars.insert({ stmt_let.ident.value.value(), Var { .stack_loc = gen->m_stack_size }});
                // Generate code to evaluate the expression and store the result.
                gen->gen_expr(stmt_let.expr);
            }
        };

        StmtVisitor visitor { .gen = this };
        // Visit the statement based on its type (exit or let) and generate corresponding assembly code.
        std::visit(visitor, stmt.var);
    }

    // Generates the complete assembly code for the entire program.
    [[nodiscard]] std::string gen_prog()
    {
        m_output << "global _start\n_start:\n"; // Define the entry point of the program.

        // Generate assembly code for each statement in the program.
        for (const NodeStmt& stmt : m_prog.stmts) {
            gen_stmt(stmt);
        }

        // Append the default exit code to ensure the program terminates correctly.
        m_output << "    mov rax, 60\n"; // System call number for exit.
        m_output << "    mov rdi, 0\n"; // Exit code 0.
        m_output << "    syscall\n"; // Make the system call to exit the program.
        return m_output.str();
    }

private:
    // Generates assembly code to push a register value onto the stack.
    void push(const std::string& reg) {
        m_output << "    push " << reg << "\n";
        m_stack_size++; // Increment the stack size.
    }

    // Generates assembly code to pop a register value from the stack.
    void pop(const std::string& reg) {
        m_output << "    pop " << reg << "\n";
        m_stack_size--; // Decrement the stack size.
    }

    // Represents a variable and its location on the stack.
    struct Var {
        size_t stack_loc; // Stack location of the variable.
    };

    const NodeProg m_prog; // The abstract syntax tree of the program.
    std::stringstream m_output; // Stream to accumulate the generated assembly code.
    size_t m_stack_size = 0; // Current stack size.
    std::unordered_map<std::string, Var> m_vars {}; // Map of variable names to their stack locations.
};
