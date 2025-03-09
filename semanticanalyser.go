package baisl

import (
	"fmt"
	"strconv"
)

type Scope struct {
	name         string
	parent       *Scope
	children     []*Scope
	declarations []Declaration
}

type SemanticAnalyser struct {
	currentScope         *Scope
	scopes               []*Scope
	resolvedDeclarations []ResolvedDeclaration
}

type ResolvedRefExpr struct {
	ExprType ExprType // Always ExprType_DECL_REF
	Value    *ResolvedDeclaration
	IsCall   bool
	Args     []ResolvedExpr
}

type ResolvedValueExpr struct {
	ExprType ExprType // Always ExprType_INT
	Value    int
}

type ResolvedExpr interface {
	GetExprType() ExprType
	GetType() Type
}

func (rr *ResolvedRefExpr) GetExprType() ExprType {
	return rr.ExprType
}

func (rr *ResolvedRefExpr) GetType() Type {
	value := *rr.Value
	switch value.(type) {
	case *ResolvedVariableDeclaration:
		return value.(*ResolvedVariableDeclaration).Type
	case *ResolvedFunctionDeclaration:
		returnStmt := value.(*ResolvedFunctionDeclaration).Body.Stmts[len(value.(*ResolvedFunctionDeclaration).Body.Stmts)-1]
		if returnStmt.Expr == nil {
			return Type_VOID
		}

		return returnStmt.Expr.GetType()
	}

	return Type_VOID
}

func (rv *ResolvedValueExpr) GetExprType() ExprType {
	return rv.ExprType
}

func (rv *ResolvedValueExpr) GetType() Type {
	return Type_INT
}

type ResolvedStatement struct {
	StmtType StmtType
	Expr     ResolvedExpr
}

type ResolvedBlock struct {
	Stmts []*ResolvedStatement
}

type ResolvedDeclaration interface {
	GetDeclType() DeclType
	GetId() string
}

type ResolvedFunctionDeclaration struct {
	Id       string
	DeclType DeclType
	Params   []ResolvedDeclaration
	Body     *ResolvedBlock
}

func (rfd *ResolvedFunctionDeclaration) GetDeclType() DeclType {
	return rfd.DeclType
}

func (rfd *ResolvedFunctionDeclaration) GetId() string {
	return rfd.Id
}

type ResolvedVariableDeclaration struct {
	Id       string
	DeclType DeclType
	Type     Type
	Value    ResolvedExpr
}

func (rvd *ResolvedVariableDeclaration) GetDeclType() DeclType {
	return rvd.DeclType
}

func (rvd *ResolvedVariableDeclaration) GetId() string {
	return rvd.Id
}

func (sa *SemanticAnalyser) EnterScope(name string) {
	newScope := &Scope{
		name:   name,
		parent: sa.currentScope,
	}
	if sa.currentScope != nil {
		sa.currentScope.children = append(sa.currentScope.children, newScope)
	}
	sa.currentScope = newScope
}

func (sa *SemanticAnalyser) ExitScope() {
	sa.currentScope = sa.currentScope.parent
}

func (sa *SemanticAnalyser) AddDeclaration(decl Declaration) error {
	for _, d := range sa.currentScope.declarations {
		if d.GetId() == decl.GetId() {
			return fmt.Errorf("Duplicate declaration of %s at %d:%d in %s", decl.GetId(), decl.GetLocation().Line, decl.GetLocation().Column, sa.currentScope.name)
		}
	}
	sa.currentScope.declarations = append(sa.currentScope.declarations, decl)
	return nil
}

func (sa *SemanticAnalyser) FindDeclaration(id string) Declaration {
	for scope := sa.currentScope; scope != nil; scope = scope.parent {
		for _, decl := range scope.declarations {
			if decl.GetId() == id {
				return decl
			}
		}
	}
	return nil
}

func (sa *SemanticAnalyser) AnalyseBlock(block *Block) error {
	for _, stmt := range block.Stmts {
		switch stmt.(type) {
		case *ReturnStmt:
			expr := stmt.(*ReturnStmt).Expr
			if expr != nil && expr.Type == ExprType_DECL_REF {
				found := sa.FindDeclaration(expr.Value)
				if found == nil {
					return fmt.Errorf("Undeclared variable %s at %d:%d in %s", expr.Value, expr.Location.Line, expr.Location.Column, sa.currentScope.name)
				}
			}
		}
	}
	return nil
}

