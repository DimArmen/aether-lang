package ast

import (
	"bytes"
	"strings"

	"github.com/aether-lang/aether/pkg/lexer"
)

// Node represents a node in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents a statement node
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ============================================================================
// Statement Types
// ============================================================================

// ResourceStatement represents a resource declaration
type ResourceStatement struct {
	Token      lexer.Token      // the 'resource' token
	Type       string           // e.g., "compute.instance"
	Name       string           // e.g., "web_server"
	Properties *BlockExpression // resource properties block
	Attributes *MapLiteral      // (deprecated, use Properties)
}

func (rs *ResourceStatement) statementNode()       {}
func (rs *ResourceStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ResourceStatement) String() string {
	var out bytes.Buffer
	out.WriteString("resource ")
	out.WriteString(rs.Type)
	out.WriteString(" \"")
	out.WriteString(rs.Name)
	out.WriteString("\" ")
	if rs.Attributes != nil {
		out.WriteString(rs.Attributes.String())
	}
	return out.String()
}

// VariableStatement represents a variable declaration
type VariableStatement struct {
	Token       lexer.Token      // the 'variable' token
	Name        *Identifier      // variable name
	Type        *Identifier      // variable type (optional)
	Default     Expression       // default value (optional)
	Description *StringLiteral   // description (optional)
}

func (vs *VariableStatement) statementNode()       {}
func (vs *VariableStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VariableStatement) String() string {
	var out bytes.Buffer
	out.WriteString("variable ")
	if vs.Name != nil {
		out.WriteString(vs.Name.String())
	}
	if vs.Type != nil {
		out.WriteString(" ")
		out.WriteString(vs.Type.String())
	}
	if vs.Default != nil {
		out.WriteString(" = ")
		out.WriteString(vs.Default.String())
	}
	return out.String()
}

// OutputStatement represents an output declaration
type OutputStatement struct {
	Token lexer.Token     // the 'output' token
	Name  string          // output name
	Value Expression      // output value
}

func (os *OutputStatement) statementNode()       {}
func (os *OutputStatement) TokenLiteral() string { return os.Token.Literal }
func (os *OutputStatement) String() string {
	var out bytes.Buffer
	out.WriteString("output \"")
	out.WriteString(os.Name)
	out.WriteString("\" ")
	if os.Value != nil {
		out.WriteString(os.Value.String())
	}
	return out.String()
}

// ModuleStatement represents a module declaration
type ModuleStatement struct {
	Token  lexer.Token    // the 'module' token
	Name   string         // module name
	Config *MapLiteral    // module configuration
}

func (ms *ModuleStatement) statementNode()       {}
func (ms *ModuleStatement) TokenLiteral() string { return ms.Token.Literal }
func (ms *ModuleStatement) String() string {
	var out bytes.Buffer
	out.WriteString("module \"")
	out.WriteString(ms.Name)
	out.WriteString("\" ")
	if ms.Config != nil {
		out.WriteString(ms.Config.String())
	}
	return out.String()
}

// AgentStatement represents an agent declaration
type AgentStatement struct {
	Token     lexer.Token    // the 'agent' token
	AgentType string         // agent type (e.g., "openai", "anthropic")
	Name      string         // agent name
	Config    *MapLiteral    // agent configuration
}

func (as *AgentStatement) statementNode()       {}
func (as *AgentStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AgentStatement) String() string {
	var out bytes.Buffer
	out.WriteString("agent ")
	out.WriteString(as.AgentType)
	out.WriteString(" \"")
	out.WriteString(as.Name)
	out.WriteString("\" ")
	if as.Config != nil {
		out.WriteString(as.Config.String())
	}
	return out.String()
}

// ScriptStatement represents a script declaration
type ScriptStatement struct {
	Token   lexer.Token    // the 'script' token
	Name    string         // script name
	Body    *BlockStatement // script body
}

func (ss *ScriptStatement) statementNode()       {}
func (ss *ScriptStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *ScriptStatement) String() string {
	var out bytes.Buffer
	out.WriteString("script \"")
	out.WriteString(ss.Name)
	out.WriteString("\" ")
	if ss.Body != nil {
		out.WriteString(ss.Body.String())
	}
	return out.String()
}

