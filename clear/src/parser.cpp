#include "parser.hpp"
#include <stdexcept>

Parser::Parser(const std::vector<Token>& tokens)
    : tokens(tokens), pos(0) {
    nextToken();
}

Token Parser::peekToken() const {
    if (pos < tokens.size()) {
        return tokens[pos];
    }
    return Token(TokenType::UNDEFINED, "", 0, 0); 
}

Token Parser::consumeToken() {
    Token token = currentToken;
    nextToken();
    return token;
}

void Parser::nextToken() {
    if (pos < tokens.size()) {
        currentToken = tokens[pos++];
    }
}

std::unique_ptr<Program> Parser::parse() {
    return parseProgram();
}

std::unique_ptr<Program> Parser::parseProgram() {
    auto program = std::make_unique<Program>();
    program->setStatements(parseBlockStatement());
    return program;
}

std::unique_ptr<BlockStatement> Parser::parseBlockStatement() {
    auto block = std::make_unique<BlockStatement>();
    consumeToken(); // Consume the '{' token
    while (currentToken.getType() != TokenType::RIGHT_BRACE && currentToken.getType() != TokenType::_EOF) {
        block->addStatement(parseStatement());
    }
    consumeToken(); // Consume the '}' token
    return block;
}

std::unique_ptr<Statement> Parser::parseStatement() {
    if (currentToken.getType() == TokenType::LET) {
        return parseLetStatement();
    } else if (currentToken.getType() == TokenType::RETURN) {
        return parseReturnStatement();
    } else if (currentToken.getType() == TokenType::WHILE) {
        return parseWhileStatement();
    } else if (currentToken.getType() == TokenType::FOR) {
        return parseForStatement();
    } else if (currentToken.getType() == TokenType::BREAK) {
        return parseBreakStatement();
    } else if (currentToken.getType() == TokenType::CONTINUE) {
        return parseContinueStatement();
    } else {
        return parseExpressionStatement();
    }
}

std::unique_ptr<Expression> Parser::parseExpression() {
    // Implement expression parsing
    return nullptr;
}

std::unique_ptr<Identifier> Parser::parseIdentifier() {
    Token token = consumeToken();
    return std::make_unique<Identifier>(token.getLiteral());
}

std::unique_ptr<IntegerLiteral> Parser::parseIntegerLiteral() {
    Token token = consumeToken();
    return std::make_unique<IntegerLiteral>(std::stoi(token.getLiteral()));
}

std::unique_ptr<FloatLiteral> Parser::parseFloatLiteral() {
    Token token = consumeToken();
    return std::make_unique<FloatLiteral>(std::stof(token.getLiteral()));
}

std::unique_ptr<StringLiteral> Parser::parseStringLiteral() {
    Token token = consumeToken();
    return std::make_unique<StringLiteral>(token.getLiteral());
}

std::unique_ptr<BooleanLiteral> Parser::parseBooleanLiteral() {
    Token token = consumeToken();
    return std::make_unique<BooleanLiteral>(token.getType() == TokenType::TRUE);
}

std::unique_ptr<LetStatement> Parser::parseLetStatement() {
    // Implement let statement parsing
    return nullptr;
}

std::unique_ptr<ReturnStatement> Parser::parseReturnStatement() {
    // Implement return statement parsing
    return nullptr;
}

std::unique_ptr<IfExpression> Parser::parseIfExpression() {
    // Implement if expression parsing
    return nullptr;
}

std::unique_ptr<WhileStatement> Parser::parseWhileStatement() {
    // Implement while statement parsing
    return nullptr;
}

std::unique_ptr<ForStatement> Parser::parseForStatement() {
    // Implement for statement parsing
    return nullptr;
}

std::unique_ptr<Statement> Parser::parseBreakStatement() {
    // Implement break statement parsing
    return nullptr;
}

std::unique_ptr<Statement> Parser::parseContinueStatement() {
    // Implement continue statement parsing
    return nullptr;
}

std::unique_ptr<ExpressionStatement> Parser::parseExpressionStatement() {
    // Implement expression statement parsing
    return nullptr;
}

std::unique_ptr<FunctionLiteral> Parser::parseFunctionLiteral() {
    // Implement function literal parsing
    return nullptr;
}

std::unique_ptr<UnnamedFunctionLiteral> Parser::parseUnnamedFunctionLiteral() {
    // Implement unnamed function literal parsing
    return nullptr;
}

std::unique_ptr<CallExpression> Parser::parseCallExpression() {
    // Implement call expression parsing
    return nullptr;
}

std::unique_ptr<MemberExpression> Parser::parseMemberExpression() {
    // Implement member expression parsing
    return nullptr;
}

std::unique_ptr<ArrayExpression> Parser::parseArrayExpression() {
    // Implement array expression parsing
    return nullptr;
}

std::unique_ptr<IndexExpression> Parser::parseIndexExpression() {
    // Implement index expression parsing
    return nullptr;
}

std::unique_ptr<AssignmentExpression> Parser::parseAssignmentExpression() {
    // Implement assignment expression parsing
    return nullptr;
}

std::unique_ptr<BinaryExpression> Parser::parseBinaryExpression() {
    // Implement binary expression parsing
    return nullptr;
}

std::unique_ptr<PrefixExpression> Parser::parsePrefixExpression() {
    // Implement prefix expression parsing
    return nullptr;
}

std::unique_ptr<IfExpression> Parser::parseIfExpression() {
    // Implement if expression parsing
    return nullptr;
}