#ifndef CLEAR_AST_HPP
#define CLEAR_AST_HPP

#include <memory>
#include <string>
#include <vector>

// Base class for all AST nodes
class ASTNode {
public:
    virtual ~ASTNode() = default;
    virtual std::string stringify() const = 0;
};

// In Clear (or any other language), the AST is composed of expressions within statements
// Statements make up the entire program, while expressions are the building blocks of the program of represent the computational aspects of the code
class Expression : public ASTNode {
public:
    virtual ~Expression() = default;
};

class Statement : public ASTNode {
public:
    virtual ~Statement() = default;
};

// ----------
// STATEMENTS
// ----------

class BlockStatement : public Statement {
public:
    virtual ~BlockStatement() = default;
    std::string stringify() const override {
        std::string result = "BlockStatement(";
        for (const auto& stmt : statements) {
            result += stmt->stringify() + "; ";
        }
        result += ")";
        return result;
    }
    void addStatement(std::unique_ptr<Statement> stmt) {
        statements.push_back(std::move(stmt));
    }

private:
    // The BlockStatement class contains a vector of Statement pointers
    // This vector will store all the statements in the block
    std::vector<std::unique_ptr<Statement>> statements;
};

// Program class represents the entire program
class Program : public ASTNode {
public:
    virtual ~Program() = default;
    std::string stringify() const override {
        return "Program(" + (statements ? statements->stringify() : "null") + ")";
    }
    BlockStatement* getStatements() const { return statements.get(); }
    void setStatements(std::unique_ptr<BlockStatement> stmts) {
        statements = std::move(stmts);
    }

private:
    // The Program class contains a block statement
    // These statements are the top-level statements in the program
    std::unique_ptr<BlockStatement> statements;
};

// Other AST nodes...

class LetStatement : public Statement {
public:
    virtual ~LetStatement() = default;
    std::string stringify() const override {
        return "LetStatement(" + name + ", " + value->stringify() + ")";
    }

private:
    // The LetStatement class contains a name and an expression
    // The name is the identifier of the variable being declared
    // The expression is the value being assigned to the variable
    std::string name;
    std::unique_ptr<Expression> value;
};

class ConstStatement : public Statement {
public:
    virtual ~ConstStatement() = default;
    std::string stringify() const override {
        return "ConstStatement(" + name + ", " + value->stringify() + ")";
    }

private:
    // The ConstStatement class contains a name and an expression
    // The name is the identifier of the variable being declared
    // The expression is the value being assigned to the variable
    std::string name;
    std::unique_ptr<Expression> value;
};

class ReturnStatement : public Statement {
public:
    virtual ~ReturnStatement() = default;
    std::string stringify() const override {
        return "ReturnStatement(" + value->stringify() + ")";
    }

private:
    // The ReturnStatement class contains an expression
    // This expression represents the value being returned
    std::unique_ptr<Expression> value;
};

class ExpressionStatement : public Statement {
public:
    virtual ~ExpressionStatement() = default;
    std::string stringify() const override {
        return "ExpressionStatement(" + expression->stringify() + ")";
    }

private:
    // The ExpressionStatement class contains an expression
    // This is necessary for expressions to be contained within a collection of
    // statements, which is possible in Clear
    std::unique_ptr<Expression> expression;
};

class WhileStatement : public Statement {
public:
    virtual ~WhileStatement() = default;
    std::string stringify() const override {
        return "WhileStatement(" + condition->stringify() + ", " +
               body->stringify() + ")";
    }

private:
    // The WhileStatement class contains a condition and a block statement
    // The condition is the expression that is evaluated to determine the
    // truthiness of the while loop The body is the block of statements that is
    // executed while the condition is true
    std::unique_ptr<Expression> condition;
    std::unique_ptr<BlockStatement> body;
};

