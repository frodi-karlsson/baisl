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

func assertTokenType(token *Token, ttypes ...TokenType) error {
	match := slices.Contains(ttypes, token.TType)

	if !match {
		if len(ttypes) == 1 {
			return fmt.Errorf("Expected token type %s, got %s at %d:%d in %s", ttypes[0].String(), token.TType.String(), token.Location.Line, token.Location.Column, token.Location.Path)
		} else {
			return fmt.Errorf("Expected token type in %v, got %s at %d:%d in %s", ttypes, token.TType.String(), token.Location.Line, token.Location.Column, token.Location.Path)
		}
	}

	return nil
}

func assertNotTokenType(token *Token, ttype TokenType) error {
	if token.TType == ttype {
		return fmt.Errorf("Expected token type different from %s, got %s at %d:%d", ttype.String(), token.TType.String(), token.Location.Line, token.Location.Column)
	}

	return nil
}

func (p *Parser) ParseExpr() (*Expr, error) {
	if p.nextToken.TType == TokenType_NUMBER {
		expr := Expr{
			Location: p.nextToken.Location,
			Type:     ExprType_INT,
			Value:    p.nextToken.Value,
		}
		p.EatNextToken()
		return &expr, nil
	}
	if p.nextToken.TType == TokenType_IDENTIFIER {
		isCall := false
		location := p.nextToken.Location
		value := p.nextToken.Value
		args := make([]*Expr, 0)
		if p.EatNextToken().TType == TokenType_LPAREN {
			err := assertTokenType(p.EatNextToken(), TokenType_RPAREN, TokenType_NUMBER, TokenType_IDENTIFIER)
			if err != nil {
				return nil, err
			}

			lastToken := p.nextToken
			for p.nextToken.TType != TokenType_RPAREN {
				if p.nextToken.TType == TokenType_COMMA {
					assertTokenType(lastToken, TokenType_NUMBER, TokenType_IDENTIFIER)
					lastToken = p.nextToken
					p.EatNextToken()
				} else {
					expr, err := p.ParseExpr()
					if err != nil {
						return nil, fmt.Errorf("Failed to parse expression argument: %v", err)
					}
					args = append(args, expr)
					lastToken = p.nextToken
				}
			}

			isCall = true
		}

		expr := Expr{
			IsCall:   isCall,
			Location: location,
			Type:     ExprType_DECL_REF,
			Value:    value,
			Args:     args,
		}
		return &expr, nil
	}
	return nil, fmt.Errorf("Unexpected token %s at %d:%d", p.nextToken.TType, p.nextToken.Location.Line, p.nextToken.Location.Column)
}

func (p *Parser) ParseReturnStmt() (Statement, error) {
	err := assertTokenType(p.nextToken, TokenType_KEYW_RETURN)
	if err != nil {
		return nil, err
	}
	err = assertTokenType(p.EatNextToken(), TokenType_NUMBER, TokenType_IDENTIFIER, TokenType_RBRACE)
	if err != nil {
		return nil, err
	}

	var expr *Expr
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

		err = assertNotTokenType(p.nextToken, TokenType_EOF)
		if err != nil {
			return nil, err
		}

		p.EatNextToken()
	}

	return &returnStmt, nil
}

func (p *Parser) ParseBlock() (*Block, error) {
	p.EatNextToken()
	err := assertTokenType(p.nextToken, TokenType_LBRACE)
	if err != nil {
		return nil, err
	}

	stmts := make([]Statement, 0)
	for p.EatNextToken().TType != TokenType_RBRACE {
		err = assertTokenType(p.nextToken, TokenType_KEYW_RETURN, TokenType_RBRACE)
		if err != nil {
			return nil, err
		}

		if p.nextToken.TType == TokenType_KEYW_RETURN {
			returnStmt, err := p.ParseReturnStmt()
			if err != nil {
				return nil, err
			}

			stmts = append(stmts, returnStmt)
			err = assertTokenType(p.nextToken, TokenType_RBRACE)
			if err != nil {
				return nil, err
			}

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
		err := assertNotTokenType(p.nextToken, TokenType_EOF)
		if err != nil {
			return nil, err
		}

		if lastToken.TType == TokenType_LPAREN || lastToken.TType == TokenType_COMMA {
			err = assertTokenType(p.nextToken, TokenType_IDENTIFIER)
			if err != nil {
				return nil, err
			}

			id := p.nextToken.Value
			initialLocation := p.nextToken.Location
			err = assertTokenType(p.EatNextToken(), TokenType_COLON)
			if err != nil {
				return nil, err
			}

			err = assertTokenType(p.EatNextToken(), TokenType_KEYW_INT)
			if err != nil {
				return nil, err
			}

			decl := VariableDecl{
				Decl: Decl{
					Id:       id,
					Location: initialLocation,
				},
				Type: Type_INT,
			}

			variables = append(variables, &decl)
		} else if lastToken.TType == TokenType_IDENTIFIER {
			err = assertTokenType(p.nextToken, TokenType_COMMA, TokenType_RPAREN)
			if err != nil {
				return nil, err
			}
		} else {
			panic(fmt.Sprintf("Unexpected token at %d:%d", p.nextToken.Location.Line, p.nextToken.Location.Column))
		}
		lastToken = p.nextToken
	}

	return variables, nil
}

func (p *Parser) ParseFunction() (*FunctionDecl, error) {
	err := assertTokenType(p.nextToken, TokenType_KEYW_FN)
	if err != nil {
		return nil, err
	}
	err = assertTokenType(p.EatNextToken(), TokenType_IDENTIFIER)
	if err != nil {
		return nil, err
	}
	fnName := p.nextToken.Value

	err = assertTokenType(p.EatNextToken(), TokenType_LPAREN, TokenType_COLON)
	if err != nil {
		return nil, err
	}
	if fnName == "main" {
		err = assertNotTokenType(p.nextToken, TokenType_LPAREN)
		if err != nil {
			return nil, err
		}
	}

	var parameters []*VariableDecl
	if p.nextToken.TType == TokenType_LPAREN {
		var err error
		parameters, err = p.ParseParameterList()
		if err != nil {
			return nil, fmt.Errorf("Failed to parse parameter list: %v", err)
		}

		err = assertTokenType(p.EatNextToken(), TokenType_COLON)
		if err != nil {
			return nil, err
		}
	}

	err = assertTokenType(p.EatNextToken(), TokenType_KEYW_INT, TokenType_KEYW_VOID)
	if err != nil {
		return nil, err
	}
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
			Id:       fnName,
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

	return declarations, nil
}
