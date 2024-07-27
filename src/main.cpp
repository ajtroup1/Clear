#include <iostream>
#include <fstream>
#include <sstream>
#include <optional>
#include <vector>

#include "tokenization.hpp"
#include "parser.hpp"
#include "generation.hpp"

int main(int argc, char* argv[]) {
    if (argc != 2) {
        std::cerr << "Invalid call... Please use:" << std::endl;
        std::cerr << "(from build dir) './clear ../example.clr'" << std::endl;
        return EXIT_FAILURE;
    }

    std::string contents;
    {
        std::ifstream input(argv[1]);
        if (!input) {
            std::cerr << "Failed to open file: " << argv[1] << std::endl;
            return EXIT_FAILURE;
        }
        std::stringstream contents_stream;
        contents_stream << input.rdbuf();
        contents = contents_stream.str();
    }

    Tokenizer tokenizer(std::move(contents));
    std::vector<Token> tokens = tokenizer.tokenize();

    Parser parser(std::move(tokens));
    std::optional<NodeProg> prog = parser.parse_prog();

    if (!prog.has_value()) {
        std::cerr << "Invalid program" << std::endl;
        exit(EXIT_FAILURE);
    }
    
    Generator generator(prog.value());

    {
        std::ofstream file("out.asm");
        if (!file) {
            std::cerr << "Failed to open output file: out.asm" << std::endl;
            return EXIT_FAILURE;
        }
        file << generator.gen_prog();
    }

    system("nasm -felf64 out.asm");
    system("ld -o out out.o");

    return EXIT_SUCCESS;
}