package ast

import "github.com/sam8helloworld/uwscgo/token"

type Node interface {
	// デバッグ or テスト用
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type DimStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ds *DimStatement) statementNode() {}

func (ds *DimStatement) TokenLiteral() string {
	return ds.Token.Literal
}

type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
