package baisl_test

import (
	"testing"

	"github.com/frodi-karlsson/baisl"
)

type testpair struct {
	path     string
	expected []baisl.Token
}

var tests = []testpair{
	{"raw/ret2.baisl", []baisl.Token{
		{
			TType:    baisl.TokenType_KEYW_FN,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_IDENTIFIER,
			HasValue: true,
			Value:    "main",
		},
		{
			TType:    baisl.TokenType_COLON,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_KEYW_VOID,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_LBRACE,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_KEYW_RETURN,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_RBRACE,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_KEYW_FN,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_IDENTIFIER,
			HasValue: true,
			Value:    "return2",
		},
		{
			TType:    baisl.TokenType_COLON,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_KEYW_INT,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_LBRACE,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_KEYW_RETURN,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_NUMBER,
			HasValue: true,
			Value:    "2",
		},
		{
			TType:    baisl.TokenType_RBRACE,
			HasValue: false,
		},
		{
			TType:    baisl.TokenType_EOF,
			HasValue: false,
		},
	},
	},
}

func TestGetNextToken(t *testing.T) {
	i := 0
	eof := false
	file, err := baisl.GetSourceFile("raw/ret2.baisl")

	if err != nil {
		t.Errorf("Error opening file")
	}

	for !eof {
		token := file.GetNextToken()
		if token.TType != tests[0].expected[i].TType {
			t.Errorf("Expected ttype %v, got ttype %v at i %d", tests[0].expected[i].TType.String(), token.TType.String(), i)
		}

		if token.Value != tests[0].expected[i].Value {
			t.Errorf("Expected value %v, got value %v at i %d", tests[0].expected[i].Value, token.Value, i)
		}

		if token.HasValue != tests[0].expected[i].HasValue {
			t.Errorf("Expected hasValue %v, got hasValue %v at i %d", tests[0].expected[i].HasValue, token.HasValue, i)
		}

		if token.TType == baisl.TokenType_EOF {
			eof = true
		}

		i++
	}
}
