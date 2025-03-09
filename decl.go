package baisl

import (
	"strings"
)

type Decl struct {
	Location SourceLocation
	id       string
}

func (d *Decl) GetId() string {
	return d.id
}

func (d *Decl) GetLocation() *SourceLocation {
	return &d.Location
}

type Declaration interface {
	GetId() string
	GetLocation() *SourceLocation
	GetKind() DeclType
	String(level int) string
}

type TypeKind int

const (
	TypeType_INT TypeKind = iota
	TypeType_VOID
	TypeType_CUSTOM
)

type Type struct {
	Kind TypeKind
	Name string
}

type DeclType int

const (
	DeclType_FUNCTION DeclType = iota
	DeclType_VARIABLE
)

func (d DeclType) String() string {
	switch d {
	case DeclType_FUNCTION:
		return "Function"
	case DeclType_VARIABLE:
		return "Variable"
	default:
		return "Unknown"
	}
}

var Type_INT = Type{TypeType_INT, "int"}
var Type_VOID = Type{TypeType_VOID, "void"}

func (t Type) String() string {
	return t.Name
}

type StmtType int

const (
	StmtType_RETURN StmtType = iota
)

func (s StmtType) String() string {
	switch s {
	case StmtType_RETURN:
		return "Return"
	default:
		return "Unknown"
	}
}

type Stmt struct {
	Location SourceLocation
	Kind     StmtType
}

type ExprType int

const (
	ExprType_DECL_REF ExprType = iota
	ExprType_INT
)

type Expr struct {
	Location SourceLocation
	Type     ExprType
	value    string
}

type ReturnStmt struct {
	Stmt
	Expr *Expr
}

type Statement interface {
	GetLocation() *SourceLocation
	GetKind() StmtType
	String(level int) string
}

func (s *ReturnStmt) GetLocation() *SourceLocation {
	return &s.Stmt.Location
}

func (s *ReturnStmt) GetKind() StmtType {
	return s.Kind
}

func (s *ReturnStmt) String(level int) string {
	if s.Expr == nil {
		return strings.Repeat("  ", level) + "Return"
	}
	return strings.Repeat("  ", level) + "Return " + s.Expr.value
}

type Block struct {
	Location SourceLocation
	Stmts    []*Statement
}

func (b *Block) String(level int) string {
	stmtStr := ""
	for _, stmt := range b.Stmts {
		resolvedStmt := *stmt
		stmtStr += resolvedStmt.String(level+1) + "\n"
	}
	return strings.Repeat("  ", level) + "Block:\n" + stmtStr
}

func (b *Block) GetLocation() *SourceLocation {
	return &b.Location
}

type FunctionDecl struct {
	Decl
	ReturnType Type
	Body       *Block
	Params     []*VariableDecl
}

func (f *FunctionDecl) GetId() string {
	return f.id
}

func (f *FunctionDecl) GetLocation() *SourceLocation {
	return &f.Location
}

func (f *FunctionDecl) GetKind() DeclType {
	return DeclType_FUNCTION
}

func (f *FunctionDecl) String(level int) string {
	body := *f.Body
	bodyStr := body.String(level + 1)
	paramsStrs := make([]string, len(f.Params))
	for i, param := range f.Params {
		paramsStrs[i] = param.GetId() + ": " + param.Type.String()
	}
	return strings.Repeat("  ", level) + "Function " + f.id + "(" + strings.Join(paramsStrs, ", ") + "): " + f.ReturnType.String() + ":\n" + bodyStr
}

type VariableDecl struct {
	Decl
	Type Type
}

func (v *VariableDecl) GetId() string {
	return v.id
}

func (v *VariableDecl) GetLocation() *SourceLocation {
	return &v.Location
}

func (v *VariableDecl) GetKind() DeclType {
	return DeclType_VARIABLE
}

func (v *VariableDecl) String(level int) string {
	return strings.Repeat("  ", level) + "Variable " + v.id + " " + v.Type.String()
}