func (sa *SemanticAnalyser) AnalyseFunctionSymbols(decl *FunctionDecl) error {
	sa.EnterScope(decl.GetId())
	for _, param := range decl.Params {
		err := sa.AddDeclaration(param)
		if err != nil {
			return fmt.Errorf("Error adding parameter %s: %s", param.GetId(), err)
		}
	}
	err := sa.AnalyseBlock(decl.Body)
	if err != nil {
		return fmt.Errorf("Error analysing block: %s", err)
	}
	sa.ExitScope()
	sa.AddDeclaration(decl)
	return nil
}

func (sa *SemanticAnalyser) AnalyseSymbols(declarations []Declaration) error {
	// Scopes will
	for _, decl := range declarations {
		switch decl.(type) {
		case *VariableDecl:
			err := sa.AddDeclaration(decl)
			if err != nil {
				return fmt.Errorf("Error adding variable %s: %s", decl.GetId(), err)
			}
		case *FunctionDecl:
			err := sa.AnalyseFunctionSymbols(decl.(*FunctionDecl))
			if err != nil {
				return fmt.Errorf("Error analysing function %s: %s", decl.GetId(), err)
			}
		}
	}
	return nil
}

func (sa *SemanticAnalyser) FindResolvedDeclaration(id string) ResolvedDeclaration {
	for _, decl := range sa.resolvedDeclarations {
		if decl.GetId() == id {
			return decl
		}
	}
	return nil
}

func (sa *SemanticAnalyser) ResolveExpr(expr *Expr) (ResolvedExpr, error) {
	switch expr.Type {
	case ExprType_DECL_REF:
		var resolvedArgs []ResolvedExpr
		for _, arg := range expr.Args {
			resolvedArg, err := sa.ResolveExpr(arg)
			if err != nil {
				return nil, fmt.Errorf("Error resolving argument: %s", err)
			}
			resolvedArgs = append(resolvedArgs, resolvedArg)
		}

		found := sa.FindResolvedDeclaration(expr.Value)
		if found == nil {
			return nil, fmt.Errorf("Undeclared variable %s at %d:%d in %s", expr.Value, expr.Location.Line, expr.Location.Column, sa.currentScope.name)
		}
		return &ResolvedRefExpr{
			ExprType: ExprType_DECL_REF,
			Value:    &found,
			IsCall:   expr.IsCall,
			Args:     resolvedArgs,
		}, nil
	case ExprType_INT:
		val, err := strconv.Atoi(expr.Value)
		if err != nil {
			return nil, fmt.Errorf("Error parsing integer %s at %d:%d in %s", expr.Value, expr.Location.Line, expr.Location.Column, sa.currentScope.name)
		}
		return &ResolvedValueExpr{
			ExprType: ExprType_INT,
			Value:    val,
		}, nil
	}
	return nil, fmt.Errorf("Unknown expression type %d at %d:%d in %s", expr.Type, expr.Location.Line, expr.Location.Column, sa.currentScope.name)
}

func (sa *SemanticAnalyser) ResolveStatement(stmt Statement) (*ResolvedStatement, error) {
	switch stmt.(type) {
	case *ReturnStmt:
		expr := stmt.(*ReturnStmt).Expr
		if expr == nil {
			return &ResolvedStatement{
				StmtType: StmtType_RETURN,
			}, nil
		}
		resolvedExpr, err := sa.ResolveExpr(expr)
		if err != nil {
			return nil, fmt.Errorf("Error resolving expression: %s", err)
		}
		return &ResolvedStatement{
			StmtType: StmtType_RETURN,
			Expr:     resolvedExpr,
		}, nil
	}
	return nil, fmt.Errorf("Unknown statement type %d at %d:%d in %s", stmt.GetKind(), stmt.GetLocation().Line, stmt.GetLocation().Column, sa.currentScope.name)
}

