#include <iostream>
#include <fstream>
#include <sstream>
#include <optional>
#include <vector>

#include <./tokenization.hpp>

enum class TokenType {
    _exit,
    int_lit,
    semi
};

struct Token {
    TokenType type;
    std::optional<std::string> value;
};

std::vector<Token> tokenize(const std::string& str) { // receives string and returns it as a series of tokens
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

std::string tokens_to_asm(const std::vector<Token>& tokens) {
    std::stringstream output;
    output << "global _start\n_start:\n";

    for (int i = 0; i < tokens.size(); i++) {
        const Token& token = tokens.at(i);
        if (token.type == TokenType::_exit) {
            if (i+1 < tokens.size() && tokens.at(i+1).type == TokenType::int_lit) {
                if (i+2 < tokens.size() && tokens.at(i+2).type == TokenType::semi) {
                    output << "    mov rax, 60\n";
                    output << "    mov rdi, " << tokens.at(i+1).value.value() << "\n";
                    output << "    syscall\n";
            }
            }
        }
    }

    return output.str();
}

int main(int argc, char* argv[]) {
    if (argc != 2) { // expects a clear file to be ran
        std::cerr << "Invalid call... Please use:" << std::endl;
        std::cerr << "(from build dir) './clear ../example.clr'" << std::endl;
        return EXIT_FAILURE;
    }

    

    std::string contents; // stream the file contents into a string
    {
        std::stringstream contents_stream;
        std::fstream input(argv[1], std::ios::in);
        contents_stream << input.rdbuf();
        contents = contents_stream.str();
    }

    // std::cout << contents << std::endl; // test

    std::vector<Token> tokens = tokenize(contents);
    {
        std::fstream file("out.asm", std::ios::out);
        file << tokens_to_asm(tokens);
    }

    system("nasm -felf64 out.asm");
    system("ld -o out out.o");

    return EXIT_SUCCESS;
}