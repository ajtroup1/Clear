package com.jclear;

import static com.jclear.TokenType.*;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

// Scanner class containing all state and functionality
public class Scanner {
  private final String source; // Entire source code squished into a single string
  private final List<Token> tokens = new ArrayList<>(); // List of tokens to be populated as the Scanner lexes the
                                                        // source code
  private int start = 0;
  private int current = 0;
  private int line = 1; // Keeps track of line numbers for mainly error handling

  private static final Map<String, TokenType> keywords; // Maps string representations of keywords to their
                                                        // corresponding token

  static {
    keywords = new HashMap<>();
    keywords.put("and", AND);
    keywords.put("class", CLASS);
    keywords.put("else", ELSE);
    keywords.put("false", FALSE);
    keywords.put("for", FOR);
    keywords.put("fun", FUN);
    keywords.put("if", IF);
    keywords.put("nil", NIL);
    keywords.put("or", OR);
    keywords.put("print", PRINT);
    keywords.put("return", RETURN);
    keywords.put("super", SUPER);
    keywords.put("this", THIS);
    keywords.put("true", TRUE);
    keywords.put("var", VAR);
    keywords.put("while", WHILE);
  }

  Scanner(String source) {
    this.source = source;
  }

  List<Token> scanTokens() {
    while (!isAtEnd()) {
      // We are at the beginning of the next lexeme.
      start = current;
      scanToken();
    }

    tokens.add(new Token(EOF, "", null, line));
    return tokens;
  }

  // Returns true if the current Scanner position is beyond or at the EOF
  private boolean isAtEnd() {
    return current >= source.length();
  }

  // Core Scanner function. Analyzes the current char and makes decisions
  // accordingly
  // Most tokens are assigned by evaluating a single char, but some require a
  // 'buffer' of multiple chars to assign
  // Ex. Idents, !=, >=, ...
  private void scanToken() {
    char c = advance();
    switch (c) {
      case '(':
        addToken(LEFT_PAREN);
        break;
      case ')':
        addToken(RIGHT_PAREN);
        break;
      case '{':
        addToken(LEFT_BRACE);
        break;
      case '}':
        addToken(RIGHT_BRACE);
        break;
      case ',':
        addToken(COMMA);
        break;
      case '.':
        addToken(DOT);
        break;
      case '-':
        addToken(MINUS);
        break;
      case '+':
        addToken(PLUS);
        break;
      case ';':
        addToken(SEMICOLON);
        break;
      case '*':
        addToken(STAR);
        break;
      case '!':
        addToken(match('=') ? BANG_EQUAL : BANG);
        break;
      case '=':
        addToken(match('=') ? EQUAL_EQUAL : EQUAL);
        break;
      case '<':
        addToken(match('=') ? LESS_EQUAL : LESS);
        break;
      case '>':
        addToken(match('=') ? GREATER_EQUAL : GREATER);
        break;
      case '/':
        if (match('/')) {
          // A comment goes until the end of the line.
          while (peek() != '\n' && !isAtEnd())
            advance(); // We don't care about anything until newline
        } else {
          addToken(SLASH);
        }
        break;

      case '"':
        string();
        break;

      // Ignore whitespace.
      case ' ':
      case '\r':
      case '\t':
        break;

      case '\n':
        line++;
        break;

      // Multi-char encounters
      default:
        if (isDigit(c)) {
          number();
        } else if (isAlpha(c)) {
          identifier();
        } else {
          JClear.error(line, "Unexpected character.");
        }
        break;
    }
  }

  // Simply return the char at the current position in the source code
  private char peek() {
    if (isAtEnd())
      return '\0';
    return source.charAt(current);
  }

  // Returns the char exactly one position ahead of the current position of the source code
  private char peekNext() {
    if (current + 1 >= source.length())
      return '\0';
    return source.charAt(current + 1);
  }

  // aka "consume"
  // Both advances the "current" state attribute and returns the new char under inspection
  private char advance() {
    return source.charAt(current++);
  }

  // Helper function to add a token to the Scanner's tokens List
  private void addToken(TokenType type) {
    addToken(type, null);
  }

  private void addToken(TokenType type, Object literal) {
    String text = source.substring(start, current);
    tokens.add(new Token(type, text, literal, line));
  }

  // Allows the Scanner to make decisions about what 2-char token to assign inline
  // Ex. case '=': addToken(match('=') ? EQUAL_EQUAL : EQUAL);
  // Assigns a EQUAL if there is not an additional '=', but if there is ('=='), it assigns EQUAL_EQUAL
  private boolean match(char expected) {
    if (isAtEnd())
      return false;
    if (source.charAt(current) != expected)
      return false;

    current++;
    return true;
  }

  private void string() {
    while (peek() != '"' && !isAtEnd()) {
      if (peek() == '\n')
        line++;
      advance();
    }

    if (isAtEnd()) {
      JClear.error(line, "Unterminated string.");
      return;
    }

    // The closing ".
    advance();

    // Trim the surrounding quotes.
    String value = source.substring(start + 1, current - 1);
    addToken(STRING, value);
  }

  private boolean isDigit(char c) {
    return c >= '0' && c <= '9';
  }

  private void number() {
    while (isDigit(peek()))
      advance();

    // Look for a fractional part.
    if (peek() == '.' && isDigit(peekNext())) {
      // Consume the "."
      advance();

      while (isDigit(peek()))
        advance();
    }

    addToken(NUMBER,
        Double.parseDouble(source.substring(start, current)));
  }

  private void identifier() {
    while (isAlphaNumeric(peek()))
      advance();

    String text = source.substring(start, current);
    TokenType type = keywords.get(text);
    if (type == null)
      type = IDENTIFIER;
    addToken(type);
  }

  private boolean isAlpha(char c) {
    return (c >= 'a' && c <= 'z') ||
        (c >= 'A' && c <= 'Z') ||
        c == '_';
  }

  private boolean isAlphaNumeric(char c) {
    return isAlpha(c) || isDigit(c);
  }
}