func (sa *SemanticAnalyser) ResolveBlock(block *Block) (*ResolvedBlock, error) {
	resolvedBlock := &ResolvedBlock{}
	for _, stmt := range block.Stmts {
		resolvedStmt, err := sa.ResolveStatement(stmt)
		if err != nil {
			return nil, fmt.Errorf("Error resolving statement %s: %s", stmt.GetKind(), err)
		}
		resolvedBlock.Stmts = append(resolvedBlock.Stmts, resolvedStmt)
	}
	return resolvedBlock, nil
}

func (sa *SemanticAnalyser) ResolveVariableDeclaration(decl *VariableDecl) (*ResolvedVariableDeclaration, error) {
	var resolvedExpr ResolvedExpr
	if decl.Value != nil {
		var err error
		resolvedExpr, err = sa.ResolveExpr(decl.Value)
		if err != nil {
			return nil, fmt.Errorf("Error resolving expression: %s", err)
		}
	}
	resolvedDeclaration := &ResolvedVariableDeclaration{
		Id:       decl.GetId(),
		DeclType: decl.GetKind(),
		Type:     decl.Type,
		Value:    resolvedExpr,
	}

	sa.resolvedDeclarations = append(sa.resolvedDeclarations, resolvedDeclaration)

	return resolvedDeclaration, nil
}

func (sa *SemanticAnalyser) ResolveFunctionDeclaration(decl *FunctionDecl) (*ResolvedFunctionDeclaration, error) {
	var resolvedParams []ResolvedDeclaration
	for _, param := range decl.Params {
		resolvedParam, err := sa.ResolveVariableDeclaration(param)
		if err != nil {
			return nil, fmt.Errorf("Error resolving parameter in %s: %s", decl.GetId(), err)
		}
		resolvedParams = append(resolvedParams, resolvedParam)
		sa.resolvedDeclarations = append(sa.resolvedDeclarations, resolvedParam)
	}

	resolvedBlock, err := sa.ResolveBlock(decl.Body)
	if err != nil {
		return nil, fmt.Errorf("Error resolving block in %s: %s", decl.GetId(), err)
	}

	returnStatement := resolvedBlock.Stmts[len(resolvedBlock.Stmts)-1]
	if returnStatement.StmtType != StmtType_RETURN {
		return nil, fmt.Errorf("Function %s does not return a value", decl.GetId())
	}

	if returnStatement.Expr == nil && decl.ReturnType != Type_VOID {
		return nil, fmt.Errorf("Function %s returns void but declared as %s", decl.GetId(), decl.ReturnType)
	} else if returnStatement.Expr != nil {
		returnType := returnStatement.Expr.GetType()
		if decl.ReturnType != returnType {
			return nil, fmt.Errorf("Function %s returns %s but declared as %s", decl.GetId(), returnType, decl.ReturnType)
		}
	}

	functionDeclaration := &ResolvedFunctionDeclaration{
		Id:       decl.GetId(),
		DeclType: decl.GetKind(),
		Params:   resolvedParams,
		Body:     resolvedBlock,
	}

	sa.resolvedDeclarations = append(sa.resolvedDeclarations, functionDeclaration)
	return functionDeclaration, nil
}

func (sa *SemanticAnalyser) ResolveSymbols(declarations []Declaration) error {
	for _, decl := range declarations {
		switch decl.(type) {
		case *VariableDecl:
			_, err := sa.ResolveVariableDeclaration(decl.(*VariableDecl))
			if err != nil {
				return fmt.Errorf("Error resolving variable %s: %s", decl.GetId(), err)
			}
		case *FunctionDecl:
			_, err := sa.ResolveFunctionDeclaration(decl.(*FunctionDecl))
			if err != nil {
				return fmt.Errorf("Error resolving function %s: %s", decl.GetId(), err)
			}
		}
	}
	return nil
}

func (sa *SemanticAnalyser) Analyse(declarations []Declaration) ([]ResolvedDeclaration, error) {
	sa.EnterScope("global")
	err := sa.AnalyseSymbols(declarations)
	if err != nil {
		return nil, err
	}

	err = sa.ResolveSymbols(declarations)
	if err != nil {
		return nil, err
	}

	hasMain := false
	for _, decl := range sa.resolvedDeclarations {
		if decl.GetId() == "main" && decl.GetDeclType() == DeclType_FUNCTION {
			hasMain = true
			break
		}
	}
	if !hasMain {
		return nil, fmt.Errorf("No main function found")
	}

	return sa.resolvedDeclarations, err
}