// ExpressionStatement represents a statement that consists of a single expression
type ExpressionStatement struct {
	Token      lexer.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// AssignmentStatement represents an assignment
type AssignmentStatement struct {
	Token lexer.Token // the identifier token
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentStatement) String() string {
	var out bytes.Buffer
	if as.Name != nil {
		out.WriteString(as.Name.String())
	}
	out.WriteString(" = ")
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	return out.String()
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Token lexer.Token // the 'return' token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString("return")
	if rs.Value != nil {
		out.WriteString(" ")
		out.WriteString(rs.Value.String())
	}
	return out.String()
}

// IfStatement represents an if/else statement
type IfStatement struct {
	Token       lexer.Token     // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // optional else block
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")
	out.WriteString(is.Consequence.String())
	if is.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(is.Alternative.String())
	}
	return out.String()
}

// ForStatement represents a for loop
type ForStatement struct {
	Token    lexer.Token     // the 'for' token
	Variable *Identifier     // loop variable
	Iterable Expression      // what to iterate over
	Body     *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	if fs.Variable != nil {
		out.WriteString(fs.Variable.String())
	}
	out.WriteString(" in ")
	if fs.Iterable != nil {
		out.WriteString(fs.Iterable.String())
	}
	out.WriteString(" ")
	if fs.Body != nil {
		out.WriteString(fs.Body.String())
	}
	return out.String()
}

// BlockStatement represents a block of statements
type BlockStatement struct {
	Token      lexer.Token // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
		out.WriteString(" ")
	}
	out.WriteString("}")
	return out.String()
}

// ============================================================================
// Expression Types
// ============================================================================

// Identifier represents an identifier expression
type Identifier struct {
	Token lexer.Token // the IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents an integer literal
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents a float literal
type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral represents a string literal
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// BooleanLiteral represents a boolean literal
type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string {
	if bl.Value {
		return "true"
	}
	return "false"
}

// ArrayLiteral represents an array literal
type ArrayLiteral struct {
	Token    lexer.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// MapLiteral represents a map/object literal
type MapLiteral struct {
	Token lexer.Token // the '{' token
	Pairs map[string]Expression
}

func (ml *MapLiteral) expressionNode()      {}
func (ml *MapLiteral) TokenLiteral() string { return ml.Token.Literal }
func (ml *MapLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for k, v := range ml.Pairs {
		pairs = append(pairs, k+" = "+v.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// BinaryExpression represents a binary operation
type BinaryExpression struct {
	Token    lexer.Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode()      {}
func (be *BinaryExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BinaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(be.Left.String())
	out.WriteString(" " + be.Operator + " ")
	out.WriteString(be.Right.String())
	out.WriteString(")")
	return out.String()
}

// UnaryExpression represents a unary operation
type UnaryExpression struct {
	Token    lexer.Token // the operator token
	Operator string
	Right    Expression
}

func (ue *UnaryExpression) expressionNode()      {}
func (ue *UnaryExpression) TokenLiteral() string { return ue.Token.Literal }
func (ue *UnaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ue.Operator)
	out.WriteString(ue.Right.String())
	out.WriteString(")")
	return out.String()
}

// CallExpression represents a function call
type CallExpression struct {
	Token     lexer.Token  // the '(' token
	Function  Expression   // Identifier or MemberExpression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// IndexExpression represents array/map indexing
type IndexExpression struct {
	Token lexer.Token // the '[' token
	Left  Expression  // the array or map
	Index Expression  // the index or key
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	return out.String()
}

// MemberExpression represents object property access
type MemberExpression struct {
	Token    lexer.Token // the '.' token
	Object   Expression  // the object
	Property string      // the property name
}

func (me *MemberExpression) expressionNode()      {}
func (me *MemberExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MemberExpression) String() string {
	var out bytes.Buffer
	out.WriteString(me.Object.String())
	out.WriteString(".")
	out.WriteString(me.Property)
	return out.String()
}

// TernaryExpression represents a ternary conditional expression
type TernaryExpression struct {
	Token       lexer.Token // the '?' token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (te *TernaryExpression) expressionNode()      {}
func (te *TernaryExpression) TokenLiteral() string { return te.Token.Literal }
func (te *TernaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(te.Condition.String())
	out.WriteString(" ? ")
	out.WriteString(te.Consequence.String())
	out.WriteString(" : ")
	out.WriteString(te.Alternative.String())
	out.WriteString(")")
	return out.String()
}