#include <iostream>
#include <fstream>
#include <sstream>
#include <optional>

enum class TokenType {
    _return,
    int_lit,
    semi
};

struct Token {
    TokenType type;
    std::optional<std::string> value;
};

int main(int argc, char* argv[]) {
    if (argc != 2) { // expects a clear file to be ran
        std::cerr << "Invalid call... Please use:" << std::endl;
        std::cerr << "(from build dir) './clear ../test.clr'" << std::endl;
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

    //START TOKENIZING @39:50

    return EXIT_SUCCESS;
}