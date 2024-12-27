#include "lexer.hpp"

#include <cctype>
#include <iostream>

// Definition of the function
std::string tokenTypeToString(TokenType type) {
  switch (type) {
    case TokenType::LEFT_PAREN:
      return "LEFT_PAREN";
    case TokenType::RIGHT_PAREN:
      return "RIGHT_PAREN";
    case TokenType::LEFT_BRACE:
      return "LEFT_BRACE";
    case TokenType::RIGHT_BRACE:
      return "RIGHT_BRACE";
    case TokenType::LEFT_BRACKET:
      return "LEFT_BRACKET";
    case TokenType::RIGHT_BRACKET:
      return "RIGHT_BRACKET";
    case TokenType::SEMI:
      return "SEMI";
    case TokenType::COMMA:
      return "COMMA";
    case TokenType::DOT:
      return "DOT";
    case TokenType::SINGLE_QUOTE:
      return "SINGLE_QUOTE";
    case TokenType::DOUBLE_QUOTE:
      return "DOUBLE_QUOTE";
    case TokenType::COLON:
      return "COLON";
    case TokenType::BACKSLASH:
      return "BACKSLASH";
    case TokenType::PLUS:
      return "PLUS";
    case TokenType::MINUS:
      return "MINUS";
    case TokenType::STAR:
      return "STAR";
    case TokenType::SLASH:
      return "SLASH";
    case TokenType::PERCENT:
      return "PERCENT";
    case TokenType::CARET:
      return "CARET";
    case TokenType::QUESTION:
      return "QUESTION";
    case TokenType::EQUAL:
      return "EQUAL";
    case TokenType::LESS:
      return "LESS";
    case TokenType::GREATER:
      return "GREATER";
    case TokenType::BANG:
      return "BANG";
    case TokenType::PLUS_EQUAL:
      return "PLUS_EQUAL";
    case TokenType::MINUS_EQUAL:
      return "MINUS_EQUAL";
    case TokenType::STAR_EQUAL:
      return "STAR_EQUAL";
    case TokenType::SLASH_EQUAL:
      return "SLASH_EQUAL";
    case TokenType::EQUAL_EQUAL:
      return "EQUAL_EQUAL";
    case TokenType::LESS_EQUAL:
      return "LESS_EQUAL";
    case TokenType::GREATER_EQUAL:
      return "GREATER_EQUAL";
    case TokenType::BANG_EQUAL:
      return "BANG_EQUAL";
    case TokenType::AND:
      return "AND";
    case TokenType::OR:
      return "OR";
    case TokenType::NUMBER:
      return "NUMBER";
    case TokenType::IDENT:
      return "IDENT";
    case TokenType::LET:
      return "LET";
    case TokenType::RETURN:
      return "RETURN";
    case TokenType::CONST:
      return "CONST";
    case TokenType::IF:
      return "IF";
    case TokenType::ELSE:
      return "ELSE";
    case TokenType::STRING:
      return "STRING";
    case TokenType::BOOL:
      return "BOOL";
    case TokenType::WHILE:
      return "WHILE";
    case TokenType::FOR:
      return "FOR";
    case TokenType::BREAK:
      return "BREAK";
    case TokenType::CONTINUE:
      return "CONTINUE";
    case TokenType::_EOF:
      return "END_OF_FILE";
    default:
      return "UNDEFINED";
  }
}

// Declaration of the keywordLookup function
// Token keywordLookup(const std::string& literal, int line, int col);

char Lexer::peek() const {
  if (pos >= src.size()) {
    return '\0';
  }
  return src[pos];
}

char Lexer::peekN(size_t n) const {
  if (pos + n >= src.size()) {
    return '\0';
  }
  return src[pos + n];
}

char Lexer::consume() {
  if (pos >= src.size()) {
    return '\0';
  }
  char c = src[pos++];
  return c;
}

void Lexer::skipWhitespace() {
  while (peek() == ' ' || peek() == '\t' || peek() == '\n' || peek() == '\r') {
    char c = consume();
    if (c == '\n' || c == '\r') {
      line++;
      column = 1;
    } else if (c == '\t') {
      column += 4;
    } else if (c == ' ') {
      column++;
    }
  }
}

// Definition of the keywordLookup function
Token keywordLookup(const std::string& literal, int line, int col) {
  if (literal == "let") {
    return Token(TokenType::LET, literal, line, col);
  } else if (literal == "return") {
    return Token(TokenType::RETURN, literal, line, col);
  } else if (literal == "const") {
    return Token(TokenType::CONST, literal, line, col);
  } else if (literal == "if") {
    return Token(TokenType::IF, literal, line, col);
  } else if (literal == "else") {
    return Token(TokenType::ELSE, literal, line, col);
  } else if (literal == "string") {
    return Token(TokenType::STRING, literal, line, col);
  } else if (literal == "bool") {
    return Token(TokenType::BOOL, literal, line, col);
  } else if (literal == "while") {
    return Token(TokenType::WHILE, literal, line, col);
  } else if (literal == "for") {
    return Token(TokenType::FOR, literal, line, col);
  } else if (literal == "break") {
    return Token(TokenType::BREAK, literal, line, col);
  } else if (literal == "continue") {
    return Token(TokenType::CONTINUE, literal, line, col);
  } else {
    return Token(TokenType::IDENT, literal, line, col);
  }
}

