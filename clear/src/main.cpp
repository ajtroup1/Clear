#include <cstdlib>
#include <fstream>
#include <iostream>

#include "lexer.hpp"

int main(int argc, char** argv) {
  bool debug = true;
  if (argc < 2) {
    std::cerr << "Usage: 'clear <input file>'" << std::endl;
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

  if (debug) {
    for (const Token& token : tokens) {
      std::cout << token.stringify() << std::endl;
    }
  }

  return EXIT_SUCCESS;
}