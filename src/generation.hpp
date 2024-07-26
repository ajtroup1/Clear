#pragma once

#include <sstream>

class Generator {
public:
    inline explicit Generator(NodeProg prog)
        : m_prog(std::move(prog))
    { 
    }

    [[nodiscard]] std::string gen_expr(const NodeExpr& expr) {
        struct ExprVisitor {
            Generator* gen;
            ExprVisitor(Generator* gen) : gen(gen) {}
            void operator()(const NodeExprIntLit& expr_int_lit) {
                gen->m_output << "mov rax, " << expr_int_lit.int_lit.value.value() << "\n";
            }
            void operator()(const NodeExprIdent& expr_ident) {

            }
        };

        ExprVisitor visitor(this);
        std::visit(visitor, expr.var);
    }

    [[nodiscard]] std::string gen_stmt(const NodeStmt& stmt) {
        struct StmtVisitor {
            Generator* gen;
            StmtVisitor(Generator* gen) : gen(gen) {}
            void operator()(const NodeStmtExit& stmt_exit) {
                gen->m_output << "    mov rax, 60\n";
                gen->m_output << "    mov rdi, 0\n";
                gen->m_output << "    syscall\n"; 
            }
            void operator()(const NodeStmtLet& stmt_let) {
                
            }
        };

        StmtVisitor visitor(this);
        std::visit(visitor, stmt.var);
    }

    [[nodiscard]] std::string gen_prog() {
        std::stringstream m_output;
        m_output << "global _start\n_start:\n";

        for (const NodeStmt& stmt : m_prog.stmts) {
            m_output << gen_stmt(stmt);
        }

        m_output << "    mov rax, 60\n";
        m_output << "    mov rdi, 0\n";
        m_output << "    syscall\n"; 
        return m_output.str();
    }

private:
    const NodeProg m_prog;
    std::stringstream m_output;
};