class ForStatement : public Statement {
public:
    virtual ~ForStatement() = default;
    std::string stringify() const override {
        return "ForStatement(" + initializer->stringify() + ", " +
               condition->stringify() + ", " + increment->stringify() + ", " +
               body->stringify() + ")";
    }

private:
    // The ForStatement class contains an initializer, a condition, an increment,
    // and a block statement The initializer is the expression that is executed
    // before the loop starts The condition is the expression that is evaluated to
    // determine the truthiness of the for loop The increment is the expression
    // that is executed after each iteration of the loop The body is the block of
    // statements that is executed while the condition is true
    std::unique_ptr<Expression> initializer;
    std::unique_ptr<Expression> condition;
    std::unique_ptr<Expression> increment;
    std::unique_ptr<BlockStatement> body;
};

// -----------
// EXPRESSIONS
// -----------

class Identifier : public Expression {
public:
    virtual ~Identifier() = default;
    std::string stringify() const override { return "Identifier(" + value + ")"; }

private:
    // The Identifier class contains a string value
    // This value represents the name or literal of the identifier (x, foo, etc.)
    std::string value;
};

class IntegerLiteral : public Expression {
public:
    virtual ~IntegerLiteral() = default;
    std::string stringify() const override {
        return "IntegerLiteral(" + std::to_string(value) + ")";
    }

private:
    // The IntegerLiteral class contains an integer value
    // This value represents the literal integer value in the program
    int value;
};

class FloatLiteral : public Expression {
public:
    virtual ~FloatLiteral() = default;
    std::string stringify() const override {
        return "FloatLiteral(" + std::to_string(value) + ")";
    }

private:
    // The FloatLiteral class contains a float value
    // This value represents the literal float value in the program
    float value;
};

class StringLiteral : public Expression {
public:
    virtual ~StringLiteral() = default;
    std::string stringify() const override { return "StringLiteral(" + value + ")"; }

private:
    // The StringLiteral class contains a string value
    // This value represents the literal string value in the program
    std::string value;
};

class BooleanLiteral : public Expression {
public:
    virtual ~BooleanLiteral() = default;
    std::string stringify() const override {
        return "Boolean(" + std::string(value ? "true" : "false") + ")";
    }

private:
    // The Boolean class contains a boolean value
    // This value represents the literal boolean value in the program
    bool value;
};

class PrefixExpression : public Expression {
public:
    virtual ~PrefixExpression() = default;
    std::string stringify() const override {
        return "PrefixExpression(" + op + ", " + right->stringify() + ")";
    }

private:
    // The PrefixExpression class contains an operator and an expression
    // The operator is the prefix operator (!, -, etc.)
    // The expression is the right-hand side of the operator
    // As opposed to binary operators, prefix operators only have a "right hand
    // side" and no "left hand side"
    std::string op;
    std::unique_ptr<Expression> right;
};

class BinaryExpression : public Expression {
public:
    virtual ~BinaryExpression() = default;
    std::string stringify() const override {
        return "BinaryExpression(" + left->stringify() + ", " + op + ", " +
               right->stringify() + ")";
    }

private:
    // The BinaryExpression class contains a left expression, an operator, and a
    // right expression The left expression is the left-hand side of the operator
    // The operator is the binary arethmetic operator (+, -, *, etc.)
    // The right expression is the right-hand side of the operator
    std::unique_ptr<Expression> left;
    std::string op;
    std::unique_ptr<Expression> right;
};

class GroupedExpression : public Expression {
public:
    virtual ~GroupedExpression() = default;
    std::string stringify() const override {
        return "GroupedExpression(" + expression->stringify() + ")";
    }

private:
    // The GroupedExpression class contains an expression
    // This expression is enclosed in parentheses
    std::unique_ptr<Expression> expression;
};

class IfExpression : public Expression {
public:
    virtual ~IfExpression() = default;
    std::string stringify() const override {
        return "IfExpression(" + condition->stringify() + ", " +
               consequence->stringify() + ", " +
               (alternative ? alternative->stringify() : "no alternative defined") +
               ")";
    }

private:
    // The IfExpression class contains a condition, a consequence, and optionally
    // an alternative The condition is the expression that is evaluated to
    // determine the truthiness of the if statement The consequence is the block
    // of statements that is executed if the condition is true The alternative is
    // the block of statements that is executed if the condition is false
    std::unique_ptr<Expression> condition;
    std::unique_ptr<BlockStatement> consequence;
    std::unique_ptr<BlockStatement> alternative;
};

