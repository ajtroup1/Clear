package evaluator

import (
	"fmt"
	"strings"

	"github.com/ajtroup1/clear/ast"
	"github.com/ajtroup1/clear/logger"
	"github.com/ajtroup1/clear/object"
)

var (
	Logger *logger.Logger
	Debug  bool
	Lines []string
)

func Init(l *logger.Logger, debug bool, lines []string) {
	Logger = l
	Debug = debug
	Lines = lines

	if Debug {
		Logger.DefineSection("Evaluation", "Evaluation is simply the traversing of the AST and executing its nodes accordingly.\n\nThe core of the evaluator is the Eval(node) function, which is called recursivly on the AST. Since the AST is a nicely formatted tree structure, it is pretty simple to traverse it recusively.\n\nI would suggest inspecting the [evaluator](../clear/evaluator/) and [object](../clear/object/) package to get a better understanding of how the evaluator works. It's very simple to understand due to its recursive nature.\n\n")
	}
}

// Define a const to easily access object types throughout the evaluator
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Core evaluation function
// Primarily is called with the Program node,
// but is called recursively for all other nodes
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Eval Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		if val == nil {
			return newError("return value is nil: %s", node.Token.Line, node.Token.Col, node.ReturnValue.TokenLiteral())
		}
		return &object.ReturnValue{Value: val, Position: object.Position{Line: node.Token.Line, Col: node.Token.Col}}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)
		return val

	case *ast.WhileStatement:
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		return evalWhileStatement(node, env)

	case *ast.ForStatement:
		return evalForStatement(node, env)
	// Eval Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value, Position: object.Position{Line: node.Token.Line, Col: node.Token.Col}}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value, Position: object.Position{Line: node.Token.Line, Col: node.Token.Col}}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value, Position: object.Position{Line: node.Token.Line, Col: node.Token.Col}}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right, node.Left.TokenLiteral(), env)

	case *ast.PostfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		return evalPostfixExpression(node.Operator, left)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body, Position: object.Position{Line: node.Token.Line, Col: node.Token.Col}}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements, Position: object.Position{Line: node.Token.Line, Col: node.Token.Col}}

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}

	return nil
}

func evalWhileStatement(stmt *ast.WhileStatement, env *object.Environment) object.Object {
	var result object.Object

	for isTruthy(Eval(stmt.Condition, env)) {
		result = evalBlockStatement(stmt.Body, env)
	}

	return result
}

func evalForStatement(stmt *ast.ForStatement, env *object.Environment) object.Object {
	var result object.Object

	if stmt.Init != nil {
		result = Eval(stmt.Init, env)
		if isError(result) {
			return result
		}
	}

	for isTruthy(Eval(stmt.Condition, env)) {
		result = evalBlockStatement(stmt.Body, env)
		if isError(result) {
			return result
		}

		if stmt.Post != nil {
			result = Eval(stmt.Post, env)
			if isError(result) {
				return result
			}
		}

		env.Set(stmt.Post.TokenLiteral(), result)
	}

	return result
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Line(), left.Col(), left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Line(), index.Col(), index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalModuleStatement(stmt *ast.ModuleStatement, env *object.Environment) object.Object {
	module, exists := env.GetModule(stmt.Name.Value)
	if !exists {
		return newError("module not found: %s", stmt.Token.Line, stmt.Token.Col, stmt.Name.Value)
	}

	// fmt.Printf("//env module: %v\n", env.Modules)

	if stmt.ImportAll {
		for name, fn := range module {
			env.Set(name, fn)
		}
	} else {
		for _, importName := range stmt.Imports {
			fn, exists := module[importName.Value]
			if !exists {
				return newError("function %s not found in module %s", importName.Token.Line, importName.Token.Col, importName.Value, stmt.Name.Value)
			}
			env.Set(importName.Value, fn)
			// fmt.Printf("env module: %v\n", env.Modules)
		}
	}

	return nil
}

// Simply iterate over all statements in the program and evaluate them
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Modules {
		result = evalModuleStatement(stmt, env)
		if isError(result) {
			return result
		}
	}

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	if operator != "+" {
		return newError("unknown operator in expression: \"%s %s %s\"", left.Line(), left.Col(),
			left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", right.Line(), right.Col(), operator, right.Type())
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
	literal string,
	env *object.Environment,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		if isCompoundOperator(operator) {
			return evalCompoundAssignment(operator, left, right, env, literal)
		}
		return evalIntegerInfixExpression(operator, left, right)
	case (left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ) || (left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ):
		if isCompoundOperator(operator) {
			return evalCompoundAssignment(operator, left, right, env, literal)
		}
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		if isCompoundOperator(operator) {
			return evalCompoundAssignment(operator, left, right, env, literal)
		}
		return evalFloatInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case operator == "-=":
		return evalInfixExpression("-", left, right, literal, env)
	case operator == "*=":
		return evalInfixExpression("*", left, right, literal, env)
	case operator == "/=":
		return evalInfixExpression("/", left, right, literal, env)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Line(), left.Col(),
			left.Type(), operator, right.Type())
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Line(), left.Col(),
			left.Type(), operator, right.Type())
	}
}

