package engine

import "strconv"

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{
		Statements: []Statement{},
	}

	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}

	return program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case LET:
		return p.parseLetStatement()
	case RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression()

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression()

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression()

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[TokenType]int{
	PLUS:     SUM,
	MINUS:    SUM,
	SLASH:    PRODUCT,
	ASTERISK: PRODUCT,
}

func (p *Parser) parseExpression() Expression {
	left := p.parsePrefixExpression()

	for !p.peekTokenIs(SEMICOLON) && p.curToken.Type != EOF {
		left = p.parseInfixExpression(left)
	}

	return left
}

func (p *Parser) parsePrefixExpression() Expression {
	switch p.curToken.Type {
	case NUMBER:
		lit := &NumberLiteral{Token: p.curToken}
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			return nil
		}
		lit.Value = value
		return lit
	case TRUE:
		return &BooleanLiteral{Token: p.curToken, Value: true}
	case FALSE:
		return &BooleanLiteral{Token: p.curToken, Value: false}
	case IDENT:
		return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case MINUS, BANG:
		expression := &PrefixExpression{
			Token:    p.curToken,
			Operator: p.curToken.Literal,
		}
		p.nextToken()
		expression.Right = p.parsePrefixExpression()
		return expression
	}
	return nil
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	if !p.isOperator(p.peekToken.Type) {
		return left
	}

	p.nextToken()
	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	currentPrecedence := p.curPrecedence()
	p.nextToken()

	for !p.peekTokenIs(SEMICOLON) && currentPrecedence < p.peekPrecedence() {
		p.nextToken()
		expression.Right = p.parsePrefixExpression()
	}

	return expression
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) isOperator(t TokenType) bool {
	return t == PLUS || t == MINUS || t == ASTERISK || t == SLASH
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
