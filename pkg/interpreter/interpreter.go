package interpreter

import (
	"fmt"

	"github.com/aether-lang/aether/pkg/ast"
	"github.com/aether-lang/aether/pkg/object"
)

// Interpreter evaluates AST nodes
type Interpreter struct {
	env *object.Environment
}

// New creates a new interpreter instance
func New() *Interpreter {
	return &Interpreter{
		env: object.NewEnvironment(),
	}
}

// Eval evaluates an AST node and returns an object
func (interp *Interpreter) Eval(node ast.Node) (object.Object, error) {
	switch n := node.(type) {
	// Program
	case *ast.Program:
		return interp.evalProgram(n)

	// Statements
	case *ast.ExpressionStatement:
		return interp.Eval(n.Expression)

	case *ast.AssignmentStatement:
		val, err := interp.Eval(n.Value)
		if err != nil {
			return nil, err
		}
		if isError(val) {
			return val, nil
		}
		interp.env.Set(n.Name.Value, val)
		return val, nil

	case *ast.LetStatement:
		val, err := interp.Eval(n.Value)
		if err != nil {
			return nil, err
		}
		if isError(val) {
			return val, nil
		}
		interp.env.Set(n.Name, val)
		return val, nil

	case *ast.ResourceStatement:
		return interp.evalResourceStatement(n)

	case *ast.VariableStatement:
		return interp.evalVariableStatement(n)

	case *ast.OutputStatement:
		return interp.evalOutputStatement(n)

	// Expressions - Literals
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}, nil

	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}, nil

	case *ast.NumberLiteral:
		return &object.Float{Value: n.Value}, nil

	case *ast.StringLiteral:
		return &object.String{Value: n.Value}, nil

	case *ast.BooleanLiteral:
		return &object.Boolean{Value: n.Value}, nil

	case *ast.Identifier:
		return interp.evalIdentifier(n)

	case *ast.ArrayLiteral:
		return interp.evalArrayLiteral(n)

	case *ast.ListLiteral:
		elements, err := interp.evalExpressions(n.Elements)
		if err != nil {
			return nil, err
		}
		return &object.Array{Elements: elements}, nil

	case *ast.MapLiteral:
		return interp.evalMapLiteral(n)

	case *ast.BinaryExpression:
		return interp.evalBinaryExpression(n)

	case *ast.UnaryExpression:
		return interp.evalUnaryExpression(n)

	case *ast.IndexExpression:
		return interp.evalIndexExpression(n)

	case *ast.BlockExpression:
		return interp.evalBlockExpression(n)

	default:
		return nil, fmt.Errorf("unknown node type: %T", node)
	}
}

// evalProgram evaluates a program
func (interp *Interpreter) evalProgram(program *ast.Program) (object.Object, error) {
	var result object.Object
	var err error

	for _, stmt := range program.Statements {
		result, err = interp.Eval(stmt)
		if err != nil {
			return nil, err
		}
		if isError(result) {
			return result, nil
		}
	}

	if result == nil {
		return &object.Null{}, nil
	}
	return result, nil
}

// evalIdentifier evaluates an identifier
func (interp *Interpreter) evalIdentifier(node *ast.Identifier) (object.Object, error) {
	val, ok := interp.env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value), nil
	}
	return val, nil
}

// evalArrayLiteral evaluates an array literal
func (interp *Interpreter) evalArrayLiteral(node *ast.ArrayLiteral) (object.Object, error) {
	elements, err := interp.evalExpressions(node.Elements)
	if err != nil {
		return nil, err
	}
	return &object.Array{Elements: elements}, nil
}

// evalMapLiteral evaluates a map literal
func (interp *Interpreter) evalMapLiteral(node *ast.MapLiteral) (object.Object, error) {
	pairs := make(map[string]object.Object)

	for key, valueExpr := range node.Pairs {
		value, err := interp.Eval(valueExpr)
		if err != nil {
			return nil, err
		}
		if isError(value) {
			return value, nil
		}
		pairs[key] = value
	}

	return &object.Map{Pairs: pairs}, nil
}

// evalBinaryExpression evaluates a binary expression
func (interp *Interpreter) evalBinaryExpression(node *ast.BinaryExpression) (object.Object, error) {
	left, err := interp.Eval(node.Left)
	if err != nil {
		return nil, err
	}
	if isError(left) {
		return left, nil
	}

	right, err := interp.Eval(node.Right)
	if err != nil {
		return nil, err
	}
	if isError(right) {
		return right, nil
	}

	return interp.evalInfixExpression(node.Operator, left, right)
}

