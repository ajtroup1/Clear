#pragma once

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
    inline Tokenizer(const std::string& src)
        : m_src(std::move(src)) {}

    std::vector<Token> tokenize() { // ****START REWRITING HERE @14:30
        std::vector<Token> tokens;
        std::string buf;
        for (int i = 0; i < str.length(); i++) {
            char c = str.at(i);
            if (std::isalpha(c)) {
                buf.push_back(c);
                i++;
                while (std::isalnum(str.at(i))) {
                    buf.push_back(str.at(i));
                    i++;
                }
                i--;

                if (buf == "exit") {
                    tokens.push_back({.type = TokenType::_exit});
                    buf.clear();
                    continue;
                } else {
                    std::cerr << "Invalid input" << std::endl;
                    exit(EXIT_FAILURE);
                }
            } else if (std::isspace(c)) {
                continue;
            } else if (std::isdigit(c)) {
                buf.push_back(c);
                while (i + 1 < str.length() && std::isdigit(str.at(i + 1))) {
                    buf.push_back(str.at(++i));
                }
                tokens.push_back({.type = TokenType::int_lit, .value = buf});
                buf.clear();
            }   else if(c == ';') {
                tokens.push_back({.type = TokenType::semi});
            } else {
                    std::cerr << "Invalid input" << std::endl;
                    exit(EXIT_FAILURE);
            }
        }

        return tokens;
    }

private:
    [[nodiscard]] std::optional<char> peak(int ahead = 1) const {
        if (m_index + ahead >= m_src.length()) {
            return {};
        } else {
            return m_src.at(m_index);
        }
    }

    char consume() {
        return m_src.at(m_index++);
    }

    const std::string m_src;
    int m_index;
};