func isCompoundOperator(operator string) bool {
	switch operator {
	case "+=", "-=", "*=", "/=":
		return true
	default:
		return false
	}
}

func evalCompoundAssignment(
	operator string,
	left, right object.Object,
	env *object.Environment,
	literal string,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		result := evalIntegerInfixExpression(operator, left, right)
		env.Set(literal, result)
		return result
	case (left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ) || (left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ):
		result := evalFloatInfixExpression(operator, left, right)
		env.Set(literal, result)
		return result
	// case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
	// 	return evalFloatInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Line(), left.Col(),
			left.Type(), operator, right.Type())
	}
}



func evalPostfixExpression(operator string, left object.Object) object.Object {
	if left.Type() != object.INTEGER_OBJ && left.Type() != object.FLOAT_OBJ {
		return newError("unknown operator: %s%s", left.Line(), left.Col(), operator, left.Type())
	}

	if left.Type() == object.INTEGER_OBJ {

		leftVal := left.(*object.Integer).Value

		switch operator {
		case "++":
			return &object.Integer{Value: leftVal + 1}
		case "--":
			return &object.Integer{Value: leftVal - 1}
		default:
			return newError("unknown operator: %s%s", left.Line(), left.Col(), operator, left.Type())
		}
	}

	if left.Type() == object.FLOAT_OBJ {
		leftVal := left.(*object.Float).Value

		switch operator {
		case "++":
			return &object.Float{Value: leftVal + 1}
		case "--":
			return &object.Float{Value: leftVal - 1}
		default:
			return newError("unknown operator: %s%s", left.Line(), left.Col(), operator, left.Type())
		}
	}

	return newError("unknown operator: %s%s", left.Line(), left.Col(), operator, left.Type())
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ && right.Type() != object.FLOAT_OBJ {
		return newError("unknown operator: -%s", right.Line(), right.Col(), right.Type())
	}

	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}
	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}

	return newError("unknown operator: -%s", right.Line(), right.Col(), right.Type())
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "+=":
		return &object.Integer{Value: leftVal + rightVal}
	case "-=":
		return &object.Integer{Value: leftVal - rightVal}
	case "*=":
		return &object.Integer{Value: leftVal * rightVal}
	case "/=":
		return &object.Integer{Value: leftVal / rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Line(), left.Col(),
			left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	var leftVal, rightVal float64

	switch l := left.(type) {
	case *object.Float:
		leftVal = l.Value
	case *object.Integer:
		leftVal = float64(l.Value)
	default:
		return newError("type mismatch: %s %s %s", left.Line(), left.Col(), left.Type(), operator, right.Type())
	}

	switch r := right.(type) {
	case *object.Float:
		rightVal = r.Value
	case *object.Integer:
		rightVal = float64(r.Value)
	default:
		return newError("type mismatch: %s %s %s", left.Line(), left.Col(), left.Type(), operator, right.Type())
	}

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "+=":
		return &object.Float{Value: leftVal + rightVal}
	case "-=":
		return &object.Float{Value: leftVal - rightVal}
	case "*=":
		return &object.Float{Value: leftVal * rightVal}
	case "/=":
		return &object.Float{Value: leftVal / rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Line(), left.Col(),
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(
	ie *ast.IfExpression,
	env *object.Environment,
) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	parts := strings.Split(node.Value, ".")
	if len(parts) == 2 {
		moduleName, functionName := parts[0], parts[1]
		if module, exists := env.Modules[moduleName]; exists {
			if builtin, found := module[functionName]; found {
				return builtin
			}
			return newError("function not found in module '%s': %s", node.Token.Line, node.Token.Col, moduleName, functionName)
		}
		return newError("module not found: %s", node.Token.Line, node.Token.Col, moduleName)
	}

	return newError("identifier not found: %s", node.Token.Line, node.Token.Col, node.Value)
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, line, col int, a ...interface{}) *object.Error {
	// for _, line := range Lines {
	// 	fmt.Printf("// %s //\n", line)
	// }
	return &object.Error{Message: fmt.Sprintf(format, a...), Position: object.Position{Line: line, Col: col}, Context: Lines[line-1]}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Line(), fn.Col(), fn.Type())
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Line(), key.Col(), key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}
