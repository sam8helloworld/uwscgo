package parser

import (
	"fmt"
	"strconv"

	"github.com/sam8helloworld/uwscgo/ast"
	"github.com/sam8helloworld/uwscgo/lexer"
	"github.com/sam8helloworld/uwscgo/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // = または <>
	LESSGREATER // > または <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X または !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQUAL_OR_ASSIGN:       EQUALS,
	token.NOT_EQUAL:             EQUALS,
	token.LESS_THAN:             LESSGREATER,
	token.LESS_THAN_OR_EQUAL:    LESSGREATER,
	token.GREATER_THAN:          LESSGREATER,
	token.GREATER_THAN_OR_EQUAL: LESSGREATER,
	token.PLUS:                  SUM,
	token.MINUS:                 SUM,
	token.SLASH:                 PRODUCT,
	token.ASTERISK:              PRODUCT,
	token.MOD:                   PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer     *lexer.Lexer
	errors    []string
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LEFT_PARENTHESIS, p.parseGroupedExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.EQUAL_OR_ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN_OR_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN_OR_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)

	// 2つのトークンを読み込むことでcurTokenとpeekTokenがセットされる
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
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
	case token.DIM:
		return p.parseDimStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.IFB:
		return p.parseIfbStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDimStatement() *ast.DimStatement {
	stmt := &ast.DimStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.EQUAL_OR_ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.EOL) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.EOL) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedure int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOL) && precedure < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
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

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfStatement() ast.Statement {
	stmt := &ast.IfStatement{
		Token: p.curToken,
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.THEN) {
		return nil
	}

	p.nextToken()
	stmt.Consequence = p.parseStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		p.nextToken()

		stmt.Alternative = p.parseStatement()
	}

	return stmt
}

func (p *Parser) parseIfbStatement() ast.Statement {
	stmt := &ast.IfbStatement{
		Token: p.curToken,
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	// MEMO: IFBはTHENを省略可能
	if p.expectPeek(token.THEN) {
		p.nextToken()
	}

	stmt.Consequence = p.parseBlockStatement()
	if p.curTokenIs(token.ELSE) {
		p.nextToken()
		stmt.Alternative = p.parseBlockStatement()
	}

	if p.curTokenIs(token.ELSEIF) {
		stmt.Alternative = p.parseIfbStatement()
	}

	if !p.curTokenIs(token.ENDIF) {
		return nil
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.curToken,
	}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !blockEndTokenIs(p.curToken.Type) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func blockEndTokenIs(t token.TokenType) bool {
	ts := []token.TokenType{
		token.ELSE,
		token.ELSEIF,
		token.ENDIF,
	}

	for _, tt := range ts {
		if tt == t {
			return true
		}
	}
	return false
}