Token Lexer::nextToken() {
  skipWhitespace();

  char c = peek();

  if (c == '\0') {
    return Token(TokenType::_EOF, "", line, column + 1);
  }

  if (isalpha(c)) {
    std::string literal;
    int startColumn = column;
    while (isalnum(peek()) || peek() == '_') {
      literal += consume();
      column++;
    }
    const Token token = keywordLookup(literal, line, startColumn);
    return token;
  }

  if (isdigit(c)) {
    std::string literal;
    bool isFloat = false;
    int startColumn = column;
    literal += consume();
    column++;

    while (isdigit(peek())) {
      literal += consume();
      column++;
    }

    if (peek() == '.' && isdigit(peekN(1))) {
      literal += consume();
      column++;
      while (isdigit(peek())) {
        literal += consume();
        column++;
      }
    }

    return Token(TokenType::NUMBER, literal, line, startColumn);
  }

  switch (c) {
    case '(':
      consume();
      return Token(TokenType::LEFT_PAREN, "(", line, column);
    case ')':
      consume();
      return Token(TokenType::RIGHT_PAREN, ")", line, column);
    case '{':
      consume();
      return Token(TokenType::LEFT_BRACE, "{", line, column);
    case '}':
      consume();
      return Token(TokenType::RIGHT_BRACE, "}", line, column);
    case '[':
      consume();
      return Token(TokenType::LEFT_BRACKET, "[", line, column);
    case ']':
      consume();
      return Token(TokenType::RIGHT_BRACKET, "]", line, column);
    case ';':
      consume();
      return Token(TokenType::SEMI, ";", line, column);
    case ',':
      consume();
      return Token(TokenType::COMMA, ",", line, column);
    case '.':
      consume();
      return Token(TokenType::DOT, ".", line, column);
    case '\'':
      consume();
      return Token(TokenType::SINGLE_QUOTE, "'", line, column);
    case '"':
      consume();
      return Token(TokenType::DOUBLE_QUOTE, "\"", line, column);
    case ':':
      consume();
      return Token(TokenType::COLON, ":", line, column);
    case '\\':
      consume();
      return Token(TokenType::BACKSLASH, "\\", line, column);
    case '+':
      if (peekN(1) == '=') {
        consume();
        consume();
        return Token(TokenType::PLUS_EQUAL, "+=", line, column);
      }
      consume();
      return Token(TokenType::PLUS, "+", line, column);
    case '-':
      if (peekN(1) == '=') {
        consume();
        consume();
        return Token(TokenType::MINUS_EQUAL, "-=", line, column);
      }
      consume();
      return Token(TokenType::MINUS, "-", line, column);
    case '*':
      if (peekN(1) == '=') {
        consume();
        consume();
        return Token(TokenType::STAR_EQUAL, "*=", line, column);
      }
      consume();
      return Token(TokenType::STAR, "*", line, column);
    case '/':
      if (peekN(1) == '=') {
        consume();
        consume();
        return Token(TokenType::SLASH_EQUAL, "/=", line, column);
      }
      consume();
      return Token(TokenType::SLASH, "/", line, column);
    case '%':
      consume();
      return Token(TokenType::PERCENT, "%", line, column);
    case '^':
      consume();
      return Token(TokenType::CARET, "^", line, column);
    case '?':
      consume();
      return Token(TokenType::QUESTION, "?", line, column);
    case '=':
      if (peekN(1) == '=') {
        consume();
        consume();
        return Token(TokenType::EQUAL_EQUAL, "==", line, column);
      }
      consume();
      return Token(TokenType::EQUAL, "=", line, column);
    case '<':
      if (peekN(1) == '=') {
        consume();
        consume();
        return Token(TokenType::LESS_EQUAL, "<=", line, column);
      }
      consume();
      return Token(TokenType::LESS, "<", line, column);
    case '>':
      if (peekN(1) == '=') {
        consume();
        consume();
        return Token(TokenType::GREATER_EQUAL, ">=", line, column);
      }
      consume();
      return Token(TokenType::GREATER, ">", line, column);
    default:
      consume();
      return Token(TokenType::UNDEFINED, std::string(1, c), line, column);
  }
}

std::vector<Token> Lexer::tokenize() {
  std::vector<Token> tokens;
  Token token = nextToken();
  while (token.getType() != TokenType::_EOF) {
    tokens.push_back(token);
    token = nextToken();
  }
  tokens.push_back(token);  // EOF token
  return tokens;
}