// evalInfixExpression evaluates infix operations
func (interp *Interpreter) evalInfixExpression(operator string, left, right object.Object) (object.Object, error) {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return interp.evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ:
		return interp.evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return interp.evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return &object.Boolean{Value: objectsEqual(left, right)}, nil
	case operator == "!=":
		return &object.Boolean{Value: !objectsEqual(left, right)}, nil
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type()), nil
	}
}

// evalIntegerInfixExpression evaluates integer operations
func (interp *Interpreter) evalIntegerInfixExpression(operator string, left, right object.Object) (object.Object, error) {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}, nil
	case "-":
		return &object.Integer{Value: leftVal - rightVal}, nil
	case "*":
		return &object.Integer{Value: leftVal * rightVal}, nil
	case "/":
		if rightVal == 0 {
			return newError("division by zero"), nil
		}
		return &object.Integer{Value: leftVal / rightVal}, nil
	case "<":
		return &object.Boolean{Value: leftVal < rightVal}, nil
	case ">":
		return &object.Boolean{Value: leftVal > rightVal}, nil
	case "<=":
		return &object.Boolean{Value: leftVal <= rightVal}, nil
	case ">=":
		return &object.Boolean{Value: leftVal >= rightVal}, nil
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}, nil
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}, nil
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type()), nil
	}
}

// evalFloatInfixExpression evaluates float operations
func (interp *Interpreter) evalFloatInfixExpression(operator string, left, right object.Object) (object.Object, error) {
	leftVal := toFloat(left)
	rightVal := toFloat(right)

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}, nil
	case "-":
		return &object.Float{Value: leftVal - rightVal}, nil
	case "*":
		return &object.Float{Value: leftVal * rightVal}, nil
	case "/":
		if rightVal == 0 {
			return newError("division by zero"), nil
		}
		return &object.Float{Value: leftVal / rightVal}, nil
	case "<":
		return &object.Boolean{Value: leftVal < rightVal}, nil
	case ">":
		return &object.Boolean{Value: leftVal > rightVal}, nil
	case "<=":
		return &object.Boolean{Value: leftVal <= rightVal}, nil
	case ">=":
		return &object.Boolean{Value: leftVal >= rightVal}, nil
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}, nil
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}, nil
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type()), nil
	}
}

// evalStringInfixExpression evaluates string operations
func (interp *Interpreter) evalStringInfixExpression(operator string, left, right object.Object) (object.Object, error) {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}, nil
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}, nil
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}, nil
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type()), nil
	}
}

// evalUnaryExpression evaluates unary operations
func (interp *Interpreter) evalUnaryExpression(node *ast.UnaryExpression) (object.Object, error) {
	right, err := interp.Eval(node.Right)
	if err != nil {
		return nil, err
	}
	if isError(right) {
		return right, nil
	}

	switch node.Operator {
	case "!":
		return &object.Boolean{Value: !isTruthy(right)}, nil
	case "-":
		switch obj := right.(type) {
		case *object.Integer:
			return &object.Integer{Value: -obj.Value}, nil
		case *object.Float:
			return &object.Float{Value: -obj.Value}, nil
		default:
			return newError("unknown operator: -%s", right.Type()), nil
		}
	default:
		return newError("unknown operator: %s%s", node.Operator, right.Type()), nil
	}
}

// evalIndexExpression evaluates index access
func (interp *Interpreter) evalIndexExpression(node *ast.IndexExpression) (object.Object, error) {
	left, err := interp.Eval(node.Left)
	if err != nil {
		return nil, err
	}
	if isError(left) {
		return left, nil
	}

	index, err := interp.Eval(node.Index)
	if err != nil {
		return nil, err
	}
	if isError(index) {
		return index, nil
	}

	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return interp.evalArrayIndexExpression(left, index)
	case left.Type() == object.MAP_OBJ && index.Type() == object.STRING_OBJ:
		return interp.evalMapIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type()), nil
	}
}

// evalArrayIndexExpression evaluates array indexing
func (interp *Interpreter) evalArrayIndexExpression(array, index object.Object) (object.Object, error) {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)

	if idx < 0 || idx > max {
		return &object.Null{}, nil
	}

	return arrayObj.Elements[idx], nil
}

// evalMapIndexExpression evaluates map indexing
func (interp *Interpreter) evalMapIndexExpression(mapObj, index object.Object) (object.Object, error) {
	m := mapObj.(*object.Map)
	key := index.(*object.String).Value

	val, ok := m.Pairs[key]
	if !ok {
		return &object.Null{}, nil
	}

	return val, nil
}

