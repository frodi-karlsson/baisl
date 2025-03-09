package baisl_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/frodi-karlsson/baisl"
)

type semanticAnalyserTest struct {
	declarations []baisl.Declaration
	expectedJson string
	name         string
}

type semanticAnalyserFailTest struct {
	declarations  []baisl.Declaration
	errorContains string
	name          string
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
	decls := make([]baisl.Declaration, 0)

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

	decls = append(decls, &baisl.FunctionDecl{
		Decl: baisl.Decl{
			Location: baisl.SourceLocation{
				Line:   1,
				Column: 1,
			},
			Id: "main",
		},
		ReturnType: baisl.Type_INT,
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
					Expr: &baisl.Expr{
						Location: baisl.SourceLocation{
							Line:   1,
							Column: 1,
						},
						Type:   baisl.ExprType_DECL_REF,
						Value:  "returnParam",
						IsCall: true,
						Args: []*baisl.Expr{
							{
								Location: baisl.SourceLocation{
									Line:   1,
									Column: 1,
								},
								Type:   baisl.ExprType_DECL_REF,
								Value:  "a",
								IsCall: false,
							},
						},
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

func getIncorrectReturnTypesFuncDeclarations() []baisl.Declaration {
	decls := []baisl.Declaration{
		&baisl.FunctionDecl{
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
						Expr: &baisl.Expr{
							Location: baisl.SourceLocation{
								Line:   1,
								Column: 1,
							},
							Type:   baisl.ExprType_INT,
							Value:  "1",
							IsCall: false,
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
		expectedJson: "[{\"Id\":\"main\",\"DeclType\":0,\"Params\":null,\"Body\":{\"Stmts\":[{\"StmtType\":0,\"Expr\":null}]}}]",
		name:         "Empty main",
	},
	{
		declarations: getReturnParamFuncDeclarations(),
		expectedJson: `[{"Id":"a","DeclType":1,"Type":{"Kind":0,"Name":"int"},"Value":null},{"Id":"a","DeclType":1,"Type":{"Kind":0,"Name":"int"},"Value":null},{"Id":"returnParam","DeclType":0,"Params":[{"Id":"a","DeclType":1,"Type":{"Kind":0,"Name":"int"},"Value":null}],"Body":{"Stmts":[{"StmtType":0,"Expr":{"ExprType":0,"Value":{"Id":"a","DeclType":1,"Type":{"Kind":0,"Name":"int"},"Value":null},"IsCall":false,"Args":null}}]}},{"Id":"main","DeclType":0,"Params":null,"Body":{"Stmts":[{"StmtType":0,"Expr":{"ExprType":0,"Value":{"Id":"returnParam","DeclType":0,"Params":[{"Id":"a","DeclType":1,"Type":{"Kind":0,"Name":"int"},"Value":null}],"Body":{"Stmts":[{"StmtType":0,"Expr":{"ExprType":0,"Value":{"Id":"a","DeclType":1,"Type":{"Kind":0,"Name":"int"},"Value":null},"IsCall":false,"Args":null}}]}},"IsCall":true,"Args":[{"ExprType":0,"Value":{"Id":"a","DeclType":1,"Type":{"Kind":0,"Name":"int"},"Value":null},"IsCall":false,"Args":null}]}}]}}]`,
		name:         "Return param",
	},
}

var semanticAnalyserFailTests = []semanticAnalyserFailTest{
	{
		declarations:  getReturnUndeclaredParamFuncDeclarations(),
		errorContains: "Undeclared variable b",
		name:          "Return undeclared param",
	},
	{
		declarations:  getIncorrectReturnTypesFuncDeclarations(),
		errorContains: "returns int but declared as void",
		name:          "Incorrect return type",
	},
}

func TestSemanticAnalyser(t *testing.T) {
	for _, test := range semanticAnalyserTests {

		analyser := baisl.SemanticAnalyser{}
		resolvedDeclarations, err := analyser.Analyse(test.declarations)
		if err != nil {
			t.Errorf("Error analysing %s: %s", test.name, err)
		}

		resolvedJson, err := json.Marshal(resolvedDeclarations)
		if err != nil {
			t.Errorf("Error marshalling resolved declarations: %s", err)
		}

		if string(resolvedJson) != test.expectedJson {
			t.Errorf("Failed test %s, expected %s, got %s", test.name, test.expectedJson, resolvedJson)
		}
	}

	for _, test := range semanticAnalyserFailTests {
		analyser := baisl.SemanticAnalyser{}

		_, err := analyser.Analyse(test.declarations)
		if err == nil {
			t.Errorf("Expected error, got none")
		}

		if !strings.Contains(err.Error(), test.errorContains) {
			t.Errorf("Expected error containing <%s>, got <%s>", test.errorContains, err)
		}
	}
}
