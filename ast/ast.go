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

type ConstStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (cs *ConstStatement) statementNode() {}

func (cs *ConstStatement) TokenLiteral() string {
	return cs.Token.Literal
}

func (cs *ConstStatement) String() string {
	var out bytes.Buffer

	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(cs.Name.String())
	out.WriteString(" = ")

	if cs.Value != nil {
		out.WriteString(cs.Value.String())
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
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
	IsProc     bool
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
	out.WriteString(fs.Name.Value)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(" ")
	out.WriteString(fs.Body.String())

	return out.String()
}

type ResultStatement struct {
	Token       token.Token // 'RESULT'トークン
	ResultValue Expression
}

func (rs *ResultStatement) statementNode() {}
func (rs *ResultStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ResultStatement) String() string {
	return rs.Token.Literal
}

type HashTableStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (hts *HashTableStatement) statementNode() {}
func (hts *HashTableStatement) TokenLiteral() string {
	return hts.Token.Literal
}
func (hts *HashTableStatement) String() string {
	var out bytes.Buffer

	out.WriteString(hts.TokenLiteral() + " ")
	out.WriteString(hts.Name.String())
	out.WriteString(" = ")

	if hts.Value != nil {
		out.WriteString(hts.Value.String())
	}

	return out.String()
}

type ForToStepStatement struct {
	Token   token.Token
	LoopVar *Identifier
	From    Expression
	To      Expression
	Step    Expression
	Block   *BlockStatement
}

func (ftss *ForToStepStatement) statementNode() {}
func (ftss *ForToStepStatement) TokenLiteral() string {
	return ftss.Token.Literal
}
func (ftss *ForToStepStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ftss.TokenLiteral() + " ")
	out.WriteString(ftss.LoopVar.String() + " ")
	out.WriteString(" = ")
	out.WriteString(ftss.From.String() + " ")
	out.WriteString("TO ")
	out.WriteString(ftss.To.String() + " ")
	out.WriteString("STEP ")
	out.WriteString(ftss.Step.String())
	out.WriteString("\n")
	out.WriteString(ftss.Block.String())
	out.WriteString("NEXT")

	return out.String()
}

type ForInStatement struct {
	Token      token.Token
	LoopVar    *Identifier
	Collection Expression
	Block      *BlockStatement
}

func (fis *ForInStatement) statementNode() {}
func (fis *ForInStatement) TokenLiteral() string {
	return fis.Token.Literal
}
func (fis *ForInStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fis.TokenLiteral() + " ")
	out.WriteString(fis.LoopVar.String() + " ")
	out.WriteString("IN ")
	out.WriteString(fis.Collection.String())
	out.WriteString("\n")
	out.WriteString(fis.Block.String())
	out.WriteString("NEXT")

	return out.String()
}

type ContinueStatement struct {
	Token token.Token
}

func (cs *ContinueStatement) statementNode() {}
func (cs *ContinueStatement) TokenLiteral() string {
	return cs.Token.Literal
}
func (cs *ContinueStatement) String() string {
	return cs.Token.Literal
}

type BreakStatement struct {
	Token token.Token
}

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BreakStatement) String() string {
	return bs.Token.Literal
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
	Token token.Token
	Left  Expression
	Value Expression
}

func (ae *AssignmentExpression) expressionNode() {}
func (ae *AssignmentExpression) TokenLiteral() string {
	return ae.Token.Literal
}
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Token.Literal)
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString(ae.Value.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // '('トークン
	Function  Expression  // Identifier
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}
func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

type ArrayLiteral struct {
	Token    token.Token
	Size     Expression // MEMO: nilかINTしか入らないから型で縛った方が良さそう
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}
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

type IndexExpression struct {
	Token  token.Token // '[' トークン
	Left   Expression
	Index  Expression
	Option Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

type EmptyArgument struct {
	Token token.Token
}

func (ea *EmptyArgument) expressionNode() {}
func (ea *EmptyArgument) TokenLiteral() string {
	return ea.Token.Literal
}
func (ea *EmptyArgument) String() string {
	return ea.Token.Literal
}

type Empty struct {
	Token token.Token
}

func (e *Empty) expressionNode() {}
func (e *Empty) TokenLiteral() string {
	return e.Token.Literal
}
func (e *Empty) String() string {
	return e.Token.Literal
}
