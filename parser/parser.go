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
	token.LEFT_PARENTHESIS:      CALL,
}

type (
	infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer     *lexer.Lexer
	errors    []string
	curToken  token.Token
	peekToken token.Token

	infixParseFns map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

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
	p.registerInfix(token.LEFT_PARENTHESIS, p.parseCallExpression)

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
	case token.FUNCTION:
		return p.parseFunctionStatement(false)
	case token.PROCEDURE:
		return p.parseFunctionStatement(true)
	case token.RESULT:
		return p.parseResultStatement()
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

	stmt.Value = p.parseExpression(LOWEST, false)

	if p.peekTokenIs(token.EOL) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST, true)

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

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedure int, isStartOfLine bool) ast.Expression {
	// prefix
	var leftExp ast.Expression
	switch p.curToken.Type {
	case token.IDENT:
		leftExp = p.parseIdentifier()
		if isStartOfLine {
			leftIdent := leftExp.(*ast.Identifier)
			if exp := p.parseAssignExpression(leftIdent); exp != nil {
				return exp
			}
		}
	case token.INT:
		leftExp = p.parseIntegerLiteral()
	case token.BANG, token.MINUS:
		leftExp = p.parsePrefixExpression()
	case token.TRUE, token.FALSE:
		leftExp = p.parseBoolean()
	case token.LEFT_PARENTHESIS:
		leftExp = p.parseGroupedExpression()
	}

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

func (p *Parser) parseAssignExpression(left *ast.Identifier) *ast.AssignmentExpression {
	tok := p.curToken
	switch p.peekToken.Type {
	case token.EQUAL_OR_ASSIGN:
		p.nextToken()
		p.nextToken()

		exp := p.parseExpression(LOWEST, false)
		return &ast.AssignmentExpression{
			Token:      tok,
			Identifier: left,
			Value:      exp,
		}
	default:
		return nil
	}
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

	expression.Right = p.parseExpression(PREFIX, false)

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
	expression.Right = p.parseExpression(precedence, false)

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

	exp := p.parseExpression(LOWEST, false)

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
	stmt.Condition = p.parseExpression(LOWEST, false)

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
	stmt.Condition = p.parseExpression(LOWEST, false)

	// MEMO: IFBはTHENを省略可能
	if p.expectPeek(token.THEN) {
		p.nextToken()
	}

	if !p.curTokenIs(token.EOL) {
		return nil
	}

	p.nextToken()

	stmt.Consequence = p.parseBlockStatement()

	if p.curTokenIs(token.ENDIF) {
		return stmt
	}

	if p.curTokenIs(token.ELSE) {
		if !p.expectPeek(token.EOL) {
			return nil
		}
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
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}

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
		token.FEND,
	}

	for _, tt := range ts {
		if tt == t {
			return true
		}
	}
	return false
}

func (p *Parser) parseFunctionStatement(isProc bool) *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{
		Token:  p.curToken,
		IsProc: isProc,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.LEFT_PARENTHESIS) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	if p.expectPeek(token.EOL) {
		p.nextToken()
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RIGHT_PARENTHESIS) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RIGHT_PARENTHESIS) {
		p.nextToken()
		return args
	}

	p.nextToken()

	args = append(args, p.parseExpression(LOWEST, false))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST, false))
	}

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		return nil
	}
	return args
}

func (p *Parser) parseResultStatement() *ast.ResultStatement {
	stmt := &ast.ResultStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.EQUAL_OR_ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.ResultValue = p.parseExpression(LOWEST, false)

	if p.peekTokenIs(token.EOL) {
		p.nextToken()
	}
	return stmt
}
