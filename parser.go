package baisl

import (
	"fmt"
	"slices"
)

type symbolTable map[string](bool)

type Parser struct {
	nextToken  *Token
	SourceFile *SourceFile
}

func (p *Parser) EatNextToken() *Token {
	nextToken := p.SourceFile.GetNextToken()
	p.nextToken = &nextToken
	return p.nextToken
}

func assertTokenType(token *Token, ttypes ...TokenType) {
	match := slices.Contains(ttypes, token.TType)

	if !match {
		if len(ttypes) == 1 {
			panic(fmt.Sprintf("Expected token type %s, got %s at %d:%d in %s", ttypes[0].String(), token.TType.String(), token.Location.Line, token.Location.Column, token.Location.Path))
		} else {
			panic(fmt.Sprintf("Expected token type in %v, got %s at %d:%d in %s", ttypes, token.TType.String(), token.Location.Line, token.Location.Column, token.Location.Path))
		}
	}
}

func assertNotTokenType(token *Token, ttype TokenType) {
	if token.TType == ttype {
		panic(fmt.Sprintf("Expected token type different from %s, got %s at %d:%d", ttype.String(), token.TType.String(), token.Location.Line, token.Location.Column))
	}
}

func (p *Parser) ParseExpr() (*Expr, error) {
	if p.nextToken.TType == TokenType_NUMBER {
		expr := Expr{
			Location: p.nextToken.Location,
			Type:     ExprType_INT,
			value:    p.nextToken.Value,
		}
		return &expr, nil
	}
	if p.nextToken.TType == TokenType_IDENTIFIER {
		expr := Expr{
			Location: p.nextToken.Location,
			Type:     ExprType_DECL_REF,
			value:    p.nextToken.Value,
		}
		return &expr, nil
	}
	return nil, fmt.Errorf("Unexpected token at %d:%d", p.nextToken.Location.Line, p.nextToken.Location.Column)
}

func (p *Parser) ParseReturnStmt() (Statement, error) {
	assertTokenType(p.nextToken, TokenType_KEYW_RETURN)
	var expr *Expr
	var err error
	assertTokenType(p.EatNextToken(), TokenType_NUMBER, TokenType_IDENTIFIER, TokenType_RBRACE)
	if p.nextToken.TType == TokenType_RBRACE {
		returnStmt := ReturnStmt{
			Stmt: Stmt{
				Location: p.nextToken.Location,
				Kind:     StmtType_RETURN,
			},
			Expr: nil,
		}

		return &returnStmt, nil
	}

	expr, err = p.ParseExpr()

	if err != nil {
		return nil, fmt.Errorf("Failed to parse expression: %v", err)
	}
	returnStmt := ReturnStmt{
		Stmt: Stmt{
			Location: p.nextToken.Location,
			Kind:     StmtType_RETURN,
		},
		Expr: expr,
	}

	for p.nextToken.TType != TokenType_RBRACE {
		if p.nextToken.TType == TokenType_RBRACE {
			break
		}

		assertNotTokenType(p.nextToken, TokenType_EOF)

		p.EatNextToken()
	}

	return &returnStmt, nil
}

func (p *Parser) ParseBlock() (*Block, error) {
	p.EatNextToken()
	assertTokenType(p.nextToken, TokenType_LBRACE)

	stmts := make([]*Statement, 0)
	for p.EatNextToken().TType != TokenType_RBRACE {
		assertTokenType(p.nextToken, TokenType_KEYW_RETURN, TokenType_RBRACE)
		if p.nextToken.TType == TokenType_KEYW_RETURN {
			returnStmt, err := p.ParseReturnStmt()
			if err != nil {
				return nil, err
			}

			stmts = append(stmts, &returnStmt)
			assertTokenType(p.nextToken, TokenType_RBRACE)
			break
		}
	}

	return &Block{
		Location: p.nextToken.Location,
		Stmts:    stmts,
	}, nil
}

func (p *Parser) ParseParameterList() ([]*VariableDecl, error) {
	assertTokenType(p.nextToken, TokenType_LPAREN)

	variables := make([]*VariableDecl, 0)
	lastToken := p.nextToken
	for p.EatNextToken().TType != TokenType_RPAREN {
		assertNotTokenType(p.nextToken, TokenType_EOF)

		if lastToken.TType == TokenType_LPAREN || lastToken.TType == TokenType_COMMA {
			assertTokenType(p.nextToken, TokenType_IDENTIFIER)
			id := p.nextToken.Value
			initialLocation := p.nextToken.Location
			assertTokenType(p.EatNextToken(), TokenType_COLON)
			assertTokenType(p.EatNextToken(), TokenType_KEYW_INT)

			decl := VariableDecl{
				Decl: Decl{
					id:       id,
					Location: initialLocation,
				},
				Type: Type_INT,
			}

			variables = append(variables, &decl)
		} else if lastToken.TType == TokenType_IDENTIFIER {
			assertTokenType(p.nextToken, TokenType_COMMA, TokenType_RPAREN)
		} else {
			panic(fmt.Sprintf("Unexpected token at %d:%d", p.nextToken.Location.Line, p.nextToken.Location.Column))
		}
		lastToken = p.nextToken
	}

	return variables, nil
}

func (p *Parser) ParseFunction() (*FunctionDecl, error) {
	assertTokenType(p.nextToken, TokenType_KEYW_FN)
	assertTokenType(p.EatNextToken(), TokenType_IDENTIFIER)
	fnName := p.nextToken.Value

	assertTokenType(p.EatNextToken(), TokenType_LPAREN, TokenType_COLON)
	if fnName == "main" {
		assertNotTokenType(p.nextToken, TokenType_LPAREN)
	}

	var parameters []*VariableDecl
	if p.nextToken.TType == TokenType_LPAREN {
		var err error
		parameters, err = p.ParseParameterList()
		if err != nil {
			return nil, fmt.Errorf("Failed to parse parameter list: %v", err)
		}

		assertTokenType(p.EatNextToken(), TokenType_COLON)
	}

	assertTokenType(p.EatNextToken(), TokenType_KEYW_INT, TokenType_KEYW_VOID)
	returnType := Type_INT
	if p.nextToken.TType == TokenType_KEYW_VOID {
		returnType = Type_VOID
	}

	block, err := p.ParseBlock()

	if err != nil {
		return nil, fmt.Errorf("Failed to parse block: %v", err)
	}

	return &FunctionDecl{
		Decl: Decl{
			id:       fnName,
			Location: p.nextToken.Location,
		},
		ReturnType: returnType,
		Body:       block,
		Params:     parameters,
	}, nil
}

func (p *Parser) Parse() ([]Declaration, error) {
	declarations := make([]Declaration, 0)

	var lastToken *Token
	next := p.EatNextToken()
	for next.TType != TokenType_EOF {

		if lastToken == nil && next.TType != TokenType_KEYW_FN {
			return nil, fmt.Errorf("Expected function declaration at %d:%d, found %v", next.Location.Line, next.Location.Column, next.TType)
		}

		fn, err := p.ParseFunction()
		if err != nil {
			return nil, fmt.Errorf("Failed to parse function: %v", err)
		}
		declarations = append(declarations, fn)

		lastToken = next
		next = p.EatNextToken()
	}

	hasMain := false
	for _, decl := range declarations {
		if decl.GetId() == "main" && decl.GetKind() == DeclType_FUNCTION {
			hasMain = true
			break
		}
	}

	if !hasMain {
		return nil, fmt.Errorf("No main function found in %s", p.SourceFile.path)
	}

	return declarations, nil
}
