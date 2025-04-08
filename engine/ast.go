package engine

type Node interface {
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

type Identifier struct {
	Token Token
	Value string
}

type NumberLiteral struct {
	Token Token
	Value float64
}

type StringLiteral struct {
	Token Token
	Value string
}

type FunctionLiteral struct {
	Token      Token
	Parameters []*Identifier
	Body       *BlockStatement
}

type CallExpression struct {
	Token     Token
	Function  Expression
	Arguments []Expression
}

type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

type LetStatement struct {
	Token Token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Token       Token
	ReturnValue Expression
}

type BlockStatement struct {
	Token      Token
	Statements []Statement
}

type IfExpression struct {
	Token       Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
