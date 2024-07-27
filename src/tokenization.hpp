#pragma once

#include <iostream>
#include <string>
#include <optional>
#include <vector>

enum class TokenType {
    _exit,
    int_lit,
    semi,
    open_paren,
    close_paren,
    ident,
    let,
    eq
};

struct Token {
    TokenType type;
    std::optional<std::string> value;
};

class Tokenizer {
public:
    inline explicit Tokenizer(const std::string& src)
        : m_src(std::move(src)), m_index(0) {}

    std::vector<Token> tokenize() { // convert source code into a string of tokens
        std::vector<Token> tokens;
        std::string buf; // tracks multi-char values
        while (peek().has_value()) {
            char current = peek().value();
            if (std::isalpha(current)) { // check for multi-char keywords / identifiers // must start as alpha and continue as alphanumeric
                buf.push_back(consume());
                while (peek().has_value() && std::isalnum(peek().value())) { // while the identifier string continues
                    buf.push_back(consume());
                }
                // is the string a keyword?
                if (buf == "exit") {
                    tokens.push_back({TokenType::_exit, {}});
                } else if (buf == "let") {
                    tokens.push_back({.type = TokenType::let});
                    buf.clear();
                    continue;
                } 
                // the string isn't a keyword, it's an identifier
                else {
                    tokens.push_back({.type = TokenType::ident, .value = buf});
                    buf.clear();
                    continue;
                }
                buf.clear();
            } else if (std::isdigit(current)) { // check for int lit
                buf.push_back(consume());
                while (peek().has_value() && std::isdigit(peek().value())) {
                    buf.push_back(consume());
                }
                tokens.push_back({TokenType::int_lit, buf});
                buf.clear();
            } else if (current == '=') {
                consume();
                tokens.push_back({TokenType::eq, {}});
                continue;
            }else if (current == ';') {
                consume();
                tokens.push_back({TokenType::semi, {}});
                continue;
            } else if (current == '(') {
                consume();
                tokens.push_back({TokenType::open_paren, {}});
                continue;
            } else if (current == ')') {
                consume();
                tokens.push_back({TokenType::close_paren, {}});
                continue;
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
    [[nodiscard]] inline std::optional<char> peek(int offset = 0) const {
        if (m_index + offset >= m_src.length()) {
            return std::nullopt;
        } else {
            return m_src.at(m_index + offset);
        }
    }

    inline char consume() {
        return m_src.at(m_index++);
    }

    std::string m_src;
    size_t m_index;
};
