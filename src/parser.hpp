#pragma once

#include <variant>
#include <vector>
#include <utility>
#include <iostream>
#include "tokenization.hpp"

struct NodeExprIntLit {
    Token int_lit;
};

struct NodeExprIdent {
    Token ident;
};

struct NodeExpr {
    std::variant<NodeExprIntLit, NodeExprIdent> var;
};

struct NodeStmtExit {
    NodeExpr expr;
};

struct NodeStmtLet {
    Token ident;
    NodeExpr expr;
};

struct NodeStmt {
    std::variant<NodeStmtExit, NodeStmtLet> var;
};

struct NodeProg {
    std::vector<NodeStmt> stmts;
};

class Parser {
public:
    explicit Parser(std::vector<Token> tokens)
        : m_tokens(std::move(tokens)), m_index(0) {}

    std::optional<NodeExpr> parse_expr() {
        if (peek().has_value() && peek().value().type == TokenType::int_lit) {
            return NodeExpr{.var = NodeExprIntLit {.int_lit = consume()} };
        } else if (peek().has_value() && peek().value().type == TokenType::ident) {
            return NodeExpr {.var = NodeExprIdent {.ident = consume()} };
        }else {
            return {};
        }
    }

    std::optional<NodeStmt> parse_stmt() {
    if (!peek().has_value()) {
        return {};
    }
    
    // Handle "exit" statements
    if (peek().value().type == TokenType::_exit && peek(1).has_value() && peek(1).value().type == TokenType::open_paren) {
        consume(); // consume 'exit'
        consume(); // consume '('
        
        NodeStmtExit stmt_exit;
        if (auto node_expr = parse_expr()) {
            stmt_exit = {.expr = node_expr.value()};
            
            if (peek().has_value() && peek().value().type == TokenType::close_paren) {
                consume(); // consume ')'
            } else {
                std::cerr << "Error: Expected closing parenthesis after expression." << std::endl;
                return {};
            }
            
            if (peek().has_value() && peek().value().type == TokenType::semi) {
                consume(); // consume ';'
                return NodeStmt{.var = stmt_exit};
            } else {
                std::cerr << "Error: Expected semicolon after 'exit' statement." << std::endl;
                return {};
            }
        } else {
            std::cerr << "Error: Unable to parse expression in 'exit' statement." << std::endl;
            return {};
        }
    }
    
    // Handle "let" statements
    if (peek().has_value() && peek().value().type == TokenType::let && peek(1).has_value() && peek(1).value().type == TokenType::ident && peek(2).has_value() && peek(2).value().type == TokenType::eq) {
        consume(); // consume 'let'
        auto stmt_let = NodeStmtLet {.ident = consume()}; // consume identifier
        
        consume(); // consume '='
        
        if (auto expr = parse_expr()) {
            stmt_let.expr = expr.value();
            
            if (peek().has_value() && peek().value().type == TokenType::semi) {
                consume(); // consume ';'
                return NodeStmt {.var = stmt_let};
            } else {
                std::cerr << "Error: Expected semicolon after 'let' statement." << std::endl;
                return {};
            }
        } else {
            std::cerr << "Error: Unable to parse expression in 'let' statement." << std::endl;
            return {};
        }
    }
    
    // If we reach here, it's an unknown or invalid statement
    std::cerr << "Error: Invalid statement encountered." << std::endl;
    return {};
}


    std::optional<NodeProg> parse_prog() {
        NodeProg prog;
        while (peek().has_value()) {
            if (auto stmt = parse_stmt()) {
                prog.stmts.push_back(stmt.value());
            } else {
                std::cerr << "Invalid statement" << std::endl;
                exit(EXIT_FAILURE);
            }
        }
        return prog;
    }

private:
    [[nodiscard]] inline std::optional<Token> peek(int offset = 0) const {
        if (m_index + offset >= m_tokens.size()) {
            return std::nullopt;
        } else {
            return m_tokens.at(m_index + offset);
        }
    }

    inline Token consume() {
        return m_tokens.at(m_index++);
    }

    const std::vector<Token> m_tokens;
    size_t m_index;
};
