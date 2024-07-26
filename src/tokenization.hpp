#pragma once

#include <iostream>
#include <string>
#include <optional>
#include <vector>

enum class TokenType {
    _exit,
    int_lit,
    semi
};

struct Token {
    TokenType type;
    std::optional<std::string> value;
};

class Tokenizer {
public:
    inline explicit Tokenizer(const std::string& src)
        : m_src(std::move(src)), m_index(0) {}

    std::vector<Token> tokenize() {
        std::vector<Token> tokens;
        std::string buf;
        while (peek().has_value()) {
            char current = peek().value();
            if (std::isalpha(current)) {
                buf.push_back(consume());
                while (peek().has_value() && std::isalnum(peek().value())) {
                    buf.push_back(consume());
                }
                if (buf == "exit") {
                    tokens.push_back({TokenType::_exit, {}});
                } else {
                    std::cerr << "Invalid keyword: " << buf << std::endl;
                    exit(EXIT_FAILURE);
                }
                buf.clear();
            } else if (std::isdigit(current)) {
                buf.push_back(consume());
                while (peek().has_value() && std::isdigit(peek().value())) {
                    buf.push_back(consume());
                }
                tokens.push_back({TokenType::int_lit, buf});
                buf.clear();
            } else if (current == ';') {
                consume();
                tokens.push_back({TokenType::semi, {}});
            } else if (std::isspace(current)) {
                consume();
            } else {
                std::cerr << "Unexpected character: " << current << std::endl;
                exit(EXIT_FAILURE);
            }
        }
        return tokens;
    }

private:
    [[nodiscard]] std::optional<char> peek(int ahead = 0) const {
        if (m_index + ahead >= m_src.length()) {
            return std::nullopt;
        } else {
            return m_src.at(m_index + ahead);
        }
    }

    char consume() {
        return m_src.at(m_index++);
    }

    std::string m_src;
    int m_index;
};
