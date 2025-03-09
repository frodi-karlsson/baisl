package baisl

import "fmt"

type Scope struct {
	name         string
	parent       *Scope
	children     []*Scope
	declarations []Declaration
}

type SemanticAnalyser struct {
	currentScope *Scope
	scopes       []*Scope
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

func (sa *SemanticAnalyser) Analyse(declarations []Declaration) error {
	sa.EnterScope("global")
	return sa.AnalyseSymbols(declarations)
}
