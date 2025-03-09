package baisl_test

import (
	"strings"
	"testing"

	"github.com/frodi-karlsson/baisl"
)

type semanticAnalyserTest struct {
	declarations []baisl.Declaration
}

type semanticAnalyserFailTest struct {
	declarations  []baisl.Declaration
	errorContains string
}

func getEmptyMainDeclarations() []baisl.Declaration {
	decls := make([]baisl.Declaration, 0)
	decls = append(decls, &baisl.FunctionDecl{
		Decl: baisl.Decl{
			Location: baisl.SourceLocation{
				Line:   1,
				Column: 1,
			},
			Id: "main",
		},
		ReturnType: baisl.Type_VOID,
		Params:     make([]*baisl.VariableDecl, 0),
		Body: &baisl.Block{
			Location: baisl.SourceLocation{
				Line:   1,
				Column: 1,
			},
			Stmts: []baisl.Statement{
				&baisl.ReturnStmt{
					Stmt: baisl.Stmt{
						Location: baisl.SourceLocation{
							Line:   1,
							Column: 1,
						},
						Kind: baisl.StmtType_RETURN,
					},
				},
			},
		},
	})

	return decls
}

func getReturnParamFuncDeclarations() []baisl.Declaration {
	decls := getEmptyMainDeclarations()
	decls = append(decls, &baisl.FunctionDecl{
		Decl: baisl.Decl{
			Location: baisl.SourceLocation{
				Line:   1,
				Column: 1,
			},
			Id: "returnParam",
		},
		ReturnType: baisl.Type_INT,
		Params: []*baisl.VariableDecl{
			{
				Decl: baisl.Decl{
					Location: baisl.SourceLocation{
						Line:   1,
						Column: 1,
					},
					Id: "a",
				},
				Type: baisl.Type_INT,
			},
		},
		Body: &baisl.Block{
			Location: baisl.SourceLocation{
				Line:   1,
				Column: 1,
			},
			Stmts: []baisl.Statement{
				&baisl.ReturnStmt{
					Stmt: baisl.Stmt{
						Location: baisl.SourceLocation{
							Line:   1,
							Column: 1,
						},
						Kind: baisl.StmtType_RETURN,
					},
					Expr: &baisl.Expr{
						Location: baisl.SourceLocation{
							Line:   1,
							Column: 1,
						},
						Type:  baisl.ExprType_DECL_REF,
						Value: "a",
					},
				},
			},
		},
	})
	return decls
}

func getReturnUndeclaredParamFuncDeclarations() []baisl.Declaration {
	decls := []baisl.Declaration{
		&baisl.FunctionDecl{
			Decl: baisl.Decl{
				Location: baisl.SourceLocation{
					Line:   1,
					Column: 1,
				},
				Id: "returnParam",
			},
			ReturnType: baisl.Type_INT,
			Params: []*baisl.VariableDecl{
				{
					Decl: baisl.Decl{
						Location: baisl.SourceLocation{
							Line:   1,
							Column: 1,
						},
						Id: "a",
					},
					Type: baisl.Type_INT,
				},
			},
			Body: &baisl.Block{
				Location: baisl.SourceLocation{
					Line:   1,
					Column: 1,
				},
				Stmts: []baisl.Statement{
					&baisl.ReturnStmt{
						Stmt: baisl.Stmt{
							Location: baisl.SourceLocation{
								Line:   1,
								Column: 1,
							},
							Kind: baisl.StmtType_RETURN,
						},
						Expr: &baisl.Expr{
							Location: baisl.SourceLocation{
								Line:   1,
								Column: 1,
							},
							Type:  baisl.ExprType_DECL_REF,
							Value: "b",
						},
					},
				},
			},
		},
	}
	return decls
}

var semanticAnalyserTests = []semanticAnalyserTest{
	{
		declarations: getEmptyMainDeclarations(),
	},
	{
		declarations: getReturnParamFuncDeclarations(),
	},
}

var semanticAnalyserFailTests = []semanticAnalyserFailTest{
	{
		declarations:  getReturnUndeclaredParamFuncDeclarations(),
		errorContains: "Undeclared variable b",
	},
}

func TestSemanticAnalyser(t *testing.T) {
	for _, test := range semanticAnalyserTests {

		analyser := baisl.SemanticAnalyser{}
		err := analyser.Analyse(test.declarations)
		if err != nil {
			t.Errorf("Error analysing: %s", err)
		}
	}

	for _, test := range semanticAnalyserFailTests {
		analyser := baisl.SemanticAnalyser{}

		err := analyser.Analyse(test.declarations)
		if err == nil {
			t.Errorf("Expected error, got none")
		}

		if !strings.Contains(err.Error(), test.errorContains) {
			t.Errorf("Expected error containing <%s>, got <%s>", test.errorContains, err)
		}
	}
}
