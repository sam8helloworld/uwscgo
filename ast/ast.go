package ast

import (
	"bytes"
	"strings"

	"github.com/sam8helloworld/uwscgo/token"
)

type Node interface {
	// デバッグ or テスト用
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

/////////////// Root
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

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

/////////////// Statement
type DimStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ds *DimStatement) statementNode() {}

func (ds *DimStatement) TokenLiteral() string {
	return ds.Token.Literal
}

func (ds *DimStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ds.TokenLiteral() + " ")
	out.WriteString(ds.Name.String())
	out.WriteString(" = ")

	if ds.Value != nil {
		out.WriteString(ds.Value.String())
	}

	return out.String()
}

type ExpressionStatement struct {
	Token token.Token // 式の最初のトークン
	Expression
}

func (es *ExpressionStatement) statementNode() {

}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IfStatement struct {
	Token       token.Token // 'if'トークン
	Condition   Expression
	Consequence Statement
	Alternative Statement
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) TokenLiteral() string {
	return is.Token.Literal
}
func (is *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("IF")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")
	out.WriteString("THEN")
	out.WriteString(" ")
	out.WriteString(is.Consequence.String())

	if is.Alternative != nil {
		out.WriteString("ELSE")
		out.WriteString(" ")
		out.WriteString(is.Alternative.String())
	}

	return out.String()
}

type IfbStatement struct {
	Token       token.Token // 'if'トークン
	Condition   Expression
	Consequence *BlockStatement
	Alternative Statement // *IfStatement(ELSEIF) or *BlockStatement(ELSE) or nil
}

func (is *IfbStatement) statementNode() {}
func (is *IfbStatement) TokenLiteral() string {
	return is.Token.Literal
}
func (is *IfbStatement) String() string {
	var out bytes.Buffer

	out.WriteString("IFB")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")
	out.WriteString("THEN")
	out.WriteString(" ")
	out.WriteString(is.Consequence.String())

	// TODO: ELSEIFを含めた出力
	if is.Alternative != nil {
		out.WriteString("ELSE")
		out.WriteString(" ")
		out.WriteString(is.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type FunctionStatement struct {
	Token      token.Token
	name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fs *FunctionStatement) statementNode() {}
func (fs *FunctionStatement) TokenLiteral() string {
	return fs.Token.Literal
}
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fs.name.Value)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(" ")
	out.WriteString(fs.Body.String())

	return out.String()
}

/////////////// Expression
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {

}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token // 前置トークン 、「!」など
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // 演算子トークン 、「+」など
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.Token.Literal
}

type AssignmentExpression struct {
	Token      token.Token
	Identifier *Identifier
	Value      Expression
}

func (as *AssignmentExpression) expressionNode() {}
func (as *AssignmentExpression) TokenLiteral() string {
	return as.Token.Literal
}
func (as *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(as.Token.Literal)
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString(as.Value.String())

	return out.String()
}
