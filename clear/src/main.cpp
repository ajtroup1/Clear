#include <iostream>
#include <fstream>
#include <sstream>
#include "lexer.hpp"
#include "parser.hpp"

int main(int argc, char* argv[]) {
    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " <source file>" << std::endl;
        return EXIT_FAILURE;
    }

    std::ifstream file(argv[1]);
    if (!file.is_open()) {
        std::cerr << "Error: Could not open file " << argv[1] << std::endl;
        return EXIT_FAILURE;
    }

    std::stringstream buffer;
    buffer << file.rdbuf();
    std::string src = buffer.str();

    Lexer lexer(src);
    std::vector<Token> tokens = lexer.tokenize();

    if (tokens.empty()) {
        std::cerr << "Error: No tokens generated from source file" << std::endl;
        return EXIT_FAILURE;
    }

    if (argc > 2 && std::string(argv[2]) == "--debug") {
        for (const Token& token : tokens) {
            std::cout << token.stringify() << std::endl;
        }
    }

    try {
        Parser parser(tokens);
        std::unique_ptr<Program> program = parser.parse();

        if (argc > 2 && std::string(argv[2]) == "--debug") {
            std::cout << program->stringify() << std::endl;
        }
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}