// evalBlockExpression evaluates a block expression
func (interp *Interpreter) evalBlockExpression(block *ast.BlockExpression) (object.Object, error) {
	var result object.Object
	var err error

	for _, stmt := range block.Statements {
		result, err = interp.Eval(stmt)
		if err != nil {
			return nil, err
		}
		if isError(result) {
			return result, nil
		}
	}

	if result == nil {
		return &object.Null{}, nil
	}
	return result, nil
}

// evalResourceStatement evaluates a resource statement
func (interp *Interpreter) evalResourceStatement(node *ast.ResourceStatement) (object.Object, error) {
	properties := make(map[string]object.Object)

	if node.Properties != nil {
		// Evaluate properties from block expression
		for _, stmt := range node.Properties.Statements {
			if assignStmt, ok := stmt.(*ast.AssignmentStatement); ok {
				val, err := interp.Eval(assignStmt.Value)
				if err != nil {
					return nil, err
				}
				if isError(val) {
					return val, nil
				}
				properties[assignStmt.Name.Value] = val
			}
		}
	} else if node.Attributes != nil {
		// Evaluate attributes from map literal (deprecated)
		for key, valueExpr := range node.Attributes.Pairs {
			value, err := interp.Eval(valueExpr)
			if err != nil {
				return nil, err
			}
			if isError(value) {
				return value, nil
			}
			properties[key] = value
		}
	}

	resource := &object.Resource{
		ResourceType: node.Type,
		Name:         node.Name,
		Properties:   properties,
	}

	// Store the resource in the environment by name
	interp.env.Set(node.Name, resource)

	return resource, nil
}

// evalVariableStatement evaluates a variable statement
func (interp *Interpreter) evalVariableStatement(node *ast.VariableStatement) (object.Object, error) {
	var defaultValue object.Object
	var varType string
	var description string

	if node.Default != nil {
		val, err := interp.Eval(node.Default)
		if err != nil {
			return nil, err
		}
		if isError(val) {
			return val, nil
		}
		defaultValue = val
	}

	if node.Type != nil {
		varType = node.Type.Value
	}

	if node.Description != nil {
		description = node.Description.Value
	}

	variable := &object.Variable{
		Name:        node.Name,
		VarType:     varType,
		Default:     defaultValue,
		Description: description,
	}

	// Store the variable in the environment
	interp.env.Set(node.Name, variable)

	return variable, nil
}

// evalOutputStatement evaluates an output statement
func (interp *Interpreter) evalOutputStatement(node *ast.OutputStatement) (object.Object, error) {
	var value object.Object
	var err error

	if node.Value != nil {
		value, err = interp.Eval(node.Value)
		if err != nil {
			return nil, err
		}
		if isError(value) {
			return value, nil
		}
	}

	output := &object.Output{
		Name:  node.Name,
		Value: value,
	}

	// Store the output in the environment
	interp.env.Set(node.Name, output)

	return output, nil
}

// evalExpressions evaluates a list of expressions
func (interp *Interpreter) evalExpressions(exps []ast.Expression) ([]object.Object, error) {
	result := make([]object.Object, 0, len(exps))

	for _, exp := range exps {
		evaluated, err := interp.Eval(exp)
		if err != nil {
			return nil, err
		}
		if isError(evaluated) {
			return nil, fmt.Errorf("%s", evaluated.Inspect())
		}
		result = append(result, evaluated)
	}

	return result, nil
}

// Helper functions

// isError checks if an object is an error
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// newError creates a new error object
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// isTruthy determines if an object is truthy
func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Null:
		return false
	case *object.Boolean:
		return obj.Value
	case *object.Integer:
		return obj.Value != 0
	case *object.Float:
		return obj.Value != 0.0
	case *object.String:
		return obj.Value != ""
	default:
		return true
	}
}

// objectsEqual checks if two objects are equal
func objectsEqual(left, right object.Object) bool {
	if left.Type() != right.Type() {
		return false
	}

	switch leftVal := left.(type) {
	case *object.Integer:
		return leftVal.Value == right.(*object.Integer).Value
	case *object.Float:
		return leftVal.Value == right.(*object.Float).Value
	case *object.Boolean:
		return leftVal.Value == right.(*object.Boolean).Value
	case *object.String:
		return leftVal.Value == right.(*object.String).Value
	case *object.Null:
		return true
	default:
		return false
	}
}

// toFloat converts an object to a float64
func toFloat(obj object.Object) float64 {
	switch obj := obj.(type) {
	case *object.Integer:
		return float64(obj.Value)
	case *object.Float:
		return obj.Value
	default:
		return 0.0
	}
}
