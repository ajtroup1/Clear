#ifndef CLEAR_PARSER_HPP
#define CLEAR_PARSER_HPP

#include "ast.hpp"
#include "lexer.hpp"
#include <vector>
#include <memory>

class Parser {
public:
    Parser(const std::vector<Token>& tokens);
    ~Parser() = default;
    std::unique_ptr<Program> parse(std::vector<Token>& tokens);

private:
    std::vector<Token> tokens;
    size_t pos = 0;
    Token currentToken;
    Token peekToken() const;
    Token consumeToken();
    void nextToken();

    std::unique_ptr<Program> parseProgram();
    std::unique_ptr<Statement> parseStatement();
    std::unique_ptr<Expression> parseExpression();
    std::unique_ptr<Identifier> parseIdentifier();
    std::unique_ptr<IntegerLiteral> parseIntegerLiteral();
    std::unique_ptr<FloatLiteral> parseFloatLiteral();
    std::unique_ptr<StringLiteral> parseStringLiteral();
    std::unique_ptr<BooleanLiteral> parseBooleanLiteral();
    std::unique_ptr<FunctionLiteral> parseFunctionLiteral();
    std::unique_ptr<UnnamedFunctionLiteral> parseUnnamedFunctionLiteral();
    std::unique_ptr<CallExpression> parseCallExpression();
    std::unique_ptr<MemberExpression> parseMemberExpression();
    std::unique_ptr<ArrayExpression> parseArrayExpression();
    std::unique_ptr<IndexExpression> parseIndexExpression();
    std::unique_ptr<AssignmentExpression> parseAssignmentExpression();
    std::unique_ptr<BinaryExpression> parseBinaryExpression();
    std::unique_ptr<PrefixExpression> parsePrefixExpression();
    std::unique_ptr<IfExpression> parseIfExpression();
    std::unique_ptr<WhileStatement> parseWhileStatement();
    std::unique_ptr<ForStatement> parseForStatement();
    std::unique_ptr<ReturnStatement> parseReturnStatement();
    std::unique_ptr<Statement> parseBreakStatement(); 
    std::unique_ptr<Statement> parseContinueStatement();
    std::unique_ptr<BlockStatement> parseBlockStatement();
    std::unique_ptr<ExpressionStatement> parseExpressionStatement();
    std::unique_ptr<LetStatement> parseLetStatement();
    std::unique_ptr<Statement> parseConstStatement();
};

#endif