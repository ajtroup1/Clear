#ifndef clear_lexer_hpp
#define clear_lexer_hpp

#include <sstream>
#include <string>
#include <vector>

enum class TokenType {
  // Single-character special tokens
  LEFT_PAREN,
  RIGHT_PAREN,
  LEFT_BRACE,
  RIGHT_BRACE,
  LEFT_BRACKET,
  RIGHT_BRACKET,
  SEMI,
  COMMA,
  DOT,
  SINGLE_QUOTE,
  DOUBLE_QUOTE,
  COLON,
  BACKSLASH,
  AMPERSAND,
  PIPE,

  // Single-character operators
  PLUS,
  MINUS,
  STAR,
  SLASH,
  PERCENT,
  CARET,
  QUESTION,
  EQUAL,
  LESS,
  GREATER,
  BANG,

  // Two-character operators
  PLUS_EQUAL,
  MINUS_EQUAL,
  STAR_EQUAL,
  SLASH_EQUAL,
  EQUAL_EQUAL,
  LESS_EQUAL,
  GREATER_EQUAL,
  BANG_EQUAL,

  // Logical operators
  AND,  // &&
  OR,   // ||

  // Keywords
  IDENT,

  LET,
  RETURN,
  CONST,
  IF,
  ELSE,
  BREAK,
  CONTINUE,

  // Data types
  NUMBER,
  STRING,
  BOOL,

  // Loops
  WHILE,
  FOR,

  // Special tokens
  _EOF,
  UNDEFINED
};

std::string tokenTypeToString(TokenType type);

class Token {
 public:
  Token(TokenType type, const std::string& literal, int line, int column)
      : type(type), literal(literal), line(line), column(column) {}

  TokenType getType() const { return type; }
  void setType(TokenType type) { this->type = type; }
  const std::string& getLiteral() const { return literal; }
  std::string stringify() const {
    std::stringstream ss;
    ss << "Token: " << tokenTypeToString(type) << " (" << literal
       << ") at [line: " << line << ", col: " << column << "]";
    return ss.str();
  }

 private:
  TokenType type;
  std::string literal;
  int line;
  int column;
};

class Lexer {
 public:
  Lexer(const std::string& src) : src(src), pos(0), line(1), column(1) {}

  std::vector<Token> tokenize();

 private:
  std::string src;
  size_t pos;
  int line;
  int column;

  char peek() const;
  char peekN(size_t n) const;
  char consume();
  void skipWhitespace();
  Token nextToken();
};

#endif