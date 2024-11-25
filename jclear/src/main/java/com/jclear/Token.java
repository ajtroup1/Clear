package com.jclear;

// Overall data class for any token in JClear
class Token {
  final TokenType type;
  final String lexeme; // Or literal value
  final Object literal; // ?
  final int line; // Line number where the token was encountered

  Token(TokenType type, String lexeme, Object literal, int line) {
    this.type = type;
    this.lexeme = lexeme;
    this.literal = literal;
    this.line = line;
  }

  @Override
  public String toString() {
    if (literal != null && !literal.toString().isEmpty()) {
      return "[TYPE: " + type + "] '" + lexeme + "' :: [LITERAL: '" + literal + "']";
    } else {
      return "[TYPE: " + type + "] '" + lexeme + "'";
    }
  }
}