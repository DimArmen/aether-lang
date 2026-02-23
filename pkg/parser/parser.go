package parser

import (
	"fmt"
	"strconv"

	"github.com/aether-lang/aether/pkg/ast"
	"github.com/aether-lang/aether/pkg/lexer"
)

// Parser parses Aether source code into an AST
type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	curToken  lexer.Token
	peekToken lexer.Token
}

// New creates a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// Errors returns parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("Line %d: %s", p.curToken.Line, msg))
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.addError(fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type))
	return false
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(lexer.TokenEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.TokenResource:
		return p.parseResourceStatement()
	case lexer.TokenVariable:
		return p.parseVariableStatement()
	case lexer.TokenOutput:
		return p.parseOutputStatement()
	case lexer.TokenModule:
		return p.parseModuleStatement()
	case lexer.TokenAgent:
		return p.parseAgentStatement()
	case lexer.TokenLet:
		return p.parseLetStatement()
	case lexer.TokenReturn:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseResourceStatement() *ast.ResourceStatement {
	stmt := &ast.ResourceStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TokenIdent) {
		return nil
	}

	// Parse resource type (e.g., "compute.instance")
	stmt.Type = p.curToken.Literal
	for p.peekTokenIs(lexer.TokenDot) {
		p.nextToken() // consume dot
		if !p.expectPeek(lexer.TokenIdent) {
			return nil
		}
		stmt.Type += "." + p.curToken.Literal
	}

	if !p.expectPeek(lexer.TokenString) {
		return nil
	}
	stmt.Name = p.curToken.Literal

	if !p.expectPeek(lexer.TokenLBrace) {
		return nil
	}

	stmt.Properties = p.parseBlockExpression()

	return stmt
}

func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	stmt := &ast.VariableStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TokenString) {
		return nil
	}
	stmt.Name = p.curToken.Literal

	if !p.expectPeek(lexer.TokenLBrace) {
		return nil
	}

	stmt.Properties = p.parseBlockExpression()

	return stmt
}

func (p *Parser) parseOutputStatement() *ast.OutputStatement {
	stmt := &ast.OutputStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TokenString) {
		return nil
	}
	stmt.Name = p.curToken.Literal

	if !p.expectPeek(lexer.TokenLBrace) {
		return nil
	}

	stmt.Properties = p.parseBlockExpression()

	return stmt
}

func (p *Parser) parseModuleStatement() *ast.ModuleStatement {
	stmt := &ast.ModuleStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TokenString) {
		return nil
	}
	stmt.Name = p.curToken.Literal

	if !p.expectPeek(lexer.TokenLBrace) {
		return nil
	}

	stmt.Properties = p.parseBlockExpression()

	return stmt
}