class FunctionLiteral : public Expression {
public:
    virtual ~FunctionLiteral() = default;
    std::string stringify() const override {
        std::string result = "FunctionLiteral(";
        result += "params: [";
        for (const auto& param : parameters) {
            result += param + ", ";
        }
        result += "], body: " + body->stringify() + ")";
        return result;
    }

private:
    // The FunctionLiteral class contains a vector of parameters and a block
    // statement The parameters are the identifiers of the function's parameters
    // The body is the block of statements that make up the function's body
    std::string name;
    std::vector<std::string> parameters;
    std::unique_ptr<BlockStatement> body;
};

class UnnamedFunctionLiteral : public Expression {
public:
    virtual ~UnnamedFunctionLiteral() = default;
    std::string stringify() const override {
        std::string result = "UnnamedFunctionLiteral(";
        result += "params: [";
        for (const auto& param : parameters) {
            result += param + ", ";
        }
        result += "], body: " + body->stringify() + ")";
        return result;
    }

private:
    // The UnnamedFunctionLiteral class contains a vector of parameters and a
    // block statement The parameters are the identifiers of the function's
    // parameters The body is the block of statements that make up the function's
    // body Unnamed functions are functions without a name, examples include: let
    // add = function(x, y) { return x + y; } <--- Clear  result = ((x) => x
    // * 2)(5); <--- JavaScript
    std::vector<std::string> parameters;
    std::unique_ptr<BlockStatement> body;
};

class CallExpression : public Expression {
public:
    virtual ~CallExpression() = default;
    std::string stringify() const override {
        std::string result =
            "CallExpression(" + function->stringify() + ", args: [";
        for (const auto& arg : arguments) {
            result += arg->stringify() + ", ";
        }
        result += "])";
        return result;
    }

private:
    // The CallExpression class contains a function expression and a vector of
    // argument expressions The function expression is the function being called
    // The argument expressions are the arguments being passed to the function
    std::unique_ptr<Expression> function;
    std::vector<std::unique_ptr<Expression>> arguments;
};

class MemberExpression : public Expression {
public:
    virtual ~MemberExpression() = default;
    std::string stringify() const override {
        return "MemberExpression(" + object->stringify() + ", " + property + ")";
    }

private:
    // The MemberExpression class contains an object expression and a property
    // The object expression is the object being accessed
    // The property is the property being returned
    std::unique_ptr<Expression> object;
    std::string property;
};

class AssignmentExpression : public Expression {
public:
    virtual ~AssignmentExpression() = default;
    std::string stringify() const override {
        return "AssignmentExpression(" + left->stringify() + ", " +
               right->stringify() + ")";
    }

private:
    // The AssignmentExpression class contains a left expression and a right
    // expression The left expression is the variable being assigned to The right
    // expression is the value being assigned to the variable This differentiates
    // it from the LetStatement, which is used for variable declaration The
    // AssignmentExpression is used for assignment of a previously declared
    // variable
    std::unique_ptr<Expression> left;
    std::unique_ptr<Expression> right;
};

class ArrayExpression : public Expression {
public:
    virtual ~ArrayExpression() = default;
    std::string stringify() const override {
        std::string result = "ArrayExpression([";
        for (const auto& elem : elements) {
            result += elem->stringify() + ", ";
        }
        result += "])";
        return result;
    }

private:
    // The ArrayExpression class contains a vector of expressions
    // These expressions are the elements of the array
    std::vector<std::unique_ptr<Expression>> elements;
};

class IndexExpression : public Expression {
public:
    virtual ~IndexExpression() = default;
    std::string stringify() const override {
        return "IndexExpression(" + left->stringify() + ", " + index->stringify() +
               ")";
    }

private:
    // The IndexExpression class contains a left expression and an index
    // expression The left expression is the object being accessed The index
    // expression is the index being accessed
    std::unique_ptr<Expression> left;
    std::unique_ptr<Expression> index;
};

#endif // CLEAR_AST_HPP