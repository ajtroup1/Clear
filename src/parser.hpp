#pragma once

#include <vector>
#include <utility>
#include <iostream>
#include "tokenization.hpp"

struct NodeExpr {
    Token int_lit;
};

struct NodeExit {
    NodeExpr expr;
};

class Parser {
public:
    explicit Parser(std::vector<Token> tokens)
        : m_tokens(std::move(tokens)), m_index(0) {}

    std::optional<NodeExpr> parse_expr() {
        if (peek().has_value() && peek().value().type == TokenType::int_lit) {
            return NodeExpr{.int_lit = consume()};
        } else {
            return {};
        }
    }

    std::optional<NodeExit> parse() {
        std::optional<NodeExit> exit_node;
        while (peek().has_value()) {
            if (peek().value().type == TokenType::_exit) {
                consume();
                if (auto node_expr = parse_expr()) {
                    if (peek().has_value() && peek().value().type == TokenType::semi) {
                        consume();
                        exit_node = NodeExit{.expr = node_expr.value()};
                    } else {
                        std::cerr << "Expected semicolon after expression" << std::endl;
                        exit(EXIT_FAILURE);
                    }
                } else {
                    std::cerr << "Invalid expression" << std::endl;
                    exit(EXIT_FAILURE);
                }
            } else {
                std::cerr << "Unexpected token" << std::endl;
                exit(EXIT_FAILURE);
            }
        }
        m_index = 0;
        return exit_node;
    }

private:
    [[nodiscard]] inline std::optional<Token> peek(int ahead = 0) const {
        if (m_index + ahead >= m_tokens.size()) {
            return std::nullopt;
        } else {
            return m_tokens.at(m_index + ahead);
        }
    }

    inline Token consume() {
        return m_tokens.at(m_index++);
    }

    const std::vector<Token> m_tokens;
    size_t m_index;
};