func (p *Parser) parseAgentStatement() *ast.AgentStatement {
	stmt := &ast.AgentStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TokenIdent) {
		return nil
	}
	stmt.Name = p.curToken.Literal

	if !p.expectPeek(lexer.TokenLBrace) {
		return nil
	}

	stmt.Properties = p.parseBlockExpression()

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TokenIdent) {
		return nil
	}
	stmt.Name = p.curToken.Literal

	if !p.expectPeek(lexer.TokenAssign) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	if !p.curTokenIs(lexer.TokenRBrace) && !p.curTokenIs(lexer.TokenEOF) {
		stmt.Value = p.parseExpression(LOWEST)
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseBlockExpression() *ast.BlockExpression {
	block := &ast.BlockExpression{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(lexer.TokenRBrace) && !p.curTokenIs(lexer.TokenEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// Expression parsing with precedence
const (
	_ int = iota
	LOWEST
	OR          // or
	AND         // and
	EQUALS      // == !=
	LESSGREATER // > < >= <=
	SUM         // + -
	PRODUCT     // * / %
	POWER       // **
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	MEMBER      // object.property
)

var precedences = map[lexer.TokenType]int{
	lexer.TokenOr:       OR,
	lexer.TokenAnd:      AND,
	lexer.TokenEq:       EQUALS,
	lexer.TokenNotEq:    EQUALS,
	lexer.TokenLt:       LESSGREATER,
	lexer.TokenGt:       LESSGREATER,
	lexer.TokenLtEq:     LESSGREATER,
	lexer.TokenGtEq:     LESSGREATER,
	lexer.TokenPlus:     SUM,
	lexer.TokenMinus:    SUM,
	lexer.TokenSlash:    PRODUCT,
	lexer.TokenStar:     PRODUCT,
	lexer.TokenPercent:  PRODUCT,
	lexer.TokenPower:    POWER,
	lexer.TokenLParen:   CALL,
	lexer.TokenLBracket: INDEX,
	lexer.TokenDot:      MEMBER,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Prefix parsing
	var leftExp ast.Expression

	switch p.curToken.Type {
	case lexer.TokenIdent:
		leftExp = p.parseIdentifier()
	case lexer.TokenNumber:
		leftExp = p.parseNumberLiteral()
	case lexer.TokenString:
		leftExp = p.parseStringLiteral()
	case lexer.TokenTrue, lexer.TokenFalse:
		leftExp = p.parseBooleanLiteral()
	case lexer.TokenLBracket:
		leftExp = p.parseListLiteral()
	case lexer.TokenLBrace:
		leftExp = p.parseMapLiteral()
	case lexer.TokenMinus, lexer.TokenNot:
		leftExp = p.parseUnaryExpression()
	case lexer.TokenIf:
		leftExp = p.parseIfExpression()
	default:
		p.addError(fmt.Sprintf("no prefix parse function for %s", p.curToken.Type))
		return nil
	}

	// Infix parsing
	for !p.peekTokenIs(lexer.TokenEOF) && precedence < p.peekPrecedence() {
		switch p.peekToken.Type {
		case lexer.TokenPlus, lexer.TokenMinus, lexer.TokenStar, lexer.TokenSlash,
			lexer.TokenPercent, lexer.TokenPower, lexer.TokenEq, lexer.TokenNotEq,
			lexer.TokenLt, lexer.TokenGt, lexer.TokenLtEq, lexer.TokenGtEq,
			lexer.TokenAnd, lexer.TokenOr:
			p.nextToken()
			leftExp = p.parseBinaryExpression(leftExp)
		case lexer.TokenDot:
			p.nextToken()
			leftExp = p.parseMemberExpression(leftExp)
		case lexer.TokenLBracket:
			p.nextToken()
			leftExp = p.parseIndexExpression(leftExp)
		case lexer.TokenLParen:
			p.nextToken()
			leftExp = p.parseCallExpression(leftExp)
		default:
			return leftExp
		}
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	lit := &ast.NumberLiteral{Token: p.curToken}
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as number", p.curToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(lexer.TokenTrue)}
}

func (p *Parser) parseListLiteral() ast.Expression {
	list := &ast.ListLiteral{Token: p.curToken}
	list.Elements = []ast.Expression{}

	p.nextToken()

	if p.curTokenIs(lexer.TokenRBracket) {
		return list
	}

	list.Elements = append(list.Elements, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.TokenComma) {
		p.nextToken() // consume comma
		p.nextToken()
		list.Elements = append(list.Elements, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.TokenRBracket) {
		return nil
	}

	return list
}

func (p *Parser) parseMapLiteral() ast.Expression {
	mapLit := &ast.MapLiteral{Token: p.curToken}
	mapLit.Pairs = make(map[string]ast.Expression)

	p.nextToken()

	if p.curTokenIs(lexer.TokenRBrace) {
		return mapLit
	}

	for !p.curTokenIs(lexer.TokenRBrace) && !p.curTokenIs(lexer.TokenEOF) {
		if !p.curTokenIs(lexer.TokenIdent) {
			p.addError("expected identifier as map key")
			return nil
		}
		key := p.curToken.Literal

		if !p.expectPeek(lexer.TokenAssign) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)
		mapLit.Pairs[key] = value

		if p.peekTokenIs(lexer.TokenComma) {
			p.nextToken()
		}
		p.nextToken()
	}

	return mapLit
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	expr := &ast.UnaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expr := &ast.BinaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	expr := &ast.MemberExpression{
		Token:  p.curToken,
		Object: left,
	}

	if !p.expectPeek(lexer.TokenIdent) {
		return nil
	}

	expr.Property = p.curToken.Literal

	return expr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.TokenRBracket) {
		return nil
	}

	return expr
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}

	expr.Arguments = []ast.Expression{}

	p.nextToken()

	if p.curTokenIs(lexer.TokenRParen) {
		return expr
	}

	expr.Arguments = append(expr.Arguments, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.TokenComma) {
		p.nextToken() // consume comma
		p.nextToken()
		expr.Arguments = append(expr.Arguments, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.TokenRParen) {
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.TokenLBrace) {
		return nil
	}

	expr.Consequence = p.parseBlockExpression()

	if p.peekTokenIs(lexer.TokenElse) {
		p.nextToken()

		if !p.expectPeek(lexer.TokenLBrace) {
			return nil
		}

		expr.Alternative = p.parseBlockExpression()
	}

	return expr